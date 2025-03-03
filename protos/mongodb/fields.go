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

package mongodb

import (
	"github.com/njcx/libbeat_v8/asset"
)

func init() {
	if err := asset.SetFields("packetbeat", "mongodb", asset.ModuleFieldsPri, AssetMongodb); err != nil {
		panic(err)
	}
}

// AssetMongodb returns asset data.
// This is the base64 encoded zlib format compressed contents of protos/mongodb.
func AssetMongodb() string {
	return "eJysVsFuGzcQvesrBoEPCSAJ6FWHAk1bAyma2LBdFD3V3N1ZkQjJ2XBIy/r7gkPuareSARWpTxY58+bN4+NwN/AVjztw5PfUNSuAaKLFHbz7nFd+ad6tADrkNpghGvI7+HEFAFB2P254wNb0pgV8QR+hN2g73sKTRsb6C5wJgQK0lhjtUdKjnnZ7CvKzIsLBBIQhUKSWrCCBNnuNASy+oK1pgvI+5+KrcoPFNTx/SxiOz6B8B88BmVJo8fkDqGGwR4g0VRCqDIrhgNZuV1AxdwK6Aa8czhURwscBd7APlIa6Ms+Z52Fudlq9KN7496lfdB7wW0KOoBVDQE42YgfGg/IFdA1RGy6FoSUflfGcERagEgoOmdUeIWBMwWMHzVFqMYYXDNsz1n2y9meyFttM9ItyeF0L+XhyLrRTsiBur4oCI/xzM62K6JXsUJGlU1E1irFEHkzUNXYBsV4USmz8HhR0FOH99sNkrkWBJbfbuYfG+Kl2TySG+lflRslpvN1VT7Rt1AWhfXINhid6/GqGGY/iL0t+f53ujxiLdAUva9ZRm5w4OxKQMxE2wFGFmCXpA7l66wLHKTbbK68Wu+XA3DljhM2i3EGjr17KQaeU8bDk7r3d7YOkfke/T1JR7gd24NSrccm91Xxz8v1blB7q/ndSukRgknSwFxQZT+Q2kLuuzp8aA46obQosHjVcKmSzjZjn1eRYrivzE/z2ePfl5IyoVcwlAnLR9fyURxVkGQ5GroKMJSCPkMcQBQS0KMqsQVmbtTpo02pwiSM4FVst105NpRf45TiNb23qyjicmY8xLtncE7Np7KnmmAk3QnINNxQ6DE3+Txsf13CDr4NVxq/lnt+wVwNriudaFkvdyth/xHzhr53zF6W1xpkqa30Jp97qyJ4sda74JS6nJ+Gy+KhafVK/Tl6ZV9SDqs+KUGNNyXbza7ScsuNAND7iHgO8KJsQfrhg9v8o08dzmeoHBs4MKHYp0OMwOuUQpKFTUQTo0GLEc1ol4n8jVQsWpw4YegpuPnmgvjLG5536xvkKcswtFAReAyMK5B8F8m7AIPE8ESjp0zgfvx0+K5+UPW+1zItP3fUjrU4Y06GPucdwMmQ16N393w+/3v/+V/44M1zucVN4FyeIRgf5iqmZE9/xYd2u/gkAAP//aSQ/9Q=="
}
