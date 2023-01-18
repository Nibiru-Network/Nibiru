package alien

import "math/big"

const (
	checkpointInterval = 360              //360        // About N hours if config.period is N

	secondsPerDay                    = 24 * 60 * 60 // Number of seconds for one day
	accumulateInspireRewardInterval = 1 * 60 * 60  // accumulate pof reward interval every day
	accumulatePofRewardInterval      = 2 * 60 * 60  // accumulate pof reward interval every day

	paySignerRewardInterval    = 0                   // pay singer reward  interval every day
	payPofRewardInterval       = 2*60*60 + 30*60     //  pay pof reward  interval every day
	payInspireRewardInterval = 1 * 60 * 60 + 30*60 //  pay bandwidth reward  interval every day
	payPofExitInterval       = 1 * 60 * 60 + 10*60
	payPOSExitInterval = 1 * 60 * 60 + 50*60  //  pay bandwidth reward  interval every day
	checkPOSAutoExit= 1 * 60 * 60 + 60*60

	signerPledgeLockParamPeriod    = 30 * 24 * 60 * 60
	signerPledgeLockParamRlsPeriod = 365 * 24 * 60 * 60
	signerPledgeLockParamInterval  = 24 * 60 * 60

	pofPledgeLockParamPeriod    = 30 * 24 * 60 * 60
	pofPledgeLockParamRlsPeriod = 365 * 24 * 60 * 60
	pofPledgeLockParamInterval  = 24 * 60 * 60

	rewardLockParamPeriod      = 30 * 24 * 60 * 60
	rewardLockParamRlsPeriod   = 365 * 24 * 60 * 60
	rewardLockParamInterval    = 24 * 60 * 60

)

var (
	minCndPledgeBalance = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(100)) // candidate pledge balance
	minSignerLockBalance    = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(0)) // signer reward lock balance
	minPofLockBalance       = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(0)) // pof reward lock balance
	minBandwidthLockBalance = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(0)) // bandwidth reward lock balance

)

func (a *Alien) blockPerDay() uint64 {
	return secondsPerDay / a.config.Period
}

func (a *Alien) blockAccumulatePofRewardInterval() uint64 {
	return accumulatePofRewardInterval / a.config.Period
}

func (a *Alien) blockAccumulateInspireRewardInterval() uint64 {
	return accumulateInspireRewardInterval / a.config.Period
}

func (a *Alien) blockPaySignerRewardInterval() uint64 {
	return paySignerRewardInterval / a.config.Period
}

func (a *Alien) blockPayPofRewardInterval() uint64 {
	return payPofRewardInterval / a.config.Period
}

func (a *Alien) isAccumulatePofRewards(number uint64) bool {
	block := a.blockAccumulatePofRewardInterval()
	heigtPerDay := a.blockPerDay()
	return block == number%heigtPerDay && block != number
}

func (a *Alien) isAccumulateInspireRewards(number uint64) bool {
	block := a.blockAccumulateInspireRewardInterval()
	blockPerDay :=  a.blockPerDay()
	return block == number%blockPerDay && block != number
}

func isPayInspireRewards(number uint64, period uint64) bool {
	block := payInspireRewardInterval / period
	blockPerDay := secondsPerDay / period
	return block == number%blockPerDay && block != number
}

func isPayPofRewards(number uint64, period uint64) bool {
	block := payPofRewardInterval / period
	blockPerDay := secondsPerDay / period
	return block == number%blockPerDay && block != number
}
func isPaySignerRewards(number uint64, period uint64) bool {
	block := paySignerRewardInterval / period
	blockPerDay := secondsPerDay / period
	return block == number%blockPerDay && block != number
}

func isPayPosExit(number uint64, period uint64) bool {
	block := payPOSExitInterval / period
	blockPerDay := secondsPerDay / period
	return block == number%blockPerDay && block != number
}
func isPayPoFExit(number uint64, period uint64) bool {
	block := payPofExitInterval / period
	blockPerDay := secondsPerDay / period
	return block == number%blockPerDay && block != number
}

func isCheckPOSAutoExit(number uint64, period uint64) bool {
	block := checkPOSAutoExit / period
	blockPerDay := secondsPerDay / period
	return block == number%blockPerDay && block != number
}

func  isDeletPofExitPledge(number uint64,period uint64) bool {
	block := accumulatePofRewardInterval / period
	heigtPerDay := secondsPerDay / period
	return block == number%heigtPerDay && block != number
}
