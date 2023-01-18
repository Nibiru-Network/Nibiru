package alien

import (
	"bytes"
	"github.com/token/common"
	"github.com/token/consensus"
	"github.com/token/core/state"
	"github.com/token/core/types"
	"github.com/token/ethdb"
	"github.com/token/log"
	"github.com/token/rlp"
	"math/big"
	"time"
)

const (
	LOCKREWARDDATA    = "reward"
	LOCKPOFDATA       = "pof"
	LOCKBANDWIDTHDATA = "inspire"
	LOCKPOSEXITDATA     = "posexit"
	LOCKPOFEXITDATA     = "pofexit"
)

type RlsLockData struct {
	LockBalance map[uint64]map[uint32]*PledgeItem // The primary key is lock number, The second key is pledge type
}

type LockData struct {
	Revenue map[common.Address]*LockBalanceData `json:"revenve"`
	CacheL1 []common.Hash                       `json:"cachel1"` // Store chceckout data
	CacheL2 common.Hash                         `json:"cachel2"` //Store data of the previous day
	Locktype string `json:"Locktype"`
}

func NewLockData(t string) *LockData {
	return &LockData{
		Revenue:  make(map[common.Address]*LockBalanceData),
		CacheL1:  []common.Hash{},
		CacheL2:  common.Hash{},
		Locktype: t,
	}
}

func (l *LockData) copy() *LockData {
	clone := &LockData{
		Revenue: make(map[common.Address]*LockBalanceData),
		CacheL1: []common.Hash{},
		CacheL2: l.CacheL2,
		Locktype: l.Locktype,
	}
	clone.CacheL1 = make([]common.Hash, len(l.CacheL1))
	copy(clone.CacheL1, l.CacheL1)
	for who, pledges := range l.Revenue {
		clone.Revenue[who] = &LockBalanceData{
			RewardBalance: make(map[uint32]*big.Int),
			LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
		}
		for which, balance := range l.Revenue[who].RewardBalance {
			clone.Revenue[who].RewardBalance[which] = new(big.Int).Set(balance)
		}
		for when, pledge1 := range pledges.LockBalance {
			clone.Revenue[who].LockBalance[when] = make(map[uint32]*PledgeItem)
			for which, pledge := range pledge1 {
				clone.Revenue[who].LockBalance[when][which] = pledge.copy()
			}
		}
	}
	return clone
}

func (s *LockData) updatePofLockData(snap *Snapshot, item LockRewardRecord, headerNumber *big.Int) {
	if _, ok := s.Revenue[item.Target]; !ok {
		s.Revenue[item.Target] = &LockBalanceData{
			RewardBalance: make(map[uint32]*big.Int),
			LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
		}
	}
	revenusTarget := s.Revenue[item.Target]
	if _, ok := revenusTarget.RewardBalance[item.IsReward]; !ok {
		revenusTarget.RewardBalance[item.IsReward] = new(big.Int).Set(item.Amount)
	} else {
		revenusTarget.RewardBalance[item.IsReward] = new(big.Int).Add(revenusTarget.RewardBalance[item.IsReward], item.Amount)
	}
	deposit := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18))
	if _, ok := snap.SystemConfig.Deposit[item.IsReward]; ok {
		deposit = new(big.Int).Set(snap.SystemConfig.Deposit[item.IsReward])
	}
	if 0 > revenusTarget.RewardBalance[item.IsReward].Cmp(deposit) {
		return
	}
	if _, ok := revenusTarget.LockBalance[headerNumber.Uint64()]; !ok {
		revenusTarget.LockBalance[headerNumber.Uint64()] = make(map[uint32]*PledgeItem)
	}
	lockBalance := revenusTarget.LockBalance[headerNumber.Uint64()]
	// use reward release
	lockPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].LockPeriod
	rlsPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].RlsPeriod
	interval := snap.SystemConfig.LockParameters[sscEnumRwdLock].Interval
	revenueAddress := item.Target
	revenueContract := common.Address{}
	multiSignature := common.Address{}

	// pof or Inspire reward
	if revenue, ok := snap.RevenuePof[item.Target]; ok {
		revenueAddress = revenue.RevenueAddress
		revenueContract = revenue.RevenueContract
		multiSignature = revenue.MultiSignature
	}else {
		if revenue2, ok2 := snap.PofPledge[item.Target]; ok2 {
			revenueAddress = revenue2.Manager
		}
	}

	if _, ok := lockBalance[item.IsReward]; !ok {
		lockBalance[item.IsReward] = &PledgeItem{
			Amount:          big.NewInt(0),
			PledgeType:      item.IsReward,
			Playment:        big.NewInt(0),
			LockPeriod:      lockPeriod,
			RlsPeriod:       rlsPeriod,
			Interval:        interval,
			StartHigh:       headerNumber.Uint64(),
			TargetAddress:   item.Target,
			RevenueAddress:  revenueAddress,
			RevenueContract: revenueContract,
			MultiSignature:  multiSignature,
			BurnAddress:common.Address{},
			BurnRatio:common.Big0,
			BurnAmount:common.Big0,
		}
	}
	lockBalance[item.IsReward].Amount = new(big.Int).Add(lockBalance[item.IsReward].Amount, revenusTarget.RewardBalance[item.IsReward])
	revenusTarget.RewardBalance[item.IsReward] = big.NewInt(0)
}
func (s *LockData) updateAllLockData(snap *Snapshot, isReward uint32, headerNumber *big.Int) {
	distribute:=make(map[common.Address]*big.Int)
	for target, revenusTarget := range s.Revenue {
		if 0 >= revenusTarget.RewardBalance[isReward].Cmp(big.NewInt(0)) {
			continue
		}
		if _, ok := revenusTarget.LockBalance[headerNumber.Uint64()]; !ok {
			revenusTarget.LockBalance[headerNumber.Uint64()] = make(map[uint32]*PledgeItem)
		}
		lockBalance := revenusTarget.LockBalance[headerNumber.Uint64()]
		// use reward release
		lockPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].LockPeriod
		rlsPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].RlsPeriod
		interval := snap.SystemConfig.LockParameters[sscEnumRwdLock].Interval
		revenueAddress := target
		revenueContract := common.Address{}
		multiSignature := common.Address{}
		// singer reward
		if revenue, ok := snap.RevenueNormal[target]; ok {
			revenueAddress = revenue.RevenueAddress
		}else{
			if revenue2, ok2 := snap.PosPledge[target]; ok2 {
				revenueAddress = revenue2.Manager
			}
		}
		if _, ok := lockBalance[isReward]; !ok {
			lockBalance[isReward] = &PledgeItem{
				Amount:          big.NewInt(0),
				PledgeType:      isReward,
				Playment:        big.NewInt(0),
				LockPeriod:      lockPeriod,
				RlsPeriod:       rlsPeriod,
				Interval:        interval,
				StartHigh:       headerNumber.Uint64(),
				TargetAddress:   target,
				RevenueAddress:  revenueAddress,
				RevenueContract: revenueContract,
				MultiSignature:  multiSignature,
				BurnAddress: common.Address{},
				BurnRatio: common.Big0,
				BurnAmount: common.Big0,
			}
		}
		posAmount:=new(big.Int).Set(revenusTarget.RewardBalance[isReward])
		if snap.PosPledge[target]!=nil&& len(snap.PosPledge[target].Detail)>0{
			posRateAmount:=new(big.Int).Mul(posAmount,snap.PosPledge[target].DisRate)
			posRateAmount=new(big.Int).Div(posRateAmount,posDistributionDefaultRate)
			posLeftAmount:=new(big.Int).Sub(posAmount,posRateAmount)
			if posLeftAmount.Cmp(common.Big0)>0{
				if _, ok2 := distribute[target]; ok2 {
					distribute[target]=new(big.Int).Add(distribute[target],posLeftAmount)
				}else{
					distribute[target]=new(big.Int).Set(posLeftAmount)
				}
			}
			if posRateAmount.Cmp(common.Big0)>0{
				lockBalance[isReward].Amount = new(big.Int).Add(lockBalance[isReward].Amount, posRateAmount)
			}
		}else{
			lockBalance[isReward].Amount = new(big.Int).Add(lockBalance[isReward].Amount, posAmount)
		}
		revenusTarget.RewardBalance[isReward] = big.NewInt(0)
	}

	for miner,amount:=range distribute{
		details:=snap.PosPledge[miner].Detail
		totalAmount:=snap.PosPledge[miner].TotalAmount
		for _,item:=range details{
			entrustAmount:=new(big.Int).Mul(amount,item.Amount)
			entrustAmount=new(big.Int).Div(entrustAmount,totalAmount)
			s.updateDistributeLockData(snap,item.Address,miner,entrustAmount,headerNumber)
		}
	}
}

func (s *LockData) updateDistributeLockData(snap *Snapshot, entrustTarget common.Address, revenueContract common.Address,Amount *big.Int,headerNumber *big.Int) {
	if _, ok := s.Revenue[entrustTarget]; !ok {
		s.Revenue[entrustTarget] = &LockBalanceData{
			RewardBalance: make(map[uint32]*big.Int),
			LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
		}
	}
	itemIsReward:=uint32(sscEnumSignerReward)
	revenusTarget := s.Revenue[entrustTarget]
	if _, ok := revenusTarget.RewardBalance[itemIsReward]; !ok {
		revenusTarget.RewardBalance[itemIsReward] = new(big.Int).Set(Amount)
	} else {
		revenusTarget.RewardBalance[itemIsReward] = new(big.Int).Add(revenusTarget.RewardBalance[itemIsReward], Amount)
	}
	if 0 >= revenusTarget.RewardBalance[itemIsReward].Cmp(common.Big0) {
		return
	}
	if _, ok := revenusTarget.LockBalance[headerNumber.Uint64()]; !ok {
		revenusTarget.LockBalance[headerNumber.Uint64()] = make(map[uint32]*PledgeItem)
	}
	lockBalance := revenusTarget.LockBalance[headerNumber.Uint64()]
	// use reward release
	lockPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].LockPeriod
	rlsPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].RlsPeriod
	interval := snap.SystemConfig.LockParameters[sscEnumRwdLock].Interval
	revenueAddress := entrustTarget
	multiSignature := common.Address{}

	if _, ok := lockBalance[itemIsReward]; !ok {
		lockBalance[itemIsReward] = &PledgeItem{
			Amount:          big.NewInt(0),
			PledgeType:      itemIsReward,
			Playment:        big.NewInt(0),
			LockPeriod:      lockPeriod,
			RlsPeriod:       rlsPeriod,
			Interval:        interval,
			StartHigh:       headerNumber.Uint64(),
			TargetAddress:   entrustTarget,
			RevenueAddress:  revenueAddress,
			RevenueContract: revenueContract,
			MultiSignature:  multiSignature,
			BurnAddress: common.Address{},
			BurnRatio: common.Big0,
			BurnAmount: common.Big0,
		}
	}
	lockBalance[itemIsReward].Amount = new(big.Int).Add(lockBalance[itemIsReward].Amount, revenusTarget.RewardBalance[itemIsReward])
	revenusTarget.RewardBalance[itemIsReward] = big.NewInt(0)
}


func (s *LockData) payProfit(hash common.Hash, db ethdb.Database, period uint64, headerNumber uint64, currentGrantProfit []consensus.GrantProfitRecord, playGrantProfit []consensus.GrantProfitRecord, header *types.Header, state *state.StateDB,payAddressAll map[common.Address]*big.Int) ([]consensus.GrantProfitRecord, []consensus.GrantProfitRecord, error) {
	timeNow := time.Now()
	rlsLockBalance := make(map[common.Address]*RlsLockData)
	err := s.saveCacheL1(db, hash)
	if err != nil {
		return currentGrantProfit, playGrantProfit, err
	}
	items, err := s.loadCacheL1(db)
	if err != nil {
		return currentGrantProfit, playGrantProfit, err
	}
	s.appendRlsLockData(rlsLockBalance, items)
	items, err = s.loadCacheL2(db)
	if err != nil {
		return currentGrantProfit, playGrantProfit, err
	}
	s.appendRlsLockData(rlsLockBalance, items)
	log.Info("payProfit load from disk","Locktype",s.Locktype,"len(rlsLockBalance)", len(rlsLockBalance), "elapsed", time.Since(timeNow))

	for address, items := range rlsLockBalance {
		for blockNumber, item1 := range items.LockBalance {
			for which, item := range item1 {
				result, amount := paymentPledge(true, item, state, header,payAddressAll)
				if 0 == result {
					playGrantProfit = append(playGrantProfit, consensus.GrantProfitRecord{
						Which:           which,
						MinerAddress:    address,
						BlockNumber:     blockNumber,
						Amount:          new(big.Int).Set(amount),
						RevenueAddress:  item.RevenueAddress,
						RevenueContract: item.RevenueContract,
						MultiSignature:  item.MultiSignature,
					})
				} else if 1 == result {
					currentGrantProfit = append(currentGrantProfit, consensus.GrantProfitRecord{
						Which:           which,
						MinerAddress:    address,
						BlockNumber:     blockNumber,
						Amount:          new(big.Int).Set(amount),
						RevenueAddress:  item.RevenueAddress,
						RevenueContract: item.RevenueContract,
						MultiSignature:  item.MultiSignature,
					})
				}
			}
		}
	}
	log.Info("payProfit ","Locktype",s.Locktype, "elapsed", time.Since(timeNow))
	return currentGrantProfit, playGrantProfit, nil
}

func (s *LockData) updateGrantProfit(grantProfit []consensus.GrantProfitRecord, db ethdb.Database, hash common.Hash,number uint64) error {

	rlsLockBalance := make(map[common.Address]*RlsLockData)

	items := []*PledgeItem{}
	for _, pledges := range s.Revenue {
		for _, pledge1 := range pledges.LockBalance {
			for _, pledge := range pledge1 {
				items = append(items, pledge)
			}
		}
	}

	s.appendRlsLockData(rlsLockBalance, items)

	items, err := s.loadCacheL1(db)
	if err != nil {
		return err
	}
	s.appendRlsLockData(rlsLockBalance, items)
	items, err = s.loadCacheL2(db)
	if err != nil {
		return err
	}
	s.appendRlsLockData(rlsLockBalance, items)

	hasChanged := false
	for _, item := range grantProfit {
		if 0 != item.BlockNumber {
			if _, ok := rlsLockBalance[item.MinerAddress]; ok {
				if _, ok = rlsLockBalance[item.MinerAddress].LockBalance[item.BlockNumber]; ok {
					if pledge, ok := rlsLockBalance[item.MinerAddress].LockBalance[item.BlockNumber][item.Which]; ok {
						pledge.Playment = new(big.Int).Add(pledge.Playment, item.Amount)
						hasChanged = true
						if 0 <= pledge.Playment.Cmp(pledge.Amount) {
							delete(rlsLockBalance[item.MinerAddress].LockBalance[item.BlockNumber], item.Which)
							if 0 >= len(rlsLockBalance[item.MinerAddress].LockBalance[item.BlockNumber]) {
								delete(rlsLockBalance[item.MinerAddress].LockBalance, item.BlockNumber)
								if 0 >= len(rlsLockBalance[item.MinerAddress].LockBalance) {
									delete(rlsLockBalance, item.MinerAddress)
								}
							}
						}
					}
				}
			}
		}
	}
	if hasChanged {
		s.saveCacheL2(db, rlsLockBalance, hash,number)
	}
	return nil
}

func (s *LockData) loadCacheL1(db ethdb.Database) ([]*PledgeItem, error) {
	result := []*PledgeItem{}
	for _, lv1 := range s.CacheL1 {
		key := append([]byte("alien-"+s.Locktype+"-l1-"), lv1[:]...)
		blob, err := db.Get(key)
		if err != nil {
			return nil, err
		}
		int := bytes.NewBuffer(blob)
		items := []*PledgeItem{}
		err = rlp.Decode(int, &items)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
		log.Info("LockProfitSnap loadCacheL1", "Locktype", s.Locktype, "cache hash", lv1, "size", len(items))
	}
	return result, nil
}

func (s *LockData) appendRlsLockData(rlsLockBalance map[common.Address]*RlsLockData, items []*PledgeItem) {
	for _, item := range items {
		if _, ok := rlsLockBalance[item.TargetAddress]; !ok {
			rlsLockBalance[item.TargetAddress] = &RlsLockData{
				LockBalance: make(map[uint64]map[uint32]*PledgeItem),
			}
		}
		flowRevenusTarget := rlsLockBalance[item.TargetAddress]
		if _, ok := flowRevenusTarget.LockBalance[item.StartHigh]; !ok {
			flowRevenusTarget.LockBalance[item.StartHigh] = make(map[uint32]*PledgeItem)
		}
		lockBalance := flowRevenusTarget.LockBalance[item.StartHigh]
		lockBalance[item.PledgeType] = item
	}
}

func (s *LockData) saveMereCacheL2(db ethdb.Database, rlsLockBalance map[common.Hash]*RlsLockData, hash common.Hash) error {
	items := []*PledgeItem{}
	for _, pledges := range rlsLockBalance {
		for _, pledge1 := range pledges.LockBalance {
			for _, pledge := range pledge1 {
				items = append(items, pledge)
			}
		}
	}
	err, buf := PledgeItemEncodeRlp(items)
	if err != nil {
		return err
	}
	err = db.Put(append([]byte("alien-"+s.Locktype+"-l2-"), hash[:]...), buf)
	if err != nil {
		return err
	}
	for _, pledges := range s.Revenue {
		pledges.LockBalance = make(map[uint64]map[uint32]*PledgeItem)
	}
	s.CacheL1 = []common.Hash{}
	s.CacheL2 = hash
	log.Info("LockProfitSnap saveMereCacheL2", "Locktype", s.Locktype, "cache hash", hash, "len", len(items))
	return nil
}

func (s *LockData) saveCacheL1(db ethdb.Database, hash common.Hash) error {
	items := []*PledgeItem{}
	for _, pledges := range s.Revenue {
		for _, pledge1 := range pledges.LockBalance {
			for _, pledge := range pledge1 {
				items = append(items, pledge)
			}
		}
		pledges.LockBalance = make(map[uint64]map[uint32]*PledgeItem)
	}
	if len(items) == 0 {
		return nil
	}
	err, buf := PledgeItemEncodeRlp(items)
	if err != nil {
		return err
	}
	err = db.Put(append([]byte("alien-"+s.Locktype+"-l1-"), hash[:]...), buf)
	if err != nil {
		return err
	}
	s.CacheL1 = append(s.CacheL1, hash)
	log.Info("LockProfitSnap saveCacheL1", "Locktype", s.Locktype, "cache hash", hash, "len", len(items))
	return nil
}

func (s *LockData) saveCacheL2(db ethdb.Database, rlsLockBalance map[common.Address]*RlsLockData, hash common.Hash,number uint64) error {
	items := []*PledgeItem{}
	for _, pledges := range rlsLockBalance {
		for _, pledge1 := range pledges.LockBalance {
			for _, pledge := range pledge1 {
				items = append(items, pledge)
			}
		}
	}
	err, buf := PledgeItemEncodeRlp(items)
	if err != nil {
		return err
	}
	err = db.Put(append([]byte("alien-"+s.Locktype+"-l2-"), hash[:]...), buf)
	if err != nil {
		return err
	}
	for _, pledges := range s.Revenue {
		pledges.LockBalance = make(map[uint64]map[uint32]*PledgeItem)
	}
	s.CacheL1 = []common.Hash{}
	s.CacheL2 = hash
	log.Info("LockProfitSnap saveCacheL2", "Locktype", s.Locktype, "cache hash", hash, "len", len(items),"number",number)
	return nil
}

func (s *LockData) loadCacheL2(db ethdb.Database) ([]*PledgeItem, error) {
	items := []*PledgeItem{}
	nilHash := common.Hash{}
	if s.CacheL2 == nilHash {
		return items, nil
	}
	key := append([]byte("alien-"+s.Locktype+"-l2-"), s.CacheL2[:]...)
	blob, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	int := bytes.NewBuffer(blob)
	err = rlp.Decode(int, &items)
	if err != nil {
		return nil, err
	}
	log.Info("LockProfitSnap loadCacheL2", "Locktype", s.Locktype, "cache hash", s.CacheL2, "len", len(items))
	return items, nil
}

type LockProfitSnap struct {
	Number        uint64      `json:"number"` // Block number where the snapshot was created
	Hash          common.Hash `json:"hash"`   // Block hash where the snapshot was created
	RewardLock    *LockData   `json:"reward"`
	PofLock       *LockData   `json:"pof"`
	PofInspireLock *LockData   `json:"inspireLock"`
	PosExitLock   *LockData   `json:"posexitlock"`
	PofExitLock   *LockData   `json:"pofexitlock"`
}

func NewLockProfitSnap() *LockProfitSnap {
	return &LockProfitSnap{
		Number:        0,
		Hash:          common.Hash{},
		RewardLock:    NewLockData(LOCKREWARDDATA),
		PofLock:       NewLockData(LOCKPOFDATA),
		PofInspireLock: NewLockData(LOCKBANDWIDTHDATA),
		PosExitLock :NewLockData(LOCKPOSEXITDATA),
		PofExitLock :NewLockData(LOCKPOFEXITDATA),
	}
}
func (s *LockProfitSnap) copy() *LockProfitSnap {
	clone := &LockProfitSnap{
		Number:        s.Number,
		Hash:          s.Hash,
		RewardLock:    s.RewardLock.copy(),
		PofLock:       s.PofLock.copy(),
		PofInspireLock: s.PofInspireLock.copy(),
		PosExitLock:s.PosExitLock.copy(),
		PofExitLock :s.PofExitLock.copy(),
	}
	return clone
}
func (s *LockData) addLockData(snap *Snapshot, item LockRewardRecord, headerNumber *big.Int) {
	if _, ok := s.Revenue[item.Target]; !ok {
		s.Revenue[item.Target] = &LockBalanceData{
			RewardBalance: make(map[uint32]*big.Int),
			LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
		}
	}
	flowRevenusTarget := s.Revenue[item.Target]
	if _, ok := flowRevenusTarget.RewardBalance[item.IsReward]; !ok {
		flowRevenusTarget.RewardBalance[item.IsReward] = new(big.Int).Set(item.Amount)
	} else {
		flowRevenusTarget.RewardBalance[item.IsReward] = new(big.Int).Add(flowRevenusTarget.RewardBalance[item.IsReward], item.Amount)
	}
}

func (s *LockProfitSnap) updateLockData(snap *Snapshot, LockReward []LockRewardRecord, headerNumber *big.Int) {
	blockNumber := headerNumber.Uint64()
	for _, item := range LockReward {
		if sscEnumSignerReward == item.IsReward {
			s.RewardLock.addLockData(snap, item, headerNumber)
		} else if sscEnumPofReward == item.IsReward {
			s.PofLock.updatePofLockData(snap, item, headerNumber)
		} else if sscEnumPofInspireReward == item.IsReward {
			s.PofInspireLock.updatePofLockData(snap, item, headerNumber)
		}
	}
	blockPerDay := snap.getBlockPreDay()
	if 0 == blockNumber%blockPerDay && blockNumber != 0 {
		s.RewardLock.updateAllLockData(snap, sscEnumSignerReward, headerNumber)
	}
}

func (s *LockProfitSnap) payProfit(db ethdb.Database, period uint64, headerNumber uint64, currentGrantProfit []consensus.GrantProfitRecord, playGrantProfit []consensus.GrantProfitRecord, header *types.Header, state *state.StateDB,payAddressAll map[common.Address]*big.Int) ([]consensus.GrantProfitRecord, []consensus.GrantProfitRecord, error) {
	number := header.Number.Uint64()
	if number == 0 {
		return currentGrantProfit, playGrantProfit, nil
	}
	if isPaySignerRewards(number, period) {
		log.Info("LockProfitSnap pay reward profit","number",number)
		return s.RewardLock.payProfit(s.Hash, db, period, headerNumber, currentGrantProfit, playGrantProfit, header, state,payAddressAll)
	}
	if isPayPofRewards(number, period) {
		log.Info("LockProfitSnap pay flow profit","number",number)
		return s.PofLock.payProfit(s.Hash, db, period, headerNumber, currentGrantProfit, playGrantProfit, header, state,payAddressAll)
	}
	if isPayInspireRewards(number, period) {
		log.Info("LockProfitSnap pay inspire profit","number",number)
		return s.PofInspireLock.payProfit(s.Hash, db, period, headerNumber, currentGrantProfit, playGrantProfit, header, state,payAddressAll)
	}
	if isPayPosExit(number, period) {
		log.Info("LockProfitSnap pay POS exit amount","number",number)
		return s.PosExitLock.payProfit(s.Hash, db, period, headerNumber, currentGrantProfit, playGrantProfit, header, state, payAddressAll)
	}
	if isPayPoFExit(number, period) {
		log.Info("LockProfitSnap pay POF exit amount","number",number)
		return s.PofExitLock.payProfit(s.Hash, db, period, headerNumber, currentGrantProfit, playGrantProfit, header, state, payAddressAll)
	}
	return currentGrantProfit, playGrantProfit, nil
}

func (snap *LockProfitSnap) updateGrantProfit(grantProfit []consensus.GrantProfitRecord, db ethdb.Database, headerHash common.Hash, number uint64) {
	shouldUpdateReward, shouldUpdateFlow, shouldUpdateBandwidth,shouldUpdatePosExit,shouldUpdatePofExit := false, false, false,false,false
	for _, item := range grantProfit {
		if 0 != item.BlockNumber {
			if item.Which == sscEnumSignerReward {
				shouldUpdateReward = true
			} else if item.Which == sscEnumPofReward {
				shouldUpdateFlow = true
			} else if item.Which == sscEnumPofInspireReward {
				shouldUpdateBandwidth = true
			}else if item.Which == sscEnumPosExitLock {
				shouldUpdatePosExit = true
			}else if item.Which == sscEnumPofLock {
				shouldUpdatePofExit = true
			}
		}
	}
	storeHash:=headerHash
	if shouldUpdateReward {
		err := snap.RewardLock.updateGrantProfit(grantProfit, db, storeHash,number)
		if err != nil {
			log.Warn("updateGrantProfit Reward Error", "err", err)
		}
	}
	if shouldUpdateFlow {
		err := snap.PofLock.updateGrantProfit(grantProfit, db, storeHash,number)
		if err != nil {
			log.Warn("updateGrantProfit Flow Error", "err", err)
		}
	}
	if shouldUpdateBandwidth {
		err := snap.PofInspireLock.updateGrantProfit(grantProfit, db, storeHash,number)
		if err != nil {
			log.Warn("updateGrantProfit Bandwidth Error", "err", err)
		}
	}
	if shouldUpdatePosExit {
		err := snap.PosExitLock.updateGrantProfit(grantProfit, db, storeHash,number)
		if err != nil {
			log.Warn("updateGrantProfit Pos pledge exit amount Error", "err", err)
		}
	}
	if shouldUpdatePofExit {
		err := snap.PofExitLock.updateGrantProfit(grantProfit, db, storeHash,number)
		if err != nil {
			log.Warn("updateGrantProfit Pof pledge exit amount Error", "err", err)
		}
	}
}


func (snap *LockProfitSnap) saveCacheL1(db ethdb.Database) error {
	err := snap.RewardLock.saveCacheL1(db, snap.Hash)
	if err != nil {
		return err
	}
	err = snap.PofLock.saveCacheL1(db, snap.Hash)
	if err != nil {
		return err
	}
	err = snap.PosExitLock.saveCacheL1(db, snap.Hash)
	if err != nil {
		return err
	}
	err = snap.PofExitLock.saveCacheL1(db, snap.Hash)
	if err != nil {
		return err
	}
	return snap.PofInspireLock.saveCacheL1(db, snap.Hash)
}

func PledgeItemEncodeRlp(items []*PledgeItem) (error, []byte) {
	out := bytes.NewBuffer(make([]byte, 0, 255))
	err := rlp.Encode(out, items)
	if err != nil {
		return err, nil
	}
	return nil, out.Bytes()
}


func (s *LockData) calPayProfit(db ethdb.Database,playGrantProfit []consensus.GrantProfitRecord, header *types.Header) ([]consensus.GrantProfitRecord, error) {
	timeNow := time.Now()

	rlsLockBalance := make(map[common.Address]*RlsLockData)
	items := []*PledgeItem{}
	for _, pledges := range s.Revenue {
		for _, pledge1 := range pledges.LockBalance {
			for _, pledge := range pledge1 {
				items = append(items, pledge)
			}
		}
	}
	s.appendRlsLockData(rlsLockBalance, items)

	items, err := s.loadCacheL1(db)
	if err != nil {
		return playGrantProfit, err
	}
	s.appendRlsLockData(rlsLockBalance, items)
	items, err = s.loadCacheL2(db)
	if err != nil {
		return playGrantProfit, err
	}
	s.appendRlsLockData(rlsLockBalance, items)

	log.Info("calPayProfit load from disk", "Locktype", s.Locktype, "len(rlsLockBalance)", len(rlsLockBalance), "elapsed", time.Since(timeNow), "number", header.Number.Uint64())

	for address, items := range rlsLockBalance {
		for blockNumber, item1 := range items.LockBalance {
			for which, item := range item1 {
				amount := calPaymentPledge( item, header)
				if nil!= amount {
					playGrantProfit = append(playGrantProfit, consensus.GrantProfitRecord{
						Which:           which,
						MinerAddress:    address,
						BlockNumber:     blockNumber,
						Amount:          new(big.Int).Set(amount),
						RevenueAddress:  item.RevenueAddress,
						RevenueContract: item.RevenueContract,
						MultiSignature:  item.MultiSignature,
					})
				}
			}
		}
	}
	log.Info("calPayProfit ", "Locktype", s.Locktype, "elapsed", time.Since(timeNow), "number", header.Number.Uint64())
	return playGrantProfit, nil
}

func (s *LockData) loadRlsLockBalance(db ethdb.Database) (map[common.Address]*RlsLockData , error) {
	rlsLockBalance := make(map[common.Address]*RlsLockData)

	items := []*PledgeItem{}
	for _, pledges := range s.Revenue {
		for _, pledge1 := range pledges.LockBalance {
			for _, pledge2 := range pledge1 {
				items = append(items, pledge2)
			}
		}
	}

	s.appendRlsLockData(rlsLockBalance, items)

	items, err := s.loadCacheL1(db)
	if err != nil {
		return rlsLockBalance,err
	}
	s.appendRlsLockData(rlsLockBalance, items)
	items, err = s.loadCacheL2(db)
	if err != nil {
		return rlsLockBalance,err
	}
	s.appendRlsLockData(rlsLockBalance, items)
	return rlsLockBalance,nil
}

func(s *LockData) setBurnRatio(lock *PledgeItem,burnRatio *big.Int) {
	if burnRatio.Cmp(common.Big0)>0{
		if lock.BurnRatio==nil{
			lock.BurnAddress=common.BigToAddress(big.NewInt(0))
			lock.BurnRatio=burnRatio
		}else if lock.BurnRatio.Cmp(burnRatio)<0{
			lock.BurnRatio=burnRatio
		}
		if lock.BurnAmount==nil{
			lock.BurnAmount=common.Big0
		}
	}
}

func (s *LockData) updateLockPofExitData(snap *Snapshot, item LockRewardRecord, headerNumber *big.Int) {
	if _, ok := s.Revenue[item.Target]; !ok {
		s.Revenue[item.Target] = &LockBalanceData{
			RewardBalance: make(map[uint32]*big.Int),
			LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
		}
	}
	pofPledge:=snap.PofPledge[item.Target]
	if pofPledge == nil {
		log.Warn("updateLockPofExitData","nof pofPledge","miner",item.Target)
		return
	}
	revenusTarget := s.Revenue[item.Target]
	if _, ok := revenusTarget.RewardBalance[item.IsReward]; !ok {
		revenusTarget.RewardBalance[item.IsReward] = new(big.Int).Set(item.Amount)
	} else {
		revenusTarget.RewardBalance[item.IsReward] = new(big.Int).Add(revenusTarget.RewardBalance[item.IsReward], item.Amount)
	}
	if _, ok := revenusTarget.LockBalance[headerNumber.Uint64()]; !ok {
		revenusTarget.LockBalance[headerNumber.Uint64()] = make(map[uint32]*PledgeItem)
	}
	lockBalance := revenusTarget.LockBalance[headerNumber.Uint64()]
	// use reward release
	lockPeriod := snap.SystemConfig.LockParameters[sscEnumPofLock].LockPeriod
	rlsPeriod := snap.SystemConfig.LockParameters[sscEnumPofLock].RlsPeriod
	interval := snap.SystemConfig.LockParameters[sscEnumPofLock].Interval
	revenueAddress := snap.PofPledge[item.Target].Manager
	revenueContract := common.Address{}
	multiSignature := common.Address{}
	if _, ok := lockBalance[item.IsReward]; !ok {
		lockBalance[item.IsReward] = &PledgeItem{
			Amount:          big.NewInt(0),
			PledgeType:      item.IsReward,
			Playment:        big.NewInt(0),
			LockPeriod:      lockPeriod,
			RlsPeriod:       rlsPeriod,
			Interval:        interval,
			StartHigh:       headerNumber.Uint64(),
			TargetAddress:   item.Target,
			RevenueAddress:  revenueAddress,
			RevenueContract: revenueContract,
			MultiSignature:  multiSignature,
			BurnAddress:common.Address{},
			BurnRatio:common.Big0,
			BurnAmount:common.Big0,
		}
	}
	lockBalance[item.IsReward].Amount = new(big.Int).Add(lockBalance[item.IsReward].Amount, revenusTarget.RewardBalance[item.IsReward])
	revenusTarget.RewardBalance[item.IsReward] = big.NewInt(0)
}