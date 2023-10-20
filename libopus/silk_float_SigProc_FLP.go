package libopus

import (
	"math"
	"unsafe"
)

const PI = 0

func silk_sigmoid(x float32) float32 {
	return float32(1.0 / (math.Exp(float64(-x)) + 1.0))
}
func silk_float2int(x float32) int32 {
	return int32(int(math.Floor(float64(x + 0.5))))
}
func silk_float2short_array(out *int16, in *float32, length int32) {
	var k int32
	for k = int32(int(length) - 1); int(k) >= 0; k-- {
		if int(int32(int(math.Floor(float64(*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(k)))+0.5))))) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = silk_int16_MAX
		} else if int(int32(int(math.Floor(float64(*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(k)))+0.5))))) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = math.MinInt16
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(int32(int(math.Floor(float64(*(*float32)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(float32(0))*uintptr(k))) + 0.5)))))
		}
	}
}
func silk_short2float_array(out *float32, in *int16, length int32) {
	var k int32
	for k = int32(int(length) - 1); int(k) >= 0; k-- {
		*(*float32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(float32(0))*uintptr(k))) = float32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))
	}
}
func silk_log2(x float64) float32 {
	return float32(math.Log10(x) * 3.32192809488736)
}
