<<<<<<< HEAD:internal/exchanges/coinbasepro/exchange.go
package coinbasepro
=======
package coinbasePro
>>>>>>> 27d5f2ff5f8f7d768344c848b2ce50316e28c857:internal/exchanges/coinbasePro/exchange.go

import (
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/log"
)

// CoinbasePro exchange struct
type CoinbasePro struct {
	ExchangeName string
	WSSAddress   string
	entry        *apexlog.Entry
}

// Defaults sets default values for exchange struct
func (c *CoinbasePro) Defaults() {
	c.ExchangeName = "CoinbasePro"
	c.WSSAddress = "wss://ws-feed.pro.coinbase.com"
	c.entry = log.LoadServiceLog(c.ExchangeName)
}

//todo: create base exchange that abstracts most methods
// in my experience since the knowledge domain of an exchange is usually quite similar you can have a model that
// essentially assists in the development of exchanges
