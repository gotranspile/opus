package libopus

import "unsafe"

func silk_NLSF_unpack(ec_ix [0]opus_int16, pred_Q8 [0]uint8, psNLSF_CB *silk_NLSF_CB_struct, CB1_index int64) {
	var (
		i          int64
		entry      uint8
		ec_sel_ptr *uint8
	)
	ec_sel_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.Ec_sel), CB1_index*int64(psNLSF_CB.Order)/2))
	for i = 0; i < int64(psNLSF_CB.Order); i += 2 {
		entry = *func() *uint8 {
			p := &ec_sel_ptr
			x := *p
			*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
		ec_ix[i] = opus_int16(opus_int32(opus_int16((int64(entry)>>1)&7)) * opus_int32(opus_int16(NLSF_QUANT_MAX_AMPLITUDE*2+1)))
		pred_Q8[i] = *(*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.Pred_Q8), i+(int64(entry)&1)*int64(psNLSF_CB.Order-1)))
		ec_ix[i+1] = opus_int16(opus_int32(opus_int16((int64(entry)>>5)&7)) * opus_int32(opus_int16(NLSF_QUANT_MAX_AMPLITUDE*2+1)))
		pred_Q8[i+1] = *(*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.Pred_Q8), i+((int64(entry)>>4)&1)*int64(psNLSF_CB.Order-1)+1))
	}
}
