package libopus

import "unsafe"

func silk_NLSF_VQ_weights_laroia(pNLSFW_Q_OUT *opus_int16, pNLSF_Q15 *opus_int16, D int64) {
	var (
		k        int64
		tmp1_int opus_int32
		tmp2_int opus_int32
	)
	tmp1_int = opus_int32(silk_max_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*0))), 1))
	tmp1_int = opus_int32(1<<(NLSF_W_Q+15)) / tmp1_int
	tmp2_int = opus_int32(silk_max_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*1))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*0))), 1))
	tmp2_int = opus_int32(1<<(NLSF_W_Q+15)) / tmp2_int
	*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(opus_int16(0))*0)) = opus_int16(silk_min_int(int64(tmp1_int+tmp2_int), silk_int16_MAX))
	for k = 1; k < D-1; k += 2 {
		tmp1_int = opus_int32(silk_max_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(k+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(k)))), 1))
		tmp1_int = opus_int32(1<<(NLSF_W_Q+15)) / tmp1_int
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(silk_min_int(int64(tmp1_int+tmp2_int), silk_int16_MAX))
		tmp2_int = opus_int32(silk_max_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(k+2)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(k+1)))), 1))
		tmp2_int = opus_int32(1<<(NLSF_W_Q+15)) / tmp2_int
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(opus_int16(0))*uintptr(k+1))) = opus_int16(silk_min_int(int64(tmp1_int+tmp2_int), silk_int16_MAX))
	}
	tmp1_int = opus_int32(silk_max_int(int64((1<<15)-*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(D-1)))), 1))
	tmp1_int = opus_int32(1<<(NLSF_W_Q+15)) / tmp1_int
	*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(opus_int16(0))*uintptr(D-1))) = opus_int16(silk_min_int(int64(tmp1_int+tmp2_int), silk_int16_MAX))
}
