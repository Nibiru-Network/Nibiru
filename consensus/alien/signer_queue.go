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
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"sort"

	"github.com/token/common"
)

type TallyItem struct {
	addr  common.Address
	stake *big.Int
}
type TallySlice []TallyItem

func (s TallySlice) Len() int      { return len(s) }
func (s TallySlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TallySlice) Less(i, j int) bool {
	//we need sort reverse, so ...
	isLess := s[i].stake.Cmp(s[j].stake)
	if isLess > 0 {
		return true

	} else if isLess < 0 {
		return false
	}
	// if the stake equal
	return bytes.Compare(s[i].addr.Bytes(), s[j].addr.Bytes()) > 0
}

type SignerItem struct {
	addr common.Address
	hash common.Hash
}
type SignerSlice []SignerItem

func (s SignerSlice) Len() int      { return len(s) }
func (s SignerSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SignerSlice) Less(i, j int) bool {
	isLess :=  bytes.Compare(s[i].hash.Bytes(), s[j].hash.Bytes())
	if isLess > 0 {
		return true
	} else if isLess < 0 {
		return false
	}
	// if the hash equal
	return bytes.Compare(s[i].addr.Bytes(), s[j].addr.Bytes()) > 0
}

// verify the SignerQueue base on block hash
func (s *Snapshot) verifySignerQueue(signerQueue []common.Address) error {

	if len(signerQueue) > int(s.config.MaxSignerCount) {
		return errInvalidSignerQueue
	}
	sq, err := s.createSignerQueue()
	if err != nil {
		return err
	}
	if len(sq) == 0 || len(sq) != len(signerQueue) {
		return errInvalidSignerQueue
	}
	for i, signer := range signerQueue {
		if signer != sq[i] {
			return errInvalidSignerQueue
		}
	}

	return nil
}

func (s *Snapshot) buildTallySlice() TallySlice {
	var tallySlice TallySlice
	for address, stake := range s.Tally {
		if !candidateNeedPD || s.isCandidate(address) {
			if s.Punished[address] < minCalSignerQueueCredit{
				tallySlice = append(tallySlice, TallyItem{address, stake})
			}

		}
	}
	if len(tallySlice) ==0 {
		for address, stake := range s.Tally {
			if !candidateNeedPD || s.isCandidate(address) {
				tallySlice = append(tallySlice, TallyItem{address, stake})
			}
		}
	}
	return tallySlice

}
func (s *Snapshot) buildTallyMiner() TallySlice {
	var tallySlice TallySlice
	for address, stake := range s.TallyMiner {
		if _, ok := s.PosPledge[address]; !ok  || s.Punished[address] >= minCalSignerQueueCredit {
			continue
		}
		if _, isok := s.Tally[address]; isok{
			continue
		}
		tallySlice = append(tallySlice, TallyItem{address, stake.Stake})
	}
	return tallySlice
}
func (s *Snapshot) createSignerQueue() ([]common.Address, error) {

	if (s.Number+1)%s.config.MaxSignerCount != 0 || s.Hash != s.HistoryHash[len(s.HistoryHash)-1] {
		fmt.Println("eeeeee",s.Number,s.Number+1,"MaxSignerCount",s.config.MaxSignerCount,(s.Number+1)%s.config.MaxSignerCount)
		return nil, errCreateSignerQueueNotAllowed
	}
	var signerSlice SignerSlice
	var topStakeAddress []common.Address

	if (s.Number+1)%(s.config.MaxSignerCount*s.LCRS) == 0 {

		mainMinerSlice := s.buildTallySlice()
		sort.Sort(TallySlice(mainMinerSlice))
		secondMinerSlice := s.buildTallyMiner()
		//sort.Sort(TallySlice(secondMinerSlice))
		queueLength := int(s.config.MaxSignerCount)
		mainSignerSliceLen := len(mainMinerSlice)
		if queueLength >= defaultOfficialMaxSignerCount {
			mainMinerNumber := (posCandidateOneNum * queueLength + defaultOfficialMaxSignerCount - 1) / defaultOfficialMaxSignerCount
			secondMinerNumber := posCandidateTwoNum * queueLength / defaultOfficialMaxSignerCount
			secondMinerSlice=s.selectSecondMinerSlice(secondMinerSlice,secondMinerNumber)
			sort.Sort(TallySlice(secondMinerSlice))
			secondLength := len(secondMinerSlice)
			if secondLength== 0 {
				mainMinerNumber = queueLength - secondLength
			}else if secondMinerNumber >= secondLength {
				signerSlice = s.selectSecondMinerInsufficient(secondMinerSlice, signerSlice,secondMinerNumber)
			} else {
				mainMinerNumber = queueLength - secondMinerNumber
				signerSlice = s.selectSecondMiner(secondMinerSlice, secondMinerNumber, signerSlice)
			}
			mainMinerSlice=s.selectMainMinerSlice(mainMinerSlice)
			sort.Sort(TallySlice(mainMinerSlice))
			mainSignerSliceLen = len(mainMinerSlice)
			// select Main Miner
			signerSlice = s.selectMainMiner(mainMinerNumber, mainSignerSliceLen, signerSlice, mainMinerSlice, secondMinerNumber)
		} else {
			if queueLength > len(mainMinerSlice) {
				queueLength = len(mainMinerSlice)
			}
			for i, tallyItem := range mainMinerSlice[:queueLength] {
				signerSlice = append(signerSlice, SignerItem{tallyItem.addr, s.HistoryHash[len(s.HistoryHash)-1-i]})
			}
		}
	} else {
		for i, signer := range s.Signers {
			signerSlice = append(signerSlice, SignerItem{*signer, s.HistoryHash[len(s.HistoryHash)-1-i]})
		}
	}

	// Set the top candidates in random order base on block hash
	sort.Sort(SignerSlice(signerSlice))
	if len(signerSlice) == 0 {
		return nil, errSignerQueueEmpty
	}
	for i := 0; i < int(s.config.MaxSignerCount); i++ {
		topStakeAddress = append(topStakeAddress, signerSlice[i%len(signerSlice)].addr)
	}
	return topStakeAddress, nil
}
func (s *Snapshot) selectSecondMinerSlice(candidatePledgeSlice TallySlice,secondMinerNumber int) TallySlice{
	if len(candidatePledgeSlice) <=secondMinerNumber {
		return  candidatePledgeSlice
	}
	maxAmount :=big.NewInt(0)
	minAmount :=big.NewInt(0)
	totalAmount :=big.NewInt(0)
	for _, item := range candidatePledgeSlice {
		if _, ok := s.TallyMiner[item.addr]; ok {
			totalAmount= new(big.Int).Add(totalAmount,item.stake)
			if maxAmount.Cmp(item.stake)==-1{
				maxAmount=new(big.Int).Set(item.stake)
			}
			if minAmount.Cmp(item.stake)==1{
				minAmount=new(big.Int).Set(item.stake)
			}
		}
	}
	avgAmount:=new(big.Int).Div(new(big.Int).Sub(new(big.Int).Sub(totalAmount,maxAmount),minAmount), big.NewInt(int64(len(candidatePledgeSlice)-2)))
	avgAmount=new(big.Int).Div(new(big.Int).Mul(avgAmount,posCandidateAvgRate),big.NewInt(100))//acgAmount * 0.75
	var preMiners TallySlice
	totalAmount =big.NewInt(0)
	for _, item := range candidatePledgeSlice {
		if _, ok := s.TallyMiner[item.addr]; ok {
			if item.stake.Cmp(avgAmount)>=0 {
				totalAmount= new(big.Int).Add(totalAmount,item.stake)
				preMiners = append(preMiners,item )
			}

		}
	}
	sort.Sort(preMiners)
	return preMiners
}

func (s *Snapshot) selectSecondMiner(candidatePledgeSlice TallySlice, secondMinerNumber int, signerSlice SignerSlice) SignerSlice {
	candidatePledgeSlice = s.reBuildMiner(candidatePledgeSlice)
	for i, tallyItem := range candidatePledgeSlice[:secondMinerNumber] {
		signerSlice = append(signerSlice, SignerItem{tallyItem.addr, s.HistoryHash[len(s.HistoryHash)-1-i]})
	}
	return signerSlice
}
func (s *Snapshot) selectSecondMinerInsufficient(tallyMiner TallySlice, signerSlice SignerSlice,secondMinerNumber int) SignerSlice {
	minerNum:= len(tallyMiner)
	catallyMiner:=s.reBuildMiner(tallyMiner)
	for i := 0; i < secondMinerNumber; i++ {
		signerSlice = append(signerSlice, SignerItem{catallyMiner[i%minerNum].addr, s.HistoryHash[len(s.HistoryHash)-1-i]})
	}
	return signerSlice
}
func (s *Snapshot) selectMainMinerSlice( miners TallySlice ) TallySlice{
	if len(miners) <= posCandidateOneNum{
		return miners
	}
	totalAmount :=big.NewInt(0)
	minAmount :=big.NewInt(0)
	maxAmount :=big.NewInt(0)
	for _, item := range miners {
		totalAmount= new(big.Int).Add(totalAmount,item.stake)
		if maxAmount.Cmp(item.stake)==-1{
			maxAmount=new(big.Int).Set(item.stake)
		}
		if minAmount.Cmp(item.stake)==1{
			minAmount=new(big.Int).Set(item.stake)
		}
	}
	avgAmount:=new(big.Int).Div(new(big.Int).Sub(new(big.Int).Sub(totalAmount,maxAmount),minAmount), big.NewInt(int64(len(miners)-2)))
	avgAmount=new(big.Int).Div(new(big.Int).Mul(avgAmount,posCandidateAvgRate),big.NewInt(100))//acgAmount * 0.75
	var tallySlice TallySlice
	for _, item := range miners {
		if item.stake.Cmp(avgAmount) >= 0 {
			tallySlice = append(tallySlice, item)
		}
	}
	sort.Sort(tallySlice)
	return tallySlice
}
func (s *Snapshot) selectMainMiner(mainMinerNumber int, mainSignerSliceLen int, signerSlice SignerSlice, mainMinerSlice TallySlice, secondMinerNumber int) SignerSlice {
	mainMinerSlice=s.reBuildMainMiner(mainMinerSlice)
	if mainMinerNumber > mainSignerSliceLen {
		//mainSignerSliceLen := len(mainMinerSlice)
		for i := 0; i < mainMinerNumber; i++ {
			signerSlice = append(signerSlice, SignerItem{mainMinerSlice[i%mainSignerSliceLen].addr, s.HistoryHash[len(s.HistoryHash)-1-i-secondMinerNumber]})
		}
	} else {
		for i := 0; i < mainMinerNumber; i++ {
			signerSlice = append(signerSlice, SignerItem{mainMinerSlice[i].addr, s.HistoryHash[len(s.HistoryHash)-1-i-secondMinerNumber]})
		}
	}
	return signerSlice
}
func  (s *Snapshot)  reBuildMainMiner(miners TallySlice) TallySlice{
	totalAmount :=big.NewInt(0)
	for _, item := range miners {
		totalAmount= new(big.Int).Add(totalAmount,item.stake)
	}
	var tallySlice TallySlice
	for _, item := range miners {
		signerNumber :=uint64(0)
		if  count,ok := s.TallySigner[item.addr];ok {
			signerNumber =count
		}
		selectParam:=s.calculateMinerState(item,totalAmount,signerNumber)
		tallySlice = append(tallySlice, TallyItem{item.addr,selectParam.BigInt() })

	}
	sort.Sort(tallySlice)
	return tallySlice
}
func  (s *Snapshot)  reBuildMiner(miners TallySlice) TallySlice{
	totalAmount :=big.NewInt(0)
	for _, item := range miners {
		if _, ok := s.TallyMiner[item.addr]; ok {
			totalAmount= new(big.Int).Add(totalAmount,item.stake)
		}
	}
	var tallySlice TallySlice
	for _, item := range miners {
		signerNumber :=uint64(0)
		if status, ok := s.TallyMiner[item.addr]; ok {
			signerNumber = status.SignerNumber
		}
		selectParam:=s.calculateMinerState(item,totalAmount,signerNumber)
		tallySlice = append(tallySlice, TallyItem{item.addr,selectParam.BigInt() })
	}
	sort.Sort(tallySlice)
	return tallySlice
}
func (s *Snapshot) calculateMinerState(item TallyItem,totalAmount *big.Int,signerNumber uint64) decimal.Decimal{

	sigerIndex :=  float64(1) /float64(signerNumber+1)
	assetsRate :=decimal.Zero
	if totalAmount.Cmp(big.NewInt(0)) > 0 {
		assetsRate =decimal.NewFromBigInt(item.stake,0).Div(decimal.NewFromBigInt(totalAmount,0))
	}
	creditWeight := uint64(defaultFullCredit)
	if _, ok := s.Punished[item.addr]; ok {
		creditWeight =defaultFullCredit - s.Punished[item.addr]
	}
	assetsRate=assetsRate.Mul(decimal.NewFromFloat(float64(creditWeight)))
	return decimal.NewFromFloat(math.Sqrt(sigerIndex )) .Mul(assetsRate).Mul(decimal.NewFromFloat(1e+18))
}