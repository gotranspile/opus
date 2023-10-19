package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_decode_parameters(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, condCoding int64) {
	var (
		i          int64
		k          int64
		Ix         int64
		pNLSF_Q15  [16]opus_int16
		pNLSF0_Q15 [16]opus_int16
		cbk_ptr_Q7 *int8
	)
	silk_gains_dequant(psDecCtrl.Gains_Q16, psDec.Indices.GainsIndices, &psDec.LastGainIndex, int64(libc.BoolToInt(condCoding == CODE_CONDITIONALLY)), psDec.Nb_subfr)
	silk_NLSF_decode(&pNLSF_Q15[0], &psDec.Indices.NLSFIndices[0], psDec.PsNLSF_CB)
	silk_NLSF2A(&psDecCtrl.PredCoef_Q12[1][0], &pNLSF_Q15[0], psDec.LPC_order, psDec.Arch)
	if psDec.First_frame_after_reset == 1 {
		psDec.Indices.NLSFInterpCoef_Q2 = 4
	}
	if int64(psDec.Indices.NLSFInterpCoef_Q2) < 4 {
		for i = 0; i < psDec.LPC_order; i++ {
			pNLSF0_Q15[i] = psDec.PrevNLSF_Q15[i] + ((opus_int16(psDec.Indices.NLSFInterpCoef_Q2) * (pNLSF_Q15[i] - psDec.PrevNLSF_Q15[i])) >> 2)
		}
		silk_NLSF2A(&psDecCtrl.PredCoef_Q12[0][0], &pNLSF0_Q15[0], psDec.LPC_order, psDec.Arch)
	} else {
		libc.MemCpy(unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[0])[0]), unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[1])[0]), int(psDec.LPC_order*int64(unsafe.Sizeof(opus_int16(0)))))
	}
	libc.MemCpy(unsafe.Pointer(&psDec.PrevNLSF_Q15[0]), unsafe.Pointer(&pNLSF_Q15[0]), int(psDec.LPC_order*int64(unsafe.Sizeof(opus_int16(0)))))
	if psDec.LossCnt != 0 {
		silk_bwexpander(&psDecCtrl.PredCoef_Q12[0][0], psDec.LPC_order, BWE_AFTER_LOSS_Q16)
		silk_bwexpander(&psDecCtrl.PredCoef_Q12[1][0], psDec.LPC_order, BWE_AFTER_LOSS_Q16)
	}
	if int64(psDec.Indices.SignalType) == TYPE_VOICED {
		silk_decode_pitch(psDec.Indices.LagIndex, psDec.Indices.ContourIndex, psDecCtrl.PitchL[:], psDec.Fs_kHz, psDec.Nb_subfr)
		cbk_ptr_Q7 = silk_LTP_vq_ptrs_Q7[psDec.Indices.PERIndex]
		for k = 0; k < psDec.Nb_subfr; k++ {
			Ix = int64(psDec.Indices.LTPIndex[k])
			for i = 0; i < LTP_ORDER; i++ {
				psDecCtrl.LTPCoef_Q14[k*LTP_ORDER+i] = opus_int16(opus_int32(opus_uint32(*(*int8)(unsafe.Add(unsafe.Pointer(cbk_ptr_Q7), Ix*LTP_ORDER+i))) << 7))
			}
		}
		Ix = int64(psDec.Indices.LTP_scaleIndex)
		psDecCtrl.LTP_scale_Q14 = int64(silk_LTPScales_table_Q14[Ix])
	} else {
		libc.MemSet(unsafe.Pointer(&psDecCtrl.PitchL[0]), 0, int(psDec.Nb_subfr*int64(unsafe.Sizeof(int64(0)))))
		libc.MemSet(unsafe.Pointer(&psDecCtrl.LTPCoef_Q14[0]), 0, int(LTP_ORDER*psDec.Nb_subfr*int64(unsafe.Sizeof(opus_int16(0)))))
		psDec.Indices.PERIndex = 0
		psDecCtrl.LTP_scale_Q14 = 0
	}
}
