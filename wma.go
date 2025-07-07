package go4ta

func WMA(close []float64, timePeriod int) ([]float64, error) {
	return MA(close, timePeriod, 2)
}
