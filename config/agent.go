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

package config

import (
	"fmt"
	"runtime"

	"github.com/njcx/packetbeat8_dpdk/procs"
	conf "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/elastic/go-ucfg"
)

type datastream struct {
	Namespace string `config:"namespace"`
	Dataset   string `config:"dataset"`
	Type      string `config:"type"`
}

type agentInput struct {
	Type       string                   `config:"type"`
	Datastream datastream               `config:"data_stream"`
	Processors []mapstr.M               `config:"processors"`
	Streams    []map[string]interface{} `config:"streams"`
}

func defaultDevice() string {
	if runtime.GOOS == "linux" {
		return "any"
	}
	return "default_route"
}

func (i agentInput) addProcessorsAndIndex(cfg *conf.C) (*conf.C, error) {
	namespace := i.Datastream.Namespace
	if namespace == "" {
		namespace = "default"
	}
	datastreamConfig := struct {
		Datastream datastream `config:"data_stream"`
	}{}
	if err := cfg.Unpack(&datastreamConfig); err != nil {
		return nil, err
	}
	mergeConfig, err := conf.NewConfigFrom(mapstr.M{
		"index": datastreamConfig.Datastream.Type + "-" + datastreamConfig.Datastream.Dataset + "-" + namespace,
		"processors": append([]mapstr.M{
			{
				"add_fields": mapstr.M{
					"target": "data_stream",
					"fields": mapstr.M{
						"type":      datastreamConfig.Datastream.Type,
						"dataset":   datastreamConfig.Datastream.Dataset,
						"namespace": namespace,
					},
				},
			},
			{
				"add_fields": mapstr.M{
					"target": "event",
					"fields": mapstr.M{
						"dataset": datastreamConfig.Datastream.Dataset,
					},
				},
			},
		}, i.Processors...),
	})
	if err != nil {
		return nil, err
	}
	if err := cfg.MergeWithOpts(mergeConfig, ucfg.FieldAppendValues("processors")); err != nil {
		return nil, err
	}
	return cfg, nil
}

func mergeProcsConfig(one, two procs.ProcsConfig) procs.ProcsConfig {
	maxProcReadFreq := one.MaxProcReadFreq
	if two.MaxProcReadFreq > maxProcReadFreq {
		maxProcReadFreq = two.MaxProcReadFreq
	}

	refreshPidsFreq := one.RefreshPidsFreq
	if two.RefreshPidsFreq < refreshPidsFreq {
		refreshPidsFreq = two.RefreshPidsFreq
	}

	return procs.ProcsConfig{
		Enabled:         true,
		MaxProcReadFreq: maxProcReadFreq,
		RefreshPidsFreq: refreshPidsFreq,
		Monitored:       append(one.Monitored, two.Monitored...),
	}
}

// NewAgentConfig allows the packetbeat configuration to understand
// agent semantics
func NewAgentConfig(cfg *conf.C) (Config, error) {
	logp.Debug("agent", "Normalizing agent configuration")
	var (
		input  agentInput
		config Config
	)
	if err := cfg.Unpack(&input); err != nil {
		return config, err
	}

	logp.Debug("agent", "Found %d inputs", len(input.Streams))
	for _, stream := range input.Streams {
		if interfaceOverride, ok := stream["interface"]; ok {
			cfg, err := conf.NewConfigFrom(interfaceOverride)
			if err != nil {
				return config, err
			}
			var iface InterfaceConfig
			if err := cfg.Unpack(&iface); err != nil {
				return config, err
			}
			config.Interfaces = append(config.Interfaces, iface)
		}

		if procsOverride, ok := stream["procs"]; ok {
			cfg, err := conf.NewConfigFrom(procsOverride)
			if err != nil {
				return config, err
			}
			var newProcsConfig procs.ProcsConfig
			if err := cfg.Unpack(&newProcsConfig); err != nil {
				return config, err
			}
			config.Procs = mergeProcsConfig(config.Procs, newProcsConfig)
		}

		if rawStreamType, ok := stream["type"]; ok {
			streamType, ok := rawStreamType.(string)
			if !ok {
				return config, fmt.Errorf("invalid input type of: '%T'", rawStreamType)
			}
			logp.Debug("agent", "Found agent configuration for %v", streamType)
			cfg, err := conf.NewConfigFrom(stream)
			if err != nil {
				return config, err
			}
			cfg, err = input.addProcessorsAndIndex(cfg)
			if err != nil {
				return config, err
			}
			switch streamType {
			case "flow":
				if err := cfg.Unpack(&config.Flows); err != nil {
					return config, err
				}
			default:
				config.ProtocolsList = append(config.ProtocolsList, cfg)
			}
		}
	}
	if len(config.Interfaces) == 0 {
		config.Interfaces = []InterfaceConfig{
			// TODO: Make this configurable rather than just using the default device.
			{Device: defaultDevice()},
		}
	}
	return config, nil
}
