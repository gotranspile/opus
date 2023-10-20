package libopus

type silk_nsq_state struct {
	Xq               [640]int16
	SLTP_shp_Q14     [640]int32
	SLPC_Q14         [96]int32
	SAR2_Q14         [24]int32
	SLF_AR_shp_Q14   int32
	SDiff_shp_Q14    int32
	LagPrev          int
	SLTP_buf_idx     int
	SLTP_shp_buf_idx int
	Rand_seed        int32
	Prev_gain_Q16    int32
	Rewhite_flag     int
}
type silk_VAD_state struct {
	AnaState        [2]int32
	AnaState1       [2]int32
	AnaState2       [2]int32
	XnrgSubfr       [4]int32
	NrgRatioSmth_Q8 [4]int32
	HPstate         int16
	NL              [4]int32
	Inv_NL          [4]int32
	NoiseLevelBias  [4]int32
	Counter         int32
}
type silk_LP_state struct {
	In_LP_State         [2]int32
	Transition_frame_no int32
	Mode                int
	Saved_fs_kHz        int32
}
type silk_NLSF_CB_struct struct {
	NVectors            int16
	Order               int16
	QuantStepSize_Q16   int16
	InvQuantStepSize_Q6 int16
	CB1_NLSF_Q8         *uint8
	CB1_Wght_Q9         *int16
	CB1_iCDF            *uint8
	Pred_Q8             *uint8
	Ec_sel              *uint8
	Ec_iCDF             *uint8
	Ec_Rates_Q5         *uint8
	DeltaMin_Q15        *int16
}
type stereo_enc_state struct {
	Pred_prev_Q13   [2]int16
	SMid            [2]int16
	SSide           [2]int16
	Mid_side_amp_Q0 [4]int32
	Smth_width_Q14  int16
	Width_prev_Q14  int16
	Silent_side_len int16
	PredIx          [3][2][3]int8
	Mid_only_flags  [3]int8
}
type stereo_dec_state struct {
	Pred_prev_Q13 [2]int16
	SMid          [2]int16
	SSide         [2]int16
}
type SideInfoIndices struct {
	GainsIndices      [4]int8
	LTPIndex          [4]int8
	NLSFIndices       [17]int8
	LagIndex          int16
	ContourIndex      int8
	SignalType        int8
	QuantOffsetType   int8
	NLSFInterpCoef_Q2 int8
	PERIndex          int8
	LTP_scaleIndex    int8
	Seed              int8
}
type silk_encoder_state struct {
	In_HP_State                   [2]int32
	Variable_HP_smth1_Q15         int32
	Variable_HP_smth2_Q15         int32
	SLP                           silk_LP_state
	SVAD                          silk_VAD_state
	SNSQ                          silk_nsq_state
	Prev_NLSFq_Q15                [16]int16
	Speech_activity_Q8            int
	Allow_bandwidth_switch        int
	LBRRprevLastGainIndex         int8
	PrevSignalType                int8
	PrevLag                       int
	Pitch_LPC_win_length          int
	Max_pitch_lag                 int
	API_fs_Hz                     int32
	Prev_API_fs_Hz                int32
	MaxInternal_fs_Hz             int
	MinInternal_fs_Hz             int
	DesiredInternal_fs_Hz         int
	Fs_kHz                        int
	Nb_subfr                      int
	Frame_length                  int
	Subfr_length                  int
	Ltp_mem_length                int
	La_pitch                      int
	La_shape                      int
	ShapeWinLength                int
	TargetRate_bps                int32
	PacketSize_ms                 int
	PacketLoss_perc               int
	FrameCounter                  int32
	Complexity                    int
	NStatesDelayedDecision        int
	UseInterpolatedNLSFs          int
	ShapingLPCOrder               int
	PredictLPCOrder               int
	PitchEstimationComplexity     int
	PitchEstimationLPCOrder       int
	PitchEstimationThreshold_Q16  int32
	Sum_log_gain_Q7               int32
	NLSF_MSVQ_Survivors           int
	First_frame_after_reset       int
	Controlled_since_last_payload int
	Warping_Q16                   int
	UseCBR                        int
	PrefillFlag                   int
	Pitch_lag_low_bits_iCDF       *uint8
	Pitch_contour_iCDF            *uint8
	PsNLSF_CB                     *silk_NLSF_CB_struct
	Input_quality_bands_Q15       [4]int
	Input_tilt_Q15                int
	SNR_dB_Q7                     int
	VAD_flags                     [3]int8
	LBRR_flag                     int8
	LBRR_flags                    [3]int
	Indices                       SideInfoIndices
	Pulses                        [320]int8
	Arch                          int
	InputBuf                      [322]int16
	InputBufIx                    int
	NFramesPerPacket              int
	NFramesEncoded                int
	NChannelsAPI                  int
	NChannelsInternal             int
	ChannelNb                     int
	Frames_since_onset            int
	Ec_prevSignalType             int
	Ec_prevLagIndex               int16
	Resampler_state               silk_resampler_state_struct
	UseDTX                        int
	InDTX                         int
	NoSpeechCounter               int
	UseInBandFEC                  int
	LBRR_enabled                  int
	LBRR_GainIncreases            int
	Indices_LBRR                  [3]SideInfoIndices
	Pulses_LBRR                   [3][320]int8
}
type silk_PLC_struct struct {
	PitchL_Q8         int32
	LTPCoef_Q14       [5]int16
	PrevLPC_Q12       [16]int16
	Last_frame_lost   int
	Rand_seed         int32
	RandScale_Q14     int16
	Conc_energy       int32
	Conc_energy_shift int
	PrevLTP_scale_Q14 int16
	PrevGain_Q16      [2]int32
	Fs_kHz            int
	Nb_subfr          int
	Subfr_length      int
}
type silk_CNG_struct struct {
	CNG_exc_buf_Q14   [320]int32
	CNG_smth_NLSF_Q15 [16]int16
	CNG_synth_state   [16]int32
	CNG_smth_Gain_Q16 int32
	Rand_seed         int32
	Fs_kHz            int
}
type silk_decoder_state struct {
	Prev_gain_Q16           int32
	Exc_Q14                 [320]int32
	SLPC_Q14_buf            [16]int32
	OutBuf                  [480]int16
	LagPrev                 int
	LastGainIndex           int8
	Fs_kHz                  int
	Fs_API_hz               int32
	Nb_subfr                int
	Frame_length            int
	Subfr_length            int
	Ltp_mem_length          int
	LPC_order               int
	PrevNLSF_Q15            [16]int16
	First_frame_after_reset int
	Pitch_lag_low_bits_iCDF *uint8
	Pitch_contour_iCDF      *uint8
	NFramesDecoded          int
	NFramesPerPacket        int
	Ec_prevSignalType       int
	Ec_prevLagIndex         int16
	VAD_flags               [3]int
	LBRR_flag               int
	LBRR_flags              [3]int
	Resampler_state         silk_resampler_state_struct
	PsNLSF_CB               *silk_NLSF_CB_struct
	Indices                 SideInfoIndices
	SCNG                    silk_CNG_struct
	LossCnt                 int
	PrevSignalType          int
	Arch                    int
	SPLC                    silk_PLC_struct
}
type silk_decoder_control struct {
	PitchL        [4]int
	Gains_Q16     [4]int32
	PredCoef_Q12  [2][16]int16
	LTPCoef_Q14   [20]int16
	LTP_scale_Q14 int
}
