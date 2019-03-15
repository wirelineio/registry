//
// Copyright 2018 Wireline, Inc.
//

package tx

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/wirelineio/registry/x/multisig/msgs"
)

// GetCmdAbortMultiSig is the CLI command for sending a AbortMultiSig transaction.
func GetCmdAbortMultiSig(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "abort [id]",
		Short: "Abort multisig contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			id := args[0]

			aliceAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := msgs.NewMsgAbortMultiSig(id, aliceAddress)
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
