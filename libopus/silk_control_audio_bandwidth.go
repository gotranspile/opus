package libopus

import "unsafe"

func silk_control_audio_bandwidth(psEncC *silk_encoder_state, encControl *silk_EncControlStruct) int64 {
	var (
		fs_kHz   int64
		orig_kHz int64
		fs_Hz    opus_int32
	)
	orig_kHz = psEncC.Fs_kHz
	if orig_kHz == 0 {
		orig_kHz = int64(psEncC.SLP.Saved_fs_kHz)
	}
	fs_kHz = orig_kHz
	fs_Hz = opus_int32(opus_int16(fs_kHz)) * 1000
	if fs_Hz == 0 {
		if psEncC.DesiredInternal_fs_Hz < int64(psEncC.API_fs_Hz) {
			fs_Hz = opus_int32(psEncC.DesiredInternal_fs_Hz)
		} else {
			fs_Hz = psEncC.API_fs_Hz
		}
		fs_kHz = int64(fs_Hz / 1000)
	} else if fs_Hz > psEncC.API_fs_Hz || fs_Hz > opus_int32(psEncC.MaxInternal_fs_Hz) || fs_Hz < opus_int32(psEncC.MinInternal_fs_Hz) {
		fs_Hz = psEncC.API_fs_Hz
		if fs_Hz < opus_int32(psEncC.MaxInternal_fs_Hz) {
			fs_Hz = fs_Hz
		} else {
			fs_Hz = opus_int32(psEncC.MaxInternal_fs_Hz)
		}
		if fs_Hz > opus_int32(psEncC.MinInternal_fs_Hz) {
			fs_Hz = fs_Hz
		} else {
			fs_Hz = opus_int32(psEncC.MinInternal_fs_Hz)
		}
		fs_kHz = int64(fs_Hz / 1000)
	} else {
		if psEncC.SLP.Transition_frame_no >= opus_int32(TRANSITION_TIME_MS/(SUB_FRAME_LENGTH_MS*MAX_NB_SUBFR)) {
			psEncC.SLP.Mode = 0
		}
		if psEncC.Allow_bandwidth_switch != 0 || encControl.OpusCanSwitch != 0 {
			if (opus_int32(opus_int16(orig_kHz)) * 1000) > opus_int32(psEncC.DesiredInternal_fs_Hz) {
				if psEncC.SLP.Mode == 0 {
					psEncC.SLP.Transition_frame_no = opus_int32(TRANSITION_TIME_MS / (SUB_FRAME_LENGTH_MS * MAX_NB_SUBFR))
					*(*[2]opus_int32)(unsafe.Pointer(&psEncC.SLP.In_LP_State[0])) = [2]opus_int32{}
				}
				if encControl.OpusCanSwitch != 0 {
					psEncC.SLP.Mode = 0
					if orig_kHz == 16 {
						fs_kHz = 12
					} else {
						fs_kHz = 8
					}
				} else {
					if psEncC.SLP.Transition_frame_no <= 0 {
						encControl.SwitchReady = 1
						encControl.MaxBits -= encControl.MaxBits * 5 / (encControl.PayloadSize_ms + 5)
					} else {
						psEncC.SLP.Mode = -2
					}
				}
			} else if (opus_int32(opus_int16(orig_kHz)) * 1000) < opus_int32(psEncC.DesiredInternal_fs_Hz) {
				if encControl.OpusCanSwitch != 0 {
					if orig_kHz == 8 {
						fs_kHz = 12
					} else {
						fs_kHz = 16
					}
					psEncC.SLP.Transition_frame_no = 0
					*(*[2]opus_int32)(unsafe.Pointer(&psEncC.SLP.In_LP_State[0])) = [2]opus_int32{}
					psEncC.SLP.Mode = 1
				} else {
					if psEncC.SLP.Mode == 0 {
						encControl.SwitchReady = 1
						encControl.MaxBits -= encControl.MaxBits * 5 / (encControl.PayloadSize_ms + 5)
					} else {
						psEncC.SLP.Mode = 1
					}
				}
			} else {
				if psEncC.SLP.Mode < 0 {
					psEncC.SLP.Mode = 1
				}
			}
		}
	}
	return fs_kHz
}
