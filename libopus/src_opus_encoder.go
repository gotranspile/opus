package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const MAX_ENCODER_BUFFER = 480
const PSEUDO_SNR_THRESHOLD = 316.23

type StereoWidthState struct {
	XX             opus_val32
	XY             opus_val32
	YY             opus_val32
	Smoothed_width opus_val16
	Max_follower   opus_val16
}
type OpusEncoder struct {
	Celt_enc_offset         int64
	Silk_enc_offset         int64
	Silk_mode               silk_EncControlStruct
	Application             int64
	Channels                int64
	Delay_compensation      int64
	Force_channels          int64
	Signal_type             int64
	User_bandwidth          int64
	Max_bandwidth           int64
	User_forced_mode        int64
	Voice_ratio             int64
	Fs                      opus_int32
	Use_vbr                 int64
	Vbr_constraint          int64
	Variable_duration       int64
	Bitrate_bps             opus_int32
	User_bitrate_bps        opus_int32
	Lsb_depth               int64
	Encoder_buffer          int64
	Lfe                     int64
	Arch                    int64
	Use_dtx                 int64
	Fec_config              int64
	Analysis                TonalityAnalysisState
	Stream_channels         int64
	Hybrid_stereo_width_Q14 opus_int16
	Variable_HP_smth2_Q15   opus_int32
	Prev_HB_gain            opus_val16
	Hp_mem                  [4]opus_val32
	Mode                    int64
	Prev_mode               int64
	Prev_channels           int64
	Prev_framesize          int64
	Bandwidth               int64
	Auto_bandwidth          int64
	Silk_bw_switch          int64
	First                   int64
	Energy_masking          *opus_val16
	Width_mem               StereoWidthState
	Delay_buffer            [960]opus_val16
	Detected_bandwidth      int64
	Nb_no_activity_ms_Q1    int64
	Peak_signal_energy      opus_val32
	Nonfinal_frame          int64
	RangeFinal              opus_uint32
}

var mono_voice_bandwidth_thresholds [8]opus_int32 = [8]opus_int32{9000, 700, 9000, 700, 13500, 1000, 14000, 2000}
var mono_music_bandwidth_thresholds [8]opus_int32 = [8]opus_int32{9000, 700, 9000, 700, 11000, 1000, 12000, 2000}
var stereo_voice_bandwidth_thresholds [8]opus_int32 = [8]opus_int32{9000, 700, 9000, 700, 13500, 1000, 14000, 2000}
var stereo_music_bandwidth_thresholds [8]opus_int32 = [8]opus_int32{9000, 700, 9000, 700, 11000, 1000, 12000, 2000}
var stereo_voice_threshold opus_int32 = 19000
var stereo_music_threshold opus_int32 = 17000
var mode_thresholds [2][2]opus_int32 = [2][2]opus_int32{{64000, 10000}, {44000, 10000}}
var fec_thresholds [10]opus_int32 = [10]opus_int32{12000, 1000, 14000, 1000, 16000, 1000, 20000, 1000, 22000, 1000}

func opus_encoder_get_size(channels int64) int64 {
	var (
		silkEncSizeBytes int64
		celtEncSizeBytes int64
		ret              int64
	)
	if channels < 1 || channels > 2 {
		return 0
	}
	ret = silk_Get_Encoder_Size(&silkEncSizeBytes)
	if ret != 0 {
		return 0
	}
	silkEncSizeBytes = align(silkEncSizeBytes)
	celtEncSizeBytes = celt_encoder_get_size(channels)
	return align(int64(unsafe.Sizeof(OpusEncoder{}))) + silkEncSizeBytes + celtEncSizeBytes
}
func opus_encoder_init(st *OpusEncoder, Fs opus_int32, channels int64, application int64) int64 {
	var (
		silk_enc         unsafe.Pointer
		celt_enc         *OpusCustomEncoder
		err              int64
		ret              int64
		silkEncSizeBytes int64
	)
	if Fs != 48000 && Fs != 24000 && Fs != 16000 && Fs != 12000 && Fs != 8000 || channels != 1 && channels != 2 || application != OPUS_APPLICATION_VOIP && application != OPUS_APPLICATION_AUDIO && application != OPUS_APPLICATION_RESTRICTED_LOWDELAY {
		return -1
	}
	libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(st))), 0, int(opus_encoder_get_size(channels)*int64(unsafe.Sizeof(byte(0)))))
	ret = silk_Get_Encoder_Size(&silkEncSizeBytes)
	if ret != 0 {
		return -1
	}
	silkEncSizeBytes = align(silkEncSizeBytes)
	st.Silk_enc_offset = align(int64(unsafe.Sizeof(OpusEncoder{})))
	st.Celt_enc_offset = st.Silk_enc_offset + silkEncSizeBytes
	silk_enc = unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_enc_offset)
	celt_enc = (*OpusCustomEncoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Celt_enc_offset))))
	st.Stream_channels = func() int64 {
		p := &st.Channels
		st.Channels = channels
		return *p
	}()
	st.Fs = Fs
	st.Arch = opus_select_arch()
	ret = silk_InitEncoder(silk_enc, st.Arch, &st.Silk_mode)
	if ret != 0 {
		return -3
	}
	st.Silk_mode.NChannelsAPI = opus_int32(channels)
	st.Silk_mode.NChannelsInternal = opus_int32(channels)
	st.Silk_mode.API_sampleRate = st.Fs
	st.Silk_mode.MaxInternalSampleRate = 16000
	st.Silk_mode.MinInternalSampleRate = 8000
	st.Silk_mode.DesiredInternalSampleRate = 16000
	st.Silk_mode.PayloadSize_ms = 20
	st.Silk_mode.BitRate = 25000
	st.Silk_mode.PacketLossPercentage = 0
	st.Silk_mode.Complexity = 9
	st.Silk_mode.UseInBandFEC = 0
	st.Silk_mode.UseDTX = 0
	st.Silk_mode.UseCBR = 0
	st.Silk_mode.ReducedDependency = 0
	err = celt_encoder_init(celt_enc, Fs, channels, st.Arch)
	if err != OPUS_OK {
		return -3
	}
	opus_custom_encoder_ctl(celt_enc, CELT_SET_SIGNALLING_REQUEST, func() int64 {
		0 == 0
		return 0
	}())
	opus_custom_encoder_ctl(celt_enc, OPUS_SET_COMPLEXITY_REQUEST, func() opus_int32 {
		st.Silk_mode.Complexity == 0
		return opus_int32(st.Silk_mode.Complexity)
	}())
	st.Use_vbr = 1
	st.Vbr_constraint = 1
	st.User_bitrate_bps = -1000
	st.Bitrate_bps = Fs*opus_int32(channels) + 3000
	st.Application = application
	st.Signal_type = -1000
	st.User_bandwidth = -1000
	st.Max_bandwidth = OPUS_BANDWIDTH_FULLBAND
	st.Force_channels = -1000
	st.User_forced_mode = -1000
	st.Voice_ratio = -1
	st.Encoder_buffer = int64(st.Fs / 100)
	st.Lsb_depth = 24
	st.Variable_duration = OPUS_FRAMESIZE_ARG
	st.Delay_compensation = int64(st.Fs / 250)
	st.Hybrid_stereo_width_Q14 = 1 << 14
	st.Prev_HB_gain = opus_val16(Q15ONE)
	st.Variable_HP_smth2_Q15 = opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)
	st.First = 1
	st.Mode = MODE_HYBRID
	st.Bandwidth = OPUS_BANDWIDTH_FULLBAND
	tonality_analysis_init(&st.Analysis, st.Fs)
	st.Analysis.Application = st.Application
	return OPUS_OK
}
func gen_toc(mode int64, framerate int64, bandwidth int64, channels int64) uint8 {
	var (
		period int64
		toc    uint8
	)
	period = 0
	for framerate < 400 {
		framerate <<= 1
		period++
	}
	if mode == MODE_SILK_ONLY {
		toc = uint8(int8((bandwidth - OPUS_BANDWIDTH_NARROWBAND) << 5))
		toc |= uint8(int8((period - 2) << 3))
	} else if mode == MODE_CELT_ONLY {
		var tmp int64 = bandwidth - OPUS_BANDWIDTH_MEDIUMBAND
		if tmp < 0 {
			tmp = 0
		}
		toc = 0x80
		toc |= uint8(int8(tmp << 5))
		toc |= uint8(int8(period << 3))
	} else {
		toc = 0x60
		toc |= uint8(int8((bandwidth - OPUS_BANDWIDTH_SUPERWIDEBAND) << 4))
		toc |= uint8(int8((period - 2) << 3))
	}
	toc |= uint8(int8(int64(libc.BoolToInt(channels == 2)) << 2))
	return toc
}
func silk_biquad_float(in *opus_val16, B_Q28 *opus_int32, A_Q28 *opus_int32, S *opus_val32, out *opus_val16, len_ opus_int32, stride int64) {
	var (
		k     int64
		vout  opus_val32
		inval opus_val32
		A     [2]opus_val32
		B     [3]opus_val32
	)
	A[0] = opus_val32(float64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*0))) * (1.0 / (1 << 28)))
	A[1] = opus_val32(float64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*1))) * (1.0 / (1 << 28)))
	B[0] = opus_val32(float64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*0))) * (1.0 / (1 << 28)))
	B[1] = opus_val32(float64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*1))) * (1.0 / (1 << 28)))
	B[2] = opus_val32(float64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*2))) * (1.0 / (1 << 28)))
	for k = 0; k < int64(len_); k++ {
		inval = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(k*stride))))
		vout = *(*opus_val32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_val32(0))*0)) + B[0]*inval
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_val32(0))*0)) = *(*opus_val32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_val32(0))*1)) - vout*A[0] + B[1]*inval
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_val32(0))*1)) = opus_val32(float64(-vout*A[1]+B[2]*inval) + VERY_SMALL)
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(k*stride))) = opus_val16(vout)
	}
}
func hp_cutoff(in *opus_val16, cutoff_Hz opus_int32, out *opus_val16, hp_mem *opus_val32, len_ int64, channels int64, Fs opus_int32, arch int64) {
	var (
		B_Q28  [3]opus_int32
		A_Q28  [2]opus_int32
		Fc_Q19 opus_int32
		r_Q28  opus_int32
		r_Q22  opus_int32
	)
	_ = arch
	Fc_Q19 = (opus_int32(opus_int16(cutoff_Hz)) * opus_int32(opus_int16(opus_int32((1.5*3.14159/1000)*(1<<19)+0.5)))) / (Fs / 1000)
	r_Q28 = (opus_int32(1.0*(1<<28) + 0.5)) - Fc_Q19*(opus_int32(0.92*(1<<9)+0.5))
	B_Q28[0] = r_Q28
	B_Q28[1] = opus_int32(opus_uint32(-r_Q28) << 1)
	B_Q28[2] = r_Q28
	r_Q22 = r_Q28 >> 6
	A_Q28[0] = opus_int32((int64(r_Q22) * int64((opus_int32((int64(Fc_Q19)*int64(Fc_Q19))>>16))-(opus_int32(2.0*(1<<22)+0.5)))) >> 16)
	A_Q28[1] = opus_int32((int64(r_Q22) * int64(r_Q22)) >> 16)
	silk_biquad_float(in, &B_Q28[0], &A_Q28[0], hp_mem, out, opus_int32(len_), channels)
	if channels == 2 {
		silk_biquad_float((*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*1)), &B_Q28[0], &A_Q28[0], (*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*2)), (*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*1)), opus_int32(len_), channels)
	}
}
func dc_reject(in *opus_val16, cutoff_Hz opus_int32, out *opus_val16, hp_mem *opus_val32, len_ int64, channels int64, Fs opus_int32) {
	var (
		i     int64
		coef  float32
		coef2 float32
	)
	coef = float32(float64(cutoff_Hz) * 6.3 / float64(Fs))
	coef2 = 1 - coef
	if channels == 2 {
		var (
			m0 float32
			m2 float32
		)
		m0 = float32(*(*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*0)))
		m2 = float32(*(*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*2)))
		for i = 0; i < len_; i++ {
			var (
				x0   opus_val32
				x1   opus_val32
				out0 opus_val32
				out1 opus_val32
			)
			x0 = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+0))))
			x1 = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+1))))
			out0 = x0 - opus_val32(m0)
			out1 = x1 - opus_val32(m2)
			m0 = float32(float64(coef*float32(x0)) + VERY_SMALL + float64(coef2*m0))
			m2 = float32(float64(coef*float32(x1)) + VERY_SMALL + float64(coef2*m2))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+0))) = opus_val16(out0)
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+1))) = opus_val16(out1)
		}
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*0)) = opus_val32(m0)
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*2)) = opus_val32(m2)
	} else {
		var m0 float32
		m0 = float32(*(*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*0)))
		for i = 0; i < len_; i++ {
			var (
				x opus_val32
				y opus_val32
			)
			x = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
			y = x - opus_val32(m0)
			m0 = float32(float64(coef*float32(x)) + VERY_SMALL + float64(coef2*m0))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(y)
		}
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(hp_mem), unsafe.Sizeof(opus_val32(0))*0)) = opus_val32(m0)
	}
}
func stereo_fade(in *opus_val16, out *opus_val16, g1 opus_val16, g2 opus_val16, overlap48 int64, frame_size int64, channels int64, window *opus_val16, Fs opus_int32) {
	var (
		i       int64
		overlap int64
		inc     int64
	)
	inc = int64(48000 / Fs)
	overlap = overlap48 / inc
	g1 = opus_val16(Q15ONE - float64(g1))
	g2 = opus_val16(Q15ONE - float64(g2))
	for i = 0; i < overlap; i++ {
		var (
			diff opus_val32
			g    opus_val16
			w    opus_val16
		)
		w = (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc)))) * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc))))
		g = opus_val16((opus_val32(w) * opus_val32(g2)) + opus_val32(Q15ONE-float64(w))*opus_val32(g1))
		diff = opus_val32(float64(opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels))))-opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+1))))) * 0.5)
		diff = opus_val32(g * opus_val16(diff))
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels))) - opus_val16(diff)
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+1))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+1))) + opus_val16(diff)
	}
	for ; i < frame_size; i++ {
		var diff opus_val32
		diff = opus_val32(float64(opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels))))-opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+1))))) * 0.5)
		diff = opus_val32(g2 * opus_val16(diff))
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels))) - opus_val16(diff)
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+1))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+1))) + opus_val16(diff)
	}
}
func gain_fade(in *opus_val16, out *opus_val16, g1 opus_val16, g2 opus_val16, overlap48 int64, frame_size int64, channels int64, window *opus_val16, Fs opus_int32) {
	var (
		i       int64
		inc     int64
		overlap int64
		c       int64
	)
	inc = int64(48000 / Fs)
	overlap = overlap48 / inc
	if channels == 1 {
		for i = 0; i < overlap; i++ {
			var (
				g opus_val16
				w opus_val16
			)
			w = (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc)))) * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc))))
			g = opus_val16((opus_val32(w) * opus_val32(g2)) + opus_val32(Q15ONE-float64(w))*opus_val32(g1))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = g * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		}
	} else {
		for i = 0; i < overlap; i++ {
			var (
				g opus_val16
				w opus_val16
			)
			w = (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc)))) * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc))))
			g = opus_val16((opus_val32(w) * opus_val32(g2)) + opus_val32(Q15ONE-float64(w))*opus_val32(g1))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*2))) = g * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*2))))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+1))) = g * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+1))))
		}
	}
	c = 0
	for {
		for i = overlap; i < frame_size; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+c))) = g2 * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+c))))
		}
		if func() int64 {
			p := &c
			*p++
			return *p
		}() >= channels {
			break
		}
	}
}
func opus_encoder_create(Fs opus_int32, channels int64, application int64, error *int64) *OpusEncoder {
	var (
		ret int64
		st  *OpusEncoder
	)
	if Fs != 48000 && Fs != 24000 && Fs != 16000 && Fs != 12000 && Fs != 8000 || channels != 1 && channels != 2 || application != OPUS_APPLICATION_VOIP && application != OPUS_APPLICATION_AUDIO && application != OPUS_APPLICATION_RESTRICTED_LOWDELAY {
		if error != nil {
			*error = -1
		}
		return nil
	}
	st = (*OpusEncoder)(libc.Malloc(int(opus_encoder_get_size(channels))))
	if st == nil {
		if error != nil {
			*error = -7
		}
		return nil
	}
	ret = opus_encoder_init(st, Fs, channels, application)
	if error != nil {
		*error = ret
	}
	if ret != OPUS_OK {
		libc.Free(unsafe.Pointer(st))
		st = nil
	}
	return st
}
func user_bitrate_to_bitrate(st *OpusEncoder, frame_size int64, max_data_bytes int64) opus_int32 {
	if frame_size == 0 {
		frame_size = int64(st.Fs / 400)
	}
	if st.User_bitrate_bps == -1000 {
		return st.Fs*60/opus_int32(frame_size) + st.Fs*opus_int32(st.Channels)
	} else if st.User_bitrate_bps == -1 {
		return opus_int32(max_data_bytes * 8 * int64(st.Fs) / frame_size)
	} else {
		return st.User_bitrate_bps
	}
}
func downmix_float(_x unsafe.Pointer, y *opus_val32, subframe int64, offset int64, c1 int64, c2 int64, C int64) {
	var (
		x *float32
		j int64
	)
	x = (*float32)(_x)
	for j = 0; j < subframe; j++ {
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val32(0))*uintptr(j))) = opus_val32(float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr((j+offset)*C+c1)))) * CELT_SIG_SCALE)
	}
	if c2 > -1 {
		for j = 0; j < subframe; j++ {
			*(*opus_val32)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val32(0))*uintptr(j))) += opus_val32(float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr((j+offset)*C+c2)))) * CELT_SIG_SCALE)
		}
	} else if c2 == -2 {
		var c int64
		for c = 1; c < C; c++ {
			for j = 0; j < subframe; j++ {
				*(*opus_val32)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val32(0))*uintptr(j))) += opus_val32(float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr((j+offset)*C+c)))) * CELT_SIG_SCALE)
			}
		}
	}
}
func downmix_int(_x unsafe.Pointer, y *opus_val32, subframe int64, offset int64, c1 int64, c2 int64, C int64) {
	var (
		x *opus_int16
		j int64
	)
	x = (*opus_int16)(_x)
	for j = 0; j < subframe; j++ {
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val32(0))*uintptr(j))) = opus_val32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr((j+offset)*C+c1))))
	}
	if c2 > -1 {
		for j = 0; j < subframe; j++ {
			*(*opus_val32)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val32(0))*uintptr(j))) += opus_val32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr((j+offset)*C+c2))))
		}
	} else if c2 == -2 {
		var c int64
		for c = 1; c < C; c++ {
			for j = 0; j < subframe; j++ {
				*(*opus_val32)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val32(0))*uintptr(j))) += opus_val32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr((j+offset)*C+c))))
			}
		}
	}
}
func frame_size_select(frame_size opus_int32, variable_duration int64, Fs opus_int32) opus_int32 {
	var new_size int64
	if frame_size < Fs/400 {
		return -1
	}
	if variable_duration == OPUS_FRAMESIZE_ARG {
		new_size = int64(frame_size)
	} else if variable_duration >= OPUS_FRAMESIZE_2_5_MS && variable_duration <= OPUS_FRAMESIZE_120_MS {
		if variable_duration <= OPUS_FRAMESIZE_40_MS {
			new_size = int64((Fs / 400) << opus_int32(variable_duration-OPUS_FRAMESIZE_2_5_MS))
		} else {
			new_size = (variable_duration - OPUS_FRAMESIZE_2_5_MS - 2) * int64(Fs) / 50
		}
	} else {
		return -1
	}
	if new_size > int64(frame_size) {
		return -1
	}
	if new_size*400 != int64(Fs) && new_size*200 != int64(Fs) && new_size*100 != int64(Fs) && new_size*50 != int64(Fs) && new_size*25 != int64(Fs) && new_size*50 != int64(Fs*3) && new_size*50 != int64(Fs*4) && new_size*50 != int64(Fs*5) && new_size*50 != int64(Fs*6) {
		return -1
	}
	return opus_int32(new_size)
}
func compute_stereo_width(pcm *opus_val16, frame_size int64, Fs opus_int32, mem *StereoWidthState) opus_val16 {
	var (
		xx          opus_val32
		xy          opus_val32
		yy          opus_val32
		sqrt_xx     opus_val16
		sqrt_yy     opus_val16
		qrrt_xx     opus_val16
		qrrt_yy     opus_val16
		frame_rate  int64
		i           int64
		short_alpha opus_val16
	)
	frame_rate = int64(Fs / opus_int32(frame_size))
	short_alpha = opus_val16(Q15ONE - float64((opus_val32(Q15ONE)*25)/opus_val32(func() int64 {
		if 50 > frame_rate {
			return 50
		}
		return frame_rate
	}())))
	xx = func() opus_val32 {
		xy = func() opus_val32 {
			yy = 0
			return yy
		}()
		return xy
	}()
	for i = 0; i < frame_size-3; i += 4 {
		var (
			pxx opus_val32 = 0
			pxy opus_val32 = 0
			pyy opus_val32 = 0
			x   opus_val16
			y   opus_val16
		)
		x = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2)))
		y = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+1)))
		pxx = opus_val32(x) * opus_val32(x)
		pxy = opus_val32(x) * opus_val32(y)
		pyy = opus_val32(y) * opus_val32(y)
		x = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+2)))
		y = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+3)))
		pxx += opus_val32(x) * opus_val32(x)
		pxy += opus_val32(x) * opus_val32(y)
		pyy += opus_val32(y) * opus_val32(y)
		x = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+4)))
		y = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+5)))
		pxx += opus_val32(x) * opus_val32(x)
		pxy += opus_val32(x) * opus_val32(y)
		pyy += opus_val32(y) * opus_val32(y)
		x = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+6)))
		y = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*2+7)))
		pxx += opus_val32(x) * opus_val32(x)
		pxy += opus_val32(x) * opus_val32(y)
		pyy += opus_val32(y) * opus_val32(y)
		xx += pxx
		xy += pxy
		yy += pyy
	}
	if float64(xx) >= 1e+09 || xx != xx || float64(yy) >= 1e+09 || yy != yy {
		xy = func() opus_val32 {
			xx = func() opus_val32 {
				yy = 0
				return yy
			}()
			return xx
		}()
	}
	mem.XX += opus_val32(short_alpha * opus_val16(xx-mem.XX))
	mem.XY += opus_val32(short_alpha * opus_val16(xy-mem.XY))
	mem.YY += opus_val32(short_alpha * opus_val16(yy-mem.YY))
	if 0 > mem.XX {
		mem.XX = 0
	} else {
		mem.XX = mem.XX
	}
	if 0 > mem.XY {
		mem.XY = 0
	} else {
		mem.XY = mem.XY
	}
	if 0 > mem.YY {
		mem.YY = 0
	} else {
		mem.YY = mem.YY
	}
	if float64(func() opus_val32 {
		if mem.XX > mem.YY {
			return mem.XX
		}
		return mem.YY
	}()) > 0.0008 {
		var (
			corr  opus_val16
			ldiff opus_val16
			width opus_val16
		)
		sqrt_xx = opus_val16(float32(math.Sqrt(float64(mem.XX))))
		sqrt_yy = opus_val16(float32(math.Sqrt(float64(mem.YY))))
		qrrt_xx = opus_val16(float32(math.Sqrt(float64(sqrt_xx))))
		qrrt_yy = opus_val16(float32(math.Sqrt(float64(sqrt_yy))))
		if mem.XY < opus_val32(sqrt_xx*sqrt_yy) {
			mem.XY = mem.XY
		} else {
			mem.XY = opus_val32(sqrt_xx * sqrt_yy)
		}
		corr = opus_val16(float64(float32(mem.XY)) / (EPSILON + float64(opus_val32(sqrt_xx)*opus_val32(sqrt_yy))))
		ldiff = opus_val16(float64(opus_val32(Q15ONE)*opus_val32(float32(math.Abs(float64(qrrt_xx-qrrt_yy))))) / (EPSILON + float64(qrrt_xx) + float64(qrrt_yy)))
		width = opus_val16((float32(math.Sqrt(1.0 - float64(opus_val32(corr)*opus_val32(corr))))) * float32(ldiff))
		mem.Smoothed_width += (width - mem.Smoothed_width) / opus_val16(frame_rate)
		if (float64(mem.Max_follower) - 0.02/float64(frame_rate)) > float64(mem.Smoothed_width) {
			mem.Max_follower = opus_val16(float64(mem.Max_follower) - 0.02/float64(frame_rate))
		} else {
			mem.Max_follower = mem.Smoothed_width
		}
	}
	if Q15ONE < float64(opus_val32(mem.Max_follower)*20) {
		return opus_val16(Q15ONE)
	}
	return opus_val16(opus_val32(mem.Max_follower) * 20)
}
func decide_fec(useInBandFEC int64, PacketLoss_perc int64, last_fec int64, mode int64, bandwidth *int64, rate opus_int32) int64 {
	var orig_bandwidth int64
	if useInBandFEC == 0 || PacketLoss_perc == 0 || mode == MODE_CELT_ONLY {
		return 0
	}
	orig_bandwidth = *bandwidth
	for {
		var (
			hysteresis          opus_int32
			LBRR_rate_thres_bps opus_int32
		)
		LBRR_rate_thres_bps = fec_thresholds[(*bandwidth-OPUS_BANDWIDTH_NARROWBAND)*2]
		hysteresis = fec_thresholds[(*bandwidth-OPUS_BANDWIDTH_NARROWBAND)*2+1]
		if last_fec == 1 {
			LBRR_rate_thres_bps -= hysteresis
		}
		if last_fec == 0 {
			LBRR_rate_thres_bps += hysteresis
		}
		LBRR_rate_thres_bps = ((LBRR_rate_thres_bps * opus_int32(125-(func() int64 {
			if PacketLoss_perc < 25 {
				return PacketLoss_perc
			}
			return 25
		}()))) * opus_int32(int64(opus_int16(opus_int32(0.01*(1<<16)+0.5))))) >> 16
		if rate > LBRR_rate_thres_bps {
			return 1
		} else if PacketLoss_perc <= 5 {
			return 0
		} else if *bandwidth > OPUS_BANDWIDTH_NARROWBAND {
			(*bandwidth)--
		} else {
			break
		}
	}
	*bandwidth = orig_bandwidth
	return 0
}
func compute_silk_rate_for_hybrid(rate int64, bandwidth int64, frame20ms int64, vbr int64, fec int64, channels int64) int64 {
	var (
		entry      int64
		i          int64
		N          int64
		silk_rate  int64
		rate_table [7][5]int64 = [7][5]int64{{}, {12000, 10000, 10000, 11000, 11000}, {16000, 13500, 13500, 15000, 15000}, {20000, 16000, 16000, 18000, 18000}, {24000, 18000, 18000, 21000, 21000}, {32000, 22000, 22000, 28000, 28000}, {64000, 38000, 38000, 50000, 50000}}
	)
	rate /= channels
	entry = frame20ms + 1 + fec*2
	N = int64(unsafe.Sizeof([7][5]int64{}) / unsafe.Sizeof([5]int64{}))
	for i = 1; i < N; i++ {
		if rate_table[i][0] > rate {
			break
		}
	}
	if i == N {
		silk_rate = rate_table[i-1][entry]
		silk_rate += (rate - rate_table[i-1][0]) / 2
	} else {
		var (
			lo opus_int32
			hi opus_int32
			x0 opus_int32
			x1 opus_int32
		)
		lo = opus_int32(rate_table[i-1][entry])
		hi = opus_int32(rate_table[i][entry])
		x0 = opus_int32(rate_table[i-1][0])
		x1 = opus_int32(rate_table[i][0])
		silk_rate = int64((lo*(x1-opus_int32(rate)) + hi*opus_int32(rate-int64(x0))) / (x1 - x0))
	}
	if vbr == 0 {
		silk_rate += 100
	}
	if bandwidth == OPUS_BANDWIDTH_SUPERWIDEBAND {
		silk_rate += 300
	}
	silk_rate *= channels
	if channels == 2 && rate >= 12000 {
		silk_rate -= 1000
	}
	return silk_rate
}
func compute_equiv_rate(bitrate opus_int32, channels int64, frame_rate int64, vbr int64, mode int64, complexity int64, loss int64) opus_int32 {
	var equiv opus_int32
	equiv = bitrate
	if frame_rate > 50 {
		equiv -= opus_int32((channels*40 + 20) * (frame_rate - 50))
	}
	if vbr == 0 {
		equiv -= equiv / 12
	}
	equiv = equiv * opus_int32(complexity+90) / 100
	if mode == MODE_SILK_ONLY || mode == MODE_HYBRID {
		if complexity < 2 {
			equiv = equiv * 4 / 5
		}
		equiv -= equiv * opus_int32(loss) / opus_int32(loss*6+10)
	} else if mode == MODE_CELT_ONLY {
		if complexity < 5 {
			equiv = equiv * 9 / 10
		}
	} else {
		equiv -= equiv * opus_int32(loss) / opus_int32(loss*12+20)
	}
	return equiv
}
func is_digital_silence(pcm *opus_val16, frame_size int64, channels int64, lsb_depth int64) int64 {
	var (
		silence    int64      = 0
		sample_max opus_val32 = 0
	)
	sample_max = celt_maxabs16(pcm, frame_size*channels)
	silence = int64(libc.BoolToInt(sample_max <= opus_val32(1/(1<<lsb_depth))))
	return silence
}
func compute_frame_energy(pcm *opus_val16, frame_size int64, channels int64, arch int64) opus_val32 {
	var len_ int64 = frame_size * channels
	return (func() opus_val32 {
		_ = arch
		return celt_inner_prod_c(pcm, pcm, len_)
	}()) / opus_val32(len_)
}
func decide_dtx_mode(activity int64, nb_no_activity_ms_Q1 *int64, frame_size_ms_Q1 int64) int64 {
	if activity == 0 {
		*nb_no_activity_ms_Q1 += frame_size_ms_Q1
		if *nb_no_activity_ms_Q1 > NB_SPEECH_FRAMES_BEFORE_DTX*20*2 {
			if *nb_no_activity_ms_Q1 <= (NB_SPEECH_FRAMES_BEFORE_DTX+MAX_CONSECUTIVE_DTX)*20*2 {
				return 1
			} else {
				*nb_no_activity_ms_Q1 = NB_SPEECH_FRAMES_BEFORE_DTX * 20 * 2
			}
		}
	} else {
		*nb_no_activity_ms_Q1 = 0
	}
	return 0
}
func encode_multiframe_packet(st *OpusEncoder, pcm *opus_val16, nb_frames int64, frame_size int64, data *uint8, out_data_bytes opus_int32, to_celt int64, lsb_depth int64, float_api int64) opus_int32 {
	var (
		i                int64
		ret              int64 = 0
		tmp_data         *uint8
		bak_mode         int64
		bak_bandwidth    int64
		bak_channels     int64
		bak_to_mono      int64
		rp               *OpusRepacketizer
		max_header_bytes int64
		bytes_per_frame  opus_int32
		cbr_bytes        opus_int32
		repacketize_len  opus_int32
		tmp_len          int64
	)
	if nb_frames == 2 {
		max_header_bytes = 3
	} else {
		max_header_bytes = (nb_frames-1)*2 + 2
	}
	if st.Use_vbr != 0 || st.User_bitrate_bps == -1 {
		repacketize_len = out_data_bytes
	} else {
		cbr_bytes = st.Bitrate_bps * 3 / (st.Fs * (3 * 8) / opus_int32(frame_size*nb_frames))
		if cbr_bytes < out_data_bytes {
			repacketize_len = cbr_bytes
		} else {
			repacketize_len = out_data_bytes
		}
	}
	if 1276 < ((repacketize_len-opus_int32(max_header_bytes))/opus_int32(nb_frames) + 1) {
		bytes_per_frame = 1276
	} else {
		bytes_per_frame = (repacketize_len-opus_int32(max_header_bytes))/opus_int32(nb_frames) + 1
	}
	tmp_data = (*uint8)(libc.Malloc(int((nb_frames * int64(bytes_per_frame)) * int64(unsafe.Sizeof(uint8(0))))))
	rp = (*OpusRepacketizer)(libc.Malloc(int(unsafe.Sizeof(OpusRepacketizer{}) * 1)))
	opus_repacketizer_init(rp)
	bak_mode = st.User_forced_mode
	bak_bandwidth = st.User_bandwidth
	bak_channels = st.Force_channels
	st.User_forced_mode = st.Mode
	st.User_bandwidth = st.Bandwidth
	st.Force_channels = st.Stream_channels
	bak_to_mono = st.Silk_mode.ToMono
	if bak_to_mono != 0 {
		st.Force_channels = 1
	} else {
		st.Prev_channels = st.Stream_channels
	}
	for i = 0; i < nb_frames; i++ {
		st.Silk_mode.ToMono = 0
		st.Nonfinal_frame = int64(libc.BoolToInt(i < (nb_frames - 1)))
		if to_celt != 0 && i == nb_frames-1 {
			st.User_forced_mode = MODE_CELT_ONLY
		}
		tmp_len = int64(opus_encode_native(st, (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i*(st.Channels*frame_size)))), frame_size, (*uint8)(unsafe.Add(unsafe.Pointer(tmp_data), i*int64(bytes_per_frame))), bytes_per_frame, lsb_depth, nil, 0, 0, 0, 0, nil, float_api))
		if tmp_len < 0 {
			return -3
		}
		ret = opus_repacketizer_cat(rp, (*uint8)(unsafe.Add(unsafe.Pointer(tmp_data), i*int64(bytes_per_frame))), opus_int32(tmp_len))
		if ret < 0 {
			return -3
		}
	}
	ret = int64(opus_repacketizer_out_range_impl(rp, 0, nb_frames, data, repacketize_len, 0, int64(libc.BoolToInt(st.Use_vbr == 0))))
	if ret < 0 {
		return -3
	}
	st.User_forced_mode = bak_mode
	st.User_bandwidth = bak_bandwidth
	st.Force_channels = bak_channels
	st.Silk_mode.ToMono = bak_to_mono
	return opus_int32(ret)
}
func compute_redundancy_bytes(max_data_bytes opus_int32, bitrate_bps opus_int32, frame_rate int64, channels int64) int64 {
	var (
		redundancy_bytes_cap int64
		redundancy_bytes     int64
		redundancy_rate      opus_int32
		base_bits            int64
		available_bits       opus_int32
	)
	base_bits = channels*40 + 20
	redundancy_rate = bitrate_bps + opus_int32(base_bits*(200-frame_rate))
	redundancy_rate = redundancy_rate * 3 / 2
	redundancy_bytes = int64(redundancy_rate / 1600)
	available_bits = max_data_bytes*8 - opus_int32(base_bits*2)
	redundancy_bytes_cap = int64((available_bits*240/opus_int32(48000/frame_rate+240) + opus_int32(base_bits)) / 8)
	if redundancy_bytes < redundancy_bytes_cap {
		redundancy_bytes = redundancy_bytes
	} else {
		redundancy_bytes = redundancy_bytes_cap
	}
	if redundancy_bytes > channels*8+4 {
		if 257 < redundancy_bytes {
			redundancy_bytes = 257
		} else {
			redundancy_bytes = redundancy_bytes
		}
	} else {
		redundancy_bytes = 0
	}
	return redundancy_bytes
}
func opus_encode_native(st *OpusEncoder, pcm *opus_val16, frame_size int64, data *uint8, out_data_bytes opus_int32, lsb_depth int64, analysis_pcm unsafe.Pointer, analysis_size opus_int32, c1 int64, c2 int64, analysis_channels int64, downmix downmix_func, float_api int64) opus_int32 {
	var (
		silk_enc                   unsafe.Pointer
		celt_enc                   *OpusCustomEncoder
		i                          int64
		ret                        int64 = 0
		nBytes                     opus_int32
		enc                        ec_enc
		bytes_target               int64
		prefill                    int64 = 0
		start_band                 int64 = 0
		redundancy                 int64 = 0
		redundancy_bytes           int64 = 0
		celt_to_silk               int64 = 0
		pcm_buf                    *opus_val16
		nb_compr_bytes             int64
		to_celt                    int64       = 0
		redundant_rng              opus_uint32 = 0
		cutoff_Hz                  int64
		hp_freq_smth1              int64
		voice_est                  int64
		equiv_rate                 opus_int32
		delay_compensation         int64
		frame_rate                 int64
		max_rate                   opus_int32
		curr_bandwidth             int64
		HB_gain                    opus_val16
		max_data_bytes             opus_int32
		total_buffer               int64
		stereo_width               opus_val16
		celt_mode                  *OpusCustomMode
		analysis_info              AnalysisInfo
		analysis_read_pos_bak      int64 = -1
		analysis_read_subframe_bak int64 = -1
		is_silence                 int64 = 0
		activity                   int64 = -1
		tmp_prefill                *opus_val16
	)
	if 1276 < out_data_bytes {
		max_data_bytes = 1276
	} else {
		max_data_bytes = out_data_bytes
	}
	st.RangeFinal = 0
	if frame_size <= 0 || max_data_bytes <= 0 {
		return -1
	}
	if max_data_bytes == 1 && st.Fs == opus_int32(frame_size*10) {
		return -2
	}
	silk_enc = unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_enc_offset)
	celt_enc = (*OpusCustomEncoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Celt_enc_offset))))
	if st.Application == OPUS_APPLICATION_RESTRICTED_LOWDELAY {
		delay_compensation = 0
	} else {
		delay_compensation = st.Delay_compensation
	}
	if lsb_depth < st.Lsb_depth {
		lsb_depth = lsb_depth
	} else {
		lsb_depth = st.Lsb_depth
	}
	opus_custom_encoder_ctl(celt_enc, CELT_GET_MODE_REQUEST, (**OpusCustomMode)(unsafe.Add(unsafe.Pointer(&celt_mode), unsafe.Sizeof((*OpusCustomMode)(nil))*uintptr(int64(uintptr(unsafe.Pointer(&celt_mode))-uintptr(unsafe.Pointer(&celt_mode)))))))
	analysis_info.Valid = 0
	if st.Silk_mode.Complexity >= 7 && st.Fs >= 16000 {
		is_silence = is_digital_silence(pcm, frame_size, st.Channels, lsb_depth)
		analysis_read_pos_bak = st.Analysis.Read_pos
		analysis_read_subframe_bak = st.Analysis.Read_subframe
		run_analysis(&st.Analysis, celt_mode, analysis_pcm, int64(analysis_size), frame_size, c1, c2, analysis_channels, st.Fs, lsb_depth, downmix, &analysis_info)
		if is_silence == 0 && float64(analysis_info.Activity_probability) > DTX_ACTIVITY_THRESHOLD {
			if (float64(st.Peak_signal_energy) * 0.999) > float64(compute_frame_energy(pcm, frame_size, st.Channels, st.Arch)) {
				st.Peak_signal_energy = opus_val32(float64(st.Peak_signal_energy) * 0.999)
			} else {
				st.Peak_signal_energy = compute_frame_energy(pcm, frame_size, st.Channels, st.Arch)
			}
		}
	} else if st.Analysis.Initialized != 0 {
		tonality_analysis_reset(&st.Analysis)
	}
	if is_silence == 0 {
		st.Voice_ratio = -1
	}
	if is_silence != 0 {
		activity = int64(libc.BoolToInt(is_silence == 0))
	} else if analysis_info.Valid != 0 {
		activity = int64(libc.BoolToInt(float64(analysis_info.Activity_probability) >= DTX_ACTIVITY_THRESHOLD))
		if activity == 0 {
			var noise_energy opus_val32 = compute_frame_energy(pcm, frame_size, st.Channels, st.Arch)
			activity = int64(libc.BoolToInt(float64(st.Peak_signal_energy) < (PSEUDO_SNR_THRESHOLD * float64(noise_energy))))
		}
	}
	st.Detected_bandwidth = 0
	if analysis_info.Valid != 0 {
		var analysis_bandwidth int64
		if st.Signal_type == -1000 {
			var prob float32
			if st.Prev_mode == 0 {
				prob = analysis_info.Music_prob
			} else if st.Prev_mode == MODE_CELT_ONLY {
				prob = analysis_info.Music_prob_max
			} else {
				prob = analysis_info.Music_prob_min
			}
			st.Voice_ratio = int64(math.Floor(float64((1-prob)*100) + 0.5))
		}
		analysis_bandwidth = analysis_info.Bandwidth
		if analysis_bandwidth <= 12 {
			st.Detected_bandwidth = OPUS_BANDWIDTH_NARROWBAND
		} else if analysis_bandwidth <= 14 {
			st.Detected_bandwidth = OPUS_BANDWIDTH_MEDIUMBAND
		} else if analysis_bandwidth <= 16 {
			st.Detected_bandwidth = OPUS_BANDWIDTH_WIDEBAND
		} else if analysis_bandwidth <= 18 {
			st.Detected_bandwidth = OPUS_BANDWIDTH_SUPERWIDEBAND
		} else {
			st.Detected_bandwidth = OPUS_BANDWIDTH_FULLBAND
		}
	}
	if st.Channels == 2 && st.Force_channels != 1 {
		stereo_width = compute_stereo_width(pcm, frame_size, st.Fs, &st.Width_mem)
	} else {
		stereo_width = 0
	}
	total_buffer = delay_compensation
	st.Bitrate_bps = user_bitrate_to_bitrate(st, frame_size, int64(max_data_bytes))
	frame_rate = int64(st.Fs / opus_int32(frame_size))
	if st.Use_vbr == 0 {
		var (
			cbrBytes     int64
			frame_rate12 int64 = int64(st.Fs * 12 / opus_int32(frame_size))
		)
		if ((st.Bitrate_bps*12/8 + opus_int32(frame_rate12/2)) / opus_int32(frame_rate12)) < max_data_bytes {
			cbrBytes = int64((st.Bitrate_bps*12/8 + opus_int32(frame_rate12/2)) / opus_int32(frame_rate12))
		} else {
			cbrBytes = int64(max_data_bytes)
		}
		st.Bitrate_bps = opus_int32(cbrBytes * int64(opus_int32(frame_rate12)) * 8 / 12)
		if 1 > cbrBytes {
			max_data_bytes = 1
		} else {
			max_data_bytes = opus_int32(cbrBytes)
		}
	}
	if max_data_bytes < 3 || st.Bitrate_bps < opus_int32(frame_rate*3*8) || frame_rate < 50 && (max_data_bytes*opus_int32(frame_rate) < 300 || st.Bitrate_bps < 2400) {
		var (
			tocmode int64 = st.Mode
			bw      int64
		)
		if st.Bandwidth == 0 {
			bw = OPUS_BANDWIDTH_NARROWBAND
		} else {
			bw = st.Bandwidth
		}
		var packet_code int64 = 0
		var num_multiframes int64 = 0
		if tocmode == 0 {
			tocmode = MODE_SILK_ONLY
		}
		if frame_rate > 100 {
			tocmode = MODE_CELT_ONLY
		}
		if frame_rate == 25 && tocmode != MODE_SILK_ONLY {
			frame_rate = 50
			packet_code = 1
		}
		if frame_rate <= 16 {
			if out_data_bytes == 1 || tocmode == MODE_SILK_ONLY && frame_rate != 10 {
				tocmode = MODE_SILK_ONLY
				packet_code = int64(libc.BoolToInt(frame_rate <= 12))
				if frame_rate == 12 {
					frame_rate = 25
				} else {
					frame_rate = 16
				}
			} else {
				num_multiframes = 50 / frame_rate
				frame_rate = 50
				packet_code = 3
			}
		}
		if tocmode == MODE_SILK_ONLY && bw > OPUS_BANDWIDTH_WIDEBAND {
			bw = OPUS_BANDWIDTH_WIDEBAND
		} else if tocmode == MODE_CELT_ONLY && bw == OPUS_BANDWIDTH_MEDIUMBAND {
			bw = OPUS_BANDWIDTH_NARROWBAND
		} else if tocmode == MODE_HYBRID && bw <= OPUS_BANDWIDTH_SUPERWIDEBAND {
			bw = OPUS_BANDWIDTH_SUPERWIDEBAND
		}
		*data = gen_toc(tocmode, frame_rate, bw, st.Stream_channels)
		*data |= uint8(int8(packet_code))
		if packet_code <= 1 {
			ret = 1
		} else {
			ret = 2
		}
		if max_data_bytes > opus_int32(ret) {
			max_data_bytes = max_data_bytes
		} else {
			max_data_bytes = opus_int32(ret)
		}
		if packet_code == 3 {
			*(*uint8)(unsafe.Add(unsafe.Pointer(data), 1)) = uint8(int8(num_multiframes))
		}
		if st.Use_vbr == 0 {
			ret = opus_packet_pad(data, opus_int32(ret), max_data_bytes)
			if ret == OPUS_OK {
				ret = int64(max_data_bytes)
			} else {
				ret = -3
			}
		}
		return opus_int32(ret)
	}
	max_rate = opus_int32(frame_rate * int64(max_data_bytes) * 8)
	equiv_rate = compute_equiv_rate(st.Bitrate_bps, st.Channels, int64(st.Fs/opus_int32(frame_size)), st.Use_vbr, 0, st.Silk_mode.Complexity, st.Silk_mode.PacketLossPercentage)
	if st.Signal_type == OPUS_SIGNAL_VOICE {
		voice_est = math.MaxInt8
	} else if st.Signal_type == OPUS_SIGNAL_MUSIC {
		voice_est = 0
	} else if st.Voice_ratio >= 0 {
		voice_est = st.Voice_ratio * 327 >> 8
		if st.Application == OPUS_APPLICATION_AUDIO {
			if voice_est < 115 {
				voice_est = voice_est
			} else {
				voice_est = 115
			}
		}
	} else if st.Application == OPUS_APPLICATION_VOIP {
		voice_est = 115
	} else {
		voice_est = 48
	}
	if st.Force_channels != -1000 && st.Channels == 2 {
		st.Stream_channels = st.Force_channels
	} else {
		if st.Channels == 2 {
			var stereo_threshold opus_int32
			stereo_threshold = stereo_music_threshold + opus_int32((voice_est*voice_est*int64(stereo_voice_threshold-stereo_music_threshold))>>14)
			if st.Stream_channels == 2 {
				stereo_threshold -= 1000
			} else {
				stereo_threshold += 1000
			}
			if equiv_rate > stereo_threshold {
				st.Stream_channels = 2
			} else {
				st.Stream_channels = 1
			}
		} else {
			st.Stream_channels = st.Channels
		}
	}
	equiv_rate = compute_equiv_rate(st.Bitrate_bps, st.Stream_channels, int64(st.Fs/opus_int32(frame_size)), st.Use_vbr, 0, st.Silk_mode.Complexity, st.Silk_mode.PacketLossPercentage)
	st.Silk_mode.UseDTX = int64(libc.BoolToInt(st.Use_dtx != 0 && (analysis_info.Valid == 0 && is_silence == 0)))
	if st.Application == OPUS_APPLICATION_RESTRICTED_LOWDELAY {
		st.Mode = MODE_CELT_ONLY
	} else if st.User_forced_mode == -1000 {
		var (
			mode_voice opus_int32
			mode_music opus_int32
			threshold  opus_int32
		)
		mode_voice = opus_int32(((Q15ONE - float64(stereo_width)) * float64(mode_thresholds[0][0])) + float64(stereo_width*opus_val16(mode_thresholds[1][0])))
		mode_music = opus_int32(((Q15ONE - float64(stereo_width)) * float64(mode_thresholds[1][1])) + float64(stereo_width*opus_val16(mode_thresholds[1][1])))
		threshold = mode_music + opus_int32((voice_est*voice_est*int64(mode_voice-mode_music))>>14)
		if st.Application == OPUS_APPLICATION_VOIP {
			threshold += 8000
		}
		if st.Prev_mode == MODE_CELT_ONLY {
			threshold -= 4000
		} else if st.Prev_mode > 0 {
			threshold += 4000
		}
		if equiv_rate >= threshold {
			st.Mode = MODE_CELT_ONLY
		} else {
			st.Mode = MODE_SILK_ONLY
		}
		if st.Silk_mode.UseInBandFEC != 0 && st.Silk_mode.PacketLossPercentage > (128-voice_est)>>4 && (st.Fec_config != 2 || voice_est > 25) {
			st.Mode = MODE_SILK_ONLY
		}
		if st.Silk_mode.UseDTX != 0 && voice_est > 100 {
			st.Mode = MODE_SILK_ONLY
		}
		if max_data_bytes < opus_int32((func() int64 {
			if frame_rate > 50 {
				return 9000
			}
			return 6000
		}())*frame_size/int64(st.Fs*8)) {
			st.Mode = MODE_CELT_ONLY
		}
	} else {
		st.Mode = st.User_forced_mode
	}
	if st.Mode != MODE_CELT_ONLY && frame_size < int64(st.Fs/100) {
		st.Mode = MODE_CELT_ONLY
	}
	if st.Lfe != 0 {
		st.Mode = MODE_CELT_ONLY
	}
	if st.Prev_mode > 0 && (st.Mode != MODE_CELT_ONLY && st.Prev_mode == MODE_CELT_ONLY || st.Mode == MODE_CELT_ONLY && st.Prev_mode != MODE_CELT_ONLY) {
		redundancy = 1
		celt_to_silk = int64(libc.BoolToInt(st.Mode != MODE_CELT_ONLY))
		if celt_to_silk == 0 {
			if frame_size >= int64(st.Fs/100) {
				st.Mode = st.Prev_mode
				to_celt = 1
			} else {
				redundancy = 0
			}
		}
	}
	if st.Stream_channels == 1 && st.Prev_channels == 2 && st.Silk_mode.ToMono == 0 && st.Mode != MODE_CELT_ONLY && st.Prev_mode != MODE_CELT_ONLY {
		st.Silk_mode.ToMono = 1
		st.Stream_channels = 2
	} else {
		st.Silk_mode.ToMono = 0
	}
	equiv_rate = compute_equiv_rate(st.Bitrate_bps, st.Stream_channels, int64(st.Fs/opus_int32(frame_size)), st.Use_vbr, st.Mode, st.Silk_mode.Complexity, st.Silk_mode.PacketLossPercentage)
	if st.Mode != MODE_CELT_ONLY && st.Prev_mode == MODE_CELT_ONLY {
		var dummy silk_EncControlStruct
		silk_InitEncoder(silk_enc, st.Arch, &dummy)
		prefill = 1
	}
	if st.Mode == MODE_CELT_ONLY || st.First != 0 || st.Silk_mode.AllowBandwidthSwitch != 0 {
		var (
			voice_bandwidth_thresholds *opus_int32
			music_bandwidth_thresholds *opus_int32
			bandwidth_thresholds       [8]opus_int32
			bandwidth                  int64 = OPUS_BANDWIDTH_FULLBAND
		)
		if st.Channels == 2 && st.Force_channels != 1 {
			voice_bandwidth_thresholds = &stereo_voice_bandwidth_thresholds[0]
			music_bandwidth_thresholds = &stereo_music_bandwidth_thresholds[0]
		} else {
			voice_bandwidth_thresholds = &mono_voice_bandwidth_thresholds[0]
			music_bandwidth_thresholds = &mono_music_bandwidth_thresholds[0]
		}
		for i = 0; i < 8; i++ {
			bandwidth_thresholds[i] = *(*opus_int32)(unsafe.Add(unsafe.Pointer(music_bandwidth_thresholds), unsafe.Sizeof(opus_int32(0))*uintptr(i))) + opus_int32((voice_est*voice_est*int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(voice_bandwidth_thresholds), unsafe.Sizeof(opus_int32(0))*uintptr(i)))-*(*opus_int32)(unsafe.Add(unsafe.Pointer(music_bandwidth_thresholds), unsafe.Sizeof(opus_int32(0))*uintptr(i)))))>>14)
		}
		for {
			{
				var (
					threshold  int64
					hysteresis int64
				)
				threshold = int64(bandwidth_thresholds[(bandwidth-OPUS_BANDWIDTH_MEDIUMBAND)*2])
				hysteresis = int64(bandwidth_thresholds[(bandwidth-OPUS_BANDWIDTH_MEDIUMBAND)*2+1])
				if st.First == 0 {
					if st.Auto_bandwidth >= bandwidth {
						threshold -= hysteresis
					} else {
						threshold += hysteresis
					}
				}
				if equiv_rate >= opus_int32(threshold) {
					break
				}
			}
			if func() int64 {
				p := &bandwidth
				*p--
				return *p
			}() <= OPUS_BANDWIDTH_NARROWBAND {
				break
			}
		}
		if bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
			bandwidth = OPUS_BANDWIDTH_WIDEBAND
		}
		st.Bandwidth = func() int64 {
			p := &st.Auto_bandwidth
			st.Auto_bandwidth = bandwidth
			return *p
		}()
		if st.First == 0 && st.Mode != MODE_CELT_ONLY && st.Silk_mode.InWBmodeWithoutVariableLP == 0 && st.Bandwidth > OPUS_BANDWIDTH_WIDEBAND {
			st.Bandwidth = OPUS_BANDWIDTH_WIDEBAND
		}
	}
	if st.Bandwidth > st.Max_bandwidth {
		st.Bandwidth = st.Max_bandwidth
	}
	if st.User_bandwidth != -1000 {
		st.Bandwidth = st.User_bandwidth
	}
	if st.Mode != MODE_CELT_ONLY && max_rate < 15000 {
		if st.Bandwidth < OPUS_BANDWIDTH_WIDEBAND {
			st.Bandwidth = st.Bandwidth
		} else {
			st.Bandwidth = OPUS_BANDWIDTH_WIDEBAND
		}
	}
	if st.Fs <= 24000 && st.Bandwidth > OPUS_BANDWIDTH_SUPERWIDEBAND {
		st.Bandwidth = OPUS_BANDWIDTH_SUPERWIDEBAND
	}
	if st.Fs <= 16000 && st.Bandwidth > OPUS_BANDWIDTH_WIDEBAND {
		st.Bandwidth = OPUS_BANDWIDTH_WIDEBAND
	}
	if st.Fs <= 12000 && st.Bandwidth > OPUS_BANDWIDTH_MEDIUMBAND {
		st.Bandwidth = OPUS_BANDWIDTH_MEDIUMBAND
	}
	if st.Fs <= 8000 && st.Bandwidth > OPUS_BANDWIDTH_NARROWBAND {
		st.Bandwidth = OPUS_BANDWIDTH_NARROWBAND
	}
	if st.Detected_bandwidth != 0 && st.User_bandwidth == -1000 {
		var min_detected_bandwidth int64
		if equiv_rate <= opus_int32(st.Stream_channels*18000) && st.Mode == MODE_CELT_ONLY {
			min_detected_bandwidth = OPUS_BANDWIDTH_NARROWBAND
		} else if equiv_rate <= opus_int32(st.Stream_channels*24000) && st.Mode == MODE_CELT_ONLY {
			min_detected_bandwidth = OPUS_BANDWIDTH_MEDIUMBAND
		} else if equiv_rate <= opus_int32(st.Stream_channels*30000) {
			min_detected_bandwidth = OPUS_BANDWIDTH_WIDEBAND
		} else if equiv_rate <= opus_int32(st.Stream_channels*44000) {
			min_detected_bandwidth = OPUS_BANDWIDTH_SUPERWIDEBAND
		} else {
			min_detected_bandwidth = OPUS_BANDWIDTH_FULLBAND
		}
		if st.Detected_bandwidth > min_detected_bandwidth {
			st.Detected_bandwidth = st.Detected_bandwidth
		} else {
			st.Detected_bandwidth = min_detected_bandwidth
		}
		if st.Bandwidth < st.Detected_bandwidth {
			st.Bandwidth = st.Bandwidth
		} else {
			st.Bandwidth = st.Detected_bandwidth
		}
	}
	st.Silk_mode.LBRR_coded = decide_fec(st.Silk_mode.UseInBandFEC, st.Silk_mode.PacketLossPercentage, st.Silk_mode.LBRR_coded, st.Mode, &st.Bandwidth, equiv_rate)
	opus_custom_encoder_ctl(celt_enc, OPUS_SET_LSB_DEPTH_REQUEST, func() opus_int32 {
		lsb_depth == 0
		return opus_int32(lsb_depth)
	}())
	if st.Mode == MODE_CELT_ONLY && st.Bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
		st.Bandwidth = OPUS_BANDWIDTH_WIDEBAND
	}
	if st.Lfe != 0 {
		st.Bandwidth = OPUS_BANDWIDTH_NARROWBAND
	}
	curr_bandwidth = st.Bandwidth
	if st.Mode == MODE_SILK_ONLY && curr_bandwidth > OPUS_BANDWIDTH_WIDEBAND {
		st.Mode = MODE_HYBRID
	}
	if st.Mode == MODE_HYBRID && curr_bandwidth <= OPUS_BANDWIDTH_WIDEBAND {
		st.Mode = MODE_SILK_ONLY
	}
	if frame_size > int64(st.Fs/50) && st.Mode != MODE_SILK_ONLY || frame_size > int64(st.Fs*3/50) {
		var (
			enc_frame_size int64
			nb_frames      int64
		)
		if st.Mode == MODE_SILK_ONLY {
			if frame_size == int64(st.Fs*2/25) {
				enc_frame_size = int64(st.Fs / 25)
			} else if frame_size == int64(st.Fs*3/25) {
				enc_frame_size = int64(st.Fs * 3 / 50)
			} else {
				enc_frame_size = int64(st.Fs / 50)
			}
		} else {
			enc_frame_size = int64(st.Fs / 50)
		}
		nb_frames = frame_size / enc_frame_size
		if analysis_read_pos_bak != -1 {
			st.Analysis.Read_pos = analysis_read_pos_bak
			st.Analysis.Read_subframe = analysis_read_subframe_bak
		}
		ret = int64(encode_multiframe_packet(st, pcm, nb_frames, enc_frame_size, data, out_data_bytes, to_celt, lsb_depth, float_api))
		return opus_int32(ret)
	}
	if st.Silk_bw_switch != 0 {
		redundancy = 1
		celt_to_silk = 1
		st.Silk_bw_switch = 0
		prefill = 2
	}
	if st.Mode == MODE_CELT_ONLY {
		redundancy = 0
	}
	if redundancy != 0 {
		redundancy_bytes = compute_redundancy_bytes(max_data_bytes, st.Bitrate_bps, frame_rate, st.Stream_channels)
		if redundancy_bytes == 0 {
			redundancy = 0
		}
	}
	bytes_target = int64((func() opus_int32 {
		if (max_data_bytes - opus_int32(redundancy_bytes)) < (st.Bitrate_bps * opus_int32(frame_size) / (st.Fs * 8)) {
			return max_data_bytes - opus_int32(redundancy_bytes)
		}
		return st.Bitrate_bps * opus_int32(frame_size) / (st.Fs * 8)
	}()) - 1)
	data = (*uint8)(unsafe.Add(unsafe.Pointer(data), 1))
	ec_enc_init(&enc, data, opus_uint32(max_data_bytes-1))
	pcm_buf = (*opus_val16)(libc.Malloc(int(((total_buffer + frame_size) * st.Channels) * int64(unsafe.Sizeof(opus_val16(0))))))
	libc.MemCpy(unsafe.Pointer(pcm_buf), unsafe.Pointer(&st.Delay_buffer[(st.Encoder_buffer-total_buffer)*st.Channels]), int((total_buffer*st.Channels)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(pcm_buf))-uintptr(unsafe.Pointer(&st.Delay_buffer[(st.Encoder_buffer-total_buffer)*st.Channels]))))*0))
	if st.Mode == MODE_CELT_ONLY {
		hp_freq_smth1 = int64(opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8))
	} else {
		hp_freq_smth1 = int64(((*silk_encoder)(silk_enc)).State_Fxx[0].SCmn.Variable_HP_smth1_Q15)
	}
	st.Variable_HP_smth2_Q15 = st.Variable_HP_smth2_Q15 + opus_int32(((hp_freq_smth1-int64(st.Variable_HP_smth2_Q15))*int64(opus_int16(opus_int32(VARIABLE_HP_SMTH_COEF2*(1<<16)+0.5))))>>16)
	cutoff_Hz = int64(silk_log2lin(st.Variable_HP_smth2_Q15 >> 8))
	if st.Application == OPUS_APPLICATION_VOIP {
		hp_cutoff(pcm, opus_int32(cutoff_Hz), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(total_buffer*st.Channels))), &st.Hp_mem[0], frame_size, st.Channels, st.Fs, st.Arch)
	} else {
		dc_reject(pcm, 3, (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(total_buffer*st.Channels))), &st.Hp_mem[0], frame_size, st.Channels, st.Fs)
	}
	if float_api != 0 {
		var sum opus_val32
		sum = func() opus_val32 {
			st.Arch
			return celt_inner_prod_c((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(total_buffer*st.Channels))), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(total_buffer*st.Channels))), frame_size*st.Channels)
		}()
		if float64(sum) >= 1e+09 || sum != sum {
			libc.MemSet(unsafe.Pointer((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(total_buffer*st.Channels)))), 0, int((frame_size*st.Channels)*int64(unsafe.Sizeof(opus_val16(0)))))
			st.Hp_mem[0] = func() opus_val32 {
				p := &st.Hp_mem[1]
				st.Hp_mem[1] = func() opus_val32 {
					p := &st.Hp_mem[2]
					st.Hp_mem[2] = func() opus_val32 {
						p := &st.Hp_mem[3]
						st.Hp_mem[3] = 0
						return *p
					}()
					return *p
				}()
				return *p
			}()
		}
	}
	HB_gain = opus_val16(Q15ONE)
	if st.Mode != MODE_CELT_ONLY {
		var (
			total_bitRate opus_int32
			celt_rate     opus_int32
			pcm_silk      *opus_int16
		)
		pcm_silk = (*opus_int16)(libc.Malloc(int((st.Channels * frame_size) * int64(unsafe.Sizeof(opus_int16(0))))))
		total_bitRate = opus_int32(bytes_target * 8 * frame_rate)
		if st.Mode == MODE_HYBRID {
			st.Silk_mode.BitRate = opus_int32(compute_silk_rate_for_hybrid(int64(total_bitRate), curr_bandwidth, int64(libc.BoolToInt(st.Fs == opus_int32(frame_size*50))), st.Use_vbr, st.Silk_mode.LBRR_coded, st.Stream_channels))
			if st.Energy_masking == nil {
				celt_rate = total_bitRate - st.Silk_mode.BitRate
				HB_gain = opus_val16(Q15ONE - float64(float32(math.Exp((float64(-celt_rate)*(1.0/1024))*0.6931471805599453))))
			}
		} else {
			st.Silk_mode.BitRate = total_bitRate
		}
		if st.Energy_masking != nil && st.Use_vbr != 0 && st.Lfe == 0 {
			var (
				mask_sum      opus_val32 = 0
				masking_depth opus_val16
				rate_offset   opus_int32
				c             int64
				end           int64      = 17
				srate         opus_int16 = 16000
			)
			if st.Bandwidth == OPUS_BANDWIDTH_NARROWBAND {
				end = 13
				srate = 8000
			} else if st.Bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
				end = 15
				srate = 12000
			}
			for c = 0; c < st.Channels; c++ {
				for i = 0; i < end; i++ {
					var mask opus_val16
					if float64(func() opus_val16 {
						if float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(st.Energy_masking), unsafe.Sizeof(opus_val16(0))*uintptr(c*21+i)))) < 0.5 {
							return *(*opus_val16)(unsafe.Add(unsafe.Pointer(st.Energy_masking), unsafe.Sizeof(opus_val16(0))*uintptr(c*21+i)))
						}
						return opus_val16(0.5)
					}()) > (-2.0) {
						if float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(st.Energy_masking), unsafe.Sizeof(opus_val16(0))*uintptr(c*21+i)))) < 0.5 {
							mask = *(*opus_val16)(unsafe.Add(unsafe.Pointer(st.Energy_masking), unsafe.Sizeof(opus_val16(0))*uintptr(c*21+i)))
						} else {
							mask = opus_val16(0.5)
						}
					} else {
						mask = opus_val16(-2.0)
					}
					if mask > 0 {
						mask = opus_val16(float64(mask) * 0.5)
					}
					mask_sum += opus_val32(mask)
				}
			}
			masking_depth = opus_val16(mask_sum / opus_val32(end) * opus_val32(st.Channels))
			masking_depth += opus_val16(0.2)
			rate_offset = opus_int32(opus_val32(srate) * opus_val32(masking_depth))
			if rate_offset > (st.Silk_mode.BitRate * (-2) / 3) {
				rate_offset = rate_offset
			} else {
				rate_offset = st.Silk_mode.BitRate * (-2) / 3
			}
			if st.Bandwidth == OPUS_BANDWIDTH_SUPERWIDEBAND || st.Bandwidth == OPUS_BANDWIDTH_FULLBAND {
				st.Silk_mode.BitRate += rate_offset * 3 / 5
			} else {
				st.Silk_mode.BitRate += rate_offset
			}
		}
		st.Silk_mode.PayloadSize_ms = frame_size * 1000 / int64(st.Fs)
		st.Silk_mode.NChannelsAPI = opus_int32(st.Channels)
		st.Silk_mode.NChannelsInternal = opus_int32(st.Stream_channels)
		if curr_bandwidth == OPUS_BANDWIDTH_NARROWBAND {
			st.Silk_mode.DesiredInternalSampleRate = 8000
		} else if curr_bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
			st.Silk_mode.DesiredInternalSampleRate = 12000
		} else {
			st.Silk_mode.DesiredInternalSampleRate = 16000
		}
		if st.Mode == MODE_HYBRID {
			st.Silk_mode.MinInternalSampleRate = 16000
		} else {
			st.Silk_mode.MinInternalSampleRate = 8000
		}
		st.Silk_mode.MaxInternalSampleRate = 16000
		if st.Mode == MODE_SILK_ONLY {
			var effective_max_rate opus_int32 = max_rate
			if frame_rate > 50 {
				effective_max_rate = effective_max_rate * 2 / 3
			}
			if effective_max_rate < 8000 {
				st.Silk_mode.MaxInternalSampleRate = 12000
				if 12000 < st.Silk_mode.DesiredInternalSampleRate {
					st.Silk_mode.DesiredInternalSampleRate = 12000
				} else {
					st.Silk_mode.DesiredInternalSampleRate = st.Silk_mode.DesiredInternalSampleRate
				}
			}
			if effective_max_rate < 7000 {
				st.Silk_mode.MaxInternalSampleRate = 8000
				if 8000 < st.Silk_mode.DesiredInternalSampleRate {
					st.Silk_mode.DesiredInternalSampleRate = 8000
				} else {
					st.Silk_mode.DesiredInternalSampleRate = st.Silk_mode.DesiredInternalSampleRate
				}
			}
		}
		st.Silk_mode.UseCBR = int64(libc.BoolToInt(st.Use_vbr == 0))
		st.Silk_mode.MaxBits = int64((max_data_bytes - 1) * 8)
		if redundancy != 0 && redundancy_bytes >= 2 {
			st.Silk_mode.MaxBits -= redundancy_bytes*8 + 1
			if st.Mode == MODE_HYBRID {
				st.Silk_mode.MaxBits -= 20
			}
		}
		if st.Silk_mode.UseCBR != 0 {
			if st.Mode == MODE_HYBRID {
				if st.Silk_mode.MaxBits < int64(st.Silk_mode.BitRate*opus_int32(frame_size)/st.Fs) {
					st.Silk_mode.MaxBits = st.Silk_mode.MaxBits
				} else {
					st.Silk_mode.MaxBits = int64(st.Silk_mode.BitRate * opus_int32(frame_size) / st.Fs)
				}
			}
		} else {
			if st.Mode == MODE_HYBRID {
				var maxBitRate opus_int32 = opus_int32(compute_silk_rate_for_hybrid(st.Silk_mode.MaxBits*int64(st.Fs)/frame_size, curr_bandwidth, int64(libc.BoolToInt(st.Fs == opus_int32(frame_size*50))), st.Use_vbr, st.Silk_mode.LBRR_coded, st.Stream_channels))
				st.Silk_mode.MaxBits = int64(maxBitRate * opus_int32(frame_size) / st.Fs)
			}
		}
		if prefill != 0 {
			var (
				zero           opus_int32 = 0
				prefill_offset int64
			)
			prefill_offset = st.Channels * (st.Encoder_buffer - st.Delay_compensation - int64(st.Fs/400))
			gain_fade(&st.Delay_buffer[prefill_offset], &st.Delay_buffer[prefill_offset], 0, opus_val16(Q15ONE), celt_mode.Overlap, int64(st.Fs/400), st.Channels, celt_mode.Window, st.Fs)
			libc.MemSet(unsafe.Pointer(&st.Delay_buffer[0]), 0, int(prefill_offset*int64(unsafe.Sizeof(opus_val16(0)))))
			for i = 0; i < st.Encoder_buffer*st.Channels; i++ {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pcm_silk), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = FLOAT2INT16(float32(st.Delay_buffer[i]))
			}
			silk_Encode(silk_enc, &st.Silk_mode, pcm_silk, st.Encoder_buffer, nil, &zero, prefill, activity)
			st.Silk_mode.OpusCanSwitch = 0
		}
		for i = 0; i < frame_size*st.Channels; i++ {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pcm_silk), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = FLOAT2INT16(float32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(total_buffer*st.Channels+i)))))
		}
		ret = silk_Encode(silk_enc, &st.Silk_mode, pcm_silk, frame_size, &enc, &nBytes, 0, activity)
		if ret != 0 {
			return -3
		}
		if st.Mode == MODE_SILK_ONLY {
			if st.Silk_mode.InternalSampleRate == 8000 {
				curr_bandwidth = OPUS_BANDWIDTH_NARROWBAND
			} else if st.Silk_mode.InternalSampleRate == 12000 {
				curr_bandwidth = OPUS_BANDWIDTH_MEDIUMBAND
			} else if st.Silk_mode.InternalSampleRate == 16000 {
				curr_bandwidth = OPUS_BANDWIDTH_WIDEBAND
			}
		} else {
		}
		st.Silk_mode.OpusCanSwitch = int64(libc.BoolToInt(st.Silk_mode.SwitchReady != 0 && st.Nonfinal_frame == 0))
		if nBytes == 0 {
			st.RangeFinal = 0
			*(*uint8)(unsafe.Add(unsafe.Pointer(data), -1)) = gen_toc(st.Mode, int64(st.Fs/opus_int32(frame_size)), curr_bandwidth, st.Stream_channels)
			return 1
		}
		if st.Silk_mode.OpusCanSwitch != 0 {
			redundancy_bytes = compute_redundancy_bytes(max_data_bytes, st.Bitrate_bps, frame_rate, st.Stream_channels)
			redundancy = int64(libc.BoolToInt(redundancy_bytes != 0))
			celt_to_silk = 0
			st.Silk_bw_switch = 1
		}
	}
	{
		var endband int64 = 21
		switch curr_bandwidth {
		case OPUS_BANDWIDTH_NARROWBAND:
			endband = 13
		case OPUS_BANDWIDTH_MEDIUMBAND:
			fallthrough
		case OPUS_BANDWIDTH_WIDEBAND:
			endband = 17
		case OPUS_BANDWIDTH_SUPERWIDEBAND:
			endband = 19
		case OPUS_BANDWIDTH_FULLBAND:
			endband = 21
		}
		opus_custom_encoder_ctl(celt_enc, CELT_SET_END_BAND_REQUEST, func() opus_int32 {
			endband == 0
			return opus_int32(endband)
		}())
		opus_custom_encoder_ctl(celt_enc, CELT_SET_CHANNELS_REQUEST, func() opus_int32 {
			st.Stream_channels == 0
			return opus_int32(st.Stream_channels)
		}())
	}
	opus_custom_encoder_ctl(celt_enc, OPUS_SET_BITRATE_REQUEST, func() opus_int32 {
		int64(-1) == 0
		return opus_int32(-1)
	}())
	if st.Mode != MODE_SILK_ONLY {
		var celt_pred opus_val32 = 2
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_REQUEST, func() int64 {
			0 == 0
			return 0
		}())
		if st.Silk_mode.ReducedDependency != 0 {
			celt_pred = 0
		}
		opus_custom_encoder_ctl(celt_enc, CELT_SET_PREDICTION_REQUEST, func() opus_int32 {
			celt_pred == 0
			return opus_int32(celt_pred)
		}())
		if st.Mode == MODE_HYBRID {
			if st.Use_vbr != 0 {
				opus_custom_encoder_ctl(celt_enc, OPUS_SET_BITRATE_REQUEST, func() opus_int32 {
					(st.Bitrate_bps - st.Silk_mode.BitRate) == 0
					return st.Bitrate_bps - st.Silk_mode.BitRate
				}())
				opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_CONSTRAINT_REQUEST, func() int64 {
					0 == 0
					return 0
				}())
			}
		} else {
			if st.Use_vbr != 0 {
				opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_REQUEST, func() int64 {
					1 == 0
					return 1
				}())
				opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_CONSTRAINT_REQUEST, func() opus_int32 {
					st.Vbr_constraint == 0
					return opus_int32(st.Vbr_constraint)
				}())
				opus_custom_encoder_ctl(celt_enc, OPUS_SET_BITRATE_REQUEST, func() opus_int32 {
					st.Bitrate_bps == 0
					return st.Bitrate_bps
				}())
			}
		}
	}
	tmp_prefill = (*opus_val16)(libc.Malloc(int((st.Channels * int64(st.Fs) / 400) * int64(unsafe.Sizeof(opus_val16(0))))))
	if st.Mode != MODE_SILK_ONLY && st.Mode != st.Prev_mode && st.Prev_mode > 0 {
		libc.MemCpy(unsafe.Pointer(tmp_prefill), unsafe.Pointer(&st.Delay_buffer[(st.Encoder_buffer-total_buffer-int64(st.Fs/400))*st.Channels]), int((st.Channels*int64(st.Fs)/400)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(tmp_prefill))-uintptr(unsafe.Pointer(&st.Delay_buffer[(st.Encoder_buffer-total_buffer-int64(st.Fs/400))*st.Channels]))))*0))
	}
	if st.Channels*(st.Encoder_buffer-(frame_size+total_buffer)) > 0 {
		libc.MemMove(unsafe.Pointer(&st.Delay_buffer[0]), unsafe.Pointer(&st.Delay_buffer[st.Channels*frame_size]), int((st.Channels*(st.Encoder_buffer-frame_size-total_buffer))*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(&st.Delay_buffer[0]))-uintptr(unsafe.Pointer(&st.Delay_buffer[st.Channels*frame_size]))))*0))
		libc.MemCpy(unsafe.Pointer(&st.Delay_buffer[st.Channels*(st.Encoder_buffer-frame_size-total_buffer)]), unsafe.Pointer((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*0))), int(((frame_size+total_buffer)*st.Channels)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(&st.Delay_buffer[st.Channels*(st.Encoder_buffer-frame_size-total_buffer)]))-uintptr(unsafe.Pointer((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*0))))))*0))
	} else {
		libc.MemCpy(unsafe.Pointer(&st.Delay_buffer[0]), unsafe.Pointer((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr((frame_size+total_buffer-st.Encoder_buffer)*st.Channels)))), int((st.Encoder_buffer*st.Channels)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(&st.Delay_buffer[0]))-uintptr(unsafe.Pointer((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr((frame_size+total_buffer-st.Encoder_buffer)*st.Channels)))))))*0))
	}
	if float64(st.Prev_HB_gain) < Q15ONE || float64(HB_gain) < Q15ONE {
		gain_fade(pcm_buf, pcm_buf, st.Prev_HB_gain, HB_gain, celt_mode.Overlap, frame_size, st.Channels, celt_mode.Window, st.Fs)
	}
	st.Prev_HB_gain = HB_gain
	if st.Mode != MODE_HYBRID || st.Stream_channels == 1 {
		if equiv_rate > 32000 {
			st.Silk_mode.StereoWidth_Q14 = 16384
		} else if equiv_rate < 16000 {
			st.Silk_mode.StereoWidth_Q14 = 0
		} else {
			st.Silk_mode.StereoWidth_Q14 = int64(16384 - (32000-equiv_rate)*2048/(equiv_rate-14000))
		}
	}
	if st.Energy_masking == nil && st.Channels == 2 {
		if st.Hybrid_stereo_width_Q14 < (1<<14) || st.Silk_mode.StereoWidth_Q14 < (1<<14) {
			var (
				g1 opus_val16
				g2 opus_val16
			)
			g1 = opus_val16(st.Hybrid_stereo_width_Q14)
			g2 = opus_val16(st.Silk_mode.StereoWidth_Q14)
			g1 *= opus_val16(1.0 / 16384)
			g2 *= opus_val16(1.0 / 16384)
			stereo_fade(pcm_buf, pcm_buf, g1, g2, celt_mode.Overlap, frame_size, st.Channels, celt_mode.Window, st.Fs)
			st.Hybrid_stereo_width_Q14 = opus_int16(st.Silk_mode.StereoWidth_Q14)
		}
	}
	if st.Mode != MODE_CELT_ONLY && ec_tell((*ec_ctx)(unsafe.Pointer(&enc)))+17+int64(libc.BoolToInt(st.Mode == MODE_HYBRID))*20 <= int64((max_data_bytes-1)*8) {
		if st.Mode == MODE_HYBRID {
			ec_enc_bit_logp(&enc, redundancy, 12)
		}
		if redundancy != 0 {
			var max_redundancy int64
			ec_enc_bit_logp(&enc, celt_to_silk, 1)
			if st.Mode == MODE_HYBRID {
				max_redundancy = int64((max_data_bytes - 1) - opus_int32((ec_tell((*ec_ctx)(unsafe.Pointer(&enc)))+8+3+7)>>3))
			} else {
				max_redundancy = int64((max_data_bytes - 1) - opus_int32((ec_tell((*ec_ctx)(unsafe.Pointer(&enc)))+7)>>3))
			}
			if max_redundancy < redundancy_bytes {
				redundancy_bytes = max_redundancy
			} else {
				redundancy_bytes = redundancy_bytes
			}
			if 257 < (func() int64 {
				if 2 > redundancy_bytes {
					return 2
				}
				return redundancy_bytes
			}()) {
				redundancy_bytes = 257
			} else if 2 > redundancy_bytes {
				redundancy_bytes = 2
			} else {
				redundancy_bytes = redundancy_bytes
			}
			if st.Mode == MODE_HYBRID {
				ec_enc_uint(&enc, opus_uint32(redundancy_bytes-2), 256)
			}
		}
	} else {
		redundancy = 0
	}
	if redundancy == 0 {
		st.Silk_bw_switch = 0
		redundancy_bytes = 0
	}
	if st.Mode != MODE_CELT_ONLY {
		start_band = 17
	}
	if st.Mode == MODE_SILK_ONLY {
		ret = (ec_tell((*ec_ctx)(unsafe.Pointer(&enc))) + 7) >> 3
		ec_enc_done(&enc)
		nb_compr_bytes = ret
	} else {
		nb_compr_bytes = int64((max_data_bytes - 1) - opus_int32(redundancy_bytes))
		ec_enc_shrink(&enc, opus_uint32(nb_compr_bytes))
	}
	if redundancy != 0 || st.Mode != MODE_SILK_ONLY {
		opus_custom_encoder_ctl(celt_enc, CELT_SET_ANALYSIS_REQUEST, (*AnalysisInfo)(unsafe.Add(unsafe.Pointer(&analysis_info), unsafe.Sizeof(AnalysisInfo{})*uintptr(int64(uintptr(unsafe.Pointer(&analysis_info))-uintptr(unsafe.Pointer(&analysis_info)))))))
	}
	if st.Mode == MODE_HYBRID {
		var info SILKInfo
		info.SignalType = st.Silk_mode.SignalType
		info.Offset = st.Silk_mode.Offset
		opus_custom_encoder_ctl(celt_enc, CELT_SET_SILK_INFO_REQUEST, (*SILKInfo)(unsafe.Add(unsafe.Pointer(&info), unsafe.Sizeof(SILKInfo{})*uintptr(int64(uintptr(unsafe.Pointer(&info))-uintptr(unsafe.Pointer(&info)))))))
	}
	if redundancy != 0 && celt_to_silk != 0 {
		var err int64
		opus_custom_encoder_ctl(celt_enc, CELT_SET_START_BAND_REQUEST, func() int64 {
			0 == 0
			return 0
		}())
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_REQUEST, func() int64 {
			0 == 0
			return 0
		}())
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_BITRATE_REQUEST, func() opus_int32 {
			int64(-1) == 0
			return opus_int32(-1)
		}())
		err = celt_encode_with_ec(celt_enc, pcm_buf, int64(st.Fs/200), (*uint8)(unsafe.Add(unsafe.Pointer(data), nb_compr_bytes)), redundancy_bytes, nil)
		if err < 0 {
			return -3
		}
		opus_custom_encoder_ctl(celt_enc, OPUS_GET_FINAL_RANGE_REQUEST, (*opus_uint32)(unsafe.Add(unsafe.Pointer(&redundant_rng), unsafe.Sizeof(opus_uint32(0))*uintptr(int64(uintptr(unsafe.Pointer(&redundant_rng))-uintptr(unsafe.Pointer(&redundant_rng)))))))
		opus_custom_encoder_ctl(celt_enc, OPUS_RESET_STATE)
	}
	opus_custom_encoder_ctl(celt_enc, CELT_SET_START_BAND_REQUEST, func() opus_int32 {
		start_band == 0
		return opus_int32(start_band)
	}())
	if st.Mode != MODE_SILK_ONLY {
		if st.Mode != st.Prev_mode && st.Prev_mode > 0 {
			var dummy [2]uint8
			opus_custom_encoder_ctl(celt_enc, OPUS_RESET_STATE)
			celt_encode_with_ec(celt_enc, tmp_prefill, int64(st.Fs/400), &dummy[0], 2, nil)
			opus_custom_encoder_ctl(celt_enc, CELT_SET_PREDICTION_REQUEST, func() int64 {
				0 == 0
				return 0
			}())
		}
		if ec_tell((*ec_ctx)(unsafe.Pointer(&enc))) <= nb_compr_bytes*8 {
			if redundancy != 0 && celt_to_silk != 0 && st.Mode == MODE_HYBRID && st.Use_vbr != 0 {
				opus_custom_encoder_ctl(celt_enc, OPUS_SET_BITRATE_REQUEST, func() opus_int32 {
					(st.Bitrate_bps - st.Silk_mode.BitRate) == 0
					return st.Bitrate_bps - st.Silk_mode.BitRate
				}())
			}
			opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_REQUEST, func() opus_int32 {
				st.Use_vbr == 0
				return opus_int32(st.Use_vbr)
			}())
			ret = celt_encode_with_ec(celt_enc, pcm_buf, frame_size, nil, nb_compr_bytes, &enc)
			if ret < 0 {
				return -3
			}
			if redundancy != 0 && celt_to_silk != 0 && st.Mode == MODE_HYBRID && st.Use_vbr != 0 {
				libc.MemMove(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), ret))), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), nb_compr_bytes))), int(redundancy_bytes*int64(unsafe.Sizeof(uint8(0)))+(int64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), ret))))-uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), nb_compr_bytes))))))*0))
				nb_compr_bytes = nb_compr_bytes + redundancy_bytes
			}
		}
	}
	if redundancy != 0 && celt_to_silk == 0 {
		var (
			err   int64
			dummy [2]uint8
			N2    int64
			N4    int64
		)
		N2 = int64(st.Fs / 200)
		N4 = int64(st.Fs / 400)
		opus_custom_encoder_ctl(celt_enc, OPUS_RESET_STATE)
		opus_custom_encoder_ctl(celt_enc, CELT_SET_START_BAND_REQUEST, func() int64 {
			0 == 0
			return 0
		}())
		opus_custom_encoder_ctl(celt_enc, CELT_SET_PREDICTION_REQUEST, func() int64 {
			0 == 0
			return 0
		}())
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_VBR_REQUEST, func() int64 {
			0 == 0
			return 0
		}())
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_BITRATE_REQUEST, func() opus_int32 {
			int64(-1) == 0
			return opus_int32(-1)
		}())
		if st.Mode == MODE_HYBRID {
			nb_compr_bytes = ret
			ec_enc_shrink(&enc, opus_uint32(nb_compr_bytes))
		}
		celt_encode_with_ec(celt_enc, (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*(frame_size-N2-N4)))), N4, &dummy[0], 2, nil)
		err = celt_encode_with_ec(celt_enc, (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_buf), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*(frame_size-N2)))), N2, (*uint8)(unsafe.Add(unsafe.Pointer(data), nb_compr_bytes)), redundancy_bytes, nil)
		if err < 0 {
			return -3
		}
		opus_custom_encoder_ctl(celt_enc, OPUS_GET_FINAL_RANGE_REQUEST, (*opus_uint32)(unsafe.Add(unsafe.Pointer(&redundant_rng), unsafe.Sizeof(opus_uint32(0))*uintptr(int64(uintptr(unsafe.Pointer(&redundant_rng))-uintptr(unsafe.Pointer(&redundant_rng)))))))
	}
	data = (*uint8)(unsafe.Add(unsafe.Pointer(data), -1))
	*data = gen_toc(st.Mode, int64(st.Fs/opus_int32(frame_size)), curr_bandwidth, st.Stream_channels)
	st.RangeFinal = enc.Rng ^ redundant_rng
	if to_celt != 0 {
		st.Prev_mode = MODE_CELT_ONLY
	} else {
		st.Prev_mode = st.Mode
	}
	st.Prev_channels = st.Stream_channels
	st.Prev_framesize = frame_size
	st.First = 0
	if st.Use_dtx != 0 && (analysis_info.Valid != 0 || is_silence != 0) {
		if decide_dtx_mode(activity, &st.Nb_no_activity_ms_Q1, frame_size*(2*1000)/int64(st.Fs)) != 0 {
			st.RangeFinal = 0
			*data = gen_toc(st.Mode, int64(st.Fs/opus_int32(frame_size)), curr_bandwidth, st.Stream_channels)
			return 1
		}
	} else {
		st.Nb_no_activity_ms_Q1 = 0
	}
	if ec_tell((*ec_ctx)(unsafe.Pointer(&enc))) > int64((max_data_bytes-1)*8) {
		if max_data_bytes < 2 {
			return -2
		}
		*(*uint8)(unsafe.Add(unsafe.Pointer(data), 1)) = 0
		ret = 1
		st.RangeFinal = 0
	} else if st.Mode == MODE_SILK_ONLY && redundancy == 0 {
		for ret > 2 && int64(*(*uint8)(unsafe.Add(unsafe.Pointer(data), ret))) == 0 {
			ret--
		}
	}
	ret += redundancy_bytes + 1
	if st.Use_vbr == 0 {
		if opus_packet_pad(data, opus_int32(ret), max_data_bytes) != OPUS_OK {
			return -3
		}
		ret = int64(max_data_bytes)
	}
	return opus_int32(ret)
}
func opus_encode(st *OpusEncoder, pcm *opus_int16, analysis_frame_size int64, data *uint8, max_data_bytes opus_int32) opus_int32 {
	var (
		i          int64
		ret        int64
		frame_size int64
		in         *float32
	)
	frame_size = int64(frame_size_select(opus_int32(analysis_frame_size), st.Variable_duration, st.Fs))
	if frame_size <= 0 {
		return -1
	}
	in = (*float32)(libc.Malloc(int((frame_size * st.Channels) * int64(unsafe.Sizeof(float32(0))))))
	for i = 0; i < frame_size*st.Channels; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(i))) = float32(float64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) * (1.0 / 32768))
	}
	ret = int64(opus_encode_native(st, (*opus_val16)(unsafe.Pointer(in)), frame_size, data, max_data_bytes, 16, unsafe.Pointer(pcm), opus_int32(analysis_frame_size), 0, -2, st.Channels, func(arg1 unsafe.Pointer, arg2 *opus_val32, arg3 int64, arg4 int64, arg5 int64, arg6 int64, arg7 int64) {
		downmix_int(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}, 0))
	return opus_int32(ret)
}
func opus_encode_float(st *OpusEncoder, pcm *float32, analysis_frame_size int64, data *uint8, out_data_bytes opus_int32) opus_int32 {
	var frame_size int64
	frame_size = int64(frame_size_select(opus_int32(analysis_frame_size), st.Variable_duration, st.Fs))
	return opus_encode_native(st, (*opus_val16)(unsafe.Pointer(pcm)), frame_size, data, out_data_bytes, 24, unsafe.Pointer(pcm), opus_int32(analysis_frame_size), 0, -2, st.Channels, func(arg1 unsafe.Pointer, arg2 *opus_val32, arg3 int64, arg4 int64, arg5 int64, arg6 int64, arg7 int64) {
		downmix_float(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}, 1)
}
func opus_encoder_ctl(st *OpusEncoder, request int64, _rest ...interface{}) int64 {
	var (
		ret      int64
		celt_enc *OpusCustomEncoder
		ap       libc.ArgList
	)
	ret = OPUS_OK
	ap.Start(request, _rest)
	celt_enc = (*OpusCustomEncoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Celt_enc_offset))))
	switch request {
	case OPUS_SET_APPLICATION_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value != OPUS_APPLICATION_VOIP && value != OPUS_APPLICATION_AUDIO && value != OPUS_APPLICATION_RESTRICTED_LOWDELAY || st.First == 0 && st.Application != int64(value) {
			ret = -1
			break
		}
		st.Application = int64(value)
		st.Analysis.Application = int64(value)
	case OPUS_GET_APPLICATION_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Application)
	case OPUS_SET_BITRATE_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value != -1000 && value != -1 {
			if value <= 0 {
				goto bad_arg
			} else if value <= 500 {
				value = 500
			} else if value > opus_int32(st.Channels*300000) {
				value = opus_int32(st.Channels * 300000)
			}
		}
		st.User_bitrate_bps = value
	case OPUS_GET_BITRATE_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = user_bitrate_to_bitrate(st, st.Prev_framesize, 1276)
	case OPUS_SET_FORCE_CHANNELS_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if (value < 1 || value > opus_int32(st.Channels)) && value != -1000 {
			goto bad_arg
		}
		st.Force_channels = int64(value)
	case OPUS_GET_FORCE_CHANNELS_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Force_channels)
	case OPUS_SET_MAX_BANDWIDTH_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < OPUS_BANDWIDTH_NARROWBAND || value > OPUS_BANDWIDTH_FULLBAND {
			goto bad_arg
		}
		st.Max_bandwidth = int64(value)
		if st.Max_bandwidth == OPUS_BANDWIDTH_NARROWBAND {
			st.Silk_mode.MaxInternalSampleRate = 8000
		} else if st.Max_bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
			st.Silk_mode.MaxInternalSampleRate = 12000
		} else {
			st.Silk_mode.MaxInternalSampleRate = 16000
		}
	case OPUS_GET_MAX_BANDWIDTH_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Max_bandwidth)
	case OPUS_SET_BANDWIDTH_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if (value < OPUS_BANDWIDTH_NARROWBAND || value > OPUS_BANDWIDTH_FULLBAND) && value != -1000 {
			goto bad_arg
		}
		st.User_bandwidth = int64(value)
		if st.User_bandwidth == OPUS_BANDWIDTH_NARROWBAND {
			st.Silk_mode.MaxInternalSampleRate = 8000
		} else if st.User_bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
			st.Silk_mode.MaxInternalSampleRate = 12000
		} else {
			st.Silk_mode.MaxInternalSampleRate = 16000
		}
	case OPUS_GET_BANDWIDTH_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Bandwidth)
	case OPUS_SET_DTX_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 1 {
			goto bad_arg
		}
		st.Use_dtx = int64(value)
	case OPUS_GET_DTX_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Use_dtx)
	case OPUS_SET_COMPLEXITY_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 10 {
			goto bad_arg
		}
		st.Silk_mode.Complexity = int64(value)
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_COMPLEXITY_REQUEST, func() opus_int32 {
			value == 0
			return value
		}())
	case OPUS_GET_COMPLEXITY_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Silk_mode.Complexity)
	case OPUS_SET_INBAND_FEC_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 2 {
			goto bad_arg
		}
		st.Fec_config = int64(value)
		st.Silk_mode.UseInBandFEC = int64(libc.BoolToInt(value != 0))
	case OPUS_GET_INBAND_FEC_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Fec_config)
	case OPUS_SET_PACKET_LOSS_PERC_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 100 {
			goto bad_arg
		}
		st.Silk_mode.PacketLossPercentage = int64(value)
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_PACKET_LOSS_PERC_REQUEST, func() opus_int32 {
			value == 0
			return value
		}())
	case OPUS_GET_PACKET_LOSS_PERC_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Silk_mode.PacketLossPercentage)
	case OPUS_SET_VBR_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 1 {
			goto bad_arg
		}
		st.Use_vbr = int64(value)
		st.Silk_mode.UseCBR = int64(1 - value)
	case OPUS_GET_VBR_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Use_vbr)
	case OPUS_SET_VOICE_RATIO_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < -1 || value > 100 {
			goto bad_arg
		}
		st.Voice_ratio = int64(value)
	case OPUS_GET_VOICE_RATIO_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Voice_ratio)
	case OPUS_SET_VBR_CONSTRAINT_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 1 {
			goto bad_arg
		}
		st.Vbr_constraint = int64(value)
	case OPUS_GET_VBR_CONSTRAINT_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Vbr_constraint)
	case OPUS_SET_SIGNAL_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value != -1000 && value != OPUS_SIGNAL_VOICE && value != OPUS_SIGNAL_MUSIC {
			goto bad_arg
		}
		st.Signal_type = int64(value)
	case OPUS_GET_SIGNAL_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Signal_type)
	case OPUS_GET_LOOKAHEAD_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = st.Fs / 400
		if st.Application != OPUS_APPLICATION_RESTRICTED_LOWDELAY {
			*value += opus_int32(st.Delay_compensation)
		}
	case OPUS_GET_SAMPLE_RATE_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = st.Fs
	case OPUS_GET_FINAL_RANGE_REQUEST:
		var value *opus_uint32 = ap.Arg().(*opus_uint32)
		if value == nil {
			goto bad_arg
		}
		*value = st.RangeFinal
	case OPUS_SET_LSB_DEPTH_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 8 || value > 24 {
			goto bad_arg
		}
		st.Lsb_depth = int64(value)
	case OPUS_GET_LSB_DEPTH_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Lsb_depth)
	case OPUS_SET_EXPERT_FRAME_DURATION_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value != OPUS_FRAMESIZE_ARG && value != OPUS_FRAMESIZE_2_5_MS && value != OPUS_FRAMESIZE_5_MS && value != OPUS_FRAMESIZE_10_MS && value != OPUS_FRAMESIZE_20_MS && value != OPUS_FRAMESIZE_40_MS && value != OPUS_FRAMESIZE_60_MS && value != OPUS_FRAMESIZE_80_MS && value != OPUS_FRAMESIZE_100_MS && value != OPUS_FRAMESIZE_120_MS {
			goto bad_arg
		}
		st.Variable_duration = int64(value)
	case OPUS_GET_EXPERT_FRAME_DURATION_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Variable_duration)
	case OPUS_SET_PREDICTION_DISABLED_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value > 1 || value < 0 {
			goto bad_arg
		}
		st.Silk_mode.ReducedDependency = int64(value)
	case OPUS_GET_PREDICTION_DISABLED_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Silk_mode.ReducedDependency)
	case OPUS_SET_PHASE_INVERSION_DISABLED_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 1 {
			goto bad_arg
		}
		opus_custom_encoder_ctl(celt_enc, OPUS_SET_PHASE_INVERSION_DISABLED_REQUEST, func() opus_int32 {
			value == 0
			return value
		}())
	case OPUS_GET_PHASE_INVERSION_DISABLED_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		opus_custom_encoder_ctl(celt_enc, OPUS_GET_PHASE_INVERSION_DISABLED_REQUEST, (*opus_int32)(unsafe.Add(unsafe.Pointer(value), unsafe.Sizeof(opus_int32(0))*uintptr(int64(uintptr(unsafe.Pointer(value))-uintptr(unsafe.Pointer(value)))))))
	case OPUS_RESET_STATE:
		var (
			silk_enc unsafe.Pointer
			dummy    silk_EncControlStruct
			start    *byte
		)
		silk_enc = unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_enc_offset)
		tonality_analysis_reset(&st.Analysis)
		start = (*byte)(unsafe.Pointer(&st.Stream_channels))
		libc.MemSet(unsafe.Pointer(start), 0, int((unsafe.Sizeof(OpusEncoder{})-uintptr(int64(uintptr(unsafe.Pointer(start))-uintptr(unsafe.Pointer((*byte)(unsafe.Pointer(st)))))))*unsafe.Sizeof(byte(0))))
		opus_custom_encoder_ctl(celt_enc, OPUS_RESET_STATE)
		silk_InitEncoder(silk_enc, st.Arch, &dummy)
		st.Stream_channels = st.Channels
		st.Hybrid_stereo_width_Q14 = 1 << 14
		st.Prev_HB_gain = opus_val16(Q15ONE)
		st.First = 1
		st.Mode = MODE_HYBRID
		st.Bandwidth = OPUS_BANDWIDTH_FULLBAND
		st.Variable_HP_smth2_Q15 = opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)
	case OPUS_SET_FORCE_MODE_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if (value < MODE_SILK_ONLY || value > MODE_CELT_ONLY) && value != -1000 {
			goto bad_arg
		}
		st.User_forced_mode = int64(value)
	case OPUS_SET_LFE_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		st.Lfe = int64(value)
		ret = opus_custom_encoder_ctl(celt_enc, OPUS_SET_LFE_REQUEST, func() opus_int32 {
			value == 0
			return value
		}())
	case OPUS_SET_ENERGY_MASK_REQUEST:
		var value *opus_val16 = ap.Arg().(*opus_val16)
		st.Energy_masking = value
		ret = opus_custom_encoder_ctl(celt_enc, OPUS_SET_ENERGY_MASK_REQUEST, (*opus_val16)(unsafe.Add(unsafe.Pointer(value), unsafe.Sizeof(opus_val16(0))*uintptr(int64(uintptr(unsafe.Pointer(value))-uintptr(unsafe.Pointer(value)))))))
	case OPUS_GET_IN_DTX_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		if st.Silk_mode.UseDTX != 0 && (st.Prev_mode == MODE_SILK_ONLY || st.Prev_mode == MODE_HYBRID) {
			var silk_enc *silk_encoder = (*silk_encoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_enc_offset))))
			*value = opus_int32(libc.BoolToInt(silk_enc.State_Fxx[0].SCmn.NoSpeechCounter >= NB_SPEECH_FRAMES_BEFORE_DTX))
			if *value == 1 && st.Silk_mode.NChannelsInternal == 2 && silk_enc.Prev_decode_only_middle == 0 {
				*value = opus_int32(libc.BoolToInt(silk_enc.State_Fxx[1].SCmn.NoSpeechCounter >= NB_SPEECH_FRAMES_BEFORE_DTX))
			}
		} else if st.Use_dtx != 0 {
			*value = opus_int32(libc.BoolToInt(st.Nb_no_activity_ms_Q1 >= NB_SPEECH_FRAMES_BEFORE_DTX*20*2))
		} else {
			*value = 0
		}
	case CELT_GET_MODE_REQUEST:
		var value **OpusCustomMode = ap.Arg().(**OpusCustomMode)
		if value == nil {
			goto bad_arg
		}
		ret = opus_custom_encoder_ctl(celt_enc, CELT_GET_MODE_REQUEST, (**OpusCustomMode)(unsafe.Add(unsafe.Pointer(value), unsafe.Sizeof((*OpusCustomMode)(nil))*uintptr(int64(uintptr(unsafe.Pointer(value))-uintptr(unsafe.Pointer(value)))))))
	default:
		ret = -5
	}
	ap.End()
	return ret
bad_arg:
	ap.End()
	return -1
}
func opus_encoder_destroy(st *OpusEncoder) {
	libc.Free(unsafe.Pointer(st))
}
