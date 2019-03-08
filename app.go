//
// Copyright 2018 Wireline, Inc.
//

package app

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/wirelineio/wirechain/x/htlc"
	"github.com/wirelineio/wirechain/x/multisig"
	msighandler "github.com/wirelineio/wirechain/x/multisig/handlers"
	"github.com/wirelineio/wirechain/x/registry"
	"github.com/wirelineio/wirechain/x/utxo"

	"github.com/wirelineio/wirechain/gql"
)

const (
	appName = "wirechain"
)

type wirechainApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyTxStore       *sdk.KVStoreKey

	keyHtlcStore     *sdk.KVStoreKey
	keyMultisigStore *sdk.KVStoreKey
	keyAccUtxoStore  *sdk.KVStoreKey
	keyUtxoStore     *sdk.KVStoreKey
	keyRegStore      *sdk.KVStoreKey

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper

	htlcKeeper     htlc.Keeper
	multisigKeeper msighandler.Keeper
	utxoKeeper     utxo.Keeper
	regKeeper      registry.Keeper
}

// NewWirechainApp is a constructor function for wirechainApp
func NewWirechainApp(logger log.Logger, db dbm.DB) *wirechainApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))
	bApp.SetMinimumFees(sdk.Coins{sdk.NewInt64Coin("wire", 1)})

	// Here you initialize your application with the store keys it requires
	var app = &wirechainApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),
		keyTxStore:       sdk.NewKVStoreKey("tx"),

		keyHtlcStore:     sdk.NewKVStoreKey("htlc"),
		keyMultisigStore: sdk.NewKVStoreKey("multisig"),
		keyAccUtxoStore:  sdk.NewKVStoreKey("acc_utxo"),
		keyUtxoStore:     sdk.NewKVStoreKey("utxo"),
		keyRegStore:      sdk.NewKVStoreKey("registry"),
	}

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper)

	// The FeeCollectionKeeper collects transaction fees and renders them to the fee distribution module
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

	app.htlcKeeper = htlc.NewKeeper(app.bankKeeper, app.keyHtlcStore, app.cdc)

	app.multisigKeeper = msighandler.NewKeeper(app.bankKeeper, app.keyMultisigStore, app.cdc)

	app.utxoKeeper = utxo.NewKeeper(app.accountKeeper, app.bankKeeper, app.keyAccUtxoStore, app.keyUtxoStore, app.keyTxStore, app.cdc)

	app.regKeeper = registry.NewKeeper(app.accountKeeper, app.bankKeeper, app.keyRegStore, app.cdc)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The app.Router is the main transaction router where each module registers its routes
	// Register the bank and wirechain routes here
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("htlc", htlc.NewHandler(app.htlcKeeper)).
		AddRoute("multisig", msighandler.NewHandler(app.multisigKeeper)).
		AddRoute("utxo", utxo.NewHandler(app.utxoKeeper)).
		AddRoute("registry", registry.NewHandler(app.regKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute("multisig", msighandler.NewQuerier(app.multisigKeeper)).
		AddRoute("utxo", registry.NewQuerier(app.regKeeper)).
		AddRoute("registry", registry.NewQuerier(app.regKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyTxStore,
		app.keyFeeCollection,

		app.keyHtlcStore,
		app.keyMultisigStore,
		app.keyAccUtxoStore,
		app.keyUtxoStore,
		app.keyRegStore,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	go gql.Server(app.BaseApp, app.regKeeper)

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState struct {
	Accounts []*auth.BaseAccount `json:"accounts"`
}

func (app *wirechainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators does the things
func (app *wirechainApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*auth.BaseAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenesisState{Accounts: accounts}
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)

	htlc.RegisterCodec(cdc)
	multisig.RegisterCodec(cdc)
	utxo.RegisterCodec(cdc)
	registry.RegisterCodec(cdc)

	codec.RegisterCrypto(cdc)

	return cdc
}
