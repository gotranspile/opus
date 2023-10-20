package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func combine_and_check(pulses_comb *int, pulses_in *int, max_pulses int, len_ int) int {
	var (
		k   int
		sum int
	)
	for k = 0; k < len_; k++ {
		sum = *(*int)(unsafe.Add(unsafe.Pointer(pulses_in), unsafe.Sizeof(int(0))*uintptr(k*2))) + *(*int)(unsafe.Add(unsafe.Pointer(pulses_in), unsafe.Sizeof(int(0))*uintptr(k*2+1)))
		if sum > max_pulses {
			return 1
		}
		*(*int)(unsafe.Add(unsafe.Pointer(pulses_comb), unsafe.Sizeof(int(0))*uintptr(k))) = sum
	}
	return 0
}
func silk_encode_pulses(psRangeEnc *ec_enc, signalType int, quantOffsetType int, pulses []int8, frame_length int) {
	var (
		i              int
		k              int
		j              int
		iter           int
		bit            int
		nLS            int
		scale_down     int
		RateLevelIndex int = 0
		abs_q          int32
		minSumBits_Q5  int32
		sumBits_Q5     int32
		abs_pulses     *int
		sum_pulses     *int
		nRshifts       *int
		pulses_comb    [8]int
		abs_pulses_ptr *int
		pulses_ptr     *int8
		cdf_ptr        *uint8
		nBits_ptr      *uint8
	)
	libc.MemSet(unsafe.Pointer(&pulses_comb[0]), 0, int(8*unsafe.Sizeof(int(0))))
	iter = frame_length >> LOG2_SHELL_CODEC_FRAME_LENGTH
	if iter*SHELL_CODEC_FRAME_LENGTH < frame_length {
		iter++
		libc.MemSet(unsafe.Pointer(&pulses[frame_length]), 0, int(SHELL_CODEC_FRAME_LENGTH*unsafe.Sizeof(int8(0))))
	}
	abs_pulses = (*int)(libc.Malloc((iter * SHELL_CODEC_FRAME_LENGTH) * int(unsafe.Sizeof(int(0)))))
	for i = 0; i < iter*SHELL_CODEC_FRAME_LENGTH; i += 4 {
		if int(pulses[i+0]) > 0 {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+0))) = int(pulses[i+0])
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+0))) = int(-(pulses[i+0]))
		}
		if int(pulses[i+1]) > 0 {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+1))) = int(pulses[i+1])
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+1))) = int(-(pulses[i+1]))
		}
		if int(pulses[i+2]) > 0 {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+2))) = int(pulses[i+2])
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+2))) = int(-(pulses[i+2]))
		}
		if int(pulses[i+3]) > 0 {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+3))) = int(pulses[i+3])
		} else {
			*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i+3))) = int(-(pulses[i+3]))
		}
	}
	sum_pulses = (*int)(libc.Malloc(iter * int(unsafe.Sizeof(int(0)))))
	nRshifts = (*int)(libc.Malloc(iter * int(unsafe.Sizeof(int(0)))))
	abs_pulses_ptr = abs_pulses
	for i = 0; i < iter; i++ {
		*(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i))) = 0
		for {
			scale_down = combine_and_check(&pulses_comb[0], abs_pulses_ptr, int(silk_max_pulses_table[0]), 8)
			scale_down += combine_and_check(&pulses_comb[0], &pulses_comb[0], int(silk_max_pulses_table[1]), 4)
			scale_down += combine_and_check(&pulses_comb[0], &pulses_comb[0], int(silk_max_pulses_table[2]), 2)
			scale_down += combine_and_check((*int)(unsafe.Add(unsafe.Pointer(sum_pulses), unsafe.Sizeof(int(0))*uintptr(i))), &pulses_comb[0], int(silk_max_pulses_table[3]), 1)
			if scale_down != 0 {
				*(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i)))++
				for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
					*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses_ptr), unsafe.Sizeof(int(0))*uintptr(k))) = (*(*int)(unsafe.Add(unsafe.Pointer(abs_pulses_ptr), unsafe.Sizeof(int(0))*uintptr(k)))) >> 1
				}
			} else {
				break
			}
		}
		abs_pulses_ptr = (*int)(unsafe.Add(unsafe.Pointer(abs_pulses_ptr), unsafe.Sizeof(int(0))*uintptr(SHELL_CODEC_FRAME_LENGTH)))
	}
	minSumBits_Q5 = silk_int32_MAX
	for k = 0; k < int(N_RATE_LEVELS-1); k++ {
		nBits_ptr = &silk_pulses_per_block_BITS_Q5[k][0]
		sumBits_Q5 = int32(silk_rate_levels_BITS_Q5[signalType>>1][k])
		for i = 0; i < iter; i++ {
			if *(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i))) > 0 {
				sumBits_Q5 += int32(*(*uint8)(unsafe.Add(unsafe.Pointer(nBits_ptr), int(SILK_MAX_PULSES+1))))
			} else {
				sumBits_Q5 += int32(*(*uint8)(unsafe.Add(unsafe.Pointer(nBits_ptr), *(*int)(unsafe.Add(unsafe.Pointer(sum_pulses), unsafe.Sizeof(int(0))*uintptr(i))))))
			}
		}
		if int(sumBits_Q5) < int(minSumBits_Q5) {
			minSumBits_Q5 = sumBits_Q5
			RateLevelIndex = k
		}
	}
	ec_enc_icdf(psRangeEnc, RateLevelIndex, silk_rate_levels_iCDF[signalType>>1][:], 8)
	cdf_ptr = &silk_pulses_per_block_iCDF[RateLevelIndex][0]
	for i = 0; i < iter; i++ {
		if *(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i))) == 0 {
			ec_enc_icdf(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(sum_pulses), unsafe.Sizeof(int(0))*uintptr(i))), []byte(cdf_ptr), 8)
		} else {
			ec_enc_icdf(psRangeEnc, int(SILK_MAX_PULSES+1), []byte(cdf_ptr), 8)
			for k = 0; k < *(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i)))-1; k++ {
				ec_enc_icdf(psRangeEnc, int(SILK_MAX_PULSES+1), silk_pulses_per_block_iCDF[int(N_RATE_LEVELS-1)][:], 8)
			}
			ec_enc_icdf(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(sum_pulses), unsafe.Sizeof(int(0))*uintptr(i))), silk_pulses_per_block_iCDF[int(N_RATE_LEVELS-1)][:], 8)
		}
	}
	for i = 0; i < iter; i++ {
		if *(*int)(unsafe.Add(unsafe.Pointer(sum_pulses), unsafe.Sizeof(int(0))*uintptr(i))) > 0 {
			silk_shell_encoder(psRangeEnc, (*int)(unsafe.Add(unsafe.Pointer(abs_pulses), unsafe.Sizeof(int(0))*uintptr(i*SHELL_CODEC_FRAME_LENGTH))))
		}
	}
	for i = 0; i < iter; i++ {
		if *(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i))) > 0 {
			pulses_ptr = &pulses[i*SHELL_CODEC_FRAME_LENGTH]
			nLS = *(*int)(unsafe.Add(unsafe.Pointer(nRshifts), unsafe.Sizeof(int(0))*uintptr(i))) - 1
			for k = 0; k < SHELL_CODEC_FRAME_LENGTH; k++ {
				if int(*(*int8)(unsafe.Add(unsafe.Pointer(pulses_ptr), k))) > 0 {
					abs_q = int32(*(*int8)(unsafe.Add(unsafe.Pointer(pulses_ptr), k)))
				} else {
					abs_q = int32(-(*(*int8)(unsafe.Add(unsafe.Pointer(pulses_ptr), k))))
				}
				for j = nLS; j > 0; j-- {
					bit = (int(abs_q) >> j) & 1
					ec_enc_icdf(psRangeEnc, bit, silk_lsb_iCDF[:], 8)
				}
				bit = int(abs_q) & 1
				ec_enc_icdf(psRangeEnc, bit, silk_lsb_iCDF[:], 8)
			}
		}
	}
	silk_encode_signs(psRangeEnc, pulses, frame_length, signalType, quantOffsetType, [20]int(sum_pulses))
}
