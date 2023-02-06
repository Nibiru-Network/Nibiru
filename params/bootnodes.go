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

package params

import "github.com/token/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main nbn network.
var MainnetBootnodes = []string{
	// nbn Foundation Go Bootnodes
	"enode://4fe24d2b2a995e6dd63912faad9e2e402fa62dc069cf3cd532d2779ec74b8fa26f00d35515604e2b23724c6d79346e0323aa98ce1af31ee140e6eb22bb3ddfad@52.69.72.248:30311",
	"enode://7d0eeb4917b02031b9be5418d50604b49c7479a88a3b33775a8c0146ad56203b1cf1483b066e7c02a3a861380e55de1a46a39673721f68eaaf0a5f990d09718c@52.193.245.223:30311",
	"enode://a6009b2bf3d5b958abd03b2653290807a02b48246a9e6128ca2d24b5b346024f1b5f62e61249d60cd4c1b5bac63482120743ed0750adeff733807efdc310029c@54.151.221.69:30311",
	"enode://653cd7f0ee2d3a1b53efb342dc062c10f5b31a8cd41c86c3bde312da94b78f1165880b4971384737d00d128f3a200dc58c2aa9cb8946ec4c80d930c8356c2739@54.168.157.167:30311",
	"enode://cf9a8e94f87cda4523bb4a427f0214a7dc1af3a632422c4933e3b6fcbc5b4c46574c612a804e585420eccfdf1247c9d47e6f0b211d632bcb5ea581bac63a6406@54.248.6.95:30311",
	"enode://b665a3b73a011411a1b45dd4de284bba21a5cbe46697798024d9469d01241f7cf2614f2b6f7793045c7d386ff897c839a42586f87f6f637a351f844147f8503d@64.62.206.242:30311",
	"enode://123b44cfb4b96e9818eed1f758a0b496240a2340f16f52ceba2a6ac2ae6943faab4a2236e36eaa63af35f72d355c78154b52cc8219e8d8d7ea4e75a9a51fe017@64.62.206.242:30312",
	"enode://cb9297882254fd1ac8422688a56a901d31ee3bd699e344e68fbabe3279fa97332d2f12d24d510b47a61d8f4db88b5c9fad81f1b467d718c1c9cc7ddeddccd5b7@64.62.206.242:30313",
	"enode://257f3624df3de84fcc46e8b38fe0014c34d5d3ff44f01b5094d85fac7fe9d8fe9be6233f1c45dcaaca46c05902a5852bec5b76ec78883cea62ae351af19040da@64.62.206.242:30314",
	"enode://dcae55e7aa217af7602a1401a7242d9754b82b05150736977586c37557d1397e07a7aefab95a32cb65d1320cbca0c01b3db5abba6241fac49101423ca7e6338d@64.62.206.242:30315",
	"enode://76ec54e884d391e102c0a175fff55d99cf677f8098b41a0cce153df667dcc323734fa786c1cbd4899e4a08455656e66a998849d21ed23f283fbfff6c5ab54598@64.62.206.242:30316",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Testnet test network.
var TestnetBootnodes = []string{
	"enode://2ffed1bb6b475259c1448dc93b639569886999e51ade144451877a706d2a9b71eff8eb067d289fde48ba4807370034d851553746fac8816af27f5a922703e2e4@127.0.0.1:30311",
}

var V5Bootnodes = []string{
}

const dnsPrefix = "enrtree://AKA3AM6LPBYEUDMVNU3BSVQJ5AD45Y7YPOHJLEF6W26QOE4VTUDPE@"

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	var net string
	switch genesis {
	case MainnetGenesisHash:
		net = "mainnet"
	case TestnetGenesisHash:
		net = "testnet"
	default:
		return ""
	}
	return dnsPrefix + protocol + "." + net + ".ethdisco.net"
}
