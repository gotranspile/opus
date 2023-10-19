package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_NLSF_encode(NLSFIndices *int8, pNLSF_Q15 *opus_int16, psNLSF_CB *silk_NLSF_CB_struct, pW_Q2 *opus_int16, NLSF_mu_Q20 int64, nSurvivors int64, signalType int64) opus_int32 {
	var (
		i            int64
		s            int64
		ind1         int64
		bestIndex    int64
		prob_Q8      int64
		bits_q7      int64
		W_tmp_Q9     opus_int32
		ret          opus_int32
		err_Q24      *opus_int32
		RD_Q25       *opus_int32
		tempIndices1 *int64
		tempIndices2 *int8
		res_Q10      [16]opus_int16
		NLSF_tmp_Q15 [16]opus_int16
		W_adj_Q5     [16]opus_int16
		pred_Q8      [16]uint8
		ec_ix        [16]opus_int16
		pCB_element  *uint8
		iCDF_ptr     *uint8
		pCB_Wght_Q9  *opus_int16
	)
	silk_NLSF_stabilize(pNLSF_Q15, psNLSF_CB.DeltaMin_Q15, int64(psNLSF_CB.Order))
	err_Q24 = (*opus_int32)(libc.Malloc(int(uintptr(psNLSF_CB.NVectors) * unsafe.Sizeof(opus_int32(0)))))
	silk_NLSF_VQ([0]opus_int32(err_Q24), [0]opus_int16(pNLSF_Q15), [0]uint8(psNLSF_CB.CB1_NLSF_Q8), [0]opus_int16(psNLSF_CB.CB1_Wght_Q9), int64(psNLSF_CB.NVectors), int64(psNLSF_CB.Order))
	tempIndices1 = (*int64)(libc.Malloc(int(nSurvivors * int64(unsafe.Sizeof(int64(0))))))
	silk_insertion_sort_increasing(err_Q24, tempIndices1, int64(psNLSF_CB.NVectors), nSurvivors)
	RD_Q25 = (*opus_int32)(libc.Malloc(int(nSurvivors * int64(unsafe.Sizeof(opus_int32(0))))))
	tempIndices2 = (*int8)(libc.Malloc(int((nSurvivors * MAX_LPC_ORDER) * int64(unsafe.Sizeof(int8(0))))))
	for s = 0; s < nSurvivors; s++ {
		ind1 = *(*int64)(unsafe.Add(unsafe.Pointer(tempIndices1), unsafe.Sizeof(int64(0))*uintptr(s)))
		pCB_element = (*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_NLSF_Q8), ind1*int64(psNLSF_CB.Order)))
		pCB_Wght_Q9 = (*opus_int16)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_Wght_Q9), unsafe.Sizeof(opus_int16(0))*uintptr(ind1*int64(psNLSF_CB.Order))))
		for i = 0; i < int64(psNLSF_CB.Order); i++ {
			NLSF_tmp_Q15[i] = opus_int16(opus_uint16(opus_int16(*(*uint8)(unsafe.Add(unsafe.Pointer(pCB_element), i)))) << 7)
			W_tmp_Q9 = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pCB_Wght_Q9), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
			res_Q10[i] = opus_int16((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i)))-NLSF_tmp_Q15[i]) * opus_int32(opus_int16(W_tmp_Q9))) >> 14)
			W_adj_Q5[i] = opus_int16(silk_DIV32_varQ(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(pW_Q2), unsafe.Sizeof(opus_int16(0))*uintptr(i)))), opus_int32(opus_int16(W_tmp_Q9))*opus_int32(opus_int16(W_tmp_Q9)), 21))
		}
		silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psNLSF_CB, ind1)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(opus_int32(0))*uintptr(s))) = silk_NLSF_del_dec_quant([0]int8((*int8)(unsafe.Add(unsafe.Pointer(tempIndices2), s*MAX_LPC_ORDER))), res_Q10[:], W_adj_Q5[:], pred_Q8[:], ec_ix[:], [0]uint8(psNLSF_CB.Ec_Rates_Q5), int64(psNLSF_CB.QuantStepSize_Q16), psNLSF_CB.InvQuantStepSize_Q6, opus_int32(NLSF_mu_Q20), psNLSF_CB.Order)
		iCDF_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.CB1_iCDF), (signalType>>1)*int64(psNLSF_CB.NVectors)))
		if ind1 == 0 {
			prob_Q8 = 256 - int64(*(*uint8)(unsafe.Add(unsafe.Pointer(iCDF_ptr), ind1)))
		} else {
			prob_Q8 = int64(*(*uint8)(unsafe.Add(unsafe.Pointer(iCDF_ptr), ind1-1))) - int64(*(*uint8)(unsafe.Add(unsafe.Pointer(iCDF_ptr), ind1)))
		}
		bits_q7 = int64((8 << 7) - silk_lin2log(opus_int32(prob_Q8)))
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(opus_int32(0))*uintptr(s))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(opus_int32(0))*uintptr(s)))) + (opus_int32(opus_int16(bits_q7)))*opus_int32(opus_int16(NLSF_mu_Q20>>2))
	}
	silk_insertion_sort_increasing(RD_Q25, &bestIndex, nSurvivors, 1)
	*NLSFIndices = int8(*(*int64)(unsafe.Add(unsafe.Pointer(tempIndices1), unsafe.Sizeof(int64(0))*uintptr(bestIndex))))
	libc.MemCpy(unsafe.Pointer((*int8)(unsafe.Add(unsafe.Pointer(NLSFIndices), 1))), unsafe.Pointer((*int8)(unsafe.Add(unsafe.Pointer(tempIndices2), bestIndex*MAX_LPC_ORDER))), int(uintptr(psNLSF_CB.Order)*unsafe.Sizeof(int8(0))))
	silk_NLSF_decode(pNLSF_Q15, NLSFIndices, psNLSF_CB)
	ret = *(*opus_int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(opus_int32(0))*0))
	return ret
}
