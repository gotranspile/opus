package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_decode_frame(psDec *silk_decoder_state, psRangeDec *ec_dec, pOut [0]opus_int16, pN *opus_int32, lostFlag int64, condCoding int64, arch int64) int64 {
	var (
		psDecCtrl *silk_decoder_control
		L         int64
		mv_len    int64
		ret       int64 = 0
	)
	L = psDec.Frame_length
	psDecCtrl = (*silk_decoder_control)(libc.Malloc(int(unsafe.Sizeof(silk_decoder_control{}) * 1)))
	psDecCtrl.LTP_scale_Q14 = 0
	if lostFlag == FLAG_DECODE_NORMAL || lostFlag == FLAG_DECODE_LBRR && psDec.LBRR_flags[psDec.NFramesDecoded] == 1 {
		var pulses *opus_int16
		pulses = (*opus_int16)(libc.Malloc(int(((L + SHELL_CODEC_FRAME_LENGTH - 1) & ^(SHELL_CODEC_FRAME_LENGTH - 1)) * int64(unsafe.Sizeof(opus_int16(0))))))
		silk_decode_indices(psDec, psRangeDec, psDec.NFramesDecoded, lostFlag, condCoding)
		silk_decode_pulses(psRangeDec, [0]opus_int16(pulses), int64(psDec.Indices.SignalType), int64(psDec.Indices.QuantOffsetType), psDec.Frame_length)
		silk_decode_parameters(psDec, psDecCtrl, condCoding)
		silk_decode_core(psDec, psDecCtrl, pOut, [320]opus_int16(pulses), arch)
		silk_PLC(psDec, psDecCtrl, pOut, 0, arch)
		psDec.LossCnt = 0
		psDec.PrevSignalType = int64(psDec.Indices.SignalType)
		psDec.First_frame_after_reset = 0
	} else {
		silk_PLC(psDec, psDecCtrl, pOut, 1, arch)
	}
	mv_len = psDec.Ltp_mem_length - psDec.Frame_length
	libc.MemMove(unsafe.Pointer(&psDec.OutBuf[0]), unsafe.Pointer(&psDec.OutBuf[psDec.Frame_length]), int(mv_len*int64(unsafe.Sizeof(opus_int16(0)))))
	libc.MemCpy(unsafe.Pointer(&psDec.OutBuf[mv_len]), unsafe.Pointer(&pOut[0]), int(psDec.Frame_length*int64(unsafe.Sizeof(opus_int16(0)))))
	silk_CNG(psDec, psDecCtrl, pOut, L)
	silk_PLC_glue_frames(psDec, pOut, L)
	psDec.LagPrev = psDecCtrl.PitchL[psDec.Nb_subfr-1]
	*pN = opus_int32(L)
	return ret
}
