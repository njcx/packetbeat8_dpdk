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

package mysql

import (
	"github.com/njcx/libbeat_v8/asset"
)

func init() {
	if err := asset.SetFields("packetbeat", "mysql", asset.ModuleFieldsPri, AssetMysql); err != nil {
		panic(err)
	}
}

// AssetMysql returns asset data.
// This is the base64 encoded zlib format compressed contents of protos/mysql.
func AssetMysql() string {
	return "eJy8k8GO0zAQhu95il974UL3AXLggnpYaUGC7j1y7XFrEdupZ0yVt0d20tJCCq2Q8CWRx/8/nz0zK3yjsYUf+dA3gDjpqcXTp3Hz5fWpAQyxTm4QF0OLDw0A1NiKB9LOOg36TkFgHfWGnxvMf209ukJQnn7alyXjQC12KeZh3rlUXKqUtaSFTJfikc/Rk0Mfw+5icwH0tF4sZE8TN3T0XgUDx+CsNTHb3L+H7B1PINAxiHKBq+jEcGUYst9SQrQoZOVbjvaKBSxKyFOQ5+a3+7jAlKRz5iHsl8+b9dc3HDKl8V5qZ2aoa2w69uNMQaawL0CG7LupIA9Rbtav649/pywREkisoukdr+yindsBiSSnQOYG4y898V8Ja9X/wFdz3Af3tqdiN03IDKeKuTKwKfpKIUkFVrro35XYIRMvNRilFFOno6F/mJZCVI1QjM7XxHacJuhmXk/MakePZnHBRszaxWw/AgAA//+Mhliv"
}
