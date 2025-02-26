package getIndStandard

import (
	"context"
	"fmt"

	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/indicator"
	"github.com/c9s/bbgo/pkg/types"
)

const ID = "getIndStandard"

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

	// 調用在StandardIndicatorSet裡的指標，還是要用指標本身的struct來接
	SMA *indicator.SMA

	// use this instance that would call indicator
	// by using the symbol defined in the strategy block of yaml file
	StandardIndicatorSet *bbgo.StandardIndicatorSet
}

func (s *Strategy) ID() string {
	return ID
}

// func (s *Strategy) InstanceID() string {
// 	return fmt.Sprintf("%s:%s:%s", ID, s.Symbol, s.Interval)
// }

func (s *Strategy) Subscribe(session *bbgo.ExchangeSession) {
	// session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.Interval})

	// 訂閱yaml檔策略區塊裡的symbol K線資料
	session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: "5m"})

}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	// now we call the indicator without symbol
	// but acutally it is using the symbol defined in the strategy block of yaml file
	s.SMA = s.StandardIndicatorSet.SMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})

	// after initialized, we can use the OnUpdate() method or the Last(i) method of the indicator
	// we don't have to bind of ourselves
	// StandardIndicatorSet have already done it for us
	s.SMA.OnUpdate(func(v float64) {

		fmt.Println("sma value from OnUpdate: ", v, " symbol:", s.StandardIndicatorSet.Symbol)
	})
	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {

		fmt.Println("sma value with Last(0): ", s.SMA.Last(0))

	}))
	return nil
}
