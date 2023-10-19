package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const _entenc_H = 1

func ec_write_byte(_this *ec_enc, _value uint64) int64 {
	if _this.Offs+_this.End_offs >= _this.Storage {
		return -1
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), func() opus_uint32 {
		p := &_this.Offs
		x := *p
		*p++
		return x
	}())) = uint8(_value)
	return 0
}
func ec_write_byte_at_end(_this *ec_enc, _value uint64) int64 {
	if _this.Offs+_this.End_offs >= _this.Storage {
		return -1
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage-func() opus_uint32 {
		p := &_this.End_offs
		*p++
		return *p
	}())) = uint8(_value)
	return 0
}
func ec_enc_carry_out(_this *ec_enc, _c int64) {
	if _c != ((1 << 8) - 1) {
		var carry int64
		carry = _c >> 8
		if _this.Rem >= 0 {
			_this.Error |= ec_write_byte(_this, uint64(_this.Rem+carry))
		}
		if _this.Ext > 0 {
			var sym uint64
			sym = uint64((carry + ((1 << 8) - 1)) & ((1 << 8) - 1))
			for {
				_this.Error |= ec_write_byte(_this, sym)
				if func() opus_uint32 {
					p := &_this.Ext
					*p--
					return *p
				}() <= 0 {
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
	for _this.Rng <= ((1 << (32 - 1)) >> 8) {
		ec_enc_carry_out(_this, int64(_this.Val>>(32-8-1)))
		_this.Val = (_this.Val << 8) & ((1 << (32 - 1)) - 1)
		_this.Rng <<= 8
		_this.Nbits_total += 8
	}
}
func ec_enc_init(_this *ec_enc, _buf *uint8, _size opus_uint32) {
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
func ec_encode(_this *ec_enc, _fl uint64, _fh uint64, _ft uint64) {
	var r opus_uint32
	r = celt_udiv(_this.Rng, opus_uint32(_ft))
	if _fl > 0 {
		_this.Val += _this.Rng - r*opus_uint32(_ft-_fl)
		_this.Rng = r * opus_uint32(_fh-_fl)
	} else {
		_this.Rng -= r * opus_uint32(_ft-_fh)
	}
	ec_enc_normalize(_this)
}
func ec_encode_bin(_this *ec_enc, _fl uint64, _fh uint64, _bits uint64) {
	var r opus_uint32
	r = _this.Rng >> opus_uint32(_bits)
	if _fl > 0 {
		_this.Val += _this.Rng - r*opus_uint32((1<<_bits)-_fl)
		_this.Rng = r * opus_uint32(_fh-_fl)
	} else {
		_this.Rng -= r * opus_uint32((1<<_bits)-_fh)
	}
	ec_enc_normalize(_this)
}
func ec_enc_bit_logp(_this *ec_enc, _val int64, _logp uint64) {
	var (
		r opus_uint32
		s opus_uint32
		l opus_uint32
	)
	r = _this.Rng
	l = _this.Val
	s = r >> opus_uint32(_logp)
	r -= s
	if _val != 0 {
		_this.Val = l + r
	}
	if _val != 0 {
		_this.Rng = s
	} else {
		_this.Rng = r
	}
	ec_enc_normalize(_this)
}
func ec_enc_icdf(_this *ec_enc, _s int64, _icdf *uint8, _ftb uint64) {
	var r opus_uint32
	r = _this.Rng >> opus_uint32(_ftb)
	if _s > 0 {
		_this.Val += _this.Rng - r*opus_uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s-1)))
		_this.Rng = r * opus_uint32(int64(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s-1)))-int64(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s))))
	} else {
		_this.Rng -= r * opus_uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), _s)))
	}
	ec_enc_normalize(_this)
}
func ec_enc_uint(_this *ec_enc, _fl opus_uint32, _ft opus_uint32) {
	var (
		ft  uint64
		fl  uint64
		ftb int64
	)
	_ft--
	ftb = ec_ilog(_ft)
	if ftb > 8 {
		ftb -= 8
		ft = uint64((_ft >> opus_uint32(ftb)) + 1)
		fl = uint64(_fl >> opus_uint32(ftb))
		ec_encode(_this, fl, fl+1, ft)
		ec_enc_bits(_this, _fl&opus_uint32((1<<ftb)-1), uint64(ftb))
	} else {
		ec_encode(_this, uint64(_fl), uint64(_fl+1), uint64(_ft+1))
	}
}
func ec_enc_bits(_this *ec_enc, _fl opus_uint32, _bits uint64) {
	var (
		window ec_window
		used   int64
	)
	window = _this.End_window
	used = _this.Nend_bits
	if uint64(used)+_bits > uint64(CHAR_BIT*int64(unsafe.Sizeof(ec_window(0)))) {
		for {
			_this.Error |= ec_write_byte_at_end(_this, uint64(window)&((1<<8)-1))
			window >>= 8
			used -= 8
			if used < 8 {
				break
			}
		}
	}
	window |= ec_window(_fl) << ec_window(used)
	used += int64(_bits)
	_this.End_window = window
	_this.Nend_bits = used
	_this.Nbits_total += int64(_bits)
}
func ec_enc_patch_initial_bits(_this *ec_enc, _val uint64, _nbits uint64) {
	var (
		shift int64
		mask  uint64
	)
	shift = int64(8 - _nbits)
	mask = ((1 << _nbits) - 1) << uint64(shift)
	if _this.Offs > 0 {
		*_this.Buf = uint8((uint64(*_this.Buf) & ^mask) | _val<<uint64(shift))
	} else if _this.Rem >= 0 {
		_this.Rem = int64((uint64(_this.Rem) & ^mask) | _val<<uint64(shift))
	} else if _this.Rng <= opus_uint32((1<<(32-1))>>_nbits) {
		_this.Val = (_this.Val & ^(opus_uint32(mask) << (32 - 8 - 1))) | opus_uint32(_val)<<opus_uint32(shift+(32-8-1))
	} else {
		_this.Error = -1
	}
}
func ec_enc_shrink(_this *ec_enc, _size opus_uint32) {
	libc.MemMove(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _size))), -int(_this.End_offs)))), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage))), -int(_this.End_offs)))), int(_this.End_offs*opus_uint32(unsafe.Sizeof(uint8(0)))+opus_uint32((int64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _size))), -int(_this.End_offs)))))-uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage))), -int(_this.End_offs)))))))*0)))
	_this.Storage = _size
}
func ec_enc_done(_this *ec_enc) {
	var (
		window ec_window
		used   int64
		msk    opus_uint32
		end    opus_uint32
		l      int64
	)
	l = 32 - ec_ilog(_this.Rng)
	msk = opus_uint32(((1 << (32 - 1)) - 1) >> l)
	end = (_this.Val + msk) & ^msk
	if (end | msk) >= _this.Val+_this.Rng {
		l++
		msk >>= 1
		end = (_this.Val + msk) & ^msk
	}
	for l > 0 {
		ec_enc_carry_out(_this, int64(end>>(32-8-1)))
		end = (end << 8) & ((1 << (32 - 1)) - 1)
		l -= 8
	}
	if _this.Rem >= 0 || _this.Ext > 0 {
		ec_enc_carry_out(_this, 0)
	}
	window = _this.End_window
	used = _this.Nend_bits
	for used >= 8 {
		_this.Error |= ec_write_byte_at_end(_this, uint64(window)&((1<<8)-1))
		window >>= 8
		used -= 8
	}
	if _this.Error == 0 {
		libc.MemSet(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Offs))), 0, int((_this.Storage-_this.Offs-_this.End_offs)*opus_uint32(unsafe.Sizeof(uint8(0)))))
		if used > 0 {
			if _this.End_offs >= _this.Storage {
				_this.Error = -1
			} else {
				l = -l
				if _this.Offs+_this.End_offs >= _this.Storage && l < used {
					window &= ec_window((1 << l) - 1)
					_this.Error = -1
				}
				*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage-_this.End_offs-1)) |= uint8(window)
			}
		}
	}
}
