package orders

import (
	"github.com/shopspring/decimal"
)

//VolumeWeightedAveragePrice calculates the VWAP for a group of orders
func (o Orders) VolumeWeightedAveragePrice() (vwap decimal.Decimal) {
	//preallocate 32 precision decimal, avoid allocation during mul
	sumQuantityMulPrice := decimal.New(0, 32)
	sumQuantity := decimal.New(0, 32)
	calculateVWAP := func(i interface{}) {
		order := (i).(*Element).Order
		if order != nil {
			sumQuantityMulPrice = sumQuantityMulPrice.Add(order.Price.Mul(order.Quantity))
			sumQuantity = sumQuantity.Add(order.Quantity)
		}
	}
	o.ring.Do(calculateVWAP)
	vwap = sumQuantityMulPrice.Div(sumQuantity)
	return
}
