package libopus

import "unsafe"

const MAX_ITERATIONS_RESIDUAL_NRG = 10
const REGULARIZATION_FACTOR = 1e-08

func silk_residual_energy_covar_FLP(c *float32, wXX *float32, wXx *float32, wxx float32, D int) float32 {
	var (
		i              int
		j              int
		k              int
		tmp            float32
		nrg            float32 = 0.0
		regularization float32
	)
	regularization = REGULARIZATION_FACTOR * (*(*float32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(float32(0))*0)) + *(*float32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(float32(0))*uintptr(D*D-1))))
	for k = 0; k < MAX_ITERATIONS_RESIDUAL_NRG; k++ {
		nrg = wxx
		tmp = 0.0
		for i = 0; i < D; i++ {
			tmp += *(*float32)(unsafe.Add(unsafe.Pointer(wXx), unsafe.Sizeof(float32(0))*uintptr(i))) * *(*float32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(float32(0))*uintptr(i)))
		}
		nrg -= tmp * 2.0
		for i = 0; i < D; i++ {
			tmp = 0.0
			for j = i + 1; j < D; j++ {
				tmp += (*((*float32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(float32(0))*uintptr(i+D*j))))) * *(*float32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(float32(0))*uintptr(j)))
			}
			nrg += *(*float32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(float32(0))*uintptr(i))) * (tmp*2.0 + (*((*float32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(float32(0))*uintptr(i+D*i)))))**(*float32)(unsafe.Add(unsafe.Pointer(c), unsafe.Sizeof(float32(0))*uintptr(i))))
		}
		if nrg > 0 {
			break
		} else {
			for i = 0; i < D; i++ {
				*((*float32)(unsafe.Add(unsafe.Pointer(wXX), unsafe.Sizeof(float32(0))*uintptr(i+D*i)))) += regularization
			}
			regularization *= 2.0
		}
	}
	if k == MAX_ITERATIONS_RESIDUAL_NRG {
		nrg = 1.0
	}
	return nrg
}
func silk_residual_energy_FLP(nrgs [4]float32, x []float32, a [2][16]float32, gains []float32, subfr_length int, nb_subfr int, LPC_order int) {
	var (
		shift       int
		LPC_res_ptr *float32
		LPC_res     [192]float32
	)
	LPC_res_ptr = &LPC_res[LPC_order]
	shift = LPC_order + subfr_length
	silk_LPC_analysis_filter_FLP(LPC_res[:], a[0][:], []float32(&x[shift*0]), shift*2, LPC_order)
	nrgs[0] = float32(float64(gains[0]*gains[0]) * silk_energy_FLP([]float32((*float32)(unsafe.Add(unsafe.Pointer(LPC_res_ptr), unsafe.Sizeof(float32(0))*uintptr(shift*0)))), subfr_length))
	nrgs[1] = float32(float64(gains[1]*gains[1]) * silk_energy_FLP([]float32((*float32)(unsafe.Add(unsafe.Pointer(LPC_res_ptr), unsafe.Sizeof(float32(0))*uintptr(shift*1)))), subfr_length))
	if nb_subfr == MAX_NB_SUBFR {
		silk_LPC_analysis_filter_FLP(LPC_res[:], a[1][:], []float32(&x[shift*2]), shift*2, LPC_order)
		nrgs[2] = float32(float64(gains[2]*gains[2]) * silk_energy_FLP([]float32((*float32)(unsafe.Add(unsafe.Pointer(LPC_res_ptr), unsafe.Sizeof(float32(0))*uintptr(shift*0)))), subfr_length))
		nrgs[3] = float32(float64(gains[3]*gains[3]) * silk_energy_FLP([]float32((*float32)(unsafe.Add(unsafe.Pointer(LPC_res_ptr), unsafe.Sizeof(float32(0))*uintptr(shift*1)))), subfr_length))
	}
}
