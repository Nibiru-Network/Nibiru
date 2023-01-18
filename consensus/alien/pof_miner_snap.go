package alien

import (
	"bytes"
	"fmt"
	"github.com/token/common"
	"github.com/token/core/types"
	"github.com/token/ethdb"
	"github.com/token/log"
	"github.com/token/rlp"
	"math/big"
)

type PofMinerSnap struct {
	DayStartTime      uint64                                             `json:"dayStartTime"`
	PofMinerPrevTotal uint64                                             `json:"pofminerPrevTotal"`
	PofMiner          map[common.Address]map[common.Hash]*PofMinerReport `json:"pofminerCurr"`
	PofMinerPrev      map[common.Address]map[common.Hash]*PofMinerReport `json:"pofminerPrev"`
	PofMinerCache     []string                                           `json:"pofminerCurCache"`
	PofMinerPrevCache []string                                           `json:"pofminerPrevCache"`
}

func NewPofMinerSnap(dayStartTime uint64) *PofMinerSnap {
	return &PofMinerSnap{
		DayStartTime:      dayStartTime,
		PofMinerPrevTotal: 0,
		PofMiner:          make(map[common.Address]map[common.Hash]*PofMinerReport),
		PofMinerPrev:      make(map[common.Address]map[common.Hash]*PofMinerReport),
		PofMinerCache:     []string{},
		PofMinerPrevCache: []string{},
	}
}

func (s *PofMinerSnap) copy() *PofMinerSnap {
	clone := &PofMinerSnap{
		DayStartTime:      s.DayStartTime,
		PofMinerPrevTotal: s.PofMinerPrevTotal,
		PofMiner:          make(map[common.Address]map[common.Hash]*PofMinerReport),
		PofMinerPrev:      make(map[common.Address]map[common.Hash]*PofMinerReport),
		PofMinerCache:     nil,
		PofMinerPrevCache: nil,
	}
	for who, item := range s.PofMiner {
		clone.PofMiner[who] = make(map[common.Hash]*PofMinerReport)
		for chain, report := range item {
			clone.PofMiner[who][chain] = report.copy()
		}
	}
	for who, item := range s.PofMinerPrev {
		clone.PofMinerPrev[who] = make(map[common.Hash]*PofMinerReport)
		for chain, report := range item {
			clone.PofMinerPrev[who][chain] = report.copy()
		}
	}
	clone.PofMinerCache = make([]string, len(s.PofMinerCache))
	copy(clone.PofMinerCache, s.PofMinerCache)
	clone.PofMinerPrevCache = make([]string, len(s.PofMinerPrevCache))
	copy(clone.PofMinerPrevCache, s.PofMinerPrevCache)
	return clone
}

func (s *PofMinerSnap) updatePofReport(rewardBlock uint64, blockPerDay uint64, pofReport []MinerPofReportRecord, headerNumber *big.Int) {
	for _, items := range pofReport {
		chain := items.ChainHash
		for _, item := range items.ReportContent {
			if _, ok := s.PofMiner[item.Target]; !ok {
				s.PofMiner[item.Target] = make(map[common.Hash]*PofMinerReport)
			}
			if _, ok := s.PofMiner[item.Target][chain]; !ok {
				s.PofMiner[item.Target][chain] = &PofMinerReport{
					Target:       item.Target,
					Hash:         chain,
					FlowValue1:   item.FlowValue1,
					FlowValue2:   item.FlowValue2,
				}
			}else {
				s.PofMiner[item.Target][chain].FlowValue1 += item.FlowValue1
				s.PofMiner[item.Target][chain].FlowValue2 += item.FlowValue2
			}
		}
	}
}

func (s *PofMinerSnap) updatePofMinerDaily(blockPerDay uint64, header *types.Header) {
	if 0 == header.Number.Uint64()%blockPerDay && 0 != header.Number.Uint64() {
		s.DayStartTime = header.Time
		s.PofMinerPrev = make(map[common.Address]map[common.Hash]*PofMinerReport)
		for address, item := range s.PofMiner {
			s.PofMinerPrev[address] = make(map[common.Hash]*PofMinerReport)
			for chain, report := range item {
				s.PofMinerPrev[address][chain] = &PofMinerReport{
					Target:       report.Target,
					Hash:         report.Hash,
					FlowValue1:   report.FlowValue1,
					FlowValue2:   report.FlowValue2,
				}
			}
		}
		s.PofMinerPrevCache = make([]string, len(s.PofMinerCache))
		copy(s.PofMinerPrevCache, s.PofMinerCache)
		s.PofMiner = make(map[common.Address]map[common.Hash]*PofMinerReport)
		s.PofMinerCache = []string{}
	}
}

func (s *PofMinerSnap) cleanPrevFlow() {
	s.PofMinerPrev = make(map[common.Address]map[common.Hash]*PofMinerReport)
	s.PofMinerPrevCache = []string{}
	s.PofMinerPrevTotal = 0
}

func (s *PofMinerSnap) setFlowPrevTotal(total uint64) {
	s.PofMinerPrevTotal = total
}

func (s *PofMinerSnap) accumulateFlows(db ethdb.Database) map[common.Address]*PofMinerReport {
	pofcensus := make(map[common.Address]*PofMinerReport)
	for minerAddress, item := range s.PofMinerPrev {
		for _, bandwidth := range item {
			if _, ok := pofcensus[minerAddress]; !ok {
				pofcensus[minerAddress] = &PofMinerReport{
					FlowValue1:   bandwidth.FlowValue1,
					FlowValue2:   bandwidth.FlowValue2,
				}
			} else {
				pofcensus[minerAddress].FlowValue1 += bandwidth.FlowValue1
				pofcensus[minerAddress].FlowValue2 += bandwidth.FlowValue2
			}
		}
	}
	for _, key := range s.PofMinerPrevCache {
		flows, err := s.load(db, key)
		if err != nil {
			log.Warn("accumulateFlows load cache error", "key", key, "err", err)
			continue
		}
		for _, flow := range flows {
			target := flow.Target
			if _, ok := pofcensus[target]; !ok {
				pofcensus[target] = &PofMinerReport{
					FlowValue1:   flow.FlowValue1,
					FlowValue2:   flow.FlowValue2,
				}
			} else {
				pofcensus[target].FlowValue1 += flow.FlowValue1
				pofcensus[target].FlowValue2 += flow.FlowValue2
			}
		}
	}
	return pofcensus
}

func (s *PofMinerSnap) load(db ethdb.Database, key string) ([]*PofMinerReport, error) {
	items := []*PofMinerReport{}
	blob, err := db.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	int := bytes.NewBuffer(blob)
	err = rlp.Decode(int, &items)
	if err != nil {
		return nil, err
	}
	log.Info("LockProfitSnap load", "key", key, "size", len(items))
	return items, nil
}

func (s *PofMinerSnap) store(db ethdb.Database, number uint64) error {
	items := []*PofMinerReport{}
	for _, flows := range s.PofMiner {
		for _, flow := range flows {
			items = append(items, flow)
		}
	}
	if len(items) == 0 {
		return nil
	}
	err, buf := PofMinerReportEncodeRlp(items)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("pof-%d", number)
	err = db.Put([]byte(key), buf)
	if err != nil {
		return err
	}
	s.PofMinerCache = append(s.PofMinerCache, key)
	s.PofMiner = make(map[common.Address]map[common.Hash]*PofMinerReport)
	log.Info("PofMinerSnap store", "key", key, "len", len(items))
	return nil
}

func PofMinerReportEncodeRlp(items []*PofMinerReport) (error, []byte) {
	out := bytes.NewBuffer(make([]byte, 0, 255))
	err := rlp.Encode(out, items)
	if err != nil {
		return err, nil
	}
	return nil, out.Bytes()
}
