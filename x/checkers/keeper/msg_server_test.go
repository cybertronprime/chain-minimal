package keeper

import (
	"context"
	"testing"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	storetypes "cosmossdk.io/core/store"
	"chain-minimal/x/checkers/types"
)

type mockMemoryStore struct {
	storetypes.KVStore
	store map[string][]byte
}

func newMockMemoryStore() *mockMemoryStore {
	return &mockMemoryStore{
		store: make(map[string][]byte),
	}
}

func (m *mockMemoryStore) Get(key []byte) ([]byte, error) { 
	if value, ok := m.store[string(key)]; ok {
		return value, nil
	}
	return nil, nil
}

func (m *mockMemoryStore) Has(key []byte) (bool, error) { 
	_, ok := m.store[string(key)]
	return ok, nil
}

func (m *mockMemoryStore) Set(key, value []byte) { 
	m.store[string(key)] = value
}

func (m *mockMemoryStore) Delete(key []byte) error { 
	delete(m.store, string(key))
	return nil
}

func (m *mockMemoryStore) Iterator(start, end []byte) storetypes.Iterator { 
	return nil
}

func (m *mockMemoryStore) ReverseIterator(start, end []byte) storetypes.Iterator { 
	return nil
}

type mockStoreService struct {
	storetypes.KVStoreService
	store *mockMemoryStore
}

func newMockStoreService() *mockStoreService {
	return &mockStoreService{
		store: newMockMemoryStore(),
	}
}

func (m *mockStoreService) OpenKVStore(ctx context.Context) storetypes.KVStore {
	return m.store.KVStore
}

func setupTestKeeper(t *testing.T) (Keeper, context.Context) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	storeService := newMockStoreService()
	k := NewKeeper(
		cdc,
		nil,
		storeService,
		"authority",
	)
	return k, context.Background()
}

func TestMsgServer_CreateGame(t *testing.T) {
	keeper, ctx := setupTestKeeper(t)
	srv := NewMsgServerImpl(keeper)

	testCases := []struct {
		name    string
		msg     *types.ReqCheckersTorram
		wantErr bool
	}{
		{
			name: "valid game creation",
			msg: &types.ReqCheckersTorram{
				Index: "game1",
				Black: "mini1black",
				Red:   "mini1red",
			},
			wantErr: false,
		},
		{
			name: "index too long",
			msg: &types.ReqCheckersTorram{
				Index: string(make([]byte, types.MaxIndexLength+1)),
				Black: "mini1black",
				Red:   "mini1red",
			},
			wantErr: true,
		},
		{
			name: "duplicate game",
			msg: &types.ReqCheckersTorram{
				Index: "duplicate",
				Black: "mini1black",
				Red:   "mini1red",
			},
			wantErr: false,
		},
		{
			name: "create duplicate game",
			msg: &types.ReqCheckersTorram{
				Index: "duplicate",
				Black: "mini1black",
				Red:   "mini1red",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := srv.CheckersCreateGm(ctx, tc.msg)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			require.NotNil(t, resp)

			if !tc.wantErr {
				storedGame, err := keeper.StoredGames.Get(ctx, tc.msg.Index)
				require.NoError(t, err)
				assert.Equal(t, tc.msg.Black, storedGame.Black)
				assert.Equal(t, tc.msg.Red, storedGame.Red)
			}
		})
	}
}