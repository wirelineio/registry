package gql

import (
	"context"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wirelineio/wirechain/x/registry"
)

// Resolver is the GQL query resolver.
type Resolver struct {
	baseApp *bam.BaseApp
	keeper  registry.Keeper
}

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) ListResources(ctx context.Context) ([]*Resource, error) {
	sdkContext := r.baseApp.NewContext(true, abci.Header{})

	resources := r.keeper.ListResources(sdkContext)
	gqlResponse := make([]*Resource, len(resources))

	for index, resource := range resources {
		gqlResponse[index] = &Resource{
			ID: string(resource.ID),
		}
	}

	return gqlResponse, nil
}
