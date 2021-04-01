package router

import (
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/insan1k/proto-order-data/internal/exchanges/coinbasepro"
	"github.com/insan1k/proto-order-data/internal/log"
	"github.com/insan1k/proto-order-data/internal/service"
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
	r.downstream, r.quit, err = service.Load(pair, 200, r.notify)
	return
}

// VWAPRouter holds routes and other relevant router things
type VWAPRouter struct {
	exchange         coinbasepro.CoinbasePro
	pairs            []string
	upstreamChan     chan order.Order
	upstreamQuit     func()
	upstreamService  coinbasepro.WebsocketSubscription
	downstreamRoutes map[string]vwapDownstreamRoute
	notificationChan chan []byte
	entry            *apexlog.Entry
	localQuit        chan struct{}
	quit             func()
}

// Start runs this router
func (v *VWAPRouter) Start(notificationChan chan []byte) {
	v.notificationChan = notificationChan
	v.exchange = coinbasepro.CoinbasePro{}
	v.exchange.Defaults()
	v.upstreamService = coinbasepro.WebsocketSubscription{}
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

// Stop ends/quits this router
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
