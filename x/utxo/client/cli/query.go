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

// GetCmdListAccountUtxo queries information about a contract.
func GetCmdListAccountUtxo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ls-account-utxo",
		Short: "List account UTXO birth records.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/ls-account-utxo", queryRoute), nil)
			if err != nil {
				fmt.Println("Not found")
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
