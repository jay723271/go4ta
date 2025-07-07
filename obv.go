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

// OBV 使用 CGO 直接调用 TA-Lib C 库来计算能量潮（On Balance Volume）。
//
// @param close   - 收盘价序列
// @param volume  - 成交量序列
// @return []float64 - OBV结果序列，与输入等长。
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func OBV(close, volume []float64) ([]float64, error) {
	if len(close) == 0 || len(volume) == 0 {
		return []float64{}, nil
	}
	if len(close) != len(volume) {
		return nil, fmt.Errorf("input slices (close, volume) must have the same length")
	}

	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	cVolume := (*C.double)(unsafe.Pointer(&volume[0]))
	output := make([]C.double, len(close))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_OBV(
		0,
		C.int(len(close)-1),
		cClose,
		cVolume,
		&outBegIdx,
		&outNBElement,
		cOutput,
	)

	if retCode != C.TA_SUCCESS {
		return nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	result := make([]float64, len(close))
	for i := 0; i < int(outNBElement); i++ {
		result[int(outBegIdx)+i] = float64(output[i])
	}

	return result, nil
}
