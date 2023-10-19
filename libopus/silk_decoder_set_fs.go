package libopus

import "unsafe"

func silk_decoder_set_fs(psDec *silk_decoder_state, fs_kHz int64, fs_API_Hz opus_int32) int64 {
	var (
		frame_length int64
		ret          int64 = 0
	)
	psDec.Subfr_length = int64(SUB_FRAME_LENGTH_MS * opus_int32(opus_int16(fs_kHz)))
	frame_length = int64(opus_int32(opus_int16(psDec.Nb_subfr)) * opus_int32(opus_int16(psDec.Subfr_length)))
	if psDec.Fs_kHz != fs_kHz || psDec.Fs_API_hz != fs_API_Hz {
		ret += silk_resampler_init(&psDec.Resampler_state, opus_int32(opus_int16(fs_kHz))*1000, fs_API_Hz, 0)
		psDec.Fs_API_hz = fs_API_Hz
	}
	if psDec.Fs_kHz != fs_kHz || frame_length != psDec.Frame_length {
		if fs_kHz == 8 {
			if psDec.Nb_subfr == MAX_NB_SUBFR {
				psDec.Pitch_contour_iCDF = &silk_pitch_contour_NB_iCDF[0]
			} else {
				psDec.Pitch_contour_iCDF = &silk_pitch_contour_10_ms_NB_iCDF[0]
			}
		} else {
			if psDec.Nb_subfr == MAX_NB_SUBFR {
				psDec.Pitch_contour_iCDF = &silk_pitch_contour_iCDF[0]
			} else {
				psDec.Pitch_contour_iCDF = &silk_pitch_contour_10_ms_iCDF[0]
			}
		}
		if psDec.Fs_kHz != fs_kHz {
			psDec.Ltp_mem_length = int64(LTP_MEM_LENGTH_MS * opus_int32(opus_int16(fs_kHz)))
			if fs_kHz == 8 || fs_kHz == 12 {
				psDec.LPC_order = MIN_LPC_ORDER
				psDec.PsNLSF_CB = &silk_NLSF_CB_NB_MB
			} else {
				psDec.LPC_order = MAX_LPC_ORDER
				psDec.PsNLSF_CB = &silk_NLSF_CB_WB
			}
			if fs_kHz == 16 {
				psDec.Pitch_lag_low_bits_iCDF = &silk_uniform8_iCDF[0]
			} else if fs_kHz == 12 {
				psDec.Pitch_lag_low_bits_iCDF = &silk_uniform6_iCDF[0]
			} else if fs_kHz == 8 {
				psDec.Pitch_lag_low_bits_iCDF = &silk_uniform4_iCDF[0]
			} else {
			}
			psDec.First_frame_after_reset = 1
			psDec.LagPrev = 100
			psDec.LastGainIndex = 10
			psDec.PrevSignalType = TYPE_NO_VOICE_ACTIVITY
			*(*[480]opus_int16)(unsafe.Pointer(&psDec.OutBuf[0])) = [480]opus_int16{}
			*(*[16]opus_int32)(unsafe.Pointer(&psDec.SLPC_Q14_buf[0])) = [16]opus_int32{}
		}
		psDec.Fs_kHz = fs_kHz
		psDec.Frame_length = frame_length
	}
	return ret
}
