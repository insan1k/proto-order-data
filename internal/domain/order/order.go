package order

import (
	"github.com/shopspring/decimal"
)

type Order struct {
	Asset            string
	Price            decimal.Decimal
	Quantity         decimal.Decimal
	priceMulQuantity decimal.Decimal
	Inf              Info
}

func (o *Order) preprocess() {
	o.Inf.init()
	dec := decimal.New(0, 20)
	o.Price, _ = decimal.RescalePair(o.Price, dec)
	o.Quantity, _ = decimal.RescalePair(o.Quantity, dec)

}

func (o *Order) GetPriceMulQuantity() decimal.Decimal {
	return o.priceMulQuantity
}

func NewOrderFloat64(price, quantity float64) (o Order) {
	o.Price = decimal.NewFromFloat(price)
	o.Quantity = decimal.NewFromFloat(quantity)
	o.preprocess()
	return
}

func NewOrderFloat32(price, quantity float32) (o Order) {
	o.Price = decimal.NewFromFloat32(price)
	o.Quantity = decimal.NewFromFloat32(quantity)
	o.preprocess()
	return
}

func NewOrderString(price, quantity string) (o Order) {
	o.Price = decimal.RequireFromString(price)
	o.Quantity = decimal.RequireFromString(quantity)
	o.preprocess()
	return
}

// EmptyOrder convenience function to return a dummy/empty orders
func EmptyOrder() (o Order) {
	o.Inf.init()
	o.Inf.SetTags(Empty)
	return
}
