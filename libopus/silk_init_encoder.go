package libopus

import "github.com/gotranspile/opus/silk"

func silk_init_encoder(psEnc *silk_encoder_state_FLP, arch int) int {
	return silk.InitEncoder(psEnc, arch)
}
