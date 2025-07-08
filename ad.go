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

// AD 指标（Accumulation/Distribution Line）
//
// @param high   - 最高价序列
// @param low    - 最低价序列
// @param close  - 收盘价序列
// @param volume - 成交量序列
// @return []float64 - AD结果序列，与输入等长。
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func AD(high, low, close, volume []float64) ([]float64, error) {
	if len(high) == 0 || len(low) == 0 || len(close) == 0 || len(volume) == 0 {
		return []float64{}, nil
	}
	if len(high) != len(low) || len(low) != len(close) || len(close) != len(volume) {
		return nil, fmt.Errorf("input slices (high, low, close, volume) must have the same length")
	}

	cHigh := (*C.double)(unsafe.Pointer(&high[0]))
	cLow := (*C.double)(unsafe.Pointer(&low[0]))
	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	cVolume := (*C.double)(unsafe.Pointer(&volume[0]))
	output := make([]C.double, len(high))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_AD(
		0,
		C.int(len(high)-1),
		cHigh,
		cLow,
		cClose,
		cVolume,
		&outBegIdx,
		&outNBElement,
		cOutput,
	)

	if retCode != C.TA_SUCCESS {
		return nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	result := make([]float64, len(high))
	for i := 0; i < int(outNBElement); i++ {
		result[int(outBegIdx)+i] = float64(output[i])
	}

	return result, nil
}
