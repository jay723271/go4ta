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

// PPO 使用 CGO 直接调用 TA-Lib C 库来计算百分比价格振荡器（PPO）。
//
// @param close        - 收盘价序列
// @param fastPeriod   - 快速均线周期
// @param slowPeriod   - 慢速均线周期
// @param maType       - 均线类型（如0=SMA，1=EMA等，见TA-Lib文档）
// @return []float64   - PPO结果序列，与输入等长，未计算部分为0。
// @return error       - 如果输入数据无效或 C 库调用失败，则返回错误。
func PPO(close []float64, fastPeriod, slowPeriod, maType int) ([]float64, error) {
	if len(close) == 0 {
		return []float64{}, nil
	}
	if len(close) < fastPeriod || len(close) < slowPeriod {
		return nil, fmt.Errorf("input data length (%d) is too small for the given periods", len(close))
	}

	cClose := (*C.double)(unsafe.Pointer(&close[0]))
	output := make([]C.double, len(close))
	cOutput := (*C.double)(unsafe.Pointer(&output[0]))

	outBegIdx := C.int(0)
	outNBElement := C.int(0)

	retCode := C.TA_PPO(
		0,
		C.int(len(close)-1),
		cClose,
		C.int(fastPeriod),
		C.int(slowPeriod),
		C.TA_MAType(maType),
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

// PPOWithSignal 计算PPO、信号线（PPO的EMA）和柱状图（PPO-信号线）
func PPOWithSignal(close []float64, fastPeriod, slowPeriod, signalPeriod, maType int) (ppo, signal, hist []float64, err error) {
	ppo, err = PPO(close, fastPeriod, slowPeriod, maType)
	if err != nil {
		return
	}
	// 找到第一个有效PPO的位置
	firstValid := 0
	for firstValid < len(ppo) && ppo[firstValid] == 0 {
		firstValid++
	}
	// 只对有效区间做EMA
	signalValid, err := EMA(ppo[firstValid:], signalPeriod)
	if err != nil {
		return
	}
	signal = make([]float64, len(ppo))
	hist = make([]float64, len(ppo))
	for i := range ppo {
		if i < firstValid {
			signal[i] = 0
			hist[i] = 0
		} else {
			signal[i] = signalValid[i-firstValid]
			hist[i] = ppo[i] - signal[i]
		}
	}
	return
}
