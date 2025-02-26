package getIndNotStandard

import (
	"context"
	"fmt"

	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/indicator"
	"github.com/c9s/bbgo/pkg/types"
)

const ID = "getIndNotStandard"

func init() {
	// Register our struct type to BBGO
	// Note that you don't need to field the fields.
	// BBGO uses reflect to parse your type information.
	bbgo.RegisterStrategy(ID, &Strategy{})
}

type Strategy struct {
	Symbol       string         `json:"symbol"`
	Interval     types.Interval `json:"interval"`
	Window       int            `json:"window"`
	priceCrossMa *indicator.PriceCrossMa

	// 加一個官方的sma來比對我目前在priceCrossMa寫的sma邏輯是否正確
	SMA *indicator.SMA
}

func (s *Strategy) ID() string {
	return ID
}

// func (s *Strategy) InstanceID() string {
// 	return fmt.Sprintf("%s:%s:%s", ID, s.Symbol, s.Interval)
// }

func (s *Strategy) Subscribe(session *bbgo.ExchangeSession) {
	// session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.Interval})
	session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: "5m"})
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	// 參考supertrend策略的做法，在外部先bindK(bind的是指定時框之收盤價)，再loadK
	// priceCrossMa指標內部要新增BindK()與LoadK()方法，這樣外部才能調用
	// 這樣居然算出來的sma就跟官方的sma一樣了
	priceCrossMaiw := types.IntervalWindow{Window: 20, Interval: s.Interval}
	s.priceCrossMa = &indicator.PriceCrossMa{IntervalWindow: priceCrossMaiw}

	s.priceCrossMa.BindK(session.MarketDataStream, s.Symbol, s.Interval)
	kLineStore, _ := session.MarketDataStore(s.Symbol)
	if klines, ok := kLineStore.KLinesOfInterval(s.Interval); ok {
		s.priceCrossMa.LoadK((*klines)[0:])
	}
	// 聽取每m5的sma值
	// 這個OnUpdate()不能移到OnKLineClosed()，不然會一直重複印出來
	s.priceCrossMa.OnUpdate(func(v float64) {
		fmt.Println("priceCrossMa: ", v)
	})

	// 官方sma調用方式
	s.SMA = session.StandardIndicatorSet(s.Symbol).SMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {

		indValue := s.priceCrossMa.Last(0)
		fmt.Println(s.priceCrossMa.Values.Length(), "PriceCrossMa sma len", "| value:", indValue)
		fmt.Println(s.SMA.Values.Length(), "官方sma len", "         | value:", s.SMA.Last(0))
		fmt.Println()

		// fmt.Println(k, "End of m5 data")
	}))
	return nil
}
