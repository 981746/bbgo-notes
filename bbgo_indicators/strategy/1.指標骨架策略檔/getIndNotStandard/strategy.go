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
	Symbol       string                  `json:"symbol"`   //交易對名稱會從yaml設定檔傳過來
	Interval     types.Interval          `json:"interval"` //這個接收過來的參數暫時沒用到，我後面直接寫死
	priceCrossMa *indicator.PriceCrossMa // 因為PriceCrossMa struct在indicator package裡，所以要這樣取

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

	priceCrossMaiw := types.IntervalWindow{Window: 20, Interval: s.Interval}
	s.priceCrossMa = &indicator.PriceCrossMa{IntervalWindow: priceCrossMaiw}

	symbol := s.Symbol
	store, _ := session.MarketDataStore(symbol)
	s.priceCrossMa.Bind(store)

	s.priceCrossMa.OnUpdate(func(v float64) {
		fmt.Println("priceCrossMa: ", v)
	})

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {
		fmt.Println(k, "End of m5 data")
	}))
	return nil
}
