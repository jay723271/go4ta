package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestSTDDEV(t *testing.T) {
	file, err := os.Open("test_data/stddev.csv")
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
	var expSTDDEV []float64
	for _, record := range records {
		if len(record) < 1 {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[0], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		if len(record) > 1 && record[1] != "" {
			stdVal, err := strconv.ParseFloat(record[1], 64)
			if err == nil {
				expSTDDEV = append(expSTDDEV, stdVal)
			} else {
				expSTDDEV = append(expSTDDEV, 0)
			}
		} else {
			expSTDDEV = append(expSTDDEV, 0)
		}
	}

	timePeriod := 20
	nbDev := 1.0
	result, err := STDDEV(closeVals, timePeriod, nbDev)
	if err != nil {
		t.Fatalf("STDDEV计算失败: %v", err)
	}
	if len(result) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际%d", len(closeVals), len(result))
	}

	if len(result) == 0 {
		t.Fatalf("STDDEV返回结果为空")
	}
	for i := 0; i < len(result) && i < len(expSTDDEV); i++ {
		actual := result[i]
		exp := expSTDDEV[i]
		if (exp == 0 && actual != 0) || (exp != 0 && (actual < exp-0.05 || actual > exp+0.05)) {
			t.Errorf("STDDEV[%d] 期望%.4f, 实际%.4f (Close=%.2f)", i, exp, actual, closeVals[i])
		}
	}

	// 边界与异常测试
	t.Run("空输入", func(t *testing.T) {
		res, err := STDDEV([]float64{}, timePeriod, nbDev)
		if err != nil {
			t.Errorf("空输入应无错误，实际: %v", err)
		}
		if len(res) != 0 {
			t.Errorf("空输入应返回空切片，实际长度: %d", len(res))
		}
	})

	t.Run("长度不足", func(t *testing.T) {
		short := make([]float64, timePeriod-1)
		_, err := STDDEV(short, timePeriod, nbDev)
		if err == nil {
			t.Errorf("长度不足应返回错误")
		}
	})
}
