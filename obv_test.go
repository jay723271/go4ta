package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestOBV(t *testing.T) {
	file, err := os.Open("test_data/obv.csv")
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

	var closeVals, volumeVals, expOBV []float64
	for i, record := range records {
		if len(record) < 5 {
			t.Fatalf("第%d行数据不足5列: %v", i+2, record)
		}
		closeVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			t.Fatalf("第%d行Close解析失败: %v", i+2, err)
		}
		closeVals = append(closeVals, closeVal)
		volumeVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			t.Fatalf("第%d行Volume解析失败: %v", i+2, err)
		}
		volumeVals = append(volumeVals, volumeVal)
		obvVal, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			t.Fatalf("第%d行OBV解析失败: %v", i+2, err)
		}
		expOBV = append(expOBV, obvVal)
	}

	result, err := OBV(closeVals, volumeVals)
	if err != nil {
		t.Fatalf("OBV计算失败: %v", err)
	}
	if len(result) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际%d", len(closeVals), len(result))
	}
	if len(result) == 0 {
		t.Fatalf("OBV返回结果为空")
	}
	for i := range expOBV {
		if i >= len(result) {
			t.Errorf("OBV结果长度不足，无法对比第%d行", i)
			continue
		}
		actual := result[i]
		exp := expOBV[i]
		if (exp == 0 && actual != 0) || (exp != 0 && (actual < exp-0.05 || actual > exp+0.05)) {
			t.Errorf("OBV[%d] 期望%.2f, 实际%.2f", i, exp, actual)
		}
	}
}
