//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetCmdListAccOutput queries all account output birth records.
func GetCmdListAccOutput(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ls-account-outputs",
		Short: "List account output birth records.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/ls-account-outputs", queryRoute), nil)
			if err != nil {
				fmt.Println("{}")
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// GetCmdList queries all UTXOs.
func GetCmdList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List unspent transaction outputs (UTXO).",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/ls", queryRoute), nil)
			if err != nil {
				fmt.Println("{}")
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
