package router

import (
	"fmt"
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/insan1k/proto-order-data/internal/exchanges/coinbase_pro"
	"github.com/insan1k/proto-order-data/internal/log"
	"github.com/insan1k/proto-order-data/internal/vwap_service"
)

type vwapRoute struct {
	asset  string
	route  func(o order.Order)
	notify chan []byte
	quit   func()
}

func newVWAPRoute(pair string) (r vwapRoute, err error) {
	r.asset = pair
	r.notify = make(chan []byte)
	r.route, r.quit, err = vwap_service.Load(pair, 200, r.notify)
	return
}

type VWAPRouter struct {
	exchange        coinbase_pro.CoinbasePro
	pairs           []string
	upstreamChan    chan order.Order
	upstreamQuit    func()
	upstreamService coinbase_pro.WebsocketSubscription
	routes          map[string]vwapRoute
	entry           *apexlog.Entry
	localQuit       chan struct{}
	quit            func()
}

func (v *VWAPRouter) Start() {
	v.exchange = coinbase_pro.CoinbasePro{}
	v.exchange.Defaults()
	v.upstreamService = coinbase_pro.WebsocketSubscription{}
	v.pairs = []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	v.upstreamChan = make(chan order.Order)
	v.entry = log.LoadServiceLog("router")
	v.localQuit = make(chan struct{})
	var err error
	v.routes = make(map[string]vwapRoute)
	for _, pair := range v.pairs {
		route, err := newVWAPRoute(pair)
		if err != nil {
			v.entry.Errorf("creating route for %v failed %v", pair, err)
			return
		}
		v.routes[pair] = route
	}
	go v.notificationHandler()
	v.upstreamQuit, err = v.upstreamService.SubscribeMatches(&v.exchange, v.pairs, v.upstreamChan)
	if err != nil {
		v.entry.Errorf("subscribe error %v", err)
		return
	}
	v.quit = func() {
		close(v.localQuit)
	}
	v.route()
}

func (v *VWAPRouter) Stop() {
	v.quit()
}

func (v *VWAPRouter) route() {
	for {
		select {
		case o := <-v.upstreamChan:
			if r, ok := v.routes[o.Asset]; ok {
				r.route(o)
			}
		case <-v.localQuit:
			v.entry.Info("quitting router")
			v.upstreamQuit()
			for _, route := range v.routes {
				route.quit()
			}
			return
		}
	}
}

func (v *VWAPRouter) notificationHandler() {
	for {
		select {
		case notification := <-v.routes["BTC-USD"].notify:
			fmt.Printf("got vwap %s\n", notification)
		case notification := <-v.routes["ETH-USD"].notify:
			fmt.Printf("got vwap %s\n", notification)
		case notification := <-v.routes["ETH-BTC"].notify:
			fmt.Printf("got vwap %s\n", notification)
		case <-v.localQuit:
			v.entry.Info("quitting notification handler")
			return
		}
	}
}
