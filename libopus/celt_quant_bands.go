package libopus

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

var eMeans [25]opus_val16 = [25]opus_val16{opus_val16(6.4375), opus_val16(6.25), opus_val16(5.75), opus_val16(5.3125), opus_val16(5.0625), opus_val16(4.8125), opus_val16(4.5), opus_val16(4.375), opus_val16(4.875), opus_val16(4.6875), opus_val16(4.5625), opus_val16(4.4375), opus_val16(4.875), opus_val16(4.625), opus_val16(4.3125), opus_val16(4.5), opus_val16(4.375), opus_val16(4.625), opus_val16(4.75), opus_val16(4.4375), opus_val16(3.75), opus_val16(3.75), opus_val16(3.75), opus_val16(3.75), opus_val16(3.75)}
var pred_coef [4]opus_val16 = [4]opus_val16{opus_val16(29440 / 32768.0), opus_val16(26112 / 32768.0), opus_val16(21248 / 32768.0), opus_val16(16384 / 32768.0)}
var beta_coef [4]opus_val16 = [4]opus_val16{opus_val16(30147 / 32768.0), opus_val16(22282 / 32768.0), opus_val16(12124 / 32768.0), opus_val16(6554 / 32768.0)}
var beta_intra opus_val16 = opus_val16(4915 / 32768.0)
var e_prob_model [4][2][42]uint8 = [4][2][42]uint8{{{72, math.MaxInt8, 65, 129, 66, 128, 65, 128, 64, 128, 62, 128, 64, 128, 64, 128, 92, 78, 92, 79, 92, 78, 90, 79, 116, 41, 115, 40, 114, 40, 132, 26, 132, 26, 145, 17, 161, 12, 176, 10, 177, 11}, {24, 179, 48, 138, 54, 135, 54, 132, 53, 134, 56, 133, 55, 132, 55, 132, 61, 114, 70, 96, 74, 88, 75, 88, 87, 74, 89, 66, 91, 67, 100, 59, 108, 50, 120, 40, 122, 37, 97, 43, 78, 50}}, {{83, 78, 84, 81, 88, 75, 86, 74, 87, 71, 90, 73, 93, 74, 93, 74, 109, 40, 114, 36, 117, 34, 117, 34, 143, 17, 145, 18, 146, 19, 162, 12, 165, 10, 178, 7, 189, 6, 190, 8, 177, 9}, {23, 178, 54, 115, 63, 102, 66, 98, 69, 99, 74, 89, 71, 91, 73, 91, 78, 89, 86, 80, 92, 66, 93, 64, 102, 59, 103, 60, 104, 60, 117, 52, 123, 44, 138, 35, 133, 31, 97, 38, 77, 45}}, {{61, 90, 93, 60, 105, 42, 107, 41, 110, 45, 116, 38, 113, 38, 112, 38, 124, 26, 132, 27, 136, 19, 140, 20, 155, 14, 159, 16, 158, 18, 170, 13, 177, 10, 187, 8, 192, 6, 175, 9, 159, 10}, {21, 178, 59, 110, 71, 86, 75, 85, 84, 83, 91, 66, 88, 73, 87, 72, 92, 75, 98, 72, 105, 58, 107, 54, 115, 52, 114, 55, 112, 56, 129, 51, 132, 40, 150, 33, 140, 29, 98, 35, 77, 42}}, {{42, 121, 96, 66, 108, 43, 111, 40, 117, 44, 123, 32, 120, 36, 119, 33, math.MaxInt8, 33, 134, 34, 139, 21, 147, 23, 152, 20, 158, 25, 154, 26, 166, 21, 173, 16, 184, 13, 184, 10, 150, 13, 139, 15}, {22, 178, 63, 114, 74, 82, 84, 83, 92, 82, 103, 62, 96, 72, 96, 67, 101, 73, 107, 72, 113, 55, 118, 52, 125, 52, 118, 52, 117, 55, 135, 49, 137, 39, 157, 32, 145, 29, 97, 33, 77, 40}}}
var small_energy_icdf [3]uint8 = [3]uint8{2, 1, 0}

func loss_distortion(eBands *opus_val16, oldEBands *opus_val16, start int64, end int64, len_ int64, C int64) opus_val32 {
	var (
		c    int64
		i    int64
		dist opus_val32 = 0
	)
	c = 0
	for {
		for i = start; i < end; i++ {
			var d opus_val16 = ((*(*opus_val16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*len_)))) - (*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*len_)))))
			dist = dist + opus_val32(d)*opus_val32(d)
		}
		if func() int64 {
			p := &c
			*p++
			return *p
		}() >= C {
			break
		}
	}
	if 200 < dist {
		return 200
	}
	return dist
}
func quant_coarse_energy_impl(m *OpusCustomMode, start int64, end int64, eBands *opus_val16, oldEBands *opus_val16, budget opus_int32, tell opus_int32, prob_model *uint8, error *opus_val16, enc *ec_enc, C int64, LM int64, intra int64, max_decay opus_val16, lfe int64) int64 {
	var (
		i       int64
		c       int64
		badness int64         = 0
		prev    [2]opus_val32 = [2]opus_val32{}
		coef    opus_val16
		beta    opus_val16
	)
	if tell+3 <= budget {
		ec_enc_bit_logp(enc, intra, 3)
	}
	if intra != 0 {
		coef = 0
		beta = beta_intra
	} else {
		beta = beta_coef[LM]
		coef = pred_coef[LM]
	}
	for i = start; i < end; i++ {
		c = 0
		for {
			{
				var (
					bits_left   int64
					qi          int64
					qi0         int64
					q           opus_val32
					x           opus_val16
					f           opus_val32
					tmp         opus_val32
					oldE        opus_val16
					decay_bound opus_val16
				)
				x = *(*opus_val16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))
				if (-9.0) > float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))) {
					oldE = opus_val16(-9.0)
				} else {
					oldE = *(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))
				}
				f = opus_val32(x - coef*oldE - opus_val16(prev[c]))
				qi = int64(math.Floor(float64(f) + 0.5))
				decay_bound = opus_val16((func() float64 {
					if (-28.0) > float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))) {
						return -28.0
					}
					return float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))))
				}()) - float64(max_decay))
				if qi < 0 && x < decay_bound {
					qi += int64(decay_bound - x)
					if qi > 0 {
						qi = 0
					}
				}
				qi0 = qi
				tell = opus_int32(ec_tell((*ec_ctx)(unsafe.Pointer(enc))))
				bits_left = int64(budget - tell - opus_int32(C*3*(end-i)))
				if i != start && bits_left < 30 {
					if bits_left < 24 {
						if 1 < qi {
							qi = 1
						} else {
							qi = qi
						}
					}
					if bits_left < 16 {
						if int64(-1) > qi {
							qi = -1
						} else {
							qi = qi
						}
					}
				}
				if lfe != 0 && i >= 2 {
					if qi < 0 {
						qi = qi
					} else {
						qi = 0
					}
				}
				if budget-tell >= 15 {
					var pi int64
					pi = (func() int64 {
						if i < 20 {
							return i
						}
						return 20
					}()) * 2
					ec_laplace_encode(enc, &qi, uint64(int64(*(*uint8)(unsafe.Add(unsafe.Pointer(prob_model), pi)))<<7), int64(*(*uint8)(unsafe.Add(unsafe.Pointer(prob_model), pi+1)))<<6)
				} else if budget-tell >= 2 {
					if int64(-1) > (func() int64 {
						if qi < 1 {
							return qi
						}
						return 1
					}()) {
						qi = -1
					} else if qi < 1 {
						qi = qi
					} else {
						qi = 1
					}
					ec_enc_icdf(enc, qi*2^(-int64(libc.BoolToInt(qi < 0))), &small_energy_icdf[0], 2)
				} else if budget-tell >= 1 {
					if 0 < qi {
						qi = 0
					} else {
						qi = qi
					}
					ec_enc_bit_logp(enc, -qi, 1)
				} else {
					qi = -1
				}
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(error), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) = opus_val16(f - opus_val32(qi))
				badness += cmath.Abs(qi0 - qi)
				q = opus_val32(qi)
				tmp = (opus_val32(coef) * opus_val32(oldE)) + prev[c] + q
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) = opus_val16(tmp)
				prev[c] = prev[c] + q - opus_val32(beta)*q
			}
			if func() int64 {
				p := &c
				*p++
				return *p
			}() >= C {
				break
			}
		}
	}
	if lfe != 0 {
		return 0
	}
	return badness
}
func quant_coarse_energy(m *OpusCustomMode, start int64, end int64, effEnd int64, eBands *opus_val16, oldEBands *opus_val16, budget opus_uint32, error *opus_val16, enc *ec_enc, C int64, LM int64, nbAvailableBytes int64, force_intra int64, delayedIntra *opus_val32, two_pass int64, loss_rate int64, lfe int64) {
	var (
		intra           int64
		max_decay       opus_val16
		oldEBands_intra *opus_val16
		error_intra     *opus_val16
		enc_start_state ec_enc
		tell            opus_uint32
		badness1        int64 = 0
		intra_bias      opus_int32
		new_distortion  opus_val32
	)
	intra = int64(libc.BoolToInt(force_intra != 0 || two_pass == 0 && *delayedIntra > opus_val32(C*2*(end-start)) && nbAvailableBytes > (end-start)*C))
	intra_bias = opus_int32((opus_val32(budget) * *delayedIntra * opus_val32(loss_rate)) / opus_val32(C*512))
	new_distortion = loss_distortion(eBands, oldEBands, start, effEnd, m.NbEBands, C)
	tell = opus_uint32(ec_tell((*ec_ctx)(unsafe.Pointer(enc))))
	if tell+3 > budget {
		two_pass = func() int64 {
			intra = 0
			return intra
		}()
	}
	max_decay = opus_val16(16.0)
	if end-start > 10 {
		if float64(max_decay) < (float64(nbAvailableBytes) * 0.125) {
			max_decay = max_decay
		} else {
			max_decay = opus_val16(float64(nbAvailableBytes) * 0.125)
		}
	}
	if lfe != 0 {
		max_decay = opus_val16(3.0)
	}
	enc_start_state = *enc
	oldEBands_intra = (*opus_val16)(libc.Malloc(int((C * m.NbEBands) * int64(unsafe.Sizeof(opus_val16(0))))))
	error_intra = (*opus_val16)(libc.Malloc(int((C * m.NbEBands) * int64(unsafe.Sizeof(opus_val16(0))))))
	libc.MemCpy(unsafe.Pointer(oldEBands_intra), unsafe.Pointer(oldEBands), int((C*m.NbEBands)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(oldEBands_intra))-uintptr(unsafe.Pointer(oldEBands))))*0))
	if two_pass != 0 || intra != 0 {
		badness1 = quant_coarse_energy_impl(m, start, end, eBands, oldEBands_intra, opus_int32(budget), opus_int32(tell), &e_prob_model[LM][1][0], error_intra, enc, C, LM, 1, max_decay, lfe)
	}
	if intra == 0 {
		var (
			intra_buf       *uint8
			enc_intra_state ec_enc
			tell_intra      opus_int32
			nstart_bytes    opus_uint32
			nintra_bytes    opus_uint32
			save_bytes      opus_uint32
			badness2        int64
			intra_bits      *uint8
		)
		tell_intra = opus_int32(ec_tell_frac((*ec_ctx)(unsafe.Pointer(enc))))
		enc_intra_state = *enc
		nstart_bytes = ec_range_bytes((*ec_ctx)(unsafe.Pointer(&enc_start_state)))
		nintra_bytes = ec_range_bytes((*ec_ctx)(unsafe.Pointer(&enc_intra_state)))
		intra_buf = (*uint8)(unsafe.Add(unsafe.Pointer(ec_get_buffer((*ec_ctx)(unsafe.Pointer(&enc_intra_state)))), nstart_bytes))
		save_bytes = nintra_bytes - nstart_bytes
		if save_bytes == 0 {
			save_bytes = ALLOC_NONE
		}
		intra_bits = (*uint8)(libc.Malloc(int(save_bytes * opus_uint32(unsafe.Sizeof(uint8(0))))))
		libc.MemCpy(unsafe.Pointer(intra_bits), unsafe.Pointer(intra_buf), int((nintra_bytes-nstart_bytes)*opus_uint32(unsafe.Sizeof(uint8(0)))+opus_uint32((int64(uintptr(unsafe.Pointer(intra_bits))-uintptr(unsafe.Pointer(intra_buf))))*0)))
		*enc = enc_start_state
		badness2 = quant_coarse_energy_impl(m, start, end, eBands, oldEBands, opus_int32(budget), opus_int32(tell), &e_prob_model[LM][intra][0], error, enc, C, LM, 0, max_decay, lfe)
		if two_pass != 0 && (badness1 < badness2 || badness1 == badness2 && (opus_int32(ec_tell_frac((*ec_ctx)(unsafe.Pointer(enc)))))+intra_bias > tell_intra) {
			*enc = enc_intra_state
			libc.MemCpy(unsafe.Pointer(intra_buf), unsafe.Pointer(intra_bits), int((nintra_bytes-nstart_bytes)*opus_uint32(unsafe.Sizeof(uint8(0)))+opus_uint32((int64(uintptr(unsafe.Pointer(intra_buf))-uintptr(unsafe.Pointer(intra_bits))))*0)))
			libc.MemCpy(unsafe.Pointer(oldEBands), unsafe.Pointer(oldEBands_intra), int((C*m.NbEBands)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(oldEBands))-uintptr(unsafe.Pointer(oldEBands_intra))))*0))
			libc.MemCpy(unsafe.Pointer(error), unsafe.Pointer(error_intra), int((C*m.NbEBands)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(error))-uintptr(unsafe.Pointer(error_intra))))*0))
			intra = 1
		}
	} else {
		libc.MemCpy(unsafe.Pointer(oldEBands), unsafe.Pointer(oldEBands_intra), int((C*m.NbEBands)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(oldEBands))-uintptr(unsafe.Pointer(oldEBands_intra))))*0))
		libc.MemCpy(unsafe.Pointer(error), unsafe.Pointer(error_intra), int((C*m.NbEBands)*int64(unsafe.Sizeof(opus_val16(0)))+(int64(uintptr(unsafe.Pointer(error))-uintptr(unsafe.Pointer(error_intra))))*0))
	}
	if intra != 0 {
		*delayedIntra = new_distortion
	} else {
		*delayedIntra = opus_val32((((pred_coef[LM]) * (pred_coef[LM])) * opus_val16(*delayedIntra)) + opus_val16(new_distortion))
	}
}
func quant_fine_energy(m *OpusCustomMode, start int64, end int64, oldEBands *opus_val16, error *opus_val16, fine_quant *int64, enc *ec_enc, C int64) {
	var (
		i int64
		c int64
	)
	for i = start; i < end; i++ {
		var frac opus_int16 = opus_int16(1 << *(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i))))
		if *(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i))) <= 0 {
			continue
		}
		c = 0
		for {
			{
				var (
					q2     int64
					offset opus_val16
				)
				q2 = int64(math.Floor((float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(error), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))) + 0.5) * float64(frac)))
				if q2 > int64(frac-1) {
					q2 = int64(frac - 1)
				}
				if q2 < 0 {
					q2 = 0
				}
				ec_enc_bits(enc, opus_uint32(q2), uint64(*(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i)))))
				offset = opus_val16((float64(q2)+0.5)*float64(1<<(14-*(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i)))))*(1.0/16384) - 0.5)
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) += offset
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(error), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) -= offset
			}
			if func() int64 {
				p := &c
				*p++
				return *p
			}() >= C {
				break
			}
		}
	}
}
func quant_energy_finalise(m *OpusCustomMode, start int64, end int64, oldEBands *opus_val16, error *opus_val16, fine_quant *int64, fine_priority *int64, bits_left int64, enc *ec_enc, C int64) {
	var (
		i    int64
		prio int64
		c    int64
	)
	for prio = 0; prio < 2; prio++ {
		for i = start; i < end && bits_left >= C; i++ {
			if *(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i))) >= MAX_FINE_BITS || *(*int64)(unsafe.Add(unsafe.Pointer(fine_priority), unsafe.Sizeof(int64(0))*uintptr(i))) != prio {
				continue
			}
			c = 0
			for {
				{
					var (
						q2     int64
						offset opus_val16
					)
					if *(*opus_val16)(unsafe.Add(unsafe.Pointer(error), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) < 0 {
						q2 = 0
					} else {
						q2 = 1
					}
					ec_enc_bits(enc, opus_uint32(q2), 1)
					offset = opus_val16((float64(q2) - 0.5) * float64(1<<(14-*(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i)))-1)) * (1.0 / 16384))
					*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) += offset
					*(*opus_val16)(unsafe.Add(unsafe.Pointer(error), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) -= offset
					bits_left--
				}
				if func() int64 {
					p := &c
					*p++
					return *p
				}() >= C {
					break
				}
			}
		}
	}
}
func unquant_coarse_energy(m *OpusCustomMode, start int64, end int64, oldEBands *opus_val16, intra int64, dec *ec_dec, C int64, LM int64) {
	var (
		prob_model *uint8 = &e_prob_model[LM][intra][0]
		i          int64
		c          int64
		prev       [2]opus_val32 = [2]opus_val32{}
		coef       opus_val16
		beta       opus_val16
		budget     opus_int32
		tell       opus_int32
	)
	if intra != 0 {
		coef = 0
		beta = beta_intra
	} else {
		beta = beta_coef[LM]
		coef = pred_coef[LM]
	}
	budget = opus_int32(dec.Storage * 8)
	for i = start; i < end; i++ {
		c = 0
		for {
			{
				var (
					qi  int64
					q   opus_val32
					tmp opus_val32
				)
				tell = opus_int32(ec_tell((*ec_ctx)(unsafe.Pointer(dec))))
				if budget-tell >= 15 {
					var pi int64
					pi = (func() int64 {
						if i < 20 {
							return i
						}
						return 20
					}()) * 2
					qi = ec_laplace_decode(dec, uint64(int64(*(*uint8)(unsafe.Add(unsafe.Pointer(prob_model), pi)))<<7), int64(*(*uint8)(unsafe.Add(unsafe.Pointer(prob_model), pi+1)))<<6)
				} else if budget-tell >= 2 {
					qi = ec_dec_icdf(dec, &small_energy_icdf[0], 2)
					qi = (qi >> 1) ^ (-(qi & 1))
				} else if budget-tell >= 1 {
					qi = -ec_dec_bit_logp(dec, 1)
				} else {
					qi = -1
				}
				q = opus_val32(qi)
				if (-9.0) > float64(*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))) {
					*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) = opus_val16(-9.0)
				} else {
					*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands)))
				}
				tmp = (opus_val32(coef) * opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))))) + prev[c] + q
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) = opus_val16(tmp)
				prev[c] = prev[c] + q - opus_val32(beta)*q
			}
			if func() int64 {
				p := &c
				*p++
				return *p
			}() >= C {
				break
			}
		}
	}
}
func unquant_fine_energy(m *OpusCustomMode, start int64, end int64, oldEBands *opus_val16, fine_quant *int64, dec *ec_dec, C int64) {
	var (
		i int64
		c int64
	)
	for i = start; i < end; i++ {
		if *(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i))) <= 0 {
			continue
		}
		c = 0
		for {
			{
				var (
					q2     int64
					offset opus_val16
				)
				q2 = int64(ec_dec_bits(dec, uint64(*(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i))))))
				offset = opus_val16((float64(q2)+0.5)*float64(1<<(14-*(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i)))))*(1.0/16384) - 0.5)
				*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) += offset
			}
			if func() int64 {
				p := &c
				*p++
				return *p
			}() >= C {
				break
			}
		}
	}
}
func unquant_energy_finalise(m *OpusCustomMode, start int64, end int64, oldEBands *opus_val16, fine_quant *int64, fine_priority *int64, bits_left int64, dec *ec_dec, C int64) {
	var (
		i    int64
		prio int64
		c    int64
	)
	for prio = 0; prio < 2; prio++ {
		for i = start; i < end && bits_left >= C; i++ {
			if *(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i))) >= MAX_FINE_BITS || *(*int64)(unsafe.Add(unsafe.Pointer(fine_priority), unsafe.Sizeof(int64(0))*uintptr(i))) != prio {
				continue
			}
			c = 0
			for {
				{
					var (
						q2     int64
						offset opus_val16
					)
					q2 = int64(ec_dec_bits(dec, 1))
					offset = opus_val16((float64(q2) - 0.5) * float64(1<<(14-*(*int64)(unsafe.Add(unsafe.Pointer(fine_quant), unsafe.Sizeof(int64(0))*uintptr(i)))-1)) * (1.0 / 16384))
					*(*opus_val16)(unsafe.Add(unsafe.Pointer(oldEBands), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) += offset
					bits_left--
				}
				if func() int64 {
					p := &c
					*p++
					return *p
				}() >= C {
					break
				}
			}
		}
	}
}
func amp2Log2(m *OpusCustomMode, effEnd int64, end int64, bandE *celt_ener, bandLogE *opus_val16, C int64) {
	var (
		c int64
		i int64
	)
	c = 0
	for {
		for i = 0; i < effEnd; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(bandLogE), unsafe.Sizeof(opus_val16(0))*uintptr(i+c*m.NbEBands))) = opus_val16((float32(math.Log(float64(*(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i+c*m.NbEBands))))) * 1.4426950408889634)) - float32(eMeans[i]))
		}
		for i = effEnd; i < end; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(bandLogE), unsafe.Sizeof(opus_val16(0))*uintptr(c*m.NbEBands+i))) = opus_val16(-14.0)
		}
		if func() int64 {
			p := &c
			*p++
			return *p
		}() >= C {
			break
		}
	}
}
