package refundable_date

import "time"

type Checker interface {
	CheckRefundable(date, today time.Time) bool
}
