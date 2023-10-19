package libopus

import "unsafe"

func silk_encode_indices(psEncC *silk_encoder_state, psRangeEnc *ec_enc, FrameIndex int64, encode_LBRR int64, condCoding int64) {
	var (
		i                        int64
		k                        int64
		typeOffset               int64
		encode_absolute_lagIndex int64
		delta_lagIndex           int64
		ec_ix                    [16]opus_int16
		pred_Q8                  [16]uint8
		psIndices                *SideInfoIndices
	)
	if encode_LBRR != 0 {
		psIndices = &psEncC.Indices_LBRR[FrameIndex]
	} else {
		psIndices = &psEncC.Indices
	}
	typeOffset = int64(psIndices.SignalType)*2 + int64(psIndices.QuantOffsetType)
	if encode_LBRR != 0 || typeOffset >= 2 {
		ec_enc_icdf(psRangeEnc, typeOffset-2, &silk_type_offset_VAD_iCDF[0], 8)
	} else {
		ec_enc_icdf(psRangeEnc, typeOffset, &silk_type_offset_no_VAD_iCDF[0], 8)
	}
	if condCoding == CODE_CONDITIONALLY {
		ec_enc_icdf(psRangeEnc, int64(psIndices.GainsIndices[0]), &silk_delta_gain_iCDF[0], 8)
	} else {
		ec_enc_icdf(psRangeEnc, int64(psIndices.GainsIndices[0])>>3, &silk_gain_iCDF[psIndices.SignalType][0], 8)
		ec_enc_icdf(psRangeEnc, int64(psIndices.GainsIndices[0])&7, &silk_uniform8_iCDF[0], 8)
	}
	for i = 1; i < psEncC.Nb_subfr; i++ {
		ec_enc_icdf(psRangeEnc, int64(psIndices.GainsIndices[i]), &silk_delta_gain_iCDF[0], 8)
	}
	ec_enc_icdf(psRangeEnc, int64(psIndices.NLSFIndices[0]), (*uint8)(unsafe.Add(unsafe.Pointer(psEncC.PsNLSF_CB.CB1_iCDF), (int64(psIndices.SignalType)>>1)*int64(psEncC.PsNLSF_CB.NVectors))), 8)
	silk_NLSF_unpack(ec_ix[:], pred_Q8[:], psEncC.PsNLSF_CB, int64(psIndices.NLSFIndices[0]))
	for i = 0; i < int64(psEncC.PsNLSF_CB.Order); i++ {
		if int64(psIndices.NLSFIndices[i+1]) >= NLSF_QUANT_MAX_AMPLITUDE {
			ec_enc_icdf(psRangeEnc, NLSF_QUANT_MAX_AMPLITUDE*2, (*uint8)(unsafe.Add(unsafe.Pointer(psEncC.PsNLSF_CB.Ec_iCDF), ec_ix[i])), 8)
			ec_enc_icdf(psRangeEnc, int64(psIndices.NLSFIndices[i+1])-NLSF_QUANT_MAX_AMPLITUDE, &silk_NLSF_EXT_iCDF[0], 8)
		} else if int64(psIndices.NLSFIndices[i+1]) <= -NLSF_QUANT_MAX_AMPLITUDE {
			ec_enc_icdf(psRangeEnc, 0, (*uint8)(unsafe.Add(unsafe.Pointer(psEncC.PsNLSF_CB.Ec_iCDF), ec_ix[i])), 8)
			ec_enc_icdf(psRangeEnc, int64(-psIndices.NLSFIndices[i+1])-NLSF_QUANT_MAX_AMPLITUDE, &silk_NLSF_EXT_iCDF[0], 8)
		} else {
			ec_enc_icdf(psRangeEnc, int64(psIndices.NLSFIndices[i+1])+NLSF_QUANT_MAX_AMPLITUDE, (*uint8)(unsafe.Add(unsafe.Pointer(psEncC.PsNLSF_CB.Ec_iCDF), ec_ix[i])), 8)
		}
	}
	if psEncC.Nb_subfr == MAX_NB_SUBFR {
		ec_enc_icdf(psRangeEnc, int64(psIndices.NLSFInterpCoef_Q2), &silk_NLSF_interpolation_factor_iCDF[0], 8)
	}
	if int64(psIndices.SignalType) == TYPE_VOICED {
		encode_absolute_lagIndex = 1
		if condCoding == CODE_CONDITIONALLY && psEncC.Ec_prevSignalType == TYPE_VOICED {
			delta_lagIndex = int64(psIndices.LagIndex - psEncC.Ec_prevLagIndex)
			if delta_lagIndex < -8 || delta_lagIndex > 11 {
				delta_lagIndex = 0
			} else {
				delta_lagIndex = delta_lagIndex + 9
				encode_absolute_lagIndex = 0
			}
			ec_enc_icdf(psRangeEnc, delta_lagIndex, &silk_pitch_delta_iCDF[0], 8)
		}
		if encode_absolute_lagIndex != 0 {
			var (
				pitch_high_bits opus_int32
				pitch_low_bits  opus_int32
			)
			pitch_high_bits = opus_int32(int64(psIndices.LagIndex) / (psEncC.Fs_kHz >> 1))
			pitch_low_bits = opus_int32(psIndices.LagIndex) - opus_int32(opus_int16(pitch_high_bits))*opus_int32(opus_int16(psEncC.Fs_kHz>>1))
			ec_enc_icdf(psRangeEnc, int64(pitch_high_bits), &silk_pitch_lag_iCDF[0], 8)
			ec_enc_icdf(psRangeEnc, int64(pitch_low_bits), psEncC.Pitch_lag_low_bits_iCDF, 8)
		}
		psEncC.Ec_prevLagIndex = psIndices.LagIndex
		ec_enc_icdf(psRangeEnc, int64(psIndices.ContourIndex), psEncC.Pitch_contour_iCDF, 8)
		ec_enc_icdf(psRangeEnc, int64(psIndices.PERIndex), &silk_LTP_per_index_iCDF[0], 8)
		for k = 0; k < psEncC.Nb_subfr; k++ {
			ec_enc_icdf(psRangeEnc, int64(psIndices.LTPIndex[k]), silk_LTP_gain_iCDF_ptrs[psIndices.PERIndex], 8)
		}
		if condCoding == CODE_INDEPENDENTLY {
			ec_enc_icdf(psRangeEnc, int64(psIndices.LTP_scaleIndex), &silk_LTPscale_iCDF[0], 8)
		}
	}
	psEncC.Ec_prevSignalType = int64(psIndices.SignalType)
	ec_enc_icdf(psRangeEnc, int64(psIndices.Seed), &silk_uniform4_iCDF[0], 8)
}
