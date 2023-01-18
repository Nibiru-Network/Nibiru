package alien

import (
	"errors"
	"fmt"
	"github.com/token/common"
	"github.com/token/consensus"
	"math/big"
	"reflect"
	"strconv"
)

const (
	lr_s="LockReward"
	en_s="ExchangeCoin"
	db_s="DeviceBind"
	cp_s   ="CandidatePunish"
	ms_s   ="MinerStake"
	cb_s   ="ClaimedBandwidth"
	cd_s   ="ConfigDeposit"
	ci_s   ="ConfigISPQOS"
	lp_s   ="LockParameters"
	ma_s   ="ManagerAddress"
	gp_s   ="GrantProfit"
	pr_s   ="PofReport"
	mprt_s ="MinerPofReportItem"
	epn_s     = "CandidatePledgeNew"
	epent_s   = "CandidatePledgeEntrust"
	epente_s  = "CandidatePEntrustExit"
	eae_s     = "CandidateAutoExit"
	ecr_s     = "CandidateChangeRate"
	ppr_s     ="PofPledgeReq"
	pmpr_s    ="PofMinerPriceReq"
	ccm_s    ="CandidateChangeManager"
)
func verifyHeaderExtern(currentExtra *HeaderExtra, verifyExtra *HeaderExtra) error {

	//ExchangeCoin               []ExchangeCoinRecord
	err := verifyExchangeCoin(currentExtra.ExchangeCoin, verifyExtra.ExchangeCoin)
	if err != nil {
		return err
	}
	//LockReward                []LockRewardRecord
	err = verifyLockReward(currentExtra.LockReward, verifyExtra.LockReward)
	if err != nil {
		return err
	}

	//DeviceBind                []DeviceBindRecord
	err = verifyDeviceBind(currentExtra.DeviceBind, verifyExtra.DeviceBind)
	if err != nil {
		return err
	}

	//CandidatePunish           []CandidatePunishRecord
	err = verifyCandidatePunish(currentExtra.CandidatePunish, verifyExtra.CandidatePunish)
	if err != nil {
		return err
	}
	//MinerStake                []MinerStakeRecord
	err = verifyMinerStake(currentExtra.MinerStake, verifyExtra.MinerStake)
	if err != nil {
		return err
	}

	//CandidateExit             []common.Address
	err = verifyExit(currentExtra.CandidateExit, verifyExtra.CandidateExit,"CandidateExit")
	if err != nil {
		return err
	}

	//ClaimedBandwidth          []ClaimedBandwidthRecord
	err = verifyClaimedBandwidth(currentExtra.ClaimedBandwidth, verifyExtra.ClaimedBandwidth)
	if err != nil {
		return err
	}

	//PofMinerExit             []common.Address
	err = verifyExit(currentExtra.PofMinerExit, verifyExtra.PofMinerExit,"PofMinerExit")
	if err != nil {
		return err
	}


	//ConfigExchRate            uint32
	err = verifyUint32Config(currentExtra.ConfigExchRate, verifyExtra.ConfigExchRate,"ConfigExchRate")
	if err != nil {
		return err
	}
	//ConfigOffLine             uint32
	err = verifyUint32Config(currentExtra.ConfigOffLine, verifyExtra.ConfigOffLine,"ConfigOffLine")
	if err != nil {
		return err
	}

	//ConfigDeposit             []ConfigDepositRecord
	err = verifyConfigDeposit(currentExtra.ConfigDeposit, verifyExtra.ConfigDeposit)
	if err != nil {
		return err
	}

	//ConfigISPQOS              []ISPQOSRecord
	err = verifyConfigISPQOS(currentExtra.ConfigISPQOS, verifyExtra.ConfigISPQOS)
	if err != nil {
		return err
	}

	//LockParameters            []LockParameterRecord
	err = verifyLockParameters(currentExtra.LockParameters, verifyExtra.LockParameters)
	if err != nil {
		return err
	}

	//ManagerAddress            []ManagerAddressRecord
	err = verifyManagerAddress(currentExtra.ManagerAddress, verifyExtra.ManagerAddress)
	if err != nil {
		return err
	}
	//PofHarvest               *big.Int
	err = verifyHarvest(currentExtra.PofHarvest, verifyExtra.PofHarvest,"PofHarvest")
	if err != nil {
		return err
	}
	//GrantProfit               []consensus.GrantProfitRecord
	err = verifyGrantProfit(currentExtra.GrantProfit, verifyExtra.GrantProfit)
	if err != nil {
		return err
	}

	//PofReport                []MinerPofReportRecord
	err = verifyPofReport(currentExtra.PofReport, verifyExtra.PofReport)
	if err != nil {
		return err
	}
	return nil

	//CoinDataRoot
	if currentExtra.CoinDataRoot != verifyExtra.CoinDataRoot {
		return errors.New("Compare CoinDataRoot, current is " + currentExtra.CoinDataRoot.String() + ". but verify is " + verifyExtra.CoinDataRoot.String())
	}
	//GrantProfitHash
	if currentExtra.GrantProfitHash != verifyExtra.GrantProfitHash {
		return errors.New("Compare GrantProfitHash, current is " + currentExtra.GrantProfitHash.String() + ". but verify is " + verifyExtra.GrantProfitHash.String())
	}

	//epn_s    = "CandidatePledgeNew"
	err = verifyCandidatePledgeNew(currentExtra.CandidatePledgeNew, verifyExtra.CandidatePledgeNew)
	if err != nil {
		return err
	}
	//epent_s   = "CandidatePledgeEntrust"
	err = verifyCandidatePledgeEntrust(currentExtra.CandidatePledgeEntrust, verifyExtra.CandidatePledgeEntrust)
	if err != nil {
		return err
	}
	//epente_s   = "CandidatePEntrustExit"
	err = verifyCandidatePEntrustExit(currentExtra.CandidatePEntrustExit, verifyExtra.CandidatePEntrustExit)
	if err != nil {
		return err
	}
	//eae_s     = "CandidateAutoExit"
	err = verifyCandidateAutoExit(currentExtra.CandidateAutoExit, verifyExtra.CandidateAutoExit)
	if err != nil {
		return err
	}
	//ecr_s     = "CandidateChangeRate"
	err = verifyCandidateChangeRate(currentExtra.CandidateChangeRate, verifyExtra.CandidateChangeRate)
	if err != nil {
		return err
	}
	//ppr_s     ="PofPledgeReq"
	err = verifyPofPledgeReq(currentExtra.PofPledgeReq, verifyExtra.PofPledgeReq)
	if err != nil {
		return err
	}
	//pmpr_s    ="PofMinerPriceReq"
	err = verifyPofMinerPriceReq(currentExtra.PofMinerPriceReq, verifyExtra.PofMinerPriceReq)
	if err != nil {
		return err
	}
	//InspireHarvest               *big.Int
	err = verifyHarvest(currentExtra.InspireHarvest, verifyExtra.InspireHarvest,"InspireHarvest")
	if err != nil {
		return err
	}
	//ccm_s    ="CandidateChangeManager"
	err = verifyCandidateChangeManager(currentExtra.CandidateChangeManager, verifyExtra.CandidateChangeManager)
	if err != nil {
		return err
	}
	return nil
}

func verifyUint32Config(current uint32, verify uint32,name string) error {
	if current!=verify{
		s:=strconv.FormatUint(uint64(current), 10)
		s2:=strconv.FormatUint(uint64(verify), 10)
		return errors.New("Compare "+name+", current is "+s+". but verify is "+s2)
	}
	return nil
}


func verifyLockReward(current []LockRewardRecord, verify []LockRewardRecord) error {
	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(lr_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(lr_s)
	}
	if len(current) != len(verify) {
		return errorsMsg3(lr_s,len(current),len(verify) )
	}
	if len(current)==0{
		return nil
	}
	err:=compareLockReward(current,verify)
	if err!=nil{
		return err
	}
	err=compareLockReward(verify,current)
	if err!=nil{
		return err
	}
	return nil
}
func compareLockReward(a []LockRewardRecord, b []LockRewardRecord) error{
	b2:= make([]LockRewardRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Amount.Cmp(v.Amount) == 0  && c.FlowValue1 == v.FlowValue1 && c.FlowValue2 == v.FlowValue2 && c.IsReward == v.IsReward && c.Target == v.Target  {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(lr_s,c)
		}
	}
	return nil
}


func verifyExchangeCoin(current []ExchangeCoinRecord, verify []ExchangeCoinRecord) error {
	arrLen, err := verifyArrayBasic(en_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err= compareExchangeCoin(current,verify)
	if err!=nil{
		return err
	}
	err= compareExchangeCoin(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareExchangeCoin(a []ExchangeCoinRecord, b []ExchangeCoinRecord) error{
	b2:= make([]ExchangeCoinRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target && c.Amount.Cmp(v.Amount)==0 {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(en_s,c)
		}
	}
	return nil
}

func verifyDeviceBind(current []DeviceBindRecord, verify []DeviceBindRecord) error {
	arrLen, err := verifyArrayBasic(db_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareDeviceBind(current,verify)
	if err!=nil{
		return err
	}
	err=compareDeviceBind(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareDeviceBind(a []DeviceBindRecord, b []DeviceBindRecord) error{
	b2:= make([]DeviceBindRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Device == v.Device  && c.Revenue == v.Revenue && c.Contract == v.Contract && c.MultiSign == v.MultiSign && c.Type == v.Type  && c.Bind == v.Bind {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(db_s,c)
		}
	}
	return nil
}

func verifyCandidatePunish(current []CandidatePunishRecord, verify []CandidatePunishRecord) error {
	arrLen, err := verifyArrayBasic(cp_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareCandidatePunish(current,verify)
	if err!=nil{
		return err
	}
	err=compareCandidatePunish(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareCandidatePunish(a []CandidatePunishRecord, b []CandidatePunishRecord) error{
	b2:= make([]CandidatePunishRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target  && c.Amount.Cmp(v.Amount)==0  && c.Credit==v.Credit{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cp_s,c)
		}
	}
	return nil
}

func verifyMinerStake(current []MinerStakeRecord, verify []MinerStakeRecord) error {
	arrLen, err := verifyArrayBasic(ms_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareMinerStake(current,verify)
	if err!=nil{
		return err
	}
	err=compareMinerStake(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareMinerStake(a []MinerStakeRecord, b []MinerStakeRecord) error{
	b2:= make([]MinerStakeRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target  && c.Stake.Cmp(v.Stake)==0{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ms_s,c)
		}
	}
	return nil
}

func verifyExit(current []common.Address, verify []common.Address,name string) error {
	arrLen, err := verifyArrayBasic(name, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareExit(current,verify,name)
	if err!=nil{
		return err
	}
	err=compareExit(verify,current,name)
	if err!=nil{
		return err
	}
	return nil
}

func compareExit(a []common.Address, b []common.Address,name string) error {
	b2:= make([]common.Address, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c == v {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(name,c)
		}
	}
	return nil
}

func verifyClaimedBandwidth(current []ClaimedBandwidthRecord, verify []ClaimedBandwidthRecord) error {
	arrLen, err := verifyArrayBasic(cb_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareClaimedBandwidth(current,verify)
	if err!=nil{
		return err
	}
	err=compareClaimedBandwidth(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareClaimedBandwidth(a []ClaimedBandwidthRecord, b []ClaimedBandwidthRecord) error {
	b2:= make([]ClaimedBandwidthRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.Amount.Cmp(v.Amount)==0&&c.Bandwidth==v.Bandwidth {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cb_s,c)
		}
	}
	return nil
}

func verifyConfigDeposit(current []ConfigDepositRecord, verify []ConfigDepositRecord) error {
	arrLen, err := verifyArrayBasic(cd_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareConfigDeposit(current,verify)
	if err!=nil{
		return err
	}
	err=compareConfigDeposit(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareConfigDeposit(a []ConfigDepositRecord, b []ConfigDepositRecord) error {
	b2:= make([]ConfigDepositRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Who == v.Who&&c.Amount.Cmp(v.Amount)==0 {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(cd_s,c)
		}
	}
	return nil
}

func verifyConfigISPQOS(current []ISPQOSRecord, verify []ISPQOSRecord) error {
	arrLen, err := verifyArrayBasic(ci_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareConfigISPQOS(current,verify)
	if err!=nil{
		return err
	}
	err=compareConfigISPQOS(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareConfigISPQOS(a []ISPQOSRecord, b []ISPQOSRecord) error {
	b2:= make([]ISPQOSRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.ISPID == v.ISPID&&c.QOS==v.QOS {
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ci_s,c)
		}
	}
	return nil
}

func verifyLockParameters(current []LockParameterRecord, verify []LockParameterRecord) error {
	arrLen, err := verifyArrayBasic(lp_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareLockParameters(current,verify)
	if err!=nil{
		return err
	}
	err=compareLockParameters(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareLockParameters(a []LockParameterRecord, b []LockParameterRecord) error {
	b2:= make([]LockParameterRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.LockPeriod == v.LockPeriod&&c.RlsPeriod==v.RlsPeriod &&c.Interval==v.Interval&&c.Who==v.Who{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(lp_s,c)
		}
	}
	return nil
}

func verifyManagerAddress(current []ManagerAddressRecord, verify []ManagerAddressRecord) error {
	arrLen, err := verifyArrayBasic(ma_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareManagerAddress(current,verify)
	if err!=nil{
		return err
	}
	err=compareManagerAddress(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareManagerAddress(a []ManagerAddressRecord, b []ManagerAddressRecord) error {
	b2:= make([]ManagerAddressRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.Who==v.Who{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ma_s,c)
		}
	}
	return nil
}

func verifyHarvest(current *big.Int, verify *big.Int,fh_s string) error {

	if current == nil && verify == nil {
		return nil
	}
	if current == nil && verify != nil {
		return errorsMsg1(fh_s)
	}
	if current != nil && verify == nil {
		return errorsMsg2(fh_s)
	}
	if current != nil && verify != nil && current.Cmp(verify)!=0 {
		return errors.New("Compare "+fh_s+", current is "+current.String()+". but verify is "+verify.String())
	}
	return nil
}

func verifyGrantProfit(current []consensus.GrantProfitRecord, verify []consensus.GrantProfitRecord) error {
	arrLen, err := verifyArrayBasic(gp_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareGrantProfit(current,verify)
	if err!=nil{
		return err
	}
	err=compareGrantProfit(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareGrantProfit(a []consensus.GrantProfitRecord, b []consensus.GrantProfitRecord)error {
	b2:= make([]consensus.GrantProfitRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Which == v.Which&&c.MinerAddress==v.MinerAddress&&c.BlockNumber==v.BlockNumber&&c.Amount.Cmp(v.Amount)==0&&c.RevenueAddress==v.RevenueAddress&&c.RevenueContract==v.RevenueContract&&c.MultiSignature==v.MultiSignature{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(gp_s,c)
		}
	}
	return nil
}


func verifyPofReport(current []MinerPofReportRecord, verify []MinerPofReportRecord) error {
	arrLen, err := verifyArrayBasic(pr_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareFlowReport(current,verify)
	if err!=nil{
		return err
	}
	err=compareFlowReport(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareFlowReport(a []MinerPofReportRecord, b []MinerPofReportRecord)error {
	b2:= make([]MinerPofReportRecord, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.ChainHash == v.ChainHash&&c.ReportTime==v.ReportTime{
				if err:=verifyMinerFlowReportItem(c.ReportContent,v.ReportContent);err==nil{
					find = true
					b2=append(b2[:i],b2[i+1:]...)
					break
				}
			}
		}
		if !find {
			return errorsMsg4(pr_s,c)
		}
	}
	return nil
}

func verifyMinerFlowReportItem(current []MinerPofReportItem, verify []MinerPofReportItem) error {
	arrLen, err := verifyArrayBasic(mprt_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err=compareMinerFlowReportItem(current,verify)
	if err!=nil{
		return err
	}
	err=compareMinerFlowReportItem(verify,current)
	if err!=nil{
		return err
	}
	return nil
}

func compareMinerFlowReportItem(a []MinerPofReportItem, b []MinerPofReportItem)error {
	b2:= make([]MinerPofReportItem, len(b))
	copy(b2,b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target&&c.FlowValue1==v.FlowValue1&&c.FlowValue2==v.FlowValue2&&c.Miner==v.Miner{
				find = true
				b2=append(b2[:i],b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(mprt_s,c)
		}
	}
	return nil
}

func errorsMsg1(name string) error {
	return errors.New("Compare "+name+" , current is nil. but verify is not nil")
}
func errorsMsg2(name string) error {
	return errors.New("Compare "+name+" , current is not nil. but verify is nil")
}
func errorsMsg3(name string,lenc int,lenv int) error {
	return errors.New(fmt.Sprintf("Compare "+name+", The array length is not equals. the current length is %d. the verify length is %d", lenc, lenv))
}
func errorsMsg4(name string,c interface{}) error {
	return errors.New(fmt.Sprintf("Compare "+name+", can't find %v in verify data", c))
}


func verifyArrayBasic(title string, current interface{}, verify interface{}) (int, error) {
	if current == nil {
		if verify == nil {
			return 0, nil
		}
		verifyLen := reflect.ValueOf(verify).Len()
		if verifyLen == 0 {
			return 0, nil
		}
		return 0, errorsMsg1(title)
	}
	currentLen := reflect.ValueOf(current).Len()
	if verify == nil {
		if currentLen == 0 {
			return 0, nil
		} else {
			return 0, errorsMsg2(title)
		}
	}
	verifyLen := reflect.ValueOf(verify).Len()
	if currentLen != verifyLen {
		return 0, errorsMsg3(title, currentLen, verifyLen)
	}
	return currentLen, nil
}

func verifyCandidatePledgeNew(current []CandidatePledgeNewRecord, verify []CandidatePledgeNewRecord) error {
	arrLen, err := verifyArrayBasic(epn_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = compareCandidatePledgeNew(current, verify)
	if err != nil {
		return err
	}
	err = compareCandidatePledgeNew(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func compareCandidatePledgeNew(a []CandidatePledgeNewRecord, b []CandidatePledgeNewRecord) error {
	b2 := make([]CandidatePledgeNewRecord, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target && c.Amount.Cmp(v.Amount) == 0&& c.Manager == v.Manager&& c.Hash == v.Hash {
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(epn_s, c)
		}
	}
	return nil
}

func verifyCandidatePledgeEntrust(current []CandidatePledgeEntrustRecord, verify []CandidatePledgeEntrustRecord) error {
	arrLen, err := verifyArrayBasic(epent_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = compareCandidatePledgeEntrust(current, verify)
	if err != nil {
		return err
	}
	err = compareCandidatePledgeEntrust(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func compareCandidatePledgeEntrust(a []CandidatePledgeEntrustRecord, b []CandidatePledgeEntrustRecord) error {
	b2 := make([]CandidatePledgeEntrustRecord, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target && c.Amount.Cmp(v.Amount) == 0&& c.Address == v.Address&& c.Hash == v.Hash {
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(epent_s, c)
		}
	}
	return nil
}

func verifyCandidatePEntrustExit(current []CandidatePEntrustExitRecord, verify []CandidatePEntrustExitRecord) error {
	arrLen, err := verifyArrayBasic(epente_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = compareCandidatePEntrustExit(current, verify)
	if err != nil {
		return err
	}
	err = compareCandidatePEntrustExit(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func compareCandidatePEntrustExit(a []CandidatePEntrustExitRecord, b []CandidatePEntrustExitRecord) error {
	b2 := make([]CandidatePEntrustExitRecord, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target && c.Amount.Cmp(v.Amount) == 0&& c.Address == v.Address&& c.Hash == v.Hash {
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(epent_s, c)
		}
	}
	return nil
}

func verifyCandidateAutoExit(current []common.Address, verify []common.Address) error {
	arrLen, err := verifyArrayBasic(eae_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = compareCandidateAutoExit(current, verify)
	if err != nil {
		return err
	}
	err = compareCandidateAutoExit(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func compareCandidateAutoExit(a []common.Address, b []common.Address) error {
	b2 := make([]common.Address, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c==v {
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(eae_s, c)
		}
	}
	return nil
}


func verifyCandidateChangeRate(current []CandidateChangeRateRecord, verify []CandidateChangeRateRecord) error {
	arrLen, err := verifyArrayBasic(ecr_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = compareCandidateChangeRate(current, verify)
	if err != nil {
		return err
	}
	err = compareCandidateChangeRate(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func compareCandidateChangeRate(a []CandidateChangeRateRecord, b []CandidateChangeRateRecord) error {
	b2 := make([]CandidateChangeRateRecord, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target && c.Rate.Cmp(v.Rate) == 0 {
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ecr_s, c)
		}
	}
	return nil
}

func verifyPofPledgeReq(current []PofPledgeReq, verify []PofPledgeReq) error {
	arrLen, err := verifyArrayBasic(ppr_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = comparePofPledgeReq(current, verify)
	if err != nil {
		return err
	}
	err = comparePofPledgeReq(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func comparePofPledgeReq(a []PofPledgeReq, b []PofPledgeReq) error {
	b2 := make([]PofPledgeReq, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.PofMiner == v.PofMiner &&c.Manager == v.Manager &&c.PledgeAmount.Cmp(v.PledgeAmount) == 0 &&c.Bandwidth==v.Bandwidth &&c.PofPrice.Cmp(v.PofPrice)==0{
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ppr_s, c)
		}
	}
	return nil
}

func verifyPofMinerPriceReq(current []PofMinerPriceRecord, verify []PofMinerPriceRecord) error {
	arrLen, err := verifyArrayBasic(pmpr_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = comparePofMinerPriceReq(current, verify)
	if err != nil {
		return err
	}
	err = comparePofMinerPriceReq(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func comparePofMinerPriceReq(a []PofMinerPriceRecord, b []PofMinerPriceRecord) error {
	b2 := make([]PofMinerPriceRecord, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target &&c.PofPrice.Cmp(v.PofPrice) == 0{
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(pmpr_s, c)
		}
	}
	return nil
}

func verifyCandidateChangeManager(current []CandidateChangeManagerRecord, verify []CandidateChangeManagerRecord) error {
	arrLen, err := verifyArrayBasic(ccm_s, current, verify)
	if err != nil {
		return err
	}
	if arrLen == 0 {
		return nil
	}
	err = compareCandidateChangeManager(current, verify)
	if err != nil {
		return err
	}
	err = compareCandidateChangeManager(verify, current)
	if err != nil {
		return err
	}
	return nil
}

func compareCandidateChangeManager(a []CandidateChangeManagerRecord, b []CandidateChangeManagerRecord) error {
	b2 := make([]CandidateChangeManagerRecord, len(b))
	copy(b2, b)
	for _, c := range a {
		find := false
		for i, v := range b2 {
			if c.Target == v.Target &&c.Manager == v.Manager{
				find = true
				b2 = append(b2[:i], b2[i+1:]...)
				break
			}
		}
		if !find {
			return errorsMsg4(ccm_s, c)
		}
	}
	return nil
}