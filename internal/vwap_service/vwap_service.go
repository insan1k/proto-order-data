package vwap_service

import (
	"fmt"
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/domain/notifications"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/insan1k/proto-order-data/internal/domain/orders"
	"github.com/insan1k/proto-order-data/internal/log"
	"time"
)

type Service struct {
	orders      orders.Orders
	localOrders chan order.Order
	localNotify chan notifications.VWAP
	localStop   chan struct{}
	entry       *apexlog.Entry
}

//Load starts the vwap goroutine
func Load(amountOfOrders int,notifyChan chan []byte) (feed func(o order.Order), stop func(), err error) {
	var service Service
	service.orders, err = orders.NewOrders(amountOfOrders)
	if err != nil {
		return
	}
	//lexical isolation of channels this is useful so that we don't run risk of someone else inadvertently using the
	//channels for this service
	service.localOrders = make(chan order.Order)
	service.localNotify = make(chan notifications.VWAP)
	service.localStop = make(chan struct{})
	service.entry = log.LoadServiceLog("vwap")
	service.entry = service.entry.WithFields(apexlog.Fields{
		"orders-amount":fmt.Sprintf("%v",amountOfOrders),
	})
	stop = func() {
		close(service.localStop)
	}
	feed = func(o order.Order) {
		service.localOrders <- o
	}
	go service.do()
	go service.notify(notifyChan)
	return
}

func (v *Service) do() {
	v.entry.Info("started calculation routine")
	for {
		select {
		case o := <-v.localOrders:
			starTime := time.Now()
			v.orders.Insert(&o)
			vwap := v.orders.VolumeWeightedAveragePrice()
			v.entry.Debugf("calculated vwap: %v",vwap.String())
			v.localNotify <- notifications.VWAP{
				OrdersAmount:    v.orders.Count(),
				TimeStart:       v.orders.TimeStart(),
				TimeEnd:         v.orders.TimeEnd(),
				TimePeriod:      v.orders.TimePeriod(),
				VWAP:            vwap,
				CalculationTime: time.Now().Sub(starTime),
			}
		case <-v.localStop:
			v.entry.Info("quitting calculation routine")
			return
		}
	}
}

func (v *Service) notify(externalNotify chan []byte) {
	v.entry.Info("started notification routine")
	for {
		select {
		case n := <-v.localNotify:
			packed, err := n.JSON()
			if err != nil {
				v.entry.Errorf("failed to send notification: %v", err)
			} else {
				externalNotify <- packed
			}
		case <-v.localStop:
			v.entry.Info("quitting notification routine")
			return
		}
	}
}
