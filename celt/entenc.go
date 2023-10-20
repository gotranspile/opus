package celt

import (
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

func (ec *ECEnc) WriteByte(_value uint) int {
	if int(ec.Offs)+int(ec.End_offs) >= int(ec.Storage) {
		return -1
	}
	ec.Buf[func() uint32 {
		p := &ec.Offs
		x := *p
		*p++
		return x
	}()] = byte(uint8(_value))
	return 0
}
func (ec *ECEnc) WriteByteAtEnd(_value uint) int {
	if int(ec.Offs)+int(ec.End_offs) >= int(ec.Storage) {
		return -1
	}
	ec.Buf[int(ec.Storage)-int(func() uint32 {
		p := &ec.End_offs
		*p++
		return *p
	}())] = byte(uint8(_value))
	return 0
}
func (ec *ECEnc) CarryOut(_c int) {
	if _c != ((1 << 8) - 1) {
		var carry int
		carry = _c >> 8
		if ec.Rem >= 0 {
			ec.Error |= ec.WriteByte(uint(ec.Rem + carry))
		}
		if int(ec.Ext) > 0 {
			var sym uint
			sym = uint((carry + ((1 << 8) - 1)) & ((1 << 8) - 1))
			for {
				ec.Error |= ec.WriteByte(sym)
				if int(func() uint32 {
					p := &ec.Ext
					*p--
					return *p
				}()) <= 0 {
					break
				}
			}
		}
		ec.Rem = _c & ((1 << 8) - 1)
	} else {
		ec.Ext++
	}
}
func (ec *ECEnc) Normalize() {
	for int(ec.Rng) <= ((1 << (32 - 1)) >> 8) {
		ec.CarryOut(int(ec.Val) >> (32 - 8 - 1))
		ec.Val = uint32(int32((int(ec.Val) << 8) & ((1 << (32 - 1)) - 1)))
		ec.Rng <<= 8
		ec.Nbits_total += 8
	}
}
func (ec *ECEnc) Init(buf []byte) {
	ec.Buf = buf
	ec.End_offs = 0
	ec.End_window = 0
	ec.Nend_bits = 0
	ec.Nbits_total = 32 + 1
	ec.Offs = 0
	ec.Rng = 1 << (32 - 1)
	ec.Rem = -1
	ec.Val = 0
	ec.Ext = 0
	ec.Storage = uint32(len(buf))
	ec.Error = 0
}
func (ec *ECEnc) Encode(_fl uint, _fh uint, _ft uint) {
	var r uint32
	r = ec.Rng / uint32(_ft)
	if _fl > 0 {
		ec.Val += uint32(uint(ec.Rng) - uint(r)*(_ft-_fl))
		ec.Rng = uint32(uint(r) * (_fh - _fl))
	} else {
		ec.Rng -= uint32(uint(r) * (_ft - _fh))
	}
	ec.Normalize()
}
func (ec *ECEnc) EncodeBin(_fl uint, _fh uint, _bits uint) {
	var r uint32
	r = uint32(uint(ec.Rng) >> _bits)
	if _fl > 0 {
		ec.Val += uint32(uint(ec.Rng) - uint(r)*((1<<_bits)-_fl))
		ec.Rng = uint32(uint(r) * (_fh - _fl))
	} else {
		ec.Rng -= uint32(uint(r) * ((1 << _bits) - _fh))
	}
	ec.Normalize()
}
func (ec *ECEnc) EncBitLogp(_val int, _logp uint) {
	var (
		r uint32
		s uint32
		l uint32
	)
	r = ec.Rng
	l = ec.Val
	s = uint32(uint(r) >> _logp)
	r -= s
	if _val != 0 {
		ec.Val = uint32(int32(int(l) + int(r)))
	}
	if _val != 0 {
		ec.Rng = s
	} else {
		ec.Rng = r
	}
	ec.Normalize()
}
func (ec *ECEnc) EncIcdf(_s int, _icdf []byte, _ftb uint) {
	var r uint32
	r = uint32(uint(ec.Rng) >> _ftb)
	if _s > 0 {
		ec.Val += uint32(int32(int(ec.Rng) - int(r*uint32(_icdf[_s-1]))))
		ec.Rng = r * uint32(_icdf[_s-1]-_icdf[_s])
	} else {
		ec.Rng -= r * uint32(_icdf[_s])
	}
	ec.Normalize()
}
func (ec *ECEnc) EncUint(_fl uint32, _ft uint32) {
	var (
		ft  uint
		fl  uint
		ftb int
	)
	_ft--
	ftb = EC_ilog(_ft)
	if ftb > 8 {
		ftb -= 8
		ft = uint((int(_ft) >> ftb) + 1)
		fl = uint(int(_fl) >> ftb)
		ec.Encode(fl, fl+1, ft)
		ec.EncBits(uint32(int32(int(_fl)&((1<<ftb)-1))), uint(ftb))
	} else {
		ec.Encode(uint(_fl), uint(int(_fl)+1), uint(int(_ft)+1))
	}
}
func (ec *ECEnc) EncBits(_fl uint32, _bits uint) {
	var (
		window ECWindow
		used   int
	)
	window = ec.End_window
	used = ec.Nend_bits
	if used+int(_bits) > (8 * int(unsafe.Sizeof(ECWindow(0)))) {
		for {
			ec.Error |= ec.WriteByteAtEnd(uint(window) & ((1 << 8) - 1))
			window >>= 8
			used -= 8
			if used < 8 {
				break
			}
		}
	}
	window |= ECWindow(int32(int(ECWindow(_fl)) << used))
	used += int(_bits)
	ec.End_window = window
	ec.Nend_bits = used
	ec.Nbits_total += int(_bits)
}
func (ec *ECEnc) EncPatchInitialBits(_val uint, _nbits uint) {
	var (
		shift int
		mask  uint
	)
	shift = int(8 - _nbits)
	mask = ((1 << _nbits) - 1) << uint(shift)
	if int(ec.Offs) > 0 {
		ec.Buf[0] = byte(uint8((uint(ec.Buf[0]) & ^mask) | _val<<uint(shift)))
	} else if ec.Rem >= 0 {
		ec.Rem = (ec.Rem & int(^mask)) | int(_val<<uint(shift))
	} else if uint(ec.Rng) <= ((1 << (32 - 1)) >> _nbits) {
		ec.Val = uint32(int32((int(ec.Val) & int(uint32(int32(^(int(uint32(mask)) << (32 - 8 - 1)))))) | int(uint32(_val))<<(shift+(32-8-1))))
	} else {
		ec.Error = -1
	}
}
func (ec *ECEnc) Shrink(_size uint32) {
	libc.MemMove(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&ec.Buf[_size]), -int(ec.End_offs)))), unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&ec.Buf[ec.Storage]), -int(ec.End_offs)))), int(uintptr(ec.End_offs)*unsafe.Sizeof(byte(0))+uintptr((int64(uintptr(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&ec.Buf[_size]), -int(ec.End_offs)))))-uintptr(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(&ec.Buf[ec.Storage]), -int(ec.End_offs)))))))*0)))
	ec.Storage = _size
}
func (ec *ECEnc) Done() {
	var (
		window ECWindow
		used   int
		msk    uint32
		end    uint32
		l      int
	)
	l = 32 - EC_ilog(ec.Rng)
	msk = uint32(int32(((1 << (32 - 1)) - 1) >> l))
	end = uint32(int32((int(ec.Val) + int(msk)) & int(^msk)))
	if (int(end) | int(msk)) >= int(ec.Val)+int(ec.Rng) {
		l++
		msk >>= 1
		end = uint32(int32((int(ec.Val) + int(msk)) & int(^msk)))
	}
	for l > 0 {
		ec.CarryOut(int(end) >> (32 - 8 - 1))
		end = uint32(int32((int(end) << 8) & ((1 << (32 - 1)) - 1)))
		l -= 8
	}
	if ec.Rem >= 0 || int(ec.Ext) > 0 {
		ec.CarryOut(0)
	}
	window = ec.End_window
	used = ec.Nend_bits
	for used >= 8 {
		ec.Error |= ec.WriteByteAtEnd(uint(window) & ((1 << 8) - 1))
		window >>= 8
		used -= 8
	}
	if ec.Error == 0 {
		libc.MemSet(unsafe.Pointer(&ec.Buf[ec.Offs]), 0, (int(ec.Storage)-int(ec.Offs)-int(ec.End_offs))*int(unsafe.Sizeof(byte(0))))
		if used > 0 {
			if int(ec.End_offs) >= int(ec.Storage) {
				ec.Error = -1
			} else {
				l = -l
				if int(ec.Offs)+int(ec.End_offs) >= int(ec.Storage) && l < used {
					window &= ECWindow(int32((1 << l) - 1))
					ec.Error = -1
				}
				ec.Buf[int(ec.Storage)-int(ec.End_offs)-1] |= byte(uint8(window))
			}
		}
	}
}
