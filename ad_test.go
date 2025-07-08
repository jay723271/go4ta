package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestAD(t *testing.T) {
	t.Run("Test with sample data", func(t *testing.T) {
		file, err := os.Open("test_data/ad.csv")
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

		var high, low, closeP, volume, expectedAD []float64
		for _, record := range records {
			if len(record) < 5 {
				continue
			}
			h, _ := strconv.ParseFloat(record[0], 64)
			l, _ := strconv.ParseFloat(record[1], 64)
			c, _ := strconv.ParseFloat(record[2], 64)
			v, _ := strconv.ParseFloat(record[3], 64)
			high = append(high, h)
			low = append(low, l)
			closeP = append(closeP, c)
			volume = append(volume, v)
			ad, _ := strconv.ParseFloat(record[4], 64)
			expectedAD = append(expectedAD, ad)
		}

		result, err := AD(high, low, closeP, volume)
		if err != nil {
			t.Fatalf("AD计算失败: %v", err)
		}
		if len(result) != len(expectedAD) {
			t.Errorf("期望结果长度%d, 实际%d", len(expectedAD), len(result))
		}
		tolerance := 1e-6 // 浮点误差容忍
		for i := 0; i < len(expectedAD) && i < len(result); i++ {
			if diff := result[i] - expectedAD[i]; diff < -tolerance || diff > tolerance {
				t.Errorf("AD[%d] 期望%.6f, 实际%.6f, diff=%.6f", i, expectedAD[i], result[i], diff)
			}
		}
	})
}
