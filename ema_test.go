package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestEMA(t *testing.T) {
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

	var closeVals, expEMA []float64
	for _, record := range records {
		if len(record) < 1 {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[0], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		if len(record) > 1 && record[1] != "" {
			emaVal, err := strconv.ParseFloat(record[1], 64)
			if err == nil {
				expEMA = append(expEMA, emaVal)
			}
		}
	}

	timePeriod := 30
	result, err := EMA(closeVals, timePeriod)
	if err != nil {
		t.Fatalf("EMA计算失败: %v", err)
	}
	if len(result) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际%d", len(closeVals), len(result))
	}

	// 只校验有期望值的部分
	startIdx := len(closeVals) - len(expEMA)
	for i, exp := range expEMA {
		idx := startIdx + i
		if idx >= len(result) {
			break
		}
		actual := result[idx]
		if (exp == 0 && actual != 0) || (exp != 0 && (actual < exp-0.05 || actual > exp+0.05)) {
			t.Errorf("EMA[%d] 期望%.2f, 实际%.2f", idx, exp, actual)
		}
	}
}
