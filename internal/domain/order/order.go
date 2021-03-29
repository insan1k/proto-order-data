package order

import (
	"github.com/shopspring/decimal"
)

type Order struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
	Inf      Info
}

func NewOrderFloat64(price, quantity float64) (o Order) {
	o.Inf.init()
	o.Price = decimal.NewFromFloat(price)
	o.Quantity = decimal.NewFromFloat(quantity)
	return
}

func NewOrderFloat32(price, quantity float32) (o Order) {
	o.Inf.init()
	o.Price = decimal.NewFromFloat32(price)
	o.Quantity = decimal.NewFromFloat32(quantity)
	return
}

func NewOrderString(price, quantity string) (o Order) {
	o.Inf.init()
	o.Price = decimal.RequireFromString(price)
	o.Quantity = decimal.RequireFromString(quantity)
	return
}

// EmptyOrder convenience function to return a dummy/empty orders
func EmptyOrder() (o Order) {
	o.Inf.init()
	o.Inf.SetTags(Empty)
	return
}
