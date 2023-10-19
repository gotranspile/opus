package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func opus_pcm_soft_clip(_x *float32, N int64, C int64, declip_mem *float32) {
	var (
		c int64
		i int64
		x *float32
	)
	if C < 1 || N < 1 || _x == nil || declip_mem == nil {
		return
	}
	for i = 0; i < N*C; i++ {
		if (-2.0) > (func() float64 {
			if 2.0 < float64(*(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i)))) {
				return 2.0
			}
			return float64(*(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i))))
		}()) {
			*(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i))) = -2.0
		} else if 2.0 < float64(*(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i)))) {
			*(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i))) = 2.0
		} else {
			*(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i))) = *(*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(i)))
		}
	}
	for c = 0; c < C; c++ {
		var (
			a    float32
			x0   float32
			curr int64
		)
		x = (*float32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(float32(0))*uintptr(c)))
		a = *(*float32)(unsafe.Add(unsafe.Pointer(declip_mem), unsafe.Sizeof(float32(0))*uintptr(c)))
		for i = 0; i < N; i++ {
			if *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))*a >= 0 {
				break
			}
			*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) = *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) + a**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))
		}
		curr = 0
		x0 = *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*0))
		for {
			var (
				start    int64
				end      int64
				maxval   float32
				special  int64 = 0
				peak_pos int64
			)
			for i = curr; i < N; i++ {
				if *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) > 1 || *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) < float32(-1) {
					break
				}
			}
			if i == N {
				a = 0
				break
			}
			peak_pos = i
			start = func() int64 {
				end = i
				return end
			}()
			maxval = float32(math.Abs(float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))))))
			for start > 0 && *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr((start-1)*C))) >= 0 {
				start--
			}
			for end < N && *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(end*C))) >= 0 {
				if (float32(math.Abs(float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(end*C))))))) > maxval {
					maxval = float32(math.Abs(float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(end*C))))))
					peak_pos = end
				}
				end++
			}
			special = int64(libc.BoolToInt(start == 0 && *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*0)) >= 0))
			a = (maxval - 1) / (maxval * maxval)
			a += float32(float64(a) * 2.4e-07)
			if *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) > 0 {
				a = -a
			}
			for i = start; i < end; i++ {
				*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) = *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) + a**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))**(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))
			}
			if special != 0 && peak_pos >= 2 {
				var (
					delta  float32
					offset float32 = x0 - *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*0))
				)
				delta = offset / float32(peak_pos)
				for i = curr; i < peak_pos; i++ {
					offset -= delta
					*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) += offset
					if (-1.0) > (func() float64 {
						if 1.0 < float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))) {
							return 1.0
						}
						return float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))))
					}()) {
						*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) = -1.0
					} else if 1.0 < float64(*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))) {
						*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) = 1.0
					} else {
						*(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C))) = *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(i*C)))
					}
				}
			}
			curr = end
			if curr == N {
				break
			}
		}
		*(*float32)(unsafe.Add(unsafe.Pointer(declip_mem), unsafe.Sizeof(float32(0))*uintptr(c))) = a
	}
}
func encode_size(size int64, data *uint8) int64 {
	if size < 252 {
		*data = uint8(int8(size))
		return 1
	} else {
		*data = uint8(int8((size & 0x3) + 252))
		*(*uint8)(unsafe.Add(unsafe.Pointer(data), 1)) = uint8(int8((size - int64(*data)) >> 2))
		return 2
	}
}
func parse_size(data *uint8, len_ opus_int32, size *opus_int16) int64 {
	if len_ < 1 {
		*size = -1
		return -1
	} else if int64(*data) < 252 {
		*size = opus_int16(*data)
		return 1
	} else if len_ < 2 {
		*size = -1
		return -1
	} else {
		*size = opus_int16(int64(*(*uint8)(unsafe.Add(unsafe.Pointer(data), 1)))*4 + int64(*data))
		return 2
	}
}
func opus_packet_get_samples_per_frame(data *uint8, Fs opus_int32) int64 {
	var audiosize int64
	if int64(*data)&0x80 != 0 {
		audiosize = (int64(*data) >> 3) & 0x3
		audiosize = int64((Fs << opus_int32(audiosize)) / 400)
	} else if (int64(*data) & 0x60) == 0x60 {
		if (int64(*data) & 0x8) != 0 {
			audiosize = int64(Fs / 50)
		} else {
			audiosize = int64(Fs / 100)
		}
	} else {
		audiosize = (int64(*data) >> 3) & 0x3
		if audiosize == 3 {
			audiosize = int64(Fs * 60 / 1000)
		} else {
			audiosize = int64((Fs << opus_int32(audiosize)) / 100)
		}
	}
	return audiosize
}
func opus_packet_parse_impl(data *uint8, len_ opus_int32, self_delimited int64, out_toc *uint8, frames [48]*uint8, size [48]opus_int16, payload_offset *int64, packet_offset *opus_int32) int64 {
	var (
		i         int64
		bytes     int64
		count     int64
		cbr       int64
		ch        uint8
		toc       uint8
		framesize int64
		last_size opus_int32
		pad       opus_int32 = 0
		data0     *uint8     = data
	)
	if size == nil || len_ < 0 {
		return -1
	}
	if len_ == 0 {
		return -4
	}
	framesize = opus_packet_get_samples_per_frame(data, 48000)
	cbr = 0
	toc = *func() *uint8 {
		p := &data
		x := *p
		*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
		return x
	}()
	len_--
	last_size = len_
	switch int64(toc) & 0x3 {
	case 0:
		count = 1
	case 1:
		count = 2
		cbr = 1
		if self_delimited == 0 {
			if len_&0x1 != 0 {
				return -4
			}
			last_size = len_ / 2
			size[0] = opus_int16(last_size)
		}
	case 2:
		count = 2
		bytes = parse_size(data, len_, &size[0])
		len_ -= opus_int32(bytes)
		if size[0] < 0 || opus_int32(size[0]) > len_ {
			return -4
		}
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), bytes))
		last_size = len_ - opus_int32(size[0])
	default:
		if len_ < 1 {
			return -4
		}
		ch = *func() *uint8 {
			p := &data
			x := *p
			*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
		count = int64(ch) & 0x3F
		if count <= 0 || framesize*int64(opus_int32(count)) > 5760 {
			return -4
		}
		len_--
		if int64(ch)&0x40 != 0 {
			var p int64
			for {
				{
					var tmp int64
					if len_ <= 0 {
						return -4
					}
					p = int64(*func() *uint8 {
						p := &data
						x := *p
						*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}())
					len_--
					if p == math.MaxUint8 {
						tmp = 254
					} else {
						tmp = p
					}
					len_ -= opus_int32(tmp)
					pad += opus_int32(tmp)
				}
				if p != math.MaxUint8 {
					break
				}
			}
		}
		if len_ < 0 {
			return -4
		}
		cbr = int64(libc.BoolToInt((int64(ch) & 0x80) == 0))
		if cbr == 0 {
			last_size = len_
			for i = 0; i < count-1; i++ {
				bytes = parse_size(data, len_, &size[i])
				len_ -= opus_int32(bytes)
				if size[i] < 0 || opus_int32(size[i]) > len_ {
					return -4
				}
				data = (*uint8)(unsafe.Add(unsafe.Pointer(data), bytes))
				last_size -= opus_int32(bytes + int64(size[i]))
			}
			if last_size < 0 {
				return -4
			}
		} else if self_delimited == 0 {
			last_size = len_ / opus_int32(count)
			if last_size*opus_int32(count) != len_ {
				return -4
			}
			for i = 0; i < count-1; i++ {
				size[i] = opus_int16(last_size)
			}
		}
	}
	if self_delimited != 0 {
		bytes = parse_size(data, len_, (*opus_int16)(unsafe.Add(unsafe.Pointer(&size[count]), -int(unsafe.Sizeof(opus_int16(0))*1))))
		len_ -= opus_int32(bytes)
		if size[count-1] < 0 || opus_int32(size[count-1]) > len_ {
			return -4
		}
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), bytes))
		if cbr != 0 {
			if int64(size[count-1])*count > int64(len_) {
				return -4
			}
			for i = 0; i < count-1; i++ {
				size[i] = size[count-1]
			}
		} else if bytes+int64(size[count-1]) > int64(last_size) {
			return -4
		}
	} else {
		if last_size > 1275 {
			return -4
		}
		size[count-1] = opus_int16(last_size)
	}
	if payload_offset != nil {
		*payload_offset = int64(uintptr(unsafe.Pointer(data)) - uintptr(unsafe.Pointer(data0)))
	}
	for i = 0; i < count; i++ {
		if frames != nil {
			frames[i] = data
		}
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), size[i]))
	}
	if packet_offset != nil {
		*packet_offset = pad + opus_int32(int64(uintptr(unsafe.Pointer(data))-uintptr(unsafe.Pointer(data0))))
	}
	if out_toc != nil {
		*out_toc = toc
	}
	return count
}
func opus_packet_parse(data *uint8, len_ opus_int32, out_toc *uint8, frames [48]*uint8, size [48]opus_int16, payload_offset *int64) int64 {
	return opus_packet_parse_impl(data, len_, 0, out_toc, frames, size, payload_offset, nil)
}
