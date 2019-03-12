package gql

import (
	"context"
	"encoding/base64"
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wirechain/x/registry"
)

// Resolver is the GQL query resolver.
type Resolver struct {
	baseApp       *bam.BaseApp
	keeper        registry.Keeper
	accountKeeper auth.AccountKeeper
}

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

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

	pubKey := base64.StdEncoding.EncodeToString(account.GetPubKey().Bytes())

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
		PubKey:  &pubKey,
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