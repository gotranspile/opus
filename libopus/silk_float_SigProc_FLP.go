package libopus

import (
	"math"
	"unsafe"
)

const PI = 0

func silk_sigmoid(x float32) float32 {
	return float32(1.0 / (math.Exp(float64(-x)) + 1.0))
}
func silk_float2int(x float32) opus_int32 {
	return opus_int32(int64(math.Floor(float64(x) + 0.5)))
}
func silk_float2short_array(out *opus_int16, in *float32, length opus_int32) {
	var k opus_int32
	for k = length - 1; k >= 0; k-- {
		if (opus_int32(int64(math.Floor(float64(*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(k)))) + 0.5)))) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (opus_int32(int64(math.Floor(float64(*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(k)))) + 0.5)))) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = math.MinInt16
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(opus_int32(int64(math.Floor(float64(*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(k)))) + 0.5))))
		}
	}
}
func silk_short2float_array(out *float32, in *opus_int16, length opus_int32) {
	var k opus_int32
	for k = length - 1; k >= 0; k-- {
		*(*float32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(float32(0))*uintptr(k))) = float32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k))))
	}
}
func silk_log2(x float64) float32 {
	return float32(math.Log10(x) * 3.32192809488736)
}
