//
// Copyright 2019 Wireline, Inc.
//

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	regcmd "github.com/wirelineio/wirechain/x/registry/client/cli"
)

// ModuleClient exports all client functionality from this module.
type ModuleClient struct {
	cdc *amino.Codec
}

// NewModuleClient is the constructor for the module client.
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetQueryCmd returns the cli query commands for this module.
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group utxo queries under a subcommand
	regQueryCmd := &cobra.Command{
		Use:   "registry",
		Short: "Querying commands for the registry module",
	}

	regQueryCmd.AddCommand(client.GetCommands(
		regcmd.GetCmdGetResource("registry", mc.cdc),
		regcmd.GetCmdList("registry", mc.cdc),
		regcmd.GetCmdGraph("registry", mc.cdc),
		regcmd.GetCmdTest("registry", mc.cdc),
	)...)

	return regQueryCmd
}

// GetTxCmd returns the transaction commands for this module.
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	regTxCmd := &cobra.Command{
		Use:   "registry",
		Short: "registry transactions subcommands",
	}

	regTxCmd.AddCommand(client.PostCommands(
		regcmd.GetCmdSetResource(mc.cdc),
	)...)

	return regTxCmd
}
