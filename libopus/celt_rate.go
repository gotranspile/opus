package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const MAX_PSEUDO = 40
const LOG_MAX_PSEUDO = 6
const CELT_MAX_PULSES = 128
const MAX_FINE_BITS = 8
const FINE_OFFSET = 21
const QTHETA_OFFSET = 4
const QTHETA_OFFSET_TWOPHASE = 16
const ALLOC_STEPS = 6

func get_pulses(i int64) int64 {
	if i < 8 {
		return i
	}
	return ((i & 7) + 8) << ((i >> 3) - 1)
}
func bits2pulses(m *OpusCustomMode, band int64, LM int64, bits int64) int64 {
	var (
		i     int64
		lo    int64
		hi    int64
		cache *uint8
	)
	LM++
	cache = (*uint8)(unsafe.Add(unsafe.Pointer(m.Cache.Bits), *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.Cache.Index), unsafe.Sizeof(opus_int16(0))*uintptr(LM*m.NbEBands+band)))))
	lo = 0
	hi = int64(*cache)
	bits--
	for i = 0; i < LOG_MAX_PSEUDO; i++ {
		var mid int64 = (lo + hi + 1) >> 1
		if int64(*(*uint8)(unsafe.Add(unsafe.Pointer(cache), mid))) >= bits {
			hi = mid
		} else {
			lo = mid
		}
	}
	if bits-(func() int64 {
		if lo == 0 {
			return -1
		}
		return int64(*(*uint8)(unsafe.Add(unsafe.Pointer(cache), lo)))
	}()) <= int64(*(*uint8)(unsafe.Add(unsafe.Pointer(cache), hi)))-bits {
		return lo
	} else {
		return hi
	}
}
func pulses2bits(m *OpusCustomMode, band int64, LM int64, pulses int64) int64 {
	var cache *uint8
	LM++
	cache = (*uint8)(unsafe.Add(unsafe.Pointer(m.Cache.Bits), *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.Cache.Index), unsafe.Sizeof(opus_int16(0))*uintptr(LM*m.NbEBands+band)))))
	if pulses == 0 {
		return 0
	}
	return int64(*(*uint8)(unsafe.Add(unsafe.Pointer(cache), pulses))) + 1
}

var LOG2_FRAC_TABLE [24]uint8 = [24]uint8{0, 8, 13, 16, 19, 21, 23, 24, 26, 27, 28, 29, 30, 31, 32, 32, 33, 34, 34, 35, 36, 36, 37, 37}

func interp_bits2pulses(m *OpusCustomMode, start int64, end int64, skip_start int64, bits1 *int64, bits2 *int64, thresh *int64, cap_ *int64, total opus_int32, _balance *opus_int32, skip_rsv int64, intensity *int64, intensity_rsv int64, dual_stereo *int64, dual_stereo_rsv int64, bits *int64, ebits *int64, fine_priority *int64, C int64, LM int64, ec *ec_ctx, encode int64, prev int64, signalBandwidth int64) int64 {
	var (
		psum        opus_int32
		lo          int64
		hi          int64
		i           int64
		j           int64
		logM        int64
		stereo      int64
		codedBands  int64 = -1
		alloc_floor int64
		left        opus_int32
		percoeff    opus_int32
		done        int64
		balance     opus_int32
	)
	alloc_floor = C << BITRES
	stereo = int64(libc.BoolToInt(C > 1))
	logM = LM << BITRES
	lo = 0
	hi = 1 << ALLOC_STEPS
	for i = 0; i < ALLOC_STEPS; i++ {
		var mid int64 = (lo + hi) >> 1
		psum = 0
		done = 0
		for j = end; func() int64 {
			p := &j
			x := *p
			*p--
			return x
		}() > start; {
			var tmp int64 = *(*int64)(unsafe.Add(unsafe.Pointer(bits1), unsafe.Sizeof(int64(0))*uintptr(j))) + (mid * int64(opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(bits2), unsafe.Sizeof(int64(0))*uintptr(j))))) >> ALLOC_STEPS)
			if tmp >= *(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j))) || done != 0 {
				done = 1
				if tmp < (*(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j)))) {
					psum += opus_int32(tmp)
				} else {
					psum += opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j))))
				}
			} else {
				if tmp >= alloc_floor {
					psum += opus_int32(alloc_floor)
				}
			}
		}
		if psum > total {
			hi = mid
		} else {
			lo = mid
		}
	}
	psum = 0
	done = 0
	for j = end; func() int64 {
		p := &j
		x := *p
		*p--
		return x
	}() > start; {
		var tmp int64 = *(*int64)(unsafe.Add(unsafe.Pointer(bits1), unsafe.Sizeof(int64(0))*uintptr(j))) + int64(opus_int32(lo)*opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(bits2), unsafe.Sizeof(int64(0))*uintptr(j))))>>ALLOC_STEPS)
		if tmp < *(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j))) && done == 0 {
			if tmp >= alloc_floor {
				tmp = alloc_floor
			} else {
				tmp = 0
			}
		} else {
			done = 1
		}
		if tmp < (*(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j)))) {
			tmp = tmp
		} else {
			tmp = *(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j)))
		}
		*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) = tmp
		psum += opus_int32(tmp)
	}
	for codedBands = end; ; codedBands-- {
		var (
			band_width int64
			band_bits  int64
			rem        int64
		)
		j = codedBands - 1
		if j <= skip_start {
			total += opus_int32(skip_rsv)
			break
		}
		left = total - psum
		percoeff = opus_int32(celt_udiv(opus_uint32(left), opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(codedBands)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(start))))))
		left -= opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(codedBands)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(start)))) * percoeff
		if (left - opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(start))))) > 0 {
			rem = int64(left - opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(start)))))
		} else {
			rem = 0
		}
		band_width = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(codedBands))) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))
		band_bits = *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) + int64(percoeff*opus_int32(band_width)) + rem
		if band_bits >= (func() int64 {
			if (*(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j)))) > (alloc_floor + (1 << BITRES)) {
				return *(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j)))
			}
			return alloc_floor + (1 << BITRES)
		}()) {
			if encode != 0 {
				var depth_threshold int64
				if codedBands > 17 {
					if j < prev {
						depth_threshold = 7
					} else {
						depth_threshold = 9
					}
				} else {
					depth_threshold = 0
				}
				if codedBands <= start+2 || band_bits > (depth_threshold*band_width<<LM<<BITRES)>>4 && j <= signalBandwidth {
					ec_enc_bit_logp((*ec_enc)(unsafe.Pointer(ec)), 1, 1)
					break
				}
				ec_enc_bit_logp((*ec_enc)(unsafe.Pointer(ec)), 0, 1)
			} else if ec_dec_bit_logp((*ec_dec)(unsafe.Pointer(ec)), 1) != 0 {
				break
			}
			psum += opus_int32(1 << BITRES)
			band_bits -= 1 << BITRES
		}
		psum -= opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) + intensity_rsv)
		if intensity_rsv > 0 {
			intensity_rsv = int64(LOG2_FRAC_TABLE[j-start])
		}
		psum += opus_int32(intensity_rsv)
		if band_bits >= alloc_floor {
			psum += opus_int32(alloc_floor)
			*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) = alloc_floor
		} else {
			*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) = 0
		}
	}
	if intensity_rsv > 0 {
		if encode != 0 {
			if (*intensity) < codedBands {
				*intensity = *intensity
			} else {
				*intensity = codedBands
			}
			ec_enc_uint((*ec_enc)(unsafe.Pointer(ec)), opus_uint32(*intensity-start), opus_uint32(codedBands+1-start))
		} else {
			*intensity = start + int64(ec_dec_uint((*ec_dec)(unsafe.Pointer(ec)), opus_uint32(codedBands+1-start)))
		}
	} else {
		*intensity = 0
	}
	if *intensity <= start {
		total += opus_int32(dual_stereo_rsv)
		dual_stereo_rsv = 0
	}
	if dual_stereo_rsv > 0 {
		if encode != 0 {
			ec_enc_bit_logp((*ec_enc)(unsafe.Pointer(ec)), *dual_stereo, 1)
		} else {
			*dual_stereo = ec_dec_bit_logp((*ec_dec)(unsafe.Pointer(ec)), 1)
		}
	} else {
		*dual_stereo = 0
	}
	left = total - psum
	percoeff = opus_int32(celt_udiv(opus_uint32(left), opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(codedBands)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(start))))))
	left -= opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(codedBands)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(start)))) * percoeff
	for j = start; j < codedBands; j++ {
		*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) += int64(percoeff) * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))
	}
	for j = start; j < codedBands; j++ {
		var tmp int64 = int64(func() opus_int32 {
			if left < opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j)))) {
				return left
			}
			return opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))
		}())
		*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) += tmp
		left -= opus_int32(tmp)
	}
	balance = 0
	for j = start; j < codedBands; j++ {
		var (
			N0     int64
			N      int64
			den    int64
			offset int64
			NClogN int64
			excess opus_int32
			bit    opus_int32
		)
		N0 = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))
		N = N0 << LM
		bit = opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j)))) + balance
		if N > 1 {
			if (bit - opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j))))) > 0 {
				excess = bit - opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j))))
			} else {
				excess = 0
			}
			*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) = int64(bit - excess)
			den = C*N + (func() int64 {
				if C == 2 && N > 2 && *dual_stereo == 0 && j < *intensity {
					return 1
				}
				return 0
			}())
			NClogN = den * (int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.LogN), unsafe.Sizeof(opus_int16(0))*uintptr(j)))) + logM)
			offset = (NClogN >> 1) - den*FINE_OFFSET
			if N == 2 {
				offset += den << BITRES >> 2
			}
			if *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j)))+offset < den*2<<BITRES {
				offset += NClogN >> 2
			} else if *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j)))+offset < den*3<<BITRES {
				offset += NClogN >> 3
			}
			if 0 > (*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) + offset + (den << (BITRES - 1))) {
				*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = 0
			} else {
				*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) + offset + (den << (BITRES - 1))
			}
			*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = int64(celt_udiv(opus_uint32(*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j)))), opus_uint32(den)) >> BITRES)
			if C**(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) > (*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) >> BITRES) {
				*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) >> stereo >> BITRES
			}
			if (*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j)))) < MAX_FINE_BITS {
				*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = *(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j)))
			} else {
				*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = MAX_FINE_BITS
			}
			*(*int64)(unsafe.Add(unsafe.Pointer(fine_priority), unsafe.Sizeof(int64(0))*uintptr(j))) = int64(libc.BoolToInt(*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j)))*(den<<BITRES) >= *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j)))+offset))
			*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) -= C * *(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) << BITRES
		} else {
			if 0 > (bit - opus_int32(C<<BITRES)) {
				excess = 0
			} else {
				excess = bit - opus_int32(C<<BITRES)
			}
			*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) = int64(bit - excess)
			*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = 0
			*(*int64)(unsafe.Add(unsafe.Pointer(fine_priority), unsafe.Sizeof(int64(0))*uintptr(j))) = 1
		}
		if excess > 0 {
			var (
				extra_fine int64
				extra_bits int64
			)
			if (excess >> opus_int32(stereo+BITRES)) < opus_int32(MAX_FINE_BITS-*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j)))) {
				extra_fine = int64(excess >> opus_int32(stereo+BITRES))
			} else {
				extra_fine = MAX_FINE_BITS - *(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j)))
			}
			*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) += extra_fine
			extra_bits = extra_fine * C << BITRES
			*(*int64)(unsafe.Add(unsafe.Pointer(fine_priority), unsafe.Sizeof(int64(0))*uintptr(j))) = int64(libc.BoolToInt(extra_bits >= int64(excess-balance)))
			excess -= opus_int32(extra_bits)
		}
		balance = excess
	}
	*_balance = balance
	for ; j < end; j++ {
		*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) = *(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) >> stereo >> BITRES
		*(*int64)(unsafe.Add(unsafe.Pointer(bits), unsafe.Sizeof(int64(0))*uintptr(j))) = 0
		*(*int64)(unsafe.Add(unsafe.Pointer(fine_priority), unsafe.Sizeof(int64(0))*uintptr(j))) = int64(libc.BoolToInt(*(*int64)(unsafe.Add(unsafe.Pointer(ebits), unsafe.Sizeof(int64(0))*uintptr(j))) < 1))
	}
	return codedBands
}
func clt_compute_allocation(m *OpusCustomMode, start int64, end int64, offsets *int64, cap_ *int64, alloc_trim int64, intensity *int64, dual_stereo *int64, total opus_int32, balance *opus_int32, pulses *int64, ebits *int64, fine_priority *int64, C int64, LM int64, ec *ec_ctx, encode int64, prev int64, signalBandwidth int64) int64 {
	var (
		lo              int64
		hi              int64
		len_            int64
		j               int64
		codedBands      int64
		skip_start      int64
		skip_rsv        int64
		intensity_rsv   int64
		dual_stereo_rsv int64
		bits1           *int64
		bits2           *int64
		thresh          *int64
		trim_offset     *int64
	)
	if total > 0 {
		total = total
	} else {
		total = 0
	}
	len_ = m.NbEBands
	skip_start = start
	if total >= opus_int32(1<<BITRES) {
		skip_rsv = 1 << BITRES
	} else {
		skip_rsv = 0
	}
	total -= opus_int32(skip_rsv)
	intensity_rsv = func() int64 {
		dual_stereo_rsv = 0
		return dual_stereo_rsv
	}()
	if C == 2 {
		intensity_rsv = int64(LOG2_FRAC_TABLE[end-start])
		if intensity_rsv > int64(total) {
			intensity_rsv = 0
		} else {
			total -= opus_int32(intensity_rsv)
			if total >= opus_int32(1<<BITRES) {
				dual_stereo_rsv = 1 << BITRES
			} else {
				dual_stereo_rsv = 0
			}
			total -= opus_int32(dual_stereo_rsv)
		}
	}
	bits1 = (*int64)(libc.Malloc(int(len_ * int64(unsafe.Sizeof(int64(0))))))
	bits2 = (*int64)(libc.Malloc(int(len_ * int64(unsafe.Sizeof(int64(0))))))
	thresh = (*int64)(libc.Malloc(int(len_ * int64(unsafe.Sizeof(int64(0))))))
	trim_offset = (*int64)(libc.Malloc(int(len_ * int64(unsafe.Sizeof(int64(0))))))
	for j = start; j < end; j++ {
		if (C << BITRES) > ((int64((*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))*3) << LM << BITRES) >> 4) {
			*(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j))) = C << BITRES
		} else {
			*(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j))) = (int64((*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))*3) << LM << BITRES) >> 4
		}
		*(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j))) = C * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j)))) * (alloc_trim - 5 - LM) * (end - j - 1) * (1 << (LM + BITRES)) >> 6
		if int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))<<LM == 1 {
			*(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j))) -= C << BITRES
		}
	}
	lo = 1
	hi = m.NbAllocVectors - 1
	for {
		{
			var (
				done int64 = 0
				psum int64 = 0
				mid  int64 = (lo + hi) >> 1
			)
			for j = end; func() int64 {
				p := &j
				x := *p
				*p--
				return x
			}() > start; {
				var (
					bitsj int64
					N     int64 = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))
				)
				bitsj = C * N * int64(*(*uint8)(unsafe.Add(unsafe.Pointer(m.AllocVectors), mid*len_+j))) << LM >> 2
				if bitsj > 0 {
					if 0 > (bitsj + *(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j)))) {
						bitsj = 0
					} else {
						bitsj = bitsj + *(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j)))
					}
				}
				bitsj += *(*int64)(unsafe.Add(unsafe.Pointer(offsets), unsafe.Sizeof(int64(0))*uintptr(j)))
				if bitsj >= *(*int64)(unsafe.Add(unsafe.Pointer(thresh), unsafe.Sizeof(int64(0))*uintptr(j))) || done != 0 {
					done = 1
					if bitsj < (*(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j)))) {
						psum += bitsj
					} else {
						psum += *(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j)))
					}
				} else {
					if bitsj >= C<<BITRES {
						psum += C << BITRES
					}
				}
			}
			if psum > int64(total) {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}
		if lo > hi {
			break
		}
	}
	hi = func() int64 {
		p := &lo
		x := *p
		*p--
		return x
	}()
	for j = start; j < end; j++ {
		var (
			bits1j int64
			bits2j int64
			N      int64 = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(j))))
		)
		bits1j = C * N * int64(*(*uint8)(unsafe.Add(unsafe.Pointer(m.AllocVectors), lo*len_+j))) << LM >> 2
		if hi >= m.NbAllocVectors {
			bits2j = *(*int64)(unsafe.Add(unsafe.Pointer(cap_), unsafe.Sizeof(int64(0))*uintptr(j)))
		} else {
			bits2j = C * N * int64(*(*uint8)(unsafe.Add(unsafe.Pointer(m.AllocVectors), hi*len_+j))) << LM >> 2
		}
		if bits1j > 0 {
			if 0 > (bits1j + *(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j)))) {
				bits1j = 0
			} else {
				bits1j = bits1j + *(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j)))
			}
		}
		if bits2j > 0 {
			if 0 > (bits2j + *(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j)))) {
				bits2j = 0
			} else {
				bits2j = bits2j + *(*int64)(unsafe.Add(unsafe.Pointer(trim_offset), unsafe.Sizeof(int64(0))*uintptr(j)))
			}
		}
		if lo > 0 {
			bits1j += *(*int64)(unsafe.Add(unsafe.Pointer(offsets), unsafe.Sizeof(int64(0))*uintptr(j)))
		}
		bits2j += *(*int64)(unsafe.Add(unsafe.Pointer(offsets), unsafe.Sizeof(int64(0))*uintptr(j)))
		if *(*int64)(unsafe.Add(unsafe.Pointer(offsets), unsafe.Sizeof(int64(0))*uintptr(j))) > 0 {
			skip_start = j
		}
		if 0 > (bits2j - bits1j) {
			bits2j = 0
		} else {
			bits2j = bits2j - bits1j
		}
		*(*int64)(unsafe.Add(unsafe.Pointer(bits1), unsafe.Sizeof(int64(0))*uintptr(j))) = bits1j
		*(*int64)(unsafe.Add(unsafe.Pointer(bits2), unsafe.Sizeof(int64(0))*uintptr(j))) = bits2j
	}
	codedBands = interp_bits2pulses(m, start, end, skip_start, bits1, bits2, thresh, cap_, total, balance, skip_rsv, intensity, intensity_rsv, dual_stereo, dual_stereo_rsv, pulses, ebits, fine_priority, C, LM, ec, encode, prev, signalBandwidth)
	return codedBands
}
