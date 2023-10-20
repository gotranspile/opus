package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const _entdec_H = 1

func ec_read_byte(_this *ec_dec) int {
	if int(_this.Offs) < int(_this.Storage) {
		return int(*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), func() uint32 {
			p := &_this.Offs
			x := *p
			*p++
			return x
		}())))
	}
	return 0
}
func ec_read_byte_from_end(_this *ec_dec) int {
	if int(_this.End_offs) < int(_this.Storage) {
		return int(*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), int(_this.Storage)-int(func() uint32 {
			p := &_this.End_offs
			*p++
			return *p
		}()))))
	}
	return 0
}
func ec_dec_normalize(_this *ec_dec) {
	for int(_this.Rng) <= ((1 << (32 - 1)) >> 8) {
		var sym int
		_this.Nbits_total += 8
		_this.Rng <<= 8
		sym = _this.Rem
		_this.Rem = ec_read_byte(_this)
		sym = (sym<<8 | _this.Rem) >> (8 - ((32-2)%8 + 1))
		_this.Val = uint32(int32(((int(_this.Val) << 8) + (^sym & ((1 << 8) - 1))) & ((1 << (32 - 1)) - 1)))
	}
}
func ec_dec_init(_this *ec_dec, _buf *uint8, _storage uint32) {
	_this.Buf = _buf
	_this.Storage = _storage
	_this.End_offs = 0
	_this.End_window = 0
	_this.Nend_bits = 0
	_this.Nbits_total = 32 + 1 - ((32-((32-2)%8+1))/8)*8
	_this.Offs = 0
	_this.Rng = 1 << ((32-2)%8 + 1)
	_this.Rem = ec_read_byte(_this)
	_this.Val = uint32(int32(int(_this.Rng) - 1 - (_this.Rem >> (8 - ((32-2)%8 + 1)))))
	_this.Error = 0
	ec_dec_normalize(_this)
}
func ec_decode(_this *ec_dec, _ft uint) uint {
	var s uint
	_this.Ext = celt_udiv(_this.Rng, uint32(_ft))
	s = uint(int(_this.Val) / int(_this.Ext))
	return _ft - ((s + 1) + ((_ft - (s + 1)) & uint(-int(libc.BoolToInt(_ft < (s+1))))))
}
func ec_decode_bin(_this *ec_dec, _bits uint) uint {
	var s uint
	_this.Ext = uint32(uint(_this.Rng) >> _bits)
	s = uint(int(_this.Val) / int(_this.Ext))
	return (1 << _bits) - ((s + 1) + (((1 << _bits) - (s + 1)) & uint(-int(libc.BoolToInt((1<<_bits) < (s+1))))))
}
func ec_dec_update(_this *ec_dec, _fl uint, _fh uint, _ft uint) {
	var s uint32
	s = uint32(uint(_this.Ext) * (_ft - _fh))
	_this.Val -= s
	if _fl > 0 {
		_this.Rng = uint32(uint(_this.Ext) * (_fh - _fl))
	} else {
		_this.Rng = uint32(int32(int(_this.Rng) - int(s)))
	}
	ec_dec_normalize(_this)
}
func ec_dec_bit_logp(_this *ec_dec, _logp uint) int {
	var (
		r   uint32
		d   uint32
		s   uint32
		ret int
	)
	r = _this.Rng
	d = _this.Val
	s = uint32(uint(r) >> _logp)
	ret = int(libc.BoolToInt(int(d) < int(s)))
	if ret == 0 {
		_this.Val = uint32(int32(int(d) - int(s)))
	}
	if ret != 0 {
		_this.Rng = s
	} else {
		_this.Rng = uint32(int32(int(r) - int(s)))
	}
	ec_dec_normalize(_this)
	return ret
}
func ec_dec_icdf(_this *ec_dec, _icdf *uint8, _ftb uint) int {
	var (
		r   uint32
		d   uint32
		s   uint32
		t   uint32
		ret int
	)
	s = _this.Rng
	d = _this.Val
	r = uint32(uint(s) >> _ftb)
	ret = -1
	for {
		t = s
		s = uint32(int32(int(r) * int(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), func() int {
			p := &ret
			*p++
			return *p
		}())))))
		if int(d) >= int(s) {
			break
		}
	}
	_this.Val = uint32(int32(int(d) - int(s)))
	_this.Rng = uint32(int32(int(t) - int(s)))
	ec_dec_normalize(_this)
	return ret
}
func ec_dec_uint(_this *ec_dec, _ft uint32) uint32 {
	var (
		ft  uint
		s   uint
		ftb int
	)
	_ft--
	ftb = ec_ilog(_ft)
	if ftb > 8 {
		var t uint32
		ftb -= 8
		ft = uint(int(_ft)>>ftb) + 1
		s = ec_decode(_this, ft)
		ec_dec_update(_this, s, s+1, ft)
		t = uint32(int32(int(uint32(s))<<ftb | int(ec_dec_bits(_this, uint(ftb)))))
		if int(t) <= int(_ft) {
			return t
		}
		_this.Error = 1
		return _ft
	} else {
		_ft++
		s = ec_decode(_this, uint(_ft))
		ec_dec_update(_this, s, s+1, uint(_ft))
		return uint32(s)
	}
}
func ec_dec_bits(_this *ec_dec, _bits uint) uint32 {
	var (
		window    ec_window
		available int
		ret       uint32
	)
	window = _this.End_window
	available = _this.Nend_bits
	if uint(available) < _bits {
		for {
			window |= ec_window(int32(int(ec_window(int32(ec_read_byte_from_end(_this)))) << available))
			available += 8
			if available > (CHAR_BIT*int(unsafe.Sizeof(ec_window(0))))-8 {
				break
			}
		}
	}
	ret = uint32(uint(uint32(window)) & ((1 << _bits) - 1))
	window >>= ec_window(_bits)
	available -= int(_bits)
	_this.End_window = window
	_this.Nend_bits = available
	_this.Nbits_total += int(_bits)
	return ret
}
