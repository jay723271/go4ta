package go4ta

import (
	"fmt"
	"math"
)

// SuperTrend indicator calculation, logic adapted to match popular library standards.
//
// @param high, low, close - 价格序列
// @param period           - ATR周期
// @param multiplier       - ATR倍数
// @return superTrend      - SuperTrend主线
// @return direction       - 方向（1=多头，-1=空头）
// @return lowerBand       - 下轨
// @return upperBand       - 上轨
// @return err             - 错误信息

// 计算基础上下轨
func calcSupTrdBasicBands(high, low, atr []float64, multiplier float64, i int) (float64, float64) {
	hl2 := (high[i] + low[i]) / 2
	basicUpper := hl2 + multiplier*atr[i]
	basicLower := hl2 - multiplier*atr[i]
	return basicUpper, basicLower
}

// 计算最终上下轨
func calcSupTrdFinalBands(prevFinal, basic float64, isUpper bool) float64 {
	if math.IsNaN(prevFinal) {
		return basic
	}
	if isUpper {
		if basic < prevFinal {
			return basic
		}
		return prevFinal
	} else {
		if basic > prevFinal {
			return basic
		}
		return prevFinal
	}
}

// SuperTrend 主函数
func SuperTrend(high, low, close []float64, period int, multiplier float64) (superTrend, direction, lowerBand, upperBand []float64, err error) {
	n := len(close)
	if n == 0 {
		return nil, nil, nil, nil, nil
	}

	atr, atrErr := ATR(high, low, close, period)
	if atrErr != nil {
		return nil, nil, nil, nil, fmt.Errorf("ATR calculation failed: %w", atrErr)
	}

	superTrend = make([]float64, n)
	direction = make([]float64, n)
	finalLowerBand := make([]float64, n)
	finalUpperBand := make([]float64, n)
	lowerBand = make([]float64, n)
	upperBand = make([]float64, n)

	for i := 0; i < n; i++ {
		superTrend[i] = math.NaN()
		direction[i] = 0
		finalLowerBand[i] = math.NaN()
		finalUpperBand[i] = math.NaN()
		lowerBand[i] = math.NaN()
		upperBand[i] = math.NaN()
	}

	for i := 0; i < n; i++ {
		if i < period || math.IsNaN(atr[i]) {
			// ATR未满周期，全部为NaN/0
			continue
		}

		basicUpper, basicLower := calcSupTrdBasicBands(high, low, atr, multiplier, i)
		if i == period {
			finalLowerBand[i] = basicLower
			finalUpperBand[i] = basicUpper
		} else {
			finalLowerBand[i] = calcSupTrdFinalBands(finalLowerBand[i-1], basicLower, false)
			finalUpperBand[i] = calcSupTrdFinalBands(finalUpperBand[i-1], basicUpper, true)
		}

		prevDir := 0
		if i > 0 {
			prevDir = int(direction[i-1])
		}

		// 默认方向为0，只有趋势成立后才赋值
		direction[i] = float64(prevDir)

		if prevDir == 0 {
			// 第一个有效点，方向初始化为1
			direction[i] = 1
		} else if prevDir == 1 && close[i] < finalLowerBand[i] {
			direction[i] = -1
		} else if prevDir == -1 && close[i] > finalUpperBand[i] {
			direction[i] = 1
		}

		// 只有方向为1或-1时才赋值主线和上下轨，否则为NaN
		switch direction[i] {
		case 1:
			superTrend[i] = finalLowerBand[i]
			lowerBand[i] = finalLowerBand[i]
			upperBand[i] = math.NaN()
			finalUpperBand[i] = math.NaN()
		case -1:
			superTrend[i] = finalUpperBand[i]
			upperBand[i] = finalUpperBand[i]
			lowerBand[i] = math.NaN()
			finalLowerBand[i] = math.NaN()
		default:
			superTrend[i] = math.NaN()
			lowerBand[i] = math.NaN()
			upperBand[i] = math.NaN()
		}
	}

	return superTrend, direction, lowerBand, upperBand, nil
}
