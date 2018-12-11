package htlc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Status represents the status of an HTLC.
type Status uint8

// HTLC status enum.
const (
	HtlcCreated  Status = 1
	HtlcRedeemed Status = 2
	HtlcFailed   Status = 3
)

// ObjHtlc is persisted in the KV store.
type ObjHtlc struct {
	Amount         sdk.Coin
	Hash           string
	Locktime       int64
	RedeemAddress  sdk.AccAddress
	TimeoutAddress sdk.AccAddress
	Status         Status
	BlockCreatedAt int64
}

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgAddHtlc:
			return handleMsgAddHtlc(ctx, keeper, msg)
		case MsgRedeemHtlc:
			return handleMsgRedeemHtlc(ctx, keeper, msg)
		case MsgFailHtlc:
			return handleMsgFailHtlc(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized htlc Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgAddHtlc
func handleMsgAddHtlc(ctx sdk.Context, keeper Keeper, msg MsgAddHtlc) sdk.Result {
	if keeper.HasHtlc(ctx, msg.Hash) {
		return sdk.ErrInternal("HTLC by that hash already exists.").Result()
	}

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.TimeoutAddress, sdk.Coins{msg.Amount})
	if err != nil {
		return sdk.ErrInsufficientCoins("Not enough coins to create HTLC.").Result()
	}

	obj := ObjHtlc{
		Amount:         msg.Amount,
		Hash:           msg.Hash,
		Locktime:       msg.Locktime,
		RedeemAddress:  msg.RedeemAddress,
		TimeoutAddress: msg.TimeoutAddress,
		Status:         HtlcCreated,
		BlockCreatedAt: ctx.BlockHeight(),
	}

	keeper.AddHtlc(ctx, obj)

	return sdk.Result{}
}

// Handle MsgRedeemHtlc
func handleMsgRedeemHtlc(ctx sdk.Context, keeper Keeper, msg MsgRedeemHtlc) sdk.Result {
	return sdk.ErrUnknownRequest("MsgRedeemHtlc: Not Implemented.").Result()
}

// Handle MsgFailHtlc
func handleMsgFailHtlc(ctx sdk.Context, keeper Keeper, msg MsgFailHtlc) sdk.Result {
	return sdk.ErrUnknownRequest("MsgFailHtlc: Not Implemented.").Result()
}
