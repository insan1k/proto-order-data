<<<<<<< HEAD:internal/exchanges/coinbasepro/matches_test.go
package coinbasepro
=======
package coinbasePro
>>>>>>> 27d5f2ff5f8f7d768344c848b2ce50316e28c857:internal/exchanges/coinbasePro/matches_test.go

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
