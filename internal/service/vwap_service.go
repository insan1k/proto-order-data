package service

import (
	"fmt"
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/domain/notifications"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/insan1k/proto-order-data/internal/domain/orders"
	"github.com/insan1k/proto-order-data/internal/log"
	"time"
)

//Service holds service information
type Service struct {
	orders      orders.Orders
	localOrders chan order.Order
	localNotify chan notifications.VWAP
	localStop   chan struct{}
	entry       *apexlog.Entry
}

//Asset returns the asset this service is responsible for
func (s Service) Asset() string {
	return s.orders.Asset()
}

//Load starts the vwap goroutine
func Load(assetName string, amountOfOrders int, notifyChan chan []byte) (feed func(o order.Order), stop func(), err error) {
	var service Service
	service.orders, err = orders.NewOrders(assetName, amountOfOrders)
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
		"orders-amount": fmt.Sprintf("%v", amountOfOrders),
		"asset-name":    assetName,
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

func (s *Service) do() {
	s.entry.Info("started calculation routine")
	for {
		select {
		case o := <-s.localOrders:
			starTime := time.Now()
			s.orders.Insert(&o)
			vwap := s.orders.VolumeWeightedAveragePrice()
			doneTime := time.Now().Sub(starTime).String()
			s.localNotify <- notifications.VWAP{
				Asset:                s.orders.Asset(),
				OrdersAmount:         s.orders.Len(),
				TimePeriodHuman:      s.orders.TimePeriod().String(),
				VWAP:                 vwap,
				CalculationTimeHuman: doneTime,
			}
		case <-s.localStop:
			s.entry.Info("quitting calculation routine")
			return
		}
	}
}

func (s *Service) notify(externalNotify chan []byte) {
	s.entry.Info("started notification routine")
	for {
		select {
		case n := <-s.localNotify:
			packed, err := n.JSON()
			if err != nil {
				s.entry.Errorf("failed to send notification: %s", err)
			} else {
				externalNotify <- packed
			}
		case <-s.localStop:
			s.entry.Info("quitting notification routine")
			return
		}
	}
}
