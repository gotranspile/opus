package libopus

import (
	"unsafe"

	"github.com/gotranspile/opus/silk"
)

func silk_Get_Decoder_Size(decSizeBytes *int) int {
	return silk.GetDecoderSize(decSizeBytes)
}
func silk_InitDecoder(decState unsafe.Pointer) int {
	return silk.InitDecoder(decState)
}
func silk_Decode(decState unsafe.Pointer, decControl *silk_DecControlStruct, lostFlag int, newPacketFlag int, psRangeDec *ec_dec, samplesOut []int16, nSamplesOut *int32, arch int) int {
	return silk.Decode(decState, decControl, lostFlag, newPacketFlag, psRangeDec, samplesOut, nSamplesOut, arch)
}
