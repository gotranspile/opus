package libopus

import "unsafe"

func silk_bwexpander_FLP(ar *float32, d int, chirp float32) {
	var (
		i    int
		cfac float32 = chirp
	)
	for i = 0; i < d-1; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(float32(0))*uintptr(i))) *= cfac
		cfac *= chirp
	}
	*(*float32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(float32(0))*uintptr(d-1))) *= cfac
}
