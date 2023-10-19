package libopus

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func xcorr_kernel_c(x *opus_val16, y *opus_val16, sum [4]opus_val32, len_ int64) {
	var (
		j   int64
		y_0 opus_val16
		y_1 opus_val16
		y_2 opus_val16
		y_3 opus_val16
	)
	y_3 = 0
	y_0 = *func() *opus_val16 {
		p := &y
		x := *p
		*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
		return x
	}()
	y_1 = *func() *opus_val16 {
		p := &y
		x := *p
		*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
		return x
	}()
	y_2 = *func() *opus_val16 {
		p := &y
		x := *p
		*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
		return x
	}()
	for j = 0; j < len_-3; j += 4 {
		var tmp opus_val16
		tmp = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_3 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_0)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_1)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_2)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_3)
		tmp = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_0 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_1)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_2)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_3)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_0)
		tmp = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_1 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_2)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_3)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_0)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_1)
		tmp = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_2 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_3)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_0)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_1)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_2)
	}
	if func() int64 {
		p := &j
		x := *p
		*p++
		return x
	}() < len_ {
		var tmp opus_val16 = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_3 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_0)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_1)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_2)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_3)
	}
	if func() int64 {
		p := &j
		x := *p
		*p++
		return x
	}() < len_ {
		var tmp opus_val16 = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_0 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_1)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_2)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_3)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_0)
	}
	if j < len_ {
		var tmp opus_val16 = *func() *opus_val16 {
			p := &x
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		y_1 = *func() *opus_val16 {
			p := &y
			x := *p
			*p = (*opus_val16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_val16(0))*1))
			return x
		}()
		sum[0] = (sum[0]) + opus_val32(tmp)*opus_val32(y_2)
		sum[1] = (sum[1]) + opus_val32(tmp)*opus_val32(y_3)
		sum[2] = (sum[2]) + opus_val32(tmp)*opus_val32(y_0)
		sum[3] = (sum[3]) + opus_val32(tmp)*opus_val32(y_1)
	}
}
func dual_inner_prod_c(x *opus_val16, y01 *opus_val16, y02 *opus_val16, N int64, xy1 *opus_val32, xy2 *opus_val32) {
	var (
		i    int64
		xy01 opus_val32 = 0
		xy02 opus_val32 = 0
	)
	for i = 0; i < N; i++ {
		xy01 = xy01 + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y01), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		xy02 = xy02 + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y02), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
	}
	*xy1 = xy01
	*xy2 = xy02
}
func celt_inner_prod_c(x *opus_val16, y *opus_val16, N int64) opus_val32 {
	var (
		i  int64
		xy opus_val32 = 0
	)
	for i = 0; i < N; i++ {
		xy = xy + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
	}
	return xy
}
func find_best_pitch(xcorr *opus_val32, y *opus_val16, len_ int64, max_pitch int64, best_pitch *int64) {
	var (
		i        int64
		j        int64
		Syy      opus_val32 = 1
		best_num [2]opus_val16
		best_den [2]opus_val32
	)
	best_num[0] = opus_val16(-1)
	best_num[1] = opus_val16(-1)
	best_den[0] = 0
	best_den[1] = 0
	*(*int64)(unsafe.Add(unsafe.Pointer(best_pitch), unsafe.Sizeof(int64(0))*0)) = 0
	*(*int64)(unsafe.Add(unsafe.Pointer(best_pitch), unsafe.Sizeof(int64(0))*1)) = 1
	for j = 0; j < len_; j++ {
		Syy = Syy + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(j))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(j))))
	}
	for i = 0; i < max_pitch; i++ {
		if *(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i))) > 0 {
			var (
				num     opus_val16
				xcorr16 opus_val32
			)
			xcorr16 = *(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i)))
			xcorr16 *= opus_val32(1e-12)
			num = opus_val16(xcorr16 * xcorr16)
			if (num * opus_val16(best_den[1])) > ((best_num[1]) * opus_val16(Syy)) {
				if (num * opus_val16(best_den[0])) > ((best_num[0]) * opus_val16(Syy)) {
					best_num[1] = best_num[0]
					best_den[1] = best_den[0]
					*(*int64)(unsafe.Add(unsafe.Pointer(best_pitch), unsafe.Sizeof(int64(0))*1)) = *(*int64)(unsafe.Add(unsafe.Pointer(best_pitch), unsafe.Sizeof(int64(0))*0))
					best_num[0] = num
					best_den[0] = Syy
					*(*int64)(unsafe.Add(unsafe.Pointer(best_pitch), unsafe.Sizeof(int64(0))*0)) = i
				} else {
					best_num[1] = num
					best_den[1] = Syy
					*(*int64)(unsafe.Add(unsafe.Pointer(best_pitch), unsafe.Sizeof(int64(0))*1)) = i
				}
			}
		}
		Syy += (opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+len_)))) * opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+len_))))) - opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		if 1 > Syy {
			Syy = 1
		} else {
			Syy = Syy
		}
	}
}
func celt_fir5(x *opus_val16, num *opus_val16, N int64) {
	var (
		i    int64
		num0 opus_val16
		num1 opus_val16
		num2 opus_val16
		num3 opus_val16
		num4 opus_val16
		mem0 opus_val32
		mem1 opus_val32
		mem2 opus_val32
		mem3 opus_val32
		mem4 opus_val32
	)
	num0 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(num), unsafe.Sizeof(opus_val16(0))*0))
	num1 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(num), unsafe.Sizeof(opus_val16(0))*1))
	num2 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(num), unsafe.Sizeof(opus_val16(0))*2))
	num3 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(num), unsafe.Sizeof(opus_val16(0))*3))
	num4 = *(*opus_val16)(unsafe.Add(unsafe.Pointer(num), unsafe.Sizeof(opus_val16(0))*4))
	mem0 = 0
	mem1 = 0
	mem2 = 0
	mem3 = 0
	mem4 = 0
	for i = 0; i < N; i++ {
		var sum opus_val32 = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		sum = sum + opus_val32(num0)*mem0
		sum = sum + opus_val32(num1)*mem1
		sum = sum + opus_val32(num2)*mem2
		sum = sum + opus_val32(num3)*mem3
		sum = sum + opus_val32(num4)*mem4
		mem4 = mem3
		mem3 = mem2
		mem2 = mem1
		mem1 = mem0
		mem0 = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(sum)
	}
}
func pitch_downsample(x [0]*celt_sig, x_lp *opus_val16, len_ int64, C int64, arch int64) {
	var (
		i    int64
		ac   [5]opus_val32
		tmp  opus_val16 = opus_val16(Q15ONE)
		lpc  [4]opus_val16
		lpc2 [5]opus_val16
		c1   opus_val16 = opus_val16(0.8)
	)
	for i = 1; i < len_>>1; i++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(x_lp), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[0]), unsafe.Sizeof(celt_sig(0))*uintptr(i*2-1))))*0.25 + float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[0]), unsafe.Sizeof(celt_sig(0))*uintptr(i*2+1))))*0.25 + float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[0]), unsafe.Sizeof(celt_sig(0))*uintptr(i*2))))*0.5)
	}
	*(*opus_val16)(unsafe.Add(unsafe.Pointer(x_lp), unsafe.Sizeof(opus_val16(0))*0)) = opus_val16(float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[0]), unsafe.Sizeof(celt_sig(0))*1)))*0.25 + float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[0]), unsafe.Sizeof(celt_sig(0))*0)))*0.5)
	if C == 2 {
		for i = 1; i < len_>>1; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(x_lp), unsafe.Sizeof(opus_val16(0))*uintptr(i))) += opus_val16(float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[1]), unsafe.Sizeof(celt_sig(0))*uintptr(i*2-1))))*0.25 + float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[1]), unsafe.Sizeof(celt_sig(0))*uintptr(i*2+1))))*0.25 + float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[1]), unsafe.Sizeof(celt_sig(0))*uintptr(i*2))))*0.5)
		}
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(x_lp), unsafe.Sizeof(opus_val16(0))*0)) += opus_val16(float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[1]), unsafe.Sizeof(celt_sig(0))*1)))*0.25 + float64(*(*celt_sig)(unsafe.Add(unsafe.Pointer(x[1]), unsafe.Sizeof(celt_sig(0))*0)))*0.5)
	}
	_celt_autocorr(x_lp, &ac[0], nil, 0, 4, len_>>1, arch)
	ac[0] *= opus_val32(1.0001)
	for i = 1; i <= 4; i++ {
		ac[i] -= opus_val32(float64(ac[i]) * (float64(i) * 0.008) * (float64(i) * 0.008))
	}
	_celt_lpc(&lpc[0], &ac[0], 4)
	for i = 0; i < 4; i++ {
		tmp = opus_val16(float64(tmp) * 0.9)
		lpc[i] = (lpc[i]) * tmp
	}
	lpc2[0] = opus_val16(float64(lpc[0]) + 0.8)
	lpc2[1] = lpc[1] + c1*(lpc[0])
	lpc2[2] = lpc[2] + c1*(lpc[1])
	lpc2[3] = lpc[3] + c1*(lpc[2])
	lpc2[4] = c1 * (lpc[3])
	celt_fir5(x_lp, &lpc2[0], len_>>1)
}
func celt_pitch_xcorr_c(_x *opus_val16, _y *opus_val16, xcorr *opus_val32, len_ int64, max_pitch int64, arch int64) {
	var i int64
	for i = 0; i < max_pitch-3; i += 4 {
		var sum [4]opus_val32 = [4]opus_val32{}
		_ = arch
		xcorr_kernel_c(_x, (*opus_val16)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val16(0))*uintptr(i))), sum, len_)
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = sum[0]
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i+1))) = sum[1]
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i+2))) = sum[2]
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i+3))) = sum[3]
	}
	for ; i < max_pitch; i++ {
		var sum opus_val32
		sum = func() opus_val32 {
			_ = arch
			return celt_inner_prod_c(_x, (*opus_val16)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val16(0))*uintptr(i))), len_)
		}()
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = sum
	}
}
func pitch_search(x_lp *opus_val16, y *opus_val16, len_ int64, max_pitch int64, pitch *int64, arch int64) {
	var (
		i          int64
		j          int64
		lag        int64
		best_pitch [2]int64 = [2]int64{}
		x_lp4      *opus_val16
		y_lp4      *opus_val16
		xcorr      *opus_val32
		offset     int64
	)
	lag = len_ + max_pitch
	x_lp4 = (*opus_val16)(libc.Malloc(int((len_ >> 2) * int64(unsafe.Sizeof(opus_val16(0))))))
	y_lp4 = (*opus_val16)(libc.Malloc(int((lag >> 2) * int64(unsafe.Sizeof(opus_val16(0))))))
	xcorr = (*opus_val32)(libc.Malloc(int((max_pitch >> 1) * int64(unsafe.Sizeof(opus_val32(0))))))
	for j = 0; j < len_>>2; j++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(x_lp4), unsafe.Sizeof(opus_val16(0))*uintptr(j))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(x_lp), unsafe.Sizeof(opus_val16(0))*uintptr(j*2)))
	}
	for j = 0; j < lag>>2; j++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y_lp4), unsafe.Sizeof(opus_val16(0))*uintptr(j))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(j*2)))
	}
	celt_pitch_xcorr_c(x_lp4, y_lp4, xcorr, len_>>2, max_pitch>>2, arch)
	find_best_pitch(xcorr, y_lp4, len_>>2, max_pitch>>2, &best_pitch[0])
	for i = 0; i < max_pitch>>1; i++ {
		var sum opus_val32
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = 0
		if cmath.Abs(i-best_pitch[0]*2) > 2 && cmath.Abs(i-best_pitch[1]*2) > 2 {
			continue
		}
		sum = func() opus_val32 {
			_ = arch
			return celt_inner_prod_c(x_lp, (*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))), len_>>1)
		}()
		if opus_val32(-1) > sum {
			*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = opus_val32(-1)
		} else {
			*(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = sum
		}
	}
	find_best_pitch(xcorr, y, len_>>1, max_pitch>>1, &best_pitch[0])
	if best_pitch[0] > 0 && best_pitch[0] < (max_pitch>>1)-1 {
		var (
			a opus_val32
			b opus_val32
			c opus_val32
		)
		a = *(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(best_pitch[0]-1)))
		b = *(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(best_pitch[0])))
		c = *(*opus_val32)(unsafe.Add(unsafe.Pointer(xcorr), unsafe.Sizeof(opus_val32(0))*uintptr(best_pitch[0]+1)))
		if float64(c-a) > (float64(b-a) * 0.7) {
			offset = 1
		} else if float64(a-c) > (float64(b-c) * 0.7) {
			offset = -1
		} else {
			offset = 0
		}
	} else {
		offset = 0
	}
	*pitch = best_pitch[0]*2 - offset
}
func compute_pitch_gain(xy opus_val32, xx opus_val32, yy opus_val32) opus_val16 {
	return opus_val16(xy / opus_val32(float32(math.Sqrt(float64(xx*yy+1)))))
}

var second_check [16]int64 = [16]int64{0, 0, 3, 2, 3, 2, 5, 2, 3, 2, 3, 2, 5, 2, 3, 2}

func remove_doubling(x *opus_val16, maxperiod int64, minperiod int64, N int64, T0_ *int64, prev_period int64, prev_gain opus_val16, arch int64) opus_val16 {
	var (
		k          int64
		i          int64
		T          int64
		T0         int64
		g          opus_val16
		g0         opus_val16
		pg         opus_val16
		xy         opus_val32
		xx         opus_val32
		yy         opus_val32
		xy2        opus_val32
		xcorr      [3]opus_val32
		best_xy    opus_val32
		best_yy    opus_val32
		offset     int64
		minperiod0 int64
		yy_lookup  *opus_val32
	)
	minperiod0 = minperiod
	maxperiod /= 2
	minperiod /= 2
	*T0_ /= 2
	prev_period /= 2
	N /= 2
	x = (*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(maxperiod)))
	if *T0_ >= maxperiod {
		*T0_ = maxperiod - 1
	}
	T = func() int64 {
		T0 = *T0_
		return T0
	}()
	yy_lookup = (*opus_val32)(libc.Malloc(int((maxperiod + 1) * int64(unsafe.Sizeof(opus_val32(0))))))
	_ = arch
	dual_inner_prod_c(x, x, (*opus_val16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(opus_val16(0))*uintptr(T0)))), N, &xx, &xy)
	*(*opus_val32)(unsafe.Add(unsafe.Pointer(yy_lookup), unsafe.Sizeof(opus_val32(0))*0)) = xx
	yy = xx
	for i = 1; i <= maxperiod; i++ {
		yy = yy + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(opus_val16(0))*uintptr(i)))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(opus_val16(0))*uintptr(i))))) - opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(N-i))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(N-i))))
		if 0 > yy {
			*(*opus_val32)(unsafe.Add(unsafe.Pointer(yy_lookup), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = 0
		} else {
			*(*opus_val32)(unsafe.Add(unsafe.Pointer(yy_lookup), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = yy
		}
	}
	yy = *(*opus_val32)(unsafe.Add(unsafe.Pointer(yy_lookup), unsafe.Sizeof(opus_val32(0))*uintptr(T0)))
	best_xy = xy
	best_yy = yy
	g = func() opus_val16 {
		g0 = compute_pitch_gain(xy, xx, yy)
		return g0
	}()
	for k = 2; k <= 15; k++ {
		var (
			T1     int64
			T1b    int64
			g1     opus_val16
			cont   opus_val16 = 0
			thresh opus_val16
		)
		T1 = int64(celt_udiv(opus_uint32(T0*2+k), opus_uint32(k*2)))
		if T1 < minperiod {
			break
		}
		if k == 2 {
			if T1+T0 > maxperiod {
				T1b = T0
			} else {
				T1b = T0 + T1
			}
		} else {
			T1b = int64(celt_udiv(opus_uint32(second_check[k]*2*T0+k), opus_uint32(k*2)))
		}
		_ = arch
		dual_inner_prod_c(x, (*opus_val16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(opus_val16(0))*uintptr(T1)))), (*opus_val16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(opus_val16(0))*uintptr(T1b)))), N, &xy, &xy2)
		xy = opus_val32(float64(xy+xy2) * 0.5)
		yy = opus_val32(float64(*(*opus_val32)(unsafe.Add(unsafe.Pointer(yy_lookup), unsafe.Sizeof(opus_val32(0))*uintptr(T1)))+*(*opus_val32)(unsafe.Add(unsafe.Pointer(yy_lookup), unsafe.Sizeof(opus_val32(0))*uintptr(T1b)))) * 0.5)
		g1 = compute_pitch_gain(xy, xx, yy)
		if cmath.Abs(T1-prev_period) <= 1 {
			cont = prev_gain
		} else if cmath.Abs(T1-prev_period) <= 2 && k*5*k < T0 {
			cont = opus_val16(float64(prev_gain) * 0.5)
		} else {
			cont = 0
		}
		if 0.3 > ((float64(g0) * 0.7) - float64(cont)) {
			thresh = opus_val16(0.3)
		} else {
			thresh = opus_val16((float64(g0) * 0.7) - float64(cont))
		}
		if T1 < minperiod*3 {
			if 0.4 > ((float64(g0) * 0.85) - float64(cont)) {
				thresh = opus_val16(0.4)
			} else {
				thresh = opus_val16((float64(g0) * 0.85) - float64(cont))
			}
		} else if T1 < minperiod*2 {
			if 0.5 > ((float64(g0) * 0.9) - float64(cont)) {
				thresh = opus_val16(0.5)
			} else {
				thresh = opus_val16((float64(g0) * 0.9) - float64(cont))
			}
		}
		if g1 > thresh {
			best_xy = xy
			best_yy = yy
			T = T1
			g = g1
		}
	}
	if 0 > best_xy {
		best_xy = 0
	} else {
		best_xy = best_xy
	}
	if best_yy <= best_xy {
		pg = opus_val16(Q15ONE)
	} else {
		pg = opus_val16(float32(best_xy) / float32(best_yy+1))
	}
	for k = 0; k < 3; k++ {
		xcorr[k] = func() opus_val32 {
			_ = arch
			return celt_inner_prod_c(x, (*opus_val16)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(opus_val16(0))*uintptr(T+k-1)))), N)
		}()
	}
	if float64(xcorr[2]-xcorr[0]) > (float64(xcorr[1]-xcorr[0]) * 0.7) {
		offset = 1
	} else if float64(xcorr[0]-xcorr[2]) > (float64(xcorr[1]-xcorr[2]) * 0.7) {
		offset = -1
	} else {
		offset = 0
	}
	if pg > g {
		pg = g
	}
	*T0_ = T*2 + offset
	if *T0_ < minperiod0 {
		*T0_ = minperiod0
	}
	return pg
}
