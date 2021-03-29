package orders

import (
	"github.com/shopspring/decimal"
)

//VolumeWeightedAveragePrice calculates the VWAP for a group of orders
func (o Orders) VolumeWeightedAveragePrice() (vwap decimal.Decimal) {
	sumQuantityMulPrice := decimal.New(0, 20)
	sumQuantity := decimal.New(0, 20)
	calculateVWAP := func(i interface{}) {
		// I don't care much about the unsafe casting here because this is the Orders structure
		// we already initialize each element with it's index, and if the order is nil we don't
		// need to do the calculation for it
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
