// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code whose
// original notices appear below.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************
// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package runner

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"golang.org/x/term"

	"github.com/ava-labs/avalanchego/app/runner"
	"github.com/ava-labs/avalanchego/node"
	"github.com/ava-labs/avalanchego/vms/rpcchainvm/grpcutils"

	"github.com/chain4travel/camino-node/app"
	"github.com/chain4travel/camino-node/app/process"

	appplugin "github.com/chain4travel/camino-node/app/plugin"
)

// Run an AvalancheGo node.
// If specified in the config, serves a hashicorp plugin that can be consumed by
// the daemon (see avalanchego/main).
func Run(config runner.Config, nodeConfig node.Config) {
	nodeApp := process.NewApp(nodeConfig) // Create node wrapper
	if config.PluginMode {                // Serve as a plugin
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: appplugin.Handshake,
			Plugins: map[string]plugin.Plugin{
				appplugin.Name: appplugin.New(nodeApp),
			},
			GRPCServer: grpcutils.NewDefaultServer, // A non-nil value here enables gRPC serving for this plugin
			Logger: hclog.New(&hclog.LoggerOptions{
				Level: hclog.Error,
			}),
		})
		return
	}

	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println(process.Header)
	}

	exitCode := app.Run(nodeApp)
	os.Exit(exitCode)
}
