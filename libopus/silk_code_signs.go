package libopus

import "unsafe"

func silk_encode_signs(psRangeEnc *ec_enc, pulses []int8, length int, signalType int, quantOffsetType int, sum_pulses [20]int) {
	var (
		i        int
		j        int
		p        int
		icdf     [2]uint8
		q_ptr    *int8
		icdf_ptr *uint8
	)
	icdf[1] = 0
	q_ptr = &pulses[0]
	i = int(int32(int16(quantOffsetType+int(int32(int(uint32(int32(signalType)))<<1))))) * 7
	icdf_ptr = &silk_sign_iCDF[i]
	length = (length + int(SHELL_CODEC_FRAME_LENGTH/2)) >> LOG2_SHELL_CODEC_FRAME_LENGTH
	for i = 0; i < length; i++ {
		p = sum_pulses[i]
		if p > 0 {
			icdf[0] = *(*uint8)(unsafe.Add(unsafe.Pointer(icdf_ptr), func() int {
				if (p & 0x1F) < 6 {
					return p & 0x1F
				}
				return 6
			}()))
			for j = 0; j < SHELL_CODEC_FRAME_LENGTH; j++ {
				if int(*(*int8)(unsafe.Add(unsafe.Pointer(q_ptr), j))) != 0 {
					ec_enc_icdf(psRangeEnc, (int(*(*int8)(unsafe.Add(unsafe.Pointer(q_ptr), j)))>>15)+1, icdf[:], 8)
				}
			}
		}
		q_ptr = (*int8)(unsafe.Add(unsafe.Pointer(q_ptr), SHELL_CODEC_FRAME_LENGTH))
	}
}
func silk_decode_signs(psRangeDec *ec_dec, pulses []int16, length int, signalType int, quantOffsetType int, sum_pulses [20]int) {
	var (
		i        int
		j        int
		p        int
		icdf     [2]uint8
		q_ptr    *int16
		icdf_ptr *uint8
	)
	icdf[1] = 0
	q_ptr = &pulses[0]
	i = int(int32(int16(quantOffsetType+int(int32(int(uint32(int32(signalType)))<<1))))) * 7
	icdf_ptr = &silk_sign_iCDF[i]
	length = (length + int(SHELL_CODEC_FRAME_LENGTH/2)) >> LOG2_SHELL_CODEC_FRAME_LENGTH
	for i = 0; i < length; i++ {
		p = sum_pulses[i]
		if p > 0 {
			icdf[0] = *(*uint8)(unsafe.Add(unsafe.Pointer(icdf_ptr), func() int {
				if (p & 0x1F) < 6 {
					return p & 0x1F
				}
				return 6
			}()))
			for j = 0; j < SHELL_CODEC_FRAME_LENGTH; j++ {
				if int(*(*int16)(unsafe.Add(unsafe.Pointer(q_ptr), unsafe.Sizeof(int16(0))*uintptr(j)))) > 0 {
					*(*int16)(unsafe.Add(unsafe.Pointer(q_ptr), unsafe.Sizeof(int16(0))*uintptr(j))) *= int16(int(int32(int(uint32(int32(ec_dec_icdf(psRangeDec, icdf[:], 8))))<<1)) - 1)
				}
			}
		}
		q_ptr = (*int16)(unsafe.Add(unsafe.Pointer(q_ptr), unsafe.Sizeof(int16(0))*uintptr(SHELL_CODEC_FRAME_LENGTH)))
	}
}
