package silk

import (
	"math"
)

var silk_TargetRate_NB_21 [107]uint8 = [107]uint8{0, 15, 39, 52, 61, 68, 74, 79, 84, 88, 92, 95, 99, 102, 105, 108, 111, 114, 117, 119, 122, 124, 126, 129, 131, 133, 135, 137, 139, 142, 143, 145, 147, 149, 151, 153, 155, 157, 158, 160, 162, 163, 165, 167, 168, 170, 171, 173, 174, 176, 177, 179, 180, 182, 183, 185, 186, 187, 189, 190, 192, 193, 194, 196, 197, 199, 200, 201, 203, 204, 205, 207, 208, 209, 211, 212, 213, 215, 216, 217, 219, 220, 221, 223, 224, 225, 227, 228, 230, 231, 232, 234, 235, 236, 238, 239, 241, 242, 243, 245, 246, 248, 249, 250, 252, 253, math.MaxUint8}
var silk_TargetRate_MB_21 [155]uint8 = [155]uint8{0, 0, 28, 43, 52, 59, 65, 70, 74, 78, 81, 85, 87, 90, 93, 95, 98, 100, 102, 105, 107, 109, 111, 113, 115, 116, 118, 120, 122, 123, 125, math.MaxInt8, 128, 130, 131, 133, 134, 136, 137, 138, 140, 141, 143, 144, 145, 147, 148, 149, 151, 152, 153, 154, 156, 157, 158, 159, 160, 162, 163, 164, 165, 166, 167, 168, 169, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, math.MaxUint8}
var silk_TargetRate_WB_21 [191]uint8 = [191]uint8{0, 0, 0, 8, 29, 41, 49, 56, 62, 66, 70, 74, 77, 80, 83, 86, 88, 91, 93, 95, 97, 99, 101, 103, 105, 107, 108, 110, 112, 113, 115, 116, 118, 119, 121, 122, 123, 125, 126, math.MaxInt8, 129, 130, 131, 132, 134, 135, 136, 137, 138, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 156, 157, 158, 159, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 171, 172, 173, 174, 175, 176, 177, 177, 178, 179, 180, 181, 181, 182, 183, 184, 185, 185, 186, 187, 188, 189, 189, 190, 191, 192, 192, 193, 194, 195, 195, 196, 197, 198, 198, 199, 200, 200, 201, 202, 203, 203, 204, 205, 206, 206, 207, 208, 209, 209, 210, 211, 211, 212, 213, 214, 214, 215, 216, 216, 217, 218, 219, 219, 220, 221, 221, 222, 223, 224, 224, 225, 226, 226, 227, 228, 229, 229, 230, 231, 232, 232, 233, 234, 234, 235, 236, 237, 237, 238, 239, 240, 240, 241, 242, 243, 243, 244, 245, 246, 246, 247, 248, 249, 249, 250, 251, 252, 253, math.MaxUint8}

func silk_control_SNR(psEncC *EncoderState, TargetRate_bps int32) int {
	psEncC.TargetRate_bps = TargetRate_bps
	if psEncC.Nb_subfr == 2 {
		TargetRate_bps -= int32(psEncC.Fs_kHz/16 + 2000)
	}
	var snr_table []byte
	if psEncC.Fs_kHz == 8 {
		snr_table = silk_TargetRate_NB_21[:]
	} else if psEncC.Fs_kHz == 12 {
		snr_table = silk_TargetRate_MB_21[:]
	} else {
		snr_table = silk_TargetRate_WB_21[:]
	}
	id := (int(TargetRate_bps) + 200) / 400
	if (id - 10) < (len(snr_table) - 1) {
		id = id - 10
	} else {
		id = len(snr_table) - 1
	}
	if id <= 0 {
		psEncC.SNR_dB_Q7 = 0
	} else {
		psEncC.SNR_dB_Q7 = int(snr_table[id]) * 21
	}
	return SILK_NO_ERROR
}
