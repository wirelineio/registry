//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"encoding/hex"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/wirelineio/wirechain/x/utxo"
)

// GetCmdBirthOutput is the CLI command for sending a BirthOutput transaction.
func GetCmdBirthOutput(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "birth [amount]",
		Short: "Birth UTXO from account funds.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			coin, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := utxo.NewMsgBirthAccOutput(coin, account)
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

// GetCmdPayToAddress creates a UTXO style payment to a given address.
func GetCmdPayToAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pay [to] [amount] [change] [hash] [index]",
		Short: "Pay to address (UTXO style).",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			change, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			hash, err := hex.DecodeString(args[3])
			if err != nil {
				return err
			}

			index, err := strconv.ParseInt(args[4], 10, 32)
			if err != nil {
				// Hack/Workaround as passing -1 as arg on the cli confuses bash.
				index = -1
			}

			msg := utxo.NewMsgPayToAddress(cdc, hash, int32(index), amount, change, from, to)
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
