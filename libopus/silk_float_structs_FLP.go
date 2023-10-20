package libopus

type silk_shape_state_FLP struct {
	LastGainIndex      int8
	HarmShapeGain_smth float32
	Tilt_smth          float32
}
type silk_encoder_state_FLP struct {
	SCmn    silk_encoder_state
	SShape  silk_shape_state_FLP
	X_buf   [720]float32
	LTPCorr float32
}
type silk_encoder_control_FLP struct {
	Gains             [4]float32
	PredCoef          [2][16]float32
	LTPCoef           [20]float32
	LTP_scale         float32
	PitchL            [4]int
	AR                [96]float32
	LF_MA_shp         [4]float32
	LF_AR_shp         [4]float32
	Tilt              [4]float32
	HarmShapeGain     [4]float32
	Lambda            float32
	Input_quality     float32
	Coding_quality    float32
	PredGain          float32
	LTPredCodGain     float32
	ResNrg            [4]float32
	GainsUnq_Q16      [4]int32
	LastGainIndexPrev int8
}
type silk_encoder struct {
	State_Fxx                 [2]silk_encoder_state_FLP
	SStereo                   stereo_enc_state
	NBitsUsedLBRR             int32
	NBitsExceeded             int32
	NChannelsAPI              int
	NChannelsInternal         int
	NPrevChannelsInternal     int
	TimeSinceSwitchAllowed_ms int
	AllowBandwidthSwitch      int
	Prev_decode_only_middle   int
}
