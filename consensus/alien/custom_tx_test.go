package alien

import (
	"github.com/token/common"
	"github.com/token/core/state"
	"github.com/token/core/types"
	"strings"
	"testing"
)
func TestAlien_checkRevenueNormalBind(t *testing.T) {
	device:=common.HexToAddress("NX1E0E2B42595Cb6046566F77Fb0c67a9D109aBE1D")
	revenue:=common.HexToAddress("NXa63b29EBe0A141B87A87e39dE17F17346e11e1b7")
	deviceBind:=DeviceBindRecord{
		Device:device,
		Revenue:revenue,
		Type:0,
	}
	alien:=&Alien{
	}
	snap :=&Snapshot{
		RevenueNormal:make(map[common.Address]*RevenueParameter),
	}
	err:=alien.checkRevenueNormalBind(deviceBind,snap)
	if err==nil{
		t.Logf("checkRevenueNormalBind pass")
	}else {
		t.Errorf("checkRevenueNormalBind fail,error msg:"+err.Error())
	}
	snap.RevenueNormal[device]=&RevenueParameter{
		RevenueAddress:revenue,
	}
	err=alien.checkRevenueNormalBind(deviceBind,snap)
	if err!=nil{
		t.Logf("checkRevenueNormalBind pass,error msg:"+err.Error())
	}else {
		t.Errorf("checkRevenueNormalBind fail,error msg is empty:")
	}
	deviceBind.Type=1
	err=alien.checkRevenueNormalBind(deviceBind,snap)
	if err==nil{
		t.Logf("checkRevenueNormalBind pass")
	}else {
		t.Errorf("checkRevenueNormalBind fail,error msg:"+err.Error())
	}
}


func TestAlien_processDeviceBind(t *testing.T) {

	currentDeviceBind:=make([]DeviceBindRecord,0)
	dev:="bec92229b1bd96919c8ffc993171fa6504121dc6"
	devAddr:=common.HexToAddress(dev)
	txData := "token:1:Bind:"+dev+":1:0000000000000000000000000000000000000000:0000000000000000000000000000000000000000"
	txDataInfo := strings.Split(txData, ":")
	txSender:= common.HexToAddress("NXa63b29EBe0A141B87A87e39dE17F17346e11e1b7")
	tx:=&types.Transaction{}
	receipts:=make([]*types.Receipt,0)
	snap :=&Snapshot{
	  RevenueNormal: make(map[common.Address]*RevenueParameter),
	  RevenuePof:    make(map[common.Address]*RevenueParameter),
	}
	alien:=&Alien{
	}
	currentDeviceBind=alien.processDeviceBind(currentDeviceBind,txDataInfo,txSender,tx,receipts,snap)
	for index := range currentDeviceBind {
		if txSender==currentDeviceBind[index].Revenue&&devAddr==currentDeviceBind[index].Device{
			t.Logf("1 pass,Revenue=%s,Device=%s" ,currentDeviceBind[index].Revenue.String(),currentDeviceBind[index].Device.String())
		}else {
			t.Errorf("1 fail,Revenue=%s,Device=%s" ,currentDeviceBind[index].Revenue.String(),currentDeviceBind[index].Device.String())
		}
	}
	txData= "token:1:Bind:"+dev+":0:0000000000000000000000000000000000000000:0000000000000000000000000000000000000000"
	txDataInfo= strings.Split(txData, ":")
	currentDeviceBind=make([]DeviceBindRecord,0)
	currentDeviceBind=alien.processDeviceBind(currentDeviceBind,txDataInfo,txSender,tx,receipts,snap)
	for index := range currentDeviceBind {
		if txSender==currentDeviceBind[index].Revenue&&devAddr==currentDeviceBind[index].Device{
			t.Logf("2 pass,Revenue=%s,Device=%s" ,currentDeviceBind[index].Revenue.String(),currentDeviceBind[index].Device.String())
		}else {
			t.Errorf("2 fail,Revenue=%s,Device=%s" ,currentDeviceBind[index].Revenue.String(),currentDeviceBind[index].Device.String())
		}
	}
	txData= "token:1:Bind:"+dev+":0:0000000000000000000000000000000000000000:0000000000000000000000000000000000000000"
	txDataInfo= strings.Split(txData, ":")
	currentDeviceBind=make([]DeviceBindRecord,0)
	currentDeviceBind=alien.processDeviceBind(currentDeviceBind,txDataInfo,txSender,tx,receipts,snap)
	if len(currentDeviceBind)==0{
		t.Logf("3 pass" )
	}else{
		t.Errorf("3 fail")
	}

	state := &state.StateDB{
	}
	currentDeviceBind=alien.processDeviceRebind(currentDeviceBind,txDataInfo,txSender,tx,receipts,state,snap)
	if len(currentDeviceBind)==0{
		t.Logf("4 pass" )
	}else{
		t.Errorf("4 fail")
	}

	newrev:="0Ff6e773Ff893fF39ed9352160889df13BDfc896"
	newrevAddr:=common.HexToAddress(newrev)
	txData= "token:1:Rebind:"+dev+":0:0000000000000000000000000000000000000000:0000000000000000000000000000000000000000:"+newrev
	txDataInfo= strings.Split(txData, ":")
	currentDeviceBind=alien.processDeviceRebind(currentDeviceBind,txDataInfo,txSender,tx,receipts,state,snap)
	for index := range currentDeviceBind {
		if newrevAddr==currentDeviceBind[index].Revenue&&devAddr==currentDeviceBind[index].Device{
			t.Logf("5 pass,Revenue=%s,Device=%s" ,currentDeviceBind[index].Revenue.String(),currentDeviceBind[index].Device.String())
		}else {
			t.Errorf("5 fail,Revenue=%s,Device=%s" ,currentDeviceBind[index].Revenue.String(),currentDeviceBind[index].Device.String())
		}
	}

}