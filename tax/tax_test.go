package tax

import (
	"fmt"
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	t.Parallel()

	okTests := []struct {
		amount   amountExclConsumpTax
		date     consumpDate
		expected amountInclConsumpTax
	}{
		{
			amount:   100,
			date:     consumpDate(time.Date(1989, 4, 1, 0, 0, 0, 0, time.Local)),
			expected: 103,
		},
		{
			amount:   100,
			date:     consumpDate(time.Date(1997, 3, 31, 0, 0, 0, 0, time.Local)),
			expected: 103,
		},
		{
			amount:   100,
			date:     consumpDate(time.Date(1997, 4, 1, 0, 0, 0, 0, time.Local)),
			expected: 105,
		},
		{
			amount:   100,
			date:     consumpDate(time.Date(2014, 3, 31, 0, 0, 0, 0, time.Local)),
			expected: 105,
		},
		{
			amount:   100,
			date:     consumpDate(time.Date(2014, 4, 1, 0, 0, 0, 0, time.Local)),
			expected: 108,
		},
		{
			amount:   100,
			date:     consumpDate(time.Date(2019, 3, 31, 0, 0, 0, 0, time.Local)),
			expected: 108,
		},
		{
			amount:   100,
			date:     consumpDate(time.Date(2019, 4, 1, 0, 0, 0, 0, time.Local)),
			expected: 110,
		},
		{
			amount:   100,
			date:     consumpDate{},
			expected: 110,
		},
		{
			amount:   100,
			date:     consumpDate(time.Time{}),
			expected: 110,
		},
		{
			amount:   100,
			date:     consumpDate(time.Now()),
			expected: 110,
		},
	}

	for _, test := range okTests {
		test := test

		t.Run(fmt.Sprintf("ok:amount,date:%v,%v", test.amount, test.date), func(t *testing.T) {
			t.Parallel()

			actual, err := consume(test.amount, test.date)
			if err != nil {
				t.Errorf("failed to consume: %v", err)
			}
			if err := actual.validate(); err != nil {
				t.Errorf("the amount should be valid: %s", err)
			}
			if actual != test.expected {
				t.Errorf("got %v, but expected %v", actual, test.expected)
			}
		})
	}

	failedTests := []struct {
		amount amountExclConsumpTax
		date   consumpDate
		reason string
	}{
		{
			amount: 100,
			date:   consumpDate(time.Date(1989, 3, 31, 0, 0, 0, 0, time.Local)),
			reason: "no consumption tax to apply",
		},
	}

	for _, test := range failedTests {
		test := test

		t.Run(fmt.Sprintf("failed:amount,date:%v,%v", test.amount, test.date), func(t *testing.T) {
			t.Parallel()

			if _, err := consume(test.amount, test.date); err == nil {
				t.Errorf("should fail to consume: %v", test.reason)
			}
		})
	}
}

func TestConsumpTax(t *testing.T) {
	t.Parallel()

	okTests := []consumpTax{
		{
			rate:       0.1,
			enforcedAt: time.Date(2019, 4, 1, 0, 0, 0, 0, time.Local),
		},
		{
			rate:       0.08,
			enforcedAt: time.Date(2014, 4, 1, 0, 0, 0, 0, time.Local),
		},
		{
			rate:       0.05,
			enforcedAt: time.Date(1997, 4, 1, 0, 0, 0, 0, time.Local),
		},
		{
			rate:       0.03,
			enforcedAt: time.Date(1989, 4, 1, 0, 0, 0, 0, time.Local),
		},
	}

	for _, test := range okTests {
		test := test

		t.Run(fmt.Sprintf("ok:%v", test), func(t *testing.T) {
			t.Parallel()

			if err := test.validate(); err != nil {
				t.Errorf("the tax should be valid: %v", err)
			}
		})
	}

	failedTests := []struct {
		tax    consumpTax
		reason string
	}{
		{
			tax: consumpTax{
				rate:       -0.05,
				enforcedAt: time.Date(2019, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "negative rate",
		},
		{
			tax: consumpTax{
				rate:       0.05,
				enforcedAt: time.Date(2019, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the rate enforced at 2019/04/01",
		},
		{
			tax: consumpTax{
				rate:       0.1,
				enforcedAt: time.Date(2014, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the enforcement date of the rate 0.1",
		},
		{
			tax: consumpTax{
				rate:       0.06,
				enforcedAt: time.Date(2014, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the rate enforced at 2014/04/01",
		},
		{
			tax: consumpTax{
				rate:       0.08,
				enforcedAt: time.Date(1997, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the enforcement date of the rate 0.08",
		},
		{
			tax: consumpTax{
				rate:       0.04,
				enforcedAt: time.Date(1997, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the rate enforced at 1997/04/01",
		},
		{
			tax: consumpTax{
				rate:       0.05,
				enforcedAt: time.Date(1989, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the enforcement date of the rate 0.05",
		},
		{
			tax: consumpTax{
				rate:       0.01,
				enforcedAt: time.Date(1989, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the rate enforced at 1989/04/01",
		},
		{
			tax: consumpTax{
				rate:       0.03,
				enforcedAt: time.Date(2019, 4, 1, 0, 0, 0, 0, time.Local),
			},
			reason: "does not match the enforcement date of the rate 0.03",
		},
	}

	for _, test := range failedTests {
		test := test

		t.Run(fmt.Sprintf("failed:%v", test.tax), func(t *testing.T) {
			t.Parallel()

			if err := test.tax.validate(); err == nil {
				t.Errorf("the tax should not be valid: %v", test.reason)
			}
		})
	}
}
