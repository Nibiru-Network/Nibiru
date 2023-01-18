package alien

import (
	"github.com/token/common"
	"math/big"
	"testing"
)

func TestAlien_PayPeriodAmount(t *testing.T) {
	 tests := []struct {
		blockNumber int64
		 playment int64
		 expectValue int64
	}{
		 {
			 blockNumber: 100,
			 playment :0,
			 expectValue:0,
		 },{
			 blockNumber: 3120,
			 playment :0,
			 expectValue:1000,
		 },{
			 blockNumber: 3140,
			 playment :1000,
			 expectValue:1000,
		 },{
			 blockNumber: 3160,
			 playment :2000,
			 expectValue:1000,
		 },{
			 blockNumber: 3180,
			 playment :3000,
			 expectValue:1000,
		 },{
			 blockNumber: 3200,
			 playment :4000,
			 expectValue:1000,
		 },{
			 blockNumber: 3220,
			 playment :5000,
			 expectValue:1000,
		 },{
			 blockNumber: 3240,
			 playment :6000,
			 expectValue:1000,
		 },{
			 blockNumber: 3260,
			 playment :7000,
			 expectValue:1000,
		 },{
			 blockNumber: 3280,
			 playment :8000,
			 expectValue:1000,
		 },{
			 blockNumber: 3300,
			 playment :9000,
			 expectValue:1000,
		 },{
			 blockNumber: 3320,
			 playment :10000,
			 expectValue:0,
		 },
	 }
	for _, tt := range tests {

		pitem := &PledgeItem{
			Amount:          big.NewInt(10000),
			PledgeType:      3,
			Playment:        big.NewInt(tt.playment),
			LockPeriod:      3000,
			RlsPeriod:       200,
			Interval:        10,
			StartHigh:       100,
			TargetAddress:   common.Address{},
			RevenueAddress:  common.Address{},
			RevenueContract: common.Address{},
			MultiSignature:  common.Address{},
		}
		amount :=caclPayPeriodAmount(pitem, big.NewInt(tt.blockNumber))
		if amount.Cmp(big.NewInt(tt.expectValue))==0  ||(tt.expectValue==0 && amount.Cmp(big.NewInt(0))<0){
			t.Logf("blocknumber=%dth, caclPayPeriodAmount success amount= %d",tt.blockNumber,amount)
		}else  {

			t.Errorf("blocknumber=%dth, caclPayPeriodAmount Failed amount= %d expectValue=%d ",tt.blockNumber,amount,tt.expectValue)
		}

	}
}
