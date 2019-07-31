package greeter

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/cosmos/hellochain/x/greeter/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the hellochain Querier
const (
	ListGreetings = "list"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListGreetings:
			return listGreetings(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown greeter query endpoint")
		}
	}
}

func listGreetings(ctx sdk.Context, params []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	greetingList := NewQueryResGreetings()
	iterator := keeper.GetGreetingsIterator(ctx)

	addr, err := sdk.AccAddressFromBech32(params[0])
	if err != nil {
		return nil, sdk.ErrInvalidAddress("invalid address query parameter")
	}

	for ; iterator.Valid(); iterator.Next() {

		var greeting Greeting

		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &greeting)

		if greeting.Recipient.Equals(addr) {
			greetingList[addr.String()] = append(greetingList[addr.String()], greeting)
		}
	}

	hellos, err2 := codec.MarshalJSONIndent(keeper.cdc, greetingList)

	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return hellos, nil
}
