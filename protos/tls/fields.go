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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package tls

import (
	"github.com/njcx/libbeat_v8/asset"
)

func init() {
	if err := asset.SetFields("packetbeat", "tls", asset.ModuleFieldsPri, AssetTls); err != nil {
		panic(err)
	}
}

// AssetTls returns asset data.
// This is the base64 encoded zlib format compressed contents of protos/tls.
func AssetTls() string {
	return "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"
}
