package libopus

import "unsafe"

func silk_autocorrelation_FLP(results *float32, inputData *float32, inputDataSize int, correlationCount int) {
	var i int
	if correlationCount > inputDataSize {
		correlationCount = inputDataSize
	}
	for i = 0; i < correlationCount; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(results), unsafe.Sizeof(float32(0))*uintptr(i))) = float32(silk_inner_product_FLP(inputData, (*float32)(unsafe.Add(unsafe.Pointer(inputData), unsafe.Sizeof(float32(0))*uintptr(i))), inputDataSize-i))
	}
}
