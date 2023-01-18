package types

import (
	"github.com/token/common"
	"math/big"
)

type SignerList []SignerTuple

type SignerTuple struct {
	V, R, S *big.Int
}

type MultiSignerTx struct {
	ChainID    *big.Int        // destination chain ID
	Nonce      uint64          // nonce of sender account
	GasPrice   *big.Int        // wei per gas
	Gas        uint64          // gas limit
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int        // wei amount
	Data       []byte          // contract invocation input data
	AccessList AccessList
	SignerList SignerList      // EIP-2930 access list
	V, R, S    *big.Int        // signature values
}

func NewMultiSignerTransaction(chainId *big.Int, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return NewTx(&MultiSignerTx{
		ChainID:  chainId,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    amount,
		Data:     data,
	})
}

func (tx *MultiSignerTx) copy() TxData {
	cpy := &MultiSignerTx{
		ChainID:    new(big.Int),
		Nonce:      tx.Nonce,
		GasPrice:   new(big.Int),
		Gas:        tx.Gas,
		To:         tx.To, // TODO: copy pointed-to address
		Value:      new(big.Int),
		Data:       common.CopyBytes(tx.Data),
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		SignerList: make(SignerList, len(tx.SignerList)),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	copy(cpy.AccessList, tx.AccessList)
	copy(cpy.SignerList, tx.SignerList)
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

func (tx *MultiSignerTx) innerCopy(v, r, s *big.Int) TxData {
	cpy := &MultiSignerTx{
		ChainID:    new(big.Int),
		Nonce:      tx.Nonce,
		GasPrice:   new(big.Int),
		Gas:        tx.Gas,
		To:         tx.To, // TODO: copy pointed-to address
		Value:      new(big.Int),
		Data:       common.CopyBytes(tx.Data),
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		SignerList: make(SignerList, len(tx.SignerList)),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	copy(cpy.AccessList, tx.AccessList)
	if v != nil {
		cpy.V.Set(v)
	}
	if r != nil {
		cpy.R.Set(r)
	}
	if s != nil {
		cpy.S.Set(s)
	}
	return cpy
}

// accessors for innerTx.
func (tx *MultiSignerTx) txType() byte                                 { return MultiSignerTxType }
func (tx *MultiSignerTx) chainID() *big.Int                            { return tx.ChainID }
func (tx *MultiSignerTx) nonce() uint64                                { return tx.Nonce }
func (tx *MultiSignerTx) gasPrice() *big.Int                           { return tx.GasPrice }
func (tx *MultiSignerTx) gasTipCap() *big.Int                          { return tx.GasPrice }
func (tx *MultiSignerTx) gasFeeCap() *big.Int                          { return tx.GasPrice }
func (tx *MultiSignerTx) gas() uint64                                  { return tx.Gas }
func (tx *MultiSignerTx) to() *common.Address                          { return tx.To }
func (tx *MultiSignerTx) value() *big.Int                              { return tx.Value }
func (tx *MultiSignerTx) data() []byte                                 { return tx.Data }
func (tx *MultiSignerTx) accessList() AccessList                       { return tx.AccessList }
func (tx *MultiSignerTx) signerList() SignerList                       { return tx.SignerList }
func (tx *MultiSignerTx) protected() bool                              { return true }
func (tx *MultiSignerTx) rawSignatureValues() (v, r, s *big.Int)       { return tx.V, tx.R, tx.S }
func (tx *MultiSignerTx) setSignatureValues(chainID, v, r, s *big.Int) {
	signer := LatestSignerForChainID(chainID)
	trans := NewTx(tx.innerCopy(v, r, s))
	currentSinger, err := signer.Sender(trans)
	if nil != err {
		return
	}
	if nil == tx.V || nil == tx.R || nil == tx.S {
		tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
	} else if 0 != chainID.Cmp(tx.ChainID) {
		return
	}

	trans = NewTx(tx.innerCopy(tx.V, tx.R, tx.S))
	initiatorSinger, err := signer.Sender(trans)
	if nil != err || currentSinger.String() == initiatorSinger.String() {
		tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
		return
	}

	for _, sign := range tx.SignerList {
		trans = NewTx(tx.innerCopy(sign.V, sign.R, sign.S))
		assistSinger, err := signer.Sender(trans)
		if nil != err || currentSinger.String() == assistSinger.String() {
			sign.V, sign.R, sign.S = v, r, s
			return
		}
	}

	tx.SignerList = append(tx.SignerList, SignerTuple{
		V: v,
		R: r,
		S: s,
	})
}
func (tx *MultiSignerTx) getAllSigners() []common.Address {
	var allSigners []common.Address
	signersMap := make(map[common.Address]bool)
	if nil == tx.V || nil == tx.R || nil == tx.S {
		return nil
	}
	signer := LatestSignerForChainID(tx.ChainID)
	trans := NewTx(tx.innerCopy(tx.V, tx.R, tx.S))
	initiatorSinger, err := signer.Sender(trans)
	if nil != err {
		return nil
	}

	allSigners = append(allSigners, initiatorSinger)
	signersMap[initiatorSinger] = true
	for _, sign := range tx.SignerList {
		trans = NewTx(tx.innerCopy(sign.V, sign.R, sign.S))
		assistSinger, err := signer.Sender(trans)
		if nil == err {
			if _, ok := signersMap[assistSinger]; !ok {
				allSigners = append(allSigners, assistSinger)
				signersMap[assistSinger] = true
			}
		}
	}

	return allSigners
}