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

package cmd

import (
	"flag"

	"github.com/spf13/pflag"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/njcx/libbeat_v8/cfgfile"
	cmd "github.com/njcx/libbeat_v8/cmd"
	"github.com/njcx/libbeat_v8/cmd/instance"
	"github.com/njcx/libbeat_v8/ecs"
	"github.com/njcx/libbeat_v8/processors"
	"github.com/njcx/libbeat_v8/publisher/processing"
	"github.com/njcx/packetbeat8_dpdk/beater"
	"github.com/njcx/packetbeat8_dpdk/include"
)

const (
	// Name of this beat.
	Name = "packetbeat"
)

// withECSVersion is a modifier that adds ecs.version to events.
var withECSVersion = processing.WithFields(mapstr.M{
	"ecs": mapstr.M{
		"version": ecs.Version,
	},
})

// RootCmd to handle beats cli
var RootCmd *cmd.BeatsRootCmd

// PacketbeatSettings contains the default settings for packetbeat
func PacketbeatSettings(globals processors.PluginConfig) instance.Settings {
	runFlags := pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.AddGoFlag(flag.CommandLine.Lookup("I"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("I")
	runFlags.AddGoFlag(flag.CommandLine.Lookup("t"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("t")
	runFlags.AddGoFlag(flag.CommandLine.Lookup("O"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("O")
	runFlags.AddGoFlag(flag.CommandLine.Lookup("l"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("l")
	runFlags.AddGoFlag(flag.CommandLine.Lookup("dump"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("dump")

	runFlags.AddGoFlag(flag.CommandLine.Lookup("dpdk_status"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("dpdk_status")
	runFlags.AddGoFlag(flag.CommandLine.Lookup("dpdk_port"))
	cfgfile.AddAllowedBackwardsCompatibleFlag("dpdk_port")

	return instance.Settings{
		RunFlags:       runFlags,
		Name:           Name,
		HasDashboards:  true,
		Processing:     processing.MakeDefaultSupport(true, globals, withECSVersion, processing.WithHost, processing.WithAgentMeta()),
		InputQueueSize: 400,
		Initialize:     []func(){include.InitializeModule},
	}
}

// Initialize initializes the entrypoint commands for packetbeat
func Initialize(settings instance.Settings) *cmd.BeatsRootCmd {
	rootCmd := cmd.GenRootCmdWithSettings(beater.New, settings)
	rootCmd.AddCommand(genDevicesCommand())
	return rootCmd
}

func init() {
	RootCmd = Initialize(PacketbeatSettings(nil))
}
