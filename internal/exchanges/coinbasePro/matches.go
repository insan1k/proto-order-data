package coinbasePro

import (
	"encoding/json"
	apexlog "github.com/apex/log"
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"github.com/insan1k/proto-order-data/internal/domain/order"
	"io"
)

//todo: separate websocket logic from exchange specific logic

// Message
type Message struct {
	Type       string    `json:"type"`
	ProductIDs []string  `json:"product_ids,omitempty"`
	Channels   []Channel `json:"channels,omitempty"`
}

// Channel
type Channel struct {
	Name       string   `json:"name"`
	ProductIDs []string `json:"product_ids,omitempty"`
}

// WebsocketSubscription
type WebsocketSubscription struct {
	exchange           *CoinbasePro
	conn               *websocket.Conn
	message            chan order.Order
	localMessage       chan []byte
	localQuit          chan struct{}
	entry              *apexlog.Entry
	subscribeMessage   Message
	unsubscribeMessage Message
}

func (w *WebsocketSubscription) dial() (err error) {
	w.conn, _, err = websocket.DefaultDialer.Dial(w.exchange.WSSAddress, nil)
	if err != nil {
		w.entry.Errorf("error %v dialing to %v", err, w.exchange.WSSAddress)
		return
	}
	return
}

func (w *WebsocketSubscription) send(message Message) (err error) {
	var msg []byte
	msg, err = json.Marshal(message)
	if err != nil {
		defer w.conn.Close()
		w.entry.Errorf("error %v parsing message json", err)
		return
	}
	err = w.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		defer w.conn.Close()
		w.entry.Errorf("error %v writing message to connection", err)
		return
	}
	return
}

func (w *WebsocketSubscription) receive() (received []byte, err error) {
	var reader io.Reader
	var msgType int
	msgType, reader, err = w.conn.NextReader()
	if err != nil {
		w.entry.Errorf("error %v reading message from connection", err)
	}
	received, err = io.ReadAll(reader)
	if err != nil {
		w.entry.Errorf("error %v reading message from reader", err)
	} else {
		w.entry.Debugf("read %v bytes from reader %v message type", len(received), msgType)
	}
	return
}

// SubscribeMatches subscribes to the matches and sends the results as orders
func (w *WebsocketSubscription) SubscribeMatches(exchange *CoinbasePro, products []string, message chan order.Order) (quit func(), err error) {
	w.message = message
	w.exchange = exchange
	w.localMessage = make(chan []byte)
	w.localQuit = make(chan struct{})
	w.entry = exchange.entry.WithFields(apexlog.Fields{
		"exchange-service": "notify-matches",
		"connection":       "websocket",
	})

	w.subscribeMessage = Message{
		Type: "subscribe",
		Channels: []Channel{{
			Name:       "matches",
			ProductIDs: products,
		}},
	}

	w.unsubscribeMessage = Message{
		Type: "unsubscribe",
		Channels: []Channel{{
			Name:       "matches",
			ProductIDs: products,
		}},
	}

	err = w.dial()
	if err != nil {
		return
	}

	err = w.send(w.subscribeMessage)
	if err != nil {
		return
	}

	quit = func() {
		close(w.localQuit)
	}

	go w.listener()
	go w.treater()

	return
}

func (w *WebsocketSubscription) treater() {
	for {
		select {
		case message := <-w.localMessage:
			w.treatMessage(message)
		case <-w.localQuit:
			w.entry.Infof("exiting websocket treater")
			return
		}
	}
}

func (w *WebsocketSubscription) listener() {
	defer w.conn.Close()
	for {
		select {
		default:
			received, err := w.receive()
			if err != nil {
				w.entry.Errorf("error receiving message %v", err)
			} else {
				w.localMessage <- received
			}
		case <-w.localQuit:
			{
				err := w.send(w.unsubscribeMessage)
				if err != nil {
					w.entry.Errorf("error %v sending message", err)
				}
				return
			}
		}
	}
}

func (w *WebsocketSubscription) treatMessage(message []byte) {
	messageType, err := jsonparser.GetString(message, "type")
	if err != nil {
		w.entry.Errorf("could not parse message type %v", err)
	}
	subscribe := func(message []byte) {
		value, _, _, err := jsonparser.Get(message, "product_ids")
		if err != nil {
			w.entry.Errorf("could not parse key %v", err)
			return
		}
		var assets []string
		_, err = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			assets = append(assets, string(value))
		})
		if err != nil {
			w.entry.Errorf("could not parse array %v", err)
			return
		}
		w.entry.Infof("confirmed subscription for %v", assets)
		return
	}
	match := func(message []byte) {
		quantity, err := jsonparser.GetString(message, "size")
		if err != nil {
			w.entry.Errorf("could not parse key %v", err)
			return
		}
		price, err := jsonparser.GetString(message, "price")
		if err != nil {
			w.entry.Errorf("could not parse key %v", err)
			return
		}
		asset, err := jsonparser.GetString(message, "product_id")
		if err != nil {
			w.entry.Errorf("could not parse key %v", err)
			return
		}
		side, err := jsonparser.GetString(message, "side")
		if err != nil {
			w.entry.Errorf("could not parse key %v", err)
			return
		}
		o := order.NewOrderString(price, quantity)
		o.Asset = asset
		if side == "sell" {
			o.Inf.SetTags(order.Sell, order.Matched)
		}
		if side == "buy" {
			o.Inf.SetTags(order.Buy, order.Matched)
		}
		//todo: if debug is enabled save the original message
		//o.Inf.SetMeta(message)
		w.message <- o
		return
	}
	handlers := map[string]func(message []byte){
		"match":         match,
		"last_match":    match,
		"subscriptions": subscribe,
	}
	if f, ok := handlers[messageType]; ok {
		go f(message)
	} else {
		w.entry.Errorf("unknown message type %v", messageType)
	}
	return
}
