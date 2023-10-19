package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_decode_pulses(psRangeDec *ec_dec, pulses [0]opus_int16, signalType int64, quantOffsetType int64, frame_length int64) {
	var (
		i              int64
		j              int64
		k              int64
		iter           int64
		abs_q          int64
		nLS            int64
		RateLevelIndex int64
		sum_pulses     [20]int64
		nLshifts       [20]int64
		pulses_ptr     *opus_int16
		cdf_ptr        *uint8
	)
	RateLevelIndex = ec_dec_icdf(psRangeDec, &silk_rate_levels_iCDF[signalType>>1][0], 8)
	iter = frame_length >> LOG2_SHELL_CODEC_FRAME_LENGTH
	if iter*SHELL_CODEC_FRAME_LENGTH < frame_length {
		iter++
	}
	cdf_ptr = &silk_pulses_per_block_iCDF[RateLevelIndex][0]
	for i = 0; i < iter; i++ {
		nLshifts[i] = 0
		sum_pulses[i] = ec_dec_icdf(psRangeDec, cdf_ptr, 8)
		for sum_pulses[i] == SILK_MAX_PULSES+1 {
			nLshifts[i]++
			sum_pulses[i] = ec_dec_icdf(psRangeDec, &silk_pulses_per_block_iCDF[N_RATE_LEVELS-1][nLshifts[i] == 10], 8)
		}
	}
	for i = 0; i < iter; i++ {
		if sum_pulses[i] > 0 {
			silk_shell_decoder(&pulses[opus_int32(opus_int16(i))*SHELL_CODEC_FRAME_LENGTH], psRangeDec, sum_pulses[i])
		} else {
			libc.MemSet(unsafe.Pointer(&pulses[opus_int32(opus_int16(i))*SHELL_CODEC_FRAME_LENGTH]), 0, int(SHELL_CODEC_FRAME_LENGTH*unsafe.Sizeof(opus_int16(0))))
		}
	}
	for i = 0; i < iter; i++ {
		if nLshifts[i] > 0 {
			nLS = nLshifts[i]
			pulses_ptr = &pulses[opus_int32(opus_int16(i))*SHELL_CODEC_FRAME_LENGTH]
			for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
				abs_q = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pulses_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(k))))
				for j = 0; j < nLS; j++ {
					abs_q = int64(opus_int32(opus_uint32(abs_q) << 1))
					abs_q += ec_dec_icdf(psRangeDec, &silk_lsb_iCDF[0], 8)
				}
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pulses_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(abs_q)
			}
			sum_pulses[i] |= nLS << 5
		}
	}
	silk_decode_signs(psRangeDec, pulses, frame_length, signalType, quantOffsetType, sum_pulses)
}
