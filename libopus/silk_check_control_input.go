package libopus

func check_control_input(encControl *silk_EncControlStruct) int64 {
	if encControl.API_sampleRate != 8000 && encControl.API_sampleRate != 12000 && encControl.API_sampleRate != 16000 && encControl.API_sampleRate != 24000 && encControl.API_sampleRate != 32000 && encControl.API_sampleRate != 44100 && encControl.API_sampleRate != 48000 || encControl.DesiredInternalSampleRate != 8000 && encControl.DesiredInternalSampleRate != 12000 && encControl.DesiredInternalSampleRate != 16000 || encControl.MaxInternalSampleRate != 8000 && encControl.MaxInternalSampleRate != 12000 && encControl.MaxInternalSampleRate != 16000 || encControl.MinInternalSampleRate != 8000 && encControl.MinInternalSampleRate != 12000 && encControl.MinInternalSampleRate != 16000 || encControl.MinInternalSampleRate > encControl.DesiredInternalSampleRate || encControl.MaxInternalSampleRate < encControl.DesiredInternalSampleRate || encControl.MinInternalSampleRate > encControl.MaxInternalSampleRate {
		return -102
	}
	if encControl.PayloadSize_ms != 10 && encControl.PayloadSize_ms != 20 && encControl.PayloadSize_ms != 40 && encControl.PayloadSize_ms != 60 {
		return -103
	}
	if encControl.PacketLossPercentage < 0 || encControl.PacketLossPercentage > 100 {
		return -105
	}
	if encControl.UseDTX < 0 || encControl.UseDTX > 1 {
		return -108
	}
	if encControl.UseCBR < 0 || encControl.UseCBR > 1 {
		return -109
	}
	if encControl.UseInBandFEC < 0 || encControl.UseInBandFEC > 1 {
		return -107
	}
	if encControl.NChannelsAPI < 1 || encControl.NChannelsAPI > ENCODER_NUM_CHANNELS {
		return -111
	}
	if encControl.NChannelsInternal < 1 || encControl.NChannelsInternal > ENCODER_NUM_CHANNELS {
		return -111
	}
	if encControl.NChannelsInternal > encControl.NChannelsAPI {
		return -111
	}
	if encControl.Complexity < 0 || encControl.Complexity > 10 {
		return -106
	}
	return SILK_NO_ERROR
}
