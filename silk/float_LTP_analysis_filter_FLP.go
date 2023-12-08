package silk

func LTP_analysis_filter_FLP(LTP_res []float32, x []float32, B [20]float32, pitchL [4]int, invGains [4]float32, subfr_length int, nb_subfr int, pre_length int) {
	var (
		Btmp     [5]float32
		inv_gain float32
	)
	x_ptr := x
	LTP_res_ptr := LTP_res
	for k := 0; k < nb_subfr; k++ {
		// FIXME
		x_lag_ptr := x_ptr[-pitchL[k]:]
		inv_gain = invGains[k]
		for i := 0; i < LTP_ORDER; i++ {
			Btmp[i] = B[k*LTP_ORDER+i]
		}
		for i := 0; i < subfr_length+pre_length; i++ {
			LTP_res_ptr[i] = x_ptr[i]
			for j := 0; j < LTP_ORDER; j++ {
				LTP_res_ptr[i] -= Btmp[j] * x_lag_ptr[(LTP_ORDER/2)-j]
			}
			LTP_res_ptr[i] *= inv_gain
			x_lag_ptr = x_lag_ptr[1:]
		}
		LTP_res_ptr = LTP_res_ptr[subfr_length+pre_length:]
		x_ptr = x_ptr[subfr_length:]
	}
}
