package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestPPO(t *testing.T) {
	file, err := os.Open("test_data/ppo.csv")
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
	var expPPOIdx []int
	var expPPO, expSignal, expHist []float64
	for i, record := range records {
		if len(record) < 5 {
			continue
		}
		// Close列在record[2]
		closeVal, err := strconv.ParseFloat(record[2], 64)
		if err == nil {
			closeVals = append(closeVals, closeVal)
		}
		// PPO在record[4]，PPO_Signal在record[5]，PPO_Hist在record[6]
		if record[4] != "" {
			ppoVal, err1 := strconv.ParseFloat(record[4], 64)
			signalVal, err2 := strconv.ParseFloat(record[5], 64)
			histVal, err3 := strconv.ParseFloat(record[6], 64)
			if err1 == nil && err2 == nil && err3 == nil {
				expPPOIdx = append(expPPOIdx, i)
				expPPO = append(expPPO, ppoVal)
				expSignal = append(expSignal, signalVal)
				expHist = append(expHist, histVal)
			}
		}
	}

	fastPeriod, slowPeriod, signalPeriod, maType := 12, 26, 9, 1 // 1=EMA
	ppo, signal, hist, err := PPOWithSignal(closeVals, fastPeriod, slowPeriod, signalPeriod, maType)
	if err != nil {
		t.Fatalf("PPOWithSignal计算失败: %v", err)
	}
	if len(ppo) != len(closeVals) || len(signal) != len(closeVals) || len(hist) != len(closeVals) {
		t.Errorf("期望结果长度%d，实际ppo=%d, signal=%d, hist=%d", len(closeVals), len(ppo), len(signal), len(hist))
	}

	if len(ppo) == 0 {
		t.Fatalf("PPO返回结果为空")
	}
	for i, idx := range expPPOIdx {
		if idx >= len(ppo) {
			t.Errorf("PPO结果长度不足，无法对比第%d行", idx)
			continue
		}
		actualPPO := ppo[idx]
		exp := expPPO[i]
		if (exp == 0 && actualPPO != 0) || (exp != 0 && (actualPPO < exp-0.05 || actualPPO > exp+0.05)) {
			t.Errorf("PPO[%d] 期望%.2f, 实际%.2f", idx, exp, actualPPO)
		}
		actualSignal := signal[idx]
		expS := expSignal[i]
		if (expS == 0 && actualSignal != 0) || (expS != 0 && (actualSignal < expS-0.05 || actualSignal > expS+0.05)) {
			t.Errorf("Signal[%d] 期望%.2f, 实际%.2f", idx, expS, actualSignal)
		}
		actualHist := hist[idx]
		expH := expHist[i]
		if (expH == 0 && actualHist != 0) || (expH != 0 && (actualHist < expH-0.05 || actualHist > expH+0.05)) {
			t.Errorf("Hist[%d] 期望%.2f, 实际%.2f", idx, expH, actualHist)
		}
	}

	// 边界与异常测试
	_, _, _, err = PPOWithSignal([]float64{}, fastPeriod, slowPeriod, signalPeriod, maType)
	if err != nil && err.Error() != "input data length (0) is too small for the given periods" {
		t.Errorf("空输入应返回无错或特定错误，实际: %v", err)
	}
	_, _, _, err = PPOWithSignal([]float64{1, 2, 3}, fastPeriod, slowPeriod, signalPeriod, maType)
	if err == nil {
		t.Errorf("输入长度不足应返回错误")
	}
}
