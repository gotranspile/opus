package celt

import (
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

type ECDec struct {
	ECCtx
}

func (ec *ECDec) readByte() int {
	if int(ec.Offs) < int(ec.Storage) {
		v := ec.Buf[ec.Offs]
		ec.Offs++
		return int(v)
	}
	return 0
}
func (ec *ECDec) readByteFromEnd() int {
	if int(ec.End_offs) < int(ec.Storage) {
		ec.End_offs++
		return int(ec.Buf[int(ec.Storage)-int(ec.End_offs)])
	}
	return 0
}
func (ec *ECDec) normalize() {
	for int(ec.Rng) <= ((1 << (32 - 1)) >> 8) {
		ec.Nbits_total += 8
		ec.Rng <<= 8
		sym := ec.Rem
		ec.Rem = ec.readByte()
		sym = (sym<<8 | ec.Rem) >> (8 - ((32-2)%8 + 1))
		ec.Val = uint32(int32(((int(ec.Val) << 8) + (^sym & ((1 << 8) - 1))) & ((1 << (32 - 1)) - 1)))
	}
}
func (ec *ECDec) Init(buf []uint8) {
	ec.Buf = buf
	ec.Storage = uint32(len(buf))
	ec.End_offs = 0
	ec.End_window = 0
	ec.Nend_bits = 0
	ec.Nbits_total = 32 + 1 - ((32-((32-2)%8+1))/8)*8
	ec.Offs = 0
	ec.Rng = 1 << ((32-2)%8 + 1)
	ec.Rem = ec.readByte()
	ec.Val = uint32(int32(int(ec.Rng) - 1 - (ec.Rem >> (8 - ((32-2)%8 + 1)))))
	ec.Error = 0
	ec.normalize()
}
func (ec *ECDec) Decode(ft uint) uint {
	ec.Ext = ec.Rng / uint32(ft)
	s := uint(int(ec.Val) / int(ec.Ext))
	return ft - ((s + 1) + ((ft - (s + 1)) & uint(-int(libc.BoolToInt(ft < (s+1))))))
}
func (ec *ECDec) DecodeBin(bits uint) uint {
	ec.Ext = uint32(uint(ec.Rng) >> bits)
	s := uint(int(ec.Val) / int(ec.Ext))
	return (1 << bits) - ((s + 1) + (((1 << bits) - (s + 1)) & uint(-int(libc.BoolToInt((1<<bits) < (s+1))))))
}
func (ec *ECDec) DecUpdate(fl uint, fh uint, ft uint) {
	s := uint32(uint(ec.Ext) * (ft - fh))
	ec.Val -= s
	if fl > 0 {
		ec.Rng = uint32(uint(ec.Ext) * (fh - fl))
	} else {
		ec.Rng = uint32(int32(int(ec.Rng) - int(s)))
	}
	ec.normalize()
}
func (ec *ECDec) DecBitLogp(logp uint) int {
	r := ec.Rng
	d := ec.Val
	s := uint32(uint(r) >> logp)
	ret := int(libc.BoolToInt(int(d) < int(s)))
	if ret == 0 {
		ec.Val = uint32(int32(int(d) - int(s)))
	}
	if ret != 0 {
		ec.Rng = s
	} else {
		ec.Rng = uint32(int32(int(r) - int(s)))
	}
	ec.normalize()
	return ret
}
func (ec *ECDec) DecIcdf(icdf []byte, ftb uint) int {
	s := ec.Rng
	d := ec.Val
	r := uint32(uint(s) >> ftb)
	ret := -1
	var t uint32
	for {
		t = s
		s = r * uint32(icdf[func() int {
			p := &ret
			*p++
			return *p
		}()])
		if int(d) >= int(s) {
			break
		}
	}
	ec.Val = uint32(int32(int(d) - int(s)))
	ec.Rng = uint32(int32(int(t) - int(s)))
	ec.normalize()
	return ret
}
func (ec *ECDec) DecUint(fta uint32) uint32 {
	fta--
	ftb := EC_ilog(fta)
	if ftb > 8 {
		var t uint32
		ftb -= 8
		ft := uint(int(fta)>>ftb) + 1
		s := ec.Decode(ft)
		ec.DecUpdate(s, s+1, ft)
		t = uint32(int32(int(uint32(s))<<ftb | int(ec.DecBits(uint(ftb)))))
		if int(t) <= int(fta) {
			return t
		}
		ec.Error = 1
		return fta
	} else {
		fta++
		s := ec.Decode(uint(fta))
		ec.DecUpdate(s, s+1, uint(fta))
		return uint32(s)
	}
}
func (ec *ECDec) DecBits(bits uint) uint32 {
	window := ec.End_window
	available := ec.Nend_bits
	if uint(available) < bits {
		for {
			window |= ECWindow(int32(int(ECWindow(int32(ec.readByteFromEnd()))) << available))
			available += 8
			if available > (8*int(unsafe.Sizeof(ECWindow(0))))-8 {
				break
			}
		}
	}
	ret := uint32(uint(window) & ((1 << bits) - 1))
	window >>= ECWindow(bits)
	available -= int(bits)
	ec.End_window = window
	ec.Nend_bits = available
	ec.Nbits_total += int(bits)
	return ret
}
