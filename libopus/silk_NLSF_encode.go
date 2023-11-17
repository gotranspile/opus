package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_NLSF_encode(NLSFIndices *int8, pNLSF_Q15 *int16, psNLSF_CB *silk_NLSF_CB_struct, pW_Q2 *int16, NLSF_mu_Q20 int, nSurvivors int, signalType int) int32 {
	var (
		i            int
		s            int
		ind1         int
		bestIndex    int
		prob_Q8      int
		bits_q7      int
		W_tmp_Q9     int32
		ret          int32
		err_Q24      *int32
		RD_Q25       *int32
		tempIndices1 *int
		tempIndices2 *int8
		res_Q10      [16]int16
		NLSF_tmp_Q15 [16]int16
		W_adj_Q5     [16]int16
		pred_Q8      [16]uint8
		ec_ix        [16]int16
		pCB_element  *uint8
		iCDF_ptr     *uint8
		pCB_Wght_Q9  *int16
	)
	silk_NLSF_stabilize([]int16(pNLSF_Q15), psNLSF_CB.DeltaMin_Q15, int(psNLSF_CB.Order))
	err_Q24 = (*int32)(libc.Malloc(int(uintptr(psNLSF_CB.NVectors) * unsafe.Sizeof(int32(0)))))
	silk_NLSF_VQ([]int32(err_Q24), []int16(pNLSF_Q15), []uint8(psNLSF_CB.CB1_NLSF_Q8), psNLSF_CB.CB1_Wght_Q9, int(psNLSF_CB.NVectors), int(psNLSF_CB.Order))
	tempIndices1 = (*int)(libc.Malloc(nSurvivors * int(unsafe.Sizeof(int(0)))))
	silk_insertion_sort_increasing([]int32(err_Q24), []int(tempIndices1), int(psNLSF_CB.NVectors), nSurvivors)
	RD_Q25 = (*int32)(libc.Malloc(nSurvivors * int(unsafe.Sizeof(int32(0)))))
	tempIndices2 = (*int8)(libc.Malloc((nSurvivors * MAX_LPC_ORDER) * int(unsafe.Sizeof(int8(0)))))
	for s = 0; s < nSurvivors; s++ {
		ind1 = *(*int)(unsafe.Add(unsafe.Pointer(tempIndices1), unsafe.Sizeof(int(0))*uintptr(s)))
		pCB_element = (*uint8)(unsafe.Pointer(&psNLSF_CB.CB1_NLSF_Q8[ind1*int(psNLSF_CB.Order)]))
		pCB_Wght_Q9 = &psNLSF_CB.CB1_Wght_Q9[ind1*int(psNLSF_CB.Order)]
		for i = 0; i < int(psNLSF_CB.Order); i++ {
			NLSF_tmp_Q15[i] = int16(int(uint16(int16(*(*uint8)(unsafe.Add(unsafe.Pointer(pCB_element), i))))) << 7)
			W_tmp_Q9 = int32(*(*int16)(unsafe.Add(unsafe.Pointer(pCB_Wght_Q9), unsafe.Sizeof(int16(0))*uintptr(i))))
			res_Q10[i] = int16((int(int32(int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(pNLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))))-int(NLSF_tmp_Q15[i])))) * int(int32(int16(W_tmp_Q9)))) >> 14)
			W_adj_Q5[i] = int16(silk_DIV32_varQ(int32(*(*int16)(unsafe.Add(unsafe.Pointer(pW_Q2), unsafe.Sizeof(int16(0))*uintptr(i)))), int32(int(int32(int16(W_tmp_Q9)))*int(int32(int16(W_tmp_Q9)))), 21))
		}
		silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psNLSF_CB, ind1)
		*(*int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(int32(0))*uintptr(s))) = silk_NLSF_del_dec_quant([]int8((*int8)(unsafe.Add(unsafe.Pointer(tempIndices2), s*MAX_LPC_ORDER))), res_Q10[:], W_adj_Q5[:], pred_Q8[:], ec_ix[:], []uint8(psNLSF_CB.Ec_Rates_Q5), int(psNLSF_CB.QuantStepSize_Q16), psNLSF_CB.InvQuantStepSize_Q6, int32(NLSF_mu_Q20), psNLSF_CB.Order)
		iCDF_ptr = (*uint8)(unsafe.Pointer(&psNLSF_CB.CB1_iCDF[(signalType>>1)*int(psNLSF_CB.NVectors)]))
		if ind1 == 0 {
			prob_Q8 = 256 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(iCDF_ptr), ind1)))
		} else {
			prob_Q8 = int(*(*uint8)(unsafe.Add(unsafe.Pointer(iCDF_ptr), ind1-1))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(iCDF_ptr), ind1)))
		}
		bits_q7 = (8 << 7) - int(silk_lin2log(int32(prob_Q8)))
		*(*int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(int32(0))*uintptr(s))) = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(int32(0))*uintptr(s)))) + int(int32(int16(bits_q7)))*int(int32(int16(NLSF_mu_Q20>>2))))
	}
	silk_insertion_sort_increasing([]int32(RD_Q25), []int(&bestIndex), nSurvivors, 1)
	*NLSFIndices = int8(*(*int)(unsafe.Add(unsafe.Pointer(tempIndices1), unsafe.Sizeof(int(0))*uintptr(bestIndex))))
	libc.MemCpy(unsafe.Pointer((*int8)(unsafe.Add(unsafe.Pointer(NLSFIndices), 1))), unsafe.Pointer((*int8)(unsafe.Add(unsafe.Pointer(tempIndices2), bestIndex*MAX_LPC_ORDER))), int(uintptr(psNLSF_CB.Order)*unsafe.Sizeof(int8(0))))
	silk_NLSF_decode([]int16(pNLSF_Q15), []int8(NLSFIndices), psNLSF_CB)
	ret = *(*int32)(unsafe.Add(unsafe.Pointer(RD_Q25), unsafe.Sizeof(int32(0))*0))
	return ret
}
