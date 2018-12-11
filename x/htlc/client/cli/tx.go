package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/wirelineio/cosmos-htlc/x/htlc"
)

// GetCmdAddHtlc is the CLI command for sending a AddHtlc transaction
func GetCmdAddHtlc(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add [amount] [hash] [locktime] [redeem-address]",
		Short: "Add HTLC.",
		Args:  cobra.ExactArgs(4),
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

			hash := args[1]

			locktime, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return err
			}

			redeemAccount, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			timeoutAccount, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := htlc.NewMsgAddHtlc(coin, hash, locktime, redeemAccount, timeoutAccount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
