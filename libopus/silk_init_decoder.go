package libopus

func silk_init_decoder(psDec *silk_decoder_state) int {
	*psDec = silk_decoder_state{}
	psDec.First_frame_after_reset = 1
	psDec.Prev_gain_Q16 = 65536
	psDec.Arch = opus_select_arch()
	silk_CNG_Reset(psDec)
	silk_PLC_Reset(psDec)
	return 0
}
