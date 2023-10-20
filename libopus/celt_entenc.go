package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const _entenc_H = 1

func ec_write_byte(_this *ec_enc, _value uint) int {
	if int(_this.Offs)+int(_this.End_offs) >= int(_this.Storage) {
		return -1
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), func() uint32 {
		p := &_this.Offs
		x := *p
		*p++
		return x
	}())) = uint8(_value)
	return 0
}
func ec_write_byte_at_end(_this *ec_enc, _value uint) int {
	if int(_this.Offs)+int(_this.End_offs) >= int(_this.Storage) {
		return -1
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), int(_this.Storage)-int(func() uint32 {
		p := &_this.End_offs
		*p++
		return *p
	}()))) = uint8(_value)
	return 0
}
func ec_enc_carry_out(_this *ec_enc, _c int) {
	if _c != ((1 << 8) - 1) {
		var carry int
		carry = _c >> 8
		if _this.Rem >= 0 {
			_this.Error |= ec_write_byte(_this, uint(_this.Rem+carry))
		}
		if int(_this.Ext) > 0 {
			var sym uint
			sym = uint((carry + ((1 << 8) - 1)) & ((1 << 8) - 1))
			for {
				_this.Error |= ec_write_byte(_this, sym)
				if int(func() uint32 {
					p := &_this.Ext
					*p--
					return *p
				}()) <= 0 {
					break
				}
			}
		}
		_this.Rem = _c & ((1 << 8) - 1)
	} else {
		_this.Ext++
	}
}
func ec_enc_normalize(_this *ec_enc) {
	for int(_this.Rng) <= ((1 << (32 - 1)) >> 8) {
		ec_enc_carry_out(_this, int(_this.Val)>>(32-8-1))
		_this.Val = uint32(int32((int(_this.Val) << 8) & ((1 << (32 - 1)) - 1)))
		_this.Rng <<= 8
		_this.Nbits_total += 8
	}
}
func ec_enc_init(_this *ec_enc, _buf *uint8, _size uint32) {
	_this.Buf = _buf
	_this.End_offs = 0
	_this.End_window = 0
	_this.Nend_bits = 0
	_this.Nbits_total = 32 + 1
	_this.Offs = 0
	_this.Rng = 1 << (32 - 1)
	_this.Rem = -1
	_this.Val = 0
	_this.Ext = 0
	_this.Storage = _size
	_this.Error = 0
}
func ec_encode(_this *ec_enc, _fl uint, _fh uint, _ft uint) {
	var r uint32
	r = celt_udiv(_this.Rng, uint32(_ft))
	if _fl > 0 {
		_this.Val += uint32(uint(_this.Rng) - uint(r)*(_ft-_fl))
		_this.Rng = uint32(uint(r) * (_fh - _fl))
	} else {
		_this.Rng -= uint32(uint(r) * (_ft - _fh))
	}
	ec_enc_normalize(_this)
}
func ec_encode_bin(_this *ec_enc, _fl uint, _fh uint, _bits uint) {
	var r uint32
	r = uint32(uint(_this.Rng) >> _bits)
	if _fl > 0 {
		_this.Val += uint32(uint(_this.Rng) - uint(r)*((1<<_bits)-_fl))
		_this.Rng = uint32(uint(r) * (_fh - _fl))
	} else {
		_this.Rng -= uint32(uint(r) * ((1 << _bits) - _fh))
	}
	ec_enc_normalize(_this)
}
func ec_enc_bit_logp(_this *ec_enc, _val int, _logp uint) {
	var (
		r uint32
		s uint32
		l uint32
	)
	r = _this.Rng
	l = _this.Val
	s = uint32(uint(r) >> _logp)
	r -= s
	if _val != 0 {
		_this.Val = uint32(int32(int(l) + int(r)))
	}
	if _val != 0 {
		_this.Rng = s
	} else {
		_this.Rng = r
	}
	ec_enc_normalize(_this)
}
func ec_enc_icdf(_this *ec_enc, _s int, _icdf *uint8, _ftb uint) {
	var r uint32
	r = uint32(uint(_this.Rng) >> _ftb)
	if _s > 0 {
		_this.Val += uint32(int32(int(_this.Rng) - int(r)*int(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s-1)))))
		_this.Rng = uint32(int32(int(r) * (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s-1))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s))))))
	} else {
		_this.Rng -= uint32(int32(int(r) * int(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s)))))
	}
	ec_enc_normalize(_this)
}
func ec_enc_uint(_this *ec_enc, _fl uint32, _ft uint32) {
	var (
		ft  uint
		fl  uint
		ftb int
	)
	_ft--
	ftb = ec_ilog(_ft)
	if ftb > 8 {
		ftb -= 8
		ft = uint((int(_ft) >> ftb) + 1)
		fl = uint(int(_fl) >> ftb)
		ec_encode(_this, fl, fl+1, ft)
		ec_enc_bits(_this, uint32(int32(int(_fl)&((1<<ftb)-1))), uint(ftb))
	} else {
		ec_encode(_this, uint(_fl), uint(int(_fl)+1), uint(int(_ft)+1))
	}
}
func ec_enc_bits(_this *ec_enc, _fl uint32, _bits uint) {
	var (
		window ec_window
		used   int
	)
	window = _this.End_window
	used = _this.Nend_bits
	if used+int(_bits) > (CHAR_BIT * int(unsafe.Sizeof(ec_window(0)))) {
		for {
			_this.Error |= ec_write_byte_at_end(_this, uint(window)&((1<<8)-1))
			window >>= 8
			used -= 8
			if used < 8 {
				break
			}
		}
	}
	window |= ec_window(int32(int(ec_window(_fl)) << used))
	used += int(_bits)
	_this.End_window = window
	_this.Nend_bits = used
	_this.Nbits_total += int(_bits)
}
func ec_enc_patch_initial_bits(_this *ec_enc, _val uint, _nbits uint) {
	var (
		shift int
		mask  uint
	)
	shift = int(8 - _nbits)
	mask = ((1 << _nbits) - 1) << uint(shift)
	if int(_this.Offs) > 0 {
		*_this.Buf = uint8((uint(*_this.Buf) & ^mask) | _val<<uint(shift))
	} else if _this.Rem >= 0 {
		_this.Rem = (_this.Rem & int(^mask)) | int(_val<<uint(shift))
	} else if uint(_this.Rng) <= ((1 << (32 - 1)) >> _nbits) {
		_this.Val = uint32(int32((int(_this.Val) & int(uint32(int32(^(int(uint32(mask)) << (32 - 8 - 1)))))) | int(uint32(_val))<<(shift+(32-8-1))))
	} else {
		_this.Error = -1
	}
}
func ec_enc_shrink(_this *ec_enc, _size uint32) {
	libc.MemMove(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _size))), -int(_this.End_offs)))), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage))), -int(_this.End_offs)))), int(uintptr(_this.End_offs)*unsafe.Sizeof(uint8(0))+uintptr((int64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _size))), -int(_this.End_offs)))))-uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage))), -int(_this.End_offs)))))))*0)))
	_this.Storage = _size
}
func ec_enc_done(_this *ec_enc) {
	var (
		window ec_window
		used   int
		msk    uint32
		end    uint32
		l      int
	)
	l = 32 - ec_ilog(_this.Rng)
	msk = uint32(int32(((1 << (32 - 1)) - 1) >> l))
	end = uint32(int32((int(_this.Val) + int(msk)) & int(^msk)))
	if (int(end) | int(msk)) >= int(_this.Val)+int(_this.Rng) {
		l++
		msk >>= 1
		end = uint32(int32((int(_this.Val) + int(msk)) & int(^msk)))
	}
	for l > 0 {
		ec_enc_carry_out(_this, int(end)>>(32-8-1))
		end = uint32(int32((int(end) << 8) & ((1 << (32 - 1)) - 1)))
		l -= 8
	}
	if _this.Rem >= 0 || int(_this.Ext) > 0 {
		ec_enc_carry_out(_this, 0)
	}
	window = _this.End_window
	used = _this.Nend_bits
	for used >= 8 {
		_this.Error |= ec_write_byte_at_end(_this, uint(window)&((1<<8)-1))
		window >>= 8
		used -= 8
	}
	if _this.Error == 0 {
		libc.MemSet(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Offs))), 0, (int(_this.Storage)-int(_this.Offs)-int(_this.End_offs))*int(unsafe.Sizeof(uint8(0))))
		if used > 0 {
			if int(_this.End_offs) >= int(_this.Storage) {
				_this.Error = -1
			} else {
				l = -l
				if int(_this.Offs)+int(_this.End_offs) >= int(_this.Storage) && l < used {
					window &= ec_window(int32((1 << l) - 1))
					_this.Error = -1
				}
				*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), int(_this.Storage)-int(_this.End_offs)-1)) |= uint8(window)
			}
		}
	}
}
