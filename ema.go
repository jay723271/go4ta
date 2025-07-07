package go4ta

func EMA(close []float64, timePeriod int) ([]float64, error) {
	return MA(close, timePeriod, 1)
}
