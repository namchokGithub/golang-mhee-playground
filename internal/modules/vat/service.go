package vat

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Calculate(amount, rate float64) (vat, total float64, err error) {
	calc, err := NewCalculation(amount, rate)
	if err != nil {
		return 0, 0, err
	}

	vat = calc.VAT()
	total = amount + vat
	return
}
