package go4ta

/*
#cgo LDFLAGS: -lta-lib -lm
#include <ta-lib/ta_libc.h>
#include <ta-lib/ta_func.h>
#include <stdlib.h>
TA_RetCode TA_BBANDS(int startIdx, int endIdx, const double inReal[], int optInTimePeriod, double optInNbDevUp, double optInNbDevDn, unsigned int optInMAType, int *outBegIdx, int *outNBElement, double *outUpperBand, double *outMiddleBand, double *outLowerBand);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// BBands 使用 CGO 直接调用 TA-Lib C 库来计算布林带（Bollinger Bands）。
//
// @param close      - 收盘价序列
// @param timePeriod - 计算周期（如20）
// @param nbDevUp    - 上轨标准差倍数（如2.0）
// @param nbDevDn    - 下轨标准差倍数（如2.0）
// @param maType     - 均线类型（如0=SMA，1=EMA等，见TA-Lib文档）
// @return upper, middle, lower - 三个与输入等长的结果序列
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func BBands(close []float64, timePeriod int, nbDevUp, nbDevDn float64, maType int) ([]float64, []float64, []float64, error) {
	if len(close) == 0 {
		return []float64{}, []float64{}, []float64{}, nil
	}
	if len(close) < timePeriod {
		return nil, nil, nil, fmt.Errorf("input data length (%d) is too small for the given timePeriod (%d)", len(close), timePeriod)
	}

	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	outUpper := make([]C.double, len(close))
	outMiddle := make([]C.double, len(close))
	outLower := make([]C.double, len(close))

	cOutUpper := (*C.double)(unsafe.Pointer(&outUpper[0]))
	cOutMiddle := (*C.double)(unsafe.Pointer(&outMiddle[0]))
	cOutLower := (*C.double)(unsafe.Pointer(&outLower[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_BBANDS(
		0,
		C.int(len(close)-1),
		cClose,
		C.int(timePeriod),
		C.double(nbDevUp),
		C.double(nbDevDn),
		C.TA_MAType(maType),
		&outBegIdx,
		&outNBElement,
		cOutUpper,
		cOutMiddle,
		cOutLower,
	)

	if retCode != C.TA_SUCCESS {
		return nil, nil, nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	upper := make([]float64, len(close))
	middle := make([]float64, len(close))
	lower := make([]float64, len(close))
	for i := 0; i < int(outNBElement); i++ {
		idx := int(outBegIdx) + i
		upper[idx] = float64(outUpper[i])
		middle[idx] = float64(outMiddle[i])
		lower[idx] = float64(outLower[i])
	}

	return upper, middle, lower, nil
}
