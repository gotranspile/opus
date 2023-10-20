package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_decode_frame(psDec *silk_decoder_state, psRangeDec *ec_dec, pOut []int16, pN *int32, lostFlag int, condCoding int, arch int) int {
	var (
		psDecCtrl *silk_decoder_control
		L         int
		mv_len    int
		ret       int = 0
	)
	L = psDec.Frame_length
	psDecCtrl = (*silk_decoder_control)(libc.Malloc(int(unsafe.Sizeof(silk_decoder_control{}) * 1)))
	psDecCtrl.LTP_scale_Q14 = 0
	if lostFlag == FLAG_DECODE_NORMAL || lostFlag == FLAG_DECODE_LBRR && psDec.LBRR_flags[psDec.NFramesDecoded] == 1 {
		var pulses *int16
		pulses = (*int16)(libc.Malloc(((L + SHELL_CODEC_FRAME_LENGTH - 1) & ^(int(SHELL_CODEC_FRAME_LENGTH - 1))) * int(unsafe.Sizeof(int16(0)))))
		silk_decode_indices(psDec, psRangeDec, psDec.NFramesDecoded, lostFlag, condCoding)
		silk_decode_pulses(psRangeDec, []int16(pulses), int(psDec.Indices.SignalType), int(psDec.Indices.QuantOffsetType), psDec.Frame_length)
		silk_decode_parameters(psDec, psDecCtrl, condCoding)
		silk_decode_core(psDec, psDecCtrl, pOut, [320]int16(pulses), arch)
		silk_PLC(psDec, psDecCtrl, pOut, 0, arch)
		psDec.LossCnt = 0
		psDec.PrevSignalType = int(psDec.Indices.SignalType)
		psDec.First_frame_after_reset = 0
	} else {
		silk_PLC(psDec, psDecCtrl, pOut, 1, arch)
	}
	mv_len = psDec.Ltp_mem_length - psDec.Frame_length
	libc.MemMove(unsafe.Pointer(&psDec.OutBuf[0]), unsafe.Pointer(&psDec.OutBuf[psDec.Frame_length]), mv_len*int(unsafe.Sizeof(int16(0))))
	libc.MemCpy(unsafe.Pointer(&psDec.OutBuf[mv_len]), unsafe.Pointer(&pOut[0]), psDec.Frame_length*int(unsafe.Sizeof(int16(0))))
	silk_CNG(psDec, psDecCtrl, pOut, L)
	silk_PLC_glue_frames(psDec, pOut, L)
	psDec.LagPrev = psDecCtrl.PitchL[psDec.Nb_subfr-1]
	*pN = int32(L)
	return ret
}
