package recreation

import (
	"testing"
	"time"
)

type MonthStartTest struct {
	Date    time.Time
	Correct time.Time
}

var (
	monthStartTests = []MonthStartTest{
		{time.Date(2022, 2, 10, 4, 14, 5, 0, time.UTC), time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(2022, 2, 1, 4, 14, 5, 0, time.UTC), time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(2022, 3, 10, 4, 14, 5, 0, time.UTC), time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)},
	}
)

func TestGetStartOfMonth(t *testing.T) {

	for _, test := range monthStartTests {
		calculatedStart := GetStartOfMonth(test.Date)
		if test.Correct != calculatedStart {
			t.Log("got wrong start", test.Correct, test.Date, calculatedStart)
		}
	}

}
