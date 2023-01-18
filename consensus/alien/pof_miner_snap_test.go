package alien

import (
	"encoding/json"
	"github.com/token/common"
	"github.com/token/core/types"
	"math/big"
	"testing"
)

func TestNewPofMinerSnap_copy(t *testing.T) {
	s := NewPofMinerSnap(100)
	s.PofMinerCache = []string{"1", "2", "3"}
	s.PofMinerPrevCache = []string{"4", "5", "6"}
	clone := s.copy()
	ss, err := json.Marshal(s)
	if err != nil {
		t.Errorf("json marshal error %v", err)
	}
	cloness, err := json.Marshal(clone)
	if err != nil {
		t.Errorf("json marshal error %v", err)
	}
	if string(ss) != string(cloness) {
		t.Error("PofMinerPrevCache not equals")
		t.Errorf("snap: %s", string(ss))
		t.Errorf("clone: %s", string(cloness))
	}
}

func TestUpdatePofMinerDaily(t *testing.T) {
	s := NewPofMinerSnap(100)
	cached := []string{"1", "2", "3"}
	s.PofMinerCache = cached
	s.PofMiner[common.Address{}] = make(map[common.Hash]*PofMinerReport)
	s.PofMiner[common.Address{}][common.Hash{}] = &PofMinerReport{
		Target:       common.Address{},
		Hash:         common.Hash{},
		FlowValue1:   2,
		FlowValue2:   3,
	}
	header := &types.Header{
		ParentHash:  common.Hash{},
		UncleHash:   common.Hash{},
		Coinbase:    common.Address{},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  nil,
		Number:      big.NewInt(100),
		GasLimit:    0,
		GasUsed:     0,
		Time:        0,
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       types.BlockNonce{},
		Initial:     nil,
		BaseFee:     nil,
	}
	s.updatePofMinerDaily(100, header)
	ss, err := json.Marshal(s)
	if err != nil {
		t.Errorf("json marshal error %v", err)
	}
	t.Log("json: ", string(ss))
	if len(s.PofMinerCache) != 0 {
		t.Errorf("error in PofMinerCache. len=%d ", len(s.PofMinerCache))
	}
	if len(s.PofMinerPrevCache) != len(cached) {
		t.Errorf("error in PofMinerPrevCache. len=%d ", len(s.PofMinerPrevCache))
	}
	for i, v := range s.PofMinerPrevCache {
		if v != cached[i] {
			t.Errorf("PofMinerPrevCache not equals, index=%d value=%v should=%v", i, v, cached[i])
		}
	}
	if len(s.PofMiner) != 0 {
		t.Errorf("error in PofMiner. len=%d ", len(s.PofMiner))
	}
	report := s.PofMinerPrev[common.Address{}][common.Hash{}]
	if  report.FlowValue1 != 2 || report.FlowValue2 != 3 {
		t.Errorf("error in PofMinerPrev data")
	}
}
