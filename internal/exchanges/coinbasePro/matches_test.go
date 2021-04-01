package coinbasePro

import (
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMatches(t *testing.T) {
	e := CoinbasePro{}
	e.Defaults()
	w := WebsocketSubscription{}
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	message := make(chan order.Order)
	quit, err := w.SubscribeMatches(&e, pairs, message)
	if err != nil {
		t.Fail()
	}
	var oo []order.Order
	tkr := time.NewTicker(10 * time.Second)
	for {
		select {
		case o := <-message:
			assert.NotEmpty(t, o)
			oo = append(oo, o)
		case <-tkr.C:
			assert.NotEmpty(t, oo)
			quit()
			return
		}
	}
}
