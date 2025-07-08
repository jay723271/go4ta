package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestSTOCH(t *testing.T) {
	file, err := os.Open("test_data/stoch.csv")
	if err != nil {
		t.Fatalf("无法打开CSV文件: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("无法读取CSV文件: %v", err)
	}

	var high, low, close []float64
	var expectK, expectD []float64

	for i, rec := range records {
		if i == 0 {
			continue // 跳过表头
		}
		h, _ := strconv.ParseFloat(rec[0], 64)
		l, _ := strconv.ParseFloat(rec[1], 64)
		c, _ := strconv.ParseFloat(rec[2], 64)
		high = append(high, h)
		low = append(low, l)
		close = append(close, c)

		// 处理SlowK
		if rec[3] != "" {
			v, _ := strconv.ParseFloat(rec[3], 64)
			expectK = append(expectK, v)
		} else {
			expectK = append(expectK, -1)
		}
		// 处理SlowD
		if rec[4] != "" {
			v, _ := strconv.ParseFloat(rec[4], 64)
			expectD = append(expectD, v)
		} else {
			expectD = append(expectD, -1)
		}
	}

	// 参数与csv一致: (5,3,3)
	k, d, err := STOCH(high, low, close, 5, 3, 3, 0, 0)
	if err != nil {
		t.Fatalf("STOCH计算失败: %v", err)
	}

	// 允许误差
	eps := 0.1
	for i := range expectK {
		if expectK[i] >= 0 {
			if abs(k[i]-expectK[i]) > eps {
				t.Errorf("SlowK[%d] 期望: %.2f, 实际: %.2f", i, expectK[i], k[i])
			}
		}
		if expectD[i] >= 0 {
			if abs(d[i]-expectD[i]) > eps {
				t.Errorf("SlowD[%d] 期望: %.2f, 实际: %.2f", i, expectD[i], d[i])
			}
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
