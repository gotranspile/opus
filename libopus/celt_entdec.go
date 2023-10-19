package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const _entdec_H = 1

func ec_read_byte(_this *ec_dec) int64 {
	if _this.Offs < _this.Storage {
		return int64(*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), func() opus_uint32 {
			p := &_this.Offs
			x := *p
			*p++
			return x
		}())))
	}
	return 0
}
func ec_read_byte_from_end(_this *ec_dec) int64 {
	if _this.End_offs < _this.Storage {
		return int64(*(*uint8)(unsafe.Add(unsafe.Pointer(_this.Buf), _this.Storage-func() opus_uint32 {
			p := &_this.End_offs
			*p++
			return *p
		}())))
	}
	return 0
}
func ec_dec_normalize(_this *ec_dec) {
	for _this.Rng <= ((1 << (32 - 1)) >> 8) {
		var sym int64
		_this.Nbits_total += 8
		_this.Rng <<= 8
		sym = _this.Rem
		_this.Rem = ec_read_byte(_this)
		sym = (sym<<8 | _this.Rem) >> (8 - ((32-2)%8 + 1))
		_this.Val = ((_this.Val << 8) + opus_uint32(^sym&((1<<8)-1))) & ((1 << (32 - 1)) - 1)
	}
}
func ec_dec_init(_this *ec_dec, _buf *uint8, _storage opus_uint32) {
	_this.Buf = _buf
	_this.Storage = _storage
	_this.End_offs = 0
	_this.End_window = 0
	_this.Nend_bits = 0
	_this.Nbits_total = 32 + 1 - ((32-((32-2)%8+1))/8)*8
	_this.Offs = 0
	_this.Rng = 1 << ((32-2)%8 + 1)
	_this.Rem = ec_read_byte(_this)
	_this.Val = _this.Rng - 1 - opus_uint32(_this.Rem>>(8-((32-2)%8+1)))
	_this.Error = 0
	ec_dec_normalize(_this)
}
func ec_decode(_this *ec_dec, _ft uint64) uint64 {
	var s uint64
	_this.Ext = celt_udiv(_this.Rng, opus_uint32(_ft))
	s = uint64(_this.Val / _this.Ext)
	return _ft - ((s + 1) + ((_ft - (s + 1)) & uint64(-int64(libc.BoolToInt(_ft < (s+1))))))
}
func ec_decode_bin(_this *ec_dec, _bits uint64) uint64 {
	var s uint64
	_this.Ext = _this.Rng >> opus_uint32(_bits)
	s = uint64(_this.Val / _this.Ext)
	return (1 << _bits) - ((s + 1) + (((1 << _bits) - (s + 1)) & uint64(-int64(libc.BoolToInt((1<<_bits) < (s+1))))))
}
func ec_dec_update(_this *ec_dec, _fl uint64, _fh uint64, _ft uint64) {
	var s opus_uint32
	s = _this.Ext * opus_uint32(_ft-_fh)
	_this.Val -= s
	if _fl > 0 {
		_this.Rng = _this.Ext * opus_uint32(_fh-_fl)
	} else {
		_this.Rng = _this.Rng - s
	}
	ec_dec_normalize(_this)
}
func ec_dec_bit_logp(_this *ec_dec, _logp uint64) int64 {
	var (
		r   opus_uint32
		d   opus_uint32
		s   opus_uint32
		ret int64
	)
	r = _this.Rng
	d = _this.Val
	s = r >> opus_uint32(_logp)
	ret = int64(libc.BoolToInt(d < s))
	if ret == 0 {
		_this.Val = d - s
	}
	if ret != 0 {
		_this.Rng = s
	} else {
		_this.Rng = r - s
	}
	ec_dec_normalize(_this)
	return ret
}
func ec_dec_icdf(_this *ec_dec, _icdf *uint8, _ftb uint64) int64 {
	var (
		r   opus_uint32
		d   opus_uint32
		s   opus_uint32
		t   opus_uint32
		ret int64
	)
	s = _this.Rng
	d = _this.Val
	r = s >> opus_uint32(_ftb)
	ret = -1
	for {
		t = s
		s = r * opus_uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(_icdf), func() int64 {
			p := &ret
			*p++
			return *p
		}())))
		if d >= s {
			break
		}
	}
	_this.Val = d - s
	_this.Rng = t - s
	ec_dec_normalize(_this)
	return ret
}
func ec_dec_uint(_this *ec_dec, _ft opus_uint32) opus_uint32 {
	var (
		ft  uint64
		s   uint64
		ftb int64
	)
	_ft--
	ftb = ec_ilog(_ft)
	if ftb > 8 {
		var t opus_uint32
		ftb -= 8
		ft = uint64(_ft>>opus_uint32(ftb)) + 1
		s = ec_decode(_this, ft)
		ec_dec_update(_this, s, s+1, ft)
		t = opus_uint32(s)<<opus_uint32(ftb) | ec_dec_bits(_this, uint64(ftb))
		if t <= _ft {
			return t
		}
		_this.Error = 1
		return _ft
	} else {
		_ft++
		s = ec_decode(_this, uint64(_ft))
		ec_dec_update(_this, s, s+1, uint64(_ft))
		return opus_uint32(s)
	}
}
func ec_dec_bits(_this *ec_dec, _bits uint64) opus_uint32 {
	var (
		window    ec_window
		available int64
		ret       opus_uint32
	)
	window = _this.End_window
	available = _this.Nend_bits
	if uint64(available) < _bits {
		for {
			window |= ec_window(ec_read_byte_from_end(_this)) << ec_window(available)
			available += 8
			if available > (CHAR_BIT*int64(unsafe.Sizeof(ec_window(0))))-8 {
				break
			}
		}
	}
	ret = opus_uint32(window) & opus_uint32((1<<_bits)-1)
	window >>= ec_window(_bits)
	available -= int64(_bits)
	_this.End_window = window
	_this.Nend_bits = available
	_this.Nbits_total += int64(_bits)
	return ret
}
