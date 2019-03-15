//
// Copyright 2018 Wireline, Inc.
//

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	multisigqry "github.com/wirelineio/registry/x/multisig/client/cli/query"
	multisigcmd "github.com/wirelineio/registry/x/multisig/client/cli/tx"
)

// ModuleClient exports all client functionality from this module.
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

// NewModuleClient is the constructor for the module client.
func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module.
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group multisig queries under a subcommand
	multisigQueryCmd := &cobra.Command{
		Use:   "multisig",
		Short: "Querying commands for the multisig module",
	}

	multisigQueryCmd.AddCommand(client.GetCommands(
		multisigqry.GetCmdView(mc.storeKey, mc.cdc),
	)...)

	return multisigQueryCmd
}

// GetTxCmd returns the transaction commands for this module.
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	multisigTxCmd := &cobra.Command{
		Use:   "multisig",
		Short: "Multisig transactions subcommands",
	}

	multisigTxCmd.AddCommand(client.PostCommands(
		multisigcmd.GetCmdInitMultiSig(mc.cdc),
		multisigcmd.GetCmdAbortMultiSig(mc.cdc),
		multisigcmd.GetCmdJoinMultiSig(mc.cdc),
		multisigcmd.GetCmdSpendMultiSig(mc.cdc),
	)...)

	return multisigTxCmd
}
