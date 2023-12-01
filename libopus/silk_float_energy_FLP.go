package libopus

import "unsafe"

func silk_energy_FLP(data *float32, dataSize int) float64 {
	var (
		i      int
		result float64
	)
	result = 0.0
	for i = 0; i < dataSize-3; i += 4 {
		result += float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+0))))*float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+0)))) + float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+1))))*float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+1)))) + float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+2))))*float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+2)))) + float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+3))))*float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i+3))))
	}
	for ; i < dataSize; i++ {
		result += float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i)))) * float64(*(*float32)(unsafe.Add(unsafe.Pointer(data), unsafe.Sizeof(float32(0))*uintptr(i))))
	}
	return result
}
