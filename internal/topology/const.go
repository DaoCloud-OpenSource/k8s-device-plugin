package topology

import "strings"

type LinkType uint

const (
	P2PLinkUnknown       LinkType = 0
	P2PLinkCrossCPU      LinkType = 10
	P2PLinkSameCPU       LinkType = 20
	P2PLinkHostBridge    LinkType = 30
	P2PLinkMultiSwitch   LinkType = 40
	P2PLinkSingleSwitch  LinkType = 50
	P2PLinkSameBoard     LinkType = 60
	SingleNVLINKLink     LinkType = 100
	TwoNVLINKLinks       LinkType = 110
	ThreeNVLINKLinks     LinkType = 120
	FourNVLINKLinks      LinkType = 130
	FiveNVLINKLinks      LinkType = 140
	SixNVLINKLinks       LinkType = 150
	SevenNVLINKLinks     LinkType = 160
	EightNVLINKLinks     LinkType = 170
	NineNVLINKLinks      LinkType = 180
	TenNVLINKLinks       LinkType = 190
	ElevenNVLINKLinks    LinkType = 200
	TwelveNVLINKLinks    LinkType = 210
	ThirteenNVLINKLinks  LinkType = 220
	FourteenNVLINKLinks  LinkType = 230
	FifteenNVLINKLinks   LinkType = 230
	SixteenNVLINKLinks   LinkType = 240
	SeventeenNVLINKLinks LinkType = 250
	EighteenNVLINKLinks  LinkType = 260
)

var (
	LinkTypeMapStr = make(map[LinkType]string)
)

func init() {
	LinkTypeMapStr[P2PLinkUnknown] = "N/A"
	LinkTypeMapStr[P2PLinkCrossCPU] = "SYS"
	LinkTypeMapStr[P2PLinkSameCPU] = "NODE"
	LinkTypeMapStr[P2PLinkHostBridge] = "PHB"
	LinkTypeMapStr[P2PLinkMultiSwitch] = "PXB"
	LinkTypeMapStr[P2PLinkSingleSwitch] = "PIX"
	LinkTypeMapStr[P2PLinkSameBoard] = "X"
	LinkTypeMapStr[SingleNVLINKLink] = "NV1"
	LinkTypeMapStr[TwoNVLINKLinks] = "NV2"
	LinkTypeMapStr[ThreeNVLINKLinks] = "NV3"
	LinkTypeMapStr[FourNVLINKLinks] = "NV4"
	LinkTypeMapStr[FiveNVLINKLinks] = "NV5"
	LinkTypeMapStr[SixNVLINKLinks] = "NV6"
	LinkTypeMapStr[SevenNVLINKLinks] = "NV7"
	LinkTypeMapStr[EightNVLINKLinks] = "NV8"
	LinkTypeMapStr[NineNVLINKLinks] = "NV9"
	LinkTypeMapStr[TenNVLINKLinks] = "NV10"
	LinkTypeMapStr[ElevenNVLINKLinks] = "NV11"
	LinkTypeMapStr[TwelveNVLINKLinks] = "NV12"
	LinkTypeMapStr[ThirteenNVLINKLinks] = "NV13"
	LinkTypeMapStr[FourteenNVLINKLinks] = "NV14"
	LinkTypeMapStr[FifteenNVLINKLinks] = "NV15"
	LinkTypeMapStr[SixteenNVLINKLinks] = "NV16"
	LinkTypeMapStr[SeventeenNVLINKLinks] = "NV17"
	LinkTypeMapStr[EighteenNVLINKLinks] = "NV18"
}

func LinkTypeToStr(link LinkType) string {
	if v, ok := LinkTypeMapStr[link]; ok {
		return v
	}
	return ""
}

func StrToLinkType(str string) LinkType {
	str = strings.TrimSpace(str)
	for k, v := range LinkTypeMapStr {
		if v == str {
			return k
		}
	}
	return P2PLinkUnknown
}
