package order

import (
	"github.com/shopspring/decimal"
)

// Order
type Order struct {
	Asset    string
	Price    decimal.Decimal
	Quantity decimal.Decimal
	Inf      Info
}

func (o *Order) preprocess() {
	// initialize stuff
	o.Inf.init()
	// make sure these prices and quantities are initialized with 32 precision decimal
	dec := decimal.New(0, 32)
	o.Price, _ = decimal.RescalePair(o.Price, dec)
	o.Quantity, _ = decimal.RescalePair(o.Quantity, dec)

}

//NewOrderFloat64 based on float64
func NewOrderFloat64(price, quantity float64) (o Order) {
	o.Price = decimal.NewFromFloat(price)
	o.Quantity = decimal.NewFromFloat(quantity)
	o.preprocess()
	return
}

//NewOrderFloat32 based on float32
func NewOrderFloat32(price, quantity float32) (o Order) {
	o.Price = decimal.NewFromFloat32(price)
	o.Quantity = decimal.NewFromFloat32(quantity)
	o.preprocess()
	return
}

//NewOrderString creates from string
func NewOrderString(price, quantity string) (o Order) {
	o.Price = decimal.RequireFromString(price)
	o.Quantity = decimal.RequireFromString(quantity)
	o.preprocess()
	return
}

//EmptyOrder convenience function to return a dummy/empty orders
func EmptyOrder() (o Order) {
	o.Inf.init()
	o.Inf.SetTags(Empty)
	return
}
