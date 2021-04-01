package router

import (
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/insan1k/proto-order-data/internal/exchanges/coinbasePro"
	"github.com/insan1k/proto-order-data/internal/log"
	"github.com/insan1k/proto-order-data/internal/vwapService"
)

type vwapDownstreamRoute struct {
	asset      string
	downstream func(o order.Order)
	notify     chan []byte
	quit       func()
}

func newVWAPRoute(pair string) (r vwapDownstreamRoute, err error) {
	r.asset = pair
	r.notify = make(chan []byte)
	r.downstream, r.quit, err = vwapService.Load(pair, 200, r.notify)
	return
}

type VWAPRouter struct {
	exchange         coinbasePro.CoinbasePro
	pairs            []string
	upstreamChan     chan order.Order
	upstreamQuit     func()
	upstreamService  coinbasePro.WebsocketSubscription
	downstreamRoutes map[string]vwapDownstreamRoute
	notificationChan chan []byte
	entry            *apexlog.Entry
	localQuit        chan struct{}
	quit             func()
}

func (v *VWAPRouter) Start(notificationChan chan []byte) {
	v.notificationChan = notificationChan
	v.exchange = coinbasePro.CoinbasePro{}
	v.exchange.Defaults()
	v.upstreamService = coinbasePro.WebsocketSubscription{}
	v.pairs = []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	v.upstreamChan = make(chan order.Order)
	v.entry = log.LoadServiceLog("router")
	v.localQuit = make(chan struct{})
	var err error
	v.downstreamRoutes = make(map[string]vwapDownstreamRoute)
	for _, pair := range v.pairs {
		route, err := newVWAPRoute(pair)
		if err != nil {
			v.entry.Errorf("creating downstream for %v failed %v", pair, err)
			return
		}
		v.downstreamRoutes[pair] = route
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
			if r, ok := v.downstreamRoutes[o.Asset]; ok {
				r.downstream(o)
			}
		case <-v.localQuit:
			v.entry.Info("quitting router")
			v.upstreamQuit()
			for _, route := range v.downstreamRoutes {
				route.quit()
			}
			return
		}
	}
}

func (v *VWAPRouter) notificationHandler() {
	for {
		select {
		case notification := <-v.downstreamRoutes["BTC-USD"].notify:
			v.notificationChan <- notification
		case notification := <-v.downstreamRoutes["ETH-USD"].notify:
			v.notificationChan <- notification
		case notification := <-v.downstreamRoutes["ETH-BTC"].notify:
			v.notificationChan <- notification
		case <-v.localQuit:
			v.entry.Info("quitting notification handler")
			return
		}
	}
}
