package alien

import (
	"errors"
	"github.com/token/common"
	"github.com/token/ethdb"
	"github.com/token/log"
	"github.com/token/rlp"
	"github.com/token/trie"
	"math/big"
)

var (
	errNotEnoughCoin = errors.New("not enough Coin")
)

type CoinAccount struct {
	Address  common.Address
	Balance  *big.Int
}

func (s *CoinAccount) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(s)
}

func decodeCoinAccount(buf []byte) *CoinAccount {
	s := &CoinAccount{}
		err := rlp.DecodeBytes(buf, s)
		if err != nil {
			return nil
		}else {
			return s
		}
}

func (s *CoinAccount) GetBalance() *big.Int {
	return s.Balance
}

func (s *CoinAccount) SetBalance(amount *big.Int) {
	s.Balance = amount
}

func (s *CoinAccount) AddBalance(amount *big.Int) {
	s.Balance = new(big.Int).Add(s.Balance, amount)
}

func (s *CoinAccount) SubBalance(amount *big.Int) error {
	if s.Balance.Cmp(amount) < 0 {
		return errNotEnoughCoin
	}
	s.Balance = new(big.Int).Sub(s.Balance, amount)
	return nil
}

//=========================================================================
type CoinTrie struct {
	trie 	*trie.SecureTrie
	db      ethdb.Database
	triedb 	*trie.Database
}

func (s *CoinTrie) GetOrNewAccount(addr common.Address) *CoinAccount {
	var obj *CoinAccount
	objData := s.trie.Get(addr.Bytes())
	obj = decodeCoinAccount(objData)
	if obj != nil {
		return obj
	}
	obj = &CoinAccount{
		Address: addr,
		Balance: common.Big0,
	}
	return obj
}

func (s *CoinTrie) TireDB() *trie.Database {
	return s.triedb
}

func (s *CoinTrie) getBalance(addr common.Address) *big.Int {
	obj := s.GetOrNewAccount(addr)
	if obj == nil {
		return common.Big0
	}
	return obj.GetBalance()
}

func (s *CoinTrie) setBalance(addr common.Address, amount *big.Int) {
	obj := s.GetOrNewAccount(addr)
	if obj == nil {
		log.Warn("cointrie setbalance", "result", "failed")
		return
	}
	obj.SetBalance(amount)
	value, _ := obj.Encode()
	s.trie.Update(addr.Bytes(), value)
}

func (s *CoinTrie) addBalance(addr common.Address, amount *big.Int) {
	obj := s.GetOrNewAccount(addr)
	if obj == nil {
		log.Warn("cointrie addbalance", "result", "failed")
		return
	}
	obj.AddBalance(amount)
	value, _ := obj.Encode()
	s.trie.Update(addr.Bytes(), value)
}

func (s *CoinTrie) subBalance(addr common.Address, amount *big.Int) error{
	obj := s.GetOrNewAccount(addr)
	if obj == nil {
		log.Warn("cointrie subbalance", "result", "failed")
		return errNotEnoughCoin
	}
	obj.SubBalance(amount)
	value, _ := obj.Encode()
	s.trie.Update(addr.Bytes(), value)
	return nil
}

func (s *CoinTrie) cmpBalance(addr common.Address, amount *big.Int) int{
	obj := s.GetOrNewAccount(addr)
	if obj == nil {
		log.Warn("cointrie cmpbalance", "error", "load account failed")
		return -1
	}

	return obj.Balance.Cmp(amount)
}

func (s *CoinTrie) Hash() common.Hash {
	return s.trie.Hash()
}

func (s *CoinTrie) commit() (root common.Hash, err error){
	hash, err := s.trie.Commit(nil)
	if err != nil {
		return common.Hash{}, err
	}
	s.triedb.Commit(hash, true, nil)
	return hash, nil
}

//====================================================================================
func NewCoinTrie(root common.Hash, db ethdb.Database) (*CoinTrie, error) {
	triedb := trie.NewDatabase(db)
	tr, err := trie.NewSecure(root, triedb)
	if err != nil {
		log.Warn("cointrie open coin trie failed", "root", root)
		return nil, err
	}

	return &CoinTrie{
		trie: tr,
		db: db,
		triedb: triedb,
	}, nil
}

func (s *CoinTrie) Get (addr common.Address) *big.Int{
	return s.getBalance(addr)
}

func (s *CoinTrie) Set (addr common.Address, amount *big.Int) {
	s.setBalance(addr, amount)
}

func (s *CoinTrie) Add (addr common.Address, amount *big.Int) {
	s.addBalance(addr, amount)
}

func (s *CoinTrie) Sub (addr common.Address, amount *big.Int) error{
	return s.subBalance(addr, amount)
}

func (s *CoinTrie) Del(addr common.Address) {
	s.setBalance(addr, common.Big0)
	s.trie.Delete(addr.Bytes())
}

func (s *CoinTrie) Copy() CoinState {
	root, _ := s.Save(nil)
	trie, _ := NewTrieCoinState(root, s.db)
	return trie
}

func (s *CoinTrie) Load(db ethdb.Database, hash common.Hash) error{
	return nil
}

func (s *CoinTrie) Save(db ethdb.Database) (common.Hash, error){
	return s.commit()
}

func (s *CoinTrie) Root() common.Hash {
	return s.Hash()
}

func (s *CoinTrie) GetAll() map[common.Address]*big.Int {
	found := make(map[common.Address]*big.Int)
	it := trie.NewIterator(s.trie.NodeIterator(nil))
	for it.Next() {
		acc := decodeCoinAccount(it.Value)
		if nil != acc {
			found[acc.Address] = acc.Balance
		}
	}
	return found
}
