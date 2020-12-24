package tax

import (
	"fmt"
	"sort"
	"time"
)

var consumpTaxes = map[float64]consumpTax{
	0.1: {
		rate:       0.1,
		enforcedAt: time.Date(2019, 4, 1, 0, 0, 0, 0, time.Local),
	},
	0.08: {
		rate:       0.08,
		enforcedAt: time.Date(2014, 4, 1, 0, 0, 0, 0, time.Local),
	},
	0.05: {
		rate:       0.05,
		enforcedAt: time.Date(1997, 4, 1, 0, 0, 0, 0, time.Local),
	},
	0.03: {
		rate:       0.03,
		enforcedAt: time.Date(1989, 4, 1, 0, 0, 0, 0, time.Local),
	},
}

var sortedConsumpRates = func() []float64 {
	rates := make([]float64, 0, len(consumpTaxes))
	for rate := range consumpTaxes {
		rates = append(rates, rate)
	}

	sort.Float64s(rates)
	return rates
}()

func main() {}

func consume(amount amountExclConsumpTax, date consumpDate) (amountInclConsumpTax, error) {
	date.initIfZero()

	if err := validate(amount, date); err != nil {
		return 0, err
	}

	tax, err := consumpTaxForDate(date)
	if err != nil {
		return 0, err
	}

	included := amountInclConsumpTax(float64(amount) + float64(amount)*float64(tax.rate))
	if err := validate(included); err != nil {
		return 0, err
	}
	return included, nil
}

type amountInclConsumpTax uint

func (a amountInclConsumpTax) validate() error {
	return nil
}

type amountExclConsumpTax uint

func (a amountExclConsumpTax) validate() error {
	return nil
}

func consumpTaxForDate(date consumpDate) (consumpTax, error) {
	date.initIfZero()

	if err := validate(date); err != nil {
		return consumpTax{}, err
	}

	for i := len(sortedConsumpRates) - 1; i >= 0; i-- {
		tax := consumpTaxes[sortedConsumpRates[i]]
		if date.unwrap().Equal(tax.enforcedAt) || date.unwrap().After(tax.enforcedAt) {
			return tax, nil
		}
	}

	return consumpTax{}, fmt.Errorf("no consumption tax for contract date")
}

type consumpTax struct {
	rate       float64
	enforcedAt time.Time
}

func (t consumpTax) validate() error {
	if t.rate < 0 {
		return fmt.Errorf("rate should be equal or more than 0")
	}

	validTax, ok := consumpTaxes[t.rate]
	if !ok {
		return fmt.Errorf("invalid rate")
	}
	if !t.enforcedAt.Equal(validTax.enforcedAt) {
		return fmt.Errorf("invalid enforcement date")
	}

	return nil
}

type consumpDate time.Time

func (d consumpDate) String() string {
	return d.unwrap().String()
}

func (d *consumpDate) initIfZero() {
	if d.unwrap().IsZero() {
		*d = consumpDate(time.Now())
	}
}

func (d consumpDate) validate() error {
	if d.unwrap().IsZero() {
		return fmt.Errorf("should be initialized")
	}
	return nil
}

func (d consumpDate) unwrap() time.Time {
	return time.Time(d)
}

func validate(vs ...validator) error {
	for _, v := range vs {
		if err := v.validate(); err != nil {
			return err
		}
	}
	return nil
}

type validator interface {
	validate() error
}
