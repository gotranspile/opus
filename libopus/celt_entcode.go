package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
)

const _entcode_H = 1
const EC_UINT_BITS = 8
const BITRES = 3

type ec_window uint32
type ec_ctx struct {
	Buf         *uint8
	Storage     uint32
	End_offs    uint32
	End_window  ec_window
	Nend_bits   int
	Nbits_total int
	Offs        uint32
	Rng         uint32
	Val         uint32
	Ext         uint32
	Rem         int
	Error       int
}
type ec_enc ec_ctx
type ec_dec ec_ctx

func ec_range_bytes(_this *ec_ctx) uint32 {
	return _this.Offs
}
func ec_get_buffer(_this *ec_ctx) *uint8 {
	return _this.Buf
}
func ec_get_error(_this *ec_ctx) int {
	return _this.Error
}
func ec_tell(_this *ec_ctx) int {
	return _this.Nbits_total - ec_ilog(_this.Rng)
}
func celt_udiv(n uint32, d uint32) uint32 {
	return uint32(int32(int(n) / int(d)))
}
func celt_sudiv(n int32, d int32) int32 {
	return int32(int(n) / int(d))
}
func ec_ilog(_v uint32) int {
	var (
		ret int
		m   int
	)
	ret = int(libc.BoolToInt(int(_v) != 0))
	m = int(libc.BoolToInt((int(_v)&0xFFFF0000) != 0)) << 4
	_v >>= uint32(int32(m))
	ret |= m
	m = int(libc.BoolToInt((int(_v)&0xFF00) != 0)) << 3
	_v >>= uint32(int32(m))
	ret |= m
	m = int(libc.BoolToInt((int(_v)&0xF0) != 0)) << 2
	_v >>= uint32(int32(m))
	ret |= m
	m = int(libc.BoolToInt((int(_v)&0xC) != 0)) << 1
	_v >>= uint32(int32(m))
	ret |= m
	ret += int(libc.BoolToInt((int(_v) & 0x2) != 0))
	return ret
}
func ec_tell_frac(_this *ec_ctx) uint32 {
	var (
		correction [8]uint = [8]uint{35733, 38967, 42495, 46340, 50535, 55109, 60097, math.MaxUint16}
		nbits      uint32
		r          uint32
		l          int
		b          uint
	)
	nbits = uint32(int32(_this.Nbits_total << BITRES))
	l = ec_ilog(_this.Rng)
	r = uint32(int32(int(_this.Rng) >> (l - 16)))
	b = uint((int(r) >> 12) - 8)
	b += uint(libc.BoolToInt(uint(r) > correction[b]))
	l = (l << 3) + int(b)
	return uint32(int32(int(nbits) - l))
}
