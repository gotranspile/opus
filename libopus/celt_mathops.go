package libopus

import "unsafe"

const PI = 3.141592653

func fast_atan2f(y float32, x float32) float32 {
	var (
		x2 float32
		y2 float32
	)
	x2 = x * x
	y2 = y * y
	if float64(x2+y2) < 1e-18 {
		return 0
	}
	if x2 < y2 {
		var den float32 = float32((float64(y2) + cB*float64(x2)) * (float64(y2) + cC*float64(x2)))
		return float32(float64(-x*y)*(float64(y2)+cA*float64(x2))/float64(den) + float64(func() float32 {
			if y < 0 {
				return -(float32(PI) / 2)
			}
			return float32(PI) / 2
		}()))
	} else {
		var den float32 = float32((float64(x2) + cB*float64(y2)) * (float64(x2) + cC*float64(y2)))
		return float32(float64(x*y)*(float64(x2)+cA*float64(y2))/float64(den) + float64(func() float32 {
			if y < 0 {
				return -(float32(PI) / 2)
			}
			return float32(PI) / 2
		}()) - float64(func() float32 {
			if x*y < 0 {
				return -(float32(PI) / 2)
			}
			return float32(PI) / 2
		}()))
	}
}
func celt_maxabs16(x *opus_val16, len_ int64) opus_val32 {
	var (
		i      int64
		maxval opus_val16 = 0
		minval opus_val16 = 0
	)
	for i = 0; i < len_; i++ {
		if maxval > (*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))) {
			maxval = maxval
		} else {
			maxval = *(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))
		}
		if minval < (*(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))) {
			minval = minval
		} else {
			minval = *(*opus_val16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_val16(0))*uintptr(i)))
		}
	}
	if maxval > (-minval) {
		return opus_val32(maxval)
	}
	return opus_val32(-minval)
}
func isqrt32(_val opus_uint32) uint64 {
	var (
		b      uint64
		g      uint64
		bshift int64
	)
	g = 0
	bshift = (ec_ilog(_val) - 1) >> 1
	b = uint64(1 << bshift)
	for {
		{
			var t opus_uint32
			t = ((opus_uint32(g) << 1) + opus_uint32(b)) << opus_uint32(bshift)
			if t <= _val {
				g += b
				_val -= t
			}
			b >>= 1
			bshift--
		}
		if bshift < 0 {
			break
		}
	}
	return g
}
