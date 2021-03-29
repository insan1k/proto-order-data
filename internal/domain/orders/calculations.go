package orders

import (
	"github.com/shopspring/decimal"
)

//VolumeWeightedAveragePrice calculates the VWAP for a group of orders
func (o Orders) VolumeWeightedAveragePrice() (vwap decimal.Decimal) {
	sumQuantityMulPrice := decimal.NewFromFloat(0)
	sumQuantity := decimal.NewFromFloat(0)
	calculateVWAP := func(i interface{}) {
		// I don't care much about the unsafe casting here because this is the Orders structure
		// we already initialize each element with it's index, and if the order is nil we don't
		// need to do the calculation for it
		if (i).(*Element).Order != nil {
			order := (i).(*Element).Order
			quantityMulPrice := order.Price.Mul(order.Quantity)
			sumQuantityMulPrice = sumQuantityMulPrice.Add(quantityMulPrice)
			sumQuantity = sumQuantity.Add(order.Quantity)
		}
	}
	o.ring.Do(calculateVWAP)
	vwap = sumQuantityMulPrice.Div(sumQuantity)
	return
}
