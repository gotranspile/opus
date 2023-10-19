package libopus

func silk_lin2log(inLin opus_int32) opus_int32 {
	var (
		lz      opus_int32
		frac_Q7 opus_int32
	)
	silk_CLZ_FRAC(inLin, &lz, &frac_Q7)
	return (frac_Q7 + (((frac_Q7 * (128 - frac_Q7)) * 179) >> 16)) + (opus_int32(opus_uint32(31-lz) << 7))
}
