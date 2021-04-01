package service

import (
	"github.com/insan1k/proto-order-data/internal/domain/notifications"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/shopspring/decimal"
	"testing"
)

type VWAPCase struct {
	Name   string
	Orders []order.Order
	Result decimal.Decimal
}

var VWAPCases = []VWAPCase{
	{
		Name: "test_easy_64",
		Orders: []order.Order{
			order.NewOrderFloat64(10, 100),
			order.NewOrderFloat64(8, 300),
			order.NewOrderFloat64(11, 200),
		},
		Result: decimal.RequireFromString("9.3333333333333333"),
	},
	{
		Name: "test_easy_32",
		Orders: []order.Order{
			order.NewOrderFloat32(10, 100),
			order.NewOrderFloat32(8, 300),
			order.NewOrderFloat32(11, 200),
		},
		Result: decimal.RequireFromString("9.3333333333333333"),
	},
	{
		Name: "test_easy_string",
		Orders: []order.Order{
			order.NewOrderString("10", "100"),
			order.NewOrderString("8", "300"),
			order.NewOrderString("11", "200"),
		},
		Result: decimal.RequireFromString("9.3333333333333333"),
	},
	{
		Name: "test_medium_string",
		Orders: []order.Order{
			order.NewOrderString("13.4993", "33.39392393"),
			order.NewOrderString("12.4324", "435.4293939"),
			order.NewOrderString("13.65454123", "44.4284892"),
			order.NewOrderString("13.865934", "55.49593953"),
			order.NewOrderString("14.39242", "525.348123"),
			order.NewOrderString("20.492302", "254234.3939"),
			order.NewOrderString("18.39499", "23454.3294"),
			order.NewOrderString("13.499459", "533.393931"),
			order.NewOrderString("12.43949", "3432.2349"),
			order.NewOrderString("13.350986", "443.42933"),
		},
		Result: decimal.RequireFromString("20.1697435628860722"),
	},
	{
		Name: "test_rollover_string",
		Orders: []order.Order{
			order.NewOrderString("10", "100"),
			order.NewOrderString("8", "300"),
			order.NewOrderString("11", "200"),
			order.NewOrderString("10", "100"),
			order.NewOrderString("8", "300"),
			order.NewOrderString("11", "200"),
			order.NewOrderString("10", "100"),
			order.NewOrderString("8", "300"),
			order.NewOrderString("11", "200"),
			order.NewOrderString("10", "100"),
			order.NewOrderString("8", "300"),
			order.NewOrderString("11", "200"),
			order.NewOrderString("10", "100"),
			order.NewOrderString("8", "300"),
			order.NewOrderString("11", "200"),
			order.NewOrderString("13.4993", "33.39392393"),
			order.NewOrderString("12.4324", "435.4293939"),
			order.NewOrderString("13.65454123", "44.4284892"),
			order.NewOrderString("13.865934", "55.49593953"),
			order.NewOrderString("14.39242", "525.348123"),
			order.NewOrderString("20.492302", "254234.3939"),
			order.NewOrderString("18.39499", "23454.3294"),
			order.NewOrderString("13.499459", "533.393931"),
			order.NewOrderString("12.43949", "3432.2349"),
			order.NewOrderString("13.350986", "443.42933"),
		},
		Result: decimal.RequireFromString("20.1697435628860722"),
	},
}

func TestVWAP(t *testing.T) {
	notifyChan := make(chan []byte)
	for _, tt := range VWAPCases {
		t.Run(tt.Name, func(t *testing.T) {
			feed, stop, err := Load("dummy-asset", 10, notifyChan)
			if err != nil {
				t.Fatalf("starting VWAP service failed: %v", err)
			}
			var Notification notifications.VWAP
			for _, o := range tt.Orders {
				feed(o)
				resp := <-notifyChan
				err := Notification.UnJSON(resp)
				if err != nil {
					t.Fatalf("could not read notification from vwap: %v", err)
				}
			}
			stop()
			got := Notification.VWAP.String()
			wanted := tt.Result.String()
			if got != wanted {
				t.Logf("case %v, failed expected %v got %v",
					tt.Name,
					tt.Result.String(),
					Notification.VWAP.String(),
				)
				t.Fail()
			}
		})
	}
}
