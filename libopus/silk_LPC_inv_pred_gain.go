package libopus

import "unsafe"

const QA = 24

func LPC_inverse_pred_gain_QA_c(A_QA [24]opus_int32, order int64) opus_int32 {
	var (
		k            int64
		n            int64
		mult2Q       int64
		invGain_Q30  opus_int32
		rc_Q31       opus_int32
		rc_mult1_Q30 opus_int32
		rc_mult2     opus_int32
		tmp1         opus_int32
		tmp2         opus_int32
	)
	invGain_Q30 = opus_int32(1*(1<<30) + 0.5)
	for k = order - 1; k > 0; k-- {
		if A_QA[k] > (opus_int32((1<<QA)*0.99975+0.5)) || A_QA[k] < -(opus_int32((1<<QA)*0.99975+0.5)) {
			return 0
		}
		rc_Q31 = -(opus_int32(opus_uint32(A_QA[k]) << opus_uint32(31-QA)))
		rc_mult1_Q30 = (opus_int32(1*(1<<30) + 0.5)) - (opus_int32((int64(rc_Q31) * int64(rc_Q31)) >> 32))
		invGain_Q30 = opus_int32(opus_uint32(opus_int32((int64(invGain_Q30)*int64(rc_mult1_Q30))>>32)) << 2)
		if invGain_Q30 < (opus_int32((1.0/MAX_PREDICTION_POWER_GAIN)*(1<<30) + 0.5)) {
			return 0
		}
		mult2Q = int64(32 - silk_CLZ32(func() opus_int32 {
			if rc_mult1_Q30 > 0 {
				return rc_mult1_Q30
			}
			return -rc_mult1_Q30
		}()))
		rc_mult2 = silk_INVERSE32_varQ(rc_mult1_Q30, mult2Q+30)
		for n = 0; n < (k+1)>>1; n++ {
			var tmp64 int64
			tmp1 = A_QA[n]
			tmp2 = A_QA[k-n-1]
			if mult2Q == 1 {
				tmp64 = ((int64(func() opus_int32 {
					if ((opus_uint32(tmp1) - opus_uint32(opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))) & 0x80000000) == 0 {
						if (tmp1 & ((opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}())) ^ 0x80000000) & 0x80000000) != 0 {
							return 0x80000000
						}
						return tmp1 - (opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}()))
					}
					if ((tmp1 ^ 0x80000000) & (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}())) & 0x80000000) != 0 {
						return silk_int32_MAX
					}
					return tmp1 - (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))
				}()) * int64(rc_mult2)) >> 1) + ((int64(func() opus_int32 {
					if ((opus_uint32(tmp1) - opus_uint32(opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))) & 0x80000000) == 0 {
						if (tmp1 & ((opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}())) ^ 0x80000000) & 0x80000000) != 0 {
							return 0x80000000
						}
						return tmp1 - (opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}()))
					}
					if ((tmp1 ^ 0x80000000) & (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}())) & 0x80000000) != 0 {
						return silk_int32_MAX
					}
					return tmp1 - (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))
				}()) * int64(rc_mult2)) & 1)
			} else {
				tmp64 = (((int64(func() opus_int32 {
					if ((opus_uint32(tmp1) - opus_uint32(opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))) & 0x80000000) == 0 {
						if (tmp1 & ((opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}())) ^ 0x80000000) & 0x80000000) != 0 {
							return 0x80000000
						}
						return tmp1 - (opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}()))
					}
					if ((tmp1 ^ 0x80000000) & (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}())) & 0x80000000) != 0 {
						return silk_int32_MAX
					}
					return tmp1 - (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp2) * int64(rc_Q31)) >> 1) + ((int64(tmp2) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp2) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))
				}()) * int64(rc_mult2)) >> (mult2Q - 1)) + 1) >> 1
			}
			if tmp64 > silk_int32_MAX || tmp64 < 0x80000000 {
				return 0
			}
			A_QA[n] = opus_int32(tmp64)
			if mult2Q == 1 {
				tmp64 = ((int64(func() opus_int32 {
					if ((opus_uint32(tmp2) - opus_uint32(opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))) & 0x80000000) == 0 {
						if (tmp2 & ((opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}())) ^ 0x80000000) & 0x80000000) != 0 {
							return 0x80000000
						}
						return tmp2 - (opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}()))
					}
					if ((tmp2 ^ 0x80000000) & (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}())) & 0x80000000) != 0 {
						return silk_int32_MAX
					}
					return tmp2 - (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))
				}()) * int64(rc_mult2)) >> 1) + ((int64(func() opus_int32 {
					if ((opus_uint32(tmp2) - opus_uint32(opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))) & 0x80000000) == 0 {
						if (tmp2 & ((opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}())) ^ 0x80000000) & 0x80000000) != 0 {
							return 0x80000000
						}
						return tmp2 - (opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}()))
					}
					if ((tmp2 ^ 0x80000000) & (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}())) & 0x80000000) != 0 {
						return silk_int32_MAX
					}
					return tmp2 - (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))
				}()) * int64(rc_mult2)) & 1)
			} else {
				tmp64 = (((int64(func() opus_int32 {
					if ((opus_uint32(tmp2) - opus_uint32(opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))) & 0x80000000) == 0 {
						if (tmp2 & ((opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}())) ^ 0x80000000) & 0x80000000) != 0 {
							return 0x80000000
						}
						return tmp2 - (opus_int32(func() int64 {
							if 31 == 1 {
								return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
							}
							return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
						}()))
					}
					if ((tmp2 ^ 0x80000000) & (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}())) & 0x80000000) != 0 {
						return silk_int32_MAX
					}
					return tmp2 - (opus_int32(func() int64 {
						if 31 == 1 {
							return ((int64(tmp1) * int64(rc_Q31)) >> 1) + ((int64(tmp1) * int64(rc_Q31)) & 1)
						}
						return (((int64(tmp1) * int64(rc_Q31)) >> (31 - 1)) + 1) >> 1
					}()))
				}()) * int64(rc_mult2)) >> (mult2Q - 1)) + 1) >> 1
			}
			if tmp64 > silk_int32_MAX || tmp64 < 0x80000000 {
				return 0
			}
			A_QA[k-n-1] = opus_int32(tmp64)
		}
	}
	if A_QA[k] > (opus_int32((1<<QA)*0.99975+0.5)) || A_QA[k] < -(opus_int32((1<<QA)*0.99975+0.5)) {
		return 0
	}
	rc_Q31 = -(opus_int32(opus_uint32(A_QA[0]) << opus_uint32(31-QA)))
	rc_mult1_Q30 = (opus_int32(1*(1<<30) + 0.5)) - (opus_int32((int64(rc_Q31) * int64(rc_Q31)) >> 32))
	invGain_Q30 = opus_int32(opus_uint32(opus_int32((int64(invGain_Q30)*int64(rc_mult1_Q30))>>32)) << 2)
	if invGain_Q30 < (opus_int32((1.0/MAX_PREDICTION_POWER_GAIN)*(1<<30) + 0.5)) {
		return 0
	}
	return invGain_Q30
}
func silk_LPC_inverse_pred_gain_c(A_Q12 *opus_int16, order int64) opus_int32 {
	var (
		k       int64
		Atmp_QA [24]opus_int32
		DC_resp opus_int32 = 0
	)
	for k = 0; k < order; k++ {
		DC_resp += opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(opus_int16(0))*uintptr(k))))
		Atmp_QA[k] = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(A_Q12), unsafe.Sizeof(opus_int16(0))*uintptr(k))))) << opus_uint32(QA-12))
	}
	if DC_resp >= 4096 {
		return 0
	}
	return LPC_inverse_pred_gain_QA_c(Atmp_QA, order)
}
