package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

type OpusDecoder struct {
	Celt_dec_offset      int64
	Silk_dec_offset      int64
	Channels             int64
	Fs                   opus_int32
	DecControl           silk_DecControlStruct
	Decode_gain          int64
	Arch                 int64
	Stream_channels      int64
	Bandwidth            int64
	Mode                 int64
	Prev_mode            int64
	Frame_size           int64
	Prev_redundancy      int64
	Last_packet_duration int64
	Softclip_mem         [2]opus_val16
	RangeFinal           opus_uint32
}

func opus_decoder_get_size(channels int64) int64 {
	var (
		silkDecSizeBytes int64
		celtDecSizeBytes int64
		ret              int64
	)
	if channels < 1 || channels > 2 {
		return 0
	}
	ret = silk_Get_Decoder_Size(&silkDecSizeBytes)
	if ret != 0 {
		return 0
	}
	silkDecSizeBytes = align(silkDecSizeBytes)
	celtDecSizeBytes = celt_decoder_get_size(channels)
	return align(int64(unsafe.Sizeof(OpusDecoder{}))) + silkDecSizeBytes + celtDecSizeBytes
}
func opus_decoder_init(st *OpusDecoder, Fs opus_int32, channels int64) int64 {
	var (
		silk_dec         unsafe.Pointer
		celt_dec         *OpusCustomDecoder
		ret              int64
		silkDecSizeBytes int64
	)
	if Fs != 48000 && Fs != 24000 && Fs != 16000 && Fs != 12000 && Fs != 8000 || channels != 1 && channels != 2 {
		return -1
	}
	libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(st))), 0, int(opus_decoder_get_size(channels)*int64(unsafe.Sizeof(byte(0)))))
	ret = silk_Get_Decoder_Size(&silkDecSizeBytes)
	if ret != 0 {
		return -3
	}
	silkDecSizeBytes = align(silkDecSizeBytes)
	st.Silk_dec_offset = align(int64(unsafe.Sizeof(OpusDecoder{})))
	st.Celt_dec_offset = st.Silk_dec_offset + silkDecSizeBytes
	silk_dec = unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_dec_offset)
	celt_dec = (*OpusCustomDecoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Celt_dec_offset))))
	st.Stream_channels = func() int64 {
		p := &st.Channels
		st.Channels = channels
		return *p
	}()
	st.Fs = Fs
	st.DecControl.API_sampleRate = st.Fs
	st.DecControl.NChannelsAPI = opus_int32(st.Channels)
	ret = silk_InitDecoder(silk_dec)
	if ret != 0 {
		return -3
	}
	ret = celt_decoder_init(celt_dec, Fs, channels)
	if ret != OPUS_OK {
		return -3
	}
	opus_custom_decoder_ctl(celt_dec, CELT_SET_SIGNALLING_REQUEST, func() int64 {
		0 == 0
		return 0
	}())
	st.Prev_mode = 0
	st.Frame_size = int64(Fs / 400)
	st.Arch = opus_select_arch()
	return OPUS_OK
}
func opus_decoder_create(Fs opus_int32, channels int64, error *int64) *OpusDecoder {
	var (
		ret int64
		st  *OpusDecoder
	)
	if Fs != 48000 && Fs != 24000 && Fs != 16000 && Fs != 12000 && Fs != 8000 || channels != 1 && channels != 2 {
		if error != nil {
			*error = -1
		}
		return nil
	}
	st = (*OpusDecoder)(libc.Malloc(int(opus_decoder_get_size(channels))))
	if st == nil {
		if error != nil {
			*error = -7
		}
		return nil
	}
	ret = opus_decoder_init(st, Fs, channels)
	if error != nil {
		*error = ret
	}
	if ret != OPUS_OK {
		libc.Free(unsafe.Pointer(st))
		st = nil
	}
	return st
}
func smooth_fade(in1 *opus_val16, in2 *opus_val16, out *opus_val16, overlap int64, channels int64, window *opus_val16, Fs opus_int32) {
	var (
		i   int64
		c   int64
		inc int64 = int64(48000 / Fs)
	)
	for c = 0; c < channels; c++ {
		for i = 0; i < overlap; i++ {
			var w opus_val16 = ((*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc)))) * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i*inc)))))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+c))) = opus_val16((opus_val32(w) * opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in2), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+c))))) + opus_val32(Q15ONE-float64(w))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(in1), unsafe.Sizeof(opus_val16(0))*uintptr(i*channels+c)))))
		}
	}
}
func opus_packet_get_mode(data *uint8) int64 {
	var mode int64
	if int64(*data)&0x80 != 0 {
		mode = MODE_CELT_ONLY
	} else if (int64(*data) & 0x60) == 0x60 {
		mode = MODE_HYBRID
	} else {
		mode = MODE_SILK_ONLY
	}
	return mode
}
func opus_decode_frame(st *OpusDecoder, data *uint8, len_ opus_int32, pcm *opus_val16, frame_size int64, decode_fec int64) int64 {
	var (
		silk_dec                 unsafe.Pointer
		celt_dec                 *OpusCustomDecoder
		i                        int64
		silk_ret                 int64 = 0
		celt_ret                 int64 = 0
		dec                      ec_dec
		silk_frame_size          opus_int32
		pcm_silk_size            int64
		pcm_silk                 *opus_int16
		pcm_transition_silk_size int64
		pcm_transition_silk      *opus_val16
		pcm_transition_celt_size int64
		pcm_transition_celt      *opus_val16
		pcm_transition           *opus_val16 = nil
		redundant_audio_size     int64
		redundant_audio          *opus_val16
		audiosize                int64
		mode                     int64
		bandwidth                int64
		transition               int64 = 0
		start_band               int64
		redundancy               int64 = 0
		redundancy_bytes         int64 = 0
		celt_to_silk             int64 = 0
		c                        int64
		F2_5                     int64
		F5                       int64
		F10                      int64
		F20                      int64
		window                   *opus_val16
		redundant_rng            opus_uint32 = 0
		celt_accum               int64
	)
	silk_dec = unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_dec_offset)
	celt_dec = (*OpusCustomDecoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Celt_dec_offset))))
	F20 = int64(st.Fs / 50)
	F10 = F20 >> 1
	F5 = F10 >> 1
	F2_5 = F5 >> 1
	if frame_size < F2_5 {
		return -2
	}
	if frame_size < int64(st.Fs/25*3) {
		frame_size = frame_size
	} else {
		frame_size = int64(st.Fs / 25 * 3)
	}
	if len_ <= 1 {
		data = nil
		if frame_size < st.Frame_size {
			frame_size = frame_size
		} else {
			frame_size = st.Frame_size
		}
	}
	if data != nil {
		audiosize = st.Frame_size
		mode = st.Mode
		bandwidth = st.Bandwidth
		ec_dec_init(&dec, data, opus_uint32(len_))
	} else {
		audiosize = frame_size
		if st.Prev_redundancy != 0 {
			mode = MODE_CELT_ONLY
		} else {
			mode = st.Prev_mode
		}
		bandwidth = 0
		if mode == 0 {
			for i = 0; i < audiosize*st.Channels; i++ {
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = 0
			}
			return audiosize
		}
		if audiosize > F20 {
			for {
				{
					var ret int64 = opus_decode_frame(st, nil, 0, pcm, func() int64 {
						if audiosize < F20 {
							return audiosize
						}
						return F20
					}(), 0)
					if ret < 0 {
						return ret
					}
					pcm = (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(ret*st.Channels)))
					audiosize -= ret
				}
				if audiosize <= 0 {
					break
				}
			}
			return frame_size
		} else if audiosize < F20 {
			if audiosize > F10 {
				audiosize = F10
			} else if mode != MODE_SILK_ONLY && audiosize > F5 && audiosize < F10 {
				audiosize = F5
			}
		}
	}
	celt_accum = 0
	pcm_transition_silk_size = ALLOC_NONE
	pcm_transition_celt_size = ALLOC_NONE
	if data != nil && st.Prev_mode > 0 && (mode == MODE_CELT_ONLY && st.Prev_mode != MODE_CELT_ONLY && st.Prev_redundancy == 0 || mode != MODE_CELT_ONLY && st.Prev_mode == MODE_CELT_ONLY) {
		transition = 1
		if mode == MODE_CELT_ONLY {
			pcm_transition_celt_size = F5 * st.Channels
		} else {
			pcm_transition_silk_size = F5 * st.Channels
		}
	}
	pcm_transition_celt = (*opus_val16)(libc.Malloc(int(pcm_transition_celt_size * int64(unsafe.Sizeof(opus_val16(0))))))
	if transition != 0 && mode == MODE_CELT_ONLY {
		pcm_transition = pcm_transition_celt
		opus_decode_frame(st, nil, 0, pcm_transition, func() int64 {
			if F5 < audiosize {
				return F5
			}
			return audiosize
		}(), 0)
	}
	if audiosize > frame_size {
		return -1
	} else {
		frame_size = audiosize
	}
	if mode != MODE_CELT_ONLY && celt_accum == 0 {
		pcm_silk_size = (func() int64 {
			if F10 > frame_size {
				return F10
			}
			return frame_size
		}()) * st.Channels
	} else {
		pcm_silk_size = ALLOC_NONE
	}
	pcm_silk = (*opus_int16)(libc.Malloc(int(pcm_silk_size * int64(unsafe.Sizeof(opus_int16(0))))))
	if mode != MODE_CELT_ONLY {
		var (
			lost_flag       int64
			decoded_samples int64
			pcm_ptr         *opus_int16
		)
		pcm_ptr = pcm_silk
		if st.Prev_mode == MODE_CELT_ONLY {
			silk_InitDecoder(silk_dec)
		}
		if 10 > (audiosize * 1000 / int64(st.Fs)) {
			st.DecControl.PayloadSize_ms = 10
		} else {
			st.DecControl.PayloadSize_ms = audiosize * 1000 / int64(st.Fs)
		}
		if data != nil {
			st.DecControl.NChannelsInternal = opus_int32(st.Stream_channels)
			if mode == MODE_SILK_ONLY {
				if bandwidth == OPUS_BANDWIDTH_NARROWBAND {
					st.DecControl.InternalSampleRate = 8000
				} else if bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
					st.DecControl.InternalSampleRate = 12000
				} else if bandwidth == OPUS_BANDWIDTH_WIDEBAND {
					st.DecControl.InternalSampleRate = 16000
				} else {
					st.DecControl.InternalSampleRate = 16000
				}
			} else {
				st.DecControl.InternalSampleRate = 16000
			}
		}
		if data == nil {
			lost_flag = 1
		} else {
			lost_flag = decode_fec * 2
		}
		decoded_samples = 0
		for {
			{
				var first_frame int64 = int64(libc.BoolToInt(decoded_samples == 0))
				silk_ret = silk_Decode(silk_dec, &st.DecControl, lost_flag, first_frame, &dec, pcm_ptr, &silk_frame_size, st.Arch)
				if silk_ret != 0 {
					if lost_flag != 0 {
						silk_frame_size = opus_int32(frame_size)
						for i = 0; i < frame_size*st.Channels; i++ {
							*(*opus_int16)(unsafe.Add(unsafe.Pointer(pcm_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = 0
						}
					} else {
						return -3
					}
				}
				pcm_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(pcm_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(silk_frame_size*opus_int32(st.Channels))))
				decoded_samples += int64(silk_frame_size)
			}
			if decoded_samples >= frame_size {
				break
			}
		}
	}
	start_band = 0
	if decode_fec == 0 && mode != MODE_CELT_ONLY && data != nil && ec_tell((*ec_ctx)(unsafe.Pointer(&dec)))+17+int64(libc.BoolToInt(mode == MODE_HYBRID))*20 <= int64(len_*8) {
		if mode == MODE_HYBRID {
			redundancy = ec_dec_bit_logp(&dec, 12)
		} else {
			redundancy = 1
		}
		if redundancy != 0 {
			celt_to_silk = ec_dec_bit_logp(&dec, 1)
			if mode == MODE_HYBRID {
				redundancy_bytes = int64(opus_int32(ec_dec_uint(&dec, 256)) + 2)
			} else {
				redundancy_bytes = int64(len_ - opus_int32((ec_tell((*ec_ctx)(unsafe.Pointer(&dec)))+7)>>3))
			}
			len_ -= opus_int32(redundancy_bytes)
			if len_*8 < opus_int32(ec_tell((*ec_ctx)(unsafe.Pointer(&dec)))) {
				len_ = 0
				redundancy_bytes = 0
				redundancy = 0
			}
			dec.Storage -= opus_uint32(redundancy_bytes)
		}
	}
	if mode != MODE_CELT_ONLY {
		start_band = 17
	}
	if redundancy != 0 {
		transition = 0
		pcm_transition_silk_size = ALLOC_NONE
	}
	pcm_transition_silk = (*opus_val16)(libc.Malloc(int(pcm_transition_silk_size * int64(unsafe.Sizeof(opus_val16(0))))))
	if transition != 0 && mode != MODE_CELT_ONLY {
		pcm_transition = pcm_transition_silk
		opus_decode_frame(st, nil, 0, pcm_transition, func() int64 {
			if F5 < audiosize {
				return F5
			}
			return audiosize
		}(), 0)
	}
	if bandwidth != 0 {
		var endband int64 = 21
		switch bandwidth {
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
		default:
		}
		for {
			if opus_custom_decoder_ctl(celt_dec, CELT_SET_END_BAND_REQUEST, func() opus_int32 {
				endband == 0
				return opus_int32(endband)
			}()) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
	}
	for {
		if opus_custom_decoder_ctl(celt_dec, CELT_SET_CHANNELS_REQUEST, func() opus_int32 {
			st.Stream_channels == 0
			return opus_int32(st.Stream_channels)
		}()) != OPUS_OK {
			return -3
		}
		if true {
			break
		}
	}
	if redundancy != 0 {
		redundant_audio_size = F5 * st.Channels
	} else {
		redundant_audio_size = ALLOC_NONE
	}
	redundant_audio = (*opus_val16)(libc.Malloc(int(redundant_audio_size * int64(unsafe.Sizeof(opus_val16(0))))))
	if redundancy != 0 && celt_to_silk != 0 {
		for {
			if opus_custom_decoder_ctl(celt_dec, CELT_SET_START_BAND_REQUEST, func() int64 {
				0 == 0
				return 0
			}()) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
		celt_decode_with_ec(celt_dec, (*uint8)(unsafe.Add(unsafe.Pointer(data), len_)), redundancy_bytes, redundant_audio, F5, nil, 0)
		for {
			if opus_custom_decoder_ctl(celt_dec, OPUS_GET_FINAL_RANGE_REQUEST, (*opus_uint32)(unsafe.Add(unsafe.Pointer(&redundant_rng), unsafe.Sizeof(opus_uint32(0))*uintptr(int64(uintptr(unsafe.Pointer(&redundant_rng))-uintptr(unsafe.Pointer(&redundant_rng))))))) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
	}
	for {
		if opus_custom_decoder_ctl(celt_dec, CELT_SET_START_BAND_REQUEST, func() opus_int32 {
			start_band == 0
			return opus_int32(start_band)
		}()) != OPUS_OK {
			return -3
		}
		if true {
			break
		}
	}
	if mode != MODE_SILK_ONLY {
		var celt_frame_size int64 = (func() int64 {
			if F20 < frame_size {
				return F20
			}
			return frame_size
		}())
		if mode != st.Prev_mode && st.Prev_mode > 0 && st.Prev_redundancy == 0 {
			for {
				if opus_custom_decoder_ctl(celt_dec, OPUS_RESET_STATE) != OPUS_OK {
					return -3
				}
				if true {
					break
				}
			}
		}
		celt_ret = celt_decode_with_ec(celt_dec, func() *uint8 {
			if decode_fec != 0 {
				return nil
			}
			return data
		}(), int64(len_), pcm, celt_frame_size, &dec, celt_accum)
	} else {
		var silence [2]uint8 = [2]uint8{math.MaxUint8, math.MaxUint8}
		if celt_accum == 0 {
			for i = 0; i < frame_size*st.Channels; i++ {
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = 0
			}
		}
		if st.Prev_mode == MODE_HYBRID && (redundancy == 0 || celt_to_silk == 0 || st.Prev_redundancy == 0) {
			for {
				if opus_custom_decoder_ctl(celt_dec, CELT_SET_START_BAND_REQUEST, func() int64 {
					0 == 0
					return 0
				}()) != OPUS_OK {
					return -3
				}
				if true {
					break
				}
			}
			celt_decode_with_ec(celt_dec, &silence[0], 2, pcm, F2_5, nil, celt_accum)
		}
	}
	if mode != MODE_CELT_ONLY && celt_accum == 0 {
		for i = 0; i < frame_size*st.Channels; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i))) + opus_val16(float64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pcm_silk), unsafe.Sizeof(opus_int16(0))*uintptr(i))))*(1.0/32768.0))
		}
	}
	{
		var celt_mode *OpusCustomMode
		for {
			if opus_custom_decoder_ctl(celt_dec, CELT_GET_MODE_REQUEST, (**OpusCustomMode)(unsafe.Add(unsafe.Pointer(&celt_mode), unsafe.Sizeof((*OpusCustomMode)(nil))*uintptr(int64(uintptr(unsafe.Pointer(&celt_mode))-uintptr(unsafe.Pointer(&celt_mode))))))) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
		window = celt_mode.Window
	}
	if redundancy != 0 && celt_to_silk == 0 {
		for {
			if opus_custom_decoder_ctl(celt_dec, OPUS_RESET_STATE) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
		for {
			if opus_custom_decoder_ctl(celt_dec, CELT_SET_START_BAND_REQUEST, func() int64 {
				0 == 0
				return 0
			}()) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
		celt_decode_with_ec(celt_dec, (*uint8)(unsafe.Add(unsafe.Pointer(data), len_)), redundancy_bytes, redundant_audio, F5, nil, 0)
		for {
			if opus_custom_decoder_ctl(celt_dec, OPUS_GET_FINAL_RANGE_REQUEST, (*opus_uint32)(unsafe.Add(unsafe.Pointer(&redundant_rng), unsafe.Sizeof(opus_uint32(0))*uintptr(int64(uintptr(unsafe.Pointer(&redundant_rng))-uintptr(unsafe.Pointer(&redundant_rng))))))) != OPUS_OK {
				return -3
			}
			if true {
				break
			}
		}
		smooth_fade((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*(frame_size-F2_5)))), (*opus_val16)(unsafe.Add(unsafe.Pointer(redundant_audio), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*(frame_size-F2_5)))), F2_5, st.Channels, window, st.Fs)
	}
	if redundancy != 0 && celt_to_silk != 0 && (st.Prev_mode != MODE_SILK_ONLY || st.Prev_redundancy != 0) {
		for c = 0; c < st.Channels; c++ {
			for i = 0; i < F2_5; i++ {
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*i+c))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(redundant_audio), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*i+c)))
			}
		}
		smooth_fade((*opus_val16)(unsafe.Add(unsafe.Pointer(redundant_audio), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), F2_5, st.Channels, window, st.Fs)
	}
	if transition != 0 {
		if audiosize >= F5 {
			for i = 0; i < st.Channels*F2_5; i++ {
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_transition), unsafe.Sizeof(opus_val16(0))*uintptr(i)))
			}
			smooth_fade((*opus_val16)(unsafe.Add(unsafe.Pointer(pcm_transition), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*F2_5))), F2_5, st.Channels, window, st.Fs)
		} else {
			smooth_fade(pcm_transition, pcm, pcm, F2_5, st.Channels, window, st.Fs)
		}
	}
	if st.Decode_gain != 0 {
		var gain opus_val32
		gain = opus_val32(float32(math.Exp((float64(st.Decode_gain) * 0.000648814081) * 0.6931471805599453)))
		for i = 0; i < frame_size*st.Channels; i++ {
			var x opus_val32
			x = opus_val32((*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i)))) * opus_val16(gain))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(x)
		}
	}
	if len_ <= 1 {
		st.RangeFinal = 0
	} else {
		st.RangeFinal = dec.Rng ^ redundant_rng
	}
	st.Prev_mode = mode
	st.Prev_redundancy = int64(libc.BoolToInt(redundancy != 0 && celt_to_silk == 0))
	if celt_ret >= 0 {
		if _opus_false() != 0 {
			for {
				if true {
					break
				}
			}
		}
	}
	if celt_ret < 0 {
		return celt_ret
	}
	return audiosize
}
func opus_decode_native(st *OpusDecoder, data *uint8, len_ opus_int32, pcm *opus_val16, frame_size int64, decode_fec int64, self_delimited int64, packet_offset *opus_int32, soft_clip int64) int64 {
	var (
		i                      int64
		nb_samples             int64
		count                  int64
		offset                 int64
		toc                    uint8
		packet_frame_size      int64
		packet_bandwidth       int64
		packet_mode            int64
		packet_stream_channels int64
		size                   [48]opus_int16
	)
	if decode_fec < 0 || decode_fec > 1 {
		return -1
	}
	if (decode_fec != 0 || len_ == 0 || data == nil) && frame_size%int64(st.Fs/400) != 0 {
		return -1
	}
	if len_ == 0 || data == nil {
		var pcm_count int64 = 0
		for {
			{
				var ret int64
				ret = opus_decode_frame(st, nil, 0, (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(pcm_count*st.Channels))), frame_size-pcm_count, 0)
				if ret < 0 {
					return ret
				}
				pcm_count += ret
			}
			if pcm_count >= frame_size {
				break
			}
		}
		if _opus_false() != 0 {
			for {
				if true {
					break
				}
			}
		}
		st.Last_packet_duration = pcm_count
		return pcm_count
	} else if len_ < 0 {
		return -1
	}
	packet_mode = opus_packet_get_mode(data)
	packet_bandwidth = opus_packet_get_bandwidth(data)
	packet_frame_size = opus_packet_get_samples_per_frame(data, st.Fs)
	packet_stream_channels = opus_packet_get_nb_channels(data)
	count = opus_packet_parse_impl(data, len_, self_delimited, &toc, ([48]*uint8)(0), size, &offset, packet_offset)
	if count < 0 {
		return count
	}
	data = (*uint8)(unsafe.Add(unsafe.Pointer(data), offset))
	if decode_fec != 0 {
		var (
			duration_copy int64
			ret           int64
		)
		if frame_size < packet_frame_size || packet_mode == MODE_CELT_ONLY || st.Mode == MODE_CELT_ONLY {
			return opus_decode_native(st, nil, 0, pcm, frame_size, 0, 0, nil, soft_clip)
		}
		duration_copy = st.Last_packet_duration
		if frame_size-packet_frame_size != 0 {
			ret = opus_decode_native(st, nil, 0, pcm, frame_size-packet_frame_size, 0, 0, nil, soft_clip)
			if ret < 0 {
				st.Last_packet_duration = duration_copy
				return ret
			}
		}
		st.Mode = packet_mode
		st.Bandwidth = packet_bandwidth
		st.Frame_size = packet_frame_size
		st.Stream_channels = packet_stream_channels
		ret = opus_decode_frame(st, data, opus_int32(size[0]), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(st.Channels*(frame_size-packet_frame_size)))), packet_frame_size, 1)
		if ret < 0 {
			return ret
		} else {
			if _opus_false() != 0 {
				for {
					if true {
						break
					}
				}
			}
			st.Last_packet_duration = frame_size
			return frame_size
		}
	}
	if count*packet_frame_size > frame_size {
		return -2
	}
	st.Mode = packet_mode
	st.Bandwidth = packet_bandwidth
	st.Frame_size = packet_frame_size
	st.Stream_channels = packet_stream_channels
	nb_samples = 0
	for i = 0; i < count; i++ {
		var ret int64
		ret = opus_decode_frame(st, data, opus_int32(size[i]), (*opus_val16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_val16(0))*uintptr(nb_samples*st.Channels))), frame_size-nb_samples, 0)
		if ret < 0 {
			return ret
		}
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), size[i]))
		nb_samples += ret
	}
	st.Last_packet_duration = nb_samples
	if _opus_false() != 0 {
		for {
			if true {
				break
			}
		}
	}
	if soft_clip != 0 {
		opus_pcm_soft_clip((*float32)(unsafe.Pointer(pcm)), nb_samples, st.Channels, (*float32)(unsafe.Pointer(&st.Softclip_mem[0])))
	} else {
		st.Softclip_mem[0] = func() opus_val16 {
			p := &st.Softclip_mem[1]
			st.Softclip_mem[1] = 0
			return *p
		}()
	}
	return nb_samples
}
func opus_decode(st *OpusDecoder, data *uint8, len_ opus_int32, pcm *opus_int16, frame_size int64, decode_fec int64) int64 {
	var (
		out        *float32
		ret        int64
		i          int64
		nb_samples int64
	)
	if frame_size <= 0 {
		return -1
	}
	if data != nil && len_ > 0 && decode_fec == 0 {
		nb_samples = opus_decoder_get_nb_samples(st, [0]uint8(data), len_)
		if nb_samples > 0 {
			if frame_size < nb_samples {
				frame_size = frame_size
			} else {
				frame_size = nb_samples
			}
		} else {
			return -4
		}
	}
	out = (*float32)(libc.Malloc(int((frame_size * st.Channels) * int64(unsafe.Sizeof(float32(0))))))
	ret = opus_decode_native(st, data, len_, (*opus_val16)(unsafe.Pointer(out)), frame_size, decode_fec, 0, nil, 1)
	if ret > 0 {
		for i = 0; i < ret*st.Channels; i++ {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pcm), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = FLOAT2INT16(*(*float32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(float32(0))*uintptr(i))))
		}
	}
	return ret
}
func opus_decode_float(st *OpusDecoder, data *uint8, len_ opus_int32, pcm *opus_val16, frame_size int64, decode_fec int64) int64 {
	if frame_size <= 0 {
		return -1
	}
	return opus_decode_native(st, data, len_, pcm, frame_size, decode_fec, 0, nil, 0)
}
func opus_decoder_ctl(st *OpusDecoder, request int64, _rest ...interface{}) int64 {
	var (
		ret      int64 = OPUS_OK
		ap       libc.ArgList
		silk_dec unsafe.Pointer
		celt_dec *OpusCustomDecoder
	)
	silk_dec = unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Silk_dec_offset)
	celt_dec = (*OpusCustomDecoder)(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Pointer(st))), st.Celt_dec_offset))))
	ap.Start(request, _rest)
	switch request {
	case OPUS_GET_BANDWIDTH_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Bandwidth)
	case OPUS_GET_FINAL_RANGE_REQUEST:
		var value *opus_uint32 = ap.Arg().(*opus_uint32)
		if value == nil {
			goto bad_arg
		}
		*value = st.RangeFinal
	case OPUS_RESET_STATE:
		libc.MemSet(unsafe.Pointer((*byte)(unsafe.Pointer(&st.Stream_channels))), 0, int((unsafe.Sizeof(OpusDecoder{})-uintptr(int64(uintptr(unsafe.Pointer((*byte)(unsafe.Pointer(&st.Stream_channels))))-uintptr(unsafe.Pointer((*byte)(unsafe.Pointer(st)))))))*unsafe.Sizeof(byte(0))))
		opus_custom_decoder_ctl(celt_dec, OPUS_RESET_STATE)
		silk_InitDecoder(silk_dec)
		st.Stream_channels = st.Channels
		st.Frame_size = int64(st.Fs / 400)
	case OPUS_GET_SAMPLE_RATE_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = st.Fs
	case OPUS_GET_PITCH_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		if st.Prev_mode == MODE_CELT_ONLY {
			ret = opus_custom_decoder_ctl(celt_dec, OPUS_GET_PITCH_REQUEST, (*opus_int32)(unsafe.Add(unsafe.Pointer(value), unsafe.Sizeof(opus_int32(0))*uintptr(int64(uintptr(unsafe.Pointer(value))-uintptr(unsafe.Pointer(value)))))))
		} else {
			*value = opus_int32(st.DecControl.PrevPitchLag)
		}
	case OPUS_GET_GAIN_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Decode_gain)
	case OPUS_SET_GAIN_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < math.MinInt16 || value > math.MaxInt16 {
			goto bad_arg
		}
		st.Decode_gain = int64(value)
	case OPUS_GET_LAST_PACKET_DURATION_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		*value = opus_int32(st.Last_packet_duration)
	case OPUS_SET_PHASE_INVERSION_DISABLED_REQUEST:
		var value opus_int32 = ap.Arg().(opus_int32)
		if value < 0 || value > 1 {
			goto bad_arg
		}
		ret = opus_custom_decoder_ctl(celt_dec, OPUS_SET_PHASE_INVERSION_DISABLED_REQUEST, func() opus_int32 {
			value == 0
			return value
		}())
	case OPUS_GET_PHASE_INVERSION_DISABLED_REQUEST:
		var value *opus_int32 = ap.Arg().(*opus_int32)
		if value == nil {
			goto bad_arg
		}
		ret = opus_custom_decoder_ctl(celt_dec, OPUS_GET_PHASE_INVERSION_DISABLED_REQUEST, (*opus_int32)(unsafe.Add(unsafe.Pointer(value), unsafe.Sizeof(opus_int32(0))*uintptr(int64(uintptr(unsafe.Pointer(value))-uintptr(unsafe.Pointer(value)))))))
	default:
		ret = -5
	}
	ap.End()
	return ret
bad_arg:
	ap.End()
	return -1
}
func opus_decoder_destroy(st *OpusDecoder) {
	libc.Free(unsafe.Pointer(st))
}
func opus_packet_get_bandwidth(data *uint8) int64 {
	var bandwidth int64
	if int64(*data)&0x80 != 0 {
		bandwidth = OPUS_BANDWIDTH_MEDIUMBAND + ((int64(*data) >> 5) & 0x3)
		if bandwidth == OPUS_BANDWIDTH_MEDIUMBAND {
			bandwidth = OPUS_BANDWIDTH_NARROWBAND
		}
	} else if (int64(*data) & 0x60) == 0x60 {
		if (int64(*data) & 0x10) != 0 {
			bandwidth = OPUS_BANDWIDTH_FULLBAND
		} else {
			bandwidth = OPUS_BANDWIDTH_SUPERWIDEBAND
		}
	} else {
		bandwidth = OPUS_BANDWIDTH_NARROWBAND + ((int64(*data) >> 5) & 0x3)
	}
	return bandwidth
}
func opus_packet_get_nb_channels(data *uint8) int64 {
	if (int64(*data) & 0x4) != 0 {
		return 2
	}
	return 1
}
func opus_packet_get_nb_frames(packet [0]uint8, len_ opus_int32) int64 {
	var count int64
	if len_ < 1 {
		return -1
	}
	count = int64(packet[0]) & 0x3
	if count == 0 {
		return 1
	} else if count != 3 {
		return 2
	} else if len_ < 2 {
		return -4
	} else {
		return int64(packet[1]) & 0x3F
	}
}
func opus_packet_get_nb_samples(packet [0]uint8, len_ opus_int32, Fs opus_int32) int64 {
	var (
		samples int64
		count   int64 = opus_packet_get_nb_frames(packet, len_)
	)
	if count < 0 {
		return count
	}
	samples = count * opus_packet_get_samples_per_frame(&packet[0], Fs)
	if samples*25 > int64(Fs*3) {
		return -4
	} else {
		return samples
	}
}
func opus_decoder_get_nb_samples(dec *OpusDecoder, packet [0]uint8, len_ opus_int32) int64 {
	return opus_packet_get_nb_samples(packet, len_, dec.Fs)
}
