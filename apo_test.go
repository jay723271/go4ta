package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestAPO(t *testing.T) {
	file, err := os.Open("test_data/apo.csv")
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
	var expAPOIdx []int
	var expAPO []float64
	for i, record := range records {
		if len(record) < 5 {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[2], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		if record[4] != "" {
			apoVal, err := strconv.ParseFloat(record[4], 64)
			if err == nil {
				expAPOIdx = append(expAPOIdx, i)
				expAPO = append(expAPO, apoVal)
			}
		}
	}

	fastPeriod, slowPeriod, maType := 12, 26, 0
	result, err := APO(closeVals, fastPeriod, slowPeriod, maType)
	if err != nil {
		t.Fatalf("APO计算失败: %v", err)
	}
	if len(result) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际%d", len(closeVals), len(result))
	}
	if len(result) == 0 {
		t.Fatalf("APO返回结果为空")
	}
	for i, idx := range expAPOIdx {
		if idx >= len(result) {
			t.Errorf("APO结果长度不足，无法对比第%d行", idx)
			continue
		}
		actual := result[idx]
		exp := expAPO[i]
		if (exp == 0 && actual != 0) || (exp != 0 && (actual < exp-0.05 || actual > exp+0.05)) {
			t.Errorf("APO[%d] 期望%.2f, 实际%.2f", idx, exp, actual)
		}
	}
}
