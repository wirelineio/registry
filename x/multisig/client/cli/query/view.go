//
// Copyright 2018 Wireline, Inc.
//

package query

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetCmdView queries information about a contract.
func GetCmdView(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "view [id]",
		Short: "View contract by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/view/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Printf("Contract not found - %s \n", string(id))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
