package vat

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Calc(amount, rate float64) (vat, total float64) {
	calc, err := NewCalculation(amount, rate)
	if err != nil {
		return 0, 0
	}

	vat = calc.VAT()
	total = amount + vat
	return
}
