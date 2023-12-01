package libopus

import "unsafe"

func silk_scale_vector_FLP(data1 *float32, gain float32, dataSize int) {
	var (
		i         int
		dataSize4 int
	)
	dataSize4 = dataSize & 0xFFFC
	for i = 0; i < dataSize4; i += 4 {
		*(*float32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(float32(0))*uintptr(i+0))) *= gain
		*(*float32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(float32(0))*uintptr(i+1))) *= gain
		*(*float32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(float32(0))*uintptr(i+2))) *= gain
		*(*float32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(float32(0))*uintptr(i+3))) *= gain
	}
	for ; i < dataSize; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(float32(0))*uintptr(i))) *= gain
	}
}
