package zerodha

import (
	"fmt"
	"log"
	"sync"

	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

type AlertChecker interface {
	CheckAlerts(tick kitemodels.Tick)
}

type ZerodhaTicker struct {
	ticker           *kiteticker.Ticker
	subscribedTokens map[uint32]bool
	mu               sync.Mutex
	checker          AlertChecker
}

func NewZerodhaTicker(apikey, accessToken string, checker AlertChecker) *ZerodhaTicker {
	zTicker := &ZerodhaTicker{
		subscribedTokens: make(map[uint32]bool),
		checker:          checker,
	}

	zTicker.ticker = kiteticker.New(apikey, accessToken)
	zTicker.ticker.OnConnect(zTicker.OnConnect)
	zTicker.ticker.OnTick(zTicker.OnTick)
	zTicker.ticker.OnError(zTicker.OnError)
	zTicker.ticker.OnClose(zTicker.onClose)

	return zTicker
}

func (z *ZerodhaTicker) Start(){
	go z.ticker.Serve()
}

func (z *ZerodhaTicker) Subscribe(token uint32) {
	z.mu.Lock()
	defer z.mu.Unlock()
	if z.subscribedTokens[token] {
		return
	}
	z.subscribedTokens[token] = true
	tokens := []uint32{token}
	if err := z.ticker.Subscribe(tokens); err != nil {
		log.Println(err)
		return
	}
	z.ticker.SetMode(kiteticker.ModeLTP, tokens)
}

func (z *ZerodhaTicker) OnError(err error) {
	// Handle error
	fmt.Println(err)
}

func (z *ZerodhaTicker) OnClose(code int, reason string) {
	// Handle close
}

func (z *ZerodhaTicker) OnConnect() {
	// Handle connectma
	log.Println("ticker connected")
	z.mu.Lock()
	defer z.mu.Unlock()
	tokens := make([]uint32, 0, len(z.subscribedTokens))
	for t := range z.subscribedTokens {
		tokens = append(tokens, t)
	}
	if len(tokens) > 0 {
		z.ticker.Subscribe(tokens)
		z.ticker.SetMode(kiteticker.ModeLTP, tokens)
	}
	// err :- z.ticker.SetMode(kiteticker.ModeLTP, )
}

func (z *ZerodhaTicker) OnTick(tick kitemodels.Tick) {
	z.checker.CheckAlerts(tick)
}

func (z *ZerodhaTicker)onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}