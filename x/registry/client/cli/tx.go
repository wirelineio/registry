//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/ghodss/yaml"
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

			filePath := args[0]
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			var payloadYaml registry.PayloadYaml
			err = yaml.Unmarshal(data, &payloadYaml)
			if err != nil {
				return err
			}

			signer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := registry.NewMsgSetResource(registry.PayloadYamlToPayload(payloadYaml), signer)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
