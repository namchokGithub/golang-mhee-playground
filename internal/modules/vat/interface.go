package vat

type Calculator interface {
	Calculate(amount, rate float64) (vat, total float64, err error)
}
