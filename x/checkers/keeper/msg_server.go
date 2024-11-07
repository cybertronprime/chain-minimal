package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"chain-minimal/x/checkers/types"
	"chain-minimal/x/checkers/rules"
)

type msgServer struct {
	k Keeper
}

var _ types.CheckersTorramServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.CheckersTorramServer {
	return &msgServer{k: keeper}
}

// CreateGame defines the handler for the MsgCreateGame message.
func (ms msgServer) CheckersCreateGm(ctx context.Context, msg *types.ReqCheckersTorram) (*types.ResCheckersTorram, error) {
	if length := len([]byte(msg.Index)); types.MaxIndexLength < length || length < 1 {
		return nil, types.ErrIndexTooLong
	}
	if _, err := ms.k.StoredGames.Get(ctx, msg.Index); err == nil || errors.Is(err, collections.ErrEncoding) {
		return nil, fmt.Errorf("game already exists at index: %s", msg.Index)
	}

	newBoard := rules.New()
	storedGame := types.StoredGame{
		Board: newBoard.String(),
		Turn:  rules.PieceStrings[newBoard.Turn],
		Black: msg.Black,
		Red:   msg.Red,
	}
	if err := storedGame.Validate(); err != nil {
		return nil, err
	}
	if err := ms.k.StoredGames.Set(ctx, msg.Index, storedGame); err != nil {
		return nil, err
	}

	return &types.ResCheckersTorram{}, nil
}