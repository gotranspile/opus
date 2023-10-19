package libopus

import (
	"math"
	"unsafe"
)

const WEIGHTS_SCALE = 0
const MAX_NEURONS = 32

type DenseLayer struct {
	Bias          *int8
	Input_weights *int8
	Nb_inputs     int64
	Nb_neurons    int64
	Sigmoid       int64
}
type GRULayer struct {
	Bias              *int8
	Input_weights     *int8
	Recurrent_weights *int8
	Nb_inputs         int64
	Nb_neurons        int64
}

func tansig_approx(x float32) float32 {
	var (
		i    int64
		y    float32
		dy   float32
		sign float32 = 1
	)
	if x >= 8 {
		return 1
	}
	if x <= float32(-8) {
		return float32(-1)
	}
	if x != x {
		return 0
	}
	if x < 0 {
		x = -x
		sign = float32(-1)
	}
	i = int64(math.Floor(float64(x*25) + 0.5))
	x -= float32(float64(i) * 0.04)
	y = tansig_table[i]
	dy = 1 - y*y
	y = y + x*dy*(1-y*x)
	return sign * y
}
func sigmoid_approx(x float32) float32 {
	return float32(float64(tansig_approx(float32(float64(x)*0.5)))*0.5 + 0.5)
}
func gemm_accum(out *float32, weights *int8, rows int64, cols int64, col_stride int64, x *float32) {
	var (
		i int64
		j int64
	)
	for i = 0; i < rows; i++ {
		for j = 0; j < cols; j++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(float32(0))*uintptr(i))) += float32(*(*int8)(unsafe.Add(unsafe.Pointer(weights), j*col_stride+i))) * *(*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(j)))
		}
	}
}
func compute_dense(layer *DenseLayer, output *float32, input *float32) {
	var (
		i      int64
		N      int64
		M      int64
		stride int64
	)
	M = layer.Nb_inputs
	N = layer.Nb_neurons
	stride = N
	for i = 0; i < N; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(float32(0))*uintptr(i))) = float32(*(*int8)(unsafe.Add(unsafe.Pointer(layer.Bias), i)))
	}
	gemm_accum(output, layer.Input_weights, N, M, stride, input)
	for i = 0; i < N; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(float32(0))*uintptr(i))) *= float32(1.0 / 128)
	}
	if layer.Sigmoid != 0 {
		for i = 0; i < N; i++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(float32(0))*uintptr(i))) = sigmoid_approx(*(*float32)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(float32(0))*uintptr(i))))
		}
	} else {
		for i = 0; i < N; i++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(float32(0))*uintptr(i))) = tansig_approx(*(*float32)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(float32(0))*uintptr(i))))
		}
	}
}
func compute_gru(gru *GRULayer, state *float32, input *float32) {
	var (
		i      int64
		N      int64
		M      int64
		stride int64
		tmp    [32]float32
		z      [32]float32
		r      [32]float32
		h      [32]float32
	)
	M = gru.Nb_inputs
	N = gru.Nb_neurons
	stride = N * 3
	for i = 0; i < N; i++ {
		z[i] = float32(*(*int8)(unsafe.Add(unsafe.Pointer(gru.Bias), i)))
	}
	gemm_accum(&z[0], gru.Input_weights, N, M, stride, input)
	gemm_accum(&z[0], gru.Recurrent_weights, N, N, stride, state)
	for i = 0; i < N; i++ {
		z[i] = sigmoid_approx(float32(float64(z[i]) * (1.0 / 128)))
	}
	for i = 0; i < N; i++ {
		r[i] = float32(*(*int8)(unsafe.Add(unsafe.Pointer(gru.Bias), N+i)))
	}
	gemm_accum(&r[0], (*int8)(unsafe.Add(unsafe.Pointer(gru.Input_weights), N)), N, M, stride, input)
	gemm_accum(&r[0], (*int8)(unsafe.Add(unsafe.Pointer(gru.Recurrent_weights), N)), N, N, stride, state)
	for i = 0; i < N; i++ {
		r[i] = sigmoid_approx(float32(float64(r[i]) * (1.0 / 128)))
	}
	for i = 0; i < N; i++ {
		h[i] = float32(*(*int8)(unsafe.Add(unsafe.Pointer(gru.Bias), N*2+i)))
	}
	for i = 0; i < N; i++ {
		tmp[i] = *(*float32)(unsafe.Add(unsafe.Pointer(state), unsafe.Sizeof(float32(0))*uintptr(i))) * r[i]
	}
	gemm_accum(&h[0], (*int8)(unsafe.Add(unsafe.Pointer(gru.Input_weights), N*2)), N, M, stride, input)
	gemm_accum(&h[0], (*int8)(unsafe.Add(unsafe.Pointer(gru.Recurrent_weights), N*2)), N, N, stride, &tmp[0])
	for i = 0; i < N; i++ {
		h[i] = z[i]**(*float32)(unsafe.Add(unsafe.Pointer(state), unsafe.Sizeof(float32(0))*uintptr(i))) + (1-z[i])*tansig_approx(float32(float64(h[i])*(1.0/128)))
	}
	for i = 0; i < N; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(state), unsafe.Sizeof(float32(0))*uintptr(i))) = h[i]
	}
}
