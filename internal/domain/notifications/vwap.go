package notifications

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

// VWAP  is the structure that we have for the VWAP notification
type VWAP struct {
	Asset                string          `json:"asset"`
	OrdersAmount         int             `json:"orders_amount"`
	TimeStart            time.Time       `json:"time_start"`
	TimeEnd              time.Time       `json:"time_end"`
	TimePeriod           time.Duration   `json:"time_period"`
	TimePeriodHuman      string          `json:"time_period_human"`
	CalculationTime      time.Duration   `json:"calculation_time"`
	CalculationTimeHuman string          `json:"calculation_time_human"`
	VWAP                 decimal.Decimal `json:"vwap"`
}

func (v *VWAP) UnJSON(in []byte) (err error) {
	err = json.Unmarshal(in, &v)
	return
}

func (v VWAP) JSON() (out []byte, err error) {
	out, err = json.Marshal(&v)
	return
}
