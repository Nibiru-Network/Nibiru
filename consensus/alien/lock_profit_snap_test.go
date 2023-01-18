package alien

import (
	"encoding/json"
	"github.com/token/common"
	"math/big"
	"testing"
)

func TestLockData_copy(t *testing.T) {
	data := NewLockData("test")
	data.CacheL1 = []common.Hash{common.Hash{}}
	item := &PledgeItem{
		Amount:          big.NewInt(100),
		PledgeType:      0,
		Playment:        big.NewInt(10),
		LockPeriod:      0,
		RlsPeriod:       0,
		Interval:        0,
		StartHigh:       0,
		TargetAddress:   common.Address{},
		RevenueAddress:  common.Address{},
		RevenueContract: common.Address{},
		MultiSignature:  common.Address{},
	}
	banace := &LockBalanceData{
		RewardBalance: make(map[uint32]*big.Int),
		LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
	}
	banace.LockBalance[1] = make(map[uint32]*PledgeItem)
	banace.LockBalance[1][3] = item
	data.Revenue[common.Address{}] = banace

	clone := data.copy()
	ss, err := json.Marshal(data)
	if err != nil {
		t.Errorf("json marshal error %v", err)
	}
	cloness, err := json.Marshal(clone)
	if err != nil {
		t.Errorf("json marshal error %v", err)
	}
	t.Log("data", string(ss))
	t.Log("clone", string(cloness))
	if string(ss) != string(cloness) {
		t.Error("LockData not equals")
	}
}
