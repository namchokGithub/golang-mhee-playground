package vat

import "errors"

const VAT_RATE = 0.0654205607

type Calculation struct {
	amount float64
	rate   float64
}

func NewCalculation(amount, rate float64) (*Calculation, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if rate <= 0 {
		rate = 7 // default VAT
	}
	return &Calculation{
		amount: amount,
		rate:   rate,
	}, nil
}

func (c *Calculation) CONST_VAT() float64 {
	return c.amount * VAT_RATE
}

func (c *Calculation) VAT() float64 {
	return c.amount * (c.rate / 100)
}

func (c *Calculation) Total() float64 {
	return c.amount + c.VAT()
}
