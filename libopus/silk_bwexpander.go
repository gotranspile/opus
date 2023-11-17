package libopus

import "github.com/gotranspile/opus/silk"

func silk_bwexpander(ar []int16, d int, chirp_Q16 int32) {
	silk.BwExpander(ar, d, chirp_Q16)
}
