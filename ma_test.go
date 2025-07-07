package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestMA(t *testing.T) {
	file, err := os.Open("test_data/ema.csv")
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

	var closeVals []float64
	var expMAIdx []int
	var expMA []float64
	for i, record := range records {
		if len(record) < 1 {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[0], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		if len(record) > 1 && record[1] != "" {
			maVal, err := strconv.ParseFloat(record[1], 64)
			if err == nil {
				expMAIdx = append(expMAIdx, i)
				expMA = append(expMA, maVal)
			}
		}
	}

	timePeriod := 30 // ema.csv为30周期EMA
	maType := 1      // 1=EMA
	result, err := MA(closeVals, timePeriod, maType)
	if err != nil {
		t.Fatalf("MA计算失败: %v", err)
	}
	if len(result) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际%d", len(closeVals), len(result))
	}

	if len(result) == 0 {
		t.Fatalf("MA返回结果为空")
	}
	for i, idx := range expMAIdx {
		if idx >= len(result) {
			t.Errorf("MA结果长度不足，无法对比第%d行", idx)
			continue
		}
		actual := result[idx]
		exp := expMA[i]
		if (exp == 0 && actual != 0) || (exp != 0 && (actual < exp-0.05 || actual > exp+0.05)) {
			t.Errorf("MA[%d] 期望%.2f, 实际%.2f", idx, exp, actual)
		}
	}
}
