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

// LinearReg 使用 CGO 直接调用 TA-Lib C 库来计算线性回归（LINEARREG）。
//
// @param close      - 收盘价序列
// @param timePeriod - 计算周期（如14）
// @return []float64 - 线性回归主值序列，与输入等长，未计算部分为0。
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func LinearReg(close []float64, timePeriod int) ([]float64, error) {
	if len(close) == 0 {
		return []float64{}, nil
	}
	if len(close) < timePeriod {
		return nil, fmt.Errorf("input data length (%d) is too small for the given timePeriod (%d)", len(close), timePeriod)
	}

	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	output := make([]C.double, len(close))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_LINEARREG(
		0,
		C.int(len(close)-1),
		cClose,
		C.int(timePeriod),
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
