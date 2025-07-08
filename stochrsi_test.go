package go4ta

import (
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"testing"
)

func TestSTOCHRSI(t *testing.T) {
	// 打开测试数据文件
	file, err := os.Open("test_data/stochrsi.csv")
	if err != nil {
		t.Fatalf("无法打开测试文件: %v", err)
	}
	defer file.Close()

	// 读取CSV数据
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	// 解析测试数据
	var closePrices, expectFastK, expectFastD []float64
	for i, record := range records {
		if i == 0 {
			continue // 跳过表头
		}

		// 解析收盘价
		close, _ := strconv.ParseFloat(record[0], 64)
		closePrices = append(closePrices, close)

		// 解析期望的FastK和FastD
		if record[1] != "" {
			fk, _ := strconv.ParseFloat(record[1], 64)
			expectFastK = append(expectFastK, fk)
		} else {
			expectFastK = append(expectFastK, -1)
		}

		if len(record) > 2 && record[2] != "" {
			fd, _ := strconv.ParseFloat(record[2], 64)
			expectFastD = append(expectFastD, fd)
		} else {
			expectFastD = append(expectFastD, -1)
		}
	}

	// 计算STOCHRSI
	fastK, fastD, err := STOCHRSI(closePrices, 14, 5, 3, 0)
	if err != nil {
		t.Fatalf("STOCHRSI计算失败: %v", err)
	}

	// 设置误差容忍度
	const epsilon = 0.3 // 允许的最大绝对误差为0.3

	// 检查FastK
	for i, exp := range expectFastK {
		if exp < 0 {
			continue // 跳过无效数据
		}
		if diff := math.Abs(fastK[i] - exp); diff > epsilon {
			t.Errorf("FastK[%d] 差异过大: 期望 %.6f, 实际 %.6f, 差异 %.6f", 
				i, exp, fastK[i], diff)
		}
	}

	// 检查FastD
	for i, exp := range expectFastD {
		if exp < 0 {
			continue // 跳过无效数据
		}
		if diff := math.Abs(fastD[i] - exp); diff > epsilon {
			t.Errorf("FastD[%d] 差异过大: 期望 %.6f, 实际 %.6f, 差异 %.6f", 
				i, exp, fastD[i], diff)
		}
	}
}