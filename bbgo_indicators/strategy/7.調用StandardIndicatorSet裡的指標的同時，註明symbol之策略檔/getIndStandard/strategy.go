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
	SMA  *indicator.SMA
	SMA2 *indicator.SMA
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

	// 特別訂閱MATICUSDT的kline
	session.Subscribe(types.KLineChannel, "MATICUSDT", types.SubscribeOptions{Interval: "5m"})
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	// 位於StandardIndicatorSet裡的指標在調用時
	// 不用在策略檔的地方使用Bind()或BindK()+LoadK()方法
	// 初始化完，就可以直接使用OnUpdate()方法或使用指標的Last(i)方法
	// 這邊使用session的StandardIndicatorSet()必須填寫symbol，再調用SMA()，
	// 所以指標與策略yaml的symbol不同，是可以的
	// 意即在同一個策略檔中，可以使用不同的symbol的資料的指標
	// 但要做(1)、(2)的設定，才能正常執行
	// (1)在yaml的backtest區塊加上其他的symbol，我多加了MATICUSDT這個symbol
	// symbols:
	// - BTCUSDT
	// - MATICUSDT
	// (2)在Subscribe()訂閱MATICUSDT的kline
	// session.Subscribe(types.KLineChannel, "MATICUSDT", types.SubscribeOptions{Interval: "5m"})
	s.SMA = session.StandardIndicatorSet(s.Symbol).SMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})
	s.SMA2 = session.StandardIndicatorSet("MATICUSDT").SMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})

	s.SMA.OnUpdate(func(v float64) {
		fmt.Println("BTCUSDT官方sma: ", v)
	})
	s.SMA2.OnUpdate(func(v float64) {
		fmt.Println("MATICUSDT官方sma: ", v)
	})

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {

		fmt.Println("BTCUSDT官方sma with Last(0): ", s.SMA.Last(0))
		fmt.Println("MATICUSDT官方sma with Last(0): ", s.SMA2.Last(0))
		fmt.Println()
	}))
	return nil
}
