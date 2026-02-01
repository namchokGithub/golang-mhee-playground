package refundable_date

import "time"

type Service struct{}

func NewService() *Service { return &Service{} }

// startOfWeek returns Monday 00:00 of the week for the given time.
func startOfWeek(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, loc)

	wd := int(midnight.Weekday())
	diff := (wd - int(time.Monday) + 7) % 7
	return midnight.AddDate(0, 0, -diff)
}

// CheckRefundable returns true if date is within Monday-Sunday of today's week (inclusive).
func (s *Service) CheckRefundable(date, today time.Time) bool {
	loc := today.Location()
	date = date.In(loc)

	weekStart := startOfWeek(today)
	if date.Before(weekStart) {
		return false
	}

	weekEnd := weekStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return !date.After(weekEnd)
}
