//
// Copyright 2019 Wireline, Inc.
//

package gql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

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

func (r *mutationResolver) BroadcastTxCommit(ctx context.Context, tx string) (*string, error) {

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

func (r *queryResolver) GetResources(ctx context.Context, ids []string) ([]*Resource, error) {
	resources := make([]*Resource, len(ids))
	for index, id := range ids {
		resource, err := r.GetResource(ctx, id)
		if err != nil {
			return nil, err
		}

		resources[index] = resource
	}

	return resources, nil
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
		gqlCoins[index] = Coin{
			Denom:  coin.Denom,
			Amount: int(coin.Amount.Int64()),
		}
	}

	return &Account{
		Address: address,
		Num:     int(account.GetAccountNumber()),
		Seq:     int(account.GetSequence()),
		PubKey:  pubKey,
		Coins:   gqlCoins,
	}, nil
}

func (r *queryResolver) GetResource(ctx context.Context, id string) (*Resource, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	dbID := registry.ID(id)
	if r.keeper.HasResource(sdkContext, dbID) {
		resource := r.keeper.GetResource(sdkContext, dbID)
		return getGQLResource(resource)
	}

	return nil, nil
}

func (r *queryResolver) ListResources(ctx context.Context) ([]*Resource, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	resources := r.keeper.ListResources(sdkContext)
	gqlResponse := make([]*Resource, len(resources))

	for index, resource := range resources {
		gqlResource, err := getGQLResource(resource)
		if err != nil {
			return nil, err
		}

		gqlResponse[index] = gqlResource
	}

	return gqlResponse, nil
}

func getGQLResource(resource registry.Resource) (*Resource, error) {
	ownerID := string(resource.Owner.ID)
	ownerAddress := string(resource.Owner.Address)

	systemAttrs, err := mapToJSONStr(resource.SystemAttributes)
	if err != nil {
		return nil, err
	}

	attrs, err := mapToJSONStr(resource.Attributes)
	if err != nil {
		return nil, err
	}

	links := make([]Link, len(resource.Links))
	for linkIndex := range resource.Links {
		linkAttrs, err := mapToJSONStr(resource.Links[linkIndex])
		if err != nil {
			return nil, err
		}

		links[linkIndex] = Link{
			ID:         resource.Links[linkIndex]["id"].(string),
			Attributes: linkAttrs,
		}
	}

	return &Resource{
		ID:   string(resource.ID),
		Type: resource.Type,
		Owner: Owner{
			ID:      &ownerID,
			Address: &ownerAddress,
		},
		SystemAttributes: systemAttrs,
		Attributes:       attrs,
		Links:            links,
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

	var msg []registry.MsgSetResource
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

func (r *queryResolver) GetBots(ctx context.Context, name []string) ([]*Bot, error) {
	bots := []*Bot{}

	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	resources := r.keeper.ListResources(sdkContext)
	for _, resource := range resources {
		if resource.Type == "Bot" && resource.Attributes != nil {
			// Name is mandatory.
			if resName, ok := resource.Attributes["name"].(string); ok {
				res, err := getGQLResource(resource)
				if err != nil {
					return nil, err
				}

				// dsinvite is optional.
				var dsinviteVal *string
				dsinvite, dsinviteOk := resource.Attributes["dsinvite"].(string)
				if dsinviteOk {
					dsinviteVal = &dsinvite
				}

				// Check for match if any names are passed as input, else return all.
				if len(name) > 0 {
					for _, iterName := range name {
						if iterName == resName {
							bots = append(bots, &Bot{
								Resource: res,
								Name:     resName,
								Dsinvite: dsinviteVal,
							})
						}
					}
				} else {
					bots = append(bots, &Bot{
						Resource: res,
						Name:     resName,
						Dsinvite: dsinviteVal,
					})
				}
			}
		}
	}

	return bots, nil

}

func (r *queryResolver) GetPseudonyms(ctx context.Context, name []string) ([]*Pseudonym, error) {
	pseudonyms := []*Pseudonym{}

	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	resources := r.keeper.ListResources(sdkContext)
	for _, resource := range resources {
		if resource.Type == "Pseudonym" && resource.Attributes != nil {
			// Name is mandatory.
			if resName, ok := resource.Attributes["name"].(string); ok {
				res, err := getGQLResource(resource)
				if err != nil {
					return nil, err
				}

				// dsinvite is optional.
				var dsinviteVal *string
				dsinvite, dsinviteOk := resource.Attributes["dsinvite"].(string)
				if dsinviteOk {
					dsinviteVal = &dsinvite
				}

				// Check for match if any names are passed as input, else return all.
				if len(name) > 0 {
					for _, iterName := range name {
						if iterName == resName {
							pseudonyms = append(pseudonyms, &Pseudonym{
								Resource: res,
								Name:     resName,
								Dsinvite: dsinviteVal,
							})
						}
					}
				} else {
					pseudonyms = append(pseudonyms, &Pseudonym{
						Resource: res,
						Name:     resName,
						Dsinvite: dsinviteVal,
					})
				}
			}
		}
	}

	return pseudonyms, nil
}
