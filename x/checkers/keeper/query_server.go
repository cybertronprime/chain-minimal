package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"chain-minimal/x/checkers"
)

var _ checkers.CheckersTorramQueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) checkers.CheckersTorramQueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// GetGame defines the handler for the Query/GetGame RPC method.
func (qs queryServer) GetCheckersTorramGm(ctx context.Context, req *checkers.ReqCheckersTorramQuery) (*checkers.ResCheckersTorramQuery, error) {
	game, err := qs.k.StoredGames.Get(ctx, req.Index)
	if err == nil {
		return &checkers.ResCheckersTorramQuery{Game: &game}, nil
	}
	if errors.Is(err, collections.ErrNotFound) {
		return &checkers.ResCheckersTorramQuery{Game: nil}, nil
	}

	return nil, status.Error(codes.Internal, err.Error())
}