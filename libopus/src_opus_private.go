package libopus

import "unsafe"

const MODE_SILK_ONLY = 1000
const MODE_HYBRID = 1001
const MODE_CELT_ONLY = 1002
const OPUS_SET_VOICE_RATIO_REQUEST = 11018
const OPUS_GET_VOICE_RATIO_REQUEST = 11019
const OPUS_SET_FORCE_MODE_REQUEST = 11002

type OpusRepacketizer struct {
	Toc       uint8
	Nb_frames int64
	Frames    [48]*uint8
	Len       [48]opus_int16
	Framesize int64
}
type ChannelLayout struct {
	Nb_channels        int64
	Nb_streams         int64
	Nb_coupled_streams int64
	Mapping            [256]uint8
}
type MappingType int64

const (
	MAPPING_TYPE_NONE = MappingType(iota)
	MAPPING_TYPE_SURROUND
	MAPPING_TYPE_AMBISONICS
)

type OpusMSEncoder struct {
	Layout            ChannelLayout
	Arch              int64
	Lfe_stream        int64
	Application       int64
	Variable_duration int64
	Mapping_type      MappingType
	Bitrate_bps       opus_int32
}
type OpusMSDecoder struct {
	Layout ChannelLayout
}
type opus_copy_channel_in_func func(dst *opus_val16, dst_stride int64, src unsafe.Pointer, src_stride int64, src_channel int64, frame_size int64, user_data unsafe.Pointer)
type opus_copy_channel_out_func func(dst unsafe.Pointer, dst_stride int64, dst_channel int64, src *opus_val16, src_stride int64, frame_size int64, user_data unsafe.Pointer)
type downmix_func func(unsafe.Pointer, *opus_val32, int64, int64, int64, int64, int64)

func align(i int64) int64 {
	type foo struct {
		C int8
		U struct {
			// union
			P unsafe.Pointer
			I opus_int32
			V opus_val32
		}
	}
	var alignment uint64 = uint64(_cxgo_offsetof(foo{}, "#member"))
	return int64(((uint64(i) + alignment - 1) / alignment) * alignment)
}
