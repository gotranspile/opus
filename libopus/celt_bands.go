package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const SPREAD_NONE = 0
const SPREAD_LIGHT = 1
const SPREAD_NORMAL = 2
const SPREAD_AGGRESSIVE = 3

func hysteresis_decision(val opus_val16, thresholds *opus_val16, hysteresis *opus_val16, N int64, prev int64) int64 {
	var i int64
	for i = 0; i < N; i++ {
		if val < *(*opus_val16)(unsafe.Add(unsafe.Pointer(thresholds), unsafe.Sizeof(opus_val16(0))*uintptr(i))) {
			break
		}
	}
	if i > prev && val < *(*opus_val16)(unsafe.Add(unsafe.Pointer(thresholds), unsafe.Sizeof(opus_val16(0))*uintptr(prev)))+*(*opus_val16)(unsafe.Add(unsafe.Pointer(hysteresis), unsafe.Sizeof(opus_val16(0))*uintptr(prev))) {
		i = prev
	}
	if i < prev && val > *(*opus_val16)(unsafe.Add(unsafe.Pointer(thresholds), unsafe.Sizeof(opus_val16(0))*uintptr(prev-1)))-*(*opus_val16)(unsafe.Add(unsafe.Pointer(hysteresis), unsafe.Sizeof(opus_val16(0))*uintptr(prev-1))) {
		i = prev
	}
	return i
}
func celt_lcg_rand(seed opus_uint32) opus_uint32 {
	return seed*1664525 + 1013904223
}
func bitexact_cos(x opus_int16) opus_int16 {
	var (
		tmp opus_int32
		x2  opus_int16
	)
	tmp = ((opus_int32(x) * opus_int32(x)) + 4096) >> 13
	x2 = opus_int16(tmp)
	x2 = opus_int16(opus_int32(math.MaxInt16-x2) + (((opus_int32(x2) * opus_int32(opus_int16((((opus_int32(x2)*opus_int32(opus_int16((((opus_int32(x2)*opus_int32(opus_int16(-626)))+16384)>>15)+8277)))+16384)>>15)+(-7651)))) + 16384) >> 15))
	return x2 + 1
}
func bitexact_log2tan(isin int64, icos int64) int64 {
	var (
		lc int64
		ls int64
	)
	lc = ec_ilog(opus_uint32(icos))
	ls = ec_ilog(opus_uint32(isin))
	icos <<= 15 - lc
	isin <<= 15 - ls
	return (ls-lc)*(1<<11) + int64(((opus_int32(opus_int16(isin))*opus_int32(opus_int16((((opus_int32(opus_int16(isin))*opus_int32(opus_int16(-2597)))+16384)>>15)+7932)))+16384)>>15) - int64(((opus_int32(opus_int16(icos))*opus_int32(opus_int16((((opus_int32(opus_int16(icos))*opus_int32(opus_int16(-2597)))+16384)>>15)+7932)))+16384)>>15)
}
func compute_band_energies(m *OpusCustomMode, X *celt_sig, bandE *celt_ener, end int64, C int64, LM int64, arch int64) {
	var (
		i      int64
		c      int64
		N      int64
		eBands *opus_int16 = m.EBands
	)
	N = m.ShortMdctSize << LM
	c = 0
	for {
		for i = 0; i < end; i++ {
			var sum opus_val32
			sum = opus_val32(float64(func() opus_val32 {
				_ = arch
				return celt_inner_prod_c((*opus_val16)(unsafe.Pointer((*celt_sig)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_sig(0))*uintptr(c*N+(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))<<LM)))))), (*opus_val16)(unsafe.Pointer((*celt_sig)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_sig(0))*uintptr(c*N+(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))<<LM)))))), int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))<<LM)
			}()) + 1e-27)
			*(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i+c*m.NbEBands))) = celt_ener(float32(math.Sqrt(float64(sum))))
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
func normalise_bands(m *OpusCustomMode, freq *celt_sig, X *celt_norm, bandE *celt_ener, end int64, C int64, M int64) {
	var (
		i      int64
		c      int64
		N      int64
		eBands *opus_int16 = m.EBands
	)
	N = M * m.ShortMdctSize
	c = 0
	for {
		for i = 0; i < end; i++ {
			var (
				j int64
				g opus_val16 = opus_val16(1.0 / (float64(*(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i+c*m.NbEBands)))) + 1e-27))
			)
			for j = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))); j < M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))); j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j+c*N))) = celt_norm(*(*celt_sig)(unsafe.Add(unsafe.Pointer(freq), unsafe.Sizeof(celt_sig(0))*uintptr(j+c*N))) * celt_sig(g))
			}
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
func denormalise_bands(m *OpusCustomMode, X *celt_norm, freq *celt_sig, bandLogE *opus_val16, start int64, end int64, M int64, downsample int64, silence int64) {
	var (
		i      int64
		N      int64
		bound  int64
		f      *celt_sig
		x      *celt_norm
		eBands *opus_int16 = m.EBands
	)
	N = M * m.ShortMdctSize
	bound = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(end))))
	if downsample != 1 {
		if bound < (N / downsample) {
			bound = bound
		} else {
			bound = N / downsample
		}
	}
	if silence != 0 {
		bound = 0
		start = func() int64 {
			end = 0
			return end
		}()
	}
	f = freq
	x = (*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start)))))))
	for i = 0; i < M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start)))); i++ {
		*func() *celt_sig {
			p := &f
			x := *p
			*p = (*celt_sig)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(celt_sig(0))*1))
			return x
		}() = 0
	}
	for i = start; i < end; i++ {
		var (
			j        int64
			band_end int64
			g        opus_val16
			lg       opus_val16
		)
		j = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
		band_end = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1))))
		lg = (*(*opus_val16)(unsafe.Add(unsafe.Pointer(bandLogE), unsafe.Sizeof(opus_val16(0))*uintptr(i)))) + opus_val16(opus_val32(eMeans[i]))
		g = opus_val16(float32(math.Exp((func() float64 {
			if 32.0 < float64(lg) {
				return 32.0
			}
			return float64(lg)
		}()) * 0.6931471805599453)))
		for {
			*func() *celt_sig {
				p := &f
				x := *p
				*p = (*celt_sig)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(celt_sig(0))*1))
				return x
			}() = celt_sig(opus_val32(*func() *celt_norm {
				p := &x
				x := *p
				*p = (*celt_norm)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(celt_norm(0))*1))
				return x
			}()) * opus_val32(g))
			if func() int64 {
				p := &j
				*p++
				return *p
			}() >= band_end {
				break
			}
		}
	}
	libc.MemSet(unsafe.Pointer((*celt_sig)(unsafe.Add(unsafe.Pointer(freq), unsafe.Sizeof(celt_sig(0))*uintptr(bound)))), 0, int((N-bound)*int64(unsafe.Sizeof(celt_sig(0)))))
}
func anti_collapse(m *OpusCustomMode, X_ *celt_norm, collapse_masks *uint8, LM int64, C int64, size int64, start int64, end int64, logE *opus_val16, prev1logE *opus_val16, prev2logE *opus_val16, pulses *int64, seed opus_uint32, arch int64) {
	var (
		c int64
		i int64
		j int64
		k int64
	)
	for i = start; i < end; i++ {
		var (
			N0     int64
			thresh opus_val16
			sqrt_1 opus_val16
			depth  int64
		)
		N0 = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1))) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
		depth = int64(celt_udiv(opus_uint32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i)))+1), opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) >> opus_uint32(LM))
		thresh = opus_val16(float64(float32(math.Exp((float64(depth)*(-0.125))*0.6931471805599453))) * 0.5)
		sqrt_1 = opus_val16(1.0 / float64(float32(math.Sqrt(float64(N0<<LM)))))
		c = 0
		for {
			{
				var (
					X           *celt_norm
					prev1       opus_val16
					prev2       opus_val16
					Ediff       opus_val32
					r           opus_val16
					renormalize int64 = 0
				)
				prev1 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(prev1logE), unsafe.Sizeof(opus_val16(0))*uintptr(c*m.NbEBands+i)))
				prev2 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(prev2logE), unsafe.Sizeof(opus_val16(0))*uintptr(c*m.NbEBands+i)))
				if C == 1 {
					if prev1 > (*(*opus_val16)(unsafe.Add(unsafe.Pointer(prev1logE), unsafe.Sizeof(opus_val16(0))*uintptr(m.NbEBands+i)))) {
						prev1 = prev1
					} else {
						prev1 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(prev1logE), unsafe.Sizeof(opus_val16(0))*uintptr(m.NbEBands+i)))
					}
					if prev2 > (*(*opus_val16)(unsafe.Add(unsafe.Pointer(prev2logE), unsafe.Sizeof(opus_val16(0))*uintptr(m.NbEBands+i)))) {
						prev2 = prev2
					} else {
						prev2 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(prev2logE), unsafe.Sizeof(opus_val16(0))*uintptr(m.NbEBands+i)))
					}
				}
				Ediff = opus_val32((*(*opus_val16)(unsafe.Add(unsafe.Pointer(logE), unsafe.Sizeof(opus_val16(0))*uintptr(c*m.NbEBands+i)))) - (func() opus_val16 {
					if prev1 < prev2 {
						return prev1
					}
					return prev2
				}()))
				if 0 > Ediff {
					Ediff = 0
				} else {
					Ediff = Ediff
				}
				r = opus_val16(float64(float32(math.Exp(float64(-Ediff)*0.6931471805599453))) * 2.0)
				if LM == 3 {
					r *= opus_val16(1.41421356)
				}
				if thresh < r {
					r = thresh
				} else {
					r = r
				}
				r = r * sqrt_1
				X = (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(X_), unsafe.Sizeof(celt_norm(0))*uintptr(c*size)))), unsafe.Sizeof(celt_norm(0))*uintptr(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.EBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))<<LM)))
				for k = 0; k < 1<<LM; k++ {
					if (int64(*(*uint8)(unsafe.Add(unsafe.Pointer(collapse_masks), i*C+c))) & (1 << k)) == 0 {
						for j = 0; j < N0; j++ {
							seed = celt_lcg_rand(seed)
							if seed&0x8000 != 0 {
								*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr((j<<LM)+k))) = celt_norm(r)
							} else {
								*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr((j<<LM)+k))) = celt_norm(-r)
							}
						}
						renormalize = 1
					}
				}
				if renormalize != 0 {
					renormalise_vector(X, N0<<LM, opus_val16(Q15ONE), arch)
				}
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
func compute_channel_weights(Ex celt_ener, Ey celt_ener, w [2]opus_val16) {
	var minE celt_ener
	if Ex < Ey {
		minE = Ex
	} else {
		minE = Ey
	}
	Ex = Ex + minE/3
	Ey = Ey + minE/3
	w[0] = opus_val16(Ex)
	w[1] = opus_val16(Ey)
}
func intensity_stereo(m *OpusCustomMode, X *celt_norm, Y *celt_norm, bandE *celt_ener, bandID int64, N int64) {
	var (
		i     int64 = bandID
		j     int64
		a1    opus_val16
		a2    opus_val16
		left  opus_val16
		right opus_val16
		norm  opus_val16
	)
	left = opus_val16(*(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i))))
	right = opus_val16(*(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i+m.NbEBands))))
	norm = opus_val16(EPSILON + float64(float32(math.Sqrt(EPSILON+float64(opus_val32(left)*opus_val32(left))+float64(opus_val32(right)*opus_val32(right))))))
	a1 = opus_val16((opus_val32(left)) / opus_val32(norm))
	a2 = opus_val16((opus_val32(right)) / opus_val32(norm))
	for j = 0; j < N; j++ {
		var (
			r celt_norm
			l celt_norm
		)
		l = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j)))
		r = *(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j)))
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm((opus_val32(a1) * opus_val32(l)) + opus_val32(a2)*opus_val32(r))
	}
}
func stereo_split(X *celt_norm, Y *celt_norm, N int64) {
	var j int64
	for j = 0; j < N; j++ {
		var (
			r opus_val32
			l opus_val32
		)
		l = opus_val32(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j)))) * opus_val32(0.70710678)
		r = opus_val32(*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j)))) * opus_val32(0.70710678)
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(l + r)
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(r - l)
	}
}
func stereo_merge(X *celt_norm, Y *celt_norm, mid opus_val16, N int64, arch int64) {
	var (
		j     int64
		xp    opus_val32 = 0
		side  opus_val32 = 0
		El    opus_val32
		Er    opus_val32
		mid2  opus_val16
		t     opus_val32
		lgain opus_val32
		rgain opus_val32
	)
	_ = arch
	dual_inner_prod_c((*opus_val16)(unsafe.Pointer(Y)), (*opus_val16)(unsafe.Pointer(X)), (*opus_val16)(unsafe.Pointer(Y)), N, &xp, &side)
	xp = opus_val32(mid * opus_val16(xp))
	mid2 = mid
	El = (opus_val32(mid2) * opus_val32(mid2)) + side - xp*2
	Er = (opus_val32(mid2) * opus_val32(mid2)) + side + xp*2
	if float64(Er) < 0.0006 || float64(El) < 0.0006 {
		libc.MemCpy(unsafe.Pointer(Y), unsafe.Pointer(X), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(Y))-uintptr(unsafe.Pointer(X))))*0))
		return
	}
	t = El
	lgain = opus_val32(1.0 / float64(float32(math.Sqrt(float64(t)))))
	t = Er
	rgain = opus_val32(1.0 / float64(float32(math.Sqrt(float64(t)))))
	for j = 0; j < N; j++ {
		var (
			r celt_norm
			l celt_norm
		)
		l = celt_norm(mid * opus_val16(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j)))))
		r = *(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j)))
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(lgain * opus_val32(l-r))
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(rgain * opus_val32(l+r))
	}
}
func spreading_decision(m *OpusCustomMode, X *celt_norm, average *int64, last_decision int64, hf_average *int64, tapset_decision *int64, update_hf int64, end int64, C int64, M int64, spread_weight *int64) int64 {
	var (
		i        int64
		c        int64
		N0       int64
		sum      int64       = 0
		nbBands  int64       = 0
		eBands   *opus_int16 = m.EBands
		decision int64
		hf_sum   int64 = 0
	)
	N0 = M * m.ShortMdctSize
	if M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(end)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(end-1)))) <= 8 {
		return 0
	}
	c = 0
	for {
		for i = 0; i < end; i++ {
			var (
				j      int64
				N      int64
				tmp    int64      = 0
				tcount [3]int64   = [3]int64{}
				x      *celt_norm = (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), unsafe.Sizeof(celt_norm(0))*uintptr(c*N0)))
			)
			N = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
			if N <= 8 {
				continue
			}
			for j = 0; j < N; j++ {
				var x2N opus_val32
				x2N = opus_val32((*(*celt_norm)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(celt_norm(0))*uintptr(j))))*(*(*celt_norm)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(celt_norm(0))*uintptr(j))))) * opus_val32(N)
				if float64(x2N) < 0.25 {
					tcount[0]++
				}
				if float64(x2N) < 0.0625 {
					tcount[1]++
				}
				if float64(x2N) < 0.015625 {
					tcount[2]++
				}
			}
			if i > m.NbEBands-4 {
				hf_sum += int64(celt_udiv(opus_uint32((tcount[1]+tcount[0])*32), opus_uint32(N)))
			}
			tmp = int64(libc.BoolToInt((tcount[2]*2 >= N) + (tcount[1]*2 >= N) + (tcount[0]*2 >= N)))
			sum += tmp * *(*int64)(unsafe.Add(unsafe.Pointer(spread_weight), unsafe.Sizeof(int64(0))*uintptr(i)))
			nbBands += *(*int64)(unsafe.Add(unsafe.Pointer(spread_weight), unsafe.Sizeof(int64(0))*uintptr(i)))
		}
		if func() int64 {
			p := &c
			*p++
			return *p
		}() >= C {
			break
		}
	}
	if update_hf != 0 {
		if hf_sum != 0 {
			hf_sum = int64(celt_udiv(opus_uint32(hf_sum), opus_uint32(C*(4-m.NbEBands+end))))
		}
		*hf_average = (*hf_average + hf_sum) >> 1
		hf_sum = *hf_average
		if *tapset_decision == 2 {
			hf_sum += 4
		} else if *tapset_decision == 0 {
			hf_sum -= 4
		}
		if hf_sum > 22 {
			*tapset_decision = 2
		} else if hf_sum > 18 {
			*tapset_decision = 1
		} else {
			*tapset_decision = 0
		}
	}
	sum = int64(celt_udiv(opus_uint32(opus_int32(sum)<<8), opus_uint32(nbBands)))
	sum = (sum + *average) >> 1
	*average = sum
	sum = (sum*3 + (((3 - last_decision) << 7) + 64) + 2) >> 2
	if sum < 80 {
		decision = 3
	} else if sum < 256 {
		decision = 2
	} else if sum < 384 {
		decision = 1
	} else {
		decision = 0
	}
	return decision
}

var ordery_table [30]int64 = [30]int64{1, 0, 3, 0, 2, 1, 7, 0, 4, 3, 6, 1, 5, 2, 15, 0, 8, 7, 12, 3, 11, 4, 14, 1, 9, 6, 13, 2, 10, 5}

func deinterleave_hadamard(X *celt_norm, N0 int64, stride int64, hadamard int64) {
	var (
		i   int64
		j   int64
		tmp *celt_norm
		N   int64
	)
	N = N0 * stride
	tmp = (*celt_norm)(libc.Malloc(int(N * int64(unsafe.Sizeof(celt_norm(0))))))
	if hadamard != 0 {
		var ordery *int64 = (*int64)(unsafe.Add(unsafe.Pointer(&ordery_table[stride]), -int(unsafe.Sizeof(int64(0))*2)))
		for i = 0; i < stride; i++ {
			for j = 0; j < N0; j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(tmp), unsafe.Sizeof(celt_norm(0))*uintptr(*(*int64)(unsafe.Add(unsafe.Pointer(ordery), unsafe.Sizeof(int64(0))*uintptr(i)))*N0+j))) = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j*stride+i)))
			}
		}
	} else {
		for i = 0; i < stride; i++ {
			for j = 0; j < N0; j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(tmp), unsafe.Sizeof(celt_norm(0))*uintptr(i*N0+j))) = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j*stride+i)))
			}
		}
	}
	libc.MemCpy(unsafe.Pointer(X), unsafe.Pointer(tmp), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(X))-uintptr(unsafe.Pointer(tmp))))*0))
}
func interleave_hadamard(X *celt_norm, N0 int64, stride int64, hadamard int64) {
	var (
		i   int64
		j   int64
		tmp *celt_norm
		N   int64
	)
	N = N0 * stride
	tmp = (*celt_norm)(libc.Malloc(int(N * int64(unsafe.Sizeof(celt_norm(0))))))
	if hadamard != 0 {
		var ordery *int64 = (*int64)(unsafe.Add(unsafe.Pointer(&ordery_table[stride]), -int(unsafe.Sizeof(int64(0))*2)))
		for i = 0; i < stride; i++ {
			for j = 0; j < N0; j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(tmp), unsafe.Sizeof(celt_norm(0))*uintptr(j*stride+i))) = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(*(*int64)(unsafe.Add(unsafe.Pointer(ordery), unsafe.Sizeof(int64(0))*uintptr(i)))*N0+j)))
			}
		}
	} else {
		for i = 0; i < stride; i++ {
			for j = 0; j < N0; j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(tmp), unsafe.Sizeof(celt_norm(0))*uintptr(j*stride+i))) = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(i*N0+j)))
			}
		}
	}
	libc.MemCpy(unsafe.Pointer(X), unsafe.Pointer(tmp), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(X))-uintptr(unsafe.Pointer(tmp))))*0))
}
func haar1(X *celt_norm, N0 int64, stride int64) {
	var (
		i int64
		j int64
	)
	N0 >>= 1
	for i = 0; i < stride; i++ {
		for j = 0; j < N0; j++ {
			var (
				tmp1 opus_val32
				tmp2 opus_val32
			)
			tmp1 = opus_val32(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(stride*2*j+i)))) * opus_val32(0.70710678)
			tmp2 = opus_val32(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(stride*(j*2+1)+i)))) * opus_val32(0.70710678)
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(stride*2*j+i))) = celt_norm(tmp1 + tmp2)
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(stride*(j*2+1)+i))) = celt_norm(tmp1 - tmp2)
		}
	}
}
func compute_qn(N int64, b int64, offset int64, pulse_cap int64, stereo int64) int64 {
	var (
		exp2_table8 [8]opus_int16 = [8]opus_int16{16384, 17866, 19483, 21247, 23170, 25267, 27554, 30048}
		qn          int64
		qb          int64
		N2          int64 = N*2 - 1
	)
	if stereo != 0 && N == 2 {
		N2--
	}
	qb = int64(celt_sudiv(opus_int32(b+N2*offset), opus_int32(N2)))
	if (b - pulse_cap - (4 << BITRES)) < qb {
		qb = b - pulse_cap - (4 << BITRES)
	} else {
		qb = qb
	}
	if (8 << BITRES) < qb {
		qb = 8 << BITRES
	} else {
		qb = qb
	}
	if qb < (1 << BITRES >> 1) {
		qn = 1
	} else {
		qn = int64(exp2_table8[qb&0x7]) >> (14 - (qb >> BITRES))
		qn = (qn + 1) >> 1 << 1
	}
	return qn
}

type band_ctx struct {
	Encode            int64
	Resynth           int64
	M                 *OpusCustomMode
	I                 int64
	Intensity         int64
	Spread            int64
	Tf_change         int64
	Ec                *ec_ctx
	Remaining_bits    opus_int32
	BandE             *celt_ener
	Seed              opus_uint32
	Arch              int64
	Theta_round       int64
	Disable_inv       int64
	Avoid_split_noise int64
}
type split_ctx struct {
	Inv    int64
	Imid   int64
	Iside  int64
	Delta  int64
	Itheta int64
	Qalloc int64
}

func compute_theta(ctx *band_ctx, sctx *split_ctx, X *celt_norm, Y *celt_norm, N int64, b *int64, B int64, B0 int64, LM int64, stereo int64, fill *int64) {
	var (
		qn        int64
		itheta    int64 = 0
		delta     int64
		imid      int64
		iside     int64
		qalloc    int64
		pulse_cap int64
		offset    int64
		tell      opus_int32
		inv       int64 = 0
		encode    int64
		m         *OpusCustomMode
		i         int64
		intensity int64
		ec        *ec_ctx
		bandE     *celt_ener
	)
	encode = ctx.Encode
	m = ctx.M
	i = ctx.I
	intensity = ctx.Intensity
	ec = ctx.Ec
	bandE = ctx.BandE
	pulse_cap = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(m.LogN), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) + LM*(1<<BITRES)
	offset = (pulse_cap >> 1) - (func() int64 {
		if stereo != 0 && N == 2 {
			return QTHETA_OFFSET_TWOPHASE
		}
		return QTHETA_OFFSET
	}())
	qn = compute_qn(N, *b, offset, pulse_cap, stereo)
	if stereo != 0 && i >= intensity {
		qn = 1
	}
	if encode != 0 {
		itheta = stereo_itheta(X, Y, stereo, N, ctx.Arch)
	}
	tell = opus_int32(ec_tell_frac(ec))
	if qn != 1 {
		if encode != 0 {
			if stereo == 0 || ctx.Theta_round == 0 {
				itheta = (itheta*int64(opus_int32(qn)) + 8192) >> 14
				if stereo == 0 && ctx.Avoid_split_noise != 0 && itheta > 0 && itheta < qn {
					var unquantized int64 = int64(celt_udiv(opus_uint32(opus_int32(itheta)*16384), opus_uint32(qn)))
					imid = int64(bitexact_cos(opus_int16(unquantized)))
					iside = int64(bitexact_cos(opus_int16(16384 - unquantized)))
					delta = int64(((opus_int32(opus_int16((N-1)<<7)) * opus_int32(opus_int16(bitexact_log2tan(iside, imid)))) + 16384) >> 15)
					if delta > *b {
						itheta = qn
					} else if delta < -*b {
						itheta = 0
					}
				}
			} else {
				var (
					down int64
					bias int64
				)
				if itheta > 8192 {
					bias = math.MaxInt16 / qn
				} else {
					bias = int64(-32767 / qn)
				}
				if (qn - 1) < (func() int64 {
					if 0 > ((itheta*int64(opus_int32(qn)) + bias) >> 14) {
						return 0
					}
					return (itheta*int64(opus_int32(qn)) + bias) >> 14
				}()) {
					down = qn - 1
				} else if 0 > ((itheta*int64(opus_int32(qn)) + bias) >> 14) {
					down = 0
				} else {
					down = (itheta*int64(opus_int32(qn)) + bias) >> 14
				}
				if ctx.Theta_round < 0 {
					itheta = down
				} else {
					itheta = down + 1
				}
			}
		}
		if stereo != 0 && N > 2 {
			var (
				p0 int64 = 3
				x  int64 = itheta
				x0 int64 = qn / 2
				ft int64 = p0*(x0+1) + x0
			)
			if encode != 0 {
				ec_encode((*ec_enc)(unsafe.Pointer(ec)), uint64(func() int64 {
					if x <= x0 {
						return p0 * x
					}
					return (x - 1 - x0) + (x0+1)*p0
				}()), uint64(func() int64 {
					if x <= x0 {
						return p0 * (x + 1)
					}
					return (x - x0) + (x0+1)*p0
				}()), uint64(ft))
			} else {
				var fs int64
				fs = int64(ec_decode((*ec_dec)(unsafe.Pointer(ec)), uint64(ft)))
				if fs < (x0+1)*p0 {
					x = fs / p0
				} else {
					x = x0 + 1 + (fs - (x0+1)*p0)
				}
				ec_dec_update((*ec_dec)(unsafe.Pointer(ec)), uint64(func() int64 {
					if x <= x0 {
						return p0 * x
					}
					return (x - 1 - x0) + (x0+1)*p0
				}()), uint64(func() int64 {
					if x <= x0 {
						return p0 * (x + 1)
					}
					return (x - x0) + (x0+1)*p0
				}()), uint64(ft))
				itheta = x
			}
		} else if B0 > 1 || stereo != 0 {
			if encode != 0 {
				ec_enc_uint((*ec_enc)(unsafe.Pointer(ec)), opus_uint32(itheta), opus_uint32(qn+1))
			} else {
				itheta = int64(ec_dec_uint((*ec_dec)(unsafe.Pointer(ec)), opus_uint32(qn+1)))
			}
		} else {
			var (
				fs int64 = 1
				ft int64
			)
			ft = ((qn >> 1) + 1) * ((qn >> 1) + 1)
			if encode != 0 {
				var fl int64
				if itheta <= (qn >> 1) {
					fs = itheta + 1
				} else {
					fs = qn + 1 - itheta
				}
				if itheta <= (qn >> 1) {
					fl = itheta * (itheta + 1) >> 1
				} else {
					fl = ft - ((qn + 1 - itheta) * (qn + 2 - itheta) >> 1)
				}
				ec_encode((*ec_enc)(unsafe.Pointer(ec)), uint64(fl), uint64(fl+fs), uint64(ft))
			} else {
				var (
					fl int64 = 0
					fm int64
				)
				fm = int64(ec_decode((*ec_dec)(unsafe.Pointer(ec)), uint64(ft)))
				if fm < ((qn >> 1) * ((qn >> 1) + 1) >> 1) {
					itheta = int64((isqrt32(opus_uint32(fm)*8+1) - 1) >> 1)
					fs = itheta + 1
					fl = itheta * (itheta + 1) >> 1
				} else {
					itheta = int64((uint64((qn+1)*2) - isqrt32(opus_uint32(ft-fm-1)*8+1)) >> 1)
					fs = qn + 1 - itheta
					fl = ft - ((qn + 1 - itheta) * (qn + 2 - itheta) >> 1)
				}
				ec_dec_update((*ec_dec)(unsafe.Pointer(ec)), uint64(fl), uint64(fl+fs), uint64(ft))
			}
		}
		itheta = int64(celt_udiv(opus_uint32(opus_int32(itheta)*16384), opus_uint32(qn)))
		if encode != 0 && stereo != 0 {
			if itheta == 0 {
				intensity_stereo(m, X, Y, bandE, i, N)
			} else {
				stereo_split(X, Y, N)
			}
		}
	} else if stereo != 0 {
		if encode != 0 {
			inv = int64(libc.BoolToInt(itheta > 8192 && ctx.Disable_inv == 0))
			if inv != 0 {
				var j int64
				for j = 0; j < N; j++ {
					*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = -*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j)))
				}
			}
			intensity_stereo(m, X, Y, bandE, i, N)
		}
		if *b > 2<<BITRES && ctx.Remaining_bits > opus_int32(2<<BITRES) {
			if encode != 0 {
				ec_enc_bit_logp((*ec_enc)(unsafe.Pointer(ec)), inv, 2)
			} else {
				inv = ec_dec_bit_logp((*ec_dec)(unsafe.Pointer(ec)), 2)
			}
		} else {
			inv = 0
		}
		if ctx.Disable_inv != 0 {
			inv = 0
		}
		itheta = 0
	}
	qalloc = int64(ec_tell_frac(ec) - opus_uint32(tell))
	*b -= qalloc
	if itheta == 0 {
		imid = math.MaxInt16
		iside = 0
		*fill &= (1 << B) - 1
		delta = -16384
	} else if itheta == 16384 {
		imid = 0
		iside = math.MaxInt16
		*fill &= ((1 << B) - 1) << B
		delta = 16384
	} else {
		imid = int64(bitexact_cos(opus_int16(itheta)))
		iside = int64(bitexact_cos(opus_int16(16384 - itheta)))
		delta = int64(((opus_int32(opus_int16((N-1)<<7)) * opus_int32(opus_int16(bitexact_log2tan(iside, imid)))) + 16384) >> 15)
	}
	sctx.Inv = inv
	sctx.Imid = imid
	sctx.Iside = iside
	sctx.Delta = delta
	sctx.Itheta = itheta
	sctx.Qalloc = qalloc
}
func quant_band_n1(ctx *band_ctx, X *celt_norm, Y *celt_norm, lowband_out *celt_norm) uint64 {
	var (
		c      int64
		stereo int64
		x      *celt_norm = X
		encode int64
		ec     *ec_ctx
	)
	encode = ctx.Encode
	ec = ctx.Ec
	stereo = int64(libc.BoolToInt(Y != nil))
	c = 0
	for {
		{
			var sign int64 = 0
			if ctx.Remaining_bits >= opus_int32(1<<BITRES) {
				if encode != 0 {
					sign = int64(libc.BoolToInt(*(*celt_norm)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(celt_norm(0))*0)) < 0))
					ec_enc_bits((*ec_enc)(unsafe.Pointer(ec)), opus_uint32(sign), 1)
				} else {
					sign = int64(ec_dec_bits((*ec_dec)(unsafe.Pointer(ec)), 1))
				}
				ctx.Remaining_bits -= opus_int32(1 << BITRES)
			}
			if ctx.Resynth != 0 {
				if sign != 0 {
					*(*celt_norm)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(celt_norm(0))*0)) = celt_norm(float32(-NORM_SCALING))
				} else {
					*(*celt_norm)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(celt_norm(0))*0)) = celt_norm(NORM_SCALING)
				}
			}
			x = Y
		}
		if func() int64 {
			p := &c
			*p++
			return *p
		}() >= stereo+1 {
			break
		}
	}
	if lowband_out != nil {
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(lowband_out), unsafe.Sizeof(celt_norm(0))*0)) = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*0))
	}
	return 1
}
func quant_partition(ctx *band_ctx, X *celt_norm, N int64, b int64, B int64, lowband *celt_norm, LM int64, gain opus_val16, fill int64) uint64 {
	var (
		cache     *uint8
		q         int64
		curr_bits int64
		imid      int64      = 0
		iside     int64      = 0
		B0        int64      = B
		mid       opus_val16 = 0
		side      opus_val16 = 0
		cm        uint64     = 0
		Y         *celt_norm = nil
		encode    int64
		m         *OpusCustomMode
		i         int64
		spread    int64
		ec        *ec_ctx
	)
	encode = ctx.Encode
	m = ctx.M
	i = ctx.I
	spread = ctx.Spread
	ec = ctx.Ec
	cache = (*uint8)(unsafe.Add(unsafe.Pointer(m.Cache.Bits), *(*opus_int16)(unsafe.Add(unsafe.Pointer(m.Cache.Index), unsafe.Sizeof(opus_int16(0))*uintptr((LM+1)*m.NbEBands+i)))))
	if LM != -1 && b > int64(*(*uint8)(unsafe.Add(unsafe.Pointer(cache), *cache)))+12 && N > 2 {
		var (
			mbits         int64
			sbits         int64
			delta         int64
			itheta        int64
			qalloc        int64
			sctx          split_ctx
			next_lowband2 *celt_norm = nil
			rebalance     opus_int32
		)
		N >>= 1
		Y = (*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(N)))
		LM -= 1
		if B == 1 {
			fill = (fill & 1) | fill<<1
		}
		B = (B + 1) >> 1
		compute_theta(ctx, &sctx, X, Y, N, &b, B, B0, LM, 0, &fill)
		imid = sctx.Imid
		iside = sctx.Iside
		delta = sctx.Delta
		itheta = sctx.Itheta
		qalloc = sctx.Qalloc
		mid = opus_val16(float64(imid) * (1.0 / 32768))
		side = opus_val16(float64(iside) * (1.0 / 32768))
		if B0 > 1 && (itheta&0x3FFF) != 0 {
			if itheta > 8192 {
				delta -= delta >> (4 - LM)
			} else if 0 < (delta + (N << BITRES >> (5 - LM))) {
				delta = 0
			} else {
				delta = delta + (N << BITRES >> (5 - LM))
			}
		}
		if 0 > (func() int64 {
			if b < ((b - delta) / 2) {
				return b
			}
			return (b - delta) / 2
		}()) {
			mbits = 0
		} else if b < ((b - delta) / 2) {
			mbits = b
		} else {
			mbits = (b - delta) / 2
		}
		sbits = b - mbits
		ctx.Remaining_bits -= opus_int32(qalloc)
		if lowband != nil {
			next_lowband2 = (*celt_norm)(unsafe.Add(unsafe.Pointer(lowband), unsafe.Sizeof(celt_norm(0))*uintptr(N)))
		}
		rebalance = ctx.Remaining_bits
		if mbits >= sbits {
			cm = quant_partition(ctx, X, N, mbits, B, lowband, LM, gain*mid, fill)
			rebalance = opus_int32(mbits - int64(rebalance-ctx.Remaining_bits))
			if rebalance > opus_int32(3<<BITRES) && itheta != 0 {
				sbits += int64(rebalance - opus_int32(3<<BITRES))
			}
			cm |= quant_partition(ctx, Y, N, sbits, B, next_lowband2, LM, gain*side, fill>>B) << uint64(B0>>1)
		} else {
			cm = quant_partition(ctx, Y, N, sbits, B, next_lowband2, LM, gain*side, fill>>B) << uint64(B0>>1)
			rebalance = opus_int32(sbits - int64(rebalance-ctx.Remaining_bits))
			if rebalance > opus_int32(3<<BITRES) && itheta != 16384 {
				mbits += int64(rebalance - opus_int32(3<<BITRES))
			}
			cm |= quant_partition(ctx, X, N, mbits, B, lowband, LM, gain*mid, fill)
		}
	} else {
		q = bits2pulses(m, i, LM, b)
		curr_bits = pulses2bits(m, i, LM, q)
		ctx.Remaining_bits -= opus_int32(curr_bits)
		for ctx.Remaining_bits < 0 && q > 0 {
			ctx.Remaining_bits += opus_int32(curr_bits)
			q--
			curr_bits = pulses2bits(m, i, LM, q)
			ctx.Remaining_bits -= opus_int32(curr_bits)
		}
		if q != 0 {
			var K int64 = get_pulses(q)
			if encode != 0 {
				cm = alg_quant(X, N, K, spread, B, (*ec_enc)(unsafe.Pointer(ec)), gain, ctx.Resynth, ctx.Arch)
			} else {
				cm = alg_unquant(X, N, K, spread, B, (*ec_dec)(unsafe.Pointer(ec)), gain)
			}
		} else {
			var j int64
			if ctx.Resynth != 0 {
				var cm_mask uint64
				cm_mask = uint64(1<<B) - 1
				fill &= int64(cm_mask)
				if fill == 0 {
					libc.MemSet(unsafe.Pointer(X), 0, int(N*int64(unsafe.Sizeof(celt_norm(0)))))
				} else {
					if lowband == nil {
						for j = 0; j < N; j++ {
							ctx.Seed = celt_lcg_rand(ctx.Seed)
							*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(opus_int32(ctx.Seed) >> 20)
						}
						cm = cm_mask
					} else {
						for j = 0; j < N; j++ {
							var tmp opus_val16
							ctx.Seed = celt_lcg_rand(ctx.Seed)
							tmp = opus_val16(1.0 / 256)
							if ctx.Seed&0x8000 != 0 {
								tmp = tmp
							} else {
								tmp = -tmp
							}
							*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = *(*celt_norm)(unsafe.Add(unsafe.Pointer(lowband), unsafe.Sizeof(celt_norm(0))*uintptr(j))) + celt_norm(tmp)
						}
						cm = uint64(fill)
					}
					renormalise_vector(X, N, gain, ctx.Arch)
				}
			}
		}
	}
	return cm
}
func quant_band(ctx *band_ctx, X *celt_norm, N int64, b int64, B int64, lowband *celt_norm, LM int64, lowband_out *celt_norm, gain opus_val16, lowband_scratch *celt_norm, fill int64) uint64 {
	var (
		N0          int64 = N
		N_B         int64 = N
		N_B0        int64
		B0          int64 = B
		time_divide int64 = 0
		recombine   int64 = 0
		longBlocks  int64
		cm          uint64 = 0
		k           int64
		encode      int64
		tf_change   int64
	)
	encode = ctx.Encode
	tf_change = ctx.Tf_change
	longBlocks = int64(libc.BoolToInt(B0 == 1))
	N_B = int64(celt_udiv(opus_uint32(N_B), opus_uint32(B)))
	if N == 1 {
		return quant_band_n1(ctx, X, nil, lowband_out)
	}
	if tf_change > 0 {
		recombine = tf_change
	}
	if lowband_scratch != nil && lowband != nil && (recombine != 0 || (N_B&1) == 0 && tf_change < 0 || B0 > 1) {
		libc.MemCpy(unsafe.Pointer(lowband_scratch), unsafe.Pointer(lowband), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(lowband_scratch))-uintptr(unsafe.Pointer(lowband))))*0))
		lowband = lowband_scratch
	}
	for k = 0; k < recombine; k++ {
		var bit_interleave_table [16]uint8 = [16]uint8{0, 1, 1, 1, 2, 3, 3, 3, 2, 3, 3, 3, 2, 3, 3, 3}
		if encode != 0 {
			haar1(X, N>>k, 1<<k)
		}
		if lowband != nil {
			haar1(lowband, N>>k, 1<<k)
		}
		fill = int64(bit_interleave_table[fill&0xF]) | int64(bit_interleave_table[fill>>4])<<2
	}
	B >>= recombine
	N_B <<= recombine
	for (N_B&1) == 0 && tf_change < 0 {
		if encode != 0 {
			haar1(X, N_B, B)
		}
		if lowband != nil {
			haar1(lowband, N_B, B)
		}
		fill |= fill << B
		B <<= 1
		N_B >>= 1
		time_divide++
		tf_change++
	}
	B0 = B
	N_B0 = N_B
	if B0 > 1 {
		if encode != 0 {
			deinterleave_hadamard(X, N_B>>recombine, B0<<recombine, longBlocks)
		}
		if lowband != nil {
			deinterleave_hadamard(lowband, N_B>>recombine, B0<<recombine, longBlocks)
		}
	}
	cm = quant_partition(ctx, X, N, b, B, lowband, LM, gain, fill)
	if ctx.Resynth != 0 {
		if B0 > 1 {
			interleave_hadamard(X, N_B>>recombine, B0<<recombine, longBlocks)
		}
		N_B = N_B0
		B = B0
		for k = 0; k < time_divide; k++ {
			B >>= 1
			N_B <<= 1
			cm |= cm >> uint64(B)
			haar1(X, N_B, B)
		}
		for k = 0; k < recombine; k++ {
			var bit_deinterleave_table [16]uint8 = [16]uint8{0x0, 0x3, 0xC, 0xF, 0x30, 0x33, 0x3C, 0x3F, 0xC0, 0xC3, 0xCC, 0xCF, 0xF0, 0xF3, 0xFC, math.MaxUint8}
			cm = uint64(bit_deinterleave_table[cm])
			haar1(X, N0>>k, 1<<k)
		}
		B <<= recombine
		if lowband_out != nil {
			var (
				j int64
				n opus_val16
			)
			n = opus_val16(float32(math.Sqrt(float64(N0))))
			for j = 0; j < N0; j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(lowband_out), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(n * opus_val16(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*uintptr(j)))))
			}
		}
		cm &= uint64((1 << B) - 1)
	}
	return cm
}
func quant_band_stereo(ctx *band_ctx, X *celt_norm, Y *celt_norm, N int64, b int64, B int64, lowband *celt_norm, LM int64, lowband_out *celt_norm, lowband_scratch *celt_norm, fill int64) uint64 {
	var (
		imid      int64      = 0
		iside     int64      = 0
		inv       int64      = 0
		mid       opus_val16 = 0
		side      opus_val16 = 0
		cm        uint64     = 0
		mbits     int64
		sbits     int64
		delta     int64
		itheta    int64
		qalloc    int64
		sctx      split_ctx
		orig_fill int64
		encode    int64
		ec        *ec_ctx
	)
	encode = ctx.Encode
	ec = ctx.Ec
	if N == 1 {
		return quant_band_n1(ctx, X, Y, lowband_out)
	}
	orig_fill = fill
	compute_theta(ctx, &sctx, X, Y, N, &b, B, B, LM, 1, &fill)
	inv = sctx.Inv
	imid = sctx.Imid
	iside = sctx.Iside
	delta = sctx.Delta
	itheta = sctx.Itheta
	qalloc = sctx.Qalloc
	mid = opus_val16(float64(imid) * (1.0 / 32768))
	side = opus_val16(float64(iside) * (1.0 / 32768))
	if N == 2 {
		var (
			c    int64
			sign int64 = 0
			x2   *celt_norm
			y2   *celt_norm
		)
		mbits = b
		sbits = 0
		if itheta != 0 && itheta != 16384 {
			sbits = 1 << BITRES
		}
		mbits -= sbits
		c = int64(libc.BoolToInt(itheta > 8192))
		ctx.Remaining_bits -= opus_int32(qalloc + sbits)
		if c != 0 {
			x2 = Y
		} else {
			x2 = X
		}
		if c != 0 {
			y2 = X
		} else {
			y2 = Y
		}
		if sbits != 0 {
			if encode != 0 {
				sign = int64(libc.BoolToInt(*(*celt_norm)(unsafe.Add(unsafe.Pointer(x2), unsafe.Sizeof(celt_norm(0))*0))**(*celt_norm)(unsafe.Add(unsafe.Pointer(y2), unsafe.Sizeof(celt_norm(0))*1))-*(*celt_norm)(unsafe.Add(unsafe.Pointer(x2), unsafe.Sizeof(celt_norm(0))*1))**(*celt_norm)(unsafe.Add(unsafe.Pointer(y2), unsafe.Sizeof(celt_norm(0))*0)) < 0))
				ec_enc_bits((*ec_enc)(unsafe.Pointer(ec)), opus_uint32(sign), 1)
			} else {
				sign = int64(ec_dec_bits((*ec_dec)(unsafe.Pointer(ec)), 1))
			}
		}
		sign = 1 - sign*2
		cm = quant_band(ctx, x2, N, mbits, B, lowband, LM, lowband_out, opus_val16(Q15ONE), lowband_scratch, orig_fill)
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(y2), unsafe.Sizeof(celt_norm(0))*0)) = celt_norm(-sign) * *(*celt_norm)(unsafe.Add(unsafe.Pointer(x2), unsafe.Sizeof(celt_norm(0))*1))
		*(*celt_norm)(unsafe.Add(unsafe.Pointer(y2), unsafe.Sizeof(celt_norm(0))*1)) = celt_norm(sign) * *(*celt_norm)(unsafe.Add(unsafe.Pointer(x2), unsafe.Sizeof(celt_norm(0))*0))
		if ctx.Resynth != 0 {
			var tmp celt_norm
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*0)) = celt_norm(mid * opus_val16(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*0))))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*1)) = celt_norm(mid * opus_val16(*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*1))))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*0)) = celt_norm(side * opus_val16(*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*0))))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*1)) = celt_norm(side * opus_val16(*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*1))))
			tmp = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*0))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*0)) = tmp - (*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*0)))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*0)) = tmp + (*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*0)))
			tmp = *(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*1))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(celt_norm(0))*1)) = tmp - (*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*1)))
			*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*1)) = tmp + (*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*1)))
		}
	} else {
		var rebalance opus_int32
		if 0 > (func() int64 {
			if b < ((b - delta) / 2) {
				return b
			}
			return (b - delta) / 2
		}()) {
			mbits = 0
		} else if b < ((b - delta) / 2) {
			mbits = b
		} else {
			mbits = (b - delta) / 2
		}
		sbits = b - mbits
		ctx.Remaining_bits -= opus_int32(qalloc)
		rebalance = ctx.Remaining_bits
		if mbits >= sbits {
			cm = quant_band(ctx, X, N, mbits, B, lowband, LM, lowband_out, opus_val16(Q15ONE), lowband_scratch, fill)
			rebalance = opus_int32(mbits - int64(rebalance-ctx.Remaining_bits))
			if rebalance > opus_int32(3<<BITRES) && itheta != 0 {
				sbits += int64(rebalance - opus_int32(3<<BITRES))
			}
			cm |= quant_band(ctx, Y, N, sbits, B, nil, LM, nil, side, nil, fill>>B)
		} else {
			cm = quant_band(ctx, Y, N, sbits, B, nil, LM, nil, side, nil, fill>>B)
			rebalance = opus_int32(sbits - int64(rebalance-ctx.Remaining_bits))
			if rebalance > opus_int32(3<<BITRES) && itheta != 16384 {
				mbits += int64(rebalance - opus_int32(3<<BITRES))
			}
			cm |= quant_band(ctx, X, N, mbits, B, lowband, LM, lowband_out, opus_val16(Q15ONE), lowband_scratch, fill)
		}
	}
	if ctx.Resynth != 0 {
		if N != 2 {
			stereo_merge(X, Y, mid, N, ctx.Arch)
		}
		if inv != 0 {
			var j int64
			for j = 0; j < N; j++ {
				*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = -*(*celt_norm)(unsafe.Add(unsafe.Pointer(Y), unsafe.Sizeof(celt_norm(0))*uintptr(j)))
			}
		}
	}
	return cm
}
func special_hybrid_folding(m *OpusCustomMode, norm *celt_norm, norm2 *celt_norm, start int64, M int64, dual_stereo int64) {
	var (
		n1     int64
		n2     int64
		eBands *opus_int16 = m.EBands
	)
	n1 = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start))))
	n2 = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start+2)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start+1))))
	libc.MemCpy(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(n1)))), unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(n1*2-n2)))), int((n2-n1)*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(n1)))))-uintptr(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(n1*2-n2)))))))*0))
	if dual_stereo != 0 {
		libc.MemCpy(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(n1)))), unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(n1*2-n2)))), int((n2-n1)*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(n1)))))-uintptr(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(n1*2-n2)))))))*0))
	}
}
func quant_all_bands(encode int64, m *OpusCustomMode, start int64, end int64, X_ *celt_norm, Y_ *celt_norm, collapse_masks *uint8, bandE *celt_ener, pulses *int64, shortBlocks int64, spread int64, dual_stereo int64, intensity int64, tf_res *int64, total_bits opus_int32, balance opus_int32, ec *ec_ctx, LM int64, codedBands int64, seed *opus_uint32, complexity int64, arch int64, disable_inv int64) {
	var (
		i                int64
		remaining_bits   opus_int32
		eBands           *opus_int16 = m.EBands
		norm             *celt_norm
		norm2            *celt_norm
		_norm            *celt_norm
		_lowband_scratch *celt_norm
		X_save           *celt_norm
		Y_save           *celt_norm
		X_save2          *celt_norm
		Y_save2          *celt_norm
		norm_save2       *celt_norm
		resynth_alloc    int64
		lowband_scratch  *celt_norm
		B                int64
		M                int64
		lowband_offset   int64
		update_lowband   int64 = 1
		C                int64
	)
	if Y_ != nil {
		C = 2
	} else {
		C = 1
	}
	var norm_offset int64
	var theta_rdo int64 = int64(libc.BoolToInt(encode != 0 && Y_ != nil && dual_stereo == 0 && complexity >= 8))
	var resynth int64 = int64(libc.BoolToInt(encode == 0 || theta_rdo != 0))
	var ctx band_ctx
	M = 1 << LM
	if shortBlocks != 0 {
		B = M
	} else {
		B = 1
	}
	norm_offset = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start))))
	_norm = (*celt_norm)(libc.Malloc(int((C * (M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(m.NbEBands-1)))) - norm_offset)) * int64(unsafe.Sizeof(celt_norm(0))))))
	norm = _norm
	norm2 = (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(m.NbEBands-1)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
	if encode != 0 && resynth != 0 {
		resynth_alloc = M * int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(m.NbEBands)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(m.NbEBands-1))))
	} else {
		resynth_alloc = ALLOC_NONE
	}
	_lowband_scratch = (*celt_norm)(libc.Malloc(int(resynth_alloc * int64(unsafe.Sizeof(celt_norm(0))))))
	if encode != 0 && resynth != 0 {
		lowband_scratch = _lowband_scratch
	} else {
		lowband_scratch = (*celt_norm)(unsafe.Add(unsafe.Pointer(X_), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(m.NbEBands-1)))))))
	}
	X_save = (*celt_norm)(libc.Malloc(int(resynth_alloc * int64(unsafe.Sizeof(celt_norm(0))))))
	Y_save = (*celt_norm)(libc.Malloc(int(resynth_alloc * int64(unsafe.Sizeof(celt_norm(0))))))
	X_save2 = (*celt_norm)(libc.Malloc(int(resynth_alloc * int64(unsafe.Sizeof(celt_norm(0))))))
	Y_save2 = (*celt_norm)(libc.Malloc(int(resynth_alloc * int64(unsafe.Sizeof(celt_norm(0))))))
	norm_save2 = (*celt_norm)(libc.Malloc(int(resynth_alloc * int64(unsafe.Sizeof(celt_norm(0))))))
	lowband_offset = 0
	ctx.BandE = bandE
	ctx.Ec = ec
	ctx.Encode = encode
	ctx.Intensity = intensity
	ctx.M = m
	ctx.Seed = *seed
	ctx.Spread = spread
	ctx.Arch = arch
	ctx.Disable_inv = disable_inv
	ctx.Resynth = resynth
	ctx.Theta_round = 0
	ctx.Avoid_split_noise = int64(libc.BoolToInt(B > 1))
	for i = start; i < end; i++ {
		var (
			tell              opus_int32
			b                 int64
			N                 int64
			curr_balance      opus_int32
			effective_lowband int64 = -1
			X                 *celt_norm
			Y                 *celt_norm
			tf_change         int64 = 0
			x_cm              uint64
			y_cm              uint64
			last              int64
		)
		ctx.I = i
		last = int64(libc.BoolToInt(i == end-1))
		X = (*celt_norm)(unsafe.Add(unsafe.Pointer(X_), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))
		if Y_ != nil {
			Y = (*celt_norm)(unsafe.Add(unsafe.Pointer(Y_), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))
		} else {
			Y = nil
		}
		N = M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))) - M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
		tell = opus_int32(ec_tell_frac(ec))
		if i != start {
			balance -= tell
		}
		remaining_bits = total_bits - tell - 1
		ctx.Remaining_bits = remaining_bits
		if i <= codedBands-1 {
			curr_balance = celt_sudiv(balance, opus_int32(func() int64 {
				if 3 < (codedBands - i) {
					return 3
				}
				return codedBands - i
			}()))
			if 0 > (func() opus_int32 {
				if 16383 < (func() opus_int32 {
					if (remaining_bits + 1) < opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i)))+int64(curr_balance)) {
						return remaining_bits + 1
					}
					return opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i))) + int64(curr_balance))
				}()) {
					return 16383
				}
				if (remaining_bits + 1) < opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i)))+int64(curr_balance)) {
					return remaining_bits + 1
				}
				return opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i))) + int64(curr_balance))
			}()) {
				b = 0
			} else if 16383 < (func() opus_int32 {
				if (remaining_bits + 1) < opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i)))+int64(curr_balance)) {
					return remaining_bits + 1
				}
				return opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i))) + int64(curr_balance))
			}()) {
				b = 16383
			} else if (remaining_bits + 1) < opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i)))+int64(curr_balance)) {
				b = int64(remaining_bits + 1)
			} else {
				b = *(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i))) + int64(curr_balance)
			}
		} else {
			b = 0
		}
		if resynth != 0 && (M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))-N >= M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(start)))) || i == start+1) && (update_lowband != 0 || lowband_offset == 0) {
			lowband_offset = i
		}
		if i == start+1 {
			special_hybrid_folding(m, norm, norm2, start, M, dual_stereo)
		}
		tf_change = *(*int64)(unsafe.Add(unsafe.Pointer(tf_res), unsafe.Sizeof(int64(0))*uintptr(i)))
		ctx.Tf_change = tf_change
		if i >= m.EffEBands {
			X = norm
			if Y_ != nil {
				Y = norm
			}
			lowband_scratch = nil
		}
		if last != 0 && theta_rdo == 0 {
			lowband_scratch = nil
		}
		if lowband_offset != 0 && (spread != 3 || B > 1 || tf_change < 0) {
			var (
				fold_start int64
				fold_end   int64
				fold_i     int64
			)
			if 0 > (M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(lowband_offset)))) - norm_offset - N) {
				effective_lowband = 0
			} else {
				effective_lowband = M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(lowband_offset)))) - norm_offset - N
			}
			fold_start = lowband_offset
			for M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(func() int64 {
				p := &fold_start
				*p--
				return *p
			}())))) > effective_lowband+norm_offset {
			}
			fold_end = lowband_offset - 1
			for func() int64 {
				p := &fold_end
				*p++
				return *p
			}() < i && M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(fold_end)))) < effective_lowband+norm_offset+N {
			}
			x_cm = func() uint64 {
				y_cm = 0
				return y_cm
			}()
			fold_i = fold_start
			for {
				x_cm |= uint64(*(*uint8)(unsafe.Add(unsafe.Pointer(collapse_masks), fold_i*C+0)))
				y_cm |= uint64(*(*uint8)(unsafe.Add(unsafe.Pointer(collapse_masks), fold_i*C+C-1)))
				if func() int64 {
					p := &fold_i
					*p++
					return *p
				}() >= fold_end {
					break
				}
			}
		} else {
			x_cm = func() uint64 {
				y_cm = uint64((1 << B) - 1)
				return y_cm
			}()
		}
		if dual_stereo != 0 && i == intensity {
			var j int64
			dual_stereo = 0
			if resynth != 0 {
				for j = 0; j < M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i))))-norm_offset; j++ {
					*(*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(j))) = celt_norm(float64(*(*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(j)))+*(*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(j)))) * 0.5)
				}
			}
		}
		if dual_stereo != 0 {
			x_cm = quant_band(&ctx, X, N, b/2, B, func() *celt_norm {
				if effective_lowband != -1 {
					return (*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(effective_lowband)))
				}
				return nil
			}(), LM, func() *celt_norm {
				if last != 0 {
					return nil
				}
				return (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
			}(), opus_val16(Q15ONE), lowband_scratch, int64(x_cm))
			y_cm = quant_band(&ctx, Y, N, b/2, B, func() *celt_norm {
				if effective_lowband != -1 {
					return (*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(effective_lowband)))
				}
				return nil
			}(), LM, func() *celt_norm {
				if last != 0 {
					return nil
				}
				return (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm2), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
			}(), opus_val16(Q15ONE), lowband_scratch, int64(y_cm))
		} else {
			if Y != nil {
				if theta_rdo != 0 && i < intensity {
					var (
						ec_save      ec_ctx
						ec_save2     ec_ctx
						ctx_save     band_ctx
						ctx_save2    band_ctx
						dist0        opus_val32
						dist1        opus_val32
						cm           uint64
						cm2          uint64
						nstart_bytes int64
						nend_bytes   int64
						save_bytes   int64
						bytes_buf    *uint8
						bytes_save   [1275]uint8
						w            [2]opus_val16
					)
					compute_channel_weights(*(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i))), *(*celt_ener)(unsafe.Add(unsafe.Pointer(bandE), unsafe.Sizeof(celt_ener(0))*uintptr(i+m.NbEBands))), w)
					cm = x_cm | y_cm
					ec_save = *ec
					ctx_save = ctx
					libc.MemCpy(unsafe.Pointer(X_save), unsafe.Pointer(X), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(X_save))-uintptr(unsafe.Pointer(X))))*0))
					libc.MemCpy(unsafe.Pointer(Y_save), unsafe.Pointer(Y), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(Y_save))-uintptr(unsafe.Pointer(Y))))*0))
					ctx.Theta_round = -1
					x_cm = quant_band_stereo(&ctx, X, Y, N, b, B, func() *celt_norm {
						if effective_lowband != -1 {
							return (*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(effective_lowband)))
						}
						return nil
					}(), LM, func() *celt_norm {
						if last != 0 {
							return nil
						}
						return (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
					}(), lowband_scratch, int64(cm))
					dist0 = opus_val32(((w[0]) * opus_val16(func() opus_val32 {
						_ = arch
						return celt_inner_prod_c((*opus_val16)(unsafe.Pointer(X_save)), (*opus_val16)(unsafe.Pointer(X)), N)
					}())) + (w[1])*opus_val16(func() opus_val32 {
						_ = arch
						return celt_inner_prod_c((*opus_val16)(unsafe.Pointer(Y_save)), (*opus_val16)(unsafe.Pointer(Y)), N)
					}()))
					cm2 = x_cm
					ec_save2 = *ec
					ctx_save2 = ctx
					libc.MemCpy(unsafe.Pointer(X_save2), unsafe.Pointer(X), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(X_save2))-uintptr(unsafe.Pointer(X))))*0))
					libc.MemCpy(unsafe.Pointer(Y_save2), unsafe.Pointer(Y), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(Y_save2))-uintptr(unsafe.Pointer(Y))))*0))
					if last == 0 {
						libc.MemCpy(unsafe.Pointer(norm_save2), unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(norm_save2))-uintptr(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))))))*0))
					}
					nstart_bytes = int64(ec_save.Offs)
					nend_bytes = int64(ec_save.Storage)
					bytes_buf = (*uint8)(unsafe.Add(unsafe.Pointer(ec_save.Buf), nstart_bytes))
					save_bytes = nend_bytes - nstart_bytes
					libc.MemCpy(unsafe.Pointer(&bytes_save[0]), unsafe.Pointer(bytes_buf), int(save_bytes*int64(unsafe.Sizeof(uint8(0)))+(int64(uintptr(unsafe.Pointer(&bytes_save[0]))-uintptr(unsafe.Pointer(bytes_buf))))*0))
					*ec = ec_save
					ctx = ctx_save
					libc.MemCpy(unsafe.Pointer(X), unsafe.Pointer(X_save), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(X))-uintptr(unsafe.Pointer(X_save))))*0))
					libc.MemCpy(unsafe.Pointer(Y), unsafe.Pointer(Y_save), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(Y))-uintptr(unsafe.Pointer(Y_save))))*0))
					if i == start+1 {
						special_hybrid_folding(m, norm, norm2, start, M, dual_stereo)
					}
					ctx.Theta_round = 1
					x_cm = quant_band_stereo(&ctx, X, Y, N, b, B, func() *celt_norm {
						if effective_lowband != -1 {
							return (*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(effective_lowband)))
						}
						return nil
					}(), LM, func() *celt_norm {
						if last != 0 {
							return nil
						}
						return (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
					}(), lowband_scratch, int64(cm))
					dist1 = opus_val32(((w[0]) * opus_val16(func() opus_val32 {
						_ = arch
						return celt_inner_prod_c((*opus_val16)(unsafe.Pointer(X_save)), (*opus_val16)(unsafe.Pointer(X)), N)
					}())) + (w[1])*opus_val16(func() opus_val32 {
						_ = arch
						return celt_inner_prod_c((*opus_val16)(unsafe.Pointer(Y_save)), (*opus_val16)(unsafe.Pointer(Y)), N)
					}()))
					if dist0 >= dist1 {
						x_cm = cm2
						*ec = ec_save2
						ctx = ctx_save2
						libc.MemCpy(unsafe.Pointer(X), unsafe.Pointer(X_save2), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(X))-uintptr(unsafe.Pointer(X_save2))))*0))
						libc.MemCpy(unsafe.Pointer(Y), unsafe.Pointer(Y_save2), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer(Y))-uintptr(unsafe.Pointer(Y_save2))))*0))
						if last == 0 {
							libc.MemCpy(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))), unsafe.Pointer(norm_save2), int(N*int64(unsafe.Sizeof(celt_norm(0)))+(int64(uintptr(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))))-uintptr(unsafe.Pointer(norm_save2))))*0))
						}
						libc.MemCpy(unsafe.Pointer(bytes_buf), unsafe.Pointer(&bytes_save[0]), int(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((bytes_buf-bytes_save)*0), save_bytes*int64(unsafe.Sizeof(uint8(0)))))))))
					}
				} else {
					ctx.Theta_round = 0
					x_cm = quant_band_stereo(&ctx, X, Y, N, b, B, func() *celt_norm {
						if effective_lowband != -1 {
							return (*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(effective_lowband)))
						}
						return nil
					}(), LM, func() *celt_norm {
						if last != 0 {
							return nil
						}
						return (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
					}(), lowband_scratch, int64(x_cm|y_cm))
				}
			} else {
				x_cm = quant_band(&ctx, X, N, b, B, func() *celt_norm {
					if effective_lowband != -1 {
						return (*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(effective_lowband)))
					}
					return nil
				}(), LM, func() *celt_norm {
					if last != 0 {
						return nil
					}
					return (*celt_norm)(unsafe.Add(unsafe.Pointer((*celt_norm)(unsafe.Add(unsafe.Pointer(norm), unsafe.Sizeof(celt_norm(0))*uintptr(M*int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(eBands), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))))), -int(unsafe.Sizeof(celt_norm(0))*uintptr(norm_offset))))
				}(), opus_val16(Q15ONE), lowband_scratch, int64(x_cm|y_cm))
			}
			y_cm = x_cm
		}
		*(*uint8)(unsafe.Add(unsafe.Pointer(collapse_masks), i*C+0)) = uint8(x_cm)
		*(*uint8)(unsafe.Add(unsafe.Pointer(collapse_masks), i*C+C-1)) = uint8(y_cm)
		balance += opus_int32(*(*int64)(unsafe.Add(unsafe.Pointer(pulses), unsafe.Sizeof(int64(0))*uintptr(i))) + int64(tell))
		update_lowband = int64(libc.BoolToInt(b > (N << BITRES)))
		ctx.Avoid_split_noise = 0
	}
	*seed = ctx.Seed
}
