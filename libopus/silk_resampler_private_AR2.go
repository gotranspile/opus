package libopus

func silk_resampler_private_AR2(S [0]opus_int32, out_Q8 [0]opus_int32, in [0]opus_int16, A_Q14 [0]opus_int16, len_ opus_int32) {
	var (
		k     opus_int32
		out32 opus_int32
	)
	for k = 0; k < len_; k++ {
		out32 = (S[0]) + (opus_int32(opus_uint32(opus_int32(in[k])) << 8))
		out_Q8[k] = out32
		out32 = opus_int32(opus_uint32(out32) << 2)
		S[0] = (S[1]) + ((out32 * opus_int32(int64(A_Q14[0]))) >> 16)
		S[1] = (out32 * opus_int32(int64(A_Q14[1]))) >> 16
	}
}
