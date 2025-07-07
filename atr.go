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

// CalculateATR 使用 CGO 直接调用 TA-Lib C 库来计算平均真实波幅 (ATR)。
//
// @param high       - 最高价序列
// @param low        - 最低价序列
// @param close      - 收盘价序列
// @param timePeriod - 计算周期 (例如 14)
// @return []float64 - ATR 结果序列。其长度与输入序列相同。
//
//	由于计算需要一定数量的初始数据，序列开头的部分值为 0。
//
// @return error     - 如果输入数据无效或 C 库调用失败，则返回错误。
func ATR(high, low, close []float64, timePeriod int) ([]float64, error) {
	// --- 输入数据校验 ---
	if len(high) != len(low) || len(low) != len(close) {
		return nil, fmt.Errorf("input slices (high, low, close) must have the same length")
	}
	if len(high) == 0 {
		return []float64{}, nil
	}
	// TA-Lib 的 ATR 函数要求输入数据长度至少为 timePeriod
	// See: https://github.com/ta-lib/ta-lib/blob/master/src/ta_func/ta_ATR.c#L206
	if len(high) < timePeriod {
		return nil, fmt.Errorf("input data length (%d) is too small for the given timePeriod (%d)", len(high), timePeriod)
	}

	// --- 准备 C 语言格式的输入数据 ---
	cHigh := (*C.double)(unsafe.Pointer(&high[0]))
	cLow := (*C.double)(unsafe.Pointer(&low[0]))
	cClose := (*C.double)(unsafe.Pointer(&close[0]))

	// --- 准备 C 语言格式的输出缓冲区 ---
	// TA-Lib 会将结果写入我们提供的缓冲区
	output := make([]C.double, len(high))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	// --- 准备用于接收 TA-Lib 输出元数据的变量 ---
	// outBegIdx 会告诉我们有效数据是从哪个索引开始的
	// outNBElement 会告诉我们输出了多少个有效数据点
	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	// --- 调用 C 函数 ---
	retCode := C.TA_ATR(
		0,                  // startIdx: 从输入数据的第一个元素开始
		C.int(len(high)-1), // endIdx: 到输入数据的最后一个元素结束
		cHigh,              // inHigh
		cLow,               // inLow
		cClose,             // inClose
		C.int(timePeriod),  // optInTimePeriod
		&outBegIdx,         // outBegIdx (输出参数)
		&outNBElement,      // outNBElement (输出参数)
		cOutput,            // outReal (输出缓冲区)
	)

	// --- 检查 C 函数调用结果 ---
	if retCode != C.TA_SUCCESS {
		return nil, fmt.Errorf("TA-Lib C call failed with exit code: %d", retCode)
	}

	// --- 将 C 输出结果转换为 Go 切片 ---
	// 创建一个与输入等长的 Go 切片，未计算部分默认为 0
	result := make([]float64, len(high))

	// TA-Lib 的输出结果是从 output[0] 开始填充的，共 outNBElement 个。
	// outBegIdx 指明了第一个有效结果对应于输入序列的哪个位置。
	// 例如，如果 outBegIdx 是 14，那么 output[0] 的值应该放到 result[14] 的位置。
	for i := 0; i < int(outNBElement); i++ {
		result[int(outBegIdx)+i] = float64(output[i])
	}

	return result, nil
}
