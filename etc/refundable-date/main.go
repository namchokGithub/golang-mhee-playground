package main

import (
	"fmt"
	"testing"
	"time"
)

func main() {
	type tc struct {
		name  string
		today time.Time
		date  time.Time
		want  bool
	}

	tests := []tc{
		// ====== กันยายน 2025 (ตัวอย่างอ้างอิง) ======
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

		// ====== เดือนอื่น ๆ / ปลายเดือน ======
		{"Oct-05 Sun: crosses Sep->Oct, in-week", L(2025, 10, 5), L(2025, 9, 30), true}, // สัปดาห์เริ่ม Mon 29 Sep
		{"Oct-05 Sun: day before weekStart not ref", L(2025, 10, 5), L(2025, 9, 28), false},

		{"Oct-27 Mon: week start refundable", L(2025, 10, 27), L(2025, 10, 27), true},
		{"Oct-27 Mon: previous day (Sun) not ref", L(2025, 10, 27), L(2025, 10, 26), false},

		{"Jan-01 Fri 2025: week spans 2024->2025 (Mon in-week)", L(2025, 1, 1), L(2024, 12, 30), true}, // weekStart = Mon 30 Dec 2024
		{"Jan-03 Fri 2025: week spans 2024->2025 (Mon in-week)", L(2025, 1, 3), L(2024, 12, 30), true},
		{"Jan-03 Fri 2025: day before weekStart", L(2025, 1, 3), L(2024, 12, 29), false},
		{"Jan-03 Fri 2025: 2024-12-28 (Sat) before weekStart", L(2025, 1, 3), L(2024, 12, 28), false},

		{"Jan-06 Mon 2025: new week, Sun not ref", L(2025, 1, 6), L(2025, 1, 5), false},

		{"Mar-04 Tue 2025: weekStart Mar-03", L(2025, 3, 4), L(2025, 3, 3), true},
		{"Mar-09 Sun 2025: prev Sun not ref", L(2025, 3, 9), L(2025, 3, 2), false},

		{"May-31 Sat 2025: start May-26 refundable", L(2025, 5, 31), L(2025, 5, 26), true},
		{"Jun-02 Mon 2025: prior Sun not refundable", L(2025, 6, 2), L(2025, 6, 1), false},

		// ====== ปีอธิกสุรทิน 2024 ======
		{"Feb-04 Sun 2024: weekStart Jan-29", L(2024, 2, 4), L(2024, 1, 31), true},
		{"Feb-04 Sun 2024: day before weekStart", L(2024, 2, 4), L(2024, 1, 28), false},

		{"Feb-29 Thu 2024: leap day inside week", L(2024, 2, 29), L(2024, 2, 28), true},
		{"Feb-29 Thu 2024: before weekStart not ref", L(2024, 2, 29), L(2024, 2, 25), false},
	}

	for _, tt := range tests {
		got := CheckRefundable(tt.date, tt.today)
		if got != tt.want {
			fmt.Printf("Today is %v, ❌ %s: got=%v, want=%v\n", tt.today.Format("2006-01-02"), tt.name, got, tt.want)
		} else {
			var result string
			if tt.want {
				result = "refundable"
			} else {
				result = "not refundable"
			}
			fmt.Printf("Today is %v, created on %v: %s\n", tt.today.Format("2006-01-02"), tt.date.Format("2006-01-02"), result)
			// fmt.Printf("Today is %v, %s\n", tt.today.Format("2006-01-02"), tt.name)
		}
	}
}

// ช่วยสร้างเวลาแบบ local ให้คงที่ (ตั้งเวลา 12:00 เพื่อตัดปัญหาเขตเวลา)
func L(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 12, 0, 0, 0, time.Local)
}

func TestCheckRefundableXXXX(t *testing.T) {
	tests := []struct {
		name  string
		today time.Time
		date  time.Time
		want  bool
	}{
		// ====== กันยายน 2025 (ตัวอย่างอ้างอิง) ======
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

		// ====== เดือนอื่น ๆ / ปลายเดือน ======
		{"Oct-05 Sun: crosses Sep->Oct, in-week", L(2025, 10, 5), L(2025, 9, 30), true}, // สัปดาห์เริ่ม Mon 29 Sep
		{"Oct-05 Sun: day before weekStart not ref", L(2025, 10, 5), L(2025, 9, 28), false},

		{"Oct-27 Mon: week start refundable", L(2025, 10, 27), L(2025, 10, 27), true},
		{"Oct-27 Mon: previous day (Sun) not ref", L(2025, 10, 27), L(2025, 10, 26), false},

		{"Jan-01 Fri 2025: week spans 2024->2025 (Mon in-week)", L(2025, 1, 1), L(2024, 12, 30), true}, // weekStart = Mon 30 Dec 2024
		{"Jan-03 Fri 2025: week spans 2024->2025 (Mon in-week)", L(2025, 1, 3), L(2024, 12, 30), true},
		{"Jan-03 Fri 2025: day before weekStart", L(2025, 1, 3), L(2024, 12, 29), false},
		{"Jan-03 Fri 2025: 2024-12-28 (Sat) before weekStart", L(2025, 1, 3), L(2024, 12, 28), false},

		{"Jan-06 Mon 2025: new week, Sun not ref", L(2025, 1, 6), L(2025, 1, 5), false},

		{"Mar-04 Tue 2025: weekStart Mar-03", L(2025, 3, 4), L(2025, 3, 3), true},
		{"Mar-09 Sun 2025: prev Sun not ref", L(2025, 3, 9), L(2025, 3, 2), false},

		{"May-31 Sat 2025: start May-26 refundable", L(2025, 5, 31), L(2025, 5, 26), true},
		{"Jun-02 Mon 2025: prior Sun not refundable", L(2025, 6, 2), L(2025, 6, 1), false},

		// ====== ปีอธิกสุรทิน 2024 ======
		{"Feb-04 Sun 2024: weekStart Jan-29", L(2024, 2, 4), L(2024, 1, 31), true},
		{"Feb-04 Sun 2024: day before weekStart", L(2024, 2, 4), L(2024, 1, 28), false},

		{"Feb-29 Thu 2024: leap day inside week", L(2024, 2, 29), L(2024, 2, 28), true},
		{"Feb-29 Thu 2024: before weekStart not ref", L(2024, 2, 29), L(2024, 2, 25), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CheckRefundable(tc.date, tc.today)
			if got != tc.want {
				t.Fatalf("CheckRefundable(%v, today=%v) = %v, want %v", tc.date, tc.today, got, tc.want)
			}
		})
	}
}

// หาวันเริ่มต้นสัปดาห์ (จันทร์ 00:00)
func startOfWeek(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()
	// reset เป็น 00:00
	midnight := time.Date(y, m, d, 0, 0, 0, 0, loc)

	// Go: Sunday=0, Monday=1,... Saturday=6
	wd := int(midnight.Weekday())
	// เราต้องการสัปดาห์เริ่มวันจันทร์
	diff := (wd - int(time.Monday) + 7) % 7

	return midnight.AddDate(0, 0, -diff)
}

// คืน true ถ้าวันที่ refund ได้
func CheckRefundable(date, today time.Time) bool {
	loc := today.Location()
	date = date.In(loc)

	// หาวันเริ่มสัปดาห์ของ today
	weekStart := startOfWeek(today)
	if date.Before(weekStart) {
		return false
	}

	// สัปดาห์จันทร์–อาทิตย์
	weekEnd := weekStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return !date.After(weekEnd)
}
