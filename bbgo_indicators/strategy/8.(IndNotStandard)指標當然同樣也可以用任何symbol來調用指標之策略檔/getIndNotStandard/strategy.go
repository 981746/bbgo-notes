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
	Symbol   string         `json:"symbol"`
	Interval types.Interval `json:"interval"`
	Window   int            `json:"window"`

	PriceAboveBelowMA *indicator.PriceAboveBelowMA
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

	// 特別訂閱DOGEUSDT的kline
	session.Subscribe(types.KLineChannel, "DOGEUSDT", types.SubscribeOptions{Interval: "5m"})
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	priceAboveBelowMAiw := types.IntervalWindow{Window: s.Window, Interval: s.Interval}
	s.PriceAboveBelowMA = &indicator.PriceAboveBelowMA{IntervalWindow: priceAboveBelowMAiw}
	s.PriceAboveBelowMA.BindK(session.MarketDataStream, "DOGEUSDT", s.Interval)
	kLineStore, _ := session.MarketDataStore("DOGEUSDT")
	if klines, ok := kLineStore.KLinesOfInterval(s.Interval); ok {
		s.PriceAboveBelowMA.LoadK((*klines)[0:])
	}
	s.PriceAboveBelowMA.OnUpdate(func(v float64) {
		// There are three possible outcomes of this indicator: 0, 1, -1
		// 0 stands price equals ma; 1 stands price above ma; -1 stands price below ma
		// fmt.Println("PriceAboveBelowMA: ", v)
	})

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {
		aboveBelow := s.PriceAboveBelowMA.Last(0)
		maValue := s.PriceAboveBelowMA.MAValues.Last(0)
		closeprice := s.PriceAboveBelowMA.PricesAboveBelowMA.Last(0)

		fmt.Println("收盤價: ", closeprice, "MA值: ", maValue, "aboveBelow: ", aboveBelow)

	}))
	return nil
}
