package libopus

import (
	"math"
	"unsafe"
)

func silk_NLSF_residual_dequant(x_Q10 []int16, indices []int8, pred_coef_Q8 []uint8, quant_step_size_Q16 int, order int16) {
	var (
		i        int
		out_Q10  int
		pred_Q10 int
	)
	out_Q10 = 0
	for i = int(order) - 1; i >= 0; i-- {
		pred_Q10 = (int(int32(int16(out_Q10))) * int(int32(int16(pred_coef_Q8[i])))) >> 8
		out_Q10 = int(int32(int(uint32(indices[i])) << 10))
		if out_Q10 > 0 {
			out_Q10 = out_Q10 - int(int32(math.Floor(NLSF_QUANT_LEVEL_ADJ*(1<<10)+0.5)))
		} else if out_Q10 < 0 {
			out_Q10 = out_Q10 + int(int32(math.Floor(NLSF_QUANT_LEVEL_ADJ*(1<<10)+0.5)))
		}
		out_Q10 = int(int32(pred_Q10 + int((int64(int32(out_Q10))*int64(int16(quant_step_size_Q16)))>>16)))
		x_Q10[i] = int16(out_Q10)
	}
}
func silk_NLSF_decode(pNLSF_Q15 *int16, NLSFIndices *int8, psNLSF_CB *silk_NLSF_CB_struct) {
	var (
		i            int
		pred_Q8      [16]uint8
		ec_ix        [16]int16
		res_Q10      [16]int16
		NLSF_Q15_tmp int32
		pCB_element  *uint8
		pCB_Wght_Q9  *int16
	)
	silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psNLSF_CB, int(*NLSFIndices))
	silk_NLSF_residual_dequant(res_Q10[:], []int8((*int8)(unsafe.Add(unsafe.Pointer(NLSFIndices), 1))), pred_Q8[:], int(psNLSF_CB.QuantStepSize_Q16), psNLSF_CB.Order)
	pCB_element = (*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_NLSF_Q8), int(*NLSFIndices)*int(psNLSF_CB.Order)))
	pCB_Wght_Q9 = (*int16)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_Wght_Q9), unsafe.Sizeof(int16(0))*uintptr(int(*NLSFIndices)*int(psNLSF_CB.Order))))
	for i = 0; i < int(psNLSF_CB.Order); i++ {
		NLSF_Q15_tmp = int32(int(int32(int(int32(int(uint32(int32(res_Q10[i])))<<14))/int(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_Wght_Q9), unsafe.Sizeof(int16(0))*uintptr(i)))))) + int(int32(int(uint32(int16(*(*uint8)(unsafe.Add(unsafe.Pointer(pCB_element), i)))))<<7)))
		if 0 > math.MaxInt16 {
			if int(NLSF_Q15_tmp) > 0 {
				*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = 0
			} else if int(NLSF_Q15_tmp) < math.MaxInt16 {
				*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = math.MaxInt16
			} else {
				*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(NLSF_Q15_tmp)
			}
		} else if int(NLSF_Q15_tmp) > math.MaxInt16 {
			*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = math.MaxInt16
		} else if int(NLSF_Q15_tmp) < 0 {
			*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = 0
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(NLSF_Q15_tmp)
		}
	}
	silk_NLSF_stabilize(pNLSF_Q15, psNLSF_CB.DeltaMin_Q15, int(psNLSF_CB.Order))
}
