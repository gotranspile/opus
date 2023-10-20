package libopus

import "unsafe"

func silk_NLSF_VQ_weights_laroia(pNLSFW_Q_OUT *int16, pNLSF_Q15 *int16, D int) {
	var (
		k        int
		tmp1_int int32
		tmp2_int int32
	)
	tmp1_int = int32(silk_max_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*0))), 1))
	tmp1_int = int32((1 << (int(NLSF_W_Q + 15))) / int(tmp1_int))
	tmp2_int = int32(silk_max_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*1)))-int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*0))), 1))
	tmp2_int = int32((1 << (int(NLSF_W_Q + 15))) / int(tmp2_int))
	*(*int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(int16(0))*0)) = int16(silk_min_int(int(tmp1_int)+int(tmp2_int), silk_int16_MAX))
	for k = 1; k < D-1; k += 2 {
		tmp1_int = int32(silk_max_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(k+1))))-int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(k)))), 1))
		tmp1_int = int32((1 << (int(NLSF_W_Q + 15))) / int(tmp1_int))
		*(*int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(silk_min_int(int(tmp1_int)+int(tmp2_int), silk_int16_MAX))
		tmp2_int = int32(silk_max_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(k+2))))-int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(k+1)))), 1))
		tmp2_int = int32((1 << (int(NLSF_W_Q + 15))) / int(tmp2_int))
		*(*int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(int16(0))*uintptr(k+1))) = int16(silk_min_int(int(tmp1_int)+int(tmp2_int), silk_int16_MAX))
	}
	tmp1_int = int32(silk_max_int((1<<15)-int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(D-1)))), 1))
	tmp1_int = int32((1 << (int(NLSF_W_Q + 15))) / int(tmp1_int))
	*(*int16)(unsafe.Add(unsafe.Pointer(pNLSFW_Q_OUT), unsafe.Sizeof(int16(0))*uintptr(D-1))) = int16(silk_min_int(int(tmp1_int)+int(tmp2_int), silk_int16_MAX))
}
