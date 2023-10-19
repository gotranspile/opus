package libopus

func silk_interpolate(xi [16]opus_int16, x0 [16]opus_int16, x1 [16]opus_int16, ifact_Q2 int64, d int64) {
	var i int64
	for i = 0; i < d; i++ {
		xi[i] = opus_int16(opus_int32(x0[i]) + ((opus_int32(x1[i]-x0[i]) * opus_int32(opus_int16(ifact_Q2))) >> 2))
	}
}
