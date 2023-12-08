package libopus

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_encode_do_VAD_FLP(psEnc *silk_encoder_state_FLP, activity int) {
	var activity_threshold int = int(int32(math.Floor(SPEECH_ACTIVITY_DTX_THRES*(1<<8) + 0.5)))
	psEnc.SCmn.Arch
	silk_VAD_GetSA_Q8_c(&psEnc.SCmn, []int16(&psEnc.SCmn.InputBuf[1]))
	if activity == VAD_NO_ACTIVITY && psEnc.SCmn.Speech_activity_Q8 >= activity_threshold {
		psEnc.SCmn.Speech_activity_Q8 = activity_threshold - 1
	}
	if psEnc.SCmn.Speech_activity_Q8 < activity_threshold {
		psEnc.SCmn.Indices.SignalType = TYPE_NO_VOICE_ACTIVITY
		psEnc.SCmn.NoSpeechCounter++
		if psEnc.SCmn.NoSpeechCounter <= NB_SPEECH_FRAMES_BEFORE_DTX {
			psEnc.SCmn.InDTX = 0
		} else if psEnc.SCmn.NoSpeechCounter > int(MAX_CONSECUTIVE_DTX+NB_SPEECH_FRAMES_BEFORE_DTX) {
			psEnc.SCmn.NoSpeechCounter = NB_SPEECH_FRAMES_BEFORE_DTX
			psEnc.SCmn.InDTX = 0
		}
		psEnc.SCmn.VAD_flags[psEnc.SCmn.NFramesEncoded] = 0
	} else {
		psEnc.SCmn.NoSpeechCounter = 0
		psEnc.SCmn.InDTX = 0
		psEnc.SCmn.Indices.SignalType = TYPE_UNVOICED
		psEnc.SCmn.VAD_flags[psEnc.SCmn.NFramesEncoded] = 1
	}
}
func silk_encode_frame_FLP(psEnc *silk_encoder_state_FLP, pnBytesOut *int32, psRangeEnc *ec_enc, condCoding int, maxBits int, useCBR int) int {
	var (
		sEncCtrl               silk_encoder_control_FLP
		i                      int
		iter                   int
		maxIter                int
		found_upper            int
		found_lower            int
		ret                    int = 0
		x_frame                *float32
		res_pitch_frame        *float32
		res_pitch              [672]float32
		sRangeEnc_copy         ec_enc
		sRangeEnc_copy2        ec_enc
		sNSQ_copy              silk_nsq_state
		sNSQ_copy2             silk_nsq_state
		seed_copy              int32
		nBits                  int32
		nBits_lower            int32
		nBits_upper            int32
		gainMult_lower         int32
		gainMult_upper         int32
		gainsID                int32
		gainsID_lower          int32
		gainsID_upper          int32
		gainMult_Q8            int16
		ec_prevLagIndex_copy   int16
		ec_prevSignalType_copy int
		LastGainIndex_copy2    int8
		pGains_Q16             [4]int32
		ec_buf_copy            [1275]uint8
		gain_lock              [4]int = [4]int{}
		best_gain_mult         [4]int16
		best_sum               [4]int
	)
	LastGainIndex_copy2 = int8(func() int32 {
		nBits_lower = func() int32 {
			nBits_upper = func() int32 {
				gainMult_lower = func() int32 {
					gainMult_upper = 0
					return gainMult_upper
				}()
				return gainMult_lower
			}()
			return nBits_upper
		}()
		return nBits_lower
	}())
	psEnc.SCmn.Indices.Seed = int8(int(func() int32 {
		p := &psEnc.SCmn.FrameCounter
		x := *p
		*p++
		return x
	}()) & 3)
	x_frame = &psEnc.X_buf[psEnc.SCmn.Ltp_mem_length]
	res_pitch_frame = &res_pitch[psEnc.SCmn.Ltp_mem_length]
	silk_LP_variable_cutoff(&psEnc.SCmn.SLP, []int16(&psEnc.SCmn.InputBuf[1]), psEnc.SCmn.Frame_length)
	silk_short2float_array([]float32((*float32)(unsafe.Add(unsafe.Pointer(x_frame), unsafe.Sizeof(float32(0))*uintptr(LA_SHAPE_MS*psEnc.SCmn.Fs_kHz)))), []int16(&psEnc.SCmn.InputBuf[1]), int32(psEnc.SCmn.Frame_length))
	for i = 0; i < 8; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(x_frame), unsafe.Sizeof(float32(0))*uintptr(LA_SHAPE_MS*psEnc.SCmn.Fs_kHz+i*(psEnc.SCmn.Frame_length>>3)))) += float32(float64(1-(i&2)) * 1e-06)
	}
	if psEnc.SCmn.PrefillFlag == 0 {
		silk_find_pitch_lags_FLP(psEnc, &sEncCtrl, res_pitch[:], []float32(x_frame), psEnc.SCmn.Arch)
		silk_noise_shape_analysis_FLP(psEnc, &sEncCtrl, res_pitch_frame, x_frame)
		silk_find_pred_coefs_FLP(psEnc, &sEncCtrl, []float32(res_pitch_frame), []float32(x_frame), condCoding)
		silk_process_gains_FLP(psEnc, &sEncCtrl, condCoding)
		silk_LBRR_encode_FLP(psEnc, &sEncCtrl, []float32(x_frame), condCoding)
		maxIter = 6
		gainMult_Q8 = int16(int32(math.Floor(1*(1<<8) + 0.5)))
		found_lower = 0
		found_upper = 0
		gainsID = silk_gains_ID(psEnc.SCmn.Indices.GainsIndices, psEnc.SCmn.Nb_subfr)
		gainsID_lower = -1
		gainsID_upper = -1
		libc.MemCpy(unsafe.Pointer(&sRangeEnc_copy), unsafe.Pointer(psRangeEnc), int(unsafe.Sizeof(ec_enc{})))
		libc.MemCpy(unsafe.Pointer(&sNSQ_copy), unsafe.Pointer(&psEnc.SCmn.SNSQ), int(unsafe.Sizeof(silk_nsq_state{})))
		seed_copy = int32(psEnc.SCmn.Indices.Seed)
		ec_prevLagIndex_copy = psEnc.SCmn.Ec_prevLagIndex
		ec_prevSignalType_copy = psEnc.SCmn.Ec_prevSignalType
		for iter = 0; ; iter++ {
			if int(gainsID) == int(gainsID_lower) {
				nBits = nBits_lower
			} else if int(gainsID) == int(gainsID_upper) {
				nBits = nBits_upper
			} else {
				if iter > 0 {
					libc.MemCpy(unsafe.Pointer(psRangeEnc), unsafe.Pointer(&sRangeEnc_copy), int(unsafe.Sizeof(ec_enc{})))
					libc.MemCpy(unsafe.Pointer(&psEnc.SCmn.SNSQ), unsafe.Pointer(&sNSQ_copy), int(unsafe.Sizeof(silk_nsq_state{})))
					psEnc.SCmn.Indices.Seed = int8(seed_copy)
					psEnc.SCmn.Ec_prevLagIndex = ec_prevLagIndex_copy
					psEnc.SCmn.Ec_prevSignalType = ec_prevSignalType_copy
				}
				silk_NSQ_wrapper_FLP(psEnc, &sEncCtrl, &psEnc.SCmn.Indices, &psEnc.SCmn.SNSQ, psEnc.SCmn.Pulses[:], []float32(x_frame))
				if iter == maxIter && found_lower == 0 {
					libc.MemCpy(unsafe.Pointer(&sRangeEnc_copy2), unsafe.Pointer(psRangeEnc), int(unsafe.Sizeof(ec_enc{})))
				}
				silk_encode_indices(&psEnc.SCmn, psRangeEnc, psEnc.SCmn.NFramesEncoded, 0, condCoding)
				silk_encode_pulses(psRangeEnc, int(psEnc.SCmn.Indices.SignalType), int(psEnc.SCmn.Indices.QuantOffsetType), psEnc.SCmn.Pulses[:], psEnc.SCmn.Frame_length)
				nBits = int32(ec_tell((*ec_ctx)(unsafe.Pointer(psRangeEnc))))
				if iter == maxIter && found_lower == 0 && int(nBits) > maxBits {
					libc.MemCpy(unsafe.Pointer(psRangeEnc), unsafe.Pointer(&sRangeEnc_copy2), int(unsafe.Sizeof(ec_enc{})))
					psEnc.SShape.LastGainIndex = sEncCtrl.LastGainIndexPrev
					for i = 0; i < psEnc.SCmn.Nb_subfr; i++ {
						psEnc.SCmn.Indices.GainsIndices[i] = 4
					}
					if condCoding != CODE_CONDITIONALLY {
						psEnc.SCmn.Indices.GainsIndices[0] = sEncCtrl.LastGainIndexPrev
					}
					psEnc.SCmn.Ec_prevLagIndex = ec_prevLagIndex_copy
					psEnc.SCmn.Ec_prevSignalType = ec_prevSignalType_copy
					for i = 0; i < psEnc.SCmn.Frame_length; i++ {
						psEnc.SCmn.Pulses[i] = 0
					}
					silk_encode_indices(&psEnc.SCmn, psRangeEnc, psEnc.SCmn.NFramesEncoded, 0, condCoding)
					silk_encode_pulses(psRangeEnc, int(psEnc.SCmn.Indices.SignalType), int(psEnc.SCmn.Indices.QuantOffsetType), psEnc.SCmn.Pulses[:], psEnc.SCmn.Frame_length)
					nBits = int32(ec_tell((*ec_ctx)(unsafe.Pointer(psRangeEnc))))
				}
				if useCBR == 0 && iter == 0 && int(nBits) <= maxBits {
					break
				}
			}
			if iter == maxIter {
				if found_lower != 0 && (int(gainsID) == int(gainsID_lower) || int(nBits) > maxBits) {
					libc.MemCpy(unsafe.Pointer(psRangeEnc), unsafe.Pointer(&sRangeEnc_copy2), int(unsafe.Sizeof(ec_enc{})))
					libc.MemCpy(unsafe.Pointer(&psRangeEnc.Buf[0]), unsafe.Pointer(&ec_buf_copy[0]), int(sRangeEnc_copy2.Offs))
					libc.MemCpy(unsafe.Pointer(&psEnc.SCmn.SNSQ), unsafe.Pointer(&sNSQ_copy2), int(unsafe.Sizeof(silk_nsq_state{})))
					psEnc.SShape.LastGainIndex = LastGainIndex_copy2
				}
				break
			}
			if int(nBits) > maxBits {
				if found_lower == 0 && iter >= 2 {
					if (sEncCtrl.Lambda * 1.5) > 1.5 {
						sEncCtrl.Lambda = sEncCtrl.Lambda * 1.5
					} else {
						sEncCtrl.Lambda = 1.5
					}
					psEnc.SCmn.Indices.QuantOffsetType = 0
					found_upper = 0
					gainsID_upper = -1
				} else {
					found_upper = 1
					nBits_upper = nBits
					gainMult_upper = int32(gainMult_Q8)
					gainsID_upper = gainsID
				}
			} else if int(nBits) < maxBits-5 {
				found_lower = 1
				nBits_lower = nBits
				gainMult_lower = int32(gainMult_Q8)
				if int(gainsID) != int(gainsID_lower) {
					gainsID_lower = gainsID
					libc.MemCpy(unsafe.Pointer(&sRangeEnc_copy2), unsafe.Pointer(psRangeEnc), int(unsafe.Sizeof(ec_enc{})))
					libc.MemCpy(unsafe.Pointer(&ec_buf_copy[0]), unsafe.Pointer(&psRangeEnc.Buf[0]), int(psRangeEnc.Offs))
					libc.MemCpy(unsafe.Pointer(&sNSQ_copy2), unsafe.Pointer(&psEnc.SCmn.SNSQ), int(unsafe.Sizeof(silk_nsq_state{})))
					LastGainIndex_copy2 = psEnc.SShape.LastGainIndex
				}
			} else {
				break
			}
			if found_lower == 0 && int(nBits) > maxBits {
				var j int
				for i = 0; i < psEnc.SCmn.Nb_subfr; i++ {
					var sum int = 0
					for j = i * psEnc.SCmn.Subfr_length; j < (i+1)*psEnc.SCmn.Subfr_length; j++ {
						sum += int(cmath.Abs(int64(psEnc.SCmn.Pulses[j])))
					}
					if iter == 0 || sum < best_sum[i] && gain_lock[i] == 0 {
						best_sum[i] = sum
						best_gain_mult[i] = gainMult_Q8
					} else {
						gain_lock[i] = 1
					}
				}
			}
			if (found_lower & found_upper) == 0 {
				if int(nBits) > maxBits {
					if int(gainMult_Q8) < 16384 {
						gainMult_Q8 *= 2
					} else {
						gainMult_Q8 = math.MaxInt16
					}
				} else {
					var gain_factor_Q16 int32
					gain_factor_Q16 = silk_log2lin(int32(int(int32(int(uint32(int32(int(nBits)-maxBits)))<<7))/psEnc.SCmn.Frame_length + int(int32(math.Floor(16*(1<<7)+0.5)))))
					gainMult_Q8 = int16(int32((int64(gain_factor_Q16) * int64(gainMult_Q8)) >> 16))
				}
			} else {
				gainMult_Q8 = int16(int(gainMult_lower) + ((int(gainMult_upper)-int(gainMult_lower))*(maxBits-int(nBits_lower)))/(int(nBits_upper)-int(nBits_lower)))
				if int(gainMult_Q8) > (int(gainMult_lower) + ((int(gainMult_upper) - int(gainMult_lower)) >> 2)) {
					gainMult_Q8 = int16(int(gainMult_lower) + ((int(gainMult_upper) - int(gainMult_lower)) >> 2))
				} else if int(gainMult_Q8) < (int(gainMult_upper) - ((int(gainMult_upper) - int(gainMult_lower)) >> 2)) {
					gainMult_Q8 = int16(int(gainMult_upper) - ((int(gainMult_upper) - int(gainMult_lower)) >> 2))
				}
			}
			for i = 0; i < psEnc.SCmn.Nb_subfr; i++ {
				var tmp int16
				if gain_lock[i] != 0 {
					tmp = best_gain_mult[i]
				} else {
					tmp = gainMult_Q8
				}
				pGains_Q16[i] = int32(int(uint32(int32(func() int {
					if (int(math.MinInt32) >> 8) > (int(silk_int32_MAX >> 8)) {
						if int(int32((int64(sEncCtrl.GainsUnq_Q16[i])*int64(tmp))>>16)) > (int(math.MinInt32) >> 8) {
							return int(math.MinInt32) >> 8
						}
						if int(int32((int64(sEncCtrl.GainsUnq_Q16[i])*int64(tmp))>>16)) < (int(silk_int32_MAX >> 8)) {
							return int(silk_int32_MAX >> 8)
						}
						return int(int32((int64(sEncCtrl.GainsUnq_Q16[i]) * int64(tmp)) >> 16))
					}
					if int(int32((int64(sEncCtrl.GainsUnq_Q16[i])*int64(tmp))>>16)) > (int(silk_int32_MAX >> 8)) {
						return int(silk_int32_MAX >> 8)
					}
					if int(int32((int64(sEncCtrl.GainsUnq_Q16[i])*int64(tmp))>>16)) < (int(math.MinInt32) >> 8) {
						return int(math.MinInt32) >> 8
					}
					return int(int32((int64(sEncCtrl.GainsUnq_Q16[i]) * int64(tmp)) >> 16))
				}()))) << 8)
			}
			psEnc.SShape.LastGainIndex = sEncCtrl.LastGainIndexPrev
			silk_gains_quant(psEnc.SCmn.Indices.GainsIndices, pGains_Q16, &psEnc.SShape.LastGainIndex, int(libc.BoolToInt(condCoding == CODE_CONDITIONALLY)), psEnc.SCmn.Nb_subfr)
			gainsID = silk_gains_ID(psEnc.SCmn.Indices.GainsIndices, psEnc.SCmn.Nb_subfr)
			for i = 0; i < psEnc.SCmn.Nb_subfr; i++ {
				sEncCtrl.Gains[i] = float32(float64(pGains_Q16[i]) / 65536.0)
			}
		}
	}
	libc.MemMove(unsafe.Pointer(&psEnc.X_buf[0]), unsafe.Pointer(&psEnc.X_buf[psEnc.SCmn.Frame_length]), (psEnc.SCmn.Ltp_mem_length+LA_SHAPE_MS*psEnc.SCmn.Fs_kHz)*int(unsafe.Sizeof(float32(0))))
	if psEnc.SCmn.PrefillFlag != 0 {
		*pnBytesOut = 0
		return ret
	}
	psEnc.SCmn.PrevLag = sEncCtrl.PitchL[psEnc.SCmn.Nb_subfr-1]
	psEnc.SCmn.PrevSignalType = psEnc.SCmn.Indices.SignalType
	psEnc.SCmn.First_frame_after_reset = 0
	*pnBytesOut = int32((ec_tell((*ec_ctx)(unsafe.Pointer(psRangeEnc))) + 7) >> 3)
	return ret
}
func silk_LBRR_encode_FLP(psEnc *silk_encoder_state_FLP, psEncCtrl *silk_encoder_control_FLP, xfw []float32, condCoding int) {
	var (
		k              int
		Gains_Q16      [4]int32
		TempGains      [4]float32
		psIndices_LBRR *SideInfoIndices = &psEnc.SCmn.Indices_LBRR[psEnc.SCmn.NFramesEncoded]
		sNSQ_LBRR      silk_nsq_state
	)
	if psEnc.SCmn.LBRR_enabled != 0 && psEnc.SCmn.Speech_activity_Q8 > int(int32(math.Floor(LBRR_SPEECH_ACTIVITY_THRES*(1<<8)+0.5))) {
		psEnc.SCmn.LBRR_flags[psEnc.SCmn.NFramesEncoded] = 1
		libc.MemCpy(unsafe.Pointer(&sNSQ_LBRR), unsafe.Pointer(&psEnc.SCmn.SNSQ), int(unsafe.Sizeof(silk_nsq_state{})))
		libc.MemCpy(unsafe.Pointer(psIndices_LBRR), unsafe.Pointer(&psEnc.SCmn.Indices), int(unsafe.Sizeof(SideInfoIndices{})))
		libc.MemCpy(unsafe.Pointer(&TempGains[0]), unsafe.Pointer(&psEncCtrl.Gains[0]), psEnc.SCmn.Nb_subfr*int(unsafe.Sizeof(float32(0))))
		if psEnc.SCmn.NFramesEncoded == 0 || psEnc.SCmn.LBRR_flags[psEnc.SCmn.NFramesEncoded-1] == 0 {
			psEnc.SCmn.LBRRprevLastGainIndex = psEnc.SShape.LastGainIndex
			psIndices_LBRR.GainsIndices[0] += int8(psEnc.SCmn.LBRR_GainIncreases)
			psIndices_LBRR.GainsIndices[0] = int8(silk_min_int(int(psIndices_LBRR.GainsIndices[0]), int(N_LEVELS_QGAIN-1)))
		}
		silk_gains_dequant(Gains_Q16, psIndices_LBRR.GainsIndices, &psEnc.SCmn.LBRRprevLastGainIndex, int(libc.BoolToInt(condCoding == CODE_CONDITIONALLY)), psEnc.SCmn.Nb_subfr)
		for k = 0; k < psEnc.SCmn.Nb_subfr; k++ {
			psEncCtrl.Gains[k] = float32(float64(Gains_Q16[k]) * (1.0 / 65536.0))
		}
		silk_NSQ_wrapper_FLP(psEnc, psEncCtrl, psIndices_LBRR, &sNSQ_LBRR, psEnc.SCmn.Pulses_LBRR[psEnc.SCmn.NFramesEncoded][:], xfw)
		libc.MemCpy(unsafe.Pointer(&psEncCtrl.Gains[0]), unsafe.Pointer(&TempGains[0]), psEnc.SCmn.Nb_subfr*int(unsafe.Sizeof(float32(0))))
	}
}
