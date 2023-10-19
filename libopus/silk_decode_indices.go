package libopus

import "unsafe"

func silk_decode_indices(psDec *silk_decoder_state, psRangeDec *ec_dec, FrameIndex int64, decode_LBRR int64, condCoding int64) {
	var (
		i                        int64
		k                        int64
		Ix                       int64
		decode_absolute_lagIndex int64
		delta_lagIndex           int64
		ec_ix                    [16]opus_int16
		pred_Q8                  [16]uint8
	)
	if decode_LBRR != 0 || psDec.VAD_flags[FrameIndex] != 0 {
		Ix = ec_dec_icdf(psRangeDec, &silk_type_offset_VAD_iCDF[0], 8) + 2
	} else {
		Ix = ec_dec_icdf(psRangeDec, &silk_type_offset_no_VAD_iCDF[0], 8)
	}
	psDec.Indices.SignalType = int8(Ix >> 1)
	psDec.Indices.QuantOffsetType = int8(Ix & 1)
	if condCoding == CODE_CONDITIONALLY {
		psDec.Indices.GainsIndices[0] = int8(ec_dec_icdf(psRangeDec, &silk_delta_gain_iCDF[0], 8))
	} else {
		psDec.Indices.GainsIndices[0] = int8(opus_int32(opus_uint32(ec_dec_icdf(psRangeDec, &silk_gain_iCDF[psDec.Indices.SignalType][0], 8)) << 3))
		psDec.Indices.GainsIndices[0] += int8(ec_dec_icdf(psRangeDec, &silk_uniform8_iCDF[0], 8))
	}
	for i = 1; i < psDec.Nb_subfr; i++ {
		psDec.Indices.GainsIndices[i] = int8(ec_dec_icdf(psRangeDec, &silk_delta_gain_iCDF[0], 8))
	}
	psDec.Indices.NLSFIndices[0] = int8(ec_dec_icdf(psRangeDec, (*uint8)(unsafe.Add(unsafe.Pointer(psDec.PsNLSF_CB.CB1_iCDF), (int64(psDec.Indices.SignalType)>>1)*int64(psDec.PsNLSF_CB.NVectors))), 8))
	silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psDec.PsNLSF_CB, int64(psDec.Indices.NLSFIndices[0]))
	for i = 0; i < int64(psDec.PsNLSF_CB.Order); i++ {
		Ix = ec_dec_icdf(psRangeDec, (*uint8)(unsafe.Add(unsafe.Pointer(psDec.PsNLSF_CB.Ec_iCDF), ec_ix[i])), 8)
		if Ix == 0 {
			Ix -= ec_dec_icdf(psRangeDec, &silk_NLSF_EXT_iCDF[0], 8)
		} else if Ix == NLSF_QUANT_MAX_AMPLITUDE*2 {
			Ix += ec_dec_icdf(psRangeDec, &silk_NLSF_EXT_iCDF[0], 8)
		}
		psDec.Indices.NLSFIndices[i+1] = int8(Ix - NLSF_QUANT_MAX_AMPLITUDE)
	}
	if psDec.Nb_subfr == MAX_NB_SUBFR {
		psDec.Indices.NLSFInterpCoef_Q2 = int8(ec_dec_icdf(psRangeDec, &silk_NLSF_interpolation_factor_iCDF[0], 8))
	} else {
		psDec.Indices.NLSFInterpCoef_Q2 = 4
	}
	if int64(psDec.Indices.SignalType) == TYPE_VOICED {
		decode_absolute_lagIndex = 1
		if condCoding == CODE_CONDITIONALLY && psDec.Ec_prevSignalType == TYPE_VOICED {
			delta_lagIndex = int64(opus_int16(ec_dec_icdf(psRangeDec, &silk_pitch_delta_iCDF[0], 8)))
			if delta_lagIndex > 0 {
				delta_lagIndex = delta_lagIndex - 9
				psDec.Indices.LagIndex = opus_int16(int64(psDec.Ec_prevLagIndex) + delta_lagIndex)
				decode_absolute_lagIndex = 0
			}
		}
		if decode_absolute_lagIndex != 0 {
			psDec.Indices.LagIndex = opus_int16(int64(opus_int16(ec_dec_icdf(psRangeDec, &silk_pitch_lag_iCDF[0], 8))) * (psDec.Fs_kHz >> 1))
			psDec.Indices.LagIndex += opus_int16(ec_dec_icdf(psRangeDec, psDec.Pitch_lag_low_bits_iCDF, 8))
		}
		psDec.Ec_prevLagIndex = psDec.Indices.LagIndex
		psDec.Indices.ContourIndex = int8(ec_dec_icdf(psRangeDec, psDec.Pitch_contour_iCDF, 8))
		psDec.Indices.PERIndex = int8(ec_dec_icdf(psRangeDec, &silk_LTP_per_index_iCDF[0], 8))
		for k = 0; k < psDec.Nb_subfr; k++ {
			psDec.Indices.LTPIndex[k] = int8(ec_dec_icdf(psRangeDec, silk_LTP_gain_iCDF_ptrs[psDec.Indices.PERIndex], 8))
		}
		if condCoding == CODE_INDEPENDENTLY {
			psDec.Indices.LTP_scaleIndex = int8(ec_dec_icdf(psRangeDec, &silk_LTPscale_iCDF[0], 8))
		} else {
			psDec.Indices.LTP_scaleIndex = 0
		}
	}
	psDec.Ec_prevSignalType = int64(psDec.Indices.SignalType)
	psDec.Indices.Seed = int8(ec_dec_icdf(psRangeDec, &silk_uniform4_iCDF[0], 8))
}
