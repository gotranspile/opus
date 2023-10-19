package libopus

const FLAG_DECODE_NORMAL = 0
const FLAG_PACKET_LOST = 1
const FLAG_DECODE_LBRR = 2

type silk_EncControlStruct struct {
	NChannelsAPI              opus_int32
	NChannelsInternal         opus_int32
	API_sampleRate            opus_int32
	MaxInternalSampleRate     opus_int32
	MinInternalSampleRate     opus_int32
	DesiredInternalSampleRate opus_int32
	PayloadSize_ms            int64
	BitRate                   opus_int32
	PacketLossPercentage      int64
	Complexity                int64
	UseInBandFEC              int64
	LBRR_coded                int64
	UseDTX                    int64
	UseCBR                    int64
	MaxBits                   int64
	ToMono                    int64
	OpusCanSwitch             int64
	ReducedDependency         int64
	InternalSampleRate        opus_int32
	AllowBandwidthSwitch      int64
	InWBmodeWithoutVariableLP int64
	StereoWidth_Q14           int64
	SwitchReady               int64
	SignalType                int64
	Offset                    int64
}
type silk_DecControlStruct struct {
	NChannelsAPI       opus_int32
	NChannelsInternal  opus_int32
	API_sampleRate     opus_int32
	InternalSampleRate opus_int32
	PayloadSize_ms     int64
	PrevPitchLag       int64
}
