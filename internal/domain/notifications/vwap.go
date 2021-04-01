package notifications

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

// VWAP  is the structure that we have for the VWAP notification
type VWAP struct {
	Asset                string          `json:"asset"`
	OrdersAmount         int             `json:"orders_amount"`
	TimePeriodHuman      string          `json:"time_period"`
	CalculationTimeHuman string          `json:"calculation_time"`
	VWAP                 decimal.Decimal `json:"vwap"`
}

<<<<<<< HEAD
//UnJSON
=======
// UnJSON
>>>>>>> 27d5f2ff5f8f7d768344c848b2ce50316e28c857
func (v *VWAP) UnJSON(in []byte) (err error) {
	err = json.Unmarshal(in, &v)
	return
}

<<<<<<< HEAD
//JSON
=======
// JSON
>>>>>>> 27d5f2ff5f8f7d768344c848b2ce50316e28c857
func (v VWAP) JSON() (out []byte, err error) {
	out, err = json.Marshal(&v)
	return
}
