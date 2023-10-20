package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

type silk_decoder struct {
	Channel_state           [2]silk_decoder_state
	SStereo                 stereo_dec_state
	NChannelsAPI            int
	NChannelsInternal       int
	Prev_decode_only_middle int
}

func silk_Get_Decoder_Size(decSizeBytes *int) int {
	var ret int = SILK_NO_ERROR
	*decSizeBytes = int(unsafe.Sizeof(silk_decoder{}))
	return ret
}
func silk_InitDecoder(decState unsafe.Pointer) int {
	var (
		n             int
		ret           int                 = SILK_NO_ERROR
		channel_state *silk_decoder_state = &((*silk_decoder)(decState)).Channel_state[0]
	)
	for n = 0; n < DECODER_NUM_CHANNELS; n++ {
		ret = silk_init_decoder((*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n))))
	}
	((*silk_decoder)(decState)).SStereo = stereo_dec_state{}
	((*silk_decoder)(decState)).Prev_decode_only_middle = 0
	return ret
}
func silk_Decode(decState unsafe.Pointer, decControl *silk_DecControlStruct, lostFlag int, newPacketFlag int, psRangeDec *ec_dec, samplesOut []int16, nSamplesOut *int32, arch int) int {
	var (
		i                        int
		n                        int
		decode_only_middle       int = 0
		ret                      int = SILK_NO_ERROR
		nSamplesOutDec           int32
		LBRR_symbol              int32
		samplesOut1_tmp          [2]*int16
		samplesOut1_tmp_storage1 *int16
		samplesOut1_tmp_storage2 *int16
		samplesOut2_tmp          *int16
		MS_pred_Q13              [2]int32 = [2]int32{}
		resample_out_ptr         *int16
		psDec                    *silk_decoder       = (*silk_decoder)(decState)
		channel_state            *silk_decoder_state = &psDec.Channel_state[0]
		has_side                 int
		stereo_to_mono           int
		delay_stack_alloc        int
	)
	if newPacketFlag != 0 {
		for n = 0; n < int(decControl.NChannelsInternal); n++ {
			(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesDecoded = 0
		}
	}
	if int(decControl.NChannelsInternal) > psDec.NChannelsInternal {
		ret += silk_init_decoder((*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1)))
	}
	stereo_to_mono = int(libc.BoolToInt(int(decControl.NChannelsInternal) == 1 && psDec.NChannelsInternal == 2 && int(decControl.InternalSampleRate) == (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Fs_kHz*1000))
	if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesDecoded == 0 {
		for n = 0; n < int(decControl.NChannelsInternal); n++ {
			var fs_kHz_dec int
			if decControl.PayloadSize_ms == 0 {
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket = 1
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Nb_subfr = 2
			} else if decControl.PayloadSize_ms == 10 {
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket = 1
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Nb_subfr = 2
			} else if decControl.PayloadSize_ms == 20 {
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket = 1
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Nb_subfr = 4
			} else if decControl.PayloadSize_ms == 40 {
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket = 2
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Nb_subfr = 4
			} else if decControl.PayloadSize_ms == 60 {
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket = 3
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Nb_subfr = 4
			} else {
				return -203
			}
			fs_kHz_dec = (int(decControl.InternalSampleRate) >> 10) + 1
			if fs_kHz_dec != 8 && fs_kHz_dec != 12 && fs_kHz_dec != 16 {
				return -200
			}
			ret += silk_decoder_set_fs((*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n))), fs_kHz_dec, decControl.API_sampleRate)
		}
	}
	if int(decControl.NChannelsAPI) == 2 && int(decControl.NChannelsInternal) == 2 && (psDec.NChannelsAPI == 1 || psDec.NChannelsInternal == 1) {
		*(*[2]int16)(unsafe.Pointer(&psDec.SStereo.Pred_prev_Q13[0])) = [2]int16{}
		*(*[2]int16)(unsafe.Pointer(&psDec.SStereo.SSide[0])) = [2]int16{}
		libc.MemCpy(unsafe.Pointer(&(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).Resampler_state), unsafe.Pointer(&(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Resampler_state), int(unsafe.Sizeof(silk_resampler_state_struct{})))
	}
	psDec.NChannelsAPI = int(decControl.NChannelsAPI)
	psDec.NChannelsInternal = int(decControl.NChannelsInternal)
	if int(decControl.API_sampleRate) > int(MAX_API_FS_KHZ*1000) || int(decControl.API_sampleRate) < 8000 {
		ret = -200
		return ret
	}
	if lostFlag != FLAG_PACKET_LOST && (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesDecoded == 0 {
		for n = 0; n < int(decControl.NChannelsInternal); n++ {
			for i = 0; i < (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket; i++ {
				(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).VAD_flags[i] = ec_dec_bit_logp(psRangeDec, 1)
			}
			(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flag = ec_dec_bit_logp(psRangeDec, 1)
		}
		for n = 0; n < int(decControl.NChannelsInternal); n++ {
			*(*[3]int)(unsafe.Pointer(&(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flags[0])) = [3]int{}
			if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flag != 0 {
				if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket == 1 {
					(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flags[0] = 1
				} else {
					LBRR_symbol = int32(ec_dec_icdf(psRangeDec, []byte(silk_LBRR_flags_iCDF_ptr[(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket-2]), 8) + 1)
					for i = 0; i < (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesPerPacket; i++ {
						(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flags[i] = (int(LBRR_symbol) >> i) & 1
					}
				}
			}
		}
		if lostFlag == FLAG_DECODE_NORMAL {
			for i = 0; i < (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesPerPacket; i++ {
				for n = 0; n < int(decControl.NChannelsInternal); n++ {
					if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flags[i] != 0 {
						var (
							pulses     [320]int16
							condCoding int
						)
						if int(decControl.NChannelsInternal) == 2 && n == 0 {
							silk_stereo_decode_pred(psRangeDec, MS_pred_Q13[:])
							if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).LBRR_flags[i] == 0 {
								silk_stereo_decode_mid_only(psRangeDec, &decode_only_middle)
							}
						}
						if i > 0 && (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flags[i-1] != 0 {
							condCoding = CODE_CONDITIONALLY
						} else {
							condCoding = CODE_INDEPENDENTLY
						}
						silk_decode_indices((*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n))), psRangeDec, i, 1, condCoding)
						silk_decode_pulses(psRangeDec, pulses[:], int((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Indices.SignalType), int((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Indices.QuantOffsetType), (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Frame_length)
					}
				}
			}
		}
	}
	if int(decControl.NChannelsInternal) == 2 {
		if lostFlag == FLAG_DECODE_NORMAL || lostFlag == FLAG_DECODE_LBRR && (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).LBRR_flags[(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesDecoded] == 1 {
			silk_stereo_decode_pred(psRangeDec, MS_pred_Q13[:])
			if lostFlag == FLAG_DECODE_NORMAL && (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).VAD_flags[(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesDecoded] == 0 || lostFlag == FLAG_DECODE_LBRR && (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).LBRR_flags[(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesDecoded] == 0 {
				silk_stereo_decode_mid_only(psRangeDec, &decode_only_middle)
			} else {
				decode_only_middle = 0
			}
		} else {
			for n = 0; n < 2; n++ {
				MS_pred_Q13[n] = int32(psDec.SStereo.Pred_prev_Q13[n])
			}
		}
	}
	if int(decControl.NChannelsInternal) == 2 && decode_only_middle == 0 && psDec.Prev_decode_only_middle == 1 {
		*(*[480]int16)(unsafe.Pointer(&psDec.Channel_state[1].OutBuf[0])) = [480]int16{}
		*(*[16]int32)(unsafe.Pointer(&psDec.Channel_state[1].SLPC_Q14_buf[0])) = [16]int32{}
		psDec.Channel_state[1].LagPrev = 100
		psDec.Channel_state[1].LastGainIndex = 10
		psDec.Channel_state[1].PrevSignalType = TYPE_NO_VOICE_ACTIVITY
		psDec.Channel_state[1].First_frame_after_reset = 1
	}
	delay_stack_alloc = int(libc.BoolToInt(int(decControl.InternalSampleRate)*int(decControl.NChannelsInternal) < int(decControl.API_sampleRate)*int(decControl.NChannelsAPI)))
	samplesOut1_tmp_storage1 = (*int16)(libc.Malloc((func() int {
		if delay_stack_alloc != 0 {
			return ALLOC_NONE
		}
		return int(decControl.NChannelsInternal) * ((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Frame_length + 2)
	}()) * int(unsafe.Sizeof(int16(0)))))
	if delay_stack_alloc != 0 {
		samplesOut1_tmp[0] = &samplesOut[0]
		samplesOut1_tmp[1] = &samplesOut[(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Frame_length+2]
	} else {
		samplesOut1_tmp[0] = samplesOut1_tmp_storage1
		samplesOut1_tmp[1] = (*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp_storage1), unsafe.Sizeof(int16(0))*uintptr((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Frame_length)))), unsafe.Sizeof(int16(0))*2))
	}
	if lostFlag == FLAG_DECODE_NORMAL {
		has_side = int(libc.BoolToInt(decode_only_middle == 0))
	} else {
		has_side = int(libc.BoolToInt(psDec.Prev_decode_only_middle == 0 || int(decControl.NChannelsInternal) == 2 && lostFlag == FLAG_DECODE_LBRR && (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).LBRR_flags[(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).NFramesDecoded] == 1))
	}
	for n = 0; n < int(decControl.NChannelsInternal); n++ {
		if n == 0 || has_side != 0 {
			var (
				FrameIndex int
				condCoding int
			)
			FrameIndex = (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).NFramesDecoded - n
			if FrameIndex <= 0 {
				condCoding = CODE_INDEPENDENTLY
			} else if lostFlag == FLAG_DECODE_LBRR {
				if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).LBRR_flags[FrameIndex-1] != 0 {
					condCoding = CODE_CONDITIONALLY
				} else {
					condCoding = CODE_INDEPENDENTLY
				}
			} else if n > 0 && psDec.Prev_decode_only_middle != 0 {
				condCoding = CODE_INDEPENDENTLY_NO_LTP_SCALING
			} else {
				condCoding = CODE_CONDITIONALLY
			}
			ret += silk_decode_frame((*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n))), psRangeDec, []int16((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp[n]), unsafe.Sizeof(int16(0))*2))), &nSamplesOutDec, lostFlag, condCoding, arch)
		} else {
			libc.MemSet(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp[n]), unsafe.Sizeof(int16(0))*2))), 0, int(uintptr(nSamplesOutDec)*unsafe.Sizeof(int16(0))))
		}
		(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).NFramesDecoded++
	}
	if int(decControl.NChannelsAPI) == 2 && int(decControl.NChannelsInternal) == 2 {
		silk_stereo_MS_to_LR(&psDec.SStereo, []int16(samplesOut1_tmp[0]), []int16(samplesOut1_tmp[1]), MS_pred_Q13[:], (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Fs_kHz, int(nSamplesOutDec))
	} else {
		libc.MemCpy(unsafe.Pointer(samplesOut1_tmp[0]), unsafe.Pointer(&psDec.SStereo.SMid[0]), int(2*unsafe.Sizeof(int16(0))))
		libc.MemCpy(unsafe.Pointer(&psDec.SStereo.SMid[0]), unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp[0]), unsafe.Sizeof(int16(0))*uintptr(nSamplesOutDec)))), int(2*unsafe.Sizeof(int16(0))))
	}
	*nSamplesOut = int32((int(nSamplesOutDec) * int(decControl.API_sampleRate)) / (int(int32(int16((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Fs_kHz))) * 1000))
	samplesOut2_tmp = (*int16)(libc.Malloc((func() int {
		if int(decControl.NChannelsAPI) == 2 {
			return int(*nSamplesOut)
		}
		return ALLOC_NONE
	}()) * int(unsafe.Sizeof(int16(0)))))
	if int(decControl.NChannelsAPI) == 2 {
		resample_out_ptr = samplesOut2_tmp
	} else {
		resample_out_ptr = &samplesOut[0]
	}
	samplesOut1_tmp_storage2 = (*int16)(libc.Malloc((func() int {
		if delay_stack_alloc != 0 {
			return int(decControl.NChannelsInternal) * ((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Frame_length + 2)
		}
		return ALLOC_NONE
	}()) * int(unsafe.Sizeof(int16(0)))))
	if delay_stack_alloc != 0 {
		libc.MemCpy(unsafe.Pointer(samplesOut1_tmp_storage2), unsafe.Pointer(&samplesOut[0]), int(uintptr(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer((samplesOut1_tmp_storage2-&samplesOut[0]*(2*2))*(0*2)), unsafe.Sizeof(int16(0))*uintptr((int(decControl.NChannelsInternal)*((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Frame_length+2))*int(unsafe.Sizeof(int16(0))))))))))
		samplesOut1_tmp[0] = samplesOut1_tmp_storage2
		samplesOut1_tmp[1] = (*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp_storage2), unsafe.Sizeof(int16(0))*uintptr((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Frame_length)))), unsafe.Sizeof(int16(0))*2))
	}
	for n = 0; n < (func() int {
		if int(decControl.NChannelsAPI) < int(decControl.NChannelsInternal) {
			return int(decControl.NChannelsAPI)
		}
		return int(decControl.NChannelsInternal)
	}()); n++ {
		ret += silk_resampler(&(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*uintptr(n)))).Resampler_state, []int16(resample_out_ptr), []int16((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp[n]), unsafe.Sizeof(int16(0))*1))), nSamplesOutDec)
		if int(decControl.NChannelsAPI) == 2 {
			for i = 0; i < int(*nSamplesOut); i++ {
				samplesOut[n+i*2] = *(*int16)(unsafe.Add(unsafe.Pointer(resample_out_ptr), unsafe.Sizeof(int16(0))*uintptr(i)))
			}
		}
	}
	if int(decControl.NChannelsAPI) == 2 && int(decControl.NChannelsInternal) == 1 {
		if stereo_to_mono != 0 {
			ret += silk_resampler(&(*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*1))).Resampler_state, []int16(resample_out_ptr), []int16((*int16)(unsafe.Add(unsafe.Pointer(samplesOut1_tmp[0]), unsafe.Sizeof(int16(0))*1))), nSamplesOutDec)
			for i = 0; i < int(*nSamplesOut); i++ {
				samplesOut[i*2+1] = *(*int16)(unsafe.Add(unsafe.Pointer(resample_out_ptr), unsafe.Sizeof(int16(0))*uintptr(i)))
			}
		} else {
			for i = 0; i < int(*nSamplesOut); i++ {
				samplesOut[i*2+1] = samplesOut[i*2+0]
			}
		}
	}
	if (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).PrevSignalType == TYPE_VOICED {
		var mult_tab [3]int = [3]int{6, 4, 3}
		decControl.PrevPitchLag = (*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).LagPrev * mult_tab[((*(*silk_decoder_state)(unsafe.Add(unsafe.Pointer(channel_state), unsafe.Sizeof(silk_decoder_state{})*0))).Fs_kHz-8)>>2]
	} else {
		decControl.PrevPitchLag = 0
	}
	if lostFlag == FLAG_PACKET_LOST {
		for i = 0; i < psDec.NChannelsInternal; i++ {
			psDec.Channel_state[i].LastGainIndex = 10
		}
	} else {
		psDec.Prev_decode_only_middle = decode_only_middle
	}
	return ret
}
