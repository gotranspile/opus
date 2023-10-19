package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"os"
	"unsafe"
)

const MAX_PACKET = 1275

func main() {
	var (
		argc             int64    = int64(len(os.Args))
		argv             [0]*byte = ([0]*byte)(libc.CStringSlice(os.Args))
		err              int64
		inFile           *byte
		outFile          *byte
		fin              *stdio.File
		fout             *stdio.File
		mode             *OpusCustomMode = nil
		enc              *OpusCustomEncoder
		dec              *OpusCustomDecoder
		len_             int64
		frame_size       opus_int32
		channels         opus_int32
		rate             opus_int32
		bytes_per_packet int64
		data             [1275]uint8
		complexity       int64
		count            int64 = 0
		skip             opus_int32
		in               *opus_int16
		out              *opus_int16
	)
	if argc != 9 && argc != 8 && argc != 7 {
		stdio.Fprintf(stdio.Stderr(), "Usage: test_opus_custom <rate> <channels> <frame size>  <bytes per packet> [<complexity> [packet loss rate]] <input> <output>\n")
		os.Exit(1)
	}
	rate = opus_int32(libc.Atoi(libc.GoString(argv[1])))
	channels = opus_int32(libc.Atoi(libc.GoString(argv[2])))
	frame_size = opus_int32(libc.Atoi(libc.GoString(argv[3])))
	mode = opus_custom_mode_create(rate, int64(frame_size), nil)
	if mode == nil {
		stdio.Fprintf(stdio.Stderr(), "failed to create a mode\n")
		os.Exit(1)
	}
	bytes_per_packet = int64(libc.Atoi(libc.GoString(argv[4])))
	if bytes_per_packet < 0 || bytes_per_packet > MAX_PACKET {
		stdio.Fprintf(stdio.Stderr(), "bytes per packet must be between 0 and %d\n", MAX_PACKET)
		os.Exit(1)
	}
	inFile = argv[argc-2]
	fin = stdio.FOpen(libc.GoString(inFile), "rb")
	if fin == nil {
		stdio.Fprintf(stdio.Stderr(), "Could not open input file %s\n", argv[argc-2])
		os.Exit(1)
	}
	outFile = argv[argc-1]
	fout = stdio.FOpen(libc.GoString(outFile), "wb+")
	if fout == nil {
		stdio.Fprintf(stdio.Stderr(), "Could not open output file %s\n", argv[argc-1])
		fin.Close()
		os.Exit(1)
	}
	enc = opus_custom_encoder_create(mode, int64(channels), &err)
	if err != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to create the encoder: %s\n", opus_strerror(err))
		fin.Close()
		fout.Close()
		os.Exit(1)
	}
	dec = opus_custom_decoder_create(mode, int64(channels), &err)
	if err != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to create the decoder: %s\n", opus_strerror(err))
		fin.Close()
		fout.Close()
		os.Exit(1)
	}
	opus_custom_decoder_ctl(dec, OPUS_GET_LOOKAHEAD_REQUEST, (*opus_int32)(unsafe.Add(unsafe.Pointer(&skip), unsafe.Sizeof(opus_int32(0))*uintptr(int64(uintptr(unsafe.Pointer(&skip))-uintptr(unsafe.Pointer(&skip)))))))
	if argc > 7 {
		complexity = int64(libc.Atoi(libc.GoString(argv[5])))
		opus_custom_encoder_ctl(enc, OPUS_SET_COMPLEXITY_REQUEST, func() opus_int32 {
			complexity == 0
			return opus_int32(complexity)
		}())
	}
	in = (*opus_int16)(libc.Malloc(int(frame_size * channels * opus_int32(unsafe.Sizeof(opus_int16(0))))))
	out = (*opus_int16)(libc.Malloc(int(frame_size * channels * opus_int32(unsafe.Sizeof(opus_int16(0))))))
	for int64(fin.IsEOF()) == 0 {
		var ret int64
		err = int64(fin.ReadN((*byte)(unsafe.Pointer(in)), int(unsafe.Sizeof(int16(0))), int(frame_size*channels)))
		if int64(fin.IsEOF()) != 0 {
			break
		}
		len_ = opus_custom_encode(enc, in, int64(frame_size), &data[0], bytes_per_packet)
		if len_ <= 0 {
			stdio.Fprintf(stdio.Stderr(), "opus_custom_encode() failed: %s\n", opus_strerror(len_))
		}
		if argc == 9 && int64(libc.Rand())%1000 < int64(libc.Atoi(libc.GoString(argv[argc-3]))) {
			ret = opus_custom_decode(dec, nil, len_, out, int64(frame_size))
		} else {
			ret = opus_custom_decode(dec, &data[0], len_, out, int64(frame_size))
		}
		if ret < 0 {
			stdio.Fprintf(stdio.Stderr(), "opus_custom_decode() failed: %s\n", opus_strerror(ret))
		}
		count++
		fout.WriteN((*byte)(unsafe.Pointer((*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(skip*channels))))), int(unsafe.Sizeof(int16(0))), int((ret-int64(skip))*int64(channels)))
		skip = 0
	}
	opus_custom_encoder_destroy(enc)
	opus_custom_decoder_destroy(dec)
	fin.Close()
	fout.Close()
	opus_custom_mode_destroy(mode)
	libc.Free(unsafe.Pointer(in))
	libc.Free(unsafe.Pointer(out))
	os.Exit(0)
}
