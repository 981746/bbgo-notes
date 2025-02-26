package getIndV2NotIndSet

import (
	"context"
	"fmt"

	"github.com/c9s/bbgo/pkg/bbgo"
	. "github.com/c9s/bbgo/pkg/indicator/v2"
	"github.com/c9s/bbgo/pkg/types"
)

const ID = "getIndV2NotIndSet"

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
}

func (s *Strategy) ID() string {
	return ID
}

func (s *Strategy) InstanceID() string {
	return fmt.Sprintf("%s:%s:%s", ID, s.Symbol, s.Interval)
}

func (s *Strategy) Subscribe(session *bbgo.ExchangeSession) {
	// session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.Interval})
	session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: "5m"})

}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	// 新版PriceAboveBelowMA v2指標
	price := session.Indicators(s.Symbol).CLOSE(types.Interval(s.Interval))
	ma := session.Indicators(s.Symbol).EWMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})
	isAboveBelow := PriceAboveBelowMA2(price, ma) // 把price stream跟ma stream傳進去
	isAboveBelow.OnUpdate(func(v float64) {
		fmt.Println("PriceAboveBelowMA v2: ", v)
	})

	// 如果是SMA，沒有在indicator set，所以要用這個方式
	// price := session.Indicators(s.Symbol).CLOSE(types.Interval(s.Interval))

	// kLines := KLines(session.MarketDataStream, s.Symbol, s.Interval)
	// ma := SMA(ClosePrices(kLines), s.Window)

	// isAboveBelow := PriceAboveBelowMA2(price, ma) // 把price stream跟ma stream傳進去
	// isAboveBelow.OnUpdate(func(v float64) {
	// 	fmt.Println("PriceAboveBelowMA v2: ", v)
	// })

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {

		// fmt.Println("isAboveBelow.Last(0): ", isAboveBelow.Last(0))

	}))
	return nil
}
