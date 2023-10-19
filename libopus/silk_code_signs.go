package libopus

import "unsafe"

func silk_encode_signs(psRangeEnc *ec_enc, pulses [0]int8, length int64, signalType int64, quantOffsetType int64, sum_pulses [20]int64) {
	var (
		i        int64
		j        int64
		p        int64
		icdf     [2]uint8
		q_ptr    *int8
		icdf_ptr *uint8
	)
	icdf[1] = 0
	q_ptr = &pulses[0]
	i = int64(opus_int32(opus_int16(quantOffsetType+int64(opus_int32(opus_uint32(signalType)<<1)))) * 7)
	icdf_ptr = &silk_sign_iCDF[i]
	length = (length + SHELL_CODEC_FRAME_LENGTH/2) >> LOG2_SHELL_CODEC_FRAME_LENGTH
	for i = 0; i < length; i++ {
		p = sum_pulses[i]
		if p > 0 {
			icdf[0] = *(*uint8)(unsafe.Add(unsafe.Pointer(icdf_ptr), func() int64 {
				if (p & 0x1F) < 6 {
					return p & 0x1F
				}
				return 6
			}()))
			for j = 0; j < SHELL_CODEC_FRAME_LENGTH; j++ {
				if int64(*(*int8)(unsafe.Add(unsafe.Pointer(q_ptr), j))) != 0 {
					ec_enc_icdf(psRangeEnc, (int64(*(*int8)(unsafe.Add(unsafe.Pointer(q_ptr), j)))>>15)+1, &icdf[0], 8)
				}
			}
		}
		q_ptr = (*int8)(unsafe.Add(unsafe.Pointer(q_ptr), SHELL_CODEC_FRAME_LENGTH))
	}
}
func silk_decode_signs(psRangeDec *ec_dec, pulses [0]opus_int16, length int64, signalType int64, quantOffsetType int64, sum_pulses [20]int64) {
	var (
		i        int64
		j        int64
		p        int64
		icdf     [2]uint8
		q_ptr    *opus_int16
		icdf_ptr *uint8
	)
	icdf[1] = 0
	q_ptr = &pulses[0]
	i = int64(opus_int32(opus_int16(quantOffsetType+int64(opus_int32(opus_uint32(signalType)<<1)))) * 7)
	icdf_ptr = &silk_sign_iCDF[i]
	length = (length + SHELL_CODEC_FRAME_LENGTH/2) >> LOG2_SHELL_CODEC_FRAME_LENGTH
	for i = 0; i < length; i++ {
		p = sum_pulses[i]
		if p > 0 {
			icdf[0] = *(*uint8)(unsafe.Add(unsafe.Pointer(icdf_ptr), func() int64 {
				if (p & 0x1F) < 6 {
					return p & 0x1F
				}
				return 6
			}()))
			for j = 0; j < SHELL_CODEC_FRAME_LENGTH; j++ {
				if *(*opus_int16)(unsafe.Add(unsafe.Pointer(q_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(j))) > 0 {
					*(*opus_int16)(unsafe.Add(unsafe.Pointer(q_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(j))) *= opus_int16((opus_int32(opus_uint32(ec_dec_icdf(psRangeDec, &icdf[0], 8)) << 1)) - 1)
				}
			}
		}
		q_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(q_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(SHELL_CODEC_FRAME_LENGTH)))
	}
}
