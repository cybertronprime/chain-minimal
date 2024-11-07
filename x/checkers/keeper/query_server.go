package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"chain-minimal/x/checkers/types"
)

var _ types.CheckersTorramQueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) types.CheckersTorramQueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// GetGame defines the handler for the Query/GetGame RPC method.
func (qs queryServer) GetCheckersTorramGm(ctx context.Context, req *types.ReqCheckersTorramQuery) (*types.ResCheckersTorramQuery, error) {
	game, err := qs.k.StoredGames.Get(ctx, req.Index)
	if err == nil {
		return &types.ResCheckersTorramQuery{Game: &game}, nil
	}
	if errors.Is(err, collections.ErrNotFound) {
		return &types.ResCheckersTorramQuery{Game: nil}, nil
	}

	return nil, status.Error(codes.Internal, err.Error())
}