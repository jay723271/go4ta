package go4ta

/*
#cgo LDFLAGS: -lta-lib -lm
#include <ta-lib/ta_libc.h>
#include <ta-lib/ta_func.h>
#include <stdlib.h>
TA_RetCode TA_MACD(int startIdx, int endIdx, const double inReal[], int optInFastPeriod, int optInSlowPeriod, int optInSignalPeriod, int *outBegIdx, int *outNBElement, double *outMACD, double *outSignal, double *outHist);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// MACD 使用 CGO 直接调用 TA-Lib C 库来计算MACD指标。
//
// @param close        - 收盘价序列
// @param fastPeriod   - 快速均线周期
// @param slowPeriod   - 慢速均线周期
// @param signalPeriod - 信号线周期
// @return macd, signal, hist - 三个与输入等长的结果序列
// @return error       - 如果输入数据无效或 C 库调用失败，则返回错误。
func MACD(close []float64, fastPeriod, slowPeriod, signalPeriod int) ([]float64, []float64, []float64, error) {
	if len(close) == 0 {
		return []float64{}, []float64{}, []float64{}, nil
	}
	if len(close) < slowPeriod || len(close) < fastPeriod || len(close) < signalPeriod {
		return nil, nil, nil, fmt.Errorf("input data length (%d) is too small for the given periods", len(close))
	}

	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	outMACD := make([]C.double, len(close))
	outSignal := make([]C.double, len(close))
	outHist := make([]C.double, len(close))

	cOutMACD := (*C.double)(unsafe.Pointer(&outMACD[0]))
	cOutSignal := (*C.double)(unsafe.Pointer(&outSignal[0]))
	cOutHist := (*C.double)(unsafe.Pointer(&outHist[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_MACD(
		0,
		C.int(len(close)-1),
		cClose,
		C.int(fastPeriod),
		C.int(slowPeriod),
		C.int(signalPeriod),
		&outBegIdx,
		&outNBElement,
		cOutMACD,
		cOutSignal,
		cOutHist,
	)

	if retCode != C.TA_SUCCESS {
		return nil, nil, nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	macd := make([]float64, len(close))
	signal := make([]float64, len(close))
	hist := make([]float64, len(close))
	for i := 0; i < int(outNBElement); i++ {
		idx := int(outBegIdx) + i
		macd[idx] = float64(outMACD[i])
		signal[idx] = float64(outSignal[i])
		hist[idx] = float64(outHist[i])
	}

	return macd, signal, hist, nil
}
