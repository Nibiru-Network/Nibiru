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
	"github.com/token/common"
	"github.com/token/consensus"
	"math/big"
	"testing"
)

const (
	addr1="0x7a4539ed8a0b8b4583ead1e5a3f604e83419f502"
	addr2="0x8C4022E7A0723A7AabA3a3b9a8425548117c06F4"
	addr3="0x0Ff6e773Ff893fF39ed9352160889df13BDfc896"
	addr4="0xa63b29EBe0A141B87A87e39dE17F17346e11e1b7"
)
var amount=big.NewInt(60)

func TestAlien_verifyHeaderExtern_verifyExchangeCoin(t *testing.T) {

	deviceone:= ExchangeCoinRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
	}
	devicetwo:= ExchangeCoinRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
	}
	currentdevices := []ExchangeCoinRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []ExchangeCoinRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		ExchangeCoin: currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		ExchangeCoin: verifydevices,
	}
	devicethree:= ExchangeCoinRecord{
		Target:common.HexToAddress(addr1),
		Amount:new (big.Int).Add(amount,big.NewInt(100)),
	}
	verifydevices2 := []ExchangeCoinRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		ExchangeCoin: verifydevices2,
	}
	verifydevices3 := []ExchangeCoinRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		ExchangeCoin: verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"verifyExchangeCoin")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"verifyExchangeCoin")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"verifyExchangeCoin")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"verifyExchangeCoin")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"verifyExchangeCoin")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"verifyExchangeCoin")

	//ABB,AAB

	verifydevices5 := []ExchangeCoinRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		ExchangeCoin: verifydevices5,
	}

	verifydevices6 := []ExchangeCoinRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		ExchangeCoin: verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyExchangeCoin")

	devicefour:= ExchangeCoinRecord{
		Target:common.HexToAddress(addr2),
		Amount:new (big.Int).Add(amount,big.NewInt(100)),
	}

	verifydevices7 := []ExchangeCoinRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		ExchangeCoin: verifydevices7,
	}

	verifydevices8 := []ExchangeCoinRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		ExchangeCoin: verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyExchangeCoin")

}
func test_verifyHeaderExtern(t *testing.T,currentExtra *HeaderExtra, verifyExtra *HeaderExtra,name string) {
	err:=verifyHeaderExtern(currentExtra,verifyExtra)
	if err==nil{
		t.Logf(name+" pass")
	}else{
		t.Errorf(name+" error,expect nil,but act %s" ,err.Error())
	}
}

func test_verifyHeaderExtern2(t *testing.T,currentExtra *HeaderExtra, verifyExtra *HeaderExtra,name string) {
	err:=verifyHeaderExtern(currentExtra,verifyExtra)
	if err!=nil{
		t.Logf(name+" pass,error msg %s",err.Error())
	}else{
		t.Errorf(" error,expect not nil,err is nil")
	}
}


func TestAlien_verifyHeaderExtern_verifyDeviceBind(t *testing.T) {

	deviceone:=DeviceBindRecord{
	        Device:common.HexToAddress(addr1),
		    Revenue:common.HexToAddress(addr1),
			Contract:common.HexToAddress(addr1),
			MultiSign:common.HexToAddress(addr1),
			Type:0,
			Bind:false,
	}
	devicetwo:=DeviceBindRecord{
		Device:common.HexToAddress(addr2),
		Revenue:common.HexToAddress(addr2),
		Contract:common.HexToAddress(addr2),
		MultiSign:common.HexToAddress(addr2),
		Type:0,
		Bind:false,
	}
	currentdevices := []DeviceBindRecord{
		deviceone,
		devicetwo,
	}

	verifydevices := []DeviceBindRecord{
		deviceone,
		devicetwo,
	}

	currentExtra:=&HeaderExtra{
		DeviceBind:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		DeviceBind:verifydevices,
	}

	devicethree:=DeviceBindRecord{
		Device:common.HexToAddress(addr3),
		Revenue:common.HexToAddress(addr3),
		Contract:common.HexToAddress(addr3),
		MultiSign:common.HexToAddress(addr3),
		Type:0,
		Bind:false,
	}
	verifydevices2 := []DeviceBindRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		DeviceBind:verifydevices2,
	}

	verifydevices3 := []DeviceBindRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		DeviceBind:verifydevices3,
	}

	verifyExtra4:=&HeaderExtra{
	}

	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"verifyDeviceBind")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"verifyDeviceBind")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"verifyDeviceBind")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"verifyDeviceBind")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"verifyDeviceBind")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"verifyDeviceBind")

	//ABB,AAB

	verifydevices5 := []DeviceBindRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		DeviceBind:verifydevices5,
	}

	verifydevices6 := []DeviceBindRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		DeviceBind:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyDeviceBind")

	devicefour:=DeviceBindRecord{
		Device:common.HexToAddress(addr3),
		Revenue:common.HexToAddress(addr3),
		Contract:common.HexToAddress(addr3),
		MultiSign:common.HexToAddress(addr3),
		Type:0,
		Bind:true,
	}

	verifydevices7 := []DeviceBindRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		DeviceBind:verifydevices7,
	}

	verifydevices8 := []DeviceBindRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		DeviceBind:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyDeviceBind")

}


func TestAlien_verifyHeaderExtern_verifyLockReward(t *testing.T) {
	deviceone:=LockRewardRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		IsReward:3,
		FlowValue1:0,
		FlowValue2:0,
	}

	devicetwo:=LockRewardRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		IsReward:3,
		FlowValue1:0,
		FlowValue2:0,
	}

	devicethree:=LockRewardRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		IsReward:4,
		FlowValue1:0,
		FlowValue2:0,
	}

	currentdevices := []LockRewardRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []LockRewardRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		LockReward:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		LockReward:verifydevices,
	}

	verifydevices2 := []LockRewardRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		LockReward:verifydevices2,
	}
	verifydevices3 := []LockRewardRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		LockReward:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"verifyLockReward")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"verifyLockReward")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"verifyLockReward")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"verifyLockReward")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"verifyLockReward")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"verifyLockReward")

	verifydevices5 := []LockRewardRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		LockReward:verifydevices5,
	}

	verifydevices6 := []LockRewardRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		LockReward:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyLockReward")

	devicefour:=LockRewardRecord{
		Target:common.HexToAddress(addr3),
		Amount:new(big.Int).Add(amount,big.NewInt(1000)),
		IsReward:8,
		FlowValue1:666888,
		FlowValue2:77888,
	}

	verifydevices7 := []LockRewardRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		LockReward:verifydevices7,
	}

	verifydevices8 := []LockRewardRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		LockReward:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyLockReward")

}

func TestAlien_verifyHeaderExtern_verifyCandidatePunish(t *testing.T) {
	deviceone:=CandidatePunishRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		Credit: 100,
	}

	devicetwo:=CandidatePunishRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		Credit: 100,
	}

	devicethree:=CandidatePunishRecord{
		Target:common.HexToAddress(addr1),
		Amount:new (big.Int).Add(amount,big.NewInt(100)),
		Credit: 200,
	}

	currentdevices := []CandidatePunishRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []CandidatePunishRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		CandidatePunish:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		CandidatePunish:verifydevices,
	}

	verifydevices2 := []CandidatePunishRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		CandidatePunish:verifydevices2,
	}
	verifydevices3 := []CandidatePunishRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		CandidatePunish:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"verifyCandidatePunish")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"verifyCandidatePunish")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"verifyCandidatePunish")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"verifyCandidatePunish")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"verifyCandidatePunish")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"verifyCandidatePunish")

	verifydevices5 := []CandidatePunishRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		CandidatePunish:verifydevices5,
	}

	verifydevices6 := []CandidatePunishRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		CandidatePunish:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyCandidatePunish")

	devicefour:=CandidatePunishRecord{
		Target:common.HexToAddress(addr3),
		Amount:new(big.Int).Add(amount,big.NewInt(1000)),
		Credit:1121212,
	}

	verifydevices7 := []CandidatePunishRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		CandidatePunish:verifydevices7,
	}

	verifydevices8 := []CandidatePunishRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		CandidatePunish:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyCandidatePunish")
}

func TestAlien_verifyHeaderExtern_verifyMinerStake(t *testing.T) {
	deviceone:=MinerStakeRecord{
		Target:common.HexToAddress(addr1),
		Stake:amount,
	}
	devicetwo:=MinerStakeRecord{
		Target:common.HexToAddress(addr1),
		Stake:amount,
	}

	devicethree:=MinerStakeRecord{
		Target:common.HexToAddress(addr1),
		Stake:new (big.Int).Add(amount,big.NewInt(100)),
	}

	currentdevices := []MinerStakeRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []MinerStakeRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		MinerStake:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		MinerStake:verifydevices,
	}

	verifydevices2 := []MinerStakeRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		MinerStake:verifydevices2,
	}
	verifydevices3 := []MinerStakeRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		MinerStake:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"verifyMinerStake")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"verifyMinerStake")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"verifyMinerStake")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"verifyMinerStake")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"verifyMinerStake")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"verifyMinerStake")

	verifydevices5 := []MinerStakeRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		MinerStake:verifydevices5,
	}

	verifydevices6 := []MinerStakeRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		MinerStake:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyMinerStake")

	devicefour:=MinerStakeRecord{
		Target:common.HexToAddress(addr3),
		Stake:new(big.Int).Add(amount,big.NewInt(9000)),
	}

	verifydevices7 := []MinerStakeRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		MinerStake:verifydevices7,
	}

	verifydevices8 := []MinerStakeRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		MinerStake:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyMinerStake")
}

func TestAlien_verifyHeaderExtern_CandidateExit(t *testing.T) {
	currentdevices := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr2),
	}
	verifydevices := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr2),
	}
	currentExtra:=&HeaderExtra{
		CandidateExit:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		CandidateExit:verifydevices,
	}

	verifydevices2 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr3),
	}
	verifyExtra2:=&HeaderExtra{
		CandidateExit:verifydevices2,
	}
	verifydevices3 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr2),
		common.HexToAddress(addr3),
	}
	verifyExtra3:=&HeaderExtra{
		CandidateExit:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"CandidateExit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"CandidateExit")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"CandidateExit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"CandidateExit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"CandidateExit")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"CandidateExit")

	verifydevices5 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr3),
		common.HexToAddress(addr3),
	}
	verifyExtra5:=&HeaderExtra{
		CandidateExit:verifydevices5,
	}

	verifydevices6 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr1),
		common.HexToAddress(addr3),
	}
	verifyExtra6:=&HeaderExtra{
		CandidateExit:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"CandidateExit")

	verifydevices7 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr1),
		common.HexToAddress(addr4),
	}
	verifyExtra7:=&HeaderExtra{
		CandidateExit:verifydevices7,
	}

	verifydevices8 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr4),
		common.HexToAddress(addr4),
	}
	verifyExtra8:=&HeaderExtra{
		CandidateExit:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"CandidateExit")
}

func TestAlien_verifyHeaderExtern_ClaimedBandwidth(t *testing.T) {
	deviceone:=ClaimedBandwidthRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		Bandwidth:100,
	}
	devicetwo:=ClaimedBandwidthRecord{
		Target:common.HexToAddress(addr1),
		Amount:amount,
		Bandwidth:100,
	}

	devicethree:=ClaimedBandwidthRecord{
		Target:common.HexToAddress(addr1),
		Amount:new (big.Int).Add(amount,big.NewInt(100)),
		Bandwidth:100,
	}

	currentdevices := []ClaimedBandwidthRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []ClaimedBandwidthRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		ClaimedBandwidth:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		ClaimedBandwidth:verifydevices,
	}

	verifydevices2 := []ClaimedBandwidthRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		ClaimedBandwidth:verifydevices2,
	}
	verifydevices3 := []ClaimedBandwidthRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		ClaimedBandwidth:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"verifyClaimedBandwidth")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"verifyClaimedBandwidth")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"verifyClaimedBandwidth")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"verifyClaimedBandwidth")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"verifyClaimedBandwidth")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"verifyClaimedBandwidth")

	verifydevices5 := []ClaimedBandwidthRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		ClaimedBandwidth:verifydevices5,
	}

	verifydevices6 := []ClaimedBandwidthRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		ClaimedBandwidth:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyClaimedBandwidth")

	devicefour:=ClaimedBandwidthRecord{
		Target:common.HexToAddress(addr3),
		Amount:new(big.Int).Add(amount,big.NewInt(1000)),
		Bandwidth:666888999,
	}

	verifydevices7 := []ClaimedBandwidthRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		ClaimedBandwidth:verifydevices7,
	}

	verifydevices8 := []ClaimedBandwidthRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		ClaimedBandwidth:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyClaimedBandwidth")
}


func TestAlien_verifyHeaderExtern_PofMinerExit(t *testing.T) {
	currentdevices := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr2),
	}
	verifydevices := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr2),
	}
	currentExtra:=&HeaderExtra{
		PofMinerExit: currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		PofMinerExit: verifydevices,
	}

	verifydevices2 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr3),
	}
	verifyExtra2:=&HeaderExtra{
		PofMinerExit: verifydevices2,
	}
	verifydevices3 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr2),
		common.HexToAddress(addr3),
	}
	verifyExtra3:=&HeaderExtra{
		PofMinerExit: verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"PofMinerExit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"PofMinerExit")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"PofMinerExit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"PofMinerExit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"PofMinerExit")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"PofMinerExit")

	verifydevices5 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr3),
		common.HexToAddress(addr3),
	}
	verifyExtra5:=&HeaderExtra{
		PofMinerExit: verifydevices5,
	}

	verifydevices6 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr1),
		common.HexToAddress(addr3),
	}
	verifyExtra6:=&HeaderExtra{
		PofMinerExit: verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"PofMinerExit")

	verifydevices7 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr1),
		common.HexToAddress(addr4),
	}
	verifyExtra7:=&HeaderExtra{
		PofMinerExit: verifydevices7,
	}

	verifydevices8 := []common.Address{
		common.HexToAddress(addr1),
		common.HexToAddress(addr4),
		common.HexToAddress(addr4),
	}
	verifyExtra8:=&HeaderExtra{
		PofMinerExit: verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"PofMinerExit")

}


func TestAlien_verifyHeaderExtern_BandwidthPunish(t *testing.T) {

}

func TestAlien_verifyHeaderExtern_ConfigExchRate(t *testing.T) {
	currentExtra:=&HeaderExtra{
		ConfigExchRate:30,
	}
	verifyExtra:=&HeaderExtra{
		ConfigExchRate:30,
	}
	verifyExtra2:=&HeaderExtra{
		ConfigExchRate:31,
	}
	verifyExtra3:=&HeaderExtra{
		ConfigExchRate:32,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"ConfigExchRate")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"ConfigExchRate")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"ConfigExchRate")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"ConfigExchRate")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"ConfigExchRate")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"ConfigExchRate")
}

func TestAlien_verifyHeaderExtern_ConfigOffLine(t *testing.T) {
	currentExtra:=&HeaderExtra{
		ConfigOffLine:30,
	}
	verifyExtra:=&HeaderExtra{
		ConfigOffLine:30,
	}
	verifyExtra2:=&HeaderExtra{
		ConfigOffLine:31,
	}
	verifyExtra3:=&HeaderExtra{
		ConfigOffLine:32,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"ConfigOffLine")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"ConfigOffLine")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"ConfigOffLine")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"ConfigOffLine")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"ConfigOffLine")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"ConfigOffLine")
}


func TestAlien_verifyHeaderExtern_ConfigDeposit(t *testing.T) {
	deviceone:=ConfigDepositRecord{
		Who:0,
		Amount:amount,
	}
	devicetwo:=ConfigDepositRecord{
		Who:0,
		Amount:amount,
	}

	devicethree:=ConfigDepositRecord{
		Who:1,
		Amount:amount,
	}

	currentdevices := []ConfigDepositRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []ConfigDepositRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		ConfigDeposit:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		ConfigDeposit:verifydevices,
	}

	verifydevices2 := []ConfigDepositRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		ConfigDeposit:verifydevices2,
	}
	verifydevices3 := []ConfigDepositRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		ConfigDeposit:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"ConfigDeposit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"ConfigDeposit")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"ConfigDeposit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"ConfigDeposit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"ConfigDeposit")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"ConfigDeposit")

	verifydevices5 := []ConfigDepositRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		ConfigDeposit:verifydevices5,
	}

	verifydevices6 := []ConfigDepositRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		ConfigDeposit:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"verifyLockReward")

	devicefour:=ConfigDepositRecord{
		Who:7,
		Amount:new(big.Int).Add(amount,big.NewInt(12200)),
	}

	verifydevices7 := []ConfigDepositRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		ConfigDeposit:verifydevices7,
	}

	verifydevices8 := []ConfigDepositRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		ConfigDeposit:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"verifyLockReward")
}


func TestAlien_verifyHeaderExtern_ConfigISPQOS(t *testing.T) {
	deviceone:=ISPQOSRecord{
		ISPID:0,
		QOS:2,
	}
	devicetwo:=ISPQOSRecord{
		ISPID:0,
		QOS:2,
	}

	devicethree:=ISPQOSRecord{
		ISPID:6,
		QOS:8,
	}

	currentdevices := []ISPQOSRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []ISPQOSRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		ConfigISPQOS:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		ConfigISPQOS:verifydevices,
	}

	verifydevices2 := []ISPQOSRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		ConfigISPQOS:verifydevices2,
	}
	verifydevices3 := []ISPQOSRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		ConfigISPQOS:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"ConfigISPQOS")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"ConfigISPQOS")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"ConfigISPQOS")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"ConfigISPQOS")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"ConfigISPQOS")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"ConfigISPQOS")

	verifydevices5 := []ISPQOSRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		ConfigISPQOS:verifydevices5,
	}

	verifydevices6 := []ISPQOSRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		ConfigISPQOS:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"ConfigISPQOS")

	devicefour:=ISPQOSRecord{
		ISPID:8,
		QOS:666888,
	}

	verifydevices7 := []ISPQOSRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		ConfigISPQOS:verifydevices7,
	}

	verifydevices8 := []ISPQOSRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		ConfigISPQOS:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"ConfigISPQOS")
}

func TestAlien_verifyHeaderExtern_LockParameters(t *testing.T) {
	deviceone:=LockParameterRecord{
		LockPeriod:0,
		RlsPeriod:2,
		Interval:3,
		Who:6,
	}
	devicetwo:=LockParameterRecord{
		LockPeriod:0,
		RlsPeriod:2,
		Interval:3,
		Who:6,
	}

	devicethree:=LockParameterRecord{
		LockPeriod:0,
		RlsPeriod:2,
		Interval:3,
		Who:8,
	}

	currentdevices := []LockParameterRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []LockParameterRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		LockParameters:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		LockParameters:verifydevices,
	}

	verifydevices2 := []LockParameterRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		LockParameters:verifydevices2,
	}
	verifydevices3 := []LockParameterRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		LockParameters:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"LockParameters")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"LockParameters")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"LockParameters")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"LockParameters")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"LockParameters")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"LockParameters")

	verifydevices5 := []LockParameterRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		LockParameters:verifydevices5,
	}

	verifydevices6 := []LockParameterRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		LockParameters:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"LockParameters")

	devicefour:=LockParameterRecord{
		LockPeriod:8,
		RlsPeriod:666888,
		Interval:77888,
		Who:77888999,
	}

	verifydevices7 := []LockParameterRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		LockParameters:verifydevices7,
	}

	verifydevices8 := []LockParameterRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		LockParameters:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"LockParameters")
}

func TestAlien_verifyHeaderExtern_ManagerAddress(t *testing.T) {
	deviceone:=ManagerAddressRecord{
		Target:common.HexToAddress(addr1),
		Who:6,
	}

	devicetwo:=ManagerAddressRecord{
		Target:common.HexToAddress(addr1),
		Who:6,
	}

	devicethree:=ManagerAddressRecord{
		Target:common.HexToAddress(addr2),
		Who:8,
	}

	currentdevices := []ManagerAddressRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []ManagerAddressRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		ManagerAddress:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		ManagerAddress:verifydevices,
	}

	verifydevices2 := []ManagerAddressRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		ManagerAddress:verifydevices2,
	}
	verifydevices3 := []ManagerAddressRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		ManagerAddress:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"ManagerAddress")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"ManagerAddress")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"ManagerAddress")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"ManagerAddress")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"ManagerAddress")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"ManagerAddress")

	verifydevices5 := []ManagerAddressRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		ManagerAddress:verifydevices5,
	}

	verifydevices6 := []ManagerAddressRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		ManagerAddress:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"ManagerAddress")

	devicefour:=ManagerAddressRecord{
		Target:common.HexToAddress(addr3),
		Who:811111,
	}

	verifydevices7 := []ManagerAddressRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		ManagerAddress:verifydevices7,
	}

	verifydevices8 := []ManagerAddressRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		ManagerAddress:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"ManagerAddress")
}


func TestAlien_verifyHeaderExtern_FlowHarvest(t *testing.T) {
	currentExtra:=&HeaderExtra{
		PofHarvest:amount,
	}
	verifyExtra:=&HeaderExtra{
		PofHarvest:amount,
	}
	verifyExtra2:=&HeaderExtra{
		PofHarvest:new(big.Int).Add(amount,amount),
	}
	verifyExtra3:=&HeaderExtra{
		PofHarvest:new(big.Int).Add(amount,big.NewInt(600)),
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"FlowHarvest")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"FlowHarvest")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"FlowHarvest")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"FlowHarvest")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"FlowHarvest")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"FlowHarvest")
}

func TestAlien_verifyHeaderExtern_GrantProfit(t *testing.T) {
	deviceone:=consensus.GrantProfitRecord{
		Which:0,
		MinerAddress:common.HexToAddress(addr1),
		BlockNumber:0,
		Amount:amount,
		RevenueAddress: common.HexToAddress(addr2),
		RevenueContract: common.HexToAddress(addr2),
		MultiSignature: common.HexToAddress(addr3),
	}

	devicetwo:=consensus.GrantProfitRecord{
		Which:0,
		MinerAddress:common.HexToAddress(addr1),
		BlockNumber:0,
		Amount:amount,
		RevenueAddress: common.HexToAddress(addr2),
		RevenueContract: common.HexToAddress(addr2),
		MultiSignature: common.HexToAddress(addr3),
	}

	devicethree:=consensus.GrantProfitRecord{
		Which:0,
		MinerAddress:common.HexToAddress(addr1),
		BlockNumber:0,
		Amount:amount,
		RevenueAddress: common.HexToAddress(addr2),
		RevenueContract: common.HexToAddress(addr2),
		MultiSignature: common.HexToAddress(addr2),
	}

	currentdevices := []consensus.GrantProfitRecord{
		deviceone,
		devicetwo,
	}
	verifydevices := []consensus.GrantProfitRecord{
		deviceone,
		devicetwo,
	}
	currentExtra:=&HeaderExtra{
		GrantProfit:currentdevices,
	}
	verifyExtra:=&HeaderExtra{
		GrantProfit:verifydevices,
	}

	verifydevices2 := []consensus.GrantProfitRecord{
		deviceone,
		devicethree,
	}
	verifyExtra2:=&HeaderExtra{
		GrantProfit:verifydevices2,
	}
	verifydevices3 := []consensus.GrantProfitRecord{
		deviceone,
		devicetwo,
		devicethree,
	}
	verifyExtra3:=&HeaderExtra{
		GrantProfit:verifydevices3,
	}
	verifyExtra4:=&HeaderExtra{
	}
	test_verifyHeaderExtern(t,currentExtra,verifyExtra,"GrantProfit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra2,"GrantProfit")
	test_verifyHeaderExtern2(t,verifyExtra2,currentExtra,"GrantProfit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra3,"GrantProfit")
	test_verifyHeaderExtern2(t,currentExtra,verifyExtra4,"GrantProfit")
	test_verifyHeaderExtern2(t,verifyExtra4,currentExtra,"GrantProfit")

	verifydevices5 := []consensus.GrantProfitRecord{
		deviceone,
		devicethree,
		devicethree,
	}
	verifyExtra5:=&HeaderExtra{
		GrantProfit:verifydevices5,
	}

	verifydevices6 := []consensus.GrantProfitRecord{
		deviceone,
		deviceone,
		devicethree,
	}
	verifyExtra6:=&HeaderExtra{
		GrantProfit:verifydevices6,
	}

	test_verifyHeaderExtern2(t,verifyExtra5,verifyExtra6,"GrantProfit")

	devicefour:=consensus.GrantProfitRecord{
		Which:8,
		MinerAddress:common.HexToAddress(addr3),
		BlockNumber:95535,
		Amount:new(big.Int).Add(amount,big.NewInt(1000)),
		RevenueAddress:common.HexToAddress(addr4),
		RevenueContract:common.HexToAddress(addr3),
		MultiSignature:common.HexToAddress(addr3),
	}

	verifydevices7 := []consensus.GrantProfitRecord{
		deviceone,
		deviceone,
		devicefour,
	}
	verifyExtra7:=&HeaderExtra{
		GrantProfit:verifydevices7,
	}

	verifydevices8 := []consensus.GrantProfitRecord{
		deviceone,
		devicefour,
		devicefour,
	}
	verifyExtra8:=&HeaderExtra{
		GrantProfit:verifydevices8,
	}

	test_verifyHeaderExtern2(t,verifyExtra7,verifyExtra8,"GrantProfit")
}

func TestAlien_verifyHeaderExtern_PofReport(t *testing.T) {

}


func TestAlien_compareExchangeCoin(t *testing.T){
	a:=[]ExchangeCoinRecord{
		{
			Target: common.HexToAddress(addr1),
			Amount: amount,
		},
	}
	b:=[]ExchangeCoinRecord{
		{
			Target: common.HexToAddress(addr1),
			Amount: amount,
		},
	}
	c:=[]ExchangeCoinRecord{
		{
			Target: common.HexToAddress(addr2),
			Amount: amount,
		},
	}
	d:=[]ExchangeCoinRecord{
		{
			Target: common.HexToAddress(addr1),
			Amount: new(big.Int).Add(amount,big.NewInt(6000)),
		},
	}

	e:=[]ExchangeCoinRecord{
		{
			Target: common.HexToAddress(addr2),
			Amount: new(big.Int).Add(amount,big.NewInt(6000)),
		},
	}
	err:= compareExchangeCoin(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}

	test_compareExchangeCoin(t,a,c)
	test_compareExchangeCoin(t,a,d)
	test_compareExchangeCoin(t,a,e)
}

func test_compareExchangeCoin(t *testing.T,a []ExchangeCoinRecord, b []ExchangeCoinRecord){
	err:= compareExchangeCoin(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareDeviceBind(t *testing.T){
	a:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr1),
			Contract: common.HexToAddress(addr1),
			MultiSign: common.HexToAddress(addr1),
			Type : 1,
			Bind:false,
		},
	}
	b:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr1),
			Contract: common.HexToAddress(addr1),
			MultiSign: common.HexToAddress(addr1),
			Type : 1,
			Bind:false,
		},
	}
	c:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr2),
			Contract: common.HexToAddress(addr1),
			MultiSign: common.HexToAddress(addr1),
			Type : 1,
			Bind:false,
		},
	}
	d:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr1),
			Contract: common.HexToAddress(addr2),
			MultiSign: common.HexToAddress(addr1),
			Type : 1,
			Bind:false,
		},
	}

	e:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr1),
			Contract: common.HexToAddress(addr1),
			MultiSign: common.HexToAddress(addr2),
			Type : 1,
			Bind:false,
		},
	}

	f:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr1),
			Contract: common.HexToAddress(addr1),
			MultiSign: common.HexToAddress(addr1),
			Type : 2,
			Bind:false,
		},
	}

	g:=[]DeviceBindRecord{
		{
			Device: common.HexToAddress(addr1),
			Revenue: common.HexToAddress(addr1),
			Contract: common.HexToAddress(addr1),
			MultiSign: common.HexToAddress(addr1),
			Type : 1,
			Bind:true,
		},
	}

	err:=compareDeviceBind(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}

	test_compareDeviceBind(t,a,c)
	test_compareDeviceBind(t,a,d)
	test_compareDeviceBind(t,a,e)
	test_compareDeviceBind(t,a,f)
	test_compareDeviceBind(t,a,g)
}

func test_compareDeviceBind(t *testing.T,a []DeviceBindRecord, b []DeviceBindRecord){
	err:=compareDeviceBind(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}


func TestAlien_compareCandidatePunish(t *testing.T){
	a:=[]CandidatePunishRecord{
		{
			Target:common.HexToAddress(addr1),
			Amount:amount,
			Credit:0,
		},
	}
	b:=[]CandidatePunishRecord{
		{
			Target:common.HexToAddress(addr1),
			Amount:amount,
			Credit:0,
		},
	}
	c:=[]CandidatePunishRecord{
		{
			Target:common.HexToAddress(addr2),
			Amount:amount,
			Credit:0,
		},
	}
	d:=[]CandidatePunishRecord{
		{
			Target:common.HexToAddress(addr1),
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			Credit:0,
		},
	}
	e:=[]CandidatePunishRecord{
		{
			Target:common.HexToAddress(addr2),
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			Credit:0,
		},
	}
	f:=[]CandidatePunishRecord{
		{
			Target:common.HexToAddress(addr2),
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			Credit:10,
		},
	}
	err:=compareCandidatePunish(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareCandidatePunish(t,a,c)
	test_compareCandidatePunish(t,a,d)
	test_compareCandidatePunish(t,a,e)
	test_compareCandidatePunish(t,a,f)
}

func test_compareCandidatePunish(t *testing.T,a []CandidatePunishRecord, b []CandidatePunishRecord){
	err:=compareCandidatePunish(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareMinerStake(t *testing.T){
	a:=[]MinerStakeRecord{
		{
			Target:common.HexToAddress(addr1),
			Stake:amount,
		},
	}
	b:=[]MinerStakeRecord{
		{
			Target:common.HexToAddress(addr1),
			Stake:amount,
		},
	}
	c:=[]MinerStakeRecord{
		{
			Target:common.HexToAddress(addr2),
			Stake:amount,
		},
	}
	d:=[]MinerStakeRecord{
		{
			Target:common.HexToAddress(addr1),
			Stake:new(big.Int).Add(amount,big.NewInt(600)),
		},
	}
	e:=[]MinerStakeRecord{
		{
			Target:common.HexToAddress(addr2),
			Stake:new(big.Int).Add(amount,big.NewInt(600)),
		},
	}
	err:=compareMinerStake(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareMinerStake(t,a,c)
	test_compareMinerStake(t,a,d)
	test_compareMinerStake(t,a,e)
}

func test_compareMinerStake(t *testing.T,a []MinerStakeRecord, b []MinerStakeRecord){
	err:=compareMinerStake(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareExit(t *testing.T){
	a:=[]common.Address{
		common.HexToAddress(addr1),
	}
	b:=[]common.Address{
		common.HexToAddress(addr1),
	}
	c:=[]common.Address{
		common.HexToAddress(addr2),
	}
	err:=compareExit(a,b,"compareExit")
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareExit(t,a,c)
}

func test_compareExit(t *testing.T,a []common.Address, b []common.Address){
	err:=compareExit(a,b,"compareExit")
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}


func TestAlien_compareClaimedBandwidth(t *testing.T){
	a:=[]ClaimedBandwidthRecord{
		{
			Target:common.HexToAddress(addr1),
			Amount:amount,
			Bandwidth:1000,
		},
	}
	b:=[]ClaimedBandwidthRecord{
		{
			Target:common.HexToAddress(addr1),
			Amount:amount,
			Bandwidth:1000,
		},
	}
	c:=[]ClaimedBandwidthRecord{
		{
			Target:common.HexToAddress(addr2),
			Amount:amount,
			Bandwidth:1000,
		},
	}
	d:=[]ClaimedBandwidthRecord{
		{
			Target:common.HexToAddress(addr1),
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			Bandwidth:1000,
		},
	}
	e:=[]ClaimedBandwidthRecord{
		{
			Target:common.HexToAddress(addr2),
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			Bandwidth:1000,
		},
	}

	f:=[]ClaimedBandwidthRecord{
		{
			Target:common.HexToAddress(addr2),
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			Bandwidth:1001,
		},
	}

	err:=compareClaimedBandwidth(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareClaimedBandwidth(t,a,c)
	test_compareClaimedBandwidth(t,a,d)
	test_compareClaimedBandwidth(t,a,e)
	test_compareClaimedBandwidth(t,a,f)
}

func test_compareClaimedBandwidth(t *testing.T,a []ClaimedBandwidthRecord, b []ClaimedBandwidthRecord){
	err:=compareClaimedBandwidth(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareConfigDeposit(t *testing.T){
	a:=[]ConfigDepositRecord{
		{
			Who:0,
			Amount:amount,
		},
	}
	b:=[]ConfigDepositRecord{
		{
			Who:0,
			Amount:amount,
		},
	}
	c:=[]ConfigDepositRecord{
		{
			Who:1,
			Amount:amount,
		},
	}
	d:=[]ConfigDepositRecord{
		{
			Who:1,
			Amount:new(big.Int).Add(amount,big.NewInt(6000)),
		},
	}
	err:=compareConfigDeposit(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareConfigDeposit(t,a,c)
	test_compareConfigDeposit(t,a,d)
}

func test_compareConfigDeposit(t *testing.T,a []ConfigDepositRecord, b []ConfigDepositRecord){
	err:=compareConfigDeposit(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}




func TestAlien_compareConfigISPQOS(t *testing.T){
	a:=[]ISPQOSRecord{
		{
			ISPID:0,
			QOS:20,
		},
	}
	b:=[]ISPQOSRecord{
		{
			ISPID:0,
			QOS:20,
		},
	}
	c:=[]ISPQOSRecord{
		{
			ISPID:1,
			QOS:20,
		},
	}
	d:=[]ISPQOSRecord{
		{
			ISPID:1,
			QOS:22,
		},
	}
	err:=compareConfigISPQOS(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareConfigISPQOS(t,a,c)
	test_compareConfigISPQOS(t,a,d)
}

func test_compareConfigISPQOS(t *testing.T,a []ISPQOSRecord, b []ISPQOSRecord){
	err:=compareConfigISPQOS(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareLockParameters(t *testing.T){
	a:=[]LockParameterRecord{
		{
			LockPeriod:0,
			RlsPeriod:20,
			Interval:20,
			Who:20,
		},
	}
	b:=[]LockParameterRecord{
		{
			LockPeriod:0,
			RlsPeriod:20,
			Interval:20,
			Who:20,
		},
	}
	c:=[]LockParameterRecord{
		{
			LockPeriod:1,
			RlsPeriod:20,
			Interval:20,
			Who:20,
		},
	}
	d:=[]LockParameterRecord{
		{
			LockPeriod:0,
			RlsPeriod:22,
			Interval:20,
			Who:20,
		},
	}
	e:=[]LockParameterRecord{
		{
			LockPeriod:0,
			RlsPeriod:20,
			Interval:30,
			Who:20,
		},
	}

	f:=[]LockParameterRecord{
		{
			LockPeriod:0,
			RlsPeriod:20,
			Interval:20,
			Who:30,
		},
	}

	g:=[]LockParameterRecord{
		{
			LockPeriod:10,
			RlsPeriod:210,
			Interval:210,
			Who:301,
		},
	}

	err:=compareLockParameters(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareLockParameters(t,a,c)
	test_compareLockParameters(t,a,d)
	test_compareLockParameters(t,a,e)
	test_compareLockParameters(t,a,f)
	test_compareLockParameters(t,a,g)
}

func test_compareLockParameters(t *testing.T,a []LockParameterRecord, b []LockParameterRecord){
	err:=compareLockParameters(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareManagerAddress(t *testing.T){
	a:=[]ManagerAddressRecord{
		{
			Target:common.HexToAddress(addr1),
			Who:20,
		},
	}
	b:=[]ManagerAddressRecord{
		{
			Target:common.HexToAddress(addr1),
			Who:20,
		},
	}
	c:=[]ManagerAddressRecord{
		{
			Target:common.HexToAddress(addr2),
			Who:20,
		},
	}
	d:=[]ManagerAddressRecord{
		{
			Target:common.HexToAddress(addr2),
			Who:21,
		},
	}

	err:=compareManagerAddress(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareManagerAddress(t,a,c)
	test_compareManagerAddress(t,a,d)
}

func test_compareManagerAddress(t *testing.T,a []ManagerAddressRecord, b []ManagerAddressRecord){
	err:=compareManagerAddress(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}




func TestAlien_compareGrantProfit(t *testing.T){
	a:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:0,
			Amount:amount,
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr1),
		},
	}
	b:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:0,
			Amount:amount,
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr1),
		},
	}
	c:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr2),
			BlockNumber:0,
			Amount:amount,
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr1),
		},
	}
	d:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:1,
			Amount:amount,
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr1),
		},
	}

	e:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:0,
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr1),
		},
	}

	f:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:0,
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			RevenueAddress:common.HexToAddress(addr2),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr1),
		},
	}
	g:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:0,
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr3),
			MultiSignature:common.HexToAddress(addr1),
		},
	}

	h:=[]consensus.GrantProfitRecord{
		{
			Which:20,
			MinerAddress:common.HexToAddress(addr1),
			BlockNumber:0,
			Amount:new(big.Int).Add(amount,big.NewInt(600)),
			RevenueAddress:common.HexToAddress(addr1),
			RevenueContract:common.HexToAddress(addr1),
			MultiSignature:common.HexToAddress(addr3),
		},
	}

	err:=compareGrantProfit(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareGrantProfit(t,a,c)
	test_compareGrantProfit(t,a,d)
	test_compareGrantProfit(t,a,e)
	test_compareGrantProfit(t,a,f)
	test_compareGrantProfit(t,a,g)
	test_compareGrantProfit(t,a,h)
}

func test_compareGrantProfit(t *testing.T,a []consensus.GrantProfitRecord, b []consensus.GrantProfitRecord){
	err:=compareGrantProfit(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareFlowReport(t *testing.T){
	a:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:99,
					FlowValue2:263,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}
	b:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:99,
					FlowValue2:263,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}
	c:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{002},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:99,
					FlowValue2:263,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}
	d:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6667,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:99,
					FlowValue2:263,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}

	e:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr2),
					FlowValue1:99,
					FlowValue2:263,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}

	f:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:99,
					FlowValue2:263,
					Miner: common.HexToAddress(addr1),
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}
	g:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:991,
					FlowValue2:263,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}

	h:=[]MinerPofReportRecord{
		{
			ChainHash:common.Hash{001},
			ReportTime:6666,
			ReportContent:[]MinerPofReportItem{
				{
					Target:common.HexToAddress(addr1),
					FlowValue1:99,
					FlowValue2:2634,
				},{
					Target:common.HexToAddress(addr1),
					FlowValue1:999,
					FlowValue2:263,
				},
			},
		},
	}

	err:=compareFlowReport(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareFlowReport(t,a,c)
	test_compareFlowReport(t,a,d)
	test_compareFlowReport(t,a,e)
	test_compareFlowReport(t,a,f)
	test_compareFlowReport(t,a,g)
	test_compareFlowReport(t,a,h)
}

func test_compareFlowReport(t *testing.T,a []MinerPofReportRecord, b []MinerPofReportRecord){
	err:=compareFlowReport(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}

func TestAlien_compareMinerFlowReportItem(t *testing.T){
	a:=[]MinerPofReportItem{
		{
			Target:common.HexToAddress(addr1),
			FlowValue1:99,
			FlowValue2:263,
		},
	}
	b:=[]MinerPofReportItem{
		{
			Target:common.HexToAddress(addr1),
			FlowValue1:99,
			FlowValue2:263,
		},
	}
	c:=[]MinerPofReportItem{
		{
			Target:common.HexToAddress(addr2),
			FlowValue1:99,
			FlowValue2:263,
		},
	}
	d:=[]MinerPofReportItem{
		{
			Target:common.HexToAddress(addr1),
			FlowValue1:99,
			FlowValue2:263,
			Miner: common.HexToAddress(addr1),
		},
	}

	e:=[]MinerPofReportItem{
		{
			Target:common.HexToAddress(addr1),
			FlowValue1:991,
			FlowValue2:263,
		},
	}

	f:=[]MinerPofReportItem{
		{
			Target:common.HexToAddress(addr1),
			FlowValue1:99,
			FlowValue2:261,
		},
	}

	err:=compareMinerFlowReportItem(a,b)
	if err==nil{
		t.Logf(" pass")
	}else{
		t.Errorf(" error,expect nil,but act %s" ,err.Error())
	}
	test_compareMinerFlowReportItem(t,a,c)
	test_compareMinerFlowReportItem(t,a,d)
	test_compareMinerFlowReportItem(t,a,e)
	test_compareMinerFlowReportItem(t,a,f)
}

func test_compareMinerFlowReportItem(t *testing.T,a []MinerPofReportItem, b []MinerPofReportItem){
	err:=compareMinerFlowReportItem(a,b)
	if err!=nil{
		t.Logf(" pass,error msg:%s" ,err.Error())
	}else{
		t.Errorf(" error,expect not nil,but act %s" ,err.Error())
	}
}



