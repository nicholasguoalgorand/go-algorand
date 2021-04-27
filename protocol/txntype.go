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

package protocol

// Transaction types indicate different types of transactions that can appear
// in a block.  They are used in the data/transaction package and the REST API.

// TxType is the type of the transaction written to the ledger
type TxType string

const (
	// PaymentTx indicates a payment transaction
	PaymentTx TxType = "pay"

	// KeyRegistrationTx indicates a transaction that registers participation keys
	KeyRegistrationTx TxType = "keyreg"

	// AssetConfigTx creates, re-configures, or destroys an asset
	AssetConfigTx TxType = "acfg"

	// AssetTransferTx transfers assets between accounts (optionally closing)
	AssetTransferTx TxType = "axfer"

	// AssetFreezeTx changes the freeze status of an asset
	AssetFreezeTx TxType = "afrz"

	// ApplicationCallTx allows creating, deleting, and interacting with an application
	ApplicationCallTx TxType = "appl"

	// CompactCertTx records a compact certificate
	CompactCertTx TxType = "cert"

	// UnknownTx signals an error
	UnknownTx TxType = "unknown"
)

// TxTypeToByte converts a TxType to byte encoding
func TxTypeToByte(t TxType) byte {
	switch t {
	case PaymentTx:
		return 0
	case KeyRegistrationTx:
		return 1
	case AssetConfigTx:
		return 2
	case AssetTransferTx:
		return 3
	case AssetFreezeTx:
		return 4
	case ApplicationCallTx:
		return 5
	case CompactCertTx:
		return 6
	default:
		return 7
	}
}

// ByteToTxType converts a byte encoding to TxType
func ByteToTxType(b byte) TxType {
	txTypes := []TxType{PaymentTx, KeyRegistrationTx, AssetConfigTx, AssetTransferTx, AssetFreezeTx, ApplicationCallTx, CompactCertTx, UnknownTx}
	return txTypes[b]
}
