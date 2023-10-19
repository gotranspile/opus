package libopus

const RESAMPLER_DOWN_ORDER_FIR0 = 18
const RESAMPLER_DOWN_ORDER_FIR1 = 24
const RESAMPLER_DOWN_ORDER_FIR2 = 36
const RESAMPLER_ORDER_FIR_12 = 8

var silk_resampler_down2_0 opus_int16 = 9872
var silk_resampler_down2_1 opus_int16 = 39809 - 65536
var silk_resampler_up2_hq_0 [3]opus_int16 = [3]opus_int16{1746, 14986, 39083 - 65536}
var silk_resampler_up2_hq_1 [3]opus_int16 = [3]opus_int16{6854, 25769, 55542 - 65536}
var silk_Resampler_3_4_COEFS [29]opus_int16 = [29]opus_int16{-20694, -13867, -49, 64, 17, -157, 353, -496, 163, 11047, 22205, -39, 6, 91, -170, 186, 23, -896, 6336, 19928, -19, -36, 102, -89, -24, 328, -951, 2568, 15909}
var silk_Resampler_2_3_COEFS [20]opus_int16 = [20]opus_int16{-14457, -14019, 64, 128, -122, 36, 310, -768, 584, 9267, 17733, 12, 128, 18, -142, 288, -117, -865, 4123, 14459}
var silk_Resampler_1_2_COEFS [14]opus_int16 = [14]opus_int16{616, -14323, -10, 39, 58, -46, -84, 120, 184, -315, -541, 1284, 5380, 9024}
var silk_Resampler_1_3_COEFS [20]opus_int16 = [20]opus_int16{16102, -15162, -13, 0, 20, 26, 5, -31, -43, -4, 65, 90, 7, -157, -248, -44, 593, 1583, 2612, 3271}
var silk_Resampler_1_4_COEFS [20]opus_int16 = [20]opus_int16{22500, -15099, 3, -14, -20, -15, 2, 25, 37, 25, -16, -71, -107, -79, 50, 292, 623, 982, 1288, 1464}
var silk_Resampler_1_6_COEFS [20]opus_int16 = [20]opus_int16{27540, -15257, 17, 12, 8, 1, -10, -22, -30, -32, -22, 3, 44, 100, 168, 243, 317, 381, 429, 455}
var silk_Resampler_2_3_COEFS_LQ [6]opus_int16 = [6]opus_int16{-2797, -6507, 4697, 10739, 1567, 8276}
var silk_resampler_frac_FIR_12 [12][4]opus_int16 = [12][4]opus_int16{{189, -600, 617, 30567}, {117, -159, -1070, 29704}, {52, 221, -2392, 28276}, {-4, 529, -3350, 26341}, {-48, 758, -3956, 23973}, {-80, 905, -4235, 21254}, {-99, 972, -4222, 18278}, {-107, 967, -3957, 15143}, {-103, 896, -3487, 11950}, {-91, 773, -2865, 8798}, {-71, 611, -2143, 5784}, {-46, 425, -1375, 2996}}
