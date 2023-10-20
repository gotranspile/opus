package libopus

func silk_decode_indices(psDec *silk_decoder_state, psRangeDec *ec_dec, FrameIndex int, decode_LBRR int, condCoding int) {
	var (
		i                        int
		k                        int
		Ix                       int
		decode_absolute_lagIndex int
		delta_lagIndex           int
		ec_ix                    [16]int16
		pred_Q8                  [16]uint8
	)
	if decode_LBRR != 0 || psDec.VAD_flags[FrameIndex] != 0 {
		Ix = ec_dec_icdf(psRangeDec, silk_type_offset_VAD_iCDF[:], 8) + 2
	} else {
		Ix = ec_dec_icdf(psRangeDec, silk_type_offset_no_VAD_iCDF[:], 8)
	}
	psDec.Indices.SignalType = int8(Ix >> 1)
	psDec.Indices.QuantOffsetType = int8(Ix & 1)
	if condCoding == CODE_CONDITIONALLY {
		psDec.Indices.GainsIndices[0] = int8(ec_dec_icdf(psRangeDec, silk_delta_gain_iCDF[:], 8))
	} else {
		psDec.Indices.GainsIndices[0] = int8(int32(int(uint32(int32(ec_dec_icdf(psRangeDec, silk_gain_iCDF[psDec.Indices.SignalType][:], 8)))) << 3))
		psDec.Indices.GainsIndices[0] += int8(ec_dec_icdf(psRangeDec, silk_uniform8_iCDF[:], 8))
	}
	for i = 1; i < psDec.Nb_subfr; i++ {
		psDec.Indices.GainsIndices[i] = int8(ec_dec_icdf(psRangeDec, silk_delta_gain_iCDF[:], 8))
	}
	psDec.Indices.NLSFIndices[0] = int8(ec_dec_icdf(psRangeDec, []byte(&psDec.PsNLSF_CB.CB1_iCDF[(int(psDec.Indices.SignalType)>>1)*int(psDec.PsNLSF_CB.NVectors)]), 8))
	silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psDec.PsNLSF_CB, int(psDec.Indices.NLSFIndices[0]))
	for i = 0; i < int(psDec.PsNLSF_CB.Order); i++ {
		Ix = ec_dec_icdf(psRangeDec, []byte(&psDec.PsNLSF_CB.Ec_iCDF[ec_ix[i]]), 8)
		if Ix == 0 {
			Ix -= ec_dec_icdf(psRangeDec, silk_NLSF_EXT_iCDF[:], 8)
		} else if Ix == int(NLSF_QUANT_MAX_AMPLITUDE*2) {
			Ix += ec_dec_icdf(psRangeDec, silk_NLSF_EXT_iCDF[:], 8)
		}
		psDec.Indices.NLSFIndices[i+1] = int8(Ix - NLSF_QUANT_MAX_AMPLITUDE)
	}
	if psDec.Nb_subfr == MAX_NB_SUBFR {
		psDec.Indices.NLSFInterpCoef_Q2 = int8(ec_dec_icdf(psRangeDec, silk_NLSF_interpolation_factor_iCDF[:], 8))
	} else {
		psDec.Indices.NLSFInterpCoef_Q2 = 4
	}
	if int(psDec.Indices.SignalType) == TYPE_VOICED {
		decode_absolute_lagIndex = 1
		if condCoding == CODE_CONDITIONALLY && psDec.Ec_prevSignalType == TYPE_VOICED {
			delta_lagIndex = int(int16(ec_dec_icdf(psRangeDec, silk_pitch_delta_iCDF[:], 8)))
			if delta_lagIndex > 0 {
				delta_lagIndex = delta_lagIndex - 9
				psDec.Indices.LagIndex = int16(int(psDec.Ec_prevLagIndex) + delta_lagIndex)
				decode_absolute_lagIndex = 0
			}
		}
		if decode_absolute_lagIndex != 0 {
			psDec.Indices.LagIndex = int16(int(int16(ec_dec_icdf(psRangeDec, silk_pitch_lag_iCDF[:], 8))) * (psDec.Fs_kHz >> 1))
			psDec.Indices.LagIndex += int16(ec_dec_icdf(psRangeDec, psDec.Pitch_lag_low_bits_iCDF, 8))
		}
		psDec.Ec_prevLagIndex = psDec.Indices.LagIndex
		psDec.Indices.ContourIndex = int8(ec_dec_icdf(psRangeDec, psDec.Pitch_contour_iCDF, 8))
		psDec.Indices.PERIndex = int8(ec_dec_icdf(psRangeDec, silk_LTP_per_index_iCDF[:], 8))
		for k = 0; k < psDec.Nb_subfr; k++ {
			psDec.Indices.LTPIndex[k] = int8(ec_dec_icdf(psRangeDec, []byte(silk_LTP_gain_iCDF_ptrs[psDec.Indices.PERIndex]), 8))
		}
		if condCoding == CODE_INDEPENDENTLY {
			psDec.Indices.LTP_scaleIndex = int8(ec_dec_icdf(psRangeDec, silk_LTPscale_iCDF[:], 8))
		} else {
			psDec.Indices.LTP_scaleIndex = 0
		}
	}
	psDec.Ec_prevSignalType = int(psDec.Indices.SignalType)
	psDec.Indices.Seed = int8(ec_dec_icdf(psRangeDec, silk_uniform4_iCDF[:], 8))
}
