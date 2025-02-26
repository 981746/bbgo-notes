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
	Symbol            string         `json:"symbol"`
	Interval          types.Interval `json:"interval"`
	Window            int            `json:"window"`
	priceCrossMa      *indicator.PriceCrossMa
	PriceAboveBelowMA *indicator.PriceAboveBelowMA

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

	// Bind()+filter interval方法
	// priceCrossMaiw := types.IntervalWindow{Window: 20, Interval: s.Interval}
	// s.priceCrossMa = &indicator.PriceCrossMa{IntervalWindow: priceCrossMaiw}

	// symbol := s.Symbol
	// store, _ := session.MarketDataStore(symbol)
	// s.priceCrossMa.Bind(store)

	// s.priceCrossMa.OnUpdate(func(v float64) {
	// 	fmt.Println("priceCrossMa: ", v)
	// })

	// 在外部BindK與LoadK方法
	// priceCrossMaiw := types.IntervalWindow{Window: s.Window, Interval: s.Interval}
	// s.priceCrossMa = &indicator.PriceCrossMa{IntervalWindow: priceCrossMaiw}

	// s.priceCrossMa.BindK(session.MarketDataStream, s.Symbol, s.Interval)
	// kLineStore, _ := session.MarketDataStore(s.Symbol)
	// if klines, ok := kLineStore.KLinesOfInterval(s.Interval); ok {
	// 	s.priceCrossMa.LoadK((*klines)[0:])
	// }

	// s.priceCrossMa.OnUpdate(func(v float64) {
	// 	// There are three possible outcomes of this indicator: 0, 1, -1
	// 	// 0 stands price equals ma; 1 stands price above ma; -1 stands price below ma
	// 	fmt.Println("priceCrossMa: ", v)
	// })

	// 調用PriceAboveBelowMA指標
	priceAboveBelowMAiw := types.IntervalWindow{Window: s.Window, Interval: s.Interval}
	s.PriceAboveBelowMA = &indicator.PriceAboveBelowMA{IntervalWindow: priceAboveBelowMAiw, MAType: "SMA"}
	s.PriceAboveBelowMA.BindK(session.MarketDataStream, s.Symbol, s.Interval)
	kLineStore, _ := session.MarketDataStore(s.Symbol)
	if klines, ok := kLineStore.KLinesOfInterval(s.Interval); ok {
		s.PriceAboveBelowMA.LoadK((*klines)[0:])
	}
	s.PriceAboveBelowMA.OnUpdate(func(v float64) {
		// There are three possible outcomes of this indicator: 0, 1, -1
		// 0 stands price equals ma; 1 stands price above ma; -1 stands price below ma
		fmt.Println("PriceAboveBelowMA: ", v)
	})

	// 官方sma調用方式
	// s.SMA = session.StandardIndicatorSet(s.Symbol).SMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})
	// s.SMA.OnUpdate(func(v float64) {

	// 	fmt.Println("官方sma: ", v)
	// })

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {

		// indValue := s.priceCrossMa.Last(0)
		// fmt.Println(s.priceCrossMa.Values.Length(), "PriceCrossMa sma len", "| value:", indValue)
		// fmt.Println(s.SMA.Values.Length(), "官方sma len", "  s       | value:", s.SMA.Last(0))
		// fmt.Println()

		// 調用指標判斷目前是價格大於還是小於EMA
		indValue := s.PriceAboveBelowMA.Last(0)
		if indValue == 1 {
			fmt.Printf("價格大於%v %v%v, endTime:%v\n", s.Interval, s.PriceAboveBelowMA.MAType, s.Window, s.PriceAboveBelowMA.EndTime)
		}
		if indValue == -1 {
			fmt.Printf("價格小於%v %v%v, endTime:%v\n", s.Interval, s.PriceAboveBelowMA.MAType, s.Window, s.PriceAboveBelowMA.EndTime)
		}

		// fmt.Println(k, "End of m5 data")
	}))
	return nil
}
