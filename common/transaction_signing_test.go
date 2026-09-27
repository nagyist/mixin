package common

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/MixinNetwork/mixin/crypto"
	"github.com/stretchr/testify/require"
)

func TestSignInputOrder(t *testing.T) {
	account := deterministicAddress(31)
	funding := NewTransactionV5(XINAssetId)
	funding.AddInput(crypto.Blake3Hash([]byte("funding")), 0)
	for i := range 2 {
		funding.AddScriptOutput([]*Address{&account}, NewThresholdScript(1), NewInteger(1), bytes.Repeat([]byte{byte(32 + i)}, 64))
	}
	utxos := funding.AsVersioned().UnspentOutputs()
	store := &campaignStore{utxos: make(map[string]*UTXOWithLock)}
	for _, u := range utxos {
		store.utxos[utxoRef(u.Hash, u.Index)] = u
	}
	for _, order := range [][]int{{0, 1}, {1, 0}, {0, 0, 1}, {0, 1, 0}, {1, 0, 1}} {
		t.Run(fmt.Sprint(order), func(t *testing.T) {
			tx := NewTransactionV5(XINAssetId)
			for _, u := range utxos {
				tx.AddInput(u.Hash, u.Index)
			}
			tx.AddScriptOutput([]*Address{&account}, NewThresholdScript(1), NewInteger(2), bytes.Repeat([]byte{34}, 64))
			ver := tx.AsVersioned()
			for _, i := range order {
				require.NoError(t, ver.SignInput(store, i, []*Address{&account}))
			}
			require.Len(t, ver.SignaturesMap, len(tx.Inputs))
			require.NoError(t, ver.Validate(store, 1, false))
			decoded, err := UnmarshalVersionedTransaction(ver.Marshal())
			require.NoError(t, err)
			require.NoError(t, decoded.Validate(store, 1, false))
		})
	}
}

func TestSignInputRejectsNegativeIndex(t *testing.T) {
	account := deterministicAddress(35)
	tx := NewTransactionV5(XINAssetId)
	tx.AddInput(crypto.Blake3Hash([]byte("input")), 0)
	ver := tx.AsVersioned()
	require.ErrorContains(t, ver.SignInput(nil, -1, []*Address{&account}), "invalid input index")
	require.Empty(t, ver.SignaturesMap)
}

func TestSignRawReplacesExistingSignature(t *testing.T) {
	account := deterministicAddress(36)
	for _, mint := range []bool{false, true} {
		t.Run(fmt.Sprintf("mint=%t", mint), func(t *testing.T) {
			tx := NewTransactionV5(XINAssetId)
			if mint {
				tx.AddUniversalMintInput(1, NewInteger(1))
			} else {
				tx.AddDepositInput(&DepositData{Chain: BitcoinAssetId, AssetKey: "btc", Transaction: "deposit", Amount: NewInteger(1)})
			}
			tx.AddScriptOutput([]*Address{&account}, NewThresholdScript(1), NewInteger(1), bytes.Repeat([]byte{37}, 64))
			ver := tx.AsVersioned()
			for range 2 {
				require.NoError(t, ver.SignInput(nil, 0, []*Address{&account}))
			}
			require.Len(t, ver.SignaturesMap, 1)
			require.Len(t, ver.SignaturesMap[0], 1)
			require.True(t, account.PublicSpendKey.Verify(ver.PayloadHash(), *ver.SignaturesMap[0][0]))
		})
	}
}
