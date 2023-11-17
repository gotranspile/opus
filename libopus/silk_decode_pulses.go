package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_decode_pulses(psRangeDec *ec_dec, pulses []int16, signalType int, quantOffsetType int, frame_length int) {
	var (
		i              int
		j              int
		k              int
		iter           int
		abs_q          int
		nLS            int
		RateLevelIndex int
		sum_pulses     [20]int
		nLshifts       [20]int
		pulses_ptr     *int16
		cdf_ptr        *uint8
	)
	RateLevelIndex = ec_dec_icdf(psRangeDec, silk_rate_levels_iCDF[signalType>>1][:], 8)
	iter = frame_length >> LOG2_SHELL_CODEC_FRAME_LENGTH
	if iter*SHELL_CODEC_FRAME_LENGTH < frame_length {
		iter++
	}
	cdf_ptr = &silk_pulses_per_block_iCDF[RateLevelIndex][0]
	for i = 0; i < iter; i++ {
		nLshifts[i] = 0
		sum_pulses[i] = ec_dec_icdf(psRangeDec, []byte(cdf_ptr), 8)
		for sum_pulses[i] == int(SILK_MAX_PULSES+1) {
			nLshifts[i]++
			sum_pulses[i] = ec_dec_icdf(psRangeDec, []byte(&silk_pulses_per_block_iCDF[int(N_RATE_LEVELS-1)][nLshifts[i] == 10]), 8)
		}
	}
	for i = 0; i < iter; i++ {
		if sum_pulses[i] > 0 {
			silk_shell_decoder([]int16(&pulses[int(int32(int16(i)))*SHELL_CODEC_FRAME_LENGTH]), psRangeDec, sum_pulses[i])
		} else {
			libc.MemSet(unsafe.Pointer(&pulses[int(int32(int16(i)))*SHELL_CODEC_FRAME_LENGTH]), 0, int(SHELL_CODEC_FRAME_LENGTH*unsafe.Sizeof(int16(0))))
		}
	}
	for i = 0; i < iter; i++ {
		if nLshifts[i] > 0 {
			nLS = nLshifts[i]
			pulses_ptr = &pulses[int(int32(int16(i)))*SHELL_CODEC_FRAME_LENGTH]
			for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
				abs_q = int(*(*int16)(unsafe.Add(unsafe.Pointer(pulses_ptr), unsafe.Sizeof(int16(0))*uintptr(k))))
				for j = 0; j < nLS; j++ {
					abs_q = int(int32(int(uint32(int32(abs_q))) << 1))
					abs_q += ec_dec_icdf(psRangeDec, silk_lsb_iCDF[:], 8)
				}
				*(*int16)(unsafe.Add(unsafe.Pointer(pulses_ptr), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(abs_q)
			}
			sum_pulses[i] |= nLS << 5
		}
	}
	silk_decode_signs(psRangeDec, pulses, frame_length, signalType, quantOffsetType, sum_pulses)
}
