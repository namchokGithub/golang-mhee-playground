package refundable_date_testing

import (
	"testing"
	"time"

	"proundmhee/internal/modules/refundable_date"

	"github.com/stretchr/testify/require"
)

func L(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 12, 0, 0, 0, time.Local)
}

func TestService_CheckRefundable(t *testing.T) {
	tests := []struct {
		name  string
		today time.Time
		date  time.Time
		want  bool
	}{
		{"Aug-12 Tue: refundable", L(2025, 8, 10), L(2025, 8, 12), false},
		{"Aug-11 Mon: refundable", L(2025, 8, 11), L(2025, 8, 12), true},
		{"Sep-08 Tue: refundable", L(2025, 9, 8), L(2025, 9, 8), true},
		{"Sep-02 Tue: same day refundable", L(2025, 9, 2), L(2025, 9, 2), true},
		{"Sep-02 Tue: before last Sun not refundable", L(2025, 9, 2), L(2025, 8, 31), false},
		{"Sep-07 Sun: Mon of week refundable", L(2025, 9, 7), L(2025, 9, 1), true},
		{"Sep-07 Sun: prev month not refundable", L(2025, 9, 7), L(2025, 8, 31), false},
		{"Sep-08 Mon: yesterday (Sun) not refundable", L(2025, 9, 8), L(2025, 9, 7), false},
		{"Sep-11 Thu: inside week refundable", L(2025, 9, 11), L(2025, 9, 8), true},
		{"Sep-11 Thu: before prior Sun not ref", L(2025, 9, 11), L(2025, 9, 7), false},
		{"Sep-14 Sun: inside week refundable", L(2025, 9, 14), L(2025, 9, 12), true},
		{"Sep-16 Tue: same day refundable", L(2025, 9, 16), L(2025, 9, 16), true},
		{"Sep-16 Tue: last Sun not refundable", L(2025, 9, 16), L(2025, 9, 14), false},

		{"Oct-05 Sun: crosses Sep->Oct, in-week", L(2025, 10, 5), L(2025, 9, 30), true},
		{"Oct-05 Sun: day before weekStart not ref", L(2025, 10, 5), L(2025, 9, 28), false},

		{"Oct-27 Mon: week start refundable", L(2025, 10, 27), L(2025, 10, 27), true},
		{"Oct-27 Mon: previous day (Sun) not ref", L(2025, 10, 27), L(2025, 10, 26), false},

		{"Jan-01 Fri 2025: week spans 2024->2025 (Mon in-week)", L(2025, 1, 1), L(2024, 12, 30), true},
		{"Jan-03 Fri 2025: week spans 2024->2025 (Mon in-week)", L(2025, 1, 3), L(2024, 12, 30), true},
		{"Jan-03 Fri 2025: day before weekStart", L(2025, 1, 3), L(2024, 12, 29), false},
		{"Jan-03 Fri 2025: 2024-12-28 (Sat) before weekStart", L(2025, 1, 3), L(2024, 12, 28), false},

		{"Jan-06 Mon 2025: new week, Sun not ref", L(2025, 1, 6), L(2025, 1, 5), false},

		{"Mar-04 Tue 2025: weekStart Mar-03", L(2025, 3, 4), L(2025, 3, 3), true},
		{"Mar-09 Sun 2025: prev Sun not ref", L(2025, 3, 9), L(2025, 3, 2), false},

		{"May-31 Sat 2025: start May-26 refundable", L(2025, 5, 31), L(2025, 5, 26), true},
		{"Jun-02 Mon 2025: prior Sun not refundable", L(2025, 6, 2), L(2025, 6, 1), false},

		{"Feb-04 Sun 2024: weekStart Jan-29", L(2024, 2, 4), L(2024, 1, 31), true},
		{"Feb-04 Sun 2024: day before weekStart", L(2024, 2, 4), L(2024, 1, 28), false},

		{"Feb-29 Thu 2024: leap day inside week", L(2024, 2, 29), L(2024, 2, 28), true},
		{"Feb-29 Thu 2024: before weekStart not ref", L(2024, 2, 29), L(2024, 2, 25), false},
	}

	svc := refundable_date.NewService()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := svc.CheckRefundable(tc.date, tc.today)
			require.Equal(t, tc.want, got)
		})
	}
}
