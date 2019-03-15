//
// Copyright 2018 Wireline, Inc.
//

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	htlccmd "github.com/wirelineio/registry/x/htlc/client/cli"
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
	// Group htlc queries under a subcommand
	htlcQueryCmd := &cobra.Command{
		Use:   "htlc",
		Short: "Querying commands for the htlc module",
	}

	htlcQueryCmd.AddCommand(client.GetCommands()...)

	return htlcQueryCmd
}

// GetTxCmd returns the transaction commands for this module.
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	htlcTxCmd := &cobra.Command{
		Use:   "htlc",
		Short: "HTLC transactions subcommands",
	}

	htlcTxCmd.AddCommand(client.PostCommands(
		htlccmd.GetCmdAddHtlc(mc.cdc),
		htlccmd.GetCmdRedeemHtlc(mc.cdc),
		htlccmd.GetCmdFailHtlc(mc.cdc),
		htlccmd.GetCmdClearHtlc(mc.cdc),
	)...)

	return htlcTxCmd
}
