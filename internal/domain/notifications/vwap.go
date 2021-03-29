package notifications

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

// VWAP  is the structure that we have for the VWAP notification
type VWAP struct {
	OrdersAmount    int             `json:"orders_amount" msgpack:"a"`
	TimeStart       time.Time       `json:"time_start" msgpack:"s"`
	TimeEnd         time.Time       `json:"time_end" msgpack:"e"`
	TimePeriod      time.Duration   `json:"time_period" msgpack:"p"`
	CalculationTime time.Duration   `json:"calculation_time" msgpack:"t"`
	VWAP            decimal.Decimal `json:"vwap" msgpack:"v"`
}

func (v *VWAP) UnJSON(in []byte) (err error) {
	err = json.Unmarshal(in, &v)
	return
}

func (v VWAP) JSON() (out []byte, err error) {
	out, err = json.Marshal(&v)
	return
}
