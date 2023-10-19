package libopus

import "unsafe"

func silk_decode_pitch(lagIndex opus_int16, contourIndex int8, pitch_lags [0]int64, Fs_kHz int64, nb_subfr int64) {
	var (
		lag        int64
		k          int64
		min_lag    int64
		max_lag    int64
		cbk_size   int64
		Lag_CB_ptr *int8
	)
	if Fs_kHz == 8 {
		if nb_subfr == PE_MAX_NB_SUBFR {
			Lag_CB_ptr = &silk_CB_lags_stage2[0][0]
			cbk_size = PE_NB_CBKS_STAGE2_EXT
		} else {
			Lag_CB_ptr = &silk_CB_lags_stage2_10_ms[0][0]
			cbk_size = PE_NB_CBKS_STAGE2_10MS
		}
	} else {
		if nb_subfr == PE_MAX_NB_SUBFR {
			Lag_CB_ptr = &silk_CB_lags_stage3[0][0]
			cbk_size = PE_NB_CBKS_STAGE3_MAX
		} else {
			Lag_CB_ptr = &silk_CB_lags_stage3_10_ms[0][0]
			cbk_size = PE_NB_CBKS_STAGE3_10MS
		}
	}
	min_lag = int64(PE_MIN_LAG_MS * opus_int32(opus_int16(Fs_kHz)))
	max_lag = int64(PE_MAX_LAG_MS * opus_int32(opus_int16(Fs_kHz)))
	lag = min_lag + int64(lagIndex)
	for k = 0; k < nb_subfr; k++ {
		pitch_lags[k] = lag + int64(*((*int8)(unsafe.Add(unsafe.Pointer(Lag_CB_ptr), k*cbk_size+int64(contourIndex)))))
		if min_lag > max_lag {
			if (pitch_lags[k]) > min_lag {
				pitch_lags[k] = min_lag
			} else if (pitch_lags[k]) < max_lag {
				pitch_lags[k] = max_lag
			} else {
				pitch_lags[k] = pitch_lags[k]
			}
		} else if (pitch_lags[k]) > max_lag {
			pitch_lags[k] = max_lag
		} else if (pitch_lags[k]) < min_lag {
			pitch_lags[k] = min_lag
		} else {
			pitch_lags[k] = pitch_lags[k]
		}
	}
}
