package main

import (
	apexlog "github.com/apex/log"
	"github.com/insan1k/proto-order-data/internal/log"
	"github.com/insan1k/proto-order-data/internal/router"
	"github.com/insan1k/proto-order-data/internal/util"
)

func main() {
	log.Load(apexlog.ErrorLevel)
	l := log.Get()
	l.Infof("starting vwap feed")
	v := router.VWAPRouter{}
	notificationChan:=make(chan []byte)
	q:=util.Printer(notificationChan)
	go util.Signal(q)
	go util.Signal(v.Stop)
	v.Start(notificationChan)
}
