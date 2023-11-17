package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_shell_encoder(psRangeEnc *ec_enc, pulses0 []int) {
	silk.ShellEncoder(psRangeEnc, pulses0)
}
func silk_shell_decoder(pulses0 []int16, psRangeDec *ec_dec, pulses4 int) {
	silk.ShellDecoder(pulses0, psRangeDec, pulses4)
}
