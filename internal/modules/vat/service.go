package vat

func Calc(amount, rate float64) (vat float64, total float64) {
	vat = amount * (rate / 100.0)
	total = amount + vat
	return
}
