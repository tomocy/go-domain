package discount

const (
	regularDiscRate = 0.4

	minSummerDiccPrice = 500
	summerDiscAmount   = 300
)

func discountAsRegular(p price) discPrice {
	d := price(float64(p) * regularDiscRate)
	return discPrice(p.sub(d))
}

func discountInSummer(p price) discPrice {
	if p < minSummerDiccPrice {
		return discPrice(p)
	}

	discounted := p.sub(300)
	if discounted < minSummerDiccPrice {
		discounted = minSummerDiccPrice
	}
	return discPrice(discounted)
}

type discPrice uint

type price uint

func (p price) sub(other price) price {
	return p - other
}
