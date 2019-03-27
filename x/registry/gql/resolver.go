//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/wirelineio/registry/x/registry"
)

// Resolver is the GQL query resolver.
type Resolver struct {
	baseApp       *bam.BaseApp
	codec         *codec.Codec
	keeper        registry.Keeper
	accountKeeper auth.AccountKeeper
}

// Account resolver.
func (r *Resolver) Account() AccountResolver {
	return &accountResolver{r}
}

type accountResolver struct{ *Resolver }

// Coin resolver.
func (r *Resolver) Coin() CoinResolver {
	return &coinResolver{r}
}

type coinResolver struct{ *Resolver }

// Mutation is the entry point to tx execution.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

// BigUInt represents a 64-bit unsigned integer.
type BigUInt uint64

func (r *accountResolver) Number(ctx context.Context, obj *Account) (string, error) {
	val := uint64(obj.Number)
	return strconv.FormatUint(val, 10), nil
}

func (r *accountResolver) Sequence(ctx context.Context, obj *Account) (string, error) {
	val := uint64(obj.Sequence)
	return strconv.FormatUint(val, 10), nil
}

func (r *coinResolver) Amount(ctx context.Context, obj *Coin) (string, error) {
	val := uint64(obj.Amount)
	return strconv.FormatUint(val, 10), nil
}

func (r *mutationResolver) Submit(ctx context.Context, tx string) (*string, error) {
	stdTx, err := decodeStdTx(tx)
	if err != nil {
		return nil, err
	}

	res, err := broadcastTx(r, stdTx)
	if err != nil {
		return nil, err
	}

	txHash := res.Hash.String()

	return &txHash, nil
}

func (r *queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*Account, error) {
	accounts := make([]*Account, len(addresses))
	for index, address := range addresses {
		account, err := r.GetAccount(ctx, address)
		if err != nil {
			return nil, err
		}

		accounts[index] = account
	}

	return accounts, nil
}

func (r *queryResolver) GetRecordsByIds(ctx context.Context, ids []string) ([]*Record, error) {
	records := make([]*Record, len(ids))
	for index, id := range ids {
		record, err := r.GetResource(ctx, id)
		if err != nil {
			return nil, err
		}

		records[index] = record
	}

	return records, nil
}

func (r *queryResolver) GetAccount(ctx context.Context, address string) (*Account, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	account := r.accountKeeper.GetAccount(sdkContext, addr)
	if account == nil {
		return nil, nil
	}

	var pubKey *string
	if account.GetPubKey() != nil {
		pubKeyStr := base64.StdEncoding.EncodeToString(account.GetPubKey().Bytes())
		pubKey = &pubKeyStr
	}

	coins := []sdk.Coin(account.GetCoins())
	gqlCoins := make([]Coin, len(coins))

	for index, coin := range account.GetCoins() {
		amount := coin.Amount.Int64()
		if amount < 0 {
			return nil, errors.New("amount cannot be negative")
		}

		gqlCoins[index] = Coin{
			Type:   coin.Denom,
			Amount: BigUInt(amount),
		}
	}

	accNum := BigUInt(account.GetAccountNumber())
	seq := BigUInt(account.GetSequence())

	return &Account{
		Address:  address,
		Number:   accNum,
		Sequence: seq,
		PubKey:   pubKey,
		Balance:  gqlCoins,
	}, nil
}

func (r *queryResolver) GetResource(ctx context.Context, id string) (*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	dbID := registry.ID(id)
	if r.keeper.HasResource(sdkContext, dbID) {
		record := r.keeper.GetResource(sdkContext, dbID)
		return getGQLResource(record)
	}

	return nil, nil
}

func (r *queryResolver) GetRecordsByAttributes(ctx context.Context, namespace *string) ([]*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	records := r.keeper.ListResources(sdkContext, namespace)
	gqlResponse := make([]*Record, len(records))

	for index, record := range records {
		gqlResource, err := getGQLResource(record)
		if err != nil {
			return nil, err
		}

		gqlResponse[index] = gqlResource
	}

	return gqlResponse, nil
}

func getGQLResource(record registry.Record) (*Record, error) {
	// systemAttrs, err := mapToJSONStr(record.SystemAttributes)
	// if err != nil {
	// 	return nil, err
	// }

	attrs, err := mapToJSONStr(record.Attributes)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:         string(record.ID),
		Type:       record.Type,
		Owner:      record.Owner,
		Attributes: attrs,
	}, nil
}

func mapToJSONStr(attrs map[string]interface{}) (*string, error) {
	if len(attrs) == 0 {
		return nil, nil
	}

	attrsJSON, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}

	attrsJSONStr := string(attrsJSON)

	return &attrsJSONStr, nil
}

func decodeStdTx(tx string) (*auth.StdTx, error) {
	bytes, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		return nil, err
	}

	// Note: json.Unmarshal doesn't known which Msg struct to use, so we do it "manually".
	// See https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(bytes, &objmap)
	if err != nil {
		return nil, err
	}

	var msg []registry.MsgSetRecord
	err = json.Unmarshal(*objmap["msg"], &msg)
	if err != nil {
		return nil, err
	}

	var fee auth.StdFee
	err = json.Unmarshal(*objmap["fee"], &fee)
	if err != nil {
		return nil, err
	}

	var sigs []*json.RawMessage
	err = json.Unmarshal(*objmap["signatures"], &sigs)
	if err != nil {
		return nil, err
	}

	var sig map[string]*json.RawMessage
	err = json.Unmarshal(*sigs[0], &sig)
	if err != nil {
		return nil, err
	}

	var pubKeyStr string
	err = json.Unmarshal(*sig["pub_key"], &pubKeyStr)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}

	pubKey, err := cryptoAmino.PubKeyFromBytes(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	var signature []byte
	err = json.Unmarshal(*sig["signature"], &signature)
	if err != nil {
		return nil, err
	}

	var accountNum uint64
	err = json.Unmarshal(*sig["account_number"], &accountNum)
	if err != nil {
		return nil, err
	}

	var sequenceNum uint64
	err = json.Unmarshal(*sig["sequence"], &sequenceNum)
	if err != nil {
		return nil, err
	}

	var memo string
	err = json.Unmarshal(*objmap["memo"], &memo)
	if err != nil {
		return nil, err
	}

	stdTx := auth.StdTx{
		Msgs: []sdk.Msg{msg[0]},
		Fee:  fee,
		Signatures: []auth.StdSignature{auth.StdSignature{
			PubKey:        pubKey,
			Signature:     signature,
			AccountNumber: accountNum,
			Sequence:      sequenceNum,
		}},
		Memo: memo,
	}

	return &stdTx, nil
}

func broadcastTx(r *mutationResolver, stdTx *auth.StdTx) (*ctypes.ResultBroadcastTxCommit, error) {
	txBytes, err := r.Resolver.codec.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return nil, err
	}

	res, err := core.BroadcastTxCommit(txBytes)
	if err != nil {
		return nil, err
	}

	if res.CheckTx.IsErr() {
		return nil, errors.New(res.CheckTx.String())
	}

	if res.DeliverTx.IsErr() {
		return nil, errors.New(res.DeliverTx.String())
	}

	return res, nil
}

func (r *queryResolver) GetBots(ctx context.Context, namespace *string, name []string) ([]*Bot, error) {
	bots := []*Bot{}

	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	records := r.keeper.ListResources(sdkContext, namespace)
	for _, record := range records {
		if record.Type == "Bot" && record.Attributes != nil {
			// Name is mandatory.
			if resName, ok := record.Attributes["name"].(string); ok {
				res, err := getGQLResource(record)
				if err != nil {
					return nil, err
				}

				// accessKey is optional.
				var accessKeyVal *string
				accessKey, accessKeyOk := record.Attributes["accessKey"].(string)
				if accessKeyOk {
					accessKeyVal = &accessKey
				}

				// Check for match if any names are passed as input, else return all.
				if len(name) > 0 {
					for _, iterName := range name {
						if iterName == resName {
							bots = append(bots, &Bot{
								Record:    res,
								Name:      resName,
								AccessKey: accessKeyVal,
							})
						}
					}
				} else {
					bots = append(bots, &Bot{
						Record:    res,
						Name:      resName,
						AccessKey: accessKeyVal,
					})
				}
			}
		}
	}

	return bots, nil

}
