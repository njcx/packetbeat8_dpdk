// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mongodb

import (
	"fmt"
	"strings"
	"time"

	"github.com/njcx/libbeat_v8/common"
	conf "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/elastic/elastic-agent-libs/monitoring"

	"github.com/njcx/packetbeat8_dpdk/pb"
	"github.com/njcx/packetbeat8_dpdk/procs"
	"github.com/njcx/packetbeat8_dpdk/protos"
	"github.com/njcx/packetbeat8_dpdk/protos/tcp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var debugf = logp.MakeDebug("mongodb")

type mongodbPlugin struct {
	// config
	ports        []int
	sendRequest  bool
	sendResponse bool
	maxDocs      int
	maxDocLength int

	requests           *common.Cache
	responses          *common.Cache
	transactionTimeout time.Duration

	results protos.Reporter
	watcher *procs.ProcessesWatcher
}

type transactionKey struct {
	tcp common.HashableTCPTuple
	id  int32
}

var unmatchedRequests = monitoring.NewInt(nil, "mongodb.unmatched_requests")

func init() {
	protos.Register("mongodb", New)
}

func New(
	testMode bool,
	results protos.Reporter,
	watcher *procs.ProcessesWatcher,
	cfg *conf.C,
) (protos.Plugin, error) {
	p := &mongodbPlugin{}
	config := defaultConfig
	if !testMode {
		if err := cfg.Unpack(&config); err != nil {
			return nil, err
		}
	}

	if err := p.init(results, watcher, &config); err != nil {
		return nil, err
	}
	return p, nil
}

func (mongodb *mongodbPlugin) init(results protos.Reporter, watcher *procs.ProcessesWatcher, config *mongodbConfig) error {
	debugf("Init a MongoDB protocol parser")
	mongodb.setFromConfig(config)

	mongodb.requests = common.NewCache(
		mongodb.transactionTimeout,
		protos.DefaultTransactionHashSize)
	mongodb.requests.StartJanitor(mongodb.transactionTimeout)
	mongodb.responses = common.NewCache(
		mongodb.transactionTimeout,
		protos.DefaultTransactionHashSize)
	mongodb.responses.StartJanitor(mongodb.transactionTimeout)
	mongodb.results = results
	mongodb.watcher = watcher

	return nil
}

func (mongodb *mongodbPlugin) setFromConfig(config *mongodbConfig) {
	mongodb.ports = config.Ports
	mongodb.sendRequest = config.SendRequest
	mongodb.sendResponse = config.SendResponse
	mongodb.maxDocs = config.MaxDocs
	mongodb.maxDocLength = config.MaxDocLength
	mongodb.transactionTimeout = config.TransactionTimeout
}

func (mongodb *mongodbPlugin) GetPorts() []int {
	return mongodb.ports
}

func (mongodb *mongodbPlugin) ConnectionTimeout() time.Duration {
	return mongodb.transactionTimeout
}

func (mongodb *mongodbPlugin) Parse(
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
	private protos.ProtocolData,
) protos.ProtocolData {
	debugf("Parse method triggered")

	conn := ensureMongodbConnection(private)
	conn = mongodb.doParse(conn, pkt, tcptuple, dir)
	if conn == nil {
		return nil
	}
	return conn
}

func ensureMongodbConnection(private protos.ProtocolData) *mongodbConnectionData {
	if private == nil {
		return &mongodbConnectionData{}
	}

	priv, ok := private.(*mongodbConnectionData)
	if !ok {
		logp.Warn("mongodb connection data type error, create new one")
		return &mongodbConnectionData{}
	}
	if priv == nil {
		debugf("Unexpected: mongodb connection data not set, create new one")
		return &mongodbConnectionData{}
	}

	return priv
}

func (mongodb *mongodbPlugin) doParse(
	conn *mongodbConnectionData,
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
) *mongodbConnectionData {
	st := conn.streams[dir]
	if st == nil {
		st = newStream(pkt, tcptuple)
		conn.streams[dir] = st
		debugf("new stream: %p (dir=%v, len=%v)", st, dir, len(pkt.Payload))
	} else {
		// concatenate bytes
		st.data = append(st.data, pkt.Payload...)
		if len(st.data) > tcp.TCPMaxDataInStream {
			debugf("Stream data too large, dropping TCP stream")
			conn.streams[dir] = nil
			return conn
		}
	}

	for len(st.data) > 0 {
		if st.message == nil {
			st.message = &mongodbMessage{ts: pkt.Ts}
		}

		ok, complete := mongodbMessageParser(st)
		if !ok {
			// drop this tcp stream. Will retry parsing with the next
			// segment in it
			conn.streams[dir] = nil
			debugf("Ignore Mongodb message. Drop tcp stream. Try parsing with the next segment")
			return conn
		}

		if !complete {
			// wait for more data
			debugf("MongoDB wait for more data before parsing message")
			break
		}

		// all ok, go to next level and reset stream for new message
		debugf("MongoDB message complete")
		mongodb.handleMongodb(conn, st.message, tcptuple, dir)
		st.PrepareForNewMessage()
	}

	return conn
}

func newStream(pkt *protos.Packet, tcptuple *common.TCPTuple) *stream {
	s := &stream{
		tcptuple: tcptuple,
		data:     pkt.Payload,
		message:  &mongodbMessage{ts: pkt.Ts},
	}
	return s
}

func (mongodb *mongodbPlugin) handleMongodb(
	conn *mongodbConnectionData,
	m *mongodbMessage,
	tcptuple *common.TCPTuple,
	dir uint8,
) {
	m.tcpTuple = *tcptuple
	m.direction = dir
	m.cmdlineTuple = mongodb.watcher.FindProcessesTupleTCP(tcptuple.IPPort())

	if m.isResponse {
		debugf("MongoDB response message")
		mongodb.onResponse(conn, m)
	} else {
		debugf("MongoDB request message")
		mongodb.onRequest(conn, m)
	}
}

func (mongodb *mongodbPlugin) onRequest(conn *mongodbConnectionData, msg *mongodbMessage) {
	// publish request only transaction
	if !awaitsReply(msg) {
		mongodb.onTransComplete(msg, nil)
		return
	}

	id := msg.requestID
	key := transactionKey{tcp: msg.tcpTuple.Hashable(), id: id}

	// try to find matching response potentially inserted before
	if v := mongodb.responses.Delete(key); v != nil {
		resp := v.(*mongodbMessage)
		mongodb.onTransComplete(msg, resp)
		return
	}

	// insert into cache for correlation
	old := mongodb.requests.Put(key, msg)
	if old != nil {
		debugf("Two requests without a Response. Dropping old request")
		unmatchedRequests.Add(1)
	}
}

func (mongodb *mongodbPlugin) onResponse(conn *mongodbConnectionData, msg *mongodbMessage) {
	id := msg.responseTo
	key := transactionKey{tcp: msg.tcpTuple.Hashable(), id: id}

	// try to find matching request
	if v := mongodb.requests.Delete(key); v != nil {
		requ := v.(*mongodbMessage)
		mongodb.onTransComplete(requ, msg)
		return
	}

	// insert into cache for correlation
	mongodb.responses.Put(key, msg)
}

func (mongodb *mongodbPlugin) onTransComplete(requ, resp *mongodbMessage) {
	trans := newTransaction(requ, resp)
	debugf("Mongodb transaction completed: %s", trans.mongodb)
	mongodb.publishTransaction(trans)
}

func newTransaction(requ, resp *mongodbMessage) *transaction {
	trans := &transaction{}

	// fill request
	if requ != nil {
		trans.mongodb = mapstr.M{}
		trans.event = requ.event
		trans.method = requ.method

		trans.cmdline = requ.cmdlineTuple
		trans.ts = requ.ts
		trans.src, trans.dst = common.MakeEndpointPair(requ.tcpTuple.BaseTuple, requ.cmdlineTuple)
		if requ.direction == tcp.TCPDirectionReverse {
			trans.src, trans.dst = trans.dst, trans.src
		}
		trans.params = requ.params
		trans.resource = requ.resource
		trans.bytesIn = int(requ.messageLength)
		trans.documents = requ.documents
		trans.requestDocuments = requ.documents // preserving request documents that contains mongodb query for the new OP_MSG based protocol
	}

	// fill response
	if resp != nil {
		for k, v := range resp.event {
			trans.event[k] = v
		}

		trans.error = resp.error
		trans.documents = resp.documents

		trans.endTime = resp.ts
		trans.bytesOut = int(resp.messageLength)

	}

	return trans
}

func (mongodb *mongodbPlugin) GapInStream(tcptuple *common.TCPTuple, dir uint8,
	nbytes int, private protos.ProtocolData) (priv protos.ProtocolData, drop bool) {
	return private, true
}

func (mongodb *mongodbPlugin) ReceivedFin(tcptuple *common.TCPTuple, dir uint8,
	private protos.ProtocolData) protos.ProtocolData {
	return private
}

func copyMapWithoutKey(d map[string]interface{}, keys ...string) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range d {
		found := false
		for _, excludeKey := range keys {
			if k == excludeKey {
				found = true
				break
			}
		}
		if !found {
			res[k] = v
		}
	}
	return res
}

func reconstructQuery(t *transaction, full bool) (query string) {
	query = t.resource + "." + t.method + "("
	var doc interface{}

	if len(t.params) > 0 {
		if !full {
			// remove the actual data.
			// TODO: review if we need to add other commands here
			switch t.method {
			case "insert":
				doc = copyMapWithoutKey(t.params, "documents")
			case "update":
				doc = copyMapWithoutKey(t.params, "updates")
			case "findandmodify":
				doc = copyMapWithoutKey(t.params, "update")
			}
		} else {
			doc = t.params
		}
	} else if len(t.requestDocuments) > 0 { // This recovers the query document from OP_MSG
		if m, ok := t.requestDocuments[0].(primitive.M); ok {
			excludeKeys := []string{"lsid"}
			if !full {
				excludeKeys = append(excludeKeys, "documents")
			}
			doc = copyMapWithoutKey(m, excludeKeys...)
		}
	}

	queryString, err := doc2str(doc)
	if err != nil {
		debugf("Error marshaling query document: %v", err)
	} else {
		query += queryString
	}

	query += ")"
	skip, _ := t.event["numberToSkip"].(int)
	if skip > 0 {
		query += fmt.Sprintf(".skip(%d)", skip)
	}

	limit, _ := t.event["numberToReturn"].(int)
	if limit > 0 && limit < 0x7fffffff {
		query += fmt.Sprintf(".limit(%d)", limit)
	}
	return query
}

func (mongodb *mongodbPlugin) publishTransaction(t *transaction) {
	if mongodb.results == nil {
		debugf("Try to publish transaction with null results")
		return
	}

	evt, pbf := pb.NewBeatEvent(t.ts)
	pbf.SetSource(&t.src)
	pbf.AddIP(t.src.IP)
	pbf.SetDestination(&t.dst)
	pbf.AddIP(t.dst.IP)
	pbf.Source.Bytes = int64(t.bytesIn)
	pbf.Destination.Bytes = int64(t.bytesOut)
	pbf.Event.Dataset = "mongodb"
	pbf.Event.Start = t.ts
	pbf.Event.End = t.endTime
	pbf.Network.Transport = "tcp"
	pbf.Network.Protocol = pbf.Event.Dataset

	fields := evt.Fields
	fields["type"] = pbf.Event.Dataset
	if t.error == "" {
		fields["status"] = common.OK_STATUS
	} else {
		t.event["error"] = t.error
		fields["status"] = common.ERROR_STATUS
	}
	fields["mongodb"] = t.event
	fields["method"] = t.method
	fields["resource"] = t.resource
	fields["query"] = reconstructQuery(t, false)

	if mongodb.sendRequest {
		fields["request"] = reconstructQuery(t, true)
	}
	if mongodb.sendResponse {
		if len(t.documents) > 0 {
			// response field needs to be a string
			docs := make([]string, 0, len(t.documents))
			for i, doc := range t.documents {
				if mongodb.maxDocs > 0 && i >= mongodb.maxDocs {
					docs = append(docs, "[...]")
					break
				}
				str, err := doc2str(doc)
				if err != nil {
					logp.Warn("Failed to JSON marshal document from Mongo: %v (error: %v)", doc, err)
				} else {
					if mongodb.maxDocLength > 0 && len(str) > mongodb.maxDocLength {
						str = str[:mongodb.maxDocLength] + " ..."
					}
					docs = append(docs, str)
				}
			}
			fields["response"] = strings.Join(docs, "\n")
		}
	}

	mongodb.results(evt)
}
