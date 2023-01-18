package alien

import (
	"fmt"
	"github.com/token/common"
	"github.com/token/params"
	"math/big"
	"testing"
)

func TestQueueElect(t *testing.T) {
	snap := &Snapshot{
		config:   &params.AlienConfig{MaxSignerCount: 29},
		LCRS:     1,
		Tally:    make(map[common.Address]*big.Int),
		TallyMiner: make(map[common.Address]*CandidateState),
		Punished: make(map[common.Address]uint64),
		PosPledge:make(map[common.Address]*PosPledgeItem),
		TallySigner: make(map[common.Address]uint64 ),
		Hash: common.HexToHash("0xbcd20b9fb6daef689393e9f2f6e65dbdc100a87b6c54d5e2f3fb5aae1073c465"),
	}
	snap.initTally()
	snap.initTallyMiner()
	snap.initTallyPunish()
	snap.initHistoryHash()
	snap.initPosPledge(1)
	for i:=1;i<300;i++{
		if i%290>0 {
			continue
		}
		//if i==3013920 {
		//	snap.Tally[common.HexToAddress("0xdc5f4def1bba66481ebf3c73ecbda7a29e278aaf")]=big.NewInt(100000)
		//	snap.TallyMiner[common.HexToAddress("0x447286d98430736f5c4e1e2f23a782e6d8da5daf")].Stake=big.NewInt(100000)
		//}
		squeue,err:=snap.elect(21,uint64(i-1))
		snap.updateSignerNumber(squeue,uint64(i+1))
		//fmt.Println("**********************************************")
		if err ==nil {
			//	for index,addr:=range squeue{
			//		ismain:=false
			//		if _,ok := snap.Tally[addr];ok {
			//			ismain=true
			//		}
			//		fmt.Println("Queue","index",index,"ismain",ismain,addr)
			//	}
		}else {
			fmt.Println("err",err)
		}
	}
	for addr,state:= range snap.TallySigner{
		fmt.Println("main",addr,state)
	}
	for addr,item:= range snap.TallyMiner{
		fmt.Println("miner",addr,item.SignerNumber)
	}
}
func (snap *Snapshot)elect(maxSignerCount uint64,number uint64) ([]common.Address, error){
	snap.Number=number
	candidateNeedPD=false
	return snap.createSignerQueue()
}
func (snap *Snapshot) initHistoryHash(){
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xfdff7a18e794e5f3e06d89f69eb11bc464d5422adda48ee43354d5764f626781"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x7522975e5bd27d11a55b4cc7a85b0bd2230e424fc4256fe3537f8a31f34c00a5"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xad55860b5eecb397f8cfedc32050d63096f08a61ed08142d49c5dab8da627df4"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x6e545846777c058d367229dfa7ec8acb035723d8dc11450fe079e6ac9554bb5b"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x97e9777f9a49a22dd78d266a1799aefefcfc9ee985eeabf763a34700cd401101"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x509d07ba48a3d8e03033fde8748d2db07ede51d5a57cbf95039cb742fa534fbb"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xdc43a160d98736aa5b55b6beecfe2f2a07e5d1e261cacd25a86f1a0923ed917d"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x0846cc49351e513f7d95ddac060966c56fa560736d6608821fac6979c684e0b4"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x32d3275325a5efac9a61aae6019709040e7d7073fee4109d7bb7b495b9a42c86"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xb0ca74f89f42c74af8a35054dbcca1e6d9bc80a3f0d05be2ffeb0406478a0a00"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xb97cd415ad5c5e607c5be2d12d9b2d946016dcc9079484c1e148e442b72bb61f"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x8355bc6da3a3c75c4b7d80da10570535c9c57bceef767319625a353ca02f94aa"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xe01427193504e93b2c5cd341d2f310ef356c9ec118beb79be3dd8936932d4ead"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xf510fbf295b4c387f19027aba63a08f0fd19210dddb9f1bddf9c7505b891b166"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xe8c08337a8d3ef79d23a4a1bab2e090c6b5e2a868a763535747f3f36d83ec871"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xe20ff60deab64515e86971196c341aa19f170a0adc776edb75e02bf7cdbbf913"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x853a00d9595438aa352ad7209ecf44a92b7416238f032d43f166c1fb7bd11865"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x41194d4e005d2fe7665c62967f4239934adfadc9e909382814df15d6edad4b7c"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x14e13fe0dc722d124e499dd1cda77d274a3ee6fcbaa9c9b6bff773335b1679d7"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xd36e7274f810a3b764636b9079815fd3690b9e962e55a3df203cb52b5a77b7d2"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xd5ef0eeaa297791e0d4b7adc2ca6e3d2cec232ceb6088576109320451f9882b7"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x7d9c8373c096960dfbc7e8f0961be1ba4faba3d908b9b9cb2410852e73d3169f"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xb944428551aee8c0204ed8a6ecf25ab7e307cff621dc4018456e1d5b232b1c3d"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x83a5ecc265f0a9bf93aacd3665a51173893929123f373fd6c9e2da4291a98a56"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xfb889656240f8a9434fda5702f84202e3661d8bb42f83284d23959b9a1bd4474"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xfb3efa5d430173d7ea5827e0fcd7291bf85df42523d346f6619de27ec7068a85"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x03b30ac11cbd56d465741861f653c54251b7fca05c1c0edec4620a38f5e877e2"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x19f03a126a794766f74dfded7942527e480431cc7adf5be5a7281b56597fc4f1"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x6bd2efef61895f86d7b8d4a84bea0c4b659e5da3f5f6af56ae491748eecabdfd"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xbf18cae8d6c3b0eabaee65ff0afd528f75a8fbba8da1eed07cf28c6a2edb4d43"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xc4de32934d543765cd37e314943a84a6dd818a2144066100c10b7ac2288c0421"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x65a9818344bed14ff03a606c9c62b3fab109f5f948c3821c02a6b4c3235cc180"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x42b044c0492b22de43cc7d61bef1f5313947aa6a1cf0a16d4922096cef4af0b0"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x904f2c86d6d6b2e6d8c1f560e55b5ed9224e45809eb57a7f2173291ec072eb5c"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xa77b79ead180aaac4409e0100e6940fbecc210eee6881820b910c1d913e3aeed"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xeedde9139dc10fc0b9d95689368fc88f5a1ab258b0bd833b7182583fa27f5132"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x5768c949c5fdf7c05422cfdd0914c70a9d60186dcb6cd698d325974654b89d67"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x2e928f6759c9d3ee5d1410ab60aea03af70a79343b087383fe5e0b53eeaf557e"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xec2b02b62b4ed456971dd60f6c5976eee883f80d0f0a90a41dcee72eff1724a3"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xac91c1296cf4bf24d01af17fc195e26d097135080aa4349d0a726182617c17d2"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0x6649845774e94c53f18fd1090a0c360022b5119a1120c72000253f58acd9e8f5"))
	snap.HistoryHash= append(snap.HistoryHash,common.HexToHash("0xbcd20b9fb6daef689393e9f2f6e65dbdc100a87b6c54d5e2f3fb5aae1073c465"))
}
func (snap *Snapshot) initTally(){

	snap.Tally[common.HexToAddress("0xf0b5e4bde593251dd9a0663add673aef7b915062")]=big.NewInt(500)
	snap.Tally[common.HexToAddress("0xa1ca30cdb6d7faf4ff380aa210a92fe925d47968")]=big.NewInt(1090)
	snap.Tally[common.HexToAddress("0xe6a6f49817bfd9fdf4d0903d82f29a6aa4d9e3bd")]=big.NewInt(1080)
	snap.Tally[common.HexToAddress("0xd6c00e9bc15d87cc68ccb0cef9d802ae20c37fe6")]=big.NewInt(1050)
	snap.Tally[common.HexToAddress("0x83d6e3461aef8bfac878b9bb37e18c78954cb781")]=big.NewInt(1040)
	snap.Tally[common.HexToAddress("0x68cc6205083fb860cf13d4c0d9814248d72cb249")]=big.NewInt(1030)
	snap.Tally[common.HexToAddress("0x3af094ba2e5f0a9e56923226e487541c0ea8ff8f")]=big.NewInt(1060)
	snap.Tally[common.HexToAddress("0xea1aca7b543b5729c36dafff27ec35d84f2d0e20")]=big.NewInt(1020)
	snap.Tally[common.HexToAddress("0xfca289398fe03747d7c0708526686a01d9cf2a06")]=big.NewInt(1008)
	snap.Tally[common.HexToAddress("0x68cc6205083fb860cf13d4c0d9814248d72cb240")]=big.NewInt(1700)
	snap.Tally[common.HexToAddress("0x7b6af8663a821ea474c08c458618e64c76dc769a")]=big.NewInt(1020)
	snap.Tally[common.HexToAddress("0x8176a171fa5b09071af3cfce878c5044aefaecbf")]=big.NewInt(1000)
	snap.Tally[common.HexToAddress("0x45144c4084b1c87e7a3cd9df68f71fb97e261745")]=big.NewInt(1000)
	snap.Tally[common.HexToAddress("0x45144c4084b1c87e7a3cd9df68f71fb97e261746")]=big.NewInt(1008)
	snap.Tally[common.HexToAddress("0x45144c4084b1c87e7a3cd9df68f71fb97e261749")]=big.NewInt(1038)
	snap.Tally[common.HexToAddress("0xdc5f4def1bba66481ebf3c73ecbda7a29e278aaf")]=big.NewInt(6)
}

func (snap *Snapshot) initTallyMiner(){

	snap.TallyMiner[common.HexToAddress("0xab7fe746021b17c3bdab8327a78792d7f34533ec")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x15f4e608027052395d5c5c422f27937781563d7d")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(9004)}
	snap.TallyMiner[common.HexToAddress("0xc7fbdb9d56f13ea4486fd872200be2090c23fde0")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1003)}
	snap.TallyMiner[common.HexToAddress("0xf6a5ed52afb23cb8a577af6948cf7a651f561181")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1002)}
	snap.TallyMiner[common.HexToAddress("0xd695e08779a69c534af13178d8cd5ee9a332d00d")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1001)}
	snap.TallyMiner[common.HexToAddress("0x30a93bdf538c06045d5fa02f32bf6e7ff0006904")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0xdf7e12d75ae88b8013c763ace91700da5ca5e872")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0xcbd9e67e2e3f50983786d2c727e0a4e7fea9bc8e")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x3d040d7246327ee9e5d7ad4ff9a0e44cb7532f10")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x3f983bbea130b8eee8a168bc0438278711000b40")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x9d7bcf504b9947941721e052392989a2d9a85116")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x6ac1463f15bc965597e88b231ae9f9c0f7f6d6cb")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x619769c8996a236babd556ba26e5c80b5cb4c75f")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x619769c8996a236babd556ba26e5c80b5cb4htyd")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(1000)}
	snap.TallyMiner[common.HexToAddress("0x447286d98430736f5c4e1e2f23a782e6d8da5daf")]=&CandidateState{SignerNumber:0,Stake:big.NewInt(0)}
}

func (snap *Snapshot) initTallyPunish(){
	snap.Punished[common.HexToAddress("0xab7fe746021b17c3bdab8327a78792d7f34533ec")]=10
}
