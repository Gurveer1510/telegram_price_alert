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
	startOnce        sync.Once
	connected        bool
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

func (z *ZerodhaTicker) Start() {
	z.startOnce.Do(func() {
		go z.ticker.Serve()
	})
}

func (z *ZerodhaTicker) Subscribe(token uint32) {
	z.mu.Lock()
	if z.subscribedTokens[token] {
		z.mu.Unlock()
		return
	}
	z.subscribedTokens[token] = true
	connected := z.connected && z.ticker.Conn != nil
	z.mu.Unlock()

	if !connected {
		log.Printf("ticker not connected yet, queued subscription for token %d", token)
		return
	}

	z.subscribeTokens([]uint32{token})
}

func (z *ZerodhaTicker) OnError(err error) {
	z.mu.Lock()
	z.connected = false
	z.mu.Unlock()
	log.Println(err)
}

func (z *ZerodhaTicker) OnConnect() {
	log.Println("ticker connected")
	z.mu.Lock()
	z.connected = true
	tokens := make([]uint32, 0, len(z.subscribedTokens))
	for t := range z.subscribedTokens {
		tokens = append(tokens, t)
	}
	z.mu.Unlock()

	z.subscribeTokens(tokens)
}

func (z *ZerodhaTicker) OnTick(tick kitemodels.Tick) {
	z.checker.CheckAlerts(tick)
}

func (z *ZerodhaTicker) onClose(code int, reason string) {
	z.mu.Lock()
	z.connected = false
	z.mu.Unlock()
	fmt.Println("Close: ", code, reason)
}

func (z *ZerodhaTicker) subscribeTokens(tokens []uint32) {
	if len(tokens) == 0 {
		return
	}

	if err := z.ticker.Subscribe(tokens); err != nil {
		log.Println(err)
		return
	}

	if err := z.ticker.SetMode(kiteticker.ModeLTP, tokens); err != nil {
		log.Println(err)
	}
}
