package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chain-minimal/x/checkers/types"
)

func TestQueryServer_GetGame(t *testing.T) {
	keeper, ctx := setupTestKeeper(t)
	srv := NewQueryServerImpl(keeper)

	// Setup a test game
	testGame := types.StoredGame{
		Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:  "b",
		Black: "mini1black",
		Red:   "mini1red",
	}
	err := keeper.StoredGames.Set(ctx, "existing-game", testGame)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *types.ReqCheckersTorramQuery
		want    *types.ResCheckersTorramQuery
		wantErr bool
	}{
		{
			name: "existing game",
			req: &types.ReqCheckersTorramQuery{
				Index: "existing-game",
			},
			want: &types.ResCheckersTorramQuery{
				Game: &testGame,
			},
			wantErr: false,
		},
		{
			name: "non-existent game",
			req: &types.ReqCheckersTorramQuery{
				Index: "not-found",
			},
			want: &types.ResCheckersTorramQuery{
				Game: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.GetCheckersTorramGm(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			if tt.want.Game == nil {
				assert.Nil(t, got.Game)
			} else {
				require.NotNil(t, got.Game)
				assert.Equal(t, tt.want.Game.Board, got.Game.Board)
				assert.Equal(t, tt.want.Game.Turn, got.Game.Turn)
				assert.Equal(t, tt.want.Game.Black, got.Game.Black)
				assert.Equal(t, tt.want.Game.Red, got.Game.Red)
			}
		})
	}
}