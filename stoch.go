package go4ta

/*
#cgo LDFLAGS: -lta-lib -lm
#include <ta-lib/ta_libc.h>
#include <ta-lib/ta_func.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// STOCH 使用 CGO 直接调用 TA-Lib C 库来计算随机指标（KDJ）。
//
// @param high        - 最高价序列
// @param low         - 最低价序列
// @param close       - 收盘价序列
// @param fastKPeriod - K线周期
// @param slowKPeriod - 慢K周期
// @param slowDPeriod - 慢D周期
// @param maTypeK     - K均线类型
// @param maTypeD     - D均线类型
// @return slowK, slowD - 两个与输入等长的结果序列
// @return error      - 如果输入数据无效或 C 库调用失败，则返回错误。
func STOCH(high, low, close []float64, fastKPeriod, slowKPeriod, slowDPeriod, maTypeK, maTypeD int) ([]float64, []float64, error) {
	if len(high) == 0 || len(low) == 0 || len(close) == 0 {
		return []float64{}, []float64{}, nil
	}
	if len(high) != len(low) || len(low) != len(close) {
		return nil, nil, fmt.Errorf("input slices (high, low, close) must have the same length")
	}

	cHigh := (*C.double)(unsafe.Pointer(&high[0]))
	cLow := (*C.double)(unsafe.Pointer(&low[0]))
	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	outSlowK := make([]C.double, len(high))
	outSlowD := make([]C.double, len(high))
	cOutSlowK := (*C.double)(unsafe.Pointer(&outSlowK[0]))
	cOutSlowD := (*C.double)(unsafe.Pointer(&outSlowD[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_STOCH(
		0,
		C.int(len(high)-1),
		cHigh,
		cLow,
		cClose,
		C.int(fastKPeriod),
		C.int(slowKPeriod),
		C.TA_MAType(maTypeK),
		C.int(slowDPeriod),
		C.TA_MAType(maTypeD),
		&outBegIdx,
		&outNBElement,
		cOutSlowK,
		cOutSlowD,
	)

	if retCode != C.TA_SUCCESS {
		return nil, nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	slowK := make([]float64, len(high))
	slowD := make([]float64, len(high))
	for i := 0; i < int(outNBElement); i++ {
		idx := int(outBegIdx) + i
		slowK[idx] = float64(outSlowK[i])
		slowD[idx] = float64(outSlowD[i])
	}

	return slowK, slowD, nil
}
