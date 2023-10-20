package celt

import "math/bits"

const EC_UINT_BITS = 8
const BITRES = 3

type ECWindow uint32
type ECCtx struct {
	Buf         []byte
	Storage     uint32
	End_offs    uint32
	End_window  ECWindow
	Nend_bits   int
	Nbits_total int
	Offs        uint32
	Rng         uint32
	Val         uint32
	Ext         uint32
	Rem         int
	Error       int
}
type ECEnc = ECCtx
type ECDec = ECCtx

func EC_ilog(v uint32) int {
	return 32 - bits.LeadingZeros32(v) // TODO: check
}
func (ec *ECCtx) RangeBytes() uint32 {
	return ec.Offs
}
func (ec *ECCtx) GetBuffer() []byte {
	return ec.Buf
}
func (ec *ECCtx) GetError() int {
	return ec.Error
}
func (ec *ECCtx) Tell() int {
	return ec.Nbits_total - EC_ilog(ec.Rng)
}
func (ec *ECCtx) TellFrac() uint32 {
	var correction = [8]uint{35733, 38967, 42495, 46340, 50535, 55109, 60097, 65535}
	nbits := ec.Nbits_total << BITRES
	l := EC_ilog(ec.Rng)
	r := uint32(int32(int(ec.Rng) >> (l - 16)))
	b := uint((int(r) >> 12) - 8)
	if uint(r) > correction[b] {
		b++
	}
	l = (l << 3) + int(b)
	return uint32(int32(nbits - l))
}
