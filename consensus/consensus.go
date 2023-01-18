// Copyright 2021 The nbn Authors
// This file is part of the nbn library.
//
// The nbn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The nbn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the nbn library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different nbn consensus engines.
package consensus

import (
	"bytes"
	"math/big"

	"github.com/token/common"
	"github.com/token/core/state"
	"github.com/token/core/types"
	"github.com/token/params"
	"github.com/token/rpc"
)

type GrantProfitRecord struct {
	Which           uint32
	MinerAddress    common.Address
	BlockNumber     uint64
	Amount          *big.Int
	RevenueAddress  common.Address
	RevenueContract common.Address
	MultiSignature  common.Address
}

type GrantProfitSlice []GrantProfitRecord

func (s GrantProfitSlice) Len() int      { return len(s) }
func (s GrantProfitSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s GrantProfitSlice) Less(i, j int) bool {
	//we need sort reverse, so ...
	isLess := s[i].Amount.Cmp(s[j].Amount)
	if isLess != 0 {
		return isLess > 0
	}
	if s[i].BlockNumber != s[j].BlockNumber {
		return s[i].BlockNumber > s[j].BlockNumber
	}
	if s[i].Which != s[j].Which {
		return s[i].Which > s[j].Which
	}
	if s[i].MinerAddress != s[j].MinerAddress {
		return bytes.Compare(s[i].MinerAddress.Bytes(), s[j].MinerAddress.Bytes()) > 0
	}
	if s[i].RevenueAddress != s[j].RevenueAddress {
		return bytes.Compare(s[i].RevenueAddress.Bytes(), s[j].RevenueAddress.Bytes()) > 0
	}
	if s[i].RevenueContract != s[j].RevenueContract {
		return bytes.Compare(s[i].RevenueContract.Bytes(), s[j].RevenueContract.Bytes()) > 0
	}
	return bytes.Compare(s[i].MultiSignature.Bytes(), s[j].MultiSignature.Bytes()) > 0
}

// ChainHeaderReader defines a small collection of methods needed to access the local
// blockchain during header verification.
type ChainHeaderReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header
}

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the nbn address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	GrantProfit (chain ChainHeaderReader, header *types.Header, state *state.StateDB) ([]GrantProfitRecord, []GrantProfitRecord)

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt, grantProfit []GrantProfitRecord, gasReward *big.Int) error

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt, grantProfit []GrantProfitRecord, gasReward *big.Int) (*types.Block, error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainHeaderReader) []rpc.API

	// Close terminates any background threads maintained by the consensus engine.
	Close() error

	ApplyGenesis(chain ChainHeaderReader, genesisHash common.Hash) error

	VerifyHeaderExtra(chain ChainHeaderReader, header *types.Header, verifyExtra []byte) error
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}
