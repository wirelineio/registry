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

// GetCmdList queries all resources.
func GetCmdList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List resources.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/list", queryRoute), nil)
			if err != nil {
				fmt.Println("{}")
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// GetCmdGetResource queries a resource record.
func GetCmdGetResource(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [ID]",
		Short: "Get resource record.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			id := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, id), nil)
			if err != nil {
				fmt.Println("{}")
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// GetCmdGraph generates a dot graph.
func GetCmdGraph(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "graph",
		Short: "Generate dot graph.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/graph", queryRoute), nil)
			if err != nil {
				fmt.Println("{}")
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
