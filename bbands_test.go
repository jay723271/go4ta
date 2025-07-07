package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestBBands(t *testing.T) {
	file, err := os.Open("test_data/bbands.csv")
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

	var closeVals, expUpper, expMiddle, expLower []float64
	for _, record := range records {
		if len(record) < 1 {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[0], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		if len(record) > 1 && record[1] != "" {
			upperVal, err := strconv.ParseFloat(record[1], 64)
			if err == nil {
				expUpper = append(expUpper, upperVal)
			}
		}
		if len(record) > 2 && record[2] != "" {
			middleVal, err := strconv.ParseFloat(record[2], 64)
			if err == nil {
				expMiddle = append(expMiddle, middleVal)
			}
		}
		if len(record) > 3 && record[3] != "" {
			lowerVal, err := strconv.ParseFloat(record[3], 64)
			if err == nil {
				expLower = append(expLower, lowerVal)
			}
		}
	}

	timePeriod := 5
	nbDevUp, nbDevDn := 2.0, 2.0
	maType := 0 // 0=SMA
	upper, middle, lower, err := BBands(closeVals, timePeriod, nbDevUp, nbDevDn, maType)
	if err != nil {
		t.Fatalf("BBands计算失败: %v", err)
	}
	if len(upper) != len(closeVals) || len(middle) != len(closeVals) || len(lower) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际upper:%d middle:%d lower:%d", len(closeVals), len(upper), len(middle), len(lower))
	}

	// 只校验有期望值的部分
	startIdx := len(closeVals) - len(expUpper)
	for i := 0; i < len(expUpper); i++ {
		idx := startIdx + i
		if idx >= len(upper) {
			break
		}
		if diff := upper[idx] - expUpper[i]; diff < -0.05 || diff > 0.05 {
			t.Errorf("UpperBand[%d] 期望%.2f, 实际%.2f", idx, expUpper[i], upper[idx])
		}
		if diff := middle[idx] - expMiddle[i]; diff < -0.05 || diff > 0.05 {
			t.Errorf("MiddleBand[%d] 期望%.2f, 实际%.2f", idx, expMiddle[i], middle[idx])
		}
		if diff := lower[idx] - expLower[i]; diff < -0.05 || diff > 0.05 {
			t.Errorf("LowerBand[%d] 期望%.2f, 实际%.2f", idx, expLower[i], lower[idx])
		}
	}
}
