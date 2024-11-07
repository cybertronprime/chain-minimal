package keeper

import (
    "context"
    "testing"
    "time"

    storetypes "cosmossdk.io/core/store"
    "github.com/cosmos/cosmos-sdk/codec"
    codectypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "chain-minimal/x/checkers/types"
	"cosmossdk.io/core/header"
)

type mockMemoryStore struct {
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

func (m *mockMemoryStore) Set(key, value []byte) error {
    m.store[string(key)] = value
    return nil
}

func (m *mockMemoryStore) Delete(key []byte) error {
    delete(m.store, string(key))
    return nil
}

func (m *mockMemoryStore) Iterator(start, end []byte) (storetypes.Iterator, error) {
    return nil, nil
}

func (m *mockMemoryStore) ReverseIterator(start, end []byte) (storetypes.Iterator, error) {
    return nil, nil
}

type mockStoreService struct {
    store storetypes.KVStore
}

func newMockStoreService() *mockStoreService {
    return &mockStoreService{
        store: newMockMemoryStore(),
    }
}

func (m *mockStoreService) OpenKVStore(ctx context.Context) storetypes.KVStore {
    return m.store
}

type mockAddressCodec struct{}

func (m mockAddressCodec) StringToBytes(text string) ([]byte, error) {
    return []byte(text), nil
}

func (m mockAddressCodec) BytesToString(bz []byte) (string, error) {
    return string(bz), nil
}

func mockSDKContext() context.Context {
    // Create HeaderInfo instead of Header
    headerInfo := header.Info{
        Height:  1,
        Time:    time.Now(),
    }
    
    sdkCtx := sdk.Context{}.WithHeaderInfo(headerInfo)
    return sdk.WrapSDKContext(sdkCtx)
}

func setupTestKeeper(t *testing.T) (Keeper, context.Context) {
    interfaceRegistry := codectypes.NewInterfaceRegistry()
    cdc := codec.NewProtoCodec(interfaceRegistry)
    storeService := newMockStoreService()
    addressCodec := mockAddressCodec{}
    
    k := NewKeeper(
        cdc,
        addressCodec,
        storeService,
        "authority",
    )
    ctx := mockSDKContext()
    return k, ctx
}

func TestMsgServer_CreateGame(t *testing.T) {
    keeper, ctx := setupTestKeeper(t)
    srv := NewMsgServerImpl(keeper)
	// Use valid bech32 addresses for testing
	validBlackAddr := "cosmos1prjj6cpa5ftwdzkn4elkkk9vsz06zrtmz8kxjp" // Example valid address
	validRedAddr := "cosmos1nl0a7pkj40dvn609ymha4lt2lkz97rn36z4xwk" // Example valid address
  

    testCases := []struct {
        name    string
        msg     *types.ReqCheckersTorram
        wantErr bool
    }{
        {
            name: "valid game creation",
            msg: &types.ReqCheckersTorram{
                Index: "game1",
                Black: validBlackAddr,
                Red:   validRedAddr,
            },
            wantErr: false,
        },
        {
            name: "index too long",
            msg: &types.ReqCheckersTorram{
                Index: string(make([]byte, types.MaxIndexLength+1)),
                Black: validBlackAddr,
                Red:   validRedAddr,
            },
            wantErr: true,
        },
        {
            name: "duplicate game",
            msg: &types.ReqCheckersTorram{
                Index: "duplicate",
                Black: validBlackAddr,
                Red:   validRedAddr,
            },
            wantErr: false,
        },
        {
            name: "create duplicate game",
            msg: &types.ReqCheckersTorram{
                Index: "duplicate",
                Black: validBlackAddr,
                Red:   validRedAddr,
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

            storedGame, err := keeper.StoredGames.Get(ctx, tc.msg.Index)
            require.NoError(t, err)
            assert.Equal(t, tc.msg.Black, storedGame.Black)
            assert.Equal(t, tc.msg.Red, storedGame.Red)
            assert.Greater(t, storedGame.StartTime, int64(0))
            assert.Equal(t, int64(0), storedGame.EndTime)
        })
    }
}