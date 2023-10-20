package libopus

import "unsafe"

func combine_pulses(out *int, in *int, len_ int) {
	var k int
	for k = 0; k < len_; k++ {
		*(*int)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int(0))*uintptr(k))) = *(*int)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int(0))*uintptr(k*2))) + *(*int)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int(0))*uintptr(k*2+1)))
	}
}
func encode_split(psRangeEnc *ec_enc, p_child1 int, p int, shell_table *uint8) {
	if p > 0 {
		ec_enc_icdf(psRangeEnc, p_child1, []byte((*uint8)(unsafe.Add(unsafe.Pointer(shell_table), silk_shell_code_table_offsets[p]))), 8)
	}
}
func decode_split(p_child1 *int16, p_child2 *int16, psRangeDec *ec_dec, p int, shell_table *uint8) {
	if p > 0 {
		*(*int16)(unsafe.Add(unsafe.Pointer(p_child1), unsafe.Sizeof(int16(0))*0)) = int16(ec_dec_icdf(psRangeDec, []byte((*uint8)(unsafe.Add(unsafe.Pointer(shell_table), silk_shell_code_table_offsets[p]))), 8))
		*(*int16)(unsafe.Add(unsafe.Pointer(p_child2), unsafe.Sizeof(int16(0))*0)) = int16(p - int(*(*int16)(unsafe.Add(unsafe.Pointer(p_child1), unsafe.Sizeof(int16(0))*0))))
	} else {
		*(*int16)(unsafe.Add(unsafe.Pointer(p_child1), unsafe.Sizeof(int16(0))*0)) = 0
		*(*int16)(unsafe.Add(unsafe.Pointer(p_child2), unsafe.Sizeof(int16(0))*0)) = 0
	}
}
func silk_shell_encoder(psRangeEnc *ec_enc, pulses0 *int) {
	var (
		pulses1 [8]int
		pulses2 [4]int
		pulses3 [2]int
		pulses4 [1]int
	)
	combine_pulses(&pulses1[0], pulses0, 8)
	combine_pulses(&pulses2[0], &pulses1[0], 4)
	combine_pulses(&pulses3[0], &pulses2[0], 2)
	combine_pulses(&pulses4[0], &pulses3[0], 1)
	encode_split(psRangeEnc, pulses3[0], pulses4[0], &silk_shell_code_table3[0])
	encode_split(psRangeEnc, pulses2[0], pulses3[0], &silk_shell_code_table2[0])
	encode_split(psRangeEnc, pulses1[0], pulses2[0], &silk_shell_code_table1[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*0)), pulses1[0], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*2)), pulses1[1], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, pulses1[2], pulses2[1], &silk_shell_code_table1[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*4)), pulses1[2], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*6)), pulses1[3], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, pulses2[2], pulses3[1], &silk_shell_code_table2[0])
	encode_split(psRangeEnc, pulses1[4], pulses2[2], &silk_shell_code_table1[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*8)), pulses1[4], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*10)), pulses1[5], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, pulses1[6], pulses2[3], &silk_shell_code_table1[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*12)), pulses1[6], &silk_shell_code_table0[0])
	encode_split(psRangeEnc, *(*int)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int(0))*14)), pulses1[7], &silk_shell_code_table0[0])
}
func silk_shell_decoder(pulses0 *int16, psRangeDec *ec_dec, pulses4 int) {
	var (
		pulses3 [2]int16
		pulses2 [4]int16
		pulses1 [8]int16
	)
	decode_split(&pulses3[0], &pulses3[1], psRangeDec, pulses4, &silk_shell_code_table3[0])
	decode_split(&pulses2[0], &pulses2[1], psRangeDec, int(pulses3[0]), &silk_shell_code_table2[0])
	decode_split(&pulses1[0], &pulses1[1], psRangeDec, int(pulses2[0]), &silk_shell_code_table1[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*0)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*1)), psRangeDec, int(pulses1[0]), &silk_shell_code_table0[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*2)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*3)), psRangeDec, int(pulses1[1]), &silk_shell_code_table0[0])
	decode_split(&pulses1[2], &pulses1[3], psRangeDec, int(pulses2[1]), &silk_shell_code_table1[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*4)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*5)), psRangeDec, int(pulses1[2]), &silk_shell_code_table0[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*6)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*7)), psRangeDec, int(pulses1[3]), &silk_shell_code_table0[0])
	decode_split(&pulses2[2], &pulses2[3], psRangeDec, int(pulses3[1]), &silk_shell_code_table2[0])
	decode_split(&pulses1[4], &pulses1[5], psRangeDec, int(pulses2[2]), &silk_shell_code_table1[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*8)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*9)), psRangeDec, int(pulses1[4]), &silk_shell_code_table0[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*10)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*11)), psRangeDec, int(pulses1[5]), &silk_shell_code_table0[0])
	decode_split(&pulses1[6], &pulses1[7], psRangeDec, int(pulses2[3]), &silk_shell_code_table1[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*12)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*13)), psRangeDec, int(pulses1[6]), &silk_shell_code_table0[0])
	decode_split((*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*14)), (*int16)(unsafe.Add(unsafe.Pointer(pulses0), unsafe.Sizeof(int16(0))*15)), psRangeDec, int(pulses1[7]), &silk_shell_code_table0[0])
}
