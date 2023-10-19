package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_NLSF_del_dec_quant(indices [0]int8, x_Q10 [0]opus_int16, w_Q5 [0]opus_int16, pred_coef_Q8 [0]uint8, ec_ix [0]opus_int16, ec_rates_Q5 [0]uint8, quant_step_size_Q16 int64, inv_quant_step_size_Q6 opus_int16, mu_Q20 opus_int32, order opus_int16) opus_int32 {
	var (
		i              int64
		j              int64
		nStates        int64
		ind_tmp        int64
		ind_min_max    int64
		ind_max_min    int64
		in_Q10         int64
		res_Q10        int64
		pred_Q10       int64
		diff_Q10       int64
		rate0_Q5       int64
		rate1_Q5       int64
		out0_Q10       opus_int16
		out1_Q10       opus_int16
		RD_tmp_Q25     opus_int32
		min_Q25        opus_int32
		min_max_Q25    opus_int32
		max_min_Q25    opus_int32
		ind_sort       [4]int64
		ind            [4][16]int8
		prev_out_Q10   [8]opus_int16
		RD_Q25         [8]opus_int32
		RD_min_Q25     [4]opus_int32
		RD_max_Q25     [4]opus_int32
		rates_Q5       *uint8
		out0_Q10_table [20]int64
		out1_Q10_table [20]int64
	)
	for i = -NLSF_QUANT_MAX_AMPLITUDE_EXT; i <= NLSF_QUANT_MAX_AMPLITUDE_EXT-1; i++ {
		out0_Q10 = opus_int16(opus_int32(opus_uint32(i) << 10))
		out1_Q10 = out0_Q10 + 1024
		if i > 0 {
			out0_Q10 = opus_int16(opus_int32(out0_Q10) - (opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10) + 0.5)))
			out1_Q10 = opus_int16(opus_int32(out1_Q10) - (opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10) + 0.5)))
		} else if i == 0 {
			out1_Q10 = opus_int16(opus_int32(out1_Q10) - (opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10) + 0.5)))
		} else if i == -1 {
			out0_Q10 = opus_int16(opus_int32(out0_Q10) + (opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10) + 0.5)))
		} else {
			out0_Q10 = opus_int16(opus_int32(out0_Q10) + (opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10) + 0.5)))
			out1_Q10 = opus_int16(opus_int32(out1_Q10) + (opus_int32(NLSF_QUANT_LEVEL_ADJ*(1<<10) + 0.5)))
		}
		out0_Q10_table[i+NLSF_QUANT_MAX_AMPLITUDE_EXT] = int64((opus_int32(out0_Q10) * opus_int32(opus_int16(quant_step_size_Q16))) >> 16)
		out1_Q10_table[i+NLSF_QUANT_MAX_AMPLITUDE_EXT] = int64((opus_int32(out1_Q10) * opus_int32(opus_int16(quant_step_size_Q16))) >> 16)
	}
	nStates = 1
	RD_Q25[0] = 0
	prev_out_Q10[0] = 0
	for i = int64(order - 1); i >= 0; i-- {
		rates_Q5 = &ec_rates_Q5[ec_ix[i]]
		in_Q10 = int64(x_Q10[i])
		for j = 0; j < nStates; j++ {
			pred_Q10 = int64((opus_int32(opus_int16(pred_coef_Q8[i])) * opus_int32(prev_out_Q10[j])) >> 8)
			res_Q10 = in_Q10 - pred_Q10
			ind_tmp = int64((opus_int32(inv_quant_step_size_Q6) * opus_int32(opus_int16(res_Q10))) >> 16)
			if (-NLSF_QUANT_MAX_AMPLITUDE_EXT) > (NLSF_QUANT_MAX_AMPLITUDE_EXT - 1) {
				if ind_tmp > (-NLSF_QUANT_MAX_AMPLITUDE_EXT) {
					ind_tmp = -NLSF_QUANT_MAX_AMPLITUDE_EXT
				} else if ind_tmp < (NLSF_QUANT_MAX_AMPLITUDE_EXT - 1) {
					ind_tmp = NLSF_QUANT_MAX_AMPLITUDE_EXT - 1
				} else {
					ind_tmp = ind_tmp
				}
			} else if ind_tmp > (NLSF_QUANT_MAX_AMPLITUDE_EXT - 1) {
				ind_tmp = NLSF_QUANT_MAX_AMPLITUDE_EXT - 1
			} else if ind_tmp < (-NLSF_QUANT_MAX_AMPLITUDE_EXT) {
				ind_tmp = -NLSF_QUANT_MAX_AMPLITUDE_EXT
			} else {
				ind_tmp = ind_tmp
			}
			ind[j][i] = int8(ind_tmp)
			out0_Q10 = opus_int16(out0_Q10_table[ind_tmp+NLSF_QUANT_MAX_AMPLITUDE_EXT])
			out1_Q10 = opus_int16(out1_Q10_table[ind_tmp+NLSF_QUANT_MAX_AMPLITUDE_EXT])
			out0_Q10 = opus_int16(int64(out0_Q10) + pred_Q10)
			out1_Q10 = opus_int16(int64(out1_Q10) + pred_Q10)
			prev_out_Q10[j] = out0_Q10
			prev_out_Q10[j+nStates] = out1_Q10
			if ind_tmp+1 >= NLSF_QUANT_MAX_AMPLITUDE {
				if ind_tmp+1 == NLSF_QUANT_MAX_AMPLITUDE {
					rate0_Q5 = int64(*(*uint8)(unsafe.Add(unsafe.Pointer(rates_Q5), ind_tmp+NLSF_QUANT_MAX_AMPLITUDE)))
					rate1_Q5 = 280
				} else {
					rate0_Q5 = int64(opus_int32(280-NLSF_QUANT_MAX_AMPLITUDE*43) + opus_int32(opus_int16(ind_tmp))*43)
					rate1_Q5 = rate0_Q5 + 43
				}
			} else if ind_tmp <= -NLSF_QUANT_MAX_AMPLITUDE {
				if ind_tmp == -NLSF_QUANT_MAX_AMPLITUDE {
					rate0_Q5 = 280
					rate1_Q5 = int64(*(*uint8)(unsafe.Add(unsafe.Pointer(rates_Q5), ind_tmp+1+NLSF_QUANT_MAX_AMPLITUDE)))
				} else {
					rate0_Q5 = int64(opus_int32(280-NLSF_QUANT_MAX_AMPLITUDE*43) + opus_int32(opus_int16(ind_tmp))*(opus_int32(opus_int16(-43))))
					rate1_Q5 = rate0_Q5 - 43
				}
			} else {
				rate0_Q5 = int64(*(*uint8)(unsafe.Add(unsafe.Pointer(rates_Q5), ind_tmp+NLSF_QUANT_MAX_AMPLITUDE)))
				rate1_Q5 = int64(*(*uint8)(unsafe.Add(unsafe.Pointer(rates_Q5), ind_tmp+1+NLSF_QUANT_MAX_AMPLITUDE)))
			}
			RD_tmp_Q25 = RD_Q25[j]
			diff_Q10 = in_Q10 - int64(out0_Q10)
			RD_Q25[j] = (RD_tmp_Q25 + (opus_int32(opus_int16(diff_Q10))*opus_int32(opus_int16(diff_Q10)))*opus_int32(w_Q5[i])) + (opus_int32(opus_int16(mu_Q20)))*opus_int32(opus_int16(rate0_Q5))
			diff_Q10 = in_Q10 - int64(out1_Q10)
			RD_Q25[j+nStates] = (RD_tmp_Q25 + (opus_int32(opus_int16(diff_Q10))*opus_int32(opus_int16(diff_Q10)))*opus_int32(w_Q5[i])) + (opus_int32(opus_int16(mu_Q20)))*opus_int32(opus_int16(rate1_Q5))
		}
		if nStates <= (1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)/2 {
			for j = 0; j < nStates; j++ {
				ind[j+nStates][i] = int8(int64(ind[j][i]) + 1)
			}
			nStates = int64(opus_int32(opus_uint32(nStates) << 1))
			for j = nStates; j < (1 << NLSF_QUANT_DEL_DEC_STATES_LOG2); j++ {
				ind[j][i] = ind[j-nStates][i]
			}
		} else {
			for j = 0; j < (1 << NLSF_QUANT_DEL_DEC_STATES_LOG2); j++ {
				if RD_Q25[j] > RD_Q25[j+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)] {
					RD_max_Q25[j] = RD_Q25[j]
					RD_min_Q25[j] = RD_Q25[j+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)]
					RD_Q25[j] = RD_min_Q25[j]
					RD_Q25[j+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)] = RD_max_Q25[j]
					out0_Q10 = prev_out_Q10[j]
					prev_out_Q10[j] = prev_out_Q10[j+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)]
					prev_out_Q10[j+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)] = out0_Q10
					ind_sort[j] = j + (1 << NLSF_QUANT_DEL_DEC_STATES_LOG2)
				} else {
					RD_min_Q25[j] = RD_Q25[j]
					RD_max_Q25[j] = RD_Q25[j+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)]
					ind_sort[j] = j
				}
			}
			for {
				min_max_Q25 = silk_int32_MAX
				max_min_Q25 = 0
				ind_min_max = 0
				ind_max_min = 0
				for j = 0; j < (1 << NLSF_QUANT_DEL_DEC_STATES_LOG2); j++ {
					if min_max_Q25 > RD_max_Q25[j] {
						min_max_Q25 = RD_max_Q25[j]
						ind_min_max = j
					}
					if max_min_Q25 < RD_min_Q25[j] {
						max_min_Q25 = RD_min_Q25[j]
						ind_max_min = j
					}
				}
				if min_max_Q25 >= max_min_Q25 {
					break
				}
				ind_sort[ind_max_min] = ind_sort[ind_min_max] ^ 1<<NLSF_QUANT_DEL_DEC_STATES_LOG2
				RD_Q25[ind_max_min] = RD_Q25[ind_min_max+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)]
				prev_out_Q10[ind_max_min] = prev_out_Q10[ind_min_max+(1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)]
				RD_min_Q25[ind_max_min] = 0
				RD_max_Q25[ind_min_max] = silk_int32_MAX
				libc.MemCpy(unsafe.Pointer(&(ind[ind_max_min])[0]), unsafe.Pointer(&(ind[ind_min_max])[0]), int(MAX_LPC_ORDER*unsafe.Sizeof(int8(0))))
			}
			for j = 0; j < (1 << NLSF_QUANT_DEL_DEC_STATES_LOG2); j++ {
				ind[j][i] += int8((ind_sort[j]) >> NLSF_QUANT_DEL_DEC_STATES_LOG2)
			}
		}
	}
	ind_tmp = 0
	min_Q25 = silk_int32_MAX
	for j = 0; j < (1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)*2; j++ {
		if min_Q25 > RD_Q25[j] {
			min_Q25 = RD_Q25[j]
			ind_tmp = j
		}
	}
	for j = 0; j < int64(order); j++ {
		indices[j] = ind[ind_tmp&((1<<NLSF_QUANT_DEL_DEC_STATES_LOG2)-1)][j]
	}
	indices[0] += int8(ind_tmp >> NLSF_QUANT_DEL_DEC_STATES_LOG2)
	return min_Q25
}
