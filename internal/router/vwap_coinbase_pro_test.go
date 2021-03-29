package router

import (
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/log"
	"github.com/insan1k/proto-order-data/internal/util"
	"testing"
	"time"
)

func TestCoinbaseProVWAPRouter(t *testing.T) {
	log.Load(apexlog.ErrorLevel)
	l := log.Get()
	l.Infof("starting vwap feed")
	v := VWAPRouter{}
	go util.Signal(v.Stop)
	go func() {
		tkr := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-tkr.C:
				v.quit()
			}
		}
	}()
	v.Start()
}
