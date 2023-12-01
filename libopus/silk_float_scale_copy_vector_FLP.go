package libopus

import "unsafe"

func silk_scale_copy_vector_FLP(data_out *float32, data_in *float32, gain float32, dataSize int) {
	var (
		i         int
		dataSize4 int
	)
	dataSize4 = dataSize & 0xFFFC
	for i = 0; i < dataSize4; i += 4 {
		*(*float32)(unsafe.Add(unsafe.Pointer(data_out), unsafe.Sizeof(float32(0))*uintptr(i+0))) = gain * *(*float32)(unsafe.Add(unsafe.Pointer(data_in), unsafe.Sizeof(float32(0))*uintptr(i+0)))
		*(*float32)(unsafe.Add(unsafe.Pointer(data_out), unsafe.Sizeof(float32(0))*uintptr(i+1))) = gain * *(*float32)(unsafe.Add(unsafe.Pointer(data_in), unsafe.Sizeof(float32(0))*uintptr(i+1)))
		*(*float32)(unsafe.Add(unsafe.Pointer(data_out), unsafe.Sizeof(float32(0))*uintptr(i+2))) = gain * *(*float32)(unsafe.Add(unsafe.Pointer(data_in), unsafe.Sizeof(float32(0))*uintptr(i+2)))
		*(*float32)(unsafe.Add(unsafe.Pointer(data_out), unsafe.Sizeof(float32(0))*uintptr(i+3))) = gain * *(*float32)(unsafe.Add(unsafe.Pointer(data_in), unsafe.Sizeof(float32(0))*uintptr(i+3)))
	}
	for ; i < dataSize; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(data_out), unsafe.Sizeof(float32(0))*uintptr(i))) = gain * *(*float32)(unsafe.Add(unsafe.Pointer(data_in), unsafe.Sizeof(float32(0))*uintptr(i)))
	}
}
