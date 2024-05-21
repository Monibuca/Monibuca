package codec

import "fmt"

type H265NALUType byte

func (H265NALUType) Parse(b byte) H265NALUType {
	return H265NALUType(b & 0x7E >> 1)
}

func ParseH265NALUType(b byte) H265NALUType {
	return H265NALUType(b & 0x7E >> 1)
}

const (
	// HEVC_VPS    = 0x40
	// HEVC_SPS    = 0x42
	// HEVC_PPS    = 0x44
	// HEVC_SEI    = 0x4E
	// HEVC_IDR    = 0x26
	// HEVC_PSLICE = 0x02

	NAL_UNIT_CODED_SLICE_TRAIL_N H265NALUType = iota // 0
	NAL_UNIT_CODED_SLICE_TRAIL_R                     // 1
	NAL_UNIT_CODED_SLICE_TSA_N                       // 2
	NAL_UNIT_CODED_SLICE_TLA                         // 3 // Current name in the spec: TSA_R
	NAL_UNIT_CODED_SLICE_STSA_N                      // 4
	NAL_UNIT_CODED_SLICE_STSA_R                      // 5
	NAL_UNIT_CODED_SLICE_RADL_N                      // 6
	NAL_UNIT_CODED_SLICE_DLP                         // 7 // Current name in the spec: RADL_R
	NAL_UNIT_CODED_SLICE_RASL_N                      // 8
	NAL_UNIT_CODED_SLICE_TFD                         // 9 // Current name in the spec: RASL_R
	NAL_UNIT_RESERVED_10
	NAL_UNIT_RESERVED_11
	NAL_UNIT_RESERVED_12
	NAL_UNIT_RESERVED_13
	NAL_UNIT_RESERVED_14
	NAL_UNIT_RESERVED_15
	NAL_UNIT_CODED_SLICE_BLA      // 16 // Current name in the spec: BLA_W_LP
	NAL_UNIT_CODED_SLICE_BLANT    // 17 // Current name in the spec: BLA_W_DLP
	NAL_UNIT_CODED_SLICE_BLA_N_LP // 18
	NAL_UNIT_CODED_SLICE_IDR      // 19// Current name in the spec: IDR_W_DLP
	NAL_UNIT_CODED_SLICE_IDR_N_LP // 20
	NAL_UNIT_CODED_SLICE_CRA      // 21
	NAL_UNIT_RESERVED_22
	NAL_UNIT_RESERVED_23
	NAL_UNIT_RESERVED_24
	NAL_UNIT_RESERVED_25
	NAL_UNIT_RESERVED_26
	NAL_UNIT_RESERVED_27
	NAL_UNIT_RESERVED_28
	NAL_UNIT_RESERVED_29
	NAL_UNIT_RESERVED_30
	NAL_UNIT_RESERVED_31
	NAL_UNIT_VPS                   // 32
	NAL_UNIT_SPS                   // 33
	NAL_UNIT_PPS                   // 34
	NAL_UNIT_ACCESS_UNIT_DELIMITER // 35
	NAL_UNIT_EOS                   // 36
	NAL_UNIT_EOB                   // 37
	NAL_UNIT_FILLER_DATA           // 38
	NAL_UNIT_SEI                   // 39 Prefix SEI
	NAL_UNIT_SEI_SUFFIX            // 40 Suffix SEI
	NAL_UNIT_RESERVED_41
	NAL_UNIT_RESERVED_42
	NAL_UNIT_RESERVED_43
	NAL_UNIT_RESERVED_44
	NAL_UNIT_RESERVED_45
	NAL_UNIT_RESERVED_46
	NAL_UNIT_RESERVED_47
	NAL_UNIT_RTP_AP
	NAL_UNIT_RTP_FU
	NAL_UNIT_UNSPECIFIED_50
	NAL_UNIT_UNSPECIFIED_51
	NAL_UNIT_UNSPECIFIED_52
	NAL_UNIT_UNSPECIFIED_53
	NAL_UNIT_UNSPECIFIED_54
	NAL_UNIT_UNSPECIFIED_55
	NAL_UNIT_UNSPECIFIED_56
	NAL_UNIT_UNSPECIFIED_57
	NAL_UNIT_UNSPECIFIED_58
	NAL_UNIT_UNSPECIFIED_59
	NAL_UNIT_UNSPECIFIED_60
	NAL_UNIT_UNSPECIFIED_61
	NAL_UNIT_UNSPECIFIED_62
	NAL_UNIT_UNSPECIFIED_63
	NAL_UNIT_INVALID
)

var AudNalu = []byte{0x00, 0x00, 0x00, 0x01, 0x46, 0x01, 0x10}

type (
	IH265Ctx interface {
		GetH265Ctx() *H265Ctx
	}
	H265Ctx struct {
		H264Ctx
		VPS [][]byte
	}
)

func (ctx *H265Ctx) GetInfo() string {
	return fmt.Sprintf("sps: % 02X,pps: % 02X,vps: % 02X", ctx.SPS[0], ctx.PPS[0], ctx.VPS[0])
}

func (h265 *H265Ctx) GetHeight() int {
	return int(h265.Height)
}

func (h265 *H265Ctx) GetWidth() int {
	return int(h265.Width)
}

func (*H265Ctx) FourCC() FourCC {
	return FourCC_H265
}

func (h265 *H265Ctx) GetH265Ctx() *H265Ctx {
	return h265
}
