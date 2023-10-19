package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const USE_CELT_FIR = 0

func silk_LPC_analysis_filter(out *opus_int16, in *opus_int16, B *opus_int16, len_ opus_int32, d opus_int32, arch int64) {
	var (
		j         int64
		ix        int64
		out32_Q12 opus_int32
		out32     opus_int32
		in_ptr    *opus_int16
	)
	_ = arch
	for ix = int64(d); ix < int64(len_); ix++ {
		in_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(ix-1)))
		out32_Q12 = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), unsafe.Sizeof(opus_int16(0))*0))) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*0)))
		out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), -int(unsafe.Sizeof(opus_int16(0))*1)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*1)))))
		out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), -int(unsafe.Sizeof(opus_int16(0))*2)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*2)))))
		out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), -int(unsafe.Sizeof(opus_int16(0))*3)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*3)))))
		out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), -int(unsafe.Sizeof(opus_int16(0))*4)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*4)))))
		out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), -int(unsafe.Sizeof(opus_int16(0))*5)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*5)))))
		for j = 6; j < int64(d); j += 2 {
			out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), -int(unsafe.Sizeof(opus_int16(0))*uintptr(j))))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*uintptr(j))))))
			out32_Q12 = opus_int32(opus_uint32(out32_Q12) + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(-j-1)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))))))
		}
		out32_Q12 = opus_int32(opus_uint32(opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in_ptr), unsafe.Sizeof(opus_int16(0))*1))))<<12)) - opus_uint32(out32_Q12))
		if 12 == 1 {
			out32 = (out32_Q12 >> 1) + (out32_Q12 & 1)
		} else {
			out32 = ((out32_Q12 >> (12 - 1)) + 1) >> 1
		}
		if out32 > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(ix))) = silk_int16_MAX
		} else if out32 < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(ix))) = math.MinInt16
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(ix))) = opus_int16(out32)
		}
	}
	libc.MemSet(unsafe.Pointer(out), 0, int(d*opus_int32(unsafe.Sizeof(opus_int16(0)))))
}
