// Copyright (C) 2019-2021 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package txnsync

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/rpcs"
)

func TestBitmask(t *testing.T) {
	b := make(bitmask, 12)
	b.SetBit(0)
	b.SetBit(2)
	b.SetBit(69)
	for i := 0; i < 80; i++ {
		exists := b.EntryExists(i)
		if i == 0 || i == 2 || i == 69 {
			require.True(t, exists)
		} else {
			require.False(t, exists)
		}
	}
	b.trimBitmask(80)
	b.expandBitmask(80)
	for i := 0; i < 80; i++ {
		exists := b.EntryExists(i)
		if i == 0 || i == 2 || i == 69 {
			require.True(t, exists)
		} else {
			require.False(t, exists)
		}
	}
}

func TestTxnGroupEncodingSmall(t *testing.T) {
	genesisHash := crypto.Hash([]byte("gh"))

	inTxnGroups := []transactions.SignedTxGroup{
		transactions.SignedTxGroup{
			Transactions: []transactions.SignedTxn{
				{
					Txn: transactions.Transaction{
						Type: protocol.PaymentTx,
						Header: transactions.Header{
							Sender:      basics.Address(crypto.Hash([]byte("2"))),
							Fee:         basics.MicroAlgos{Raw: 100},
							GenesisHash: genesisHash,
						},
						PaymentTxnFields: transactions.PaymentTxnFields{
							Receiver: basics.Address(crypto.Hash([]byte("4"))),
							Amount:   basics.MicroAlgos{Raw: 1000},
						},
					},
					Sig: crypto.Signature{1},
				},
			},
		},
		transactions.SignedTxGroup{
			Transactions: []transactions.SignedTxn{
				{
					Txn: transactions.Transaction{
						Type: protocol.PaymentTx,
						Header: transactions.Header{
							Sender:      basics.Address(crypto.Hash([]byte("1"))),
							Fee:         basics.MicroAlgos{Raw: 100},
							GenesisHash: genesisHash,
						},
						PaymentTxnFields: transactions.PaymentTxnFields{
							Receiver: basics.Address(crypto.Hash([]byte("2"))),
							Amount:   basics.MicroAlgos{Raw: 1000},
						},
					},
					Sig: crypto.Signature{2},
				},
				{
					Txn: transactions.Transaction{
						Type: protocol.KeyRegistrationTx,
						Header: transactions.Header{
							Sender:      basics.Address(crypto.Hash([]byte("1"))),
							GenesisHash: genesisHash,
						},
					},
					Sig: crypto.Signature{3},
				},
			},
		},
		transactions.SignedTxGroup{
			Transactions: []transactions.SignedTxn{
				{
					Txn: transactions.Transaction{
						Type: protocol.AssetConfigTx,
						Header: transactions.Header{
							Sender:      basics.Address(crypto.Hash([]byte("1"))),
							Fee:         basics.MicroAlgos{Raw: 100},
							GenesisHash: genesisHash,
						},
					},
					Sig: crypto.Signature{4},
				},
				{
					Txn: transactions.Transaction{
						Type: protocol.AssetFreezeTx,
						Header: transactions.Header{
							Sender:      basics.Address(crypto.Hash([]byte("1"))),
							GenesisHash: genesisHash,
						},
					},
					Sig: crypto.Signature{5},
				},
				{
					Txn: transactions.Transaction{
						Type: protocol.CompactCertTx,
						Header: transactions.Header{
							Sender:      basics.Address(crypto.Hash([]byte("1"))),
							GenesisHash: genesisHash,
						},
					},
					Msig: crypto.MultisigSig{Version: 1},
				},
			},
		},
	}
	addGroupHashes(inTxnGroups, 6, []byte{1})
	encodedGroupsBytes := encodeTransactionGroups(inTxnGroups)
	fmt.Println(len(encodedGroupsBytes))
	fmt.Println(string(encodedGroupsBytes))
	out, err := decodeTransactionGroups(encodedGroupsBytes)
	require.NoError(t, err)
	require.ElementsMatch(t, inTxnGroups, out)
}

func txnGroupsData() ([]transactions.SignedTxGroup, error) {
	dat, err := ioutil.ReadFile("txns.txt")
	if err != nil {
		return nil, err
	}
	dec := protocol.NewDecoderBytes(dat)
	ntx := 0
	blocksData := make([]rpcs.EncodedBlockCert, 1)
	for {
		if len(blocksData) == ntx {
			n := make([]rpcs.EncodedBlockCert, len(blocksData)*2)
			copy(n, blocksData)
			blocksData = n
		}

		err := dec.Decode(&blocksData[ntx])
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ntx++
	}
	blocksData = blocksData[:ntx]

	var txnGroups []transactions.SignedTxGroup
	for _, blockData := range blocksData {
		block := blockData.Block
		payset, err := block.DecodePaysetGroups()
		if err != nil {
			return nil, err
		}
		for _, txns := range payset {
			var txnGroup transactions.SignedTxGroup
			for _, txn := range txns {
				txnGroup.Transactions = append(txnGroup.Transactions, txn.SignedTxn)
			}
			txnGroups = append(txnGroups, txnGroup)
		}
	}
	return txnGroups, nil
}

func TestTxnGroupEncodingLarge(t *testing.T) {
	txnGroups, err := txnGroupsData()
	require.NoError(t, err)

	encodedGroupsBytes := encodeTransactionGroups(txnGroups)
	out, err := decodeTransactionGroups(encodedGroupsBytes)
	require.NoError(t, err)
	require.ElementsMatch(t, txnGroups, out)

	encodedGroupsBytes = encodeTransactionGroupsOld(txnGroups)
	out, err = decodeTransactionGroupsOld(encodedGroupsBytes)
	require.NoError(t, err)
	require.ElementsMatch(t, txnGroups, out)

	count := make(map[protocol.TxType]int)
	for _, txg := range txnGroups {
		for _, txn := range txg.Transactions {
			count[txn.Txn.Type] += 1
		}
	}
	fmt.Println(count)
}

func BenchmarkTxnGroupEncoding (b *testing.B) {
	txnGroups, err := txnGroupsData()
	require.NoError(b, err)
	var encodedGroupsBytes []byte

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodedGroupsBytes = encodeTransactionGroups(txnGroups)
		releaseEncodedTransactionGroups(encodedGroupsBytes)
	}

	fmt.Println("new data: ", len(encodedGroupsBytes))
}

func BenchmarkTxnGroupDecoding (b *testing.B) {
	txnGroups, err := txnGroupsData()
	require.NoError(b, err)

	encodedGroupsBytes := encodeTransactionGroups(txnGroups)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = decodeTransactionGroups(encodedGroupsBytes)
		require.NoError(b, err)
	}
}

func BenchmarkTxnGroupEncodingOld (b *testing.B) {
	txnGroups, err := txnGroupsData()
	require.NoError(b, err)
	var encodedGroupsBytes []byte

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodedGroupsBytes = encodeTransactionGroupsOld(txnGroups)
		releaseEncodedTransactionGroups(encodedGroupsBytes)
	}

	fmt.Println("old data: ", len(encodedGroupsBytes))
}

func BenchmarkTxnGroupDecodingOld (b *testing.B) {
	txnGroups, err := txnGroupsData()
	require.NoError(b, err)

	encodedGroupsBytes := encodeTransactionGroupsOld(txnGroups)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = decodeTransactionGroupsOld(encodedGroupsBytes)
		require.NoError(b, err)
	}
}