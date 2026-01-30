package vat

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Calc(amount, rate float64) (vat, total float64) {
	vat = amount * (rate / 100)
	total = amount + vat
	return
}
