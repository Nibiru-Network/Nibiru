package alien

import (
	"github.com/token/common"
	"github.com/token/common/hexutil"
	"github.com/token/consensus"
	"github.com/token/core/state"
	"github.com/token/core/types"
	"github.com/token/ethdb"
	"github.com/token/log"
	"math/big"
)

const (
	categoryCandEntrust   = "CandEntrust"
	categoryCandEntrustExit = "CandETExit"
	categoryCandChangeRate = "CandChaRate"
	categoryCandChangeManager = "CandChaMan"

	sscEnumPosCommitPeriod=12
	sscEnumPosBeyondCommitPeriod=13
	sscEnumPosWithinCommitPeriod=14

	sscEnumPosExitLock = 7

)

var (
	minCndEntrustPledgeBalance = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(1))
	maxPosContinueDayFail      =uint64(30)
	posDistributionDefaultRate =big.NewInt(10000)
	BurnBase=big.NewInt(10000)
	posCommitPeriod=big.NewInt(365) //1 year
	posBeyondCommitPeriod=big.NewInt(7) //7 day
	posWithinCommitPeriod=big.NewInt(7) //7 day
	posCandidateOneNum=9
	posCandidateTwoNum=12
	clearSignNumberPerid = uint64(60480)  //7day
	posCandidateAvgRate  = big.NewInt(70) //70%
)
type CandidatePledgeNewRecord struct {
	Target common.Address
	Amount *big.Int
	Manager common.Address
	Hash common.Hash
}

type CandidatePledgeEntrustRecord struct {
	Target common.Address
	Amount *big.Int
	Address common.Address
	Hash common.Hash
}
type CandidatePEntrustExitRecord struct {
	Target common.Address
	Hash common.Hash
	Address common.Address
	Amount *big.Int
}
type CandidateChangeRateRecord struct {
	Target common.Address
	Rate *big.Int
}

type PosPledgeItem struct {
	Manager     common.Address                `json:"manager"`
	Active      uint64                        `json:"active"`
	TotalAmount *big.Int                      `json:"totalamount"`
	Detail      map[common.Hash]*PledgeDetail `json:"detail"`
	LastPunish  uint64                        `json:"lastpunish"`
	DisRate     *big.Int                      `json:"distributerate"`
}

type PledgeDetail struct {
	Address common.Address `json:"address"`
	Height  uint64         `json:"height"`
	Amount  *big.Int       `json:"amount"`
}

type CandidateChangeManagerRecord struct {
	Target common.Address
	Manager common.Address
}

func (a *Alien) processPosCustomTx(txDataInfo []string, headerExtra HeaderExtra, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, snap *Snapshot, number *big.Int, state*state.StateDB, chain consensus.ChainHeaderReader, coinBalances map[common.Address]*big.Int) HeaderExtra {
	if txDataInfo[posCategory] == tokenCategoryCandReq {
		headerExtra.CandidatePledgeNew= a.processCandidatePledgeNew(headerExtra.CandidatePledgeNew,txDataInfo,txSender,tx,receipts,state,snap,number.Uint64())
	}
	if txDataInfo[posCategory] == tokenCategoryCandExit {
		headerExtra.CandidatePEntrustExit,headerExtra.CandidateExit = a.processCandidateExitNew (headerExtra.CandidatePEntrustExit,headerExtra.CandidateExit, txDataInfo, txSender, tx, receipts, state, snap,number.Uint64())
	}
	if txDataInfo[posCategory] == categoryCandEntrust {
		headerExtra.CandidatePledgeEntrust = a.processCandidatePledgeEntrust(headerExtra.CandidatePledgeEntrust, txDataInfo, txSender, tx, receipts, state, snap, number.Uint64())
	}
	if txDataInfo[posCategory] == categoryCandEntrustExit {
		headerExtra.CandidatePEntrustExit= a.processCandidatePEntrustExit(headerExtra.CandidatePEntrustExit,txDataInfo,txSender,tx,receipts,state,snap,number.Uint64())
	}
	if txDataInfo[posCategory] == categoryCandChangeRate {
		headerExtra.CandidateChangeRate= a.processCandidateChangeRate(headerExtra.CandidateChangeRate,txDataInfo,txSender,tx,receipts,state,snap,number.Uint64())
	}
	if txDataInfo[posCategory] == categoryCandChangeManager {
		if isGEPosChangeManagerNumber(number.Uint64()) {
			headerExtra.CandidateChangeManager = a.processCandidateChangeManager(headerExtra.CandidateChangeManager, txDataInfo, txSender, tx, receipts, state, snap, number.Uint64())
		}
	}


	return headerExtra
}

func (a *Alien) processCandidatePledgeEntrust(currentCandidatePledge []CandidatePledgeEntrustRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot, number uint64) []CandidatePledgeEntrustRecord {

	if len(txDataInfo) <= 4 {
		log.Warn("Candidate Entrust", "parameter number", len(txDataInfo))
		return currentCandidatePledge
	}
	candidatePledge := CandidatePledgeEntrustRecord{
		Target: common.Address{},
		Amount: new(big.Int).Set(minCndEntrustPledgeBalance),
		Address: txSender,
		Hash: tx.Hash(),
	}
	postion := 3
	if err := candidatePledge.Target.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate Entrust", "miner address", txDataInfo[postion])
		return currentCandidatePledge
	}

	if _, ok := snap.PosPledge[candidatePledge.Target]; !ok {
		log.Warn("Candidate Entrust", "candidate is not exist", candidatePledge.Target)
		return currentCandidatePledge
	}

	if _, ok := snap.PosPledge[candidatePledge.Address]; ok {
		log.Warn("Candidate Entrust", "txSender is miner address", candidatePledge.Address)
		return currentCandidatePledge
	}
	postion++
	var err error
	if candidatePledge.Amount, err = hexutil.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate Entrust", "number", txDataInfo[postion])
		return currentCandidatePledge
	}
	if candidatePledge.Amount.Cmp(minCndEntrustPledgeBalance)<0{
		log.Warn("Candidate Entrust", "Amount less than 1 ", txDataInfo[postion])
		return currentCandidatePledge
	}
	targetMiner:=snap.findPosTargetMiner(txSender)
	nilAddr := common.Address{}
	if targetMiner!=nilAddr&&targetMiner!=candidatePledge.Target{
		log.Warn("Candidate Entrust", "one address can only pledge one miner ", targetMiner)
		return currentCandidatePledge
	}

	if isGEPosChangeManagerNumber(number){
		if snap.isPosOtherMinerManager(txSender,candidatePledge.Target){
			log.Warn("Candidate Entrust", "txSender is other miner manager ", txSender,"target",candidatePledge.Target)
			return currentCandidatePledge
		}
	}

	if state.GetBalance(txSender).Cmp(candidatePledge.Amount) < 0 {
		log.Warn("Candidate Entrust", "balance", state.GetBalance(txSender))
		return currentCandidatePledge
	}
	state.SubBalance(txSender, candidatePledge.Amount)
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xdcadcdae40a91d6ed79cf78187b18f2d3b9c49f7ff68799d06850a8d35b2fd7e")) //web3.sha3("PledgeEntrust(address,uint256)")
	topics[1].SetBytes(candidatePledge.Target.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumCndLock).Bytes())
	data := common.Hash{}
	data.SetBytes(candidatePledge.Amount.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, data.Bytes())
	currentCandidatePledge = append(currentCandidatePledge, candidatePledge)
	return currentCandidatePledge
}
func (a *Alien) processCandidatePEntrustExit(currentCandidatePEntrustExit []CandidatePEntrustExitRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot, number uint64) []CandidatePEntrustExitRecord {

	if len(txDataInfo) <= 4 {
		log.Warn("Candidate PEntrustExit", "parameter number", len(txDataInfo))
		return currentCandidatePEntrustExit
	}
	candidatePledge := CandidatePEntrustExitRecord{
		Target: common.Address{},
		Hash: common.Hash{},
		Address:common.Address{},
		Amount:common.Big0,
	}
	postion := 3
	if err := candidatePledge.Target.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate PEntrustExit", "miner address", txDataInfo[postion])
		return currentCandidatePEntrustExit
	}
	postion++
	candidatePledge.Hash = common.HexToHash(txDataInfo[postion])

	if _, ok := snap.PosPledge[candidatePledge.Target]; !ok {
		log.Warn("Candidate PEntrustExit", "candidate is not exist", candidatePledge.Target)
		return currentCandidatePEntrustExit
	}

	if _, ok := snap.PosPledge[candidatePledge.Target].Detail[candidatePledge.Hash]; !ok {
		log.Warn("Candidate PEntrustExit", "Hash is not exist", candidatePledge.Hash)
		return currentCandidatePEntrustExit
	}else{
		pledgeDetail:=snap.PosPledge[candidatePledge.Target].Detail[candidatePledge.Hash]
		if pledgeDetail.Address!=txSender {
			log.Warn("Candidate PEntrustExit", "txSender is not right", txSender)
			return currentCandidatePEntrustExit
		}
		candidatePledge.Address=pledgeDetail.Address
		candidatePledge.Amount=pledgeDetail.Amount
	}

	if isInCurrentCandidatePEntrustExit(currentCandidatePEntrustExit,candidatePledge.Hash){
		log.Warn("Candidate PEntrustExit", "Hash is in currentCandidatePEntrustExit", candidatePledge.Hash)
		return currentCandidatePEntrustExit
	}

	if snap.isInPosCommitPeriod(candidatePledge.Target,number){
		if txSender==snap.PosPledge[candidatePledge.Target].Manager {
			log.Warn("Candidate exit New", "minerAddress is in commit period", candidatePledge.Target)
			return currentCandidatePEntrustExit
		}else{
			if snap.isInPosCommitPeriodPass(candidatePledge.Target,number,candidatePledge.Hash,snap.SystemConfig.Deposit[sscEnumPosWithinCommitPeriod].Uint64()){
				log.Warn("Candidate exit New", "hash is not BeyondCommitPeriod", candidatePledge.Hash)
				return currentCandidatePEntrustExit
			}
		}
	}else {
		if snap.isInPosCommitPeriodPass(candidatePledge.Target,number,candidatePledge.Hash,snap.SystemConfig.Deposit[sscEnumPosBeyondCommitPeriod].Uint64()){
			log.Warn("Candidate exit New", "hash is not BeyondCommitPeriod", candidatePledge.Hash)
			return currentCandidatePEntrustExit
		}
	}

	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xaf7cf8bf073a87df843f342229f11fc2ecc069751926bc402836a4f1f2a52403")) //web3.sha3("PledgeEntrustExit(address,bytes32)")
	topics[1].SetBytes(candidatePledge.Target.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumCndLock).Bytes())
	data := common.Hash{}
	data.SetBytes(candidatePledge.Hash.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, data.Bytes())
	currentCandidatePEntrustExit = append(currentCandidatePEntrustExit, candidatePledge)

	if txSender==snap.PosPledge[candidatePledge.Target].Manager && !snap.isInTally(candidatePledge.Target){
		managerPledge:=make([]common.Hash,0)
		for hash,item:=range snap.PosPledge[candidatePledge.Target].Detail{
			if item.Address==txSender{
				managerPledge=append(managerPledge,hash)
			}
		}
		exitPledge:=make([]common.Hash,0)
		for _,item:=range currentCandidatePEntrustExit{
			address:=snap.PosPledge[candidatePledge.Target].Detail[item.Hash].Address
			if address==txSender{
				exitPledge=append(exitPledge,item.Hash)
			}
		}
		if len(exitPledge)>0&&len(exitPledge)==len(managerPledge){
			for hash,item:=range snap.PosPledge[candidatePledge.Target].Detail{
				if !isInCurrentCandidatePEntrustExit(currentCandidatePEntrustExit,hash){
					candidateExitPledge := CandidatePEntrustExitRecord{
						Target: candidatePledge.Target,
						Hash: hash,
						Address:item.Address,
						Amount:item.Amount,
					}
					currentCandidatePEntrustExit = append(currentCandidatePEntrustExit, candidateExitPledge)
				}
			}
		}
	}
	return currentCandidatePEntrustExit
}

func (a *Alien) processCandidateExitNew(currentCandidatePEntrustExit []CandidatePEntrustExitRecord, currentCandidateExit []common.Address, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot, number uint64) ([]CandidatePEntrustExitRecord, []common.Address) {
	if len(txDataInfo) <= tokenPosMinerAddress {
		log.Warn("Candidate exit New", "parameter number", len(txDataInfo))
		return currentCandidatePEntrustExit, currentCandidateExit
	}
	minerAddress := common.Address{}
	if err := minerAddress.UnmarshalText1([]byte(txDataInfo[tokenPosMinerAddress])); err != nil {
		log.Warn("Candidate exit New", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentCandidatePEntrustExit, currentCandidateExit
	}
	if oldBind, ok := snap.PosPledge[minerAddress]; ok {
		if oldBind.Manager !=txSender &&!(snap.isSystemManagerAndInTally(txSender,minerAddress)){
			log.Warn("Candidate exit New", "Manager address is not txSender", txSender)
			return currentCandidatePEntrustExit, currentCandidateExit
		}
	}else{
		log.Warn("Candidate exit New", "minerAddress is not exist", minerAddress)
		return currentCandidatePEntrustExit, currentCandidateExit
	}

	if snap.isInPosCommitPeriod(minerAddress,number){
		log.Warn("Candidate exit New", "minerAddress is in commit period", minerAddress)
		return currentCandidatePEntrustExit, currentCandidateExit
	}

	if snap.isInTally(minerAddress)&&!snap.isSystemManager(txSender){
		log.Warn("Candidate exit New", "minerAddress is in tally", minerAddress)
		return currentCandidatePEntrustExit, currentCandidateExit
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x9489b96ebcb056332b79de467a2645c56a999089b730c99fead37b20420d58e7")) //web3.sha3("PledgeExit(address)")
	topics[1].SetBytes(minerAddress.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumCndLock).Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	for hash,item:=range snap.PosPledge[minerAddress].Detail{
		candidateExitPledge := CandidatePEntrustExitRecord{
			Target: minerAddress,
			Hash: hash,
			Address:item.Address,
			Amount:item.Amount,
		}
		currentCandidatePEntrustExit = append(currentCandidatePEntrustExit, candidateExitPledge)
	}
	if snap.isSystemManagerAndInTally(txSender,minerAddress){
		currentCandidateExit=append(currentCandidateExit,minerAddress)
	}
	return currentCandidatePEntrustExit, currentCandidateExit
}

func (a *Alien) isPosManager(snap *Snapshot,deviceBind DeviceBindRecord,txSender common.Address,txDataInfo []string) bool {
	if _, ok := snap.PosPledge[deviceBind.Device]; ok {
		if snap.PosPledge[deviceBind.Device].Manager != txSender {
			log.Warn("isPosManager", "txSender is not manager", txSender)
			return false
		}else{
			return true
		}
	} else {
		log.Warn("isPosManager", "PosPledge is not exist", txDataInfo[tokenPosMinerAddress])
		return false
	}
}
func (a *Alien) isPofManager(snap *Snapshot,deviceBind DeviceBindRecord,txSender common.Address,txDataInfo []string) bool {
	if pofPledge,ok :=snap.PofPledge[deviceBind.Device];ok{
		if pofPledge.Manager != txSender {
			log.Warn("isPofManager", "txSender no role",txSender,"miner", txDataInfo[tokenPosMinerAddress])
			return false
		}
	}else{
		log.Warn("isPofManager", "pof miner not exit", txDataInfo[tokenPosMinerAddress])
		return false
	}
	return true
}
func (a *Alien) processCandidateChangeRate(currentCandidateRate []CandidateChangeRateRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot, number uint64) []CandidateChangeRateRecord {
	if len(txDataInfo) <= 4 {
		log.Warn("Candidate ChangeRate", "parameter number", len(txDataInfo))
		return currentCandidateRate
	}
	postion := 3
	minerAddress := common.Address{}
	if err := minerAddress.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate ChangeRate", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentCandidateRate
	}
	if oldBind, ok := snap.PosPledge[minerAddress]; ok {
		if oldBind.Manager !=txSender {
			log.Warn("Candidate ChangeRate", "Manager address is not txSender", txSender)
			return currentCandidateRate
		}
	}else{
		log.Warn("Candidate ChangeRate", "minerAddress is not exist", minerAddress)
		return currentCandidateRate
	}
	candidateChangeRate := CandidateChangeRateRecord{
		Target: minerAddress,
		Rate:common.Big0,
	}
	postion++
	var err error
	if candidateChangeRate.Rate, err = hexutil.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate ChangeRate", "number", txDataInfo[postion])
		return currentCandidateRate
	}
	if candidateChangeRate.Rate.Cmp(posDistributionDefaultRate)>0 {
		log.Warn("Candidate ChangeRate", "Rate greater than posDistributionDefaultRate ", txDataInfo[postion])
		return currentCandidateRate
	}
	if candidateChangeRate.Rate.Cmp(common.Big0)<=0 {
		log.Warn("Candidate ChangeRate", "Rate Less than or equal to 0 ", txDataInfo[postion])
		return currentCandidateRate
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x4c3b40c94758c0e27ceddbf9149c59c96f3694815f7b7dcb267fd8db56762bcf")) //web3.sha3("PledgeChaRate(address,uint256)")
	topics[1].SetBytes(minerAddress.Bytes())
	topics[2].SetBytes(candidateChangeRate.Rate.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentCandidateRate=append(currentCandidateRate,candidateChangeRate)
	return currentCandidateRate
}

func isInCurrentCandidatePEntrustExit(currentCandidatePEntrustExit []CandidatePEntrustExitRecord, hash common.Hash) bool {
	has:=false
	for _,currentItem:=range currentCandidatePEntrustExit{
		if currentItem.Hash==hash {
			has=true
			break
		}
	}
	return has
}

func (snap *Snapshot) updateCandidatePledgeNew(candidatePledge []CandidatePledgeNewRecord, number uint64) {
	for _, item := range candidatePledge {
		if _, ok := snap.PosPledge[item.Target]; !ok {
			pledgeItem := &PosPledgeItem{
				Manager:item.Manager,
				Active:number,
				TotalAmount:new(big.Int).Set(item.Amount),
				Detail:make(map[common.Hash]*PledgeDetail,0),
				LastPunish: uint64(0),
				DisRate: new(big.Int).Set(posDistributionDefaultRate),
			}
			pledgeItem.Detail[item.Hash]=&PledgeDetail{
				Address:item.Manager,
				Height:number,
				Amount:new(big.Int).Set(item.Amount),
			}
			snap.PosPledge[item.Target] = pledgeItem
		}
		if _, ok := snap.TallyMiner[item.Target]; !ok {
			snap.TallyMiner[item.Target] = &CandidateState{
				SignerNumber: 0,
				Stake:        big.NewInt(0),
			}
		}
	}
}

func (snap *Snapshot) updateCandidatePledgeEntrust(candidatePledge []CandidatePledgeEntrustRecord, number uint64) {
	for _, item := range candidatePledge {
		if _, ok := snap.PosPledge[item.Target]; ok {
			snap.PosPledge[item.Target].Detail[item.Hash]=&PledgeDetail{
				Address:item.Address,
				Height:number,
				Amount:item.Amount,
			}
			snap.PosPledge[item.Target].TotalAmount=new(big.Int).Add(snap.PosPledge[item.Target].TotalAmount,item.Amount)
		}
	}
}

func (snap *Snapshot) updateCandidatePEntrustExit(candidatePledge []CandidatePEntrustExitRecord, headerNumber *big.Int) {
	for _, item := range candidatePledge {
		if _, ok := snap.PosPledge[item.Target]; ok {
			if _, ok2 := snap.PosPledge[item.Target].Detail[item.Hash]; ok2 {
				snap.PosPledge[item.Target].TotalAmount=new(big.Int).Sub(snap.PosPledge[item.Target].TotalAmount,item.Amount)
				delete(snap.PosPledge[item.Target].Detail,item.Hash)
				snap.Revenue.PosExitLock.updatePosExitLockData(snap, item, headerNumber)
			}
			if !snap.isInTally(item.Target)&&snap.PosPledge[item.Target].TotalAmount.Cmp(common.Big0)<=0{
				snap.removePosPledge(item.Target)
			}
		}
	}
}
func (snap *Snapshot) checkCandidateAutoExit(number uint64, candidateAutoExit []common.Address, state *state.StateDB, candidatePEntrustExit []CandidatePEntrustExitRecord) ([]common.Address, []CandidatePEntrustExitRecord) {
	burnAmount:=common.Big0
	if isCheckPOSAutoExit(number, snap.Period) {
		for miner,item:=range snap.PosPledge{
			if item.LastPunish>0&&(number-item.LastPunish)>=maxPosContinueDayFail*snap.getBlockPreDay(){
				candidateAutoExit=append(candidateAutoExit,miner)
				for hash,detail:=range snap.PosPledge[miner].Detail{
					if detail.Address==item.Manager{
						burnAmount=new(big.Int).Add(burnAmount,detail.Amount)
					}else{
						candidatePEntrustExit=append(candidatePEntrustExit,CandidatePEntrustExitRecord{
							Target:miner,
							Hash:hash,
							Address: detail.Address,
							Amount:new (big.Int).Set(detail.Amount),
						})
					}
				}

			}
		}
	}
	if burnAmount.Cmp(common.Big0)>0{
		state.AddBalance(common.BigToAddress(big.NewInt(0)),burnAmount)
	}
	return candidateAutoExit, candidatePEntrustExit
}
func (snap *Snapshot) updateCandidateAutoExit(candidateAutoExit []common.Address, header *types.Header, db ethdb.Database) {
	if candidateAutoExit==nil ||len(candidateAutoExit)==0 {
		return
	}
	for _, miner := range candidateAutoExit {
		snap.removePosPledge(miner)
		snap.removeTally(miner)
	}
	err:= snap.Revenue.RewardLock.setRewardRemovePunish(candidateAutoExit, db, header.Hash(),header.Number.Uint64())
	if err != nil {
		log.Warn("setRewardRemovePunish RewardLock Error", "err", err)
	}
}

func (snap *Snapshot) isInPosCommitPeriod(minerAddress common.Address, number uint64) bool {
	if (number-snap.PosPledge[minerAddress].Active)<=(snap.SystemConfig.Deposit[sscEnumPosCommitPeriod].Uint64()*snap.getBlockPreDay()) {
		return true
	}
	return false
}

func (snap *Snapshot) isInPosCommitPeriodPass(minerAddress common.Address, number uint64, hash common.Hash, setting uint64) bool {
	pledgeDetail:=snap.PosPledge[minerAddress].Detail[hash]
	if (number-pledgeDetail.Height)<=(setting*snap.getBlockPreDay()){
		return true
	}
	return false
}

func (snap *Snapshot) findPosTargetMiner(txSender common.Address) common.Address {
	for miner,item:=range snap.PosPledge{
		details:=item.Detail
		for _,detail:=range details{
			if detail.Address==txSender{
				return miner
			}
		}
	}
	return common.Address{}
}

func (snap *Snapshot) updateCandidateChangeRate(candidateChangeRate []CandidateChangeRateRecord, header *types.Header, db ethdb.Database) {
	for _, item := range candidateChangeRate {
		if _, ok := snap.PosPledge[item.Target]; ok {
			snap.PosPledge[item.Target].DisRate = new(big.Int).Set(item.Rate)
		}
	}
}


func (s *Snapshot) isInTally(minerAddress common.Address) bool {
	if _,ok:=s.Tally[minerAddress];ok {
		return true
	}
	return false
}

func (snap *Snapshot) removeTally(miner common.Address) {
	if _, ok := snap.Tally[miner]; ok {
		delete(snap.Tally, miner)
	}
	if _, ok := snap.Votes[miner]; ok {
		delete(snap.Votes, miner)
	}
	if _, ok := snap.Voters[miner]; ok {
		delete(snap.Voters, miner)
	}
	if _, ok := snap.Candidates[miner]; ok {
		delete(snap.Candidates, miner)
	}
}

func (snap *Snapshot) isSystemManager(txSender common.Address) bool {
	return snap.SystemConfig.ManagerAddress[sscEnumSystem] == txSender
}

func (snap *Snapshot) isSystemManagerAndInTally(txSender common.Address,minerAddress common.Address) bool {
	return snap.isSystemManager(txSender)&&snap.isInTally(minerAddress)
}

func (snap *Snapshot) updateCandidateExit2(candidateExit []common.Address, number *big.Int) {
	if candidateExit==nil || len(candidateExit)==0 {
		return
	}
	for _, miner := range candidateExit {
		snap.removePosPledge(miner)
		snap.removeTally(miner)
	}
}

func (snap *Snapshot) removePosPledge(miner common.Address) {
	if _, ok := snap.PosPledge[miner]; ok {
		delete(snap.PosPledge, miner)
	}
	if _, ok := snap.RevenueNormal[miner]; ok {
		delete(snap.RevenueNormal, miner)
	}
	if _, ok := snap.TallyMiner[miner]; ok {
		delete(snap.TallyMiner, miner)
	}
}

func (s *LockData) updatePosExitLockData(snap *Snapshot, item CandidatePEntrustExitRecord, headerNumber *big.Int) {
	if _, ok := s.Revenue[item.Address]; !ok {
		s.Revenue[item.Address] = &LockBalanceData{
			RewardBalance: make(map[uint32]*big.Int),
			LockBalance:   make(map[uint64]map[uint32]*PledgeItem),
		}
	}
	itemIsReward:=uint32(sscEnumPosExitLock)
	flowRevenusTarget := s.Revenue[item.Address]
	if _, ok := flowRevenusTarget.RewardBalance[itemIsReward]; !ok {
		flowRevenusTarget.RewardBalance[itemIsReward] = new(big.Int).Set(item.Amount)
	} else {
		flowRevenusTarget.RewardBalance[itemIsReward] = new(big.Int).Add(flowRevenusTarget.RewardBalance[itemIsReward], item.Amount)
	}
	if _, ok := flowRevenusTarget.LockBalance[headerNumber.Uint64()]; !ok {
		flowRevenusTarget.LockBalance[headerNumber.Uint64()] = make(map[uint32]*PledgeItem)
	}
	lockBalance := flowRevenusTarget.LockBalance[headerNumber.Uint64()]
	lockPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].LockPeriod
	rlsPeriod := snap.SystemConfig.LockParameters[sscEnumRwdLock].RlsPeriod
	interval := snap.SystemConfig.LockParameters[sscEnumRwdLock].Interval
	revenueAddress := item.Address
	revenueContract := item.Target
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
			TargetAddress:   item.Address,
			RevenueAddress:  revenueAddress,
			RevenueContract: revenueContract,
			MultiSignature:  multiSignature,
			BurnAddress: common.Address{},
			BurnRatio: common.Big0,
			BurnAmount: common.Big0,
		}
	}
	lockBalance[itemIsReward].Amount = new(big.Int).Add(lockBalance[itemIsReward].Amount, flowRevenusTarget.RewardBalance[itemIsReward])
	flowRevenusTarget.RewardBalance[itemIsReward] = big.NewInt(0)
}


func (s *LockData) setRewardRemovePunish(pledge []common.Address, db ethdb.Database, hash common.Hash, number uint64) error {
	rlsLockBalance,err:=s.loadRlsLockBalance(db)
	if err != nil {
		return err
	}
	pledgeAddrs := make(map[common.Address]uint64)
	for _, sPAddrs := range pledge {
		pledgeAddrs[sPAddrs] = 1
	}
	hasChanged := false
	burnRatio:=new(big.Int).Set(BurnBase)
	for minerAddress,itemRlsLock:=range rlsLockBalance{
		lockBalance:=itemRlsLock.LockBalance
		if _, ok := pledgeAddrs[minerAddress]; ok {
			hasChanged=true
			for _,itemBlockLock:=range lockBalance{
				for _,itemWhichLock:=range itemBlockLock{
					s.setBurnRatio(itemWhichLock,burnRatio)
				}
			}
		}
	}
	if hasChanged{
		s.saveCacheL2(db, rlsLockBalance, hash,number)
	}
	return nil
}


func (snap *Snapshot) posApply(headerExtra HeaderExtra, header *types.Header, db ethdb.Database) (*Snapshot, error) {
	if header.Number.Uint64()==1 {
		snap.initPosPledge(header.Number.Uint64())
	}
	snap.updateCandidatePledgeNew(headerExtra.CandidatePledgeNew, header.Number.Uint64())
	snap.updateCandidatePledgeEntrust(headerExtra.CandidatePledgeEntrust, header.Number.Uint64())
	snap.updateCandidatePEntrustExit(headerExtra.CandidatePEntrustExit, header.Number)
	snap.updateCandidateAutoExit(headerExtra.CandidateAutoExit, header,db)
	snap.updateCandidateChangeRate(headerExtra.CandidateChangeRate, header,db)
	if header.Number.Uint64()%(snap.config.MaxSignerCount*snap.LCRS) == 0 {
		snap.updateSignerNumber(headerExtra.SignerQueue,header.Number.Uint64())
	}
	if len(headerExtra.MinerStake) >0 || len(headerExtra.ModifyPredecessorVotes)  >0 {
		snap.deletePunishByPosExit(header.Number.Uint64())
	}
	if isGEPosChangeManagerNumber(header.Number.Uint64()){
		snap.updateCandidateChangeManager(headerExtra.CandidateChangeManager, header.Number)
	}
	snap.updateCandidateExit2(headerExtra.CandidateExit, header.Number)
	return snap, nil
}

func (s *Snapshot) deletePunishByPosExit(headerNumber uint64){
	for punishAddr,_:= range s.Punished {
		if _, ok := s.PosPledge[punishAddr];!ok {
			delete(s.Punished,punishAddr)
		}
	}
}

func (a *Alien) processCandidatePledgeNew(currentCandidatePledge []CandidatePledgeNewRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot, number uint64) []CandidatePledgeNewRecord {
	if len(txDataInfo) <= tokenPosMinerAddress {
		log.Warn("Candidate pledgeNew", "parameter number", len(txDataInfo))
		return currentCandidatePledge
	}
	candidatePledge := CandidatePledgeNewRecord{
		Target: common.Address{},
		Amount: new(big.Int).Set(minCndPledgeBalance),
		Manager: txSender,
		Hash: tx.Hash(),
	}
	if deposit, ok := snap.SystemConfig.Deposit[0]; ok {
		candidatePledge.Amount = new(big.Int).Set(deposit)
	}
	if err := candidatePledge.Target.UnmarshalText1([]byte(txDataInfo[tokenPosMinerAddress])); err != nil {
		log.Warn("Candidate pledgeNew", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentCandidatePledge
	}
	if  candidatePledge.Target==txSender {
		log.Warn("Candidate pledgeNew", "miner address is txSender", candidatePledge.Target)
		return currentCandidatePledge
	}

	if _, ok := snap.PosPledge[candidatePledge.Target]; ok {
		log.Warn("Candidate pledgeNew", "candidate already exist", candidatePledge.Target)
		return currentCandidatePledge
	}

	targetMiner:=snap.findPosTargetMiner(candidatePledge.Manager)
	nilAddr := common.Address{}
	if targetMiner!=nilAddr{
		log.Warn("Candidate pledgeNew", "one address can only pledge one miner ", targetMiner)
		return currentCandidatePledge
	}
	entrustMiner:=snap.findPosTargetMiner(candidatePledge.Target)
	if entrustMiner!=nilAddr{
		log.Warn("Candidate pledgeNew", "miner has pledge one miner ", candidatePledge.Target)
		return currentCandidatePledge
	}

	if snap.isPosMinerManager(candidatePledge.Manager){
		log.Warn("Candidate pledgeNew", "manager is pos manager", candidatePledge.Manager)
		return currentCandidatePledge
	}

	if snap.isPosMinerManager(candidatePledge.Target){
		log.Warn("Candidate pledgeNew", "miner is pos manager", candidatePledge.Target)
		return currentCandidatePledge
	}

	if state.GetBalance(txSender).Cmp(candidatePledge.Amount) < 0 {
		log.Warn("Candidate pledgeNew", "balance", state.GetBalance(txSender))
		return currentCandidatePledge
	}
	state.SubBalance(txSender, candidatePledge.Amount)
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x61edf63329be99ab5b931ab93890ea08164175f1bce7446645ba4c1c7bdae3a8")) //web3.sha3("PledgeLock(address,uint256)")
	topics[1].SetBytes(candidatePledge.Target.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumCndLock).Bytes())
	data := common.Hash{}
	data.SetBytes(candidatePledge.Amount.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, data.Bytes())
	currentCandidatePledge = append(currentCandidatePledge, candidatePledge)
	return currentCandidatePledge
}

func (snap *Snapshot) isPosMinerManager(target common.Address) bool {
	for _,item:=range snap.PosPledge{
		if item.Manager==target{
			return true
		}
	}
	return false
}

func (s *Snapshot) updateTallyState() [] Vote {
	var tallyVote [] Vote
	for tallyAddress,vote :=range s.Votes {
		amount :=big.NewInt(0)
		if item, ok := s.PosPledge[tallyAddress]; ok {
			amount = new(big.Int).Add(amount, item.TotalAmount)
		}
		tallyVote = append(tallyVote, Vote{
			Voter :  vote.Voter,
			Candidate : vote.Candidate,
			Stake : amount,
		})

	}
	return tallyVote
}
func (s *Snapshot) updateMinerState() []MinerStakeRecord {
	var tallyMiner []MinerStakeRecord
	for minerAddress, pledge := range s.PosPledge {
		if _,ok:=s.Tally[minerAddress];ok {
			continue
		}
		if credit, ok := s.Punished[minerAddress]; ok && defaultFullCredit-minCalSignerQueueCredit >= credit {
			continue
		}
		amount := pledge.TotalAmount
		if _, ok := s.TallyMiner[minerAddress]; ok {
			s.TallyMiner[minerAddress].Stake = new(big.Int).Set(amount)
		} else {
			s.TallyMiner[minerAddress] = &CandidateState{
				SignerNumber: 0,
				Stake:        new(big.Int).Set(amount),
			}
		}
		tallyMiner = append(tallyMiner, MinerStakeRecord{
			Target: minerAddress,
			Stake:  new(big.Int).Set(amount),
		})
	}
	return tallyMiner
}

func (s *Snapshot) updateSignerNumber(sigers []common.Address, headerNumber uint64) {
	for _, minerAddress := range sigers {
		if _, ok := s.TallyMiner[minerAddress]; ok {
			s.TallyMiner[minerAddress].SignerNumber += 1
		}
		if _,ok := s.Tally[minerAddress];ok {
			if _,isOk := s.TallySigner[minerAddress];isOk{
				s.TallySigner[minerAddress] =s.TallySigner[minerAddress]+1
			}else{
				s.TallySigner[minerAddress]=1
			}
		}

	}
	if headerNumber%clearSignNumberPerid == 0 {
		for address, _ := range s.TallySigner {
			s.TallySigner[address] = 0
		}
		for address, _ := range s.TallyMiner {
			s.TallyMiner[address].SignerNumber = 0
		}
	}

}
func (s *Snapshot) checkPosPledgePunish(address common.Address,headerNumber uint64){
		if pledge,ok1:=s.PosPledge[address];ok1{
			if _,ok2:=s.Punished[address];ok2{
				if pledge.LastPunish == 0{
					pledge.LastPunish=headerNumber
				}
			}else{
				if pledge.LastPunish > 0{
					pledge.LastPunish=0
				}
			}
		}
}

func (s *Snapshot) updateSnapshotForPunish(signerMissing []common.Address, headerNumber *big.Int, coinbase common.Address) {

	for _, signerEach := range signerMissing {
		if _, ok := s.Punished[signerEach]; ok {
			// 10 times of defaultFullCredit is big enough for calculate signer order
			if s.Punished[signerEach]+missingPublishCredit <= defaultFullCredit {
				s.Punished[signerEach] += missingPublishCredit
			} else {
				s.Punished[signerEach] = defaultFullCredit
			}
		} else {
			s.Punished[signerEach] = missingPublishCredit
		}
	}
	s.SignerMissing = make([]common.Address, len(signerMissing))
	copy(s.SignerMissing, signerMissing)
	// reduce the punish of sign signer
	if _, ok := s.Punished[coinbase]; ok {
		if s.Punished[coinbase] > signRewardCredit {
			s.Punished[coinbase] -= signRewardCredit
		} else {
			delete(s.Punished, coinbase)
		}
	}
	// reduce the punish for all punished
	for _, signerEach := range s.Signers {
		sigerAddr := common.HexToAddress(signerEach.String())
		if _, ok := s.Punished[sigerAddr]; ok {
			if s.Punished[sigerAddr] > autoRewardCredit {
				s.Punished[sigerAddr] -= autoRewardCredit
				s.updatePosPledgePunish(sigerAddr,headerNumber.Uint64(), headerNumber.Uint64())
			} else {
				delete(s.Punished, sigerAddr)
				s.updatePosPledgePunish(sigerAddr,0, headerNumber.Uint64())
			}
		}else{
			s.updatePosPledgePunish(sigerAddr,0, headerNumber.Uint64())
		}

	}
	// clear all punish score at the beginning of trantor block
	if s.config.IsTrantor(headerNumber) && !s.config.IsTrantor(new(big.Int).Sub(headerNumber, big.NewInt(1))) {
		s.Punished = make(map[common.Address]uint64)
	}
}
func (s *Snapshot) updatePosPledgePunish(address common.Address, punishNumber uint64,headerNumber uint64){
		if item,ok:=s.PosPledge[address];ok{
			if punishNumber == 0 && item.LastPunish >0 {
				item.LastPunish=0
			}
			if punishNumber > 0 && item.LastPunish == 0{
				item.LastPunish=punishNumber
			}
		}
}

func (s *Snapshot) initPosPledge(number uint64) {
	for addr, _ := range s.Tally {
		if _, ok := s.PosPledge[addr]; !ok {
			lastPunish := uint64(0)
			if _, ok1 := s.Punished[addr]; ok1 {
				lastPunish = number
			}
			managerAddr := addr
			if revenue, ok1 := s.RevenueNormal[addr]; ok1 {
				managerAddr = revenue.RevenueAddress
			}

			s.PosPledge[addr] = &PosPledgeItem{
				Manager:     managerAddr,
				Active:      number,
				TotalAmount: big.NewInt(0),
				LastPunish:  lastPunish,
				DisRate:     new(big.Int).Set(posDistributionDefaultRate),
				Detail:      make(map[common.Hash]*PledgeDetail),
			}
			if _, ok2 := s.TallyMiner[addr]; ok2 {
				delete(s.TallyMiner, addr)
			}
		}
	}
}

func (a *Alien) processCandidateChangeManager(currentCandidateManager []CandidateChangeManagerRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot, number uint64) []CandidateChangeManagerRecord {
	if len(txDataInfo) <= 4 {
		log.Warn("Candidate ChangeManager", "parameter number", len(txDataInfo))
		return currentCandidateManager
	}
	postion := 3
	minerAddress := common.Address{}
	if err := minerAddress.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate ChangeManager", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentCandidateManager
	}
	var posPledgeItem *PosPledgeItem
	if oldBind, ok := snap.PosPledge[minerAddress]; ok {
		if oldBind.Manager !=txSender {
			log.Warn("Candidate ChangeManager", "manager address is not txSender", txSender)
			return currentCandidateManager
		}
		posPledgeItem=oldBind
	}else{
		log.Warn("Candidate ChangeManager", "minerAddress is not exist", minerAddress)
		return currentCandidateManager
	}
	candidateManager := CandidateChangeManagerRecord{
		Target: minerAddress,
		Manager:common.Address{},
	}
	postion++
	if err := candidateManager.Manager.UnmarshalText1([]byte(txDataInfo[postion])); err != nil {
		log.Warn("Candidate ChangeManager", "manager address", txDataInfo[postion])
		return currentCandidateManager
	}

	targetMiner:=snap.findPosTargetMiner(candidateManager.Manager)
	nilAddr := common.Address{}
	if targetMiner==nilAddr||targetMiner!=minerAddress{
		log.Warn("Candidate ChangeManager", "manager address has not pledge the miner ", targetMiner,"minerAddress",minerAddress,"manager",candidateManager.Manager)
		return currentCandidateManager
	}

	if snap.isPosMinerManager(candidateManager.Manager){
		log.Warn("Candidate ChangeManager", "manager is pos manager", candidateManager.Manager)
		return currentCandidateManager
	}

	if _, ok := snap.PosPledge[candidateManager.Manager]; ok {
		log.Warn("Candidate ChangeManager", "manager is pos miner", candidateManager.Manager)
		return currentCandidateManager
	}

	managerAmount:=snap.findPosTargetAmount(posPledgeItem,candidateManager.Manager)
	if managerAmount.Cmp(common.Big0)<=0 {
		log.Warn("Candidate ChangeManager","managerAmount small or equal than 0",managerAmount)
		return currentCandidateManager
	}
	if deposit, ok := snap.SystemConfig.Deposit[0]; !ok {
		log.Warn("Candidate ChangeManager","deposit is nil",deposit)
		return currentCandidateManager
	}else{
		if managerAmount.Cmp(deposit)<0 {
			log.Warn("Candidate ChangeManager","managerAmount small than deposit",managerAmount,"deposit",deposit)
			return currentCandidateManager
		}
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x2bbe766662a69775b1eae5e5a51f2d78677c71262ad6d442e7ac2879a9e7e461")) //web3.sha3("CandidateChangeManager(address,address)")
	topics[1].SetBytes(minerAddress.Bytes())
	topics[2].SetBytes(candidateManager.Manager.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentCandidateManager=append(currentCandidateManager,candidateManager)
	return currentCandidateManager
}

func (snap *Snapshot) updateCandidateChangeManager(candidateChangeManager []CandidateChangeManagerRecord, number *big.Int) {
	for _, item := range candidateChangeManager {
		if _, ok := snap.PosPledge[item.Target]; ok {
			snap.PosPledge[item.Target].Manager = item.Manager
		}
	}
}


func (snap *Snapshot) findPosTargetAmount(item *PosPledgeItem,sender common.Address) *big.Int {
	amount:=common.Big0
	details:=item.Detail
	for _,detail:=range details{
		if detail.Address==sender{
			amount=new(big.Int).Add(amount,detail.Amount)
		}
	}
	return amount
}

func (snap *Snapshot) isPosOtherMinerManager(txSender common.Address, target common.Address) bool {
	for miner,item:=range snap.PosPledge{
		if item.Manager==txSender{
			if miner!=target{
				return true
			}
		}
	}
	return false
}