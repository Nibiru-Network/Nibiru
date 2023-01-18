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

// Package alien implements the delegated-proof-of-stake consensus engine.

package alien

import (
	"bytes"
	"container/list"
	"github.com/token/common"
	"github.com/token/consensus"
	"github.com/token/core/types"
	"github.com/token/ethdb"
	"github.com/token/log"
	"github.com/token/rlp"
	"github.com/token/rpc"
	"math/big"
	"sync"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the delegated-proof-of-stake scheme.
type API struct {
	chain consensus.ChainHeaderReader
	alien *Alien
	sCache *list.List
	lock sync.RWMutex
}

type SnapCache struct {
	number uint64
	s *Snapshot
}

type MinerSlice []common.Address

func (s MinerSlice) Len() int      { return len(s) }
func (s MinerSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s MinerSlice) Less(i, j int) bool {
	return bytes.Compare(s[i].Bytes(), s[j].Bytes()) > 0
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.getSnapshotCache(header)
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.getSnapshotCache(header)
}

// GetSnapshotAtNumber retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtNumber(number uint64) (*Snapshot, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.getSnapshotCache(header)
}

//y add method
func (api *API) GetSnapshotSignerAtNumber(number uint64) (*SnapshotSign, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.alien.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
	if err != nil {
		log.Warn("Fail to GetSnapshotSignAtNumber", "err", err)
		return nil, errUnknownBlock
	}
	snapshotSign := &SnapshotSign{
		LoopStartTime:snapshot.LoopStartTime,
		Signers: snapshot.Signers,
		Punished: snapshot.Punished,
		SignPledge:make(map[common.Address]*SignPledgeItem),
	}
	for miner,item:=range snapshot.PosPledge{
		snapshotSign.SignPledge[miner]=&SignPledgeItem{
			TotalAmount: new(big.Int).Set(item.TotalAmount),
			LastPunish:item.LastPunish,
			DisRate:new(big.Int).Set(item.DisRate),
		}
	}
	return snapshotSign, err
}


type SnapshotSign struct {
	LoopStartTime   uint64                                              `json:"loopStartTime"`
	Signers         []*common.Address                                   `json:"signers"`
	Punished        map[common.Address]uint64                           `json:"punished"`
	SignPledge      map[common.Address]*SignPledgeItem                  `json:"signpledge"`
}
type SignPledgeItem struct {
	TotalAmount *big.Int                      `json:"totalamount"`
	LastPunish  uint64                        `json:"lastpunish"`
	DisRate     *big.Int                      `json:"distributerate"`
}

func (api *API) GetSnapshotReleaseAtNumber(number uint64,part string) (*SnapshotRelease, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.getSnapshotCache(header)
	if err != nil {
		log.Warn("Fail to GetSnapshotSignAtNumber", "err", err)
		return nil, errUnknownBlock
	}
	snapshotRelease := &SnapshotRelease{
		Revenue:         make(map[common.Address]*LockBalanceData),
	}
	if part!=""{
		if part =="pofexit"{
			snapshotRelease.appendFRlockData(snapshot.Revenue.PofExitLock,api.alien.db)
		}else if part =="rewardlock"{
			snapshotRelease.appendFRlockData(snapshot.Revenue.RewardLock,api.alien.db)
		}else if part =="poflock"{
			snapshotRelease.appendFRlockData(snapshot.Revenue.PofLock,api.alien.db)
		}else if part =="inspirelock"{
			snapshotRelease.appendFRlockData(snapshot.Revenue.PofInspireLock,api.alien.db)
		}else if part =="posexit"{
			snapshotRelease.appendFRlockData(snapshot.Revenue.PosExitLock,api.alien.db)
		}
	}else{
		snapshotRelease.appendFRlockData(snapshot.Revenue.PofExitLock,api.alien.db)
		snapshotRelease.appendFRlockData(snapshot.Revenue.RewardLock,api.alien.db)
		snapshotRelease.appendFRlockData(snapshot.Revenue.PofLock,api.alien.db)
		snapshotRelease.appendFRlockData(snapshot.Revenue.PofInspireLock,api.alien.db)
		snapshotRelease.appendFRlockData(snapshot.Revenue.PosExitLock,api.alien.db)
	}
	return snapshotRelease, err
}

func (s *SnapshotRelease) appendFRItems(items []*PledgeItem) {
	for _, item := range items {
		if _, ok := s.Revenue[item.TargetAddress]; !ok {
			s.Revenue[item.TargetAddress] = &LockBalanceData{
				RewardBalance:make(map[uint32]*big.Int),
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		revenusTarget := s.Revenue[item.TargetAddress]
		if _, ok := revenusTarget.LockBalance[item.StartHigh]; !ok {
			revenusTarget.LockBalance[item.StartHigh] = make(map[uint32]*PledgeItem)
		}
		lockBalance := revenusTarget.LockBalance[item.StartHigh]
		lockBalance[item.PledgeType] = item
	}
}

func (sr *SnapshotRelease) appendFR(Revenue map[common.Address]*LockBalanceData) (error) {
	fr1:= Revenue
	for t1, item1 := range fr1 {
		if _, ok := sr.Revenue[t1]; !ok {
			sr.Revenue[t1] = &LockBalanceData{
				RewardBalance:make(map[uint32]*big.Int),
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		rewardBalance:=item1.RewardBalance
		for t2, item2 := range rewardBalance {
			sr.Revenue[t1].RewardBalance[t2]=item2
		}
		lockBalance:=item1.LockBalance
		for t3, item3 := range lockBalance {
			if _, ok := sr.Revenue[t1].LockBalance[t3]; !ok {
				sr.Revenue[t1].LockBalance[t3] = make(map[uint32]*PledgeItem)
			}
			sr.Revenue[t1].LockBalance[t3]=item3
		}
	}
	return nil
}


func (sr *SnapshotRelease) appendFRlockData(lockData *LockData,db ethdb.Database) (error) {
	sr.appendFR(lockData.Revenue)
	items, err := lockData.loadCacheL1(db)
	if err == nil {
		sr.appendFRItems(items)
	}
	items, err = lockData.loadCacheL2(db)
	if err == nil {
		sr.appendFRItems(items)
	}
	return nil
}


type SnapshotRelease struct {
	Revenue         map[common.Address]*LockBalanceData `json:"revenve"`
}

func (api *API) GetSnapshotPofAtNumber(number uint64) (*SnapshotPof, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	lockReward:=make([]LockRecord,0)
	if len(headerExtra.LockReward)>0 {
		for _, item := range headerExtra.LockReward {
				lockReward=append(lockReward, LockRecord{
					Target: item.Target,
					Amount: item.Amount,
					FlowValue1: item.FlowValue1,
					FlowValue2: item.FlowValue2,
				})
		}
	}
	snapshotPof := &SnapshotPof{
		LockReward: lockReward,
	}
	return snapshotPof, err
}

type SnapshotPof struct {
	LockReward  []LockRecord `json:"lockrecords"`
}

type LockRecord struct {
	Target   common.Address
	Amount   *big.Int
	FlowValue1 uint64 `json:"realFlowvalue"`
	FlowValue2 uint64 `json:"validFlowvalue"`
}

func (api *API) GetSnapshotPofMinerAtNumber(number uint64) (*SnapshotPofMiner, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.getSnapshotCache(header)
	if err != nil {
		log.Warn("Fail to GetSnapshotPofMinerAtNumber", "err", err)
		return nil, errUnknownBlock
	}
	flowMiner := &SnapshotPofMiner{
		DayStartTime:snapshot.PofMiner.DayStartTime,
		FlowMinerPrevTotal: snapshot.PofMiner.PofMinerPrevTotal,
		FlowMiner: snapshot.PofMiner.PofMiner,
		FlowMinerPrev:snapshot.PofMiner.PofMinerPrev,
		FlowMinerReport:[]*PofMinerReport{},
		FlowMinerPrevReport:[]*PofMinerReport{},
	}
	fMiner:=snapshot.PofMiner
	db:=api.alien.db
	items:=flowMiner.loadFlowMinerCache(fMiner,fMiner.PofMinerCache,db)
	flowMiner.FlowMinerReport=append(flowMiner.FlowMinerReport,items...)
	items=flowMiner.loadFlowMinerCache(fMiner,fMiner.PofMinerPrevCache,db)
	flowMiner.FlowMinerPrevReport=append(flowMiner.FlowMinerPrevReport,items...)
	return flowMiner, err
}


type SnapshotPofMiner struct {
	DayStartTime       uint64                                             `json:"dayStartTime"`
	FlowMinerPrevTotal uint64                                             `json:"flowminerPrevTotal"`
	FlowMiner          map[common.Address]map[common.Hash]*PofMinerReport `json:"flowminerCurr"`
	FlowMinerReport    []*PofMinerReport                                  `json:"flowminerReport"`
	FlowMinerPrev      map[common.Address]map[common.Hash]*PofMinerReport `json:"flowminerPrev"`
	FlowMinerPrevReport    []*PofMinerReport                              `json:"flowminerPrevReport"`
}

func (sf *SnapshotPofMiner) loadFlowMinerCache(fMiner *PofMinerSnap,flowMinerCache []string,db ethdb.Database) ([]*PofMinerReport) {
	item:=[]*PofMinerReport{}
	for _, key := range flowMinerCache {
		flows, err := fMiner.load(db, key)
		if err != nil {
			log.Warn("appendFlowMinerCache load cache error", "key", key, "err", err)
			continue
		}
		item=append(item,flows...)
	}
	return item
}



func (api *API) GetSnapshotPofReportAtNumber(number uint64) (*SnapshotPofReport, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	flowReport:=make([]MinerPofReportRecord,0)
	if len(headerExtra.PofReport)>0 {
		flowReport=append(flowReport,headerExtra.PofReport...)
	}
	snapshotPofReport := &SnapshotPofReport{
		FlowReport: flowReport,
	}
	return snapshotPofReport, err
}

type SnapshotPofReport struct {
	FlowReport []MinerPofReportRecord `json:"flowreport"`
}

type SnapshotAddrCoin struct {
	AddrCoinBal *big.Int `json:"addrcoinbal"`
}

func (api *API) GetCoinBalanceAtNumber(address common.Address,number uint64) (*SnapshotAddrCoin,error) {
	log.Info("api GetCoinBalanceAtNumber", "address",address,"number", number)
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.getSnapshotCache(header)
	if err != nil {
		log.Warn("Fail to GetCoinBalanceAtNumber", "err", err)
		return nil, errUnknownBlock
	}

	snapshotAddrCoin :=&SnapshotAddrCoin{
		AddrCoinBal: common.Big0,
	}
	if snapshot.Coin !=nil{
		snapshotAddrCoin.AddrCoinBal = snapshot.Coin.Get(address)
	}

	return snapshotAddrCoin,nil
}

func (api *API) GetCoinBalance(address common.Address) (*SnapshotAddrCoin,error) {
	log.Info("api GetCoinBalance", "address",address)
	header := api.chain.CurrentHeader()
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.GetCoinBalanceAtNumber(address,header.Number.Uint64())
}

func (api *API) getSnapshotCache(header *types.Header) (*Snapshot, error) {
	number:=header.Number.Uint64()
	s:=api.findInSnapCache(number)
	if nil!=s{
		return s,nil
	}
	return api.getSnapshotByHeader(header)
}

func (api *API)findInSnapCache(number uint64) *Snapshot {
	for i := api.sCache.Front(); i != nil; i = i.Next() {
		v:=i.Value.(SnapCache)
		if v.number==number{
			return v.s
		}
	}
	return nil
}

func (api *API) getSnapshotByHeader(header *types.Header) (*Snapshot,error) {
	api.lock.Lock()
	defer api.lock.Unlock()
	number:=header.Number.Uint64()
	s:=api.findInSnapCache(number)
	if nil!=s{
		return s,nil
	}
	cacheSize:=32
	snapshot,err:= api.alien.snapshot(api.chain, number, header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
	if err != nil {
		log.Warn("Fail to getSnapshotByHeader", "err", err)
		return nil, errUnknownBlock
	}
	api.sCache.PushBack(SnapCache{
		number: number,
		s:snapshot,
	})
	if api.sCache.Len()>cacheSize{
		api.sCache.Remove(api.sCache.Front())
	}
	return snapshot,nil
}

func (api *API) GetCoinBalAtNumber(number uint64) (*SnapshotCoin, error) {
	log.Info("api GetCoinBalAtNumber", "number", number)
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.getSnapshotCache(header)
	if err != nil {
		log.Warn("Fail to GetCoinBalAtNumber", "err", err)
		return nil, errUnknownBlock
	}

	snapshotCoin :=&SnapshotCoin{
		CoinBal: make(map[common.Address]*big.Int),
	}
	if snapshot.Coin !=nil{
		coinBal := snapshot.Coin.GetAll()
		if err==nil{
			snapshotCoin.CoinBal = coinBal
		}
	}
	return snapshotCoin, err
}
type SnapshotCoin struct {
	CoinBal map[common.Address]*big.Int `json:"coinbal"`
}

func (api *API) GetSnapshotReleaseAtNumber2(number uint64,part string,startLNum uint64,endLNum uint64) (*SnapshotRelease, error) {
	log.Info("api GetSnapshotReleaseAtNumber2", "number",number,"part",part,"startLNum",startLNum,"endLNum",endLNum)
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	snapshot,err:= api.getSnapshotCache(header)
	if err != nil {
		log.Warn("Fail to GetSnapshotReleaseAtNumber2", "err", err)
		return nil, errUnknownBlock
	}
	snapshotRelease := &SnapshotRelease{
		Revenue: make(map[common.Address]*LockBalanceData),
	}
	if part!=""{
		if part =="pofexit"{
			snapshotRelease.appendFRlockData2(snapshot.Revenue.PofExitLock,api.alien.db,startLNum,endLNum)
		}else if part =="rewardlock"{
			snapshotRelease.appendFRlockData2(snapshot.Revenue.RewardLock,api.alien.db,startLNum,endLNum)
		}else if part =="poflock"{
			snapshotRelease.appendFRlockData2(snapshot.Revenue.PofLock,api.alien.db,startLNum,endLNum)
		}else if part =="inspirelock"{
			snapshotRelease.appendFRlockData2(snapshot.Revenue.PofInspireLock,api.alien.db,startLNum,endLNum)
		}else if part =="posexit"{
			snapshotRelease.appendFRlockData2(snapshot.Revenue.PosExitLock,api.alien.db,startLNum,endLNum)
		}
	}else{
		snapshotRelease.appendFRlockData2(snapshot.Revenue.PofExitLock,api.alien.db,startLNum,endLNum)
		snapshotRelease.appendFRlockData2(snapshot.Revenue.RewardLock,api.alien.db,startLNum,endLNum)
		snapshotRelease.appendFRlockData2(snapshot.Revenue.PofLock,api.alien.db,startLNum,endLNum)
		snapshotRelease.appendFRlockData2(snapshot.Revenue.PofInspireLock,api.alien.db,startLNum,endLNum)
		snapshotRelease.appendFRlockData2(snapshot.Revenue.PosExitLock,api.alien.db,startLNum,endLNum)
	}
	return snapshotRelease, err
}

func (sr *SnapshotRelease) appendFRlockData2(lockData *LockData,db ethdb.Database,startLNum uint64,endLNum uint64) (error) {
	sr.appendFR2(lockData.Revenue,startLNum,endLNum)
	items, err := lockData.loadCacheL1(db)
	if err == nil {
		sr.appendFRItems2(items,startLNum,endLNum)
	}
	items, err = lockData.loadCacheL2(db)
	if err == nil {
		sr.appendFRItems2(items,startLNum,endLNum)
	}
	return nil
}
func (s *SnapshotRelease) appendFRItems2(items []*PledgeItem,startLNum uint64,endLNum uint64) {
	for _, item := range items {
		if _, ok := s.Revenue[item.TargetAddress]; !ok {
			s.Revenue[item.TargetAddress] = &LockBalanceData{
				RewardBalance:make(map[uint32]*big.Int),
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		if inLNumScope(item.StartHigh,startLNum,endLNum){
			flowRevenusTarget := s.Revenue[item.TargetAddress]
			if _, ok := flowRevenusTarget.LockBalance[item.StartHigh]; !ok {
				flowRevenusTarget.LockBalance[item.StartHigh] = make(map[uint32]*PledgeItem)
			}
			lockBalance := flowRevenusTarget.LockBalance[item.StartHigh]
			lockBalance[item.PledgeType] = item
		}
	}
}

func (sr *SnapshotRelease) appendFR2(FlowRevenue map[common.Address]*LockBalanceData,startLNum uint64,endLNum uint64) (error) {
	fr1:=FlowRevenue
	for t1, item1 := range fr1 {
		if _, ok := sr.Revenue[t1]; !ok {
			sr.Revenue[t1] = &LockBalanceData{
				RewardBalance:make(map[uint32]*big.Int),
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		rewardBalance:=item1.RewardBalance
		for t2, item2 := range rewardBalance {
			sr.Revenue[t1].RewardBalance[t2]=item2
		}
		lockBalance:=item1.LockBalance
		for t3, item3 := range lockBalance {
			if inLNumScope(t3,startLNum,endLNum){
				if _, ok := sr.Revenue[t1].LockBalance[t3]; !ok {
					sr.Revenue[t1].LockBalance[t3] = make(map[uint32]*PledgeItem)
				}
				t3LockBalance:=sr.Revenue[t1].LockBalance[t3]
				for t4,item4:=range item3{
					if _, ok := t3LockBalance[t4]; !ok {
						t3LockBalance[t4] = item4
					}
				}
			}
		}
	}
	return nil
}

func inLNumScope(num uint64, startLNum uint64, endLNum uint64) bool {
	if num>=startLNum&&num<=endLNum {
		return true
	}
	return false
}

type SnapCanAutoExit struct {
	CandidateAutoExit  []common.Address `json:"candidateautoexit"`
}

func (api *API) GetCandidateAutoExitAtNumber(number uint64) (*SnapCanAutoExit, error) {
	log.Info("api GetCandidateAutoExitAtNumber", "number", number)
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	snapCanAutoExit:=&SnapCanAutoExit{
		CandidateAutoExit:make([]common.Address,0),
	}
	if len(headerExtra.CandidateAutoExit)>0 {
		snapCanAutoExit.CandidateAutoExit=append(snapCanAutoExit.CandidateAutoExit,headerExtra.CandidateAutoExit...)
	}
	return snapCanAutoExit, err
}

func (api *API) GetLockRewardAtNumber(number uint64) ([]LockRewardRecord, error) {
	log.Info("api GetLockRewardAtNumber", "number", number)
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	headerExtra := HeaderExtra{}
	err := rlp.DecodeBytes(header.Extra[extraVanity:len(header.Extra)-extraSeal], &headerExtra)
	if err != nil {
		log.Info("Fail to decode header Extra", "err", err)
		return nil,err
	}
	LockReward:=make([]LockRewardRecord,0)
	if len(headerExtra.LockReward)>0 {
		LockReward=append(LockReward,headerExtra.LockReward...)
	}
	return LockReward, err
}