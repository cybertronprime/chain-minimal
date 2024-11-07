package keeper

import (
	"context"

	"chain-minimal/x/checkers/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	for _, indexedStoredGame := range data.IndexedStoredGameList {
		if err := k.StoredGames.Set(ctx, indexedStoredGame.Index, indexedStoredGame.StoredGame); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var indexedStoredGames []types.IndexedStoredGame
	if err := k.StoredGames.Walk(ctx, nil, func(index string, storedGame types.StoredGame) (bool, error) {
		indexedStoredGames = append(indexedStoredGames, types.IndexedStoredGame{
			Index:      index,
			StoredGame: storedGame,
		})
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Params:                params,
		IndexedStoredGameList: indexedStoredGames,
	}, nil
}