package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_shell_encoder(psRangeEnc *ec_enc, pulses0 []int) {
	silk.ShellEncoder(psRangeEnc, pulses0)
}
