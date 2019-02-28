//
// Copyright 2019 Wireline, Inc.
//

package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wirelineio/wirechain/x/registry"
)

// GetCmdSetResource is the CLI command for sending a BirthOutput transaction.
func GetCmdSetResource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [resource file path]",
		Short: "Set resource.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc).WithChainID(registry.WirelineChainID)

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

			cliCtx.PrintResponse = true

			signOnly := viper.GetBool("sign-only")

			if signOnly {
				name := viper.GetString("from")

				sigBytes, pubKey, err := registry.GetResourceSignature(payloadYaml.Resource, name)
				if err != nil {
					return err
				}

				fmt.Println("Address:")
				fmt.Println(registry.GetAddressFromPubKey(pubKey))

				fmt.Println("PubKey:")
				fmt.Println(registry.BytesToBase64(pubKey.Bytes()))

				fmt.Println("Signature:")
				fmt.Println(registry.BytesToBase64(sigBytes))

				return nil
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

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Bool("sign-only", false, "Only sign the transaction payload.")

	return cmd
}
