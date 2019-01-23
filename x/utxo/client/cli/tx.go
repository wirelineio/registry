//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/wirelineio/wirechain/x/utxo"
	utxoutils "github.com/wirelineio/wirechain/x/utxo/utils"
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
	cmd := &cobra.Command{
		Use:   "pay [from] [to] [amount] [change] [hash] [index] [sig]",
		Short: "Pay to address (UTXO style).",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			from, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			change, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			hash, err := hex.DecodeString(args[4])
			if err != nil {
				return err
			}

			index, err := strconv.ParseInt(args[5], 10, 32)
			if err != nil {
				// Hack/Workaround as passing -1 as arg on the cli confuses bash.
				index = -1
			}

			sig, err := hex.DecodeString(args[6])
			if err != nil {
				return err
			}

			tx := utxo.NewTxPayToAddress(cdc, sig, hash, int32(index), amount, change, from, to)

			cliCtx.PrintResponse = true

			signOnly := viper.GetBool("sign-only")

			if signOnly {
				name := viper.GetString("from")

				sigBytes, err := utxo.GetTxSignature(cdc, tx, name)
				if err != nil {
					return err
				}

				fmt.Println(utxoutils.BytesToHex(sigBytes))

				return nil
			}

			signer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			msg := utxo.NewMsgTx(tx, signer)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Bool("sign-only", false, "Only sign the transaction payload.")

	return cmd
}
