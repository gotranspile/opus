package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const LPC_ORDER = 24

func _celt_lpc(_lpc *opus_val16, ac *opus_val32, p int64) {
	var (
		i     int64
		j     int64
		r     opus_val32
		error opus_val32 = *(*opus_val32)(unsafe.Add(unsafe.Pointer(ac), unsafe.Sizeof(opus_val32(0))*0))
		lpc   *float32   = (*float32)(unsafe.Pointer(_lpc))
	)
	libc.MemSet(unsafe.Pointer(lpc), 0, int(p*int64(unsafe.Sizeof(float32(0)))))
	if float64(*(*opus_val32)(unsafe.Add(unsafe.Pointer(ac), unsafe.Sizeof(opus_val32(0))*0))) > 1e-10 {
		for i = 0; i < p; i++ {
			var rr opus_val32 = 0
			for j = 0; j < i; j++ {
				rr += opus_val32((*(*float32)(unsafe.Add(unsafe.Pointer(lpc), unsafe.Sizeof(float32(0))*uintptr(j)))) * float32(*(*opus_val32)(unsafe.Add(unsafe.Pointer(ac), unsafe.Sizeof(opus_val32(0))*uintptr(i-j)))))
			}
			rr += *(*opus_val32)(unsafe.Add(unsafe.Pointer(ac), unsafe.Sizeof(opus_val32(0))*uintptr(i+1)))
			r = opus_val32(-(float32(rr) / float32(error)))
			*(*float32)(unsafe.Add(unsafe.Pointer(lpc), unsafe.Sizeof(float32(0))*uintptr(i))) = float32(r)
			for j = 0; j < (i+1)>>1; j++ {
				var (
					tmp1 opus_val32
					tmp2 opus_val32
				)
				tmp1 = opus_val32(*(*float32)(unsafe.Add(unsafe.Pointer(lpc), unsafe.Sizeof(float32(0))*uintptr(j))))
				tmp2 = opus_val32(*(*float32)(unsafe.Add(unsafe.Pointer(lpc), unsafe.Sizeof(float32(0))*uintptr(i-1-j))))
				*(*float32)(unsafe.Add(unsafe.Pointer(lpc), unsafe.Sizeof(float32(0))*uintptr(j))) = float32(tmp1 + r*tmp2)
				*(*float32)(unsafe.Add(unsafe.Pointer(lpc), unsafe.Sizeof(float32(0))*uintptr(i-1-j))) = float32(tmp2 + r*tmp1)
			}
			error = error - (r*r)*error
			if float64(error) <= float64(*(*opus_val32)(unsafe.Add(unsafe.Pointer(ac), unsafe.Sizeof(opus_val32(0))*0)))*0.001 {
				break
			}
		}
	}
}
func celt_fir_c(x *opus_val16, num *opus_val16, y *opus_val16, N int64, ord int64, arch int64) {
	var (
		i    int64
		j    int64
		rnum *opus_val16
	)
	rnum = (*opus_val16)(libc.Malloc(int(ord * int64(unsafe.Sizeof(opus_val16(0))))))
	for i = 0; i < ord; i++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(rnum), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(num), unsafe.Sizeof(opus_val16(0))*uintptr(ord-i-1)))
	}
	for i = 0; i < N-3; i += 4 {
		var sum [4]opus_val32
		sum[0] = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		sum[1] = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i+1))))
		sum[2] = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i+2))))
		sum[3] = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i+3))))
		_ = arch
		xcorr_kernel_c(rnum, (*opus_val16)(unsafe.Add(unsafe.Pointer((*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))), -int(unsafe.Sizeof(opus_val16(0))*uintptr(ord)))), sum, ord)
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(sum[0])
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+1))) = opus_val16(sum[1])
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+2))) = opus_val16(sum[2])
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+3))) = opus_val16(sum[3])
	}
	for ; i < N; i++ {
		var sum opus_val32 = opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		for j = 0; j < ord; j++ {
			sum = sum + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(rnum), unsafe.Sizeof(opus_val16(0))*uintptr(j))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i+j-ord))))
		}
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(sum)
	}
}
func celt_iir(_x *opus_val32, den *opus_val16, _y *opus_val32, N int64, ord int64, mem *opus_val16, arch int64) {
	var (
		i    int64
		j    int64
		rden *opus_val16
		y    *opus_val16
	)
	rden = (*opus_val16)(libc.Malloc(int(ord * int64(unsafe.Sizeof(opus_val16(0))))))
	y = (*opus_val16)(libc.Malloc(int((N + ord) * int64(unsafe.Sizeof(opus_val16(0))))))
	for i = 0; i < ord; i++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(rden), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*uintptr(ord-i-1)))
	}
	for i = 0; i < ord; i++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = -*(*opus_val16)(unsafe.Add(unsafe.Pointer(mem), unsafe.Sizeof(opus_val16(0))*uintptr(ord-i-1)))
	}
	for ; i < N+ord; i++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = 0
	}
	for i = 0; i < N-3; i += 4 {
		var sum [4]opus_val32
		sum[0] = *(*opus_val32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(opus_val32(0))*uintptr(i)))
		sum[1] = *(*opus_val32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(opus_val32(0))*uintptr(i+1)))
		sum[2] = *(*opus_val32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(opus_val32(0))*uintptr(i+2)))
		sum[3] = *(*opus_val32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(opus_val32(0))*uintptr(i+3)))
		_ = arch
		xcorr_kernel_c(rden, (*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i))), sum, ord)
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord))) = opus_val16(-(sum[0]))
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = sum[0]
		sum[1] = (sum[1]) + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*0)))
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord+1))) = opus_val16(-(sum[1]))
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val32(0))*uintptr(i+1))) = sum[1]
		sum[2] = (sum[2]) + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord+1))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*0)))
		sum[2] = (sum[2]) + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*1)))
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord+2))) = opus_val16(-(sum[2]))
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val32(0))*uintptr(i+2))) = sum[2]
		sum[3] = (sum[3]) + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord+2))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*0)))
		sum[3] = (sum[3]) + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord+1))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*1)))
		sum[3] = (sum[3]) + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(den), unsafe.Sizeof(opus_val16(0))*2)))
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord+3))) = opus_val16(-(sum[3]))
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val32(0))*uintptr(i+3))) = sum[3]
	}
	for ; i < N; i++ {
		var sum opus_val32 = *(*opus_val32)(unsafe.Add(unsafe.Pointer(_x), unsafe.Sizeof(opus_val32(0))*uintptr(i)))
		for j = 0; j < ord; j++ {
			sum -= opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(rden), unsafe.Sizeof(opus_val16(0))*uintptr(j)))) * opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+j))))
		}
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(y), unsafe.Sizeof(opus_val16(0))*uintptr(i+ord))) = opus_val16(sum)
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val32(0))*uintptr(i))) = sum
	}
	for i = 0; i < ord; i++ {
		*(*opus_val16)(unsafe.Add(unsafe.Pointer(mem), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = opus_val16(*(*opus_val32)(unsafe.Add(unsafe.Pointer(_y), unsafe.Sizeof(opus_val32(0))*uintptr(N-i-1))))
	}
}
func _celt_autocorr(x *opus_val16, ac *opus_val32, window *opus_val16, overlap int64, lag int64, n int64, arch int64) int64 {
	var (
		d     opus_val32
		i     int64
		k     int64
		fastN int64 = n - lag
		shift int64
		xptr  *opus_val16
		xx    *opus_val16
	)
	xx = (*opus_val16)(libc.Malloc(int(n * int64(unsafe.Sizeof(opus_val16(0))))))
	if overlap == 0 {
		xptr = x
	} else {
		for i = 0; i < n; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(xx), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = *(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))
		}
		for i = 0; i < overlap; i++ {
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(xx), unsafe.Sizeof(opus_val16(0))*uintptr(i))) = (*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))) * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
			*(*opus_val16)(unsafe.Add(unsafe.Pointer(xx), unsafe.Sizeof(opus_val16(0))*uintptr(n-i-1))) = (*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(n-i-1)))) * (*(*opus_val16)(unsafe.Add(unsafe.Pointer(window), unsafe.Sizeof(opus_val16(0))*uintptr(i))))
		}
		xptr = xx
	}
	shift = 0
	celt_pitch_xcorr_c(xptr, xptr, ac, fastN, lag+1, arch)
	for k = 0; k <= lag; k++ {
		for func() opus_val32 {
			i = k + fastN
			return func() opus_val32 {
				d = 0
				return d
			}()
		}(); i < n; i++ {
			d = d + opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(xptr), unsafe.Sizeof(opus_val16(0))*uintptr(i))))*opus_val32(*(*opus_val16)(unsafe.Add(unsafe.Pointer(xptr), unsafe.Sizeof(opus_val16(0))*uintptr(i-k))))
		}
		*(*opus_val32)(unsafe.Add(unsafe.Pointer(ac), unsafe.Sizeof(opus_val32(0))*uintptr(k))) += d
	}
	return shift
}
