package swapservice

import (
	"log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the swapservice Querier
const (
	QueryPoolStruct = "poolstruct"
	QueryPoolDatas  = "pooldatas"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPoolStruct:
			return queryPoolStruct(ctx, path[1:], req, keeper)
		case QueryPoolDatas:
			return queryPoolDatas(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown swapservice query endpoint")
		}
	}
}

// nolint: unparam
func queryPoolStruct(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	log.Printf("Got here")
	log.Printf("Path: %s", path[0])
	poolstruct := keeper.GetPoolStruct(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, poolstruct)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryPoolDatas(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var pooldatasList QueryResPoolDatas

	iterator := keeper.GetPoolDatasIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		pooldatasList = append(pooldatasList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, pooldatasList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
