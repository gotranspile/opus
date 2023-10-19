package libopus

const SILK_RESAMPLER_MAX_FIR_ORDER = 36
const SILK_RESAMPLER_MAX_IIR_ORDER = 6

type _silk_resampler_state_struct struct {
	SIIR [6]int32
	SFIR struct {
		// union
		I32 [36]int32
		I16 [36]int16
	}
	DelayBuf           [48]int16
	Resampler_function int64
	BatchSize          int64
	InvRatio_Q16       int32
	FIR_Order          int64
	FIR_Fracs          int64
	Fs_in_kHz          int64
	Fs_out_kHz         int64
	InputDelay         int64
	Coefs              *int16
}
type silk_resampler_state_struct _silk_resampler_state_struct
