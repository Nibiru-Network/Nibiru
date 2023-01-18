package alien

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/token/common"
	"github.com/token/consensus"
	"github.com/token/core/state"
	"github.com/token/core/types"
	"github.com/token/crypto"
	"github.com/token/ethdb"
	"github.com/token/log"
	"golang.org/x/crypto/sha3"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

const (
	sscEnumPofWithinBasePricePeriod = 15
	sscEnumPofServicePeriod = 16
)
var (
	yearPeridDay=uint64(365)
	pofBasePledgeAmount = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(4))
	//pofPledgeBwAjFactor = big
	pofstatus_normal = uint64(1)
	pofstatus_exit = uint64(2)
	pofDefaultBasePrice = new(big.Int).Mul(big.NewInt(1e+5), big.NewInt(9765625)) // token/Mpbs
	flowAdjustmentFactor=decimal.NewFromFloat(1)
	trafficPricingFactor=decimal.NewFromFloat(0.05)
)
type PofPledgeItem struct {
	Manager     common.Address                `json:"manager"`
	Active      uint64                        `json:"active"`
	PledgeAmount *big.Int                      `json:"pledgeamount"`
	Bandwidth  uint64                         `json:"bandwidth"`   //Declared bandwidth size
	LastBwValid  uint64                       `json:"lastbwvalid"`  //Bandwidth last verification time
	PofPrice    *big.Int                    `json:"pofprice"`  //Flow unit price    Coin / Mpbs
	PledgeStatus   uint64                      `json:"status"`  // status  1 normal 2 exit
}
type PofPledgeReq struct {
	PofMiner     common.Address
	Manager     common.Address
	PledgeAmount *big.Int
	Bandwidth  uint64
	PofPrice    *big.Int
}
type ClaimedBandwidthRecord struct {
	Target    common.Address
	Amount    *big.Int
	Bandwidth uint64
}
type PofMinerPriceRecord struct {
	Target    common.Address
	PofPrice    *big.Int
}
func (a *Alien) processPofCustomTx(txDataInfo []string, headerExtra HeaderExtra, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, snapCache *Snapshot, number *big.Int, state *state.StateDB, chain consensus.ChainHeaderReader, coinBalances map[common.Address]*big.Int) HeaderExtra {
	if  txDataInfo[posCategory]== tokenEventPofReportEn {
		headerExtra.PofReport = a.processPofReportEn(headerExtra.PofReport, txDataInfo,number.Uint64(),snapCache,txSender, tx, receipts, coinBalances)
	}else if txDataInfo[posCategory] == tokenCategoryPofReq {
		headerExtra.PofPledgeReq = a.processPofPledge (headerExtra.PofPledgeReq, txDataInfo, txSender, tx, receipts, state, snapCache)
	} else if txDataInfo[posCategory] == tokenCategoryPofExit {
		headerExtra.PofMinerExit = a.processMinerExit (headerExtra.PofMinerExit, txDataInfo, txSender, tx, receipts, state, snapCache)
	}else if txDataInfo[posCategory] == tokenEventPofChBw {
		headerExtra.ClaimedBandwidth = a.processChangeBandwidth (headerExtra.ClaimedBandwidth, txDataInfo, txSender, tx, receipts, state, snapCache)
	}else if txDataInfo[posCategory] == tokenEventPofprice {
		headerExtra.PofMinerPriceReq = a.processSetPrice(headerExtra.PofMinerPriceReq, txDataInfo, txSender, tx, receipts, state, snapCache)
	}
	return headerExtra
}
func (snap *Snapshot) pofApply(headerExtra HeaderExtra, header *types.Header, db ethdb.Database) (*Snapshot, error) {
	snap.updateFlowReport(headerExtra.PofReport, header.Number)
	snap.updatePofPledgeReq(headerExtra.PofPledgeReq,header.Number.Uint64())
	snap.updatePofPledgeBandwidth(headerExtra.ClaimedBandwidth,header.Number)
	snap.updatePofPledgePrice(headerExtra.PofMinerPriceReq,header.Number)
	snap.updateFlowMinerExit(headerExtra.PofMinerExit,header.Number)
	snap.pofExitLock(header.Number)
	return snap, nil
}

func (a *Alien) processPofReportEn(pofReport []MinerPofReportRecord, txDataInfo []string, number uint64, snap *Snapshot, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, coinBalances map[common.Address]*big.Int) []MinerPofReportRecord {
	if len(txDataInfo) <= 4 {
		log.Warn("En Pof report", "parameter number", len(txDataInfo))
		return pofReport
	}
	position := 3
	enAddr := txSender
	if _, ok := snap.PofPledge[enAddr]; !ok {
		log.Warn("En Pof report", "enAddr is not in PofPledge", enAddr)
		return pofReport
	}
	position++
	verifyArr :=strings.Split(txDataInfo[position],"|")
	if len(verifyArr)==0 {
		log.Warn("En Pof report", " verifyArr len = 0", txSender, len(verifyArr))
		return pofReport
	}
	census := MinerPofReportRecord{
		ChainHash: common.Hash{},
		ReportTime: number,
		ReportContent: []MinerPofReportItem{},
	}
	zeroAddr:=common.Address{}
	enServerFlowValue:=uint64(0)
	var verifyResult []int
	flowrecordlen:=4
	for index,verifydata :=range verifyArr {
		if verifydata == "" {
			continue
		}
		flowrecord :=strings.Split(verifydata,",")
		if len(flowrecord)!=flowrecordlen{
			log.Warn("En Pof report ", "flowrecordlen", len(flowrecord),"index", index)
			continue
		}
		i:=uint8(0)
		reportNumber, err := decimal.NewFromString(flowrecord[i])
		if err !=nil|| reportNumber.Cmp(decimal.Zero)<=0{
			log.Warn("En Pof report ", "reportNumber is wrong", flowrecord[i],"index", index)
			continue
		}
		if !reportNumber.BigInt().IsUint64()||reportNumber.BigInt().Uint64()==uint64(0){
			log.Warn("En Pof report ", "reportNumber is not uint64 or is zero",reportNumber,"checkReportNumber index", index)
			continue
		}
		if !snap.checkReportNumber(reportNumber,number) {
			log.Warn("En Pof report ", "checkReportNumber index", index)
			continue
		}
		i++
		deviceId, err := decimal.NewFromString(flowrecord[i])
		if err !=nil|| deviceId.Cmp(decimal.Zero)<=0{
			log.Warn("En Pof report ", "deviceId is wrong", flowrecord[i],"index", index)
			continue
		}
		if !deviceId.BigInt().IsUint64(){
			log.Warn("En Pof report ", "deviceId is not uint64", flowrecord[i],"index", index)
			continue
		}
		i++
		flowValue, err := decimal.NewFromString(flowrecord[i])
		if err !=nil||flowValue.Cmp(decimal.Zero)<=0{
			log.Warn("En Pof report ", "flowValue is wrong", flowrecord[i],"index", index)
			continue
		}
		if !flowValue.BigInt().IsUint64()||flowValue.BigInt().Uint64()==uint64(0){
			log.Warn("En Pof report ", "flowValue is not uint64 or is zero", flowrecord[i],"index", index)
			continue
		}
		i++
		rsv:=flowrecord[i]
		if len(rsv)==0{
			log.Warn("En Pof report ", "len(vrs) is zero", flowrecord[i],"index", index)
			continue
		}
		sig:=common.FromHex(rsv)
		if len(sig) != crypto.SignatureLength {
			log.Warn("En Pof report ", "wrong size for signature", flowrecord[i],"index", index)
			continue
		}
		var from common.Address
		if singer,ok:=isCheckFlowRecordSign(reportNumber,deviceId,enAddr,flowValue,sig);!ok{
			log.Warn("En Pof report ", "checkFlowRecordSign index", index)
			continue
		}else{
			from=singer
		}
		if from==zeroAddr{
			log.Warn("En Pof report ","index",index)
			continue
		}
		if _,ok:= coinBalances[from];!ok{
			coinBal :=snap.Coin.Get(from)
			coinBalances[from]=new(big.Int).Set(coinBal)
		}
		costCoin :=common.Big0
		if coin,ok:=snap.checkCoinEnoughItem(flowValue.BigInt().Uint64(), coinBalances[from],enAddr) ;!ok{
			log.Warn("En Pof report ", "CheckCoinEnoughItem index", index)
			continue
		}else{
			costCoin =new(big.Int).Set(coin)
		}
		flowValueM:=flowValue.BigInt().Uint64()
		flowReportItem2:= MinerPofReportItem{
			Target:from,
			FlowValue1:0,
			FlowValue2:flowValueM,
			Miner:enAddr,
		}
		census.ReportContent=append(census.ReportContent,flowReportItem2)
		verifyResult=append(verifyResult, index)
		coinBalances[from]=new(big.Int).Sub(coinBalances[from], costCoin)
		enServerFlowValue=enServerFlowValue+flowValueM
	}
	if len(census.ReportContent)>0{
		flowReportItem:= MinerPofReportItem{
			Target:enAddr,
			FlowValue1:enServerFlowValue,
			FlowValue2:0,
			Miner: common.Address{},
		}
		census.ReportContent=append(census.ReportContent,flowReportItem)
		pofReport = append(pofReport, census)
		topicdata := ""
		sort.Ints(verifyResult)
		for _, val := range verifyResult {
			if topicdata == "" {
				topicdata =fmt.Sprintf("%d", val)
			} else {
				topicdata += "," + fmt.Sprintf("%d", val)
			}
		}
		topics := make([]common.Hash, 1)
		topics[0].UnmarshalText([]byte("0xea40f050c9c577748d5ddcdb6a19aab17cacb2fa5f63f3747c516b06b597afd1"))//web3.sha3("Flwrpten(address,uint256)")
		a.addCustomerTxLog(tx, receipts, topics, []byte(topicdata))
	}
	return pofReport
}

func isCheckFlowRecordSign(reportNumber decimal.Decimal,deviceId decimal.Decimal, toAddress common.Address, flowValue decimal.Decimal,sig []byte) (common.Address,bool) {
	zeroAddr:=common.Address{}
	var hash common.Hash
	hasher := sha3.NewLegacyKeccak256()
	toAddressStr:=strings.ToLower(toAddress.String())
	msg := toAddressStr[2:]+reportNumber.String()+deviceId.String()+flowValue.String()
	hasher.Write([]byte(msg))
	hasher.Sum(hash[:0])
	var rBig=new(big.Int).SetBytes(sig[:32])
	var sBig=new(big.Int).SetBytes(sig[32:64])
	if !crypto.ValidateSignatureValues(sig[64], rBig, sBig, true) {
		log.Warn("isCheckFlowRecordSign", "crypto validateSignatureValues fail ", "sign wrong")
		return zeroAddr,false
	}
	pubkey, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		log.Warn("isCheckFlowRecordSign", "crypto.Ecrecover", err)
		return zeroAddr,false
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	return signer,true
}

func (s *Snapshot) calCoinHashVer(roothash common.Hash, number uint64, db ethdb.Database) (*Snapshot,error) {
	if s.Coin.Root() != roothash {
		return s, errors.New("Coin root hash is not same,head:" + roothash.String() + "cal:" + s.Coin.Root().String())
	}
	return s,nil
}

func (snap *Snapshot) updateCoinBalanceCost(flowReport []MinerPofReportRecord, headerNumber *big.Int) {
	for _, items := range flowReport {
		for _, flow := range items.ReportContent {
			if flow.FlowValue2 > 0 {
				address :=flow.Target
				if cost,ok:=snap.calCostCoin(flow.FlowValue2,flow.Miner);ok{
					balance:=snap.Coin.Get(address)
					if balance.Cmp(cost)>0 {
						snap.Coin.Sub(address,cost)
					}else{
						snap.Coin.Del(address)
					}
				}
			}
		}
	}
}

func (s *Snapshot) calCostCoin(value uint64,miner common.Address) (*big.Int,bool) {
	flowValue:=value
	minerPrice:=big.NewInt(0)
	if _,ok:= s.PofPledge[miner];ok{
		minerPrice = s.PofPledge[miner].PofPrice
	}else{
		log.Warn("calCostCoin", "miner is not in PofPledge", miner)
		return nil,false
	}
	cost := new(big.Int).Mul(new(big.Int).SetUint64(flowValue), minerPrice)
	return cost,true
}

func (snap *Snapshot) checkCoinEnoughItem(flowValue uint64, fromCoinBalance *big.Int,miner common.Address) (*big.Int,bool) {
	if costCoin ,ok:=snap.calCostCoin(flowValue,miner);ok{
		coinEnough :=true
		if fromCoinBalance.Cmp(costCoin)<0 {
			coinEnough =false
		}
		return costCoin, coinEnough
	}else{
		return nil,false
	}
}

func (s *Snapshot) checkReportNumber(reportNumber decimal.Decimal, number uint64) bool {
	reportDay:=reportNumber.BigInt().Uint64()/s.getBlockPreDay()
	blockDay:=number/s.getBlockPreDay()
	return (reportDay==blockDay||reportDay==(blockDay-1))&&reportNumber.BigInt().Uint64()<=number
}
func (a *Alien) processPofPledge (currentPofPledgeReq [] PofPledgeReq, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) [] PofPledgeReq {
	if len(txDataInfo) < tokenPofprice {
		log.Warn("pof pledge", "parameter number", len(txDataInfo))
		return currentPofPledgeReq
	}
	pofPledgeReq := PofPledgeReq{
		PofMiner: common.Address{},
		Manager:txSender,
		PledgeAmount:big.NewInt(0),
		Bandwidth:0,
		PofPrice:big.NewInt(0),
	}
	if err := pofPledgeReq.PofMiner.UnmarshalText1([]byte(txDataInfo[tokenPofMinerAddress])); err != nil {
		log.Warn("pof pledge", "miner address", txDataInfo[tokenPofMinerAddress])
		return currentPofPledgeReq
	}
	if _, ok := snap.PofPledge[pofPledgeReq.PofMiner]; ok {
		log.Warn("pof pledge", "miner exiting", pofPledgeReq.PofMiner)
		return currentPofPledgeReq
	}

	if bandwidth, err := strconv.ParseUint(txDataInfo[tokenPofBandwidth], 16, 32); err != nil {
		log.Warn("pof pledge", "bandwidth", txDataInfo[tokenPofBandwidth])
		return currentPofPledgeReq
	} else {
		pofPledgeReq.Bandwidth = bandwidth
	}
	pofPrice,error:= decimal.NewFromString(txDataInfo[tokenPofprice])
	if error !=nil {
		log.Warn("pof pledge", "pofPrice format error", txDataInfo[tokenPofprice])
		return currentPofPledgeReq
	}else{
		pofPledgeReq.PofPrice = pofPrice.BigInt()
	}
	minPrice := new(big.Int).Div(snap.SystemConfig.Deposit[sscEnumPofWithinBasePricePeriod],big.NewInt(10))
	maxPrice := new(big.Int).Mul(snap.SystemConfig.Deposit[sscEnumPofWithinBasePricePeriod],big.NewInt(10))
	if pofPledgeReq.PofPrice.Cmp(minPrice) < 0  || pofPledgeReq.PofPrice.Cmp(maxPrice) > 0{
		log.Warn("pof pledge", "pofPrice error", pofPledgeReq.PofPrice,"minPrice",minPrice,"maxPrice",maxPrice)
		return currentPofPledgeReq
	}
	pofPledgeReq.PledgeAmount= snap.getPofPlegeAmount(pofPledgeReq.Bandwidth)
	if pofPledgeReq.PledgeAmount.Cmp(state.GetBalance(txSender)) >= 0 {
		log.Warn("pof pledge", "balance not  enough", txSender  ,"miner",pofPledgeReq.PofMiner,"txSender",state.GetBalance(txSender))
		return currentPofPledgeReq
	}
	state.SetBalance(txSender, new(big.Int).Sub(state.GetBalance(txSender), pofPledgeReq.PledgeAmount))
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x041e56787332f2495a47171278fa0f1ddb21961f702d0ba53c2bb2c079ccd418")) //web3.sha3("ClaimedBandwidth(address,uint32,uint32)")
	topics[1].SetBytes(pofPledgeReq.PofMiner.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumBndwdthClaimed).Bytes())
	dataList := make([]common.Hash, 3)
	dataList[0].SetBytes(big.NewInt(int64(pofPledgeReq.Bandwidth)).Bytes())
	dataList[1].SetBytes(pofPledgeReq.PofPrice.Bytes())
	dataList[2].SetBytes(pofPledgeReq.PledgeAmount.Bytes())
	data := dataList[0].Bytes()
	data = append(data, dataList[1].Bytes()...)
	data = append(data, dataList[2].Bytes()...)
	a.addCustomerTxLog (tx, receipts, topics, data)
	currentPofPledgeReq = append(currentPofPledgeReq, pofPledgeReq)
	return currentPofPledgeReq
}

func (s *Snapshot) checkRepeatManagerAddr(txSender common.Address) bool{
	for _,item:=range s.PofPledge {
		if item.Manager == txSender {
			return false
		}
	}
	return true
}

func (snap *Snapshot)  getPofPlegeAmount(bandwidth uint64)  *big.Int{
	pledgePrice := pofBasePledgeAmount
	if snap.Number >snap.getBlockPreDay()*yearPeridDay  {
		totalBandWidth:=uint64(0)
		for _,pofItem:=range snap.PofPledge {
			totalBandWidth=totalBandWidth+pofItem.Bandwidth
		}
		if totalBandWidth> 0 {
			allPofCreateToken:=new(big.Int).Add(snap.InspireHarvest,snap.PofHarvest)
			calPledgePrice:=new(big.Int).Div(allPofCreateToken,big.NewInt(int64(totalBandWidth)))
			if calPledgePrice.Cmp(common.Big0) >0 && calPledgePrice.Cmp(pofBasePledgeAmount) < 0 {
				pledgePrice=new(big.Int).Set(calPledgePrice)
			}
		}
	}
	return new(big.Int).Mul(pledgePrice,big.NewInt(int64(bandwidth)))


}
func (snap *Snapshot) updatePofPledgeReq(currentPofPledgeReq [] PofPledgeReq,number uint64) {
	for _, item := range currentPofPledgeReq {
		if _, ok := snap.PofPledge[item.PofMiner]; !ok {
			snap.PofPledge[item.PofMiner] = &PofPledgeItem{
				Manager:item.Manager,
				Active:number,
				PledgeAmount:item.PledgeAmount,
				Bandwidth:item.Bandwidth,
				LastBwValid:0,
				PofPrice:item.PofPrice,
				PledgeStatus:pofstatus_normal,
			}
		}
	}
}

func (a *Alien) processMinerExit (currentPofMinerExit []common.Address, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) []common.Address {
	if len(txDataInfo) < tokenPofMinerAddress+1 {
		log.Warn("pof miner exit", "parameter number", len(txDataInfo))
		return currentPofMinerExit
	}
	minerAddress := common.Address{}
	if err := minerAddress.UnmarshalText1([]byte(txDataInfo[tokenPofMinerAddress])); err != nil {
		log.Warn("pof miner exit", "miner address", txDataInfo[tokenPofMinerAddress])
		return currentPofMinerExit
	}
	if pledgeItem, ok := snap.PofPledge[minerAddress]; !ok {
		log.Warn("pof miner exit", "miner not exist", minerAddress)
		return currentPofMinerExit
	} else {
		if pledgeItem.Manager != txSender {
			log.Warn("Flow miner exit", "txSender no role txSender", txSender ,"miner",minerAddress)
			return currentPofMinerExit
		}
		if pledgeItem.PledgeStatus == pofstatus_exit {
			log.Warn("Flow miner exit", "miner is existed", minerAddress)
			return currentPofMinerExit
		}
		yearBlockNum:=snap.SystemConfig.Deposit[sscEnumPofServicePeriod].Uint64() * a.blockPerDay()
		if snap.Number -pledgeItem.Active <= yearBlockNum{
			log.Warn("Flow miner exit", "Exit Restrictions", minerAddress,"number",snap.Number,"Active",pledgeItem.Active)
			return currentPofMinerExit
		}
	}

	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x9489b96ebcb056332b79de467a2645c56a999089b730c99fead37b20420d58e7")) //web3.sha3("PledgeExit(address)")
	topics[1].SetBytes(minerAddress.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumPofLock).Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentPofMinerExit = append(currentPofMinerExit, minerAddress)
	return currentPofMinerExit
}

func (snap *Snapshot) updateFlowMinerExit(flowMinerExit []common.Address, headerNumber *big.Int) {
	if len(flowMinerExit) == 0{
      return
	}
	for _, miner := range flowMinerExit {
		if pofPledge, ok := snap.PofPledge[miner]; ok {
			pofPledge.PledgeStatus = pofstatus_exit
			if pofPledge.PledgeAmount.Cmp(common.Big0)>0{
				LockReward :=LockRewardRecord{
					Target:miner,
					Amount:new(big.Int).Set(pofPledge.PledgeAmount),
					IsReward:sscEnumPofLock,
				}
				snap.Revenue.PofExitLock.updateLockPofExitData(snap,LockReward,headerNumber)
			}
		}
	}

}

func (snap *Snapshot)  pofExitLock(headerNumber *big.Int){
	if isDeletPofExitPledge(headerNumber.Uint64(),snap.Period){
		for miner,pofItem:= range  snap.PofPledge{
			if pofItem.PledgeStatus == pofstatus_exit {
				delete(snap.PofPledge, miner)
				delete(snap.RevenuePof, miner)
			}
		}
	}

}

func (a *Alien) processChangeBandwidth (claimedBandwidth [] ClaimedBandwidthRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) [] ClaimedBandwidthRecord {
	if len(txDataInfo) <= tokenPofBandwidth {
		log.Warn("changeBandwidth", "parameter error", len(txDataInfo))
		return claimedBandwidth
	}
	minerAddress := common.Address{}
	if err := minerAddress.UnmarshalText1([]byte(txDataInfo[tokenPofMinerAddress])); err != nil {
		log.Warn("changeBandwidth", "miner format error", txDataInfo[tokenPofMinerAddress])
		return claimedBandwidth
	}
	pofItem:=snap.PofPledge[minerAddress]
	if pofItem ==nil  {
		log.Warn("changeBandwidth", "miner not exit", minerAddress)
		return claimedBandwidth
	}
	if pofItem.Manager!=txSender  {
		log.Warn("changeBandwidth", "txSender no role",txSender,"minerAddress", minerAddress)
		return claimedBandwidth
	}
	bwrecord :=ClaimedBandwidthRecord{
		Target:minerAddress,
		Amount:new(big.Int).Set(pofItem.PledgeAmount),
		Bandwidth:0,
	}
	if bandwidth, err := strconv.ParseUint(txDataInfo[tokenPofBandwidth], 16, 32); err != nil {
		log.Warn("changeBandwidth", "bandwidth", txDataInfo[tokenPofBandwidth])
		return claimedBandwidth
	} else {
		if bandwidth  <= 0 {
			log.Warn("changeBandwidth", "minerAddress bandwidth at least 0 ", minerAddress,"bandwidth", bandwidth)
			return claimedBandwidth
		}
		bwrecord.Bandwidth=bandwidth
	}
	totalAmount := snap.getPofPlegeAmount(bwrecord.Bandwidth)
	if totalAmount.Cmp(pofItem.PledgeAmount) > 0 {
		payAmount:=new(big.Int).Sub(totalAmount,pofItem.PledgeAmount)
		if payAmount.Cmp(state.GetBalance(txSender)) > 0 {
			log.Warn("changeBandwidth", "txSender not enough ", txSender,"pledge",totalAmount,"balance",state.GetBalance(txSender))
			return claimedBandwidth
		}
		state.SubBalance(txSender,payAmount)
		bwrecord.Amount=totalAmount
	}

	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xb12bf5b909b60bb08c3e990dcb437a238072a91629c666541b667da82b3ee422"))
	topics[1].SetBytes(minerAddress.Bytes())
	topics[2].SetBytes([]byte(txDataInfo[4]))
	reData:=bwrecord.Amount.Bytes()
	a.addCustomerTxLog(tx, receipts, topics, reData)
	claimedBandwidth=append(claimedBandwidth,bwrecord )
	return claimedBandwidth
}

func  (snap *Snapshot)  updatePofPledgeBandwidth(claimedBandwidth [] ClaimedBandwidthRecord,headerNumber *big.Int){
	if len(claimedBandwidth)==0 {
		return
	}
	for _,item:=range claimedBandwidth {
		if pofItem,ok:= snap.PofPledge[item.Target];ok {
			pofItem.PledgeAmount=item.Amount
			pofItem.Bandwidth=item.Bandwidth
		}
	}
}

func (a *Alien) processSetPrice (pofMinerPrice [] PofMinerPriceRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) [] PofMinerPriceRecord {
	if len(txDataInfo) <= tokenchangePofprice {
		log.Warn("PofSetPrice", "parameter error", len(txDataInfo))
		return pofMinerPrice
	}
	minerAddress := common.Address{}
	if err := minerAddress.UnmarshalText1([]byte(txDataInfo[tokenPofMinerAddress])); err != nil {
		log.Warn("PofSetPrice", "miner format error", txDataInfo[tokenPofMinerAddress])
		return pofMinerPrice
	}
	pofItem:=snap.PofPledge[minerAddress]
	if pofItem ==nil  {
		log.Warn("PofSetPrice", "miner not exit", minerAddress)
		return pofMinerPrice
	}
	if pofItem.Manager!=txSender  {
		log.Warn("PofSetPrice", "txSender no role",txSender,"minerAddress", minerAddress)
		return pofMinerPrice
	}
	pofPriceRecord := PofMinerPriceRecord {
		Target:minerAddress,
		PofPrice:big.NewInt(0),
	}
	unitPrice,err:=decimal.NewFromString(txDataInfo[tokenchangePofprice])
	if err != nil {
		log.Warn("PofSetPrice","input price error", txDataInfo[tokenchangePofprice],"err",err)
		return pofMinerPrice
	} else {
		if unitPrice.Cmp(decimal.Zero) <= 0{
			log.Warn("PofSetPrice", "pofPrice at least 0 ","minerAddress", minerAddress,"pofPrice", unitPrice)
			return pofMinerPrice
		}
		pofPriceRecord.PofPrice=unitPrice.BigInt()
	}
	minPrice := new(big.Int).Div(snap.SystemConfig.Deposit[sscEnumPofWithinBasePricePeriod],big.NewInt(10))
	maxPrice := new(big.Int).Mul(snap.SystemConfig.Deposit[sscEnumPofWithinBasePricePeriod],big.NewInt(10))
	if pofPriceRecord.PofPrice.Cmp(minPrice) < 0 || pofPriceRecord.PofPrice.Cmp(maxPrice) > 0{
			log.Warn("PofSetPrice", "PofPrice value error ", minerAddress,"PofPrice ",pofPriceRecord.PofPrice,"minPrice",minPrice,"maxPrice",maxPrice)
			return pofMinerPrice
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xb12bf5b909b60bb08c3e990dcb437a238072a91629c666541b667da82b3ee889"))
	topics[1].SetBytes(minerAddress.Bytes())
	topics[2].SetBytes([]byte(txDataInfo[tokenchangePofprice]))
	reData:=[]byte(txDataInfo[tokenchangePofprice])
	a.addCustomerTxLog(tx, receipts, topics, reData)
	pofMinerPrice=append(pofMinerPrice,pofPriceRecord )
	return pofMinerPrice
}

func  (snap *Snapshot)  updatePofPledgePrice(pofMinerPrice [] PofMinerPriceRecord,headerNumber *big.Int){
	if len(pofMinerPrice)==0 {
		return
	}
	for _,item:=range pofMinerPrice {
		if pofItem,ok:= snap.PofPledge[item.Target];ok {
			pofItem.PofPrice=new(big.Int).Set(item.PofPrice)
		}
	}
}