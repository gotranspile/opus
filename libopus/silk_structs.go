package libopus

type silk_nsq_state struct {
	Xq               [640]opus_int16
	SLTP_shp_Q14     [640]opus_int32
	SLPC_Q14         [96]opus_int32
	SAR2_Q14         [24]opus_int32
	SLF_AR_shp_Q14   opus_int32
	SDiff_shp_Q14    opus_int32
	LagPrev          int64
	SLTP_buf_idx     int64
	SLTP_shp_buf_idx int64
	Rand_seed        opus_int32
	Prev_gain_Q16    opus_int32
	Rewhite_flag     int64
}
type silk_VAD_state struct {
	AnaState        [2]opus_int32
	AnaState1       [2]opus_int32
	AnaState2       [2]opus_int32
	XnrgSubfr       [4]opus_int32
	NrgRatioSmth_Q8 [4]opus_int32
	HPstate         opus_int16
	NL              [4]opus_int32
	Inv_NL          [4]opus_int32
	NoiseLevelBias  [4]opus_int32
	Counter         opus_int32
}
type silk_LP_state struct {
	In_LP_State         [2]opus_int32
	Transition_frame_no opus_int32
	Mode                int64
	Saved_fs_kHz        opus_int32
}
type silk_NLSF_CB_struct struct {
	NVectors            opus_int16
	Order               opus_int16
	QuantStepSize_Q16   opus_int16
	InvQuantStepSize_Q6 opus_int16
	CB1_NLSF_Q8         *uint8
	CB1_Wght_Q9         *opus_int16
	CB1_iCDF            *uint8
	Pred_Q8             *uint8
	Ec_sel              *uint8
	Ec_iCDF             *uint8
	Ec_Rates_Q5         *uint8
	DeltaMin_Q15        *opus_int16
}
type stereo_enc_state struct {
	Pred_prev_Q13   [2]opus_int16
	SMid            [2]opus_int16
	SSide           [2]opus_int16
	Mid_side_amp_Q0 [4]opus_int32
	Smth_width_Q14  opus_int16
	Width_prev_Q14  opus_int16
	Silent_side_len opus_int16
	PredIx          [3][2][3]int8
	Mid_only_flags  [3]int8
}
type stereo_dec_state struct {
	Pred_prev_Q13 [2]opus_int16
	SMid          [2]opus_int16
	SSide         [2]opus_int16
}
type SideInfoIndices struct {
	GainsIndices      [4]int8
	LTPIndex          [4]int8
	NLSFIndices       [17]int8
	LagIndex          opus_int16
	ContourIndex      int8
	SignalType        int8
	QuantOffsetType   int8
	NLSFInterpCoef_Q2 int8
	PERIndex          int8
	LTP_scaleIndex    int8
	Seed              int8
}
type silk_encoder_state struct {
	In_HP_State                   [2]opus_int32
	Variable_HP_smth1_Q15         opus_int32
	Variable_HP_smth2_Q15         opus_int32
	SLP                           silk_LP_state
	SVAD                          silk_VAD_state
	SNSQ                          silk_nsq_state
	Prev_NLSFq_Q15                [16]opus_int16
	Speech_activity_Q8            int64
	Allow_bandwidth_switch        int64
	LBRRprevLastGainIndex         int8
	PrevSignalType                int8
	PrevLag                       int64
	Pitch_LPC_win_length          int64
	Max_pitch_lag                 int64
	API_fs_Hz                     opus_int32
	Prev_API_fs_Hz                opus_int32
	MaxInternal_fs_Hz             int64
	MinInternal_fs_Hz             int64
	DesiredInternal_fs_Hz         int64
	Fs_kHz                        int64
	Nb_subfr                      int64
	Frame_length                  int64
	Subfr_length                  int64
	Ltp_mem_length                int64
	La_pitch                      int64
	La_shape                      int64
	ShapeWinLength                int64
	TargetRate_bps                opus_int32
	PacketSize_ms                 int64
	PacketLoss_perc               int64
	FrameCounter                  opus_int32
	Complexity                    int64
	NStatesDelayedDecision        int64
	UseInterpolatedNLSFs          int64
	ShapingLPCOrder               int64
	PredictLPCOrder               int64
	PitchEstimationComplexity     int64
	PitchEstimationLPCOrder       int64
	PitchEstimationThreshold_Q16  opus_int32
	Sum_log_gain_Q7               opus_int32
	NLSF_MSVQ_Survivors           int64
	First_frame_after_reset       int64
	Controlled_since_last_payload int64
	Warping_Q16                   int64
	UseCBR                        int64
	PrefillFlag                   int64
	Pitch_lag_low_bits_iCDF       *uint8
	Pitch_contour_iCDF            *uint8
	PsNLSF_CB                     *silk_NLSF_CB_struct
	Input_quality_bands_Q15       [4]int64
	Input_tilt_Q15                int64
	SNR_dB_Q7                     int64
	VAD_flags                     [3]int8
	LBRR_flag                     int8
	LBRR_flags                    [3]int64
	Indices                       SideInfoIndices
	Pulses                        [320]int8
	Arch                          int64
	InputBuf                      [322]opus_int16
	InputBufIx                    int64
	NFramesPerPacket              int64
	NFramesEncoded                int64
	NChannelsAPI                  int64
	NChannelsInternal             int64
	ChannelNb                     int64
	Frames_since_onset            int64
	Ec_prevSignalType             int64
	Ec_prevLagIndex               opus_int16
	Resampler_state               silk_resampler_state_struct
	UseDTX                        int64
	InDTX                         int64
	NoSpeechCounter               int64
	UseInBandFEC                  int64
	LBRR_enabled                  int64
	LBRR_GainIncreases            int64
	Indices_LBRR                  [3]SideInfoIndices
	Pulses_LBRR                   [3][320]int8
}
type silk_PLC_struct struct {
	PitchL_Q8         opus_int32
	LTPCoef_Q14       [5]opus_int16
	PrevLPC_Q12       [16]opus_int16
	Last_frame_lost   int64
	Rand_seed         opus_int32
	RandScale_Q14     opus_int16
	Conc_energy       opus_int32
	Conc_energy_shift int64
	PrevLTP_scale_Q14 opus_int16
	PrevGain_Q16      [2]opus_int32
	Fs_kHz            int64
	Nb_subfr          int64
	Subfr_length      int64
}
type silk_CNG_struct struct {
	CNG_exc_buf_Q14   [320]opus_int32
	CNG_smth_NLSF_Q15 [16]opus_int16
	CNG_synth_state   [16]opus_int32
	CNG_smth_Gain_Q16 opus_int32
	Rand_seed         opus_int32
	Fs_kHz            int64
}
type silk_decoder_state struct {
	Prev_gain_Q16           opus_int32
	Exc_Q14                 [320]opus_int32
	SLPC_Q14_buf            [16]opus_int32
	OutBuf                  [480]opus_int16
	LagPrev                 int64
	LastGainIndex           int8
	Fs_kHz                  int64
	Fs_API_hz               opus_int32
	Nb_subfr                int64
	Frame_length            int64
	Subfr_length            int64
	Ltp_mem_length          int64
	LPC_order               int64
	PrevNLSF_Q15            [16]opus_int16
	First_frame_after_reset int64
	Pitch_lag_low_bits_iCDF *uint8
	Pitch_contour_iCDF      *uint8
	NFramesDecoded          int64
	NFramesPerPacket        int64
	Ec_prevSignalType       int64
	Ec_prevLagIndex         opus_int16
	VAD_flags               [3]int64
	LBRR_flag               int64
	LBRR_flags              [3]int64
	Resampler_state         silk_resampler_state_struct
	PsNLSF_CB               *silk_NLSF_CB_struct
	Indices                 SideInfoIndices
	SCNG                    silk_CNG_struct
	LossCnt                 int64
	PrevSignalType          int64
	Arch                    int64
	SPLC                    silk_PLC_struct
}
type silk_decoder_control struct {
	PitchL        [4]int64
	Gains_Q16     [4]opus_int32
	PredCoef_Q12  [2][16]opus_int16
	LTPCoef_Q14   [20]opus_int16
	LTP_scale_Q14 int64
}
