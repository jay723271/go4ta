package go4ta

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestADX(t *testing.T) {
	t.Run("Test with sample data", func(t *testing.T) {
		file, err := os.Open("test_data/adx.csv")
		if err != nil {
			t.Fatalf("Failed to open test data file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV data: %v", err)
		}

		// 跳过表头
		records = records[1:]

		var high, low, closeP, expectedADX []float64
		for _, record := range records {
			h, _ := strconv.ParseFloat(record[0], 64)
			l, _ := strconv.ParseFloat(record[1], 64)
			c, _ := strconv.ParseFloat(record[2], 64)
			high = append(high, h)
			low = append(low, l)
			closeP = append(closeP, c)
			// ADX在第3列（下标3），有些为空
			if len(record) > 3 && record[3] != "" {
				a, _ := strconv.ParseFloat(record[3], 64)
				expectedADX = append(expectedADX, a)
			}
		}

		timePeriod := 14
		result, err := ADX(high, low, closeP, timePeriod)
		if err != nil {
			t.Fatalf("ADX calculation failed: %v", err)
		}
		if len(result) != len(high) {
			t.Errorf("Expected result length %d, got %d", len(high), len(result))
		}
		// 前timePeriod-1个应为0
		for i := 0; i < timePeriod-1; i++ {
			if result[i] != 0 {
				t.Errorf("Expected ADX[%d] to be 0, got %f", i, result[i])
			}
		}
		// 检查与CSV的ADX列对齐的部分
		startIdx := len(high) - len(expectedADX)
		tolerance := 0.02 // 浮点误差容忍
		for i := 0; i < len(expectedADX); i++ {
			idx := startIdx + i
			if idx >= 0 && idx < len(result) {
				diff := result[idx] - expectedADX[i]
				if diff < -tolerance || diff > tolerance {
					t.Errorf("ADX[%d] expected ~%f, got %f (diff: %f)", idx, expectedADX[i], result[idx], diff)
				}
			}
		}
	})

	t.Run("Test edge cases", func(t *testing.T) {
		t.Run("Empty input", func(t *testing.T) {
			result, err := ADX([]float64{}, []float64{}, []float64{}, 14)
			if err != nil {
				t.Fatalf("Expected no error for empty input, got %v", err)
			}
			if len(result) != 0 {
				t.Errorf("Expected empty result for empty input, got length %d", len(result))
			}
		})
		t.Run("Input shorter than time period", func(t *testing.T) {
			high := []float64{1.0, 2.0, 3.0}
			low := []float64{0.5, 1.5, 2.5}
			closeP := []float64{0.8, 1.8, 2.8}
			_, err := ADX(high, low, closeP, 5)
			if err == nil {
				t.Error("Expected error for input shorter than time period, got nil")
			}
		})
		t.Run("Mismatched input lengths", func(t *testing.T) {
			high := []float64{1.0, 2.0, 3.0}
			low := []float64{0.5, 1.5}
			closeP := []float64{0.8, 1.8, 2.8}
			_, err := ADX(high, low, closeP, 2)
			if err == nil {
				t.Error("Expected error for mismatched input lengths, got nil")
			}
		})
	})
}
