package libopus

import "unsafe"

func silk_NLSF_unpack(ec_ix []int16, pred_Q8 []uint8, psNLSF_CB *silk_NLSF_CB_struct, CB1_index int) {
	var (
		i          int
		entry      uint8
		ec_sel_ptr *uint8
	)
	ec_sel_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.Ec_sel), CB1_index*int(psNLSF_CB.Order)/2))
	for i = 0; i < int(psNLSF_CB.Order); i += 2 {
		entry = *func() *uint8 {
			p := &ec_sel_ptr
			x := *p
			*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
		ec_ix[i] = int16(int(int32(int16((int(entry)>>1)&7))) * int(int32(int16(int(NLSF_QUANT_MAX_AMPLITUDE*2)+1))))
		pred_Q8[i] = *(*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.Pred_Q8), i+(int(entry)&1)*(int(psNLSF_CB.Order)-1)))
		ec_ix[i+1] = int16(int(int32(int16((int(entry)>>5)&7))) * int(int32(int16(int(NLSF_QUANT_MAX_AMPLITUDE*2)+1))))
		pred_Q8[i+1] = *(*uint8)(unsafe.Add(unsafe.Pointer(psNLSF_CB.Pred_Q8), i+((int(entry)>>4)&1)*(int(psNLSF_CB.Order)-1)+1))
	}
}
