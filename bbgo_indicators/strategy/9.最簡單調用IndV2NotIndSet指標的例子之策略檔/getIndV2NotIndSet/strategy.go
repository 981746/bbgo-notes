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

	// define indicator and its interval window
	ewma  *EWMAStream
	EMAIW *types.IntervalWindow `json:"emaIW"`
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

	// 如果yaml檔的emaIW有設定，就訂閱emaIW的kline
	if s.EMAIW != nil {
		session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.EMAIW.Interval})
	}
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	// 如果yaml檔的emaIW有設定，我們才去調用這個指標
	// 如果我們想要讓某個指標為optional的話，可以用這樣的方式
	// 如果有在yaml填寫相關設定，就去調用指標，如果沒有填寫，程式還是可以正常運行
	if s.EMAIW != nil {
		kLines := KLines(session.MarketDataStream, s.Symbol, s.EMAIW.Interval)
		s.ewma = EWMA2(ClosePrices(kLines), s.EMAIW.Window)
		if store, ok := session.MarketDataStore(s.Symbol); ok {
			if kLinesData, ok := store.KLinesOfInterval(s.EMAIW.Interval); ok {
				for _, k := range *kLinesData {
					kLines.EmitUpdate(k)

				}
			}
		}
		s.ewma.OnUpdate(func(v float64) {

			fmt.Println("OnUpdate v: ", v)
		})

	}

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {
		// fmt.Println("EWMA: ", s.ewma.Last(0))

	}))
	return nil
}
