package go4ta

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"testing"
)

func parseFloatOrNaN(s string) float64 {
	if s == "" || s == "nan" || s == "NaN" {
		return math.NaN()
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return v
}

func TestSuperTrend(t *testing.T) {
	file, err := os.Open("test_data/super_trend.csv")
	if err != nil {
		t.Fatalf("无法打开测试数据文件: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("无法读取CSV数据: %v", err)
	}

	// 跳过表头
	records = records[1:]

	var high, low, close, expSuper, expLower, expUpper []float64
	var expDir []int
	for _, record := range records {
		high = append(high, parseFloatOrNaN(record[0]))
		low = append(low, parseFloatOrNaN(record[1]))
		close = append(close, parseFloatOrNaN(record[2]))
		expSuper = append(expSuper, parseFloatOrNaN(record[3]))

		// CSV中的���向值是 "1.00" 这样的浮点数格式，因此先按浮点数解析
		dirFloat := parseFloatOrNaN(record[4])
		if math.IsNaN(dirFloat) {
			expDir = append(expDir, 0)
		} else {
			expDir = append(expDir, int(dirFloat))
		}

		expLower = append(expLower, parseFloatOrNaN(record[5]))
		expUpper = append(expUpper, parseFloatOrNaN(record[6]))
	}

	period := 7
	multiplier := 3.0
	// 更新函数调用以处理错误返回值
	super, dir, lower, upper, err := SuperTrend(high, low, close, period, multiplier)
	if err != nil {
		t.Fatalf("SuperTrend 函数返回错误: %v", err)
	}

	if len(super) != len(records) {
		t.Fatalf("输出长度不匹配: 期望 %d, 实际 %d", len(records), len(super))
	}

	for i := range close {
		// 检查 SuperTrend
		if math.IsNaN(expSuper[i]) {
			if !math.IsNaN(super[i]) {
				t.Errorf("SuperTrend[%d]: 期望 NaN, 实际 %.2f", i, super[i])
			}
		} else if math.IsNaN(super[i]) || math.Abs(super[i]-expSuper[i]) > 0.05 {
			t.Errorf("SuperTrend[%d]: 期望 %.2f, 实际 %.2f", i, expSuper[i], super[i])
		}

		// 检查 Direction
		if dir[i] != float64(expDir[i]) {
			t.Errorf("Direction[%d]: 期望 %d, 实际 %.0f", i, expDir[i], dir[i])
		}

		// 检查 LowerBand
		if math.IsNaN(expLower[i]) {
			if !math.IsNaN(lower[i]) {
				t.Errorf("LowerBand[%d]: 期望 NaN, 实际 %.2f", i, lower[i])
			}
		} else if math.IsNaN(lower[i]) || math.Abs(lower[i]-expLower[i]) > 0.05 {
			t.Errorf("LowerBand[%d]: 期望 %.2f, 实际 %.2f", i, expLower[i], lower[i])
		}

		// 检查 UpperBand
		if math.IsNaN(expUpper[i]) {
			if !math.IsNaN(upper[i]) {
				t.Errorf("UpperBand[%d]: 期望 NaN, 实际 %.2f", i, upper[i])
			}
		} else if math.IsNaN(upper[i]) || math.Abs(upper[i]-expUpper[i]) > 0.05 {
			t.Errorf("UpperBand[%d]: 期望 %.2f, 实际 %.2f", i, expUpper[i], upper[i])
		}
	}
}