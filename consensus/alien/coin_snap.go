package alien

import (
	"github.com/token/common"
	"github.com/token/ethdb"
	"math/big"
)

type CoinState interface {
	Set(addr common.Address, amount *big.Int)
	Add(addr common.Address, amount *big.Int)
	Sub(addr common.Address, amount *big.Int) error
	Get(addr common.Address) *big.Int
	Del(addr common.Address)
	Copy() CoinState
	Load(db ethdb.Database, hash common.Hash) error
	Save(db ethdb.Database) (common.Hash, error)
	Root() common.Hash
	GetAll() map[common.Address]*big.Int
}

func NewCoin(root common.Hash,db ethdb.Database) (CoinState,error) {
	state ,err:= NewTrieCoinState(root,db)
	return state,err
}

func NewTrieCoinState(root common.Hash, db ethdb.Database) (CoinState,error) {
	return NewCoinTrie(root,db)
}