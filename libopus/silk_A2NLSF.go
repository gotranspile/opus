package libopus

import "unsafe"

const BIN_DIV_STEPS_A2NLSF_FIX = 3
const MAX_ITERATIONS_A2NLSF_FIX = 16

func silk_A2NLSF_trans_poly(p *opus_int32, dd int64) {
	var (
		k int64
		n int64
	)
	for k = 2; k <= dd; k++ {
		for n = dd; n > k; n-- {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*uintptr(n-2))) -= *(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*uintptr(n)))
		}
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*uintptr(k-2))) -= opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) << 1)
	}
}
func silk_A2NLSF_eval_poly(p *opus_int32, x opus_int32, dd int64) opus_int32 {
	var (
		n     int64
		x_Q16 opus_int32
		y32   opus_int32
	)
	y32 = *(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*uintptr(dd)))
	x_Q16 = opus_int32(opus_uint32(x) << 4)
	if dd == 8 {
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*7))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*6))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*5))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*4))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*3))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*2))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*1))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*0))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
	} else {
		for n = dd - 1; n >= 0; n-- {
			y32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(p), unsafe.Sizeof(opus_int32(0))*uintptr(n)))) + opus_int32((int64(y32)*int64(x_Q16))>>16)
		}
	}
	return y32
}
func silk_A2NLSF_init(a_Q16 *opus_int32, P *opus_int32, Q *opus_int32, dd int64) {
	var k int64
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(P), unsafe.Sizeof(opus_int32(0))*uintptr(dd))) = 1 << 16
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(Q), unsafe.Sizeof(opus_int32(0))*uintptr(dd))) = 1 << 16
	for k = 0; k < dd; k++ {
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(P), unsafe.Sizeof(opus_int32(0))*uintptr(k))) = -*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_Q16), unsafe.Sizeof(opus_int32(0))*uintptr(dd-k-1))) - *(*opus_int32)(unsafe.Add(unsafe.Pointer(a_Q16), unsafe.Sizeof(opus_int32(0))*uintptr(dd+k)))
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(Q), unsafe.Sizeof(opus_int32(0))*uintptr(k))) = -*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_Q16), unsafe.Sizeof(opus_int32(0))*uintptr(dd-k-1))) + *(*opus_int32)(unsafe.Add(unsafe.Pointer(a_Q16), unsafe.Sizeof(opus_int32(0))*uintptr(dd+k)))
	}
	for k = dd; k > 0; k-- {
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(P), unsafe.Sizeof(opus_int32(0))*uintptr(k-1))) -= *(*opus_int32)(unsafe.Add(unsafe.Pointer(P), unsafe.Sizeof(opus_int32(0))*uintptr(k)))
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(Q), unsafe.Sizeof(opus_int32(0))*uintptr(k-1))) += *(*opus_int32)(unsafe.Add(unsafe.Pointer(Q), unsafe.Sizeof(opus_int32(0))*uintptr(k)))
	}
	silk_A2NLSF_trans_poly(P, dd)
	silk_A2NLSF_trans_poly(Q, dd)
}
func silk_A2NLSF(NLSF *opus_int16, a_Q16 *opus_int32, d int64) {
	var (
		i       int64
		k       int64
		m       int64
		dd      int64
		root_ix int64
		ffrac   int64
		xlo     opus_int32
		xhi     opus_int32
		xmid    opus_int32
		ylo     opus_int32
		yhi     opus_int32
		ymid    opus_int32
		thr     opus_int32
		nom     opus_int32
		den     opus_int32
		P       [13]opus_int32
		Q       [13]opus_int32
		PQ      [2]*opus_int32
		p       *opus_int32
	)
	PQ[0] = &P[0]
	PQ[1] = &Q[0]
	dd = d >> 1
	silk_A2NLSF_init(a_Q16, &P[0], &Q[0], dd)
	p = &P[0]
	xlo = opus_int32(silk_LSFCosTab_FIX_Q12[0])
	ylo = silk_A2NLSF_eval_poly(p, xlo, dd)
	if ylo < 0 {
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*0)) = 0
		p = &Q[0]
		ylo = silk_A2NLSF_eval_poly(p, xlo, dd)
		root_ix = 1
	} else {
		root_ix = 0
	}
	k = 1
	i = 0
	thr = 0
	for {
		xhi = opus_int32(silk_LSFCosTab_FIX_Q12[k])
		yhi = silk_A2NLSF_eval_poly(p, xhi, dd)
		if ylo <= 0 && yhi >= thr || ylo >= 0 && yhi <= -thr {
			if yhi == 0 {
				thr = 1
			} else {
				thr = 0
			}
			ffrac = -256
			for m = 0; m < BIN_DIV_STEPS_A2NLSF_FIX; m++ {
				if 1 == 1 {
					xmid = ((xlo + xhi) >> 1) + ((xlo + xhi) & 1)
				} else {
					xmid = (((xlo + xhi) >> (1 - 1)) + 1) >> 1
				}
				ymid = silk_A2NLSF_eval_poly(p, xmid, dd)
				if ylo <= 0 && ymid >= 0 || ylo >= 0 && ymid <= 0 {
					xhi = xmid
					yhi = ymid
				} else {
					xlo = xmid
					ylo = ymid
					ffrac = ffrac + (128 >> m)
				}
			}
			if (func() opus_int32 {
				if ylo > 0 {
					return ylo
				}
				return -ylo
			}()) < 65536 {
				den = ylo - yhi
				nom = (opus_int32(opus_uint32(ylo) << opus_uint32(8-BIN_DIV_STEPS_A2NLSF_FIX))) + (den >> 1)
				if den != 0 {
					ffrac += int64(nom / den)
				}
			} else {
				ffrac += int64(ylo / ((ylo - yhi) >> opus_int32(8-BIN_DIV_STEPS_A2NLSF_FIX)))
			}
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*uintptr(root_ix))) = opus_int16(silk_min_32((opus_int32(opus_uint32(opus_int32(k))<<8))+opus_int32(ffrac), silk_int16_MAX))
			root_ix++
			if root_ix >= d {
				break
			}
			p = PQ[root_ix&1]
			xlo = opus_int32(silk_LSFCosTab_FIX_Q12[k-1])
			ylo = opus_int32(opus_uint32(1-(root_ix&2)) << 12)
		} else {
			k++
			xlo = xhi
			ylo = yhi
			thr = 0
			if k > LSF_COS_TAB_SZ_FIX {
				i++
				if i > MAX_ITERATIONS_A2NLSF_FIX {
					*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*0)) = opus_int16(opus_int32((1 << 15) / (d + 1)))
					for k = 1; k < d; k++ {
						*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = (*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*uintptr(k-1)))) + (*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*0)))
					}
					return
				}
				silk_bwexpander_32(a_Q16, d, 65536-(opus_int32(1<<i)))
				silk_A2NLSF_init(a_Q16, &P[0], &Q[0], dd)
				p = &P[0]
				xlo = opus_int32(silk_LSFCosTab_FIX_Q12[0])
				ylo = silk_A2NLSF_eval_poly(p, xlo, dd)
				if ylo < 0 {
					*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*0)) = 0
					p = &Q[0]
					ylo = silk_A2NLSF_eval_poly(p, xlo, dd)
					root_ix = 1
				} else {
					root_ix = 0
				}
				k = 1
			}
		}
	}
}
