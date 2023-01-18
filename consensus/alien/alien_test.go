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

package alien

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/token/core/rawdb"
	"github.com/token/core/types"
	"github.com/token/params"
	"math/big"
	"testing"

	"github.com/token/common"
)
func TestAlien_Reward(t *testing.T) {
	arpValue :=getInspireArp(800,60)
	pledgeAmount,_:=decimal.NewFromString("600000000000000000000")
	inspireReward  :=new(big.Int).Div(new(big.Int).Mul(arpValue,pledgeAmount.BigInt()),big.NewInt(int64(yearPeridDay*100)))
	fmt.Println(inspireReward)
}
func TestAlien_PenaltyTrantor(t *testing.T) {

}
func    TestAlien_FwReward(t *testing.T) {
	var ebval =float64(1099511627776)
	tests := []struct {
		flowTotal decimal.Decimal
		expectValue decimal.Decimal
	}{
		{
			flowTotal:decimal.NewFromFloat(10995116),
			expectValue:decimal.NewFromFloat(628.607),
		},{
			flowTotal:decimal.NewFromFloat(ebval+1024),
			expectValue:decimal.NewFromFloat(634.077),
		},{
			flowTotal:decimal.NewFromFloat(ebval*2+1024),
			expectValue:decimal.NewFromFloat(639.595),
		},{
			flowTotal:decimal.NewFromFloat(ebval*3+1024),
			expectValue:decimal.NewFromFloat(645.161),
		},{
			flowTotal:decimal.NewFromFloat(ebval*4+1024),
			expectValue:decimal.NewFromFloat(650.775),
		},
	}
	for _, tt := range tests {
		fwreward := getFlowRewardScale(tt.flowTotal)
		rewardgb:= decimal.NewFromFloat(1).Div(fwreward.Mul(decimal.NewFromFloat(1024)).Div(decimal.NewFromFloat(1e+18))).Round(3)
		totalEb :=tt.flowTotal.Div(decimal.NewFromInt(1099511627776))
		var nebCount=totalEb.Round(0)
		if totalEb.Cmp(nebCount)>0 {
			nebCount= nebCount.Add(decimal.NewFromInt(1))
		}
		if rewardgb.Cmp(tt.expectValue)==0 {
			fmt.Println("Flow mining reward test pass ，",nebCount,"th EB，1 TOKEN=",rewardgb,"GB flow")
		}else{
			t.Errorf("test: Flow mining reward test failed,theory 1 TOKEN need %d GB act need %d GB",tt.expectValue.BigFloat(),rewardgb.BigFloat())
		}

	}

}
func  TestAlien_accumulatePofRewards(t *testing.T) {
	flowTotal:=[29]int64{1 ,2 ,3 ,4 ,5 ,6 ,7 ,8 ,9 ,10 ,11 ,12 ,13 ,14 ,15 ,16 ,17 ,18 ,19 ,20 ,21 ,22 ,23 ,24 ,25 ,26 ,27 ,28 ,29}
    dayRewards:=[29]float64{628.607 ,634.077 ,639.595 ,645.161 ,650.775 ,656.438 ,662.150 ,667.912 ,673.724 ,679.587 ,685.501 ,691.466 ,697.483 ,703.553 ,709.675 ,715.851 ,722.080 ,728.363 ,734.702 ,741.095 ,747.544 ,754.049 ,760.611 ,767.230 ,773.906 ,780.641 ,787.434 ,794.286 ,801.198 }
	for i,item:=range flowTotal{
		dayReward:=dayRewards[i]*1024
		testAlien_accumulatePofRewards(t,dayReward,new(big.Int).Mul(big.NewInt(1099511627780),big.NewInt(item)))
	}
}
func  testAlien_accumulatePofRewards(t *testing.T,dayReward1 float64,flowTotal *big.Int) {
	var period =uint64(10)
	FlowValue1:=uint64(160000)
	FlowValue2:=uint64(480000)
	bandwidth1:=uint64(100)
	bandwidth2:=uint64(400)
	tests := []struct {
		maxSignerCount uint64
		number         uint64
		coinbase    common.Address
		Period      uint64
		expectValue  decimal.Decimal
		FlowValue1 uint64
		bandwidth uint64
		PofPrice *big.Int
	}{
		{
			maxSignerCount :3,
			number     :1,
			coinbase : common.HexToAddress("0x7a4539ed8a0b8b4583ead1e5a3f604e83419f502"),
			Period :period,
			FlowValue1:FlowValue1,
			expectValue: decimal.NewFromFloat(float64(FlowValue1)/float64(dayReward1)).Round(10),
			bandwidth:bandwidth1,
			PofPrice:new(big.Int).Set(pofDefaultBasePrice),
		},{
			maxSignerCount :3,
			number     :1,
			coinbase : common.HexToAddress("0x2aD0559Afade09a22364F6380f52BF57E9057B8D"),
			Period :period,
			bandwidth:bandwidth2,
			FlowValue1:FlowValue2,
			expectValue: decimal.NewFromFloat(float64(FlowValue2)/float64(dayReward1)).Round(10),
			PofPrice:new(big.Int).Set(pofDefaultBasePrice),
		},
	}

	snap := &Snapshot{
		LCRS:      1,
		Tally:     make(map[common.Address]*big.Int),
		Punished:  make(map[common.Address]uint64),
		PofMiner:&PofMinerSnap{
			PofMinerPrev: make(map[common.Address]map[common.Hash]*PofMinerReport),
		},
		FlowTotal:new(big.Int).Set(flowTotal),
		PofPledge:make(map[common.Address]*PofPledgeItem),
		SystemConfig:SystemParameter{
			Deposit:make(map[uint32]*big.Int),
		},
	}
	var LockRewardRecord[] LockRewardRecord
	for _, tt := range tests {
		snap.config = &params.AlienConfig{
			MaxSignerCount: tt.maxSignerCount,
			Period:         tt.Period,
		}
		snap.Number = tt.number
		snap.PofMiner.PofMinerPrev[tt.coinbase]=make(map[common.Hash]*PofMinerReport)
		snap.PofMiner.PofMinerPrev[tt.coinbase][common.Hash{}]=&PofMinerReport{
			FlowValue1:tt.FlowValue1,
		}
		snap.PofPledge[tt.coinbase]=&PofPledgeItem{
			PledgeAmount:new(big.Int).SetUint64(tt.FlowValue1),
			PledgeStatus: pofstatus_normal,
			Bandwidth:tt.bandwidth,
			PofPrice:tt.PofPrice,
		}
	}
	snap.SystemConfig.Deposit[sscEnumPofWithinBasePricePeriod]=new(big.Int).Set(pofDefaultBasePrice)
	db:=rawdb.NewMemoryDatabase()
	currentHeaderExtra := HeaderExtra{}
	LockRewardRecord, _ = accumulatePofRewards(currentHeaderExtra.LockReward,  snap, db)
	for _, tt := range tests {
		for index := range LockRewardRecord {
			if tt.coinbase==LockRewardRecord[index].Target{
				actReward:=decimal.NewFromBigInt(LockRewardRecord[index].Amount,0).Div(decimal.NewFromInt(1E+18)).Round(10)
				cut:=actReward.Sub(tt.expectValue).Abs()
				if cut.Cmp(decimal.NewFromFloat(float64(0.02)))<=0{
					t.Logf("coinbase %s,Bandwidth reward calculation pass,expect %d nbn,but act %d,cut is %d" ,(LockRewardRecord[index].Target).Hex(),tt.expectValue.BigFloat(),actReward.BigFloat(),cut.BigFloat())
				}else {
					t.Errorf("coinbase %s,Bandwidth reward calculation error,expect %d nbn,but act %d,cut is %d" ,(LockRewardRecord[index].Target).Hex(),tt.expectValue.BigFloat(),actReward.BigFloat(),cut.BigFloat())
				}
			}
		}
	}
}


func  TestAlien_accumulateInspireRewards(t *testing.T) {
	var period =uint64(10)
	dayReward1:=decimal.NewFromInt(1)
	dayReward1=dayReward1.Mul(decimal.NewFromInt(2500))
	dayReward1=dayReward1.Div(decimal.NewFromInt(10000))
	dayReward1=dayReward1.Div(decimal.NewFromInt(365))

	FlowValue1:=uint64(100000)
	FlowValue2:=uint64(400000)
	bandwidth1:=uint64(100)
	bandwidth2:=uint64(400)
	tests := []struct {
		maxSignerCount uint64
		number         uint64
		coinbase    common.Address
		Period      uint64
		expectValue  decimal.Decimal
		FlowValue1 uint64
		bandwidth uint64
	}{
		{
			maxSignerCount :3,
			number     :1,
			coinbase : common.HexToAddress("0x7a4539ed8a0b8b4583ead1e5a3f604e83419f502"),
			Period :period,
			FlowValue1:FlowValue1,
			expectValue: decimal.NewFromFloat(float64(FlowValue1)).Mul(dayReward1).Round(10),
			bandwidth:bandwidth1,
		},{
			maxSignerCount :3,
			number     :1,
			coinbase : common.HexToAddress("0x2aD0559Afade09a22364F6380f52BF57E9057B8D"),
			Period :period,
			bandwidth:bandwidth2,
			FlowValue1:FlowValue2,
			expectValue: decimal.NewFromFloat(float64(FlowValue2)).Mul(dayReward1).Round(10),
		},
	}

	snap := &Snapshot{
		LCRS:      1,
		Tally:     make(map[common.Address]*big.Int),
		Punished:  make(map[common.Address]uint64),
		PofMiner:&PofMinerSnap{
			PofMinerPrev: make(map[common.Address]map[common.Hash]*PofMinerReport),
		},
		FlowTotal:big.NewInt(0),
		PofPledge:make(map[common.Address]*PofPledgeItem),
		InspireHarvest:common.Big0,
	}

	var LockRewardRecord[] LockRewardRecord
	for _, tt := range tests {
		snap.config = &params.AlienConfig{
			MaxSignerCount: tt.maxSignerCount,
			Period:         tt.Period,
		}
		snap.Number = tt.number
		snap.PofMiner.PofMinerPrev[tt.coinbase]=make(map[common.Hash]*PofMinerReport)
		snap.PofMiner.PofMinerPrev[tt.coinbase][common.Hash{}]=&PofMinerReport{
			FlowValue1:tt.FlowValue1,
		}
		snap.PofPledge[tt.coinbase]=&PofPledgeItem{
			PledgeAmount:new(big.Int).SetUint64(tt.FlowValue1),
			PledgeStatus: pofstatus_normal,
			Bandwidth: tt.bandwidth,
		}
	}
	db:=rawdb.NewMemoryDatabase()
	currentHeaderExtra := HeaderExtra{}
	header:=&types.Header{
		Number:big.NewInt(9999),
	}
	LockRewardRecord, _ = accumulateInspireRewards(currentHeaderExtra.LockReward, header, snap,db)
	for _, tt := range tests {
		for index := range LockRewardRecord {
			if tt.coinbase==LockRewardRecord[index].Target{
				actReward:=decimal.NewFromBigInt(LockRewardRecord[index].Amount,0).Round(10)
				cut:=actReward.Sub(tt.expectValue).Abs()
				if cut.Cmp(decimal.NewFromFloat(float64(0.99)))<=0{
					t.Logf("coinbase %s,Bandwidth reward calculation pass,expect %d nbn,but act %d,cut is %d" ,(LockRewardRecord[index].Target).Hex(),tt.expectValue.BigFloat(),actReward.BigFloat(),cut.BigFloat())
				}else {
					t.Errorf("coinbase %s,Bandwidth reward calculation error,expect %d nbn,but act %d,cut is %d" ,(LockRewardRecord[index].Target).Hex(),tt.expectValue.BigFloat(),actReward.BigFloat(),cut.BigFloat())
				}
			}
		}
	}
}

func  TestAlien_getPofPlegeAmount(t *testing.T) {
	var period =uint64(10)
	inspireHarvest:=big.NewInt(9000000)
	pofHarvest:=big.NewInt(10000000000)
	pofHarvest=new(big.Int).Mul(big.NewInt(1e+18),pofHarvest)
	bandwidth1:=uint64(100)
	bandwidth2:=uint64(400)
	tests := []struct {
		number         uint64
		expectValue  decimal.Decimal
		bandwidth uint64
		coinbase    common.Address
	}{
		{
			number     :1,
			expectValue: decimal.NewFromFloat(float64(bandwidth1)).Mul(decimal.NewFromBigInt(pofBasePledgeAmount,0)).Round(10),
			bandwidth:bandwidth1,
			coinbase : common.HexToAddress("0x7a4539ed8a0b8b4583ead1e5a3f604e83419f502"),
		},{
			number     :1,
			expectValue: decimal.NewFromFloat(float64(bandwidth2)).Mul(decimal.NewFromBigInt(pofBasePledgeAmount,0)).Round(10),
			bandwidth:bandwidth2,
		},
		{
			number     :365*24*60*60/period+100,
			//(inspireHarvest+pofHarvest)/bandwidth1+bandwidth2+bandwidth2
			//expectValue: decimal.NewFromFloat(float64(bandwidth2)).Mul(decimal.NewFromFloat(float64(11121111))).Round(10),
			expectValue:decimal.NewFromFloat(float64(bandwidth2)).Mul(decimal.NewFromBigInt(pofBasePledgeAmount,0)).Round(10),
			bandwidth:bandwidth2,
			coinbase : common.HexToAddress("0x2aD0559Afade09a22364F6380f52BF57E9057B8D"),
		},
		{
			number     :365*24*60*60/period+1900,
			//(inspireHarvest+pofHarvest)/bandwidth1+bandwidth2+bandwidth2
			expectValue:decimal.NewFromFloat(float64(bandwidth2)).Mul(decimal.NewFromBigInt(pofBasePledgeAmount,0)).Round(10),
			bandwidth:bandwidth2,
			coinbase : common.HexToAddress("0x2aD0559Afade09a22364F6380f52BF57E9057B8D"),
		},
	}
	snap := &Snapshot{
		LCRS:      1,
		Tally:     make(map[common.Address]*big.Int),
		Punished:  make(map[common.Address]uint64),
		PofMiner:&PofMinerSnap{
			PofMinerPrev: make(map[common.Address]map[common.Hash]*PofMinerReport),
		},
		FlowTotal:big.NewInt(0),
		PofPledge:make(map[common.Address]*PofPledgeItem),
		InspireHarvest:new(big.Int).Set(inspireHarvest),
		PofHarvest: new(big.Int).Set(pofHarvest),
	}
	snap.config = &params.AlienConfig{
		Period:period,
	}
	for _, tt := range tests {
		snap.PofPledge[tt.coinbase]=&PofPledgeItem{
			PledgeStatus: pofstatus_normal,
			Bandwidth:tt.bandwidth,
		}
	}

	for _, tt := range tests {
		snap.Number=tt.number
		amountBig:= snap.getPofPlegeAmount(tt.bandwidth)
		actReward:=decimal.NewFromBigInt(amountBig,0)
		cut:=actReward.Sub(tt.expectValue).Abs()
		if cut.Cmp(decimal.NewFromFloat(float64(0.99)))<=0{
			t.Logf("Bandwidth %d reward calculation pass,expect %d nbn,but act %d,cut is %d" ,tt.bandwidth,tt.expectValue.BigFloat(),actReward.BigFloat(),cut.BigFloat())
		}else {
			t.Errorf("Bandwidth %d reward calculation error,expect %d nbn,but act %d,cut is %d" ,tt.bandwidth,tt.expectValue.BigFloat(),actReward.BigFloat(),cut.BigFloat())
		}
	}
}