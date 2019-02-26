//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/wirelineio/wirechain/x/registry"
)

// GetCmdSetResource is the CLI command for sending a BirthOutput transaction.
func GetCmdSetResource(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set [resource file path]",
		Short: "Set resource.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			// TODO(ashwin): Read file and construct payload.
			// filePath := args[0]
			payload := registry.Payload{}

			signer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := registry.NewMsgSetResource(payload, signer)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
