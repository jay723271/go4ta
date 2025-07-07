package go4ta

func SMA(close []float64, timePeriod int) ([]float64, error) {
	return MA(close, timePeriod, 0)
}
