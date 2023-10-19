package libopus

import (
	"math"
	"unsafe"
)

func silk_NLSF_residual_dequant(x_Q10 [0]opus_int16, indices [0]int8, pred_coef_Q8 [0]uint8, quant_step_size_Q16 int64, order opus_int16) {
	var (
		i        int64
		out_Q10  int64
		pred_Q10 int64
	)
	out_Q10 = 0
	for i = int64(order - 1); i >= 0; i-- {
		pred_Q10 = int64((opus_int32(opus_int16(out_Q10)) * opus_int32(opus_int16(pred_coef_Q8[i]))) >> 8)
		out_Q10 = int64(opus_int32(opus_uint32(indices[i]) << 10))
		if out_Q10 > 0 {
			out_Q10 = out_Q10 - int64(opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10)+0.5))
		} else if out_Q10 < 0 {
			out_Q10 = out_Q10 + int64(opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10)+0.5))
		}
		out_Q10 = int64(opus_int32(pred_Q10 + int64(((opus_int32(out_Q10))*opus_int32(int64(opus_int16(quant_step_size_Q16))))>>16)))
		x_Q10[i] = opus_int16(out_Q10)
	}
}
func silk_NLSF_decode(pNLSF_Q15 *opus_int16, NLSFIndices *int8, psNLSF_CB *silk_NLSF_CB_struct) {
	var (
		i            int64
		pred_Q8      [16]uint8
		ec_ix        [16]opus_int16
		res_Q10      [16]opus_int16
		NLSF_Q15_tmp opus_int32
		pCB_element  *uint8
		pCB_Wght_Q9  *opus_int16
	)
	silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psNLSF_CB, int64(*NLSFIndices))
	silk_NLSF_residual_dequant(res_Q10[:], [0]int8((*int8)(unsafe.Add(unsafe.Pointer(NLSFIndices), 1))), pred_Q8[:], int64(psNLSF_CB.QuantStepSize_Q16), psNLSF_CB.Order)
	pCB_element = (*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_NLSF_Q8), opus_int16(*NLSFIndices)*psNLSF_CB.Order))
	pCB_Wght_Q9 = (*opus_int16)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_Wght_Q9), unsafe.Sizeof(opus_int16(0))*uintptr(opus_int16(*NLSFIndices)*psNLSF_CB.Order)))
	for i = 0; i < int64(psNLSF_CB.Order); i++ {
		NLSF_Q15_tmp = ((opus_int32(opus_uint32(opus_int32(res_Q10[i])) << 14)) / opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pCB_Wght_Q9), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) + (opus_int32(opus_uint32(opus_int16(*(*uint8)(unsafe.Add(unsafe.Pointer(pCB_element), i)))) << 7))
		if 0 > math.MaxInt16 {
			if NLSF_Q15_tmp > 0 {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = 0
			} else if NLSF_Q15_tmp < math.MaxInt16 {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = math.MaxInt16
			} else {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16(NLSF_Q15_tmp)
			}
		} else if NLSF_Q15_tmp > math.MaxInt16 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = math.MaxInt16
		} else if NLSF_Q15_tmp < 0 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = 0
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16(NLSF_Q15_tmp)
		}
	}
	silk_NLSF_stabilize(pNLSF_Q15, psNLSF_CB.DeltaMin_Q15, int64(psNLSF_CB.Order))
}
