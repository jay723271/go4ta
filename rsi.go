package go4ta

/*
#cgo LDFLAGS: -lta-lib -lm
#include <ta-lib/ta_libc.h>
#include <ta-lib/ta_func.h>
#include <stdlib.h>
TA_RetCode TA_RSI(int startIdx, int endIdx, const double inReal[], int optInTimePeriod, int *outBegIdx, int *outNBElement, double outReal[]);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// CalculateRSI 使用 CGO 直接调用 TA-Lib C 库来计算相对强弱指数 (RSI)。
//
// @param close      - 收盘价序列
// @param timePeriod - 计算周期 (例如 14)
// @return []float64 - RSI 结果序列。其长度与输入序列相同。
//
//	由于计算需要一定数量的初始数据，序列开头的部分值为 0。
//
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func RSI(close []float64, timePeriod int) ([]float64, error) {
	// --- 输入数据校验 ---
	if len(close) == 0 {
		return []float64{}, nil
	}
	if len(close) < timePeriod {
		return nil, fmt.Errorf("input data length (%d) is too small for the given timePeriod (%d)", len(close), timePeriod)
	}

	// --- 准备 C 语言格式的输入数据 ---
	cClose := (*C.double)(unsafe.Pointer(&close[0]))

	// --- 准备 C 语言格式的输出缓冲区 ---
	output := make([]C.double, len(close))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	// --- 准备用于接收 TA-Lib 输出元数据的变量 ---
	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	// --- 调用 C 函数 ---
	retCode := C.TA_RSI(
		0,                   // startIdx: 从输入数据的第一个元素开始
		C.int(len(close)-1), // endIdx: 到输入数据的最后一个元素结束
		cClose,              // inReal
		C.int(timePeriod),   // optInTimePeriod
		&outBegIdx,          // outBegIdx (输出参数)
		&outNBElement,       // outNBElement (输出参数)
		cOutput,             // outReal (输出缓冲区)
	)

	// --- 检查 C 函数调用结果 ---
	if retCode != C.TA_SUCCESS {
		return nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	// --- 将 C 输出结果转换为 Go 切片 ---
	result := make([]float64, len(close))
	for i := 0; i < int(outNBElement); i++ {
		result[int(outBegIdx)+i] = float64(output[i])
	}

	return result, nil
}
