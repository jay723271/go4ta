package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestMACD(t *testing.T) {
	file, err := os.Open("test_data/macd.csv")
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

	var closeVals, expMACD, expSignal, expHist []float64
	for _, record := range records {
		if len(record) < 1 {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[0], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		if len(record) > 1 && record[1] != "" {
			macdVal, err := strconv.ParseFloat(record[1], 64)
			if err == nil {
				expMACD = append(expMACD, macdVal)
			}
		}
		if len(record) > 2 && record[2] != "" {
			signalVal, err := strconv.ParseFloat(record[2], 64)
			if err == nil {
				expSignal = append(expSignal, signalVal)
			}
		}
		if len(record) > 3 && record[3] != "" {
			histVal, err := strconv.ParseFloat(record[3], 64)
			if err == nil {
				expHist = append(expHist, histVal)
			}
		}
	}

	fastPeriod, slowPeriod, signalPeriod := 12, 26, 9
	macd, signal, hist, err := MACD(closeVals, fastPeriod, slowPeriod, signalPeriod)
	if err != nil {
		t.Fatalf("MACD计算失败: %v", err)
	}
	if len(macd) != len(closeVals) || len(signal) != len(closeVals) || len(hist) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际macd:%d signal:%d hist:%d", len(closeVals), len(macd), len(signal), len(hist))
	}

	// 只校验有期望值的部分
	startIdx := len(closeVals) - len(expMACD)
	for i := 0; i < len(expMACD); i++ {
		idx := startIdx + i
		if idx >= len(macd) {
			break
		}
		if diff := macd[idx] - expMACD[i]; diff < -0.05 || diff > 0.05 {
			t.Errorf("MACD[%d] 期望%.2f, 实际%.2f", idx, expMACD[i], macd[idx])
		}
		if diff := signal[idx] - expSignal[i]; diff < -0.05 || diff > 0.05 {
			t.Errorf("Signal[%d] 期望%.2f, 实际%.2f", idx, expSignal[i], signal[idx])
		}
		if diff := hist[idx] - expHist[i]; diff < -0.05 || diff > 0.05 {
			t.Errorf("Hist[%d] 期望%.2f, 实际%.2f", idx, expHist[i], hist[idx])
		}
	}
}
