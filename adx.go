package go4ta

/*
#cgo LDFLAGS: -lta-lib -lm
#include <ta-lib/ta_libc.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ADX 使用 CGO 直接调用 TA-Lib C 库来计算平均趋向指数 (ADX)。
//
// @param high       - 最高价序列
// @param low        - 最低价序列
// @param close      - 收盘价序列
// @param timePeriod - 计算周期 (例如 14)
// @return []float64 - ADX 结果序列。其长度与输入序列相同。
//
//	由于计算需要一定数量的初始数据，序列开头的部分值为 0。
//
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func ADX(high, low, close []float64, timePeriod int) ([]float64, error) {
	// --- 输入数据校验 ---
	if len(high) != len(low) || len(low) != len(close) {
		return nil, fmt.Errorf("input slices (high, low, close) must have the same length")
	}
	if len(high) == 0 {
		return []float64{}, nil
	}
	if len(high) < timePeriod {
		return nil, fmt.Errorf("input data length (%d) is too small for the given timePeriod (%d)", len(high), timePeriod)
	}

	// --- 准备 C 语言格式的输入数据 ---
	cHigh := (*C.double)(unsafe.Pointer(&high[0]))
	cLow := (*C.double)(unsafe.Pointer(&low[0]))
	cClose := (*C.double)(unsafe.Pointer(&close[0]))

	// --- 准备 C 语言格式的输出缓冲区 ---
	output := make([]C.double, len(high))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	// --- 准备用于接收 TA-Lib 输出元数据的变量 ---
	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	// --- 调用 C 函数 ---
	retCode := C.TA_ADX(
		0,                  // startIdx
		C.int(len(high)-1), // endIdx
		cHigh,              // inHigh
		cLow,               // inLow
		cClose,             // inClose
		C.int(timePeriod),  // optInTimePeriod
		&outBegIdx,         // outBegIdx
		&outNBElement,      // outNBElement
		cOutput,            // outReal
	)

	// --- 检查 C 函数调用结果 ---
	if retCode != C.TA_SUCCESS {
		return nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	// --- 将 C 输出结果转换为 Go 切片 ---
	result := make([]float64, len(high))
	for i := 0; i < int(outNBElement); i++ {
		result[int(outBegIdx)+i] = float64(output[i])
	}

	return result, nil
}
