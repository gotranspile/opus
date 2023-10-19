package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_control_encoder(psEnc *silk_encoder_state_FLP, encControl *silk_EncControlStruct, allow_bw_switch int64, channelNb int64, force_fs_kHz int64) int64 {
	var (
		fs_kHz int64
		ret    int64 = 0
	)
	psEnc.SCmn.UseDTX = encControl.UseDTX
	psEnc.SCmn.UseCBR = encControl.UseCBR
	psEnc.SCmn.API_fs_Hz = encControl.API_sampleRate
	psEnc.SCmn.MaxInternal_fs_Hz = int64(encControl.MaxInternalSampleRate)
	psEnc.SCmn.MinInternal_fs_Hz = int64(encControl.MinInternalSampleRate)
	psEnc.SCmn.DesiredInternal_fs_Hz = int64(encControl.DesiredInternalSampleRate)
	psEnc.SCmn.UseInBandFEC = encControl.UseInBandFEC
	psEnc.SCmn.NChannelsAPI = int64(encControl.NChannelsAPI)
	psEnc.SCmn.NChannelsInternal = int64(encControl.NChannelsInternal)
	psEnc.SCmn.Allow_bandwidth_switch = allow_bw_switch
	psEnc.SCmn.ChannelNb = channelNb
	if psEnc.SCmn.Controlled_since_last_payload != 0 && psEnc.SCmn.PrefillFlag == 0 {
		if psEnc.SCmn.API_fs_Hz != psEnc.SCmn.Prev_API_fs_Hz && psEnc.SCmn.Fs_kHz > 0 {
			ret += silk_setup_resamplers(psEnc, psEnc.SCmn.Fs_kHz)
		}
		return ret
	}
	fs_kHz = silk_control_audio_bandwidth(&psEnc.SCmn, encControl)
	if force_fs_kHz != 0 {
		fs_kHz = force_fs_kHz
	}
	ret += silk_setup_resamplers(psEnc, fs_kHz)
	ret += silk_setup_fs(psEnc, fs_kHz, encControl.PayloadSize_ms)
	ret += silk_setup_complexity(&psEnc.SCmn, encControl.Complexity)
	psEnc.SCmn.PacketLoss_perc = encControl.PacketLossPercentage
	ret += silk_setup_LBRR(&psEnc.SCmn, encControl)
	psEnc.SCmn.Controlled_since_last_payload = 1
	return ret
}
func silk_setup_resamplers(psEnc *silk_encoder_state_FLP, fs_kHz int64) int64 {
	var ret int64 = SILK_NO_ERROR
	if psEnc.SCmn.Fs_kHz != fs_kHz || psEnc.SCmn.Prev_API_fs_Hz != psEnc.SCmn.API_fs_Hz {
		if psEnc.SCmn.Fs_kHz == 0 {
			ret += silk_resampler_init(&psEnc.SCmn.Resampler_state, psEnc.SCmn.API_fs_Hz, opus_int32(fs_kHz*1000), 1)
		} else {
			var (
				x_buf_API_fs_Hz      *opus_int16
				temp_resampler_state *silk_resampler_state_struct
				x_bufFIX             *opus_int16
				new_buf_samples      opus_int32
				api_buf_samples      opus_int32
				old_buf_samples      opus_int32
				buf_length_ms        opus_int32
			)
			buf_length_ms = (opus_int32(opus_uint32(psEnc.SCmn.Nb_subfr*5) << 1)) + LA_SHAPE_MS
			old_buf_samples = buf_length_ms * opus_int32(psEnc.SCmn.Fs_kHz)
			new_buf_samples = buf_length_ms * opus_int32(fs_kHz)
			x_bufFIX = (*opus_int16)(libc.Malloc(int((func() opus_int32 {
				if old_buf_samples > new_buf_samples {
					return old_buf_samples
				}
				return new_buf_samples
			}()) * opus_int32(unsafe.Sizeof(opus_int16(0))))))
			silk_float2short_array(x_bufFIX, &psEnc.X_buf[0], old_buf_samples)
			temp_resampler_state = (*silk_resampler_state_struct)(libc.Malloc(int(unsafe.Sizeof(silk_resampler_state_struct{}) * 1)))
			ret += silk_resampler_init(temp_resampler_state, opus_int32(opus_int16(psEnc.SCmn.Fs_kHz))*1000, psEnc.SCmn.API_fs_Hz, 0)
			api_buf_samples = buf_length_ms * (psEnc.SCmn.API_fs_Hz / 1000)
			x_buf_API_fs_Hz = (*opus_int16)(libc.Malloc(int(api_buf_samples * opus_int32(unsafe.Sizeof(opus_int16(0))))))
			ret += silk_resampler(temp_resampler_state, [0]opus_int16(x_buf_API_fs_Hz), [0]opus_int16(x_bufFIX), old_buf_samples)
			ret += silk_resampler_init(&psEnc.SCmn.Resampler_state, psEnc.SCmn.API_fs_Hz, opus_int32(opus_int16(fs_kHz))*1000, 1)
			ret += silk_resampler(&psEnc.SCmn.Resampler_state, [0]opus_int16(x_bufFIX), [0]opus_int16(x_buf_API_fs_Hz), api_buf_samples)
			silk_short2float_array(&psEnc.X_buf[0], x_bufFIX, new_buf_samples)
		}
	}
	psEnc.SCmn.Prev_API_fs_Hz = psEnc.SCmn.API_fs_Hz
	return ret
}
func silk_setup_fs(psEnc *silk_encoder_state_FLP, fs_kHz int64, PacketSize_ms int64) int64 {
	var ret int64 = SILK_NO_ERROR
	if PacketSize_ms != psEnc.SCmn.PacketSize_ms {
		if PacketSize_ms != 10 && PacketSize_ms != 20 && PacketSize_ms != 40 && PacketSize_ms != 60 {
			ret = -103
		}
		if PacketSize_ms <= 10 {
			psEnc.SCmn.NFramesPerPacket = 1
			if PacketSize_ms == 10 {
				psEnc.SCmn.Nb_subfr = 2
			} else {
				psEnc.SCmn.Nb_subfr = 1
			}
			psEnc.SCmn.Frame_length = int64(opus_int32(opus_int16(PacketSize_ms)) * opus_int32(opus_int16(fs_kHz)))
			psEnc.SCmn.Pitch_LPC_win_length = int64(opus_int32(opus_int16((LA_PITCH_MS<<1)+10)) * opus_int32(opus_int16(fs_kHz)))
			if psEnc.SCmn.Fs_kHz == 8 {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_10_ms_NB_iCDF[0]
			} else {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_10_ms_iCDF[0]
			}
		} else {
			psEnc.SCmn.NFramesPerPacket = int64(opus_int32(PacketSize_ms / (SUB_FRAME_LENGTH_MS * MAX_NB_SUBFR)))
			psEnc.SCmn.Nb_subfr = MAX_NB_SUBFR
			psEnc.SCmn.Frame_length = int64(opus_int32(opus_int16(fs_kHz)) * 20)
			psEnc.SCmn.Pitch_LPC_win_length = int64(opus_int32(opus_int16((LA_PITCH_MS<<1)+20)) * opus_int32(opus_int16(fs_kHz)))
			if psEnc.SCmn.Fs_kHz == 8 {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_NB_iCDF[0]
			} else {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_iCDF[0]
			}
		}
		psEnc.SCmn.PacketSize_ms = PacketSize_ms
		psEnc.SCmn.TargetRate_bps = 0
	}
	if psEnc.SCmn.Fs_kHz != fs_kHz {
		psEnc.SShape = silk_shape_state_FLP{}
		psEnc.SCmn.SNSQ = silk_nsq_state{}
		*(*[16]opus_int16)(unsafe.Pointer(&psEnc.SCmn.Prev_NLSFq_Q15[0])) = [16]opus_int16{}
		*(*[2]opus_int32)(unsafe.Pointer(&psEnc.SCmn.SLP.In_LP_State[0])) = [2]opus_int32{}
		psEnc.SCmn.InputBufIx = 0
		psEnc.SCmn.NFramesEncoded = 0
		psEnc.SCmn.TargetRate_bps = 0
		psEnc.SCmn.PrevLag = 100
		psEnc.SCmn.First_frame_after_reset = 1
		psEnc.SShape.LastGainIndex = 10
		psEnc.SCmn.SNSQ.LagPrev = 100
		psEnc.SCmn.SNSQ.Prev_gain_Q16 = 65536
		psEnc.SCmn.PrevSignalType = TYPE_NO_VOICE_ACTIVITY
		psEnc.SCmn.Fs_kHz = fs_kHz
		if psEnc.SCmn.Fs_kHz == 8 {
			if psEnc.SCmn.Nb_subfr == MAX_NB_SUBFR {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_NB_iCDF[0]
			} else {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_10_ms_NB_iCDF[0]
			}
		} else {
			if psEnc.SCmn.Nb_subfr == MAX_NB_SUBFR {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_iCDF[0]
			} else {
				psEnc.SCmn.Pitch_contour_iCDF = &silk_pitch_contour_10_ms_iCDF[0]
			}
		}
		if psEnc.SCmn.Fs_kHz == 8 || psEnc.SCmn.Fs_kHz == 12 {
			psEnc.SCmn.PredictLPCOrder = MIN_LPC_ORDER
			psEnc.SCmn.PsNLSF_CB = &silk_NLSF_CB_NB_MB
		} else {
			psEnc.SCmn.PredictLPCOrder = MAX_LPC_ORDER
			psEnc.SCmn.PsNLSF_CB = &silk_NLSF_CB_WB
		}
		psEnc.SCmn.Subfr_length = SUB_FRAME_LENGTH_MS * fs_kHz
		psEnc.SCmn.Frame_length = int64(opus_int32(opus_int16(psEnc.SCmn.Subfr_length)) * opus_int32(opus_int16(psEnc.SCmn.Nb_subfr)))
		psEnc.SCmn.Ltp_mem_length = int64(LTP_MEM_LENGTH_MS * opus_int32(opus_int16(fs_kHz)))
		psEnc.SCmn.La_pitch = int64(LA_PITCH_MS * opus_int32(opus_int16(fs_kHz)))
		psEnc.SCmn.Max_pitch_lag = int64(opus_int32(opus_int16(fs_kHz)) * 18)
		if psEnc.SCmn.Nb_subfr == MAX_NB_SUBFR {
			psEnc.SCmn.Pitch_LPC_win_length = int64(opus_int32(opus_int16((LA_PITCH_MS<<1)+20)) * opus_int32(opus_int16(fs_kHz)))
		} else {
			psEnc.SCmn.Pitch_LPC_win_length = int64(opus_int32(opus_int16((LA_PITCH_MS<<1)+10)) * opus_int32(opus_int16(fs_kHz)))
		}
		if psEnc.SCmn.Fs_kHz == 16 {
			psEnc.SCmn.Pitch_lag_low_bits_iCDF = &silk_uniform8_iCDF[0]
		} else if psEnc.SCmn.Fs_kHz == 12 {
			psEnc.SCmn.Pitch_lag_low_bits_iCDF = &silk_uniform6_iCDF[0]
		} else {
			psEnc.SCmn.Pitch_lag_low_bits_iCDF = &silk_uniform4_iCDF[0]
		}
	}
	return ret
}
func silk_setup_complexity(psEncC *silk_encoder_state, Complexity int64) int64 {
	var ret int64 = 0
	if Complexity < 1 {
		psEncC.PitchEstimationComplexity = SILK_PE_MIN_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.8*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 6
		psEncC.ShapingLPCOrder = 12
		psEncC.La_shape = psEncC.Fs_kHz * 3
		psEncC.NStatesDelayedDecision = 1
		psEncC.UseInterpolatedNLSFs = 0
		psEncC.NLSF_MSVQ_Survivors = 2
		psEncC.Warping_Q16 = 0
	} else if Complexity < 2 {
		psEncC.PitchEstimationComplexity = SILK_PE_MID_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.76*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 8
		psEncC.ShapingLPCOrder = 14
		psEncC.La_shape = psEncC.Fs_kHz * 5
		psEncC.NStatesDelayedDecision = 1
		psEncC.UseInterpolatedNLSFs = 0
		psEncC.NLSF_MSVQ_Survivors = 3
		psEncC.Warping_Q16 = 0
	} else if Complexity < 3 {
		psEncC.PitchEstimationComplexity = SILK_PE_MIN_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.8*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 6
		psEncC.ShapingLPCOrder = 12
		psEncC.La_shape = psEncC.Fs_kHz * 3
		psEncC.NStatesDelayedDecision = 2
		psEncC.UseInterpolatedNLSFs = 0
		psEncC.NLSF_MSVQ_Survivors = 2
		psEncC.Warping_Q16 = 0
	} else if Complexity < 4 {
		psEncC.PitchEstimationComplexity = SILK_PE_MID_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.76*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 8
		psEncC.ShapingLPCOrder = 14
		psEncC.La_shape = psEncC.Fs_kHz * 5
		psEncC.NStatesDelayedDecision = 2
		psEncC.UseInterpolatedNLSFs = 0
		psEncC.NLSF_MSVQ_Survivors = 4
		psEncC.Warping_Q16 = 0
	} else if Complexity < 6 {
		psEncC.PitchEstimationComplexity = SILK_PE_MID_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.74*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 10
		psEncC.ShapingLPCOrder = 16
		psEncC.La_shape = psEncC.Fs_kHz * 5
		psEncC.NStatesDelayedDecision = 2
		psEncC.UseInterpolatedNLSFs = 1
		psEncC.NLSF_MSVQ_Survivors = 6
		psEncC.Warping_Q16 = psEncC.Fs_kHz * int64(opus_int32(WARPING_MULTIPLIER*(1<<16)+0.5))
	} else if Complexity < 8 {
		psEncC.PitchEstimationComplexity = SILK_PE_MID_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.72*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 12
		psEncC.ShapingLPCOrder = 20
		psEncC.La_shape = psEncC.Fs_kHz * 5
		psEncC.NStatesDelayedDecision = 3
		psEncC.UseInterpolatedNLSFs = 1
		psEncC.NLSF_MSVQ_Survivors = 8
		psEncC.Warping_Q16 = psEncC.Fs_kHz * int64(opus_int32(WARPING_MULTIPLIER*(1<<16)+0.5))
	} else {
		psEncC.PitchEstimationComplexity = SILK_PE_MAX_COMPLEX
		psEncC.PitchEstimationThreshold_Q16 = opus_int32(0.7*(1<<16) + 0.5)
		psEncC.PitchEstimationLPCOrder = 16
		psEncC.ShapingLPCOrder = 24
		psEncC.La_shape = psEncC.Fs_kHz * 5
		psEncC.NStatesDelayedDecision = MAX_DEL_DEC_STATES
		psEncC.UseInterpolatedNLSFs = 1
		psEncC.NLSF_MSVQ_Survivors = 16
		psEncC.Warping_Q16 = psEncC.Fs_kHz * int64(opus_int32(WARPING_MULTIPLIER*(1<<16)+0.5))
	}
	psEncC.PitchEstimationLPCOrder = silk_min_int(psEncC.PitchEstimationLPCOrder, psEncC.PredictLPCOrder)
	psEncC.ShapeWinLength = SUB_FRAME_LENGTH_MS*psEncC.Fs_kHz + psEncC.La_shape*2
	psEncC.Complexity = Complexity
	return ret
}
func silk_setup_LBRR(psEncC *silk_encoder_state, encControl *silk_EncControlStruct) int64 {
	var (
		LBRR_in_previous_packet int64
		ret                     int64 = SILK_NO_ERROR
	)
	LBRR_in_previous_packet = psEncC.LBRR_enabled
	psEncC.LBRR_enabled = encControl.LBRR_coded
	if psEncC.LBRR_enabled != 0 {
		if LBRR_in_previous_packet == 0 {
			psEncC.LBRR_GainIncreases = 7
		} else {
			psEncC.LBRR_GainIncreases = silk_max_int(int64(7-(((opus_int32(psEncC.PacketLoss_perc))*opus_int32(int64(opus_int16(opus_int32(0.2*(1<<16)+0.5)))))>>16)), 3)
		}
	}
	return ret
}
