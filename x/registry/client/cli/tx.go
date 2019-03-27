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
	"github.com/wirelineio/registry/x/registry"
)

// GetCmdSetResource is the CLI command for creating/updating a record.
func GetCmdSetResource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [payload file path]",
		Short: "Set record.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc).WithChainID(registry.WirelineChainID)

			payload, err := getPayloadFromFile(args[0])
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			signOnly := viper.GetBool("sign-only")
			if signOnly {
				return signResource(payload)
			}

			signer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := registry.NewMsgSetRecord(registry.PayloadToPayloadObj(payload), signer)
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

// GetCmdDeleteResource is the CLI command for deleting a record.
func GetCmdDeleteResource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [payload file path]",
		Short: "Delete record.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc).WithChainID(registry.WirelineChainID)

			payload, err := getPayloadFromFile(args[0])
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			signOnly := viper.GetBool("sign-only")
			if signOnly {
				return signResource(payload)
			}

			signer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := registry.NewMsgDeleteRecord(registry.PayloadToPayloadObj(payload), signer)
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

// GetCmdClearResources is the CLI command for clearing all records.
// NOTE: FOR LOCAL TESTING PURPOSES ONLY!
func GetCmdClearResources(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear records.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc).WithChainID(registry.WirelineChainID)

			cliCtx.PrintResponse = true

			signer, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := registry.NewMsgClearRecords(signer)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// Load payload object from YAML file.
func getPayloadFromFile(filePath string) (registry.Payload, error) {
	var payload registry.Payload

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return payload, err
	}

	err = yaml.Unmarshal(data, &payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

// Sign payload object.
func signResource(payload registry.Payload) error {
	name := viper.GetString("from")

	sigBytes, pubKey, err := registry.GetResourceSignature(payload.Record, name)
	if err != nil {
		return err
	}

	fmt.Println("Address   :", registry.GetAddressFromPubKey(pubKey))
	fmt.Println("PubKey    :", registry.BytesToBase64(pubKey.Bytes()))
	fmt.Println("Signature :", registry.BytesToBase64(sigBytes))

	return nil
}
