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

// STOCHRSI 使用 CGO 直接调用 TA-Lib C 库来计算随机RSI（Stochastic RSI）。
//
// @param close        - 收盘价序列
// @param timePeriod   - RSI周期
// @param fastKPeriod  - K线周期
// @param fastDPeriod  - D线周期
// @param maType       - 均线类型
// @return fastK, fastD - 两个与输入等长的结果序列
// @return error       - 如果输入数据无效或 C 库调用失败，则返回错误。
func STOCHRSI(close []float64, timePeriod, fastKPeriod, fastDPeriod, maType int) ([]float64, []float64, error) {
	if len(close) == 0 {
		return []float64{}, []float64{}, nil
	}
	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	outFastK := make([]C.double, len(close))
	outFastD := make([]C.double, len(close))
	cOutFastK := (*C.double)(unsafe.Pointer(&outFastK[0]))
	cOutFastD := (*C.double)(unsafe.Pointer(&outFastD[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_STOCHRSI(
		0,
		C.int(len(close)-1),
		cClose,
		C.int(timePeriod),
		C.int(fastKPeriod),
		C.int(fastDPeriod),
		C.TA_MAType(maType),
		&outBegIdx,
		&outNBElement,
		cOutFastK,
		cOutFastD,
	)

	if retCode != C.TA_SUCCESS {
		return nil, nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	fastK := make([]float64, len(close))
	fastD := make([]float64, len(close))
	for i := 0; i < int(outNBElement); i++ {
		idx := int(outBegIdx) + i
		fastK[idx] = float64(outFastK[i])
		fastD[idx] = float64(outFastD[i])
	}

	return fastK, fastD, nil
}
