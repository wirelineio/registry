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

// GetCmdInitMultiSig is the CLI command for sending a InitMultiSig transaction.
func GetCmdInitMultiSig(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "init [id] [alice-amount] [bob-amount] [bob-address]",
		Short: "Init multisig contract.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			id := args[0]

			aliceAmount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			bobAmount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			aliceAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			bobAddress, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			msg := msgs.NewMsgInitMultiSig(id, aliceAmount, aliceAddress, bobAmount, bobAddress)
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
