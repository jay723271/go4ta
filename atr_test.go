package go4ta

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestATR(t *testing.T) {
	// Load test data from CSV file
	t.Run("Test with sample data", func(t *testing.T) {
		// Read test data from CSV
		file, err := os.Open("test_data/atr.csv")
		if err != nil {
			t.Fatalf("Failed to open test data file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read CSV data: %v", err)
		}

		// Skip header row
		records = records[1:]

		// Prepare slices for high, low, close prices and expected ATR values
		var high, low, closeP, expectedATR []float64

		for _, record := range records {
			h, _ := strconv.ParseFloat(record[0], 64)
			l, _ := strconv.ParseFloat(record[1], 64)
			c, _ := strconv.ParseFloat(record[2], 64)

			high = append(high, h)
			low = append(low, l)
			closeP = append(closeP, c)

			// Some ATR values might be empty in the test data
			if len(record) > 3 && record[3] != "" {
				atr, _ := strconv.ParseFloat(record[3], 64)
				expectedATR = append(expectedATR, atr)
			}
		}

		// Test with different time periods
		timePeriods := []int{14, 5, 10}
		for _, period := range timePeriods {
			t.Run(fmt.Sprintf("Period_%d", period), func(t *testing.T) {
				result, err := ATR(high, low, closeP, period)
				if err != nil {
					t.Fatalf("ATR calculation failed: %v", err)
				}

				// Check that the result length matches input length
				if len(result) != len(high) {
					t.Errorf("Expected result length %d, got %d", len(high), len(result))
				}

				// For the first (period-1) elements, the ATR should be 0
				for i := 0; i < period-1; i++ {
					if result[i] != 0 {
						t.Errorf("Expected ATR[%d] to be 0, got %f", i, result[i])
					}
				}

				// For period=14, we can verify against known values
				if period == 14 && len(expectedATR) > 0 {
					// We'll check the last few values as they're more stable
					startIdx := len(high) - len(expectedATR)
					for i := 0; i < len(expectedATR); i++ {
						idx := startIdx + i
						tolerance := 0.01 // 0.01 is a reasonable tolerance for floating point comparison
						if idx >= 0 && idx < len(result) {
							diff := result[idx] - expectedATR[i]
							if diff < -tolerance || diff > tolerance {
								t.Errorf("ATR[%d] expected ~%f, got %f (diff: %f)", 
									idx, expectedATR[i], result[idx], diff)
							}
						}
					}
				}
			})
		}
	})

	// Test edge cases
	t.Run("Test edge cases", func(t *testing.T) {
		// Test empty input
		t.Run("Empty input", func(t *testing.T) {
			result, err := ATR([]float64{}, []float64{}, []float64{}, 14)
			if err != nil {
				t.Fatalf("Expected no error for empty input, got %v", err)
			}
			if len(result) != 0 {
				t.Errorf("Expected empty result for empty input, got length %d", len(result))
			}
		})

		// Test input shorter than time period
		t.Run("Input shorter than time period", func(t *testing.T) {
			high := []float64{1.0, 2.0, 3.0}
			low := []float64{0.5, 1.5, 2.5}
			closeP := []float64{0.8, 1.8, 2.8}

			_, err := ATR(high, low, closeP, 5)
			if err == nil {
				t.Error("Expected error for input shorter than time period, got nil")
			}
		})

		// Test mismatched input lengths
		t.Run("Mismatched input lengths", func(t *testing.T) {
			high := []float64{1.0, 2.0, 3.0}
			low := []float64{0.5, 1.5}
			closeP := []float64{0.8, 1.8, 2.8}

			_, err := ATR(high, low, closeP, 2)
			if err == nil {
				t.Error("Expected error for mismatched input lengths, got nil")
			}
		})
	})
}
