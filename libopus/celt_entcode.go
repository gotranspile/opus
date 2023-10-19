package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
)

const _entcode_H = 1
const EC_UINT_BITS = 8
const BITRES = 3

type ec_window opus_uint32
type ec_ctx struct {
	Buf         *uint8
	Storage     opus_uint32
	End_offs    opus_uint32
	End_window  ec_window
	Nend_bits   int64
	Nbits_total int64
	Offs        opus_uint32
	Rng         opus_uint32
	Val         opus_uint32
	Ext         opus_uint32
	Rem         int64
	Error       int64
}
type ec_enc ec_ctx
type ec_dec ec_ctx

func ec_range_bytes(_this *ec_ctx) opus_uint32 {
	return _this.Offs
}
func ec_get_buffer(_this *ec_ctx) *uint8 {
	return _this.Buf
}
func ec_get_error(_this *ec_ctx) int64 {
	return _this.Error
}
func ec_tell(_this *ec_ctx) int64 {
	return _this.Nbits_total - ec_ilog(_this.Rng)
}
func celt_udiv(n opus_uint32, d opus_uint32) opus_uint32 {
	return n / d
}
func celt_sudiv(n opus_int32, d opus_int32) opus_int32 {
	return n / d
}
func ec_ilog(_v opus_uint32) int64 {
	var (
		ret int64
		m   int64
	)
	ret = int64(libc.BoolToInt(_v != 0))
	m = int64(libc.BoolToInt((_v&0xFFFF0000) != 0)) << 4
	_v >>= opus_uint32(m)
	ret |= m
	m = int64(libc.BoolToInt((_v&0xFF00) != 0)) << 3
	_v >>= opus_uint32(m)
	ret |= m
	m = int64(libc.BoolToInt((_v&0xF0) != 0)) << 2
	_v >>= opus_uint32(m)
	ret |= m
	m = int64(libc.BoolToInt((_v&0xC) != 0)) << 1
	_v >>= opus_uint32(m)
	ret |= m
	ret += int64(libc.BoolToInt((_v & 0x2) != 0))
	return ret
}
func ec_tell_frac(_this *ec_ctx) opus_uint32 {
	var (
		correction [8]uint64 = [8]uint64{35733, 38967, 42495, 46340, 50535, 55109, 60097, math.MaxUint16}
		nbits      opus_uint32
		r          opus_uint32
		l          int64
		b          uint64
	)
	nbits = opus_uint32(_this.Nbits_total << BITRES)
	l = ec_ilog(_this.Rng)
	r = _this.Rng >> opus_uint32(l-16)
	b = uint64((r >> 12) - 8)
	b += uint64(libc.BoolToInt(r > opus_uint32(correction[b])))
	l = int64(uint64(l<<3) + b)
	return nbits - opus_uint32(l)
}
