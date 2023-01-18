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
	"errors"
	"fmt"
	"github.com/token/common/hexutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/token/common"
	"github.com/token/consensus"
	"github.com/token/core/state"
	"github.com/token/core/types"
	"github.com/token/log"
	"github.com/token/params"
	"github.com/token/rlp"
)

const (
	ufoVersion            = "1"

	ufoPrefix   = "ufo"
	tokenPrefix = "token"
	sscPrefix   = "SSC"

	ufoCategoryEvent      = "event"
	ufoCategoryLog        = "oplog"
	ufoCategorySC         = "sc"
	ufoEventVote          = "vote"
	ufoEventConfirm       = "confirm"
	ufoEventPorposal      = "proposal"
	ufoEventDeclare       = "declare"

	tokenCategoryExch     = "Exch"
	tokenCategoryBind     = "Bind"
	tokenCategoryUnbind   = "Unbind"
	tokenCategoryRebind   = "Rebind"
	tokenCategoryCandReq  = "CandReq"
	tokenCategoryCandExit = "CandExit"
	tokenCategoryCandPnsh = "CandPnsh"
	tokenCategoryPofReq   = "pofReq"
	tokenCategoryPofExit  = "pofExit"
	tokenEventPofReportEn = "pofrpten"
	tokenEventPofChBw = "pofchbw"
	tokenEventPofprice = "pofprice"
	sscCategoryExchRate = "ExchRate"
	sscCategoryDeposit  = "Deposit"
	sscCategoryCndLock  = "CndLock"
	sscCategoryPofLock  = "PofLock"
	sscCategoryRwdLock  = "RwdLock"
	sscCategoryOffLine  = "OffLine"
	sscCategoryManager  = "Manager"

	ufoMinSplitLen        = 3

	posPrefix             = 0
	posVersion            = 1
	posCategory           = 2

	posEventVote          = 3
	posEventConfirm       = 3
	posEventProposal      = 3
	posEventDeclare       = 3
	posEventConfirmNumber = 4

	tokenPosExchAddress    = 3
	tokenPosExchValue      =	4
	tokenPosMinerAddress   = 3
	tokenPosRevenueType    = 4
	tokenPosMiltiSign      = 6
	tokenPosRevenueAddress = 7
	tokenPofMinerAddress   = 3
	tokenPofBandwidth = 4
	tokenPofprice  = 5
	tokenchangePofprice  = 4

	sscPosExchRate        = 3
	sscPosDeposit         = 3
	sscPosDepositWho      = 4
	sscPosLockPeriod      = 3
	sscPosRlsPeriod       = 4
	sscPosInterval        = 5
	sscPosOffLine         = 3
	sscPosQosID           = 3
	sscPosQosValue        = 4
	sscPosWdthPnsh        = 4
	sscPosManagerID       = 3
	sscPosManagerAddress  = 4

	sscEnumCndLock = 0
	sscEnumPofLock = 1
	sscEnumRwdLock = 2

	sscEnumSignerReward    = 3
	sscEnumPofReward       = 4
	sscEnumPofInspireReward = 5

	sscEnumBndwdthClaimed = 0
	sscEnumBndwdthPunish  = 1
	sscEnumExchRate       = 0
	sscEnumSystem         = 1
	sscEnumWdthPnsh       = 2
	sscEnumFlowReport     = 3

	/*
	 *  proposal type
	 */
	proposalTypeCandidateAdd                  = 1
	proposalTypeCandidateRemove               = 2
	proposalTypeMinerRewardDistributionModify = 3 // count in one thousand
	proposalTypeSideChainAdd                  = 4
	proposalTypeSideChainRemove               = 5
	proposalTypeMinVoterBalanceModify         = 6
	proposalTypeProposalDepositModify         = 7
	proposalTypeRentSideChain                 = 8 // use TTC to buy coin on side chain

	/*
	 * proposal related
	 */
	maxValidationLoopCnt     = 12342                   // About one month if period = 10 & 21 super nodes
	minValidationLoopCnt     = 4                       // just for test, Note: 12350  About three days if seal each block per second & 21 super nodes
	defaultValidationLoopCnt = 2880                    // About one week if period = 10 & 21 super nodes
	maxProposalDeposit       = 100000                  // If no limit on max proposal deposit and 1 billion TTC deposit success passed, then no new proposal.
	minSCRentFee             = 100                     // 100 TTC
	minSCRentLength          = 259200                  // number of block about 1 month if period is 10
	defaultSCRentLength      = minSCRentLength * 3     // number of block about 3 month if period is 10
	maxSCRentLength          = defaultSCRentLength * 4 // number of block about 1 year if period is 10

	/*
	 * notice related
	 */
	noticeTypeGasCharging = 1
)

// RefundGas :
// refund gas to tx sender
type RefundGas map[common.Address]*big.Int

// RefundPair :
type RefundPair struct {
	Sender   common.Address
	GasPrice *big.Int
}

// RefundHash :
type RefundHash map[common.Hash]RefundPair

// Vote :
// vote come from custom tx which data like "ufo:1:event:vote"
// Sender of tx is Voter, the tx.to is Candidate
// Stake is the balance of Voter when create this vote
type Vote struct {
	Voter     common.Address `json:"voter"`
	Candidate common.Address `json:"candidate"`
	Stake     *big.Int       `json:"stake"`
}

// Confirmation :
// confirmation come  from custom tx which data like "ufo:1:event:confirm:123"
// 123 is the block number be confirmed
// Sender of tx is Signer only if the signer in the SignerQueue for block number 123
type Confirmation struct {
	Signer      common.Address
	BlockNumber *big.Int
}

// Proposal :
// proposal come from  custom tx which data like "ufo:1:event:proposal:candidate:add:address" or "ufo:1:event:proposal:percentage:60"
// proposal only come from the current candidates
// not only candidate add/remove , current signer can proposal for params modify like percentage of reward distribution ...
type Proposal struct {
	Hash                   common.Hash    `json:"hash"`                   // tx hash
	ReceivedNumber         *big.Int       `json:"receivenumber"`          // block number of proposal received
	CurrentDeposit         *big.Int       `json:"currentdeposit"`         // received deposit for this proposal
	ValidationLoopCnt      uint64         `json:"validationloopcount"`    // validation block number length of this proposal from the received block number
	ProposalType           uint64         `json:"proposaltype"`           // type of proposal 1 - add candidate 2 - remove candidate ...
	Proposer               common.Address `json:"proposer"`               // proposer
	TargetAddress          common.Address `json:"candidateaddress"`       // candidate need to add/remove if candidateNeedPD == true
	MinerRewardPerThousand uint64         `json:"minerrewardperthousand"` // reward of miner + side chain miner
	SCHash                 common.Hash    `json:"schash"`                 // side chain genesis parent hash need to add/remove
	SCBlockCountPerPeriod  uint64         `json:"scblockcountperpersiod"` // the number block sealed by this side chain per period, default 1
	SCBlockRewardPerPeriod uint64         `json:"scblockrewardperperiod"` // the reward of this side chain per period if SCBlockCountPerPeriod reach, default 0. SCBlockRewardPerPeriod/1000 * MinerRewardPerThousand/1000 * BlockReward is the reward for this side chain
	Declares               []*Declare     `json:"declares"`               // Declare this proposal received (always empty in block header)
	MinVoterBalance        uint64         `json:"minvoterbalance"`        // value of minVoterBalance , need to mul big.Int(1e+18)
	ProposalDeposit        uint64         `json:"proposaldeposit"`        // The deposit need to be frozen during before the proposal get final conclusion. (TTC)
	SCRentFee              uint64         `json:"screntfee"`              // number of TTC coin, not wei
	SCRentRate             uint64         `json:"screntrate"`             // how many coin you want for 1 TTC on main chain
	SCRentLength           uint64         `json:"screntlength"`           // minimize block number of main chain , the rent fee will be used as reward of side chain miner.
}

// Declare :
// declare come from custom tx which data like "ufo:1:event:declare:hash:yes"
// proposal only come from the current candidates
// hash is the hash of proposal tx
type Declare struct {
	ProposalHash common.Hash
	Declarer     common.Address
	Decision     bool
}

// SCConfirmation is the confirmed tx send by side chain super node
type SCConfirmation struct {
	Hash     common.Hash
	Coinbase common.Address // the side chain signer , may be diff from signer in main chain
	Number   uint64
	LoopInfo []string
}

// SCSetCoinbase is the tx send by main chain super node which can set coinbase for side chain
type SCSetCoinbase struct {
	Hash     common.Hash
	Signer   common.Address
	Coinbase common.Address
	Type     bool
}

type GasCharging struct {
	Target common.Address `json:"address"` // target address on side chain
	Volume uint64         `json:"volume"`  // volume of gas need charge (unit is ttc)
	Hash   common.Hash    `json:"hash"`    // the hash of proposal, use as id of this proposal
}

type ExchangeCoinRecord struct {
	Target common.Address
	Amount *big.Int
}

type DeviceBindRecord struct {
	Device    common.Address
	Revenue   common.Address
	Contract  common.Address
	MultiSign common.Address
	Type      uint32
	Bind      bool
}

type CandidatePunishRecord struct {
	Target common.Address
	Amount *big.Int
	Credit uint32
}

type BandwidthPunishRecord struct {
	Target   common.Address
	WdthPnsh uint32
	LastBwValid uint64
	PofPrice    *big.Int
}

type ISPQOSRecord struct {
	ISPID uint32
	QOS   uint32
}

type ManagerAddressRecord struct {
	Target common.Address
	Who    uint32
}

type LockParameterRecord struct {
	LockPeriod uint32
	RlsPeriod  uint32
	Interval   uint32
	Who        uint32
}

type MinerStakeRecord struct {
	Target common.Address
	Stake  *big.Int
}

type LockRewardRecord struct {
	Target   common.Address
	Amount   *big.Int
	IsReward uint32
	FlowValue1 uint64 `rlp:"optional"` //Real Flow Value
	FlowValue2 uint64 `rlp:"optional"` //valid Flow Value
}

type MinerPofReportItem struct {
	Target       common.Address
	FlowValue1   uint64
	FlowValue2   uint64
	Miner        common.Address
}

type MinerPofReportRecord struct {
	ChainHash     common.Hash
	ReportTime    uint64
	ReportContent []MinerPofReportItem
}

type ConfigDepositRecord struct {
	Who    uint32
	Amount *big.Int
}

// HeaderExtra is the struct of info in header.Extra[extraVanity:len(header.extra)-extraSeal]
// HeaderExtra is the current struct
type HeaderExtra struct {
	CurrentBlockConfirmations []Confirmation
	CurrentBlockVotes         []Vote
	CurrentBlockProposals     []Proposal
	CurrentBlockDeclares      []Declare
	ModifyPredecessorVotes    []Vote
	LoopStartTime             uint64
	SignerQueue               []common.Address
	SignerMissing             []common.Address
	ConfirmedBlockNumber      uint64
	SideChainConfirmations    []SCConfirmation
	SideChainSetCoinbases     []SCSetCoinbase
	SideChainNoticeConfirmed  []SCConfirmation
	SideChainCharging         []GasCharging //This only exist in side chain's header.Extra

	ExchangeCoin     []ExchangeCoinRecord
	DeviceBind       []DeviceBindRecord
	CandidatePunish  []CandidatePunishRecord
	MinerStake       []MinerStakeRecord
	CandidateExit    []common.Address
	ClaimedBandwidth []ClaimedBandwidthRecord
	PofMinerExit     []common.Address
	ConfigExchRate   uint32
	ConfigOffLine    uint32
	ConfigDeposit    []ConfigDepositRecord
	ConfigISPQOS     []ISPQOSRecord
	LockParameters   []LockParameterRecord
	ManagerAddress   []ManagerAddressRecord
	PofHarvest      *big.Int
	LockReward       []LockRewardRecord
	GrantProfit      []consensus.GrantProfitRecord
	PofReport       []MinerPofReportRecord
	CoinDataRoot    common.Hash
	GrantProfitHash common.Hash
	CandidatePledgeNew []CandidatePledgeNewRecord
	CandidatePledgeEntrust []CandidatePledgeEntrustRecord
	CandidatePEntrustExit  []CandidatePEntrustExitRecord
	CandidateAutoExit      []common.Address
	CandidateChangeRate    []CandidateChangeRateRecord
	PofPledgeReq     [] PofPledgeReq
	PofMinerPriceReq [] PofMinerPriceRecord
	InspireHarvest      *big.Int
	CandidateChangeManager []CandidateChangeManagerRecord
}
//side chain related
var minSCSetCoinbaseValue = big.NewInt(5e+18)

func (p *Proposal) copy() *Proposal {
	cpy := &Proposal{
		Hash:                   p.Hash,
		ReceivedNumber:         new(big.Int).Set(p.ReceivedNumber),
		CurrentDeposit:         new(big.Int).Set(p.CurrentDeposit),
		ValidationLoopCnt:      p.ValidationLoopCnt,
		ProposalType:           p.ProposalType,
		Proposer:               p.Proposer,
		TargetAddress:          p.TargetAddress,
		MinerRewardPerThousand: p.MinerRewardPerThousand,
		SCHash:                 p.SCHash,
		SCBlockCountPerPeriod:  p.SCBlockCountPerPeriod,
		SCBlockRewardPerPeriod: p.SCBlockRewardPerPeriod,
		Declares:               make([]*Declare, len(p.Declares)),
		MinVoterBalance:        p.MinVoterBalance,
		ProposalDeposit:        p.ProposalDeposit,
		SCRentFee:              p.SCRentFee,
		SCRentRate:             p.SCRentRate,
		SCRentLength:           p.SCRentLength,
	}

	copy(cpy.Declares, p.Declares)
	return cpy
}

func (s *SCConfirmation) copy() *SCConfirmation {
	cpy := &SCConfirmation{
		Hash:     s.Hash,
		Coinbase: s.Coinbase,
		Number:   s.Number,
		LoopInfo: make([]string, len(s.LoopInfo)),
	}
	copy(cpy.LoopInfo, s.LoopInfo)
	return cpy
}

// Encode HeaderExtra
func encodeHeaderExtra(config *params.AlienConfig, number *big.Int, val HeaderExtra) ([]byte, error) {

	var headerExtra interface{}
	switch {
	//case config.IsTrantor(number):

	default:
		headerExtra = val
	}
	return rlp.EncodeToBytes(headerExtra)

}

// Decode HeaderExtra
func decodeHeaderExtra(config *params.AlienConfig, number *big.Int, b []byte, val *HeaderExtra) error {
	var err error
	switch {
	//case config.IsTrantor(number):
	default:
		err = rlp.DecodeBytes(b, val)
	}
	return err
}
// Build side chain confirm data
func (a *Alien) buildSCEventConfirmData(scHash common.Hash, headerNumber *big.Int, headerTime *big.Int, lastLoopInfo string, chargingInfo string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:%s:%s:%d:%d:%s:%s",
		ufoPrefix, ufoVersion, ufoCategorySC, ufoEventConfirm,
		scHash.Hex(), headerNumber.Uint64(), headerTime.Uint64(), lastLoopInfo, chargingInfo))

}

// Calculate Votes from transaction in this block, write into header.Extra
func (a *Alien) processCustomTx(headerExtra HeaderExtra, chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (HeaderExtra, RefundGas, error) {
	// if predecessor voter make transaction and vote in this block,
	// just process as vote, do it in snapshot.apply
	var (
		snap       *Snapshot
		snapCache  *Snapshot
		err        error
		number     uint64
		refundGas  RefundGas
		refundHash RefundHash
	)
	refundGas = make(map[common.Address]*big.Int)
	refundHash = make(map[common.Hash]RefundPair)
	number = header.Number.Uint64()
	if number >= 1 {
		snap, err = a.snapshot(chain, number-1, header.ParentHash, nil, nil, defaultLoopCntRecalculateSigners)
		if err != nil {
			return headerExtra, nil, err
		}
		snapCache = snap.copy()
	}
	coinBalances := make(map[common.Address]*big.Int)
	for _, tx := range txs {
		txSender, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
		if err != nil {
			continue
		}

		if len(string(tx.Data())) >= len(ufoPrefix) {
			txData := string(tx.Data())
			txDataInfo := strings.Split(txData, ":")
			if len(txDataInfo) >= ufoMinSplitLen {
				if txDataInfo[posPrefix] == ufoPrefix {
					if txDataInfo[posVersion] == ufoVersion {
						// process vote event
						if txDataInfo[posCategory] == ufoCategoryEvent {
							if len(txDataInfo) > ufoMinSplitLen {
								// check is vote or not
								if txDataInfo[posEventVote] == ufoEventVote && (!candidateNeedPD || snap.isCandidate(*tx.To())) && state.GetBalance(txSender).Cmp(snap.MinVB) > 0 {
									headerExtra.CurrentBlockVotes = a.processEventVote(headerExtra.CurrentBlockVotes, state, tx, txSender)
								} else if txDataInfo[posEventConfirm] == ufoEventConfirm && snap.isCandidate(txSender) {
									headerExtra.CurrentBlockConfirmations, refundHash = a.processEventConfirm(headerExtra.CurrentBlockConfirmations, chain, txDataInfo, number, tx, txSender, refundHash)
								} else if txDataInfo[posEventProposal] == ufoEventPorposal {
									headerExtra.CurrentBlockProposals = a.processEventProposal(headerExtra.CurrentBlockProposals, txDataInfo, state, tx, txSender, snap)
								} else if txDataInfo[posEventDeclare] == ufoEventDeclare && snap.isCandidate(txSender) {
									headerExtra.CurrentBlockDeclares = a.processEventDeclare(headerExtra.CurrentBlockDeclares, txDataInfo, tx, txSender)
								}
							} else {
								// todo : something wrong, leave this transaction to process as normal transaction
							}
						} else if txDataInfo[posCategory] == ufoCategoryLog {
							// todo :
						} else if txDataInfo[posCategory] == ufoCategorySC {
							if len(txDataInfo) > ufoMinSplitLen {
								if txDataInfo[posEventConfirm] == ufoEventConfirm {
									if len(txDataInfo) > ufoMinSplitLen+5 {
										number := new(big.Int)
										if err := number.UnmarshalText([]byte(txDataInfo[ufoMinSplitLen+2])); err != nil {
											log.Trace("Side chain confirm info fail", "number", txDataInfo[ufoMinSplitLen+2])
											continue
										}
										if err := new(big.Int).UnmarshalText([]byte(txDataInfo[ufoMinSplitLen+3])); err != nil {
											log.Trace("Side chain confirm info fail", "time", txDataInfo[ufoMinSplitLen+3])
											continue
										}
										loopInfo := txDataInfo[ufoMinSplitLen+4]
										scHash := common.HexToHash(txDataInfo[ufoMinSplitLen+1])
										headerExtra.SideChainConfirmations, refundHash = a.processSCEventConfirm(headerExtra.SideChainConfirmations,
											scHash, number.Uint64(), loopInfo, tx, txSender, refundHash)

										chargingInfo := txDataInfo[ufoMinSplitLen+5]
										headerExtra.SideChainNoticeConfirmed = a.processSCEventNoticeConfirm(headerExtra.SideChainNoticeConfirmed,
											scHash, number.Uint64(), chargingInfo, txSender)

									}
								}
							}
						}
					}
				} else if txDataInfo[posPrefix] == tokenPrefix {
					if txDataInfo[posVersion] == ufoVersion {
						if txDataInfo[posCategory] == tokenCategoryExch {
							headerExtra.ExchangeCoin = a.processExchangeCoin(headerExtra.ExchangeCoin, txDataInfo, txSender, tx, receipts, state, snap)
						} else if txDataInfo[posCategory] == tokenCategoryBind {
							headerExtra.DeviceBind = a.processDeviceBind (headerExtra.DeviceBind, txDataInfo, txSender, tx, receipts, snapCache)
						} else if txDataInfo[posCategory] == tokenCategoryUnbind {
							headerExtra.DeviceBind = a.processDeviceUnbind (headerExtra.DeviceBind, txDataInfo, txSender, tx, receipts, state, snapCache)
						} else if txDataInfo[posCategory] == tokenCategoryRebind {
							headerExtra.DeviceBind = a.processDeviceRebind (headerExtra.DeviceBind, txDataInfo, txSender, tx, receipts, state, snapCache)
						}else if txDataInfo[posCategory] == tokenCategoryCandPnsh {
							headerExtra.CandidatePunish = a.processCandidatePunish (headerExtra.CandidatePunish, txDataInfo, txSender, tx, receipts, state, snapCache)
						}
						headerExtra=a.processPofCustomTx(txDataInfo,headerExtra,txSender, tx, receipts, snapCache, header.Number,state,chain, coinBalances)
						headerExtra=a.processPosCustomTx(txDataInfo,headerExtra,txSender, tx, receipts, snapCache, header.Number,state,chain, coinBalances)

					}
				}  else if txDataInfo[posPrefix] == sscPrefix {
					if txDataInfo[posVersion] == ufoVersion {
						if txDataInfo[posCategory] == sscCategoryExchRate {
							headerExtra.ConfigExchRate = a.processExchRate (txDataInfo, txSender, snapCache,tx,receipts)
						} else if txDataInfo[posCategory] == sscCategoryDeposit {
							headerExtra.ConfigDeposit = a.processCandidateDeposit (headerExtra.ConfigDeposit, txDataInfo, txSender, snapCache,tx,receipts)
						} else if txDataInfo[posCategory] == sscCategoryCndLock {
							headerExtra.LockParameters = a.processCndLockConfig (headerExtra.LockParameters, txDataInfo, txSender, snapCache,tx,receipts)
						} else if txDataInfo[posCategory] == sscCategoryPofLock {
							headerExtra.LockParameters = a.processPofLockConfig(headerExtra.LockParameters, txDataInfo, txSender, snapCache,tx,receipts)
						} else if txDataInfo[posCategory] == sscCategoryRwdLock {
							headerExtra.LockParameters = a.processRwdLockConfig (headerExtra.LockParameters, txDataInfo, txSender, snapCache,tx,receipts)
						} else if txDataInfo[posCategory] == sscCategoryOffLine {
							headerExtra.ConfigOffLine = a.processOffLine (txDataInfo, txSender, snapCache,tx,receipts)
						}   else if txDataInfo[posCategory] == sscCategoryManager {
							headerExtra.ManagerAddress = a.processManagerAddress (headerExtra.ManagerAddress, txDataInfo, txSender, snapCache,tx,receipts)
						}
					}
				}
			}
		}
		// check each address
		if number > 1 {
			headerExtra.ModifyPredecessorVotes = a.processPredecessorVoter(headerExtra.ModifyPredecessorVotes, state, tx, txSender, snap)
		}
	}

	for _, receipt := range receipts {
		if pair, ok := refundHash[receipt.TxHash]; ok && receipt.Status == 1 {
			pair.GasPrice.Mul(pair.GasPrice, big.NewInt(int64(receipt.GasUsed)))
			refundGas = a.refundAddGas(refundGas, pair.Sender, pair.GasPrice)
		}
	}
	return headerExtra, refundGas, nil
}

func (a *Alien) refundAddGas(refundGas RefundGas, address common.Address, value *big.Int) RefundGas {
	if _, ok := refundGas[address]; ok {
		refundGas[address].Add(refundGas[address], value)
	} else {
		refundGas[address] = value
	}

	return refundGas
}

func (a *Alien) processSCEventNoticeConfirm(scEventNoticeConfirm []SCConfirmation, hash common.Hash, number uint64, chargingInfo string, txSender common.Address) []SCConfirmation {
	if chargingInfo != "" {
		scEventNoticeConfirm = append(scEventNoticeConfirm, SCConfirmation{
			Hash:     hash,
			Coinbase: txSender,
			Number:   number,
			LoopInfo: strings.Split(chargingInfo, "#"),
		})
	}
	return scEventNoticeConfirm
}

func (a *Alien) processSCEventConfirm(scEventConfirmaions []SCConfirmation, hash common.Hash, number uint64, loopInfo string, tx *types.Transaction, txSender common.Address, refundHash RefundHash) ([]SCConfirmation, RefundHash) {
	scEventConfirmaions = append(scEventConfirmaions, SCConfirmation{
		Hash:     hash,
		Coinbase: txSender,
		Number:   number,
		LoopInfo: strings.Split(loopInfo, "#"),
	})
	refundHash[tx.Hash()] = RefundPair{txSender, tx.GasPrice()}
	return scEventConfirmaions, refundHash
}

func (a *Alien) processEventProposal(currentBlockProposals []Proposal, txDataInfo []string, state *state.StateDB, tx *types.Transaction, proposer common.Address, snap *Snapshot) []Proposal {
	// sample for add side chain proposal
	// eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:0,data:web3.toHex("ufo:1:event:proposal:proposal_type:4:sccount:2:screward:50:schash:0x3210000000000000000000000000000000000000000000000000000000000000:vlcnt:4")})
	// sample for declare
	// eth.sendTransaction({from:eth.accounts[0],to:eth.accounts[0],value:0,data:web3.toHex("ufo:1:event:declare:hash:0x853e10706e6b9d39c5f4719018aa2417e8b852dec8ad18f9c592d526db64c725:decision:yes")})
	if len(txDataInfo) <= posEventProposal+2 {
		return currentBlockProposals
	}

	proposal := Proposal{
		Hash:                   tx.Hash(),
		ReceivedNumber:         big.NewInt(0),
		CurrentDeposit:         proposalDeposit, // for all type of deposit
		ValidationLoopCnt:      defaultValidationLoopCnt,
		ProposalType:           proposalTypeCandidateAdd,
		Proposer:               proposer,
		TargetAddress:          common.Address{},
		SCHash:                 common.Hash{},
		SCBlockCountPerPeriod:  1,
		SCBlockRewardPerPeriod: 0,
		MinerRewardPerThousand: minerRewardPerThousand,
		Declares:               []*Declare{},
		MinVoterBalance:        new(big.Int).Div(minVoterBalance, big.NewInt(1e+18)).Uint64(),
		ProposalDeposit:        new(big.Int).Div(proposalDeposit, big.NewInt(1e+18)).Uint64(), // default value
		SCRentFee:              0,
		SCRentRate:             1,
		SCRentLength:           defaultSCRentLength,
	}

	for i := 0; i < len(txDataInfo[posEventProposal+1:])/2; i++ {
		k, v := txDataInfo[posEventProposal+1+i*2], txDataInfo[posEventProposal+2+i*2]
		switch k {
		case "vlcnt":
			// If vlcnt is missing then user default value, but if the vlcnt is beyond the min/max value then ignore this proposal
			if validationLoopCnt, err := strconv.Atoi(v); err != nil || validationLoopCnt < minValidationLoopCnt || validationLoopCnt > maxValidationLoopCnt {
				return currentBlockProposals
			} else {
				proposal.ValidationLoopCnt = uint64(validationLoopCnt)
			}
		case "schash":
			proposal.SCHash.UnmarshalText([]byte(v))
		case "sccount":
			if scBlockCountPerPeriod, err := strconv.Atoi(v); err != nil {
				return currentBlockProposals
			} else {
				proposal.SCBlockCountPerPeriod = uint64(scBlockCountPerPeriod)
			}
		case "screward":
			if scBlockRewardPerPeriod, err := strconv.Atoi(v); err != nil {
				return currentBlockProposals
			} else {
				proposal.SCBlockRewardPerPeriod = uint64(scBlockRewardPerPeriod)
			}
		case "proposal_type":
			if proposalType, err := strconv.Atoi(v); err != nil {
				return currentBlockProposals
			} else {
				proposal.ProposalType = uint64(proposalType)
			}
		case "candidate":
			// not check here
			proposal.TargetAddress.UnmarshalText([]byte(v))
		case "mrpt":
			// miner reward per thousand
			if mrpt, err := strconv.Atoi(v); err != nil || mrpt <= 0 || mrpt > 1000 {
				return currentBlockProposals
			} else {
				proposal.MinerRewardPerThousand = uint64(mrpt)
			}
		case "mvb":
			// minVoterBalance
			if mvb, err := strconv.Atoi(v); err != nil || mvb <= 0 {
				return currentBlockProposals
			} else {
				proposal.MinVoterBalance = uint64(mvb)
			}
		case "mpd":
			// proposalDeposit
			if mpd, err := strconv.Atoi(v); err != nil || mpd <= 0 || mpd > maxProposalDeposit {
				return currentBlockProposals
			} else {
				proposal.ProposalDeposit = uint64(mpd)
			}
		case "scrt":
			// target address on side chain to charge gas
			proposal.TargetAddress.UnmarshalText([]byte(v))
		case "scrf":
			// side chain rent fee
			if scrf, err := strconv.Atoi(v); err != nil || scrf < minSCRentFee {
				return currentBlockProposals
			} else {
				proposal.SCRentFee = uint64(scrf)
			}
		case "scrr":
			// side chain rent rate
			if scrr, err := strconv.Atoi(v); err != nil || scrr <= 0 {
				return currentBlockProposals
			} else {
				proposal.SCRentRate = uint64(scrr)
			}
		case "scrl":
			// side chain rent length
			if scrl, err := strconv.Atoi(v); err != nil || scrl < minSCRentLength || scrl > maxSCRentLength {
				return currentBlockProposals
			} else {
				proposal.SCRentLength = uint64(scrl)
			}
		}
	}
	// now the proposal is built
	currentProposalPay := new(big.Int).Set(proposalDeposit)
	if proposal.ProposalType == proposalTypeRentSideChain {
		// check if the proposal target side chain exist
		if !snap.isSideChainExist(proposal.SCHash) {
			return currentBlockProposals
		}
		if (proposal.TargetAddress == common.Address{}) {
			return currentBlockProposals
		}
		currentProposalPay.Add(currentProposalPay, new(big.Int).Mul(new(big.Int).SetUint64(proposal.SCRentFee), big.NewInt(1e+18)))
	}
	// check enough balance for deposit
	if state.GetBalance(proposer).Cmp(currentProposalPay) < 0 {
		return currentBlockProposals
	}
	// collection the fee for this proposal (deposit and other fee , sc rent fee ...)
	state.SetBalance(proposer, new(big.Int).Sub(state.GetBalance(proposer), currentProposalPay))

	return append(currentBlockProposals, proposal)
}

func (a *Alien) processEventDeclare(currentBlockDeclares []Declare, txDataInfo []string, tx *types.Transaction, declarer common.Address) []Declare {
	if len(txDataInfo) <= posEventDeclare+2 {
		return currentBlockDeclares
	}
	declare := Declare{
		ProposalHash: common.Hash{},
		Declarer:     declarer,
		Decision:     true,
	}
	for i := 0; i < len(txDataInfo[posEventDeclare+1:])/2; i++ {
		k, v := txDataInfo[posEventDeclare+1+i*2], txDataInfo[posEventDeclare+2+i*2]
		switch k {
		case "hash":
			declare.ProposalHash.UnmarshalText([]byte(v))
		case "decision":
			if v == "yes" {
				declare.Decision = true
			} else if v == "no" {
				declare.Decision = false
			} else {
				return currentBlockDeclares
			}
		}
	}

	return append(currentBlockDeclares, declare)
}

func (a *Alien) processEventVote(currentBlockVotes []Vote, state *state.StateDB, tx *types.Transaction, voter common.Address) []Vote {

	a.lock.RLock()
	stake := state.GetBalance(voter)
	a.lock.RUnlock()

	currentBlockVotes = append(currentBlockVotes, Vote{
		Voter:     voter,
		Candidate: *tx.To(),
		Stake:     stake,
	})

	return currentBlockVotes
}

func (a *Alien) processEventConfirm(currentBlockConfirmations []Confirmation, chain consensus.ChainHeaderReader, txDataInfo []string, number uint64, tx *types.Transaction, confirmer common.Address, refundHash RefundHash) ([]Confirmation, RefundHash) {
	if len(txDataInfo) > posEventConfirmNumber {
		confirmedBlockNumber := new(big.Int)
		err := confirmedBlockNumber.UnmarshalText([]byte(txDataInfo[posEventConfirmNumber]))
		if err != nil || number-confirmedBlockNumber.Uint64() > a.config.MaxSignerCount || number-confirmedBlockNumber.Uint64() < 0 {
			return currentBlockConfirmations, refundHash
		}
		// check if the voter is in block
		confirmedHeader := chain.GetHeaderByNumber(confirmedBlockNumber.Uint64())
		if confirmedHeader == nil {
			//log.Info("Fail to get confirmedHeader")
			return currentBlockConfirmations, refundHash
		}
		confirmedHeaderExtra := HeaderExtra{}
		if extraVanity+extraSeal > len(confirmedHeader.Extra) {
			return currentBlockConfirmations, refundHash
		}
		err = decodeHeaderExtra(a.config, confirmedBlockNumber, confirmedHeader.Extra[extraVanity:len(confirmedHeader.Extra)-extraSeal], &confirmedHeaderExtra)
		if err != nil {
			log.Info("Fail to decode parent header", "err", err)
			return currentBlockConfirmations, refundHash
		}
		for _, s := range confirmedHeaderExtra.SignerQueue {
			if s == confirmer {
				currentBlockConfirmations = append(currentBlockConfirmations, Confirmation{
					Signer:      confirmer,
					BlockNumber: new(big.Int).Set(confirmedBlockNumber),
				})
				refundHash[tx.Hash()] = RefundPair{confirmer, tx.GasPrice()}
				break
			}
		}
	}

	return currentBlockConfirmations, refundHash
}

func (a *Alien) processPredecessorVoter(modifyPredecessorVotes []Vote, state *state.StateDB, tx *types.Transaction, voter common.Address, snap *Snapshot) []Vote {
	// process normal transaction which relate to voter
	if tx.Value().Cmp(big.NewInt(0)) > 0 && tx.To() != nil {
		if snap.isVoter(voter) {
			a.lock.RLock()
			stake := state.GetBalance(voter)
			a.lock.RUnlock()
			modifyPredecessorVotes = append(modifyPredecessorVotes, Vote{
				Voter:     voter,
				Candidate: common.Address{},
				Stake:     stake,
			})
		}
		if snap.isVoter(*tx.To()) {
			a.lock.RLock()
			stake := state.GetBalance(*tx.To())
			a.lock.RUnlock()
			modifyPredecessorVotes = append(modifyPredecessorVotes, Vote{
				Voter:     *tx.To(),
				Candidate: common.Address{},
				Stake:     stake,
			})
		}

	}
	return modifyPredecessorVotes
}

func (a *Alien) addCustomerTxLog (tx *types.Transaction, receipts []*types.Receipt, topics []common.Hash, data []byte) bool {
	for _, receipt := range receipts {
		if receipt.TxHash != tx.Hash() {
			continue
		}
		if receipt.Status == types.ReceiptStatusFailed {
			return false
		}
		log := &types.Log{
			Address: common.Address{},
			Topics:  topics,
			Data:    data,
			BlockNumber: receipt.BlockNumber.Uint64(),
			TxHash: tx.Hash(),
			TxIndex: receipt.TransactionIndex,
			BlockHash: receipt.BlockHash,
			Index: uint(len(receipt.Logs)),
			Removed: false,
		}
		receipt.Logs = append(receipt.Logs, log)
		receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
		return true
	}
	return false
}

func (a *Alien) processExchangeCoin(currentExchangeCoin []ExchangeCoinRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) []ExchangeCoinRecord {
	if len(txDataInfo) <= tokenPosExchValue {
		log.Warn("Exchange TOKEN to COIN fail", "parameter number", len(txDataInfo))
		return currentExchangeCoin
	}
	exchangeCoin := ExchangeCoinRecord{
		Target: common.Address{},
		Amount: big.NewInt(0),
	}
	if err := exchangeCoin.Target.UnmarshalText1([]byte(txDataInfo[tokenPosExchAddress])); err != nil {
		log.Warn("Exchange TOKEN to COIN fail", "address", txDataInfo[tokenPosExchAddress])
		return currentExchangeCoin
	}
	amount := big.NewInt(0)
	var err error
	if amount, err = hexutil.UnmarshalText1([]byte(txDataInfo[tokenPosExchValue])); err != nil {
		log.Warn("Exchange TOKEN to COIN fail", "number", txDataInfo[tokenPosExchValue])
		return currentExchangeCoin
	}
	if state.GetBalance(txSender).Cmp(amount) < 0 {
		log.Warn("Exchange TOKEN to COIN fail", "balance", state.GetBalance(txSender))
		return currentExchangeCoin
	}
	exchangeCoin.Amount = new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(int64(snap.SystemConfig.ExchRate))),big.NewInt(10000))
	state.SetBalance(txSender, new(big.Int).Sub(state.GetBalance(txSender), amount))
	state.AddBalance(common.BigToAddress(big.NewInt(0)),amount)
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xdd6398517e51250c7ea4c550bdbec4246ce3cd80eac986e8ebbbb0eda27dcf4c")) //web3.sha3("ExchangeCoin(address,uint256)")
	topics[1].SetBytes(txSender.Bytes())
	topics[2].SetBytes(exchangeCoin.Target.Bytes())
	dataList := make([]common.Hash, 2)
	dataList[0].SetBytes(amount.Bytes())
	dataList[1].SetBytes(exchangeCoin.Amount.Bytes())
	data := dataList[0].Bytes()
	data = append(data, dataList[1].Bytes()...)
	a.addCustomerTxLog (tx, receipts, topics, data)
	currentExchangeCoin = append(currentExchangeCoin, exchangeCoin)
	return currentExchangeCoin
}

func (a *Alien) processDeviceBind (currentDeviceBind []DeviceBindRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, snap *Snapshot) []DeviceBindRecord {
	if len(txDataInfo) <= tokenPosMiltiSign {
		log.Warn("Device bind Revenue", "parameter number", len(txDataInfo))
		return currentDeviceBind
	}
	deviceBind := DeviceBindRecord {
		Device: common.Address{},
		Revenue: txSender,
		Contract: common.Address{},
		MultiSign: common.Address{},
		Type: 0,
		Bind: true,
	}
	if err := deviceBind.Device.UnmarshalText1([]byte(txDataInfo[tokenPosMinerAddress])); err != nil {
		log.Warn("Device bind Revenue", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentDeviceBind
	}
	if revenueType, err := strconv.ParseUint(txDataInfo[tokenPosRevenueType], 10, 32); err == nil {
		if revenueType == 0 {
			if _, ok := snap.RevenueNormal[deviceBind.Device]; ok {
				log.Warn("Device bind Revenue", "device already bond", txDataInfo[tokenPosMinerAddress])
				return currentDeviceBind
			}
			if !a.isPosManager(snap,deviceBind,txSender,txDataInfo){
				return currentDeviceBind
			}
		} else {
			if _, ok := snap.RevenuePof[deviceBind.Device]; ok {
				log.Warn("Device bind Revenue", "device already bond", txDataInfo[tokenPosMinerAddress])
				return currentDeviceBind
			}
			if !a.isPofManager(snap,deviceBind,txSender,txDataInfo){
				return currentDeviceBind
			}
		}
		deviceBind.Type = uint32(revenueType)
	} else {
		log.Warn("Device bind Revenue", "type", txDataInfo[tokenPosRevenueType])
		return currentDeviceBind
	}
	if len(txDataInfo) > tokenPosRevenueAddress {
		if 0 < len(txDataInfo[tokenPosRevenueAddress]) {
			if err := deviceBind.Revenue.UnmarshalText1([]byte(txDataInfo[tokenPosRevenueAddress])); err != nil {
				log.Warn("Device bind revenue", "Revenue address", txDataInfo[tokenPosRevenueAddress])
				return currentDeviceBind
			}
		}
	}
	if err := a.checkRevenueNormalBind(deviceBind,snap); err != nil {
		log.Warn("Device bind Revenue", "checkRevenueNormalBind", err.Error())
		return currentDeviceBind
	}

	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xf061654231b0035280bd8dd06084a38aa871445d0b7311be8cc2605c5672a6e3")) //web3.sha3("DeviceBind(uint32,byte32,byte32,address)")
	topics[1].SetBytes(deviceBind.Device.Bytes())
	topics[2].SetBytes(big.NewInt(int64(deviceBind.Type)).Bytes())
	dataList := make([]common.Hash, 3)
	dataList[0].SetBytes(deviceBind.Revenue.Bytes())
	dataList[1] = deviceBind.Contract.Hash()
	dataList[2] = deviceBind.MultiSign.Hash()
	data := dataList[0].Bytes()
	data = append(data, dataList[1].Bytes()...)
	data = append(data, dataList[2].Bytes()...)
	a.addCustomerTxLog (tx, receipts, topics, data)
	currentDeviceBind = append (currentDeviceBind, deviceBind)
	if deviceBind.Type == 0 {
		snap.RevenueNormal[deviceBind.Device] = &RevenueParameter{
			RevenueAddress: deviceBind.Revenue,
			RevenueContract: deviceBind.Contract,
			MultiSignature: deviceBind.MultiSign,
		}
	} else {
		snap.RevenuePof[deviceBind.Device] = &RevenueParameter{
			RevenueAddress: deviceBind.Revenue,
			RevenueContract: deviceBind.Contract,
			MultiSignature: deviceBind.MultiSign,
		}
	}
	return currentDeviceBind
}

func (a *Alien) processDeviceUnbind (currentDeviceBind []DeviceBindRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) []DeviceBindRecord {
	if len(txDataInfo) <= tokenPosRevenueType {
		log.Warn("Device unbind Revenue", "parameter number", len(txDataInfo))
		return currentDeviceBind
	}
	deviceBind := DeviceBindRecord {
		Device: common.Address{},
		Revenue: common.Address{},
		Contract: common.Address{},
		MultiSign: common.Address{},
		Type: 0,
		Bind: false,
	}
	if err := deviceBind.Device.UnmarshalText1([]byte(txDataInfo[tokenPosMinerAddress])); err != nil {
		log.Warn("Device unbind Revenue", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentDeviceBind
	}
	if revenueType, err := strconv.ParseUint(txDataInfo[tokenPosRevenueType], 10, 32); err == nil {
		if revenueType == 0 {
			if _, ok := snap.RevenueNormal[deviceBind.Device]; !ok {
				log.Warn("Device unbind Revenue", "device never bond", txDataInfo[tokenPosMinerAddress])
				return currentDeviceBind
			} else {
				if !a.isPosManager(snap,deviceBind,txSender,txDataInfo){
					return currentDeviceBind
				}
			}
		} else {
			if _, ok := snap.RevenuePof[deviceBind.Device]; !ok {
				log.Warn("Device unbind Revenue", "device never bond", txDataInfo[tokenPosMinerAddress])
				return currentDeviceBind
			} else {
				if !a.isPofManager(snap,deviceBind,txSender,txDataInfo) {
					return currentDeviceBind
				}
			}
		}
		deviceBind.Type = uint32(revenueType)
	} else {
		log.Warn("Device unbind Revenue", "type", txDataInfo[tokenPosRevenueType])
		return currentDeviceBind
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xf061654231b0035280bd8dd06084a38aa871445d0b7311be8cc2605c5672a6e3")) //web3.sha3("DeviceBind(uint32,byte32,byte32,address)")
	topics[1].SetBytes(deviceBind.Device.Bytes())
	topics[2].SetBytes(big.NewInt(int64(deviceBind.Type)).Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentDeviceBind = append (currentDeviceBind, deviceBind)
	return currentDeviceBind
}

func (a *Alien) processDeviceRebind (currentDeviceBind []DeviceBindRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) []DeviceBindRecord {
	if len(txDataInfo) <= tokenPosRevenueAddress {
		log.Warn("Device rebind Revenue", "parameter number", len(txDataInfo))
		return currentDeviceBind
	}
	deviceBind := DeviceBindRecord {
		Device: common.Address{},
		Revenue: txSender,
		Contract: common.Address{},
		MultiSign: common.Address{},
		Type: 0,
		Bind: true,
	}
	if err := deviceBind.Device.UnmarshalText1([]byte(txDataInfo[tokenPosMinerAddress])); err != nil {
		log.Warn("Device rebind Revenue", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentDeviceBind
	}
	if err := deviceBind.Revenue.UnmarshalText1([]byte(txDataInfo[tokenPosRevenueAddress])); err != nil {
		log.Warn("Device rebind Revenue", "Revenue address", txDataInfo[tokenPosMinerAddress])
		return currentDeviceBind
	}
	if revenueType, err := strconv.ParseUint(txDataInfo[tokenPosRevenueType], 10, 32); err == nil {
		if revenueType == 0 {
			if _, ok := snap.RevenueNormal[deviceBind.Device]; ok {
				if !a.isPosManager(snap,deviceBind,txSender,txDataInfo){
					return currentDeviceBind
				}
			} else {
				log.Warn("Device rebind Revenue", "device cnnnot bind", deviceBind.Revenue)
				return currentDeviceBind
			}
		} else {
			if _, ok := snap.RevenuePof[deviceBind.Device]; ok {
				if !a.isPofManager(snap,deviceBind,txSender,txDataInfo) {
					return currentDeviceBind
				}
			} else {
				log.Warn("Device rebind Revenue", "device cnnnot bind", deviceBind.Revenue)
				return currentDeviceBind
			}
		}
		deviceBind.Type = uint32(revenueType)
	} else {
		log.Warn("Device rebind Revenue", "type", txDataInfo[tokenPosRevenueType])
		return currentDeviceBind
	}
	if err := a.checkRevenueNormalBind(deviceBind,snap); err != nil {
		log.Warn("Device rebind Revenue", "checkRevenueNormalBind", err.Error())
		return currentDeviceBind
	}

	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xf061654231b0035280bd8dd06084a38aa871445d0b7311be8cc2605c5672a6e3")) //web3.sha3("DeviceBind(uint32,byte32,byte32,address)")
	topics[1].SetBytes(deviceBind.Device.Bytes())
	topics[2].SetBytes(big.NewInt(int64(deviceBind.Type)).Bytes())
	dataList := make([]common.Hash, 3)
	dataList[0].SetBytes(deviceBind.Revenue.Bytes())
	dataList[1] = deviceBind.Contract.Hash()
	dataList[2] = deviceBind.MultiSign.Hash()
	data := dataList[0].Bytes()
	data = append(data, dataList[1].Bytes()...)
	data = append(data, dataList[2].Bytes()...)
	a.addCustomerTxLog (tx, receipts, topics, data)
	currentDeviceBind = append (currentDeviceBind, deviceBind)
	return currentDeviceBind
}

func (a *Alien) processCandidatePunish (currentCandidatePunish []CandidatePunishRecord, txDataInfo []string, txSender common.Address, tx *types.Transaction, receipts []*types.Receipt, state *state.StateDB, snap *Snapshot) []CandidatePunishRecord {
	if len(txDataInfo) <= tokenPosMinerAddress {
		log.Warn("Candidate punish", "parameter number", len(txDataInfo))
		return currentCandidatePunish
	}
	candidatePunish := CandidatePunishRecord{
		Target: common.Address{},
		Amount: big.NewInt(0),
		Credit: 0,
	}
	if err := candidatePunish.Target.UnmarshalText1([]byte(txDataInfo[tokenPosMinerAddress])); err != nil {
		log.Warn("Candidate punish", "miner address", txDataInfo[tokenPosMinerAddress])
		return currentCandidatePunish
	}
	if candidateCredit, ok := snap.Punished[candidatePunish.Target]; !ok {
		log.Warn("Candidate punish", "not punish", candidatePunish.Target)
		return currentCandidatePunish
	} else {
		candidatePunish.Credit = uint32(candidateCredit)
		deposit := new(big.Int).Set(minCndPledgeBalance)
		if _, ok := snap.SystemConfig.Deposit[0]; ok {
			deposit = new(big.Int).Set(snap.SystemConfig.Deposit[0])
		}
		candidatePunish.Amount = new(big.Int).Div(new(big.Int).Mul(deposit, big.NewInt(int64(candidateCredit))), big.NewInt(int64(defaultFullCredit)))
	}
	if state.GetBalance(txSender).Cmp(candidatePunish.Amount) < 0 {
		log.Warn("Candidate punish", "balance", state.GetBalance(txSender))
		return currentCandidatePunish
	}
	if _, ok := snap.PosPledge[candidatePunish.Target]; !ok {
		log.Warn("Candidate punish", "candidate isnot exist", candidatePunish.Target)
		return currentCandidatePunish
	}
	state.SetBalance(txSender, new(big.Int).Sub(state.GetBalance(txSender), candidatePunish.Amount))
	state.AddBalance(common.BigToAddress(big.NewInt(0)),candidatePunish.Amount)
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xd67fe14bb06aa8656e0e7c3230831d68e8ce49bb4a4f71448f98a998d2674621")) //web3.sha3("PledgePunish(address,uint32)")
	topics[1].SetBytes(candidatePunish.Target.Bytes())
	topics[2].SetBytes(big.NewInt(sscEnumCndLock).Bytes())
	dataList := make([]common.Hash, 2)
	dataList[0].SetBytes(big.NewInt(int64(candidatePunish.Credit)).Bytes())
	dataList[1].SetBytes(candidatePunish.Amount.Bytes())
	data := dataList[0].Bytes()
	data = append(data, dataList[1].Bytes()...)
	a.addCustomerTxLog (tx, receipts, topics, data)
	currentCandidatePunish = append(currentCandidatePunish, candidatePunish)
	return currentCandidatePunish
}


func (a *Alien) processExchRate (txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) uint32 {
	if len(txDataInfo) <= sscPosExchRate {
		log.Warn("Config exchrate", "parameter number", len(txDataInfo))
		return 0
	}
	if exchRate, err := strconv.ParseUint(txDataInfo[sscPosExchRate], 10, 32); err != nil {
		log.Warn("Config exchrate", "exchrate", txDataInfo[sscPosExchRate])
		return 0
	} else {
		if snap.SystemConfig.ManagerAddress[sscEnumExchRate].String() != txSender.String() {
			log.Warn("Config exchrate", "manager address", txSender)
			return 0
		}
		topics := make([]common.Hash, 3)
		topics[0].UnmarshalText([]byte("0xd67fe14bb06aa8656e0e7c3230831d68e8ce49bb4a4f71448f98a998d2674628")) //web3.sha3("PledgePunish(address,uint32)")
		topics[1].SetBytes(txSender.Bytes())

		topics[2].SetBytes(txSender.Bytes())
		a.addCustomerTxLog (tx, receipts, topics, nil)
		return uint32(exchRate)
	}
}

func (a *Alien) processCandidateDeposit (currentDeposit []ConfigDepositRecord, txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) []ConfigDepositRecord {
	if len(txDataInfo) <= sscPosDepositWho {
		log.Warn("Config candidate deposit", "parameter number", len(txDataInfo))
		return currentDeposit
	}
	deposit := ConfigDepositRecord{
		Who: 0,
		Amount: big.NewInt(0),
	}
	var err error
	if deposit.Amount, err = hexutil.UnmarshalText1([]byte(txDataInfo[sscPosDeposit])); err != nil {
		log.Warn("Config candidate deposit", "deposit", txDataInfo[sscPosDeposit])
		return currentDeposit
	}
	if id, err := strconv.ParseUint(txDataInfo[sscPosDepositWho], 10, 32); err != nil {
		log.Warn("Config manager", "id", txDataInfo[sscPosDepositWho])
		return currentDeposit
	} else {
		deposit.Who = uint32(id)
	}
	if snap.SystemConfig.ManagerAddress[sscEnumSystem].String() != txSender.String() {
		log.Warn("Config candidate deposit", "manager address", txSender)
		return currentDeposit
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xd67fe14bb06aa8656e0e7c3230831d68e8ce49bb4a4f71448f98a998d2675424")) //web3.sha3("PledgePunish(address,uint32)")
	topics[1].SetBytes(txSender.Bytes())
	topics[2].SetBytes(deposit.Amount.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentDeposit = append(currentDeposit, deposit)
	return currentDeposit
}

func (a *Alien) processCndLockConfig (currentLockParameters []LockParameterRecord, txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) []LockParameterRecord {
	if len(txDataInfo) <= sscPosInterval {
		log.Warn("Config candidate lock", "parameter number", len(txDataInfo))
		return currentLockParameters
	}
	lockParameter := LockParameterRecord{
		Who: sscEnumCndLock,
		LockPeriod: uint32(30 * 24 * 60 * 60 / a.config.Period),
		RlsPeriod: uint32(210 * 24 * 60 * 60 / a.config.Period),
		Interval: uint32(24 * 60 * 60 / a.config.Period),
	}

	if lockPeriod, err := strconv.ParseUint(txDataInfo[sscPosLockPeriod], 16, 32); err != nil {
		log.Warn("Config candidate lock", "lock period", txDataInfo[sscPosLockPeriod])
		return currentLockParameters
	} else {
		lockParameter.LockPeriod = uint32(lockPeriod)
	}
	if releasePeriod, err := strconv.ParseUint(txDataInfo[sscPosRlsPeriod], 16, 32); err != nil {
		log.Warn("Config candidate lock", "release period", txDataInfo[sscPosRlsPeriod])
		return currentLockParameters
	} else {
		lockParameter.RlsPeriod = uint32(releasePeriod)
	}
	if interval, err := strconv.ParseUint(txDataInfo[sscPosInterval], 16, 32); err != nil {
		log.Warn("Config candidate lock", "release interval", txDataInfo[sscPosInterval])
		return currentLockParameters
	} else {
		lockParameter.Interval = uint32(interval)
	}
	if snap.SystemConfig.ManagerAddress[sscEnumSystem].String() != txSender.String() {
		log.Warn("Config candidate lock", "manager address", txSender)
		return currentLockParameters
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xd9629881a36a841c8f83ac60455b24494cb4cbaec14b18386a5f9a2a345ceec0")) //web3.sha3("PledgePunish(address,uint32)")
	topics[1].SetBytes(txSender.Bytes())
	topics[2].SetBytes(txSender.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentLockParameters = append(currentLockParameters, lockParameter)
	return currentLockParameters
}

func (a *Alien) processPofLockConfig(currentLockParameters []LockParameterRecord, txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) []LockParameterRecord {
	if len(txDataInfo) <= sscPosInterval {
		log.Warn("Config miner lock", "parameter number", len(txDataInfo))
		return currentLockParameters
	}
	lockParameter := LockParameterRecord{
		Who:        sscEnumPofLock,
		LockPeriod: uint32(30 * 24 * 60 * 60 / a.config.Period),
		RlsPeriod: uint32(210 * 24 * 60 * 60 / a.config.Period),
		Interval: uint32(24 * 60 * 60 / a.config.Period),
	}
	if lockPeriod, err := strconv.ParseUint(txDataInfo[sscPosLockPeriod], 16, 32); err != nil {
		log.Warn("Config miner lock", "lock period", txDataInfo[sscPosLockPeriod])
		return currentLockParameters
	} else {
		lockParameter.LockPeriod = uint32(lockPeriod)
	}
	if releasePeriod, err := strconv.ParseUint(txDataInfo[sscPosRlsPeriod], 16, 32); err != nil {
		log.Warn("Config miner lock", "release period", txDataInfo[sscPosRlsPeriod])
		return currentLockParameters
	} else {
		lockParameter.RlsPeriod = uint32(releasePeriod)
	}
	if interval, err := strconv.ParseUint(txDataInfo[sscPosInterval], 16, 32); err != nil {
		log.Warn("Config miner lock", "release interval", txDataInfo[sscPosInterval])
		return currentLockParameters
	} else {
		lockParameter.Interval = uint32(interval)
	}
	if snap.SystemConfig.ManagerAddress[sscEnumSystem].String() != txSender.String() {
		log.Warn("Config miner lock", "manager address", txSender)
		return currentLockParameters
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xe0ae901b00a40c79895a0a7647e4f4fa61de84d8934e86ae36a1c850e69e0196")) //web3.sha3("PledgePunish(address,uint32)")
	topics[1].SetBytes(txSender.Bytes())
	topics[2].SetBytes(txSender.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentLockParameters = append(currentLockParameters, lockParameter)
	return currentLockParameters
}

func (a *Alien) processRwdLockConfig (currentLockParameters []LockParameterRecord, txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) []LockParameterRecord {
	if len(txDataInfo) <= sscPosInterval {
		log.Warn("Config reward lock", "parameter number", len(txDataInfo))
		return currentLockParameters
	}
	lockParameter := LockParameterRecord{
		Who: sscEnumRwdLock,
		LockPeriod: uint32(30 * 24 * 60 * 60 / a.config.Period),
		RlsPeriod: uint32(210 * 24 * 60 * 60 / a.config.Period),
		Interval: uint32(24 * 60 * 60 / a.config.Period),
	}
	if lockPeriod, err := strconv.ParseUint(txDataInfo[sscPosLockPeriod], 16, 32); err != nil {
		log.Warn("Config reward lock", "lock period", txDataInfo[sscPosLockPeriod])
		return currentLockParameters
	} else {
		lockParameter.LockPeriod = uint32(lockPeriod)
	}
	if releasePeriod, err := strconv.ParseUint(txDataInfo[sscPosRlsPeriod], 16, 32); err != nil {
		log.Warn("Config reward lock", "release period", txDataInfo[sscPosRlsPeriod])
		return currentLockParameters
	} else {
		lockParameter.RlsPeriod = uint32(releasePeriod)
	}
	if interval, err := strconv.ParseUint(txDataInfo[sscPosInterval], 16, 32); err != nil {
		log.Warn("Config reward lock", "release interval", txDataInfo[sscPosInterval])
		return currentLockParameters
	} else {
		lockParameter.Interval = uint32(interval)
	}
	if snap.SystemConfig.ManagerAddress[sscEnumSystem].String() != txSender.String() {
		log.Warn("Config reward lock", "manager address", txSender)
		return currentLockParameters
	}
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0x4300e22752a84b16b2444fcdc36c9c510c66808994f669c642a9e67c3cca6841"))
	topics[1].SetBytes(txSender.Bytes())
	topics[2].SetBytes(txSender.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	currentLockParameters = append(currentLockParameters, lockParameter)
	return currentLockParameters
}

func (a *Alien) processOffLine (txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) uint32 {
	if len(txDataInfo) <= sscPosOffLine {
		log.Warn("Config offLine", "parameter number", len(txDataInfo))
		return 0
	}
	if offline, err := strconv.ParseUint(txDataInfo[sscPosOffLine], 10, 32); err != nil {
		log.Warn("Config offline", "offline", txDataInfo[sscPosOffLine])
		return 0
	} else {
		if snap.SystemConfig.ManagerAddress[sscEnumSystem].String() != txSender.String() {
			log.Warn("Config offLine", "manager address", txSender)
			return 0
		}
		topics := make([]common.Hash, 3)
		topics[0].UnmarshalText([]byte("0xf7b4afb0a90423ae38d2ffbd45ffcb4405df2c66d1a460535320f6f60cd5c3c3"))
		topics[1].SetBytes(txSender.Bytes())
		topics[2].SetBytes(txSender.Bytes())
		a.addCustomerTxLog (tx, receipts, topics, nil)
		return uint32(offline)
	}
}

func (a *Alien) processISPQos (currentISPQOS []ISPQOSRecord, txDataInfo []string, txSender common.Address, snap *Snapshot) []ISPQOSRecord {
	if len(txDataInfo) <= sscPosQosValue {
		log.Warn("Config isp qos", "parameter number", len(txDataInfo))
		return currentISPQOS
	}
	ISPQOS := ISPQOSRecord{
		ISPID: 0,
		QOS: 0,
	}
	if id, err := strconv.ParseUint(txDataInfo[sscPosQosID], 10, 32); err != nil {
		log.Warn("Config isp qos", "isp id", txDataInfo[sscPosQosID])
		return currentISPQOS
	} else {
		ISPQOS.ISPID = uint32(id)
	}
	if qos, err := strconv.ParseUint(txDataInfo[sscPosQosValue], 10, 32); err != nil {
		log.Warn("Config isp qos", "qos", txDataInfo[sscPosQosValue])
		return currentISPQOS
	} else {
		ISPQOS.QOS = uint32(qos)
	}
	if snap.SystemConfig.ManagerAddress[sscEnumSystem].String() != txSender.String() {
		log.Warn("Config isp qos", "manager address", txSender)
		return currentISPQOS
	}
	currentISPQOS = append(currentISPQOS, ISPQOS)
	return currentISPQOS
}

func (a *Alien) processManagerAddress (currentManagerAddress []ManagerAddressRecord, txDataInfo []string, txSender common.Address, snap *Snapshot,tx *types.Transaction, receipts []*types.Receipt) []ManagerAddressRecord {
	if len(txDataInfo) <= sscPosManagerAddress {
		log.Warn("Config manager", "parameter number", len(txDataInfo))
		return currentManagerAddress
	}
	if txSender.String() != managerAddressManager.String() {
		log.Warn("Config manager", "manager", txSender)
		return currentManagerAddress
	}
	managerAddress := ManagerAddressRecord{
		Target: common.Address{},
		Who: 0,
	}
	if id, err := strconv.ParseUint(txDataInfo[sscPosManagerID], 10, 32); err != nil {
		log.Warn("Config manager", "id", txDataInfo[sscPosManagerID])
		return currentManagerAddress
	} else {
		managerAddress.Who = uint32(id)
	}
	if err := managerAddress.Target.UnmarshalText1([]byte(txDataInfo[sscPosManagerAddress])); err != nil {
		log.Warn("Config manager", "address", txDataInfo[sscPosManagerAddress])
		return currentManagerAddress
	}
	snap.SystemConfig.ManagerAddress[managerAddress.Who] = managerAddress.Target
	currentManagerAddress = append(currentManagerAddress, managerAddress)
	topics := make([]common.Hash, 3)
	topics[0].UnmarshalText([]byte("0xeee4ff0e47c9256ad3906607ab3405f0a6885d004132930a6967156b8dfbb782"))
	topics[1].SetBytes(txSender.Bytes())
	topics[2].SetBytes(txSender.Bytes())
	a.addCustomerTxLog (tx, receipts, topics, nil)
	return currentManagerAddress
}

func (a *Alien) checkRevenueNormalBind (deviceBind DeviceBindRecord,snap *Snapshot) error {
	if deviceBind.Type == 0 {
		find := false
		for _, revenue := range snap.RevenueNormal {
			revenueAddress:=revenue.RevenueAddress
			if deviceBind.Revenue==revenueAddress {
				find = true
				break
			}
		}
		if find {
			return errors.New("revenueAddress is already bind a Normal device")
		}
	}
	return nil
}

func (a *Alien) isManagerAddressFlowReport(txSender common.Address, snap *Snapshot) bool {
	if txSender.String() == snap.SystemConfig.ManagerAddress[sscEnumFlowReport].String() {
		return true
	}
	return false
}

