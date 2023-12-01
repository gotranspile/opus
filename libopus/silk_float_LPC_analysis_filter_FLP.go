package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_LPC_analysis_filter16_FLP(r_LPC []float32, PredCoef []float32, s []float32, length int) {
	var (
		ix       int
		LPC_pred float32
		s_ptr    *float32
	)
	for ix = 16; ix < length; ix++ {
		s_ptr = &s[ix-1]
		LPC_pred = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*0))*PredCoef[0] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*1)))*PredCoef[1] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*2)))*PredCoef[2] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*3)))*PredCoef[3] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*4)))*PredCoef[4] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*5)))*PredCoef[5] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*6)))*PredCoef[6] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*7)))*PredCoef[7] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*8)))*PredCoef[8] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*9)))*PredCoef[9] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*10)))*PredCoef[10] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*11)))*PredCoef[11] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*12)))*PredCoef[12] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*13)))*PredCoef[13] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*14)))*PredCoef[14] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*15)))*PredCoef[15]
		r_LPC[ix] = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*1)) - LPC_pred
	}
}
func silk_LPC_analysis_filter12_FLP(r_LPC []float32, PredCoef []float32, s []float32, length int) {
	var (
		ix       int
		LPC_pred float32
		s_ptr    *float32
	)
	for ix = 12; ix < length; ix++ {
		s_ptr = &s[ix-1]
		LPC_pred = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*0))*PredCoef[0] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*1)))*PredCoef[1] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*2)))*PredCoef[2] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*3)))*PredCoef[3] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*4)))*PredCoef[4] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*5)))*PredCoef[5] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*6)))*PredCoef[6] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*7)))*PredCoef[7] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*8)))*PredCoef[8] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*9)))*PredCoef[9] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*10)))*PredCoef[10] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*11)))*PredCoef[11]
		r_LPC[ix] = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*1)) - LPC_pred
	}
}
func silk_LPC_analysis_filter10_FLP(r_LPC []float32, PredCoef []float32, s []float32, length int) {
	var (
		ix       int
		LPC_pred float32
		s_ptr    *float32
	)
	for ix = 10; ix < length; ix++ {
		s_ptr = &s[ix-1]
		LPC_pred = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*0))*PredCoef[0] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*1)))*PredCoef[1] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*2)))*PredCoef[2] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*3)))*PredCoef[3] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*4)))*PredCoef[4] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*5)))*PredCoef[5] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*6)))*PredCoef[6] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*7)))*PredCoef[7] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*8)))*PredCoef[8] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*9)))*PredCoef[9]
		r_LPC[ix] = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*1)) - LPC_pred
	}
}
func silk_LPC_analysis_filter8_FLP(r_LPC []float32, PredCoef []float32, s []float32, length int) {
	var (
		ix       int
		LPC_pred float32
		s_ptr    *float32
	)
	for ix = 8; ix < length; ix++ {
		s_ptr = &s[ix-1]
		LPC_pred = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*0))*PredCoef[0] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*1)))*PredCoef[1] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*2)))*PredCoef[2] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*3)))*PredCoef[3] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*4)))*PredCoef[4] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*5)))*PredCoef[5] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*6)))*PredCoef[6] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*7)))*PredCoef[7]
		r_LPC[ix] = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*1)) - LPC_pred
	}
}
func silk_LPC_analysis_filter6_FLP(r_LPC []float32, PredCoef []float32, s []float32, length int) {
	var (
		ix       int
		LPC_pred float32
		s_ptr    *float32
	)
	for ix = 6; ix < length; ix++ {
		s_ptr = &s[ix-1]
		LPC_pred = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*0))*PredCoef[0] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*1)))*PredCoef[1] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*2)))*PredCoef[2] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*3)))*PredCoef[3] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*4)))*PredCoef[4] + *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), -int(unsafe.Sizeof(float32(0))*5)))*PredCoef[5]
		r_LPC[ix] = *(*float32)(unsafe.Add(unsafe.Pointer(s_ptr), unsafe.Sizeof(float32(0))*1)) - LPC_pred
	}
}
func silk_LPC_analysis_filter_FLP(r_LPC []float32, PredCoef []float32, s []float32, length int, Order int) {
	switch Order {
	case 6:
		silk_LPC_analysis_filter6_FLP(r_LPC, PredCoef, s, length)
	case 8:
		silk_LPC_analysis_filter8_FLP(r_LPC, PredCoef, s, length)
	case 10:
		silk_LPC_analysis_filter10_FLP(r_LPC, PredCoef, s, length)
	case 12:
		silk_LPC_analysis_filter12_FLP(r_LPC, PredCoef, s, length)
	case 16:
		silk_LPC_analysis_filter16_FLP(r_LPC, PredCoef, s, length)
	default:
	}
	libc.MemSet(unsafe.Pointer(&r_LPC[0]), 0, Order*int(unsafe.Sizeof(float32(0))))
}
