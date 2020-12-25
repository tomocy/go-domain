package discount

import (
	"fmt"
	"testing"
)

func TestDiscountAsRegular(t *testing.T) {
	t.Parallel()

	tests := []struct {
		price    price
		expected discPrice
	}{
		{
			price:    1000,
			expected: 600,
		},
		{
			price:    100,
			expected: 60,
		},
		{
			price:    99,
			expected: 60,
		},
		{
			price:    10,
			expected: 6,
		},
		{
			price:    2,
			expected: 2,
		},
		{
			price:    1,
			expected: 1,
		},
		{
			price:    0,
			expected: 0,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprint(test.price), func(t *testing.T) {
			t.Parallel()

			actual := discountAsRegular(test.price)
			if actual != test.expected {
				t.Errorf("got %v, but expected %v", actual, test.expected)
			}
		})
	}
}

func TestDiscountInSummer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		price    price
		expected discPrice
	}{
		{
			price:    1000,
			expected: 700,
		},
		{
			price:    800,
			expected: 500,
		},
		{
			price:    700,
			expected: 500,
		},
		{
			price:    500,
			expected: 500,
		},
		{
			price:    499,
			expected: 499,
		},
		{
			price:    100,
			expected: 100,
		},
		{
			price:    99,
			expected: 99,
		},
		{
			price:    10,
			expected: 10,
		},
		{
			price:    2,
			expected: 2,
		},
		{
			price:    1,
			expected: 1,
		},
		{
			price:    0,
			expected: 0,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprint(test.price), func(t *testing.T) {
			t.Parallel()

			actual := discountInSummer(test.price)
			if actual != test.expected {
				t.Errorf("got %v, but expected %v", actual, test.expected)
			}
		})
	}
}
