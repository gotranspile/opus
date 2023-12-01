package silk

const BITRESERVOIR_DECAY_TIME_MS = 500
const FIND_PITCH_WHITE_NOISE_FRACTION = 0.001
const FIND_PITCH_BANDWIDTH_EXPANSION = 0.99
const FIND_LPC_COND_FAC = 1e-05
const MAX_SUM_LOG_GAIN_DB = 250.0
const LTP_CORR_INV_MAX = 0.03
const VARIABLE_HP_SMTH_COEF1 = 0.1
const VARIABLE_HP_SMTH_COEF2 = 0.015
const VARIABLE_HP_MAX_DELTA_FREQ = 0.4
const VARIABLE_HP_MIN_CUTOFF_HZ = 60
const VARIABLE_HP_MAX_CUTOFF_HZ = 100
const SPEECH_ACTIVITY_DTX_THRES = 0.05
const LBRR_SPEECH_ACTIVITY_THRES = 0.3
const BG_SNR_DECR_dB = 2.0
const HARM_SNR_INCR_dB = 2.0
const SPARSE_SNR_INCR_dB = 2.0
const ENERGY_VARIATION_THRESHOLD_QNT_OFFSET = 0.6
const WARPING_MULTIPLIER = 0.015
const SHAPE_WHITE_NOISE_FRACTION = 3e-05
const BANDWIDTH_EXPANSION = 0.94
const HARMONIC_SHAPING = 0.3
const HIGH_RATE_OR_LOW_QUALITY_HARMONIC_SHAPING = 0.2
const HP_NOISE_COEF = 0.25
const HARM_HP_NOISE_COEF = 0.35
const INPUT_TILT = 0.05
const HIGH_RATE_INPUT_TILT = 0.1
const LOW_FREQ_SHAPING = 4.0
const LOW_QUALITY_LOW_FREQ_SHAPING_DECR = 0.5
const SUBFR_SMTH_COEF = 0.4
const LAMBDA_OFFSET = 1.2
const LAMBDA_SPEECH_ACT = 0
const LAMBDA_DELAYED_DECISIONS = 0
const LAMBDA_INPUT_QUALITY = 0
const LAMBDA_CODING_QUALITY = 0
const LAMBDA_QUANT_OFFSET = 0.8
const REDUCE_BITRATE_10_MS_BPS = 2200
const MAX_BANDWIDTH_SWITCH_DELAY_MS = 5000