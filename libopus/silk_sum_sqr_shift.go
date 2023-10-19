package libopus

import "unsafe"

func silk_sum_sqr_shift(energy *opus_int32, shift *int64, x *opus_int16, len_ int64) {
	var (
		i       int64
		shft    int64
		nrg_tmp opus_uint32
		nrg     opus_int32
	)
	shft = int64(31 - silk_CLZ32(opus_int32(len_)))
	nrg = opus_int32(len_)
	for i = 0; i < len_-1; i += 2 {
		nrg_tmp = opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))
		nrg_tmp = opus_uint32(opus_int32(nrg_tmp + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))))))
		nrg = nrg + opus_int32(nrg_tmp>>opus_uint32(shft))
	}
	if i < len_ {
		nrg_tmp = opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))
		nrg = nrg + opus_int32(nrg_tmp>>opus_uint32(shft))
	}
	shft = int64(silk_max_32(0, opus_int32(shft+3-int64(silk_CLZ32(nrg)))))
	nrg = 0
	for i = 0; i < len_-1; i += 2 {
		nrg_tmp = opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))
		nrg_tmp = opus_uint32(opus_int32(nrg_tmp + opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))))*opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))))))
		nrg = nrg + opus_int32(nrg_tmp>>opus_uint32(shft))
	}
	if i < len_ {
		nrg_tmp = opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))
		nrg = nrg + opus_int32(nrg_tmp>>opus_uint32(shft))
	}
	*shift = shft
	*energy = nrg
}
