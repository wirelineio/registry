//
// Copyright 2019 Wireline, Inc.
//

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	utxocmd "github.com/wirelineio/wirechain/x/utxo/client/cli"
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
	utxoQueryCmd := &cobra.Command{
		Use:   "utxo",
		Short: "Querying commands for the utxo module",
	}

	utxoQueryCmd.AddCommand(client.GetCommands(
		utxocmd.GetCmdListAccOutput("utxo", mc.cdc),
		utxocmd.GetCmdList("utxo", mc.cdc),
	)...)

	return utxoQueryCmd
}

// GetTxCmd returns the transaction commands for this module.
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	utxoTxCmd := &cobra.Command{
		Use:   "utxo",
		Short: "utxo transactions subcommands",
	}

	utxoTxCmd.AddCommand(client.PostCommands(
		utxocmd.GetCmdBirthOutput(mc.cdc),
	)...)

	return utxoTxCmd
}
