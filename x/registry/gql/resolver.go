//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
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

// WireRegistryTypeBot => Bot.
const WireRegistryTypeBot = "wrn:registry-type:bot"

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
		return getGQLRecord(record)
	}

	return nil, nil
}

func (r *queryResolver) GetRecordsByAttributes(ctx context.Context, attributes []*KeyValueInput) ([]*Record, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	records := r.keeper.ListResources(sdkContext)
	gqlResponse := []*Record{}

	for _, record := range records {
		gqlRecord, err := getGQLRecord(record)
		if err != nil {
			return nil, err
		}

		if matchesOnAttributes(&record, attributes) {
			gqlResponse = append(gqlResponse, gqlRecord)
		}
	}

	return gqlResponse, nil
}

func matchesOnAttributes(record *registry.Record, attributes []*KeyValueInput) bool {
	recAttrs := record.Attributes

	for _, attr := range attributes {
		recAttrVal, recAttrFound := recAttrs[attr.Key]
		if !recAttrFound {
			return false
		}

		if attr.Value.Int != nil {
			recAttrValInt, ok := recAttrVal.(int)
			if !ok || *attr.Value.Int != recAttrValInt {
				return false
			}
		}

		if attr.Value.Float != nil {
			recAttrValFloat, ok := recAttrVal.(float64)
			if !ok || *attr.Value.Float != recAttrValFloat {
				return false
			}
		}

		if attr.Value.String != nil {
			recAttrValString, ok := recAttrVal.(string)
			if !ok || *attr.Value.String != recAttrValString {
				return false
			}
		}

		if attr.Value.Boolean != nil {
			recAttrValBool, ok := recAttrVal.(bool)
			if !ok || *attr.Value.Boolean != recAttrValBool {
				return false
			}
		}

		// TODO(ashwin): Handle arrays.
	}

	return true
}

func getGQLRecord(record registry.Record) (*Record, error) {
	// systemAttrs, err := mapToJSONStr(record.SystemAttributes)
	// if err != nil {
	// 	return nil, err
	// }

	attrs, err := mapToKeyValuePairs(record.Attributes)
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

func mapToKeyValuePairs(attrs map[string]interface{}) ([]*KeyValue, error) {
	kvPairs := []*KeyValue{}

	trueVal := true
	falseVal := false

	for key, value := range attrs {

		kvPair := &KeyValue{
			Key: key,
		}

		switch val := value.(type) {
		case nil:
			kvPair.Value.Null = &trueVal
		case int:
			kvPair.Value.Int = &val
		case float64:
			kvPair.Value.Float = &val
		case string:
			kvPair.Value.String = &val
		case bool:
			kvPair.Value.Boolean = &val
		}

		if kvPair.Value.Null == nil {
			kvPair.Value.Null = &falseVal
		}

		valueType := reflect.ValueOf(value)
		if valueType.Kind() == reflect.Slice {
			// TODO(ashwin): Handle arrays.
		}

		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs, nil
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

	var operationStr = "set"
	if objmap["operation"] != nil {
		err = json.Unmarshal(*objmap["operation"], &operationStr)
		if err != nil {
			return nil, err
		}
	}

	var msgs []sdk.Msg

	switch operationStr {
	case "set":
		{
			var setMsg []registry.MsgSetRecord
			err = json.Unmarshal(*objmap["msg"], &setMsg)
			if err != nil {
				return nil, err
			}
			msgs = []sdk.Msg{setMsg[0]}
		}
	case "delete":
		{
			var deleteMsg []registry.MsgDeleteRecord
			err = json.Unmarshal(*objmap["msg"], &deleteMsg)
			if err != nil {
				return nil, err
			}
			msgs = []sdk.Msg{deleteMsg[0]}
		}
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
		Msgs: msgs,
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

func (r *queryResolver) GetBotsByAttributes(ctx context.Context, attributes []*KeyValueInput) ([]*Bot, error) {
	bots := []*Bot{}

	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	records := r.keeper.ListResources(sdkContext)
	for _, record := range records {
		if record.Type == WireRegistryTypeBot && record.Attributes != nil {
			// Name is mandatory.
			if name, ok := record.Attributes["name"].(string); ok {

				// accessKey is optional.
				var accessKeyVal *string
				accessKey, accessKeyOk := record.Attributes["accessKey"].(string)
				if accessKeyOk {
					accessKeyVal = &accessKey
				}

				if matchesOnAttributes(&record, attributes) {
					res, err := getGQLRecord(record)
					if err != nil {
						return nil, err
					}

					bots = append(bots, &Bot{
						Record:    res,
						Name:      name,
						AccessKey: accessKeyVal,
					})
				}

			}
		}
	}

	return bots, nil

}

// GetStatus returns the registry status.
func (r *queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	return &Status{Version: RegistryVersion}, nil
}
