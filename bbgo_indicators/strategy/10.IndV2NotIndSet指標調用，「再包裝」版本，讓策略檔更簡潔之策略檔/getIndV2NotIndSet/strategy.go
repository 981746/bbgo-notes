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

type MACDSetting struct {
	Interval     types.Interval `json:"interval"`
	ShortWindow  int            `json:"shortWindow"`
	LongWindow   int            `json:"longWindow"`
	SignalWindow int            `json:"signalWindow"`
}

type Strategy struct {
	Symbol   string         `json:"symbol"`
	Interval types.Interval `json:"interval"`
	Window   int            `json:"window"`

	// define indicator and its interval window
	ewma  *EWMAStream
	EMAIW *types.IntervalWindow `json:"emaIW"` // 類型是*types.IntervalWindow，後面就可以直接s.EMAIW.Interval、s.EMAIW.Window

	rsi   *RSIStream
	RSIIW *types.IntervalWindow `json:"rsiIW"`

	// 因為macd指標的參數較多，並不只有interval和window，所以我們另外定義一個struct來處理
	macd        *MACDStream
	MACDSetting *MACDSetting `json:"macdSetting"` // yaml檔裡面的macdSetting欄位接過來，會用新增的MACDSetting struct來處理
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
	if s.RSIIW != nil {
		session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.RSIIW.Interval})
	}
	if s.MACDSetting != nil {
		session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: s.MACDSetting.Interval})
	}
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {

	// 如果yaml檔的emaIW有設定，我們才去調用這個指標
	if s.EMAIW != nil {
		s.initializeEMA(session)
	}
	if s.RSIIW != nil {
		s.initializeRSI(session)
	}
	if s.MACDSetting != nil {
		s.initializeMACD(session)
	}

	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {
		if s.EMAIW != nil {
			fmt.Println("EWMA: ", s.ewma.Last(0))
		}
		if s.RSIIW != nil {
			fmt.Println("RIS: ", s.rsi.Last(0))
		}
		if s.MACDSetting != nil {
			fmt.Println("fast: ", s.macd.FastEWMA.Last(0), "slow:", s.macd.SlowEWMA.Last(0), "hist:", s.macd.Histogram.Last(0))
		}
	}))
	return nil
}

func (s *Strategy) preloadKLines(
	inc *KLineStream, session *bbgo.ExchangeSession, symbol string, interval types.Interval,
) {
	if store, ok := session.MarketDataStore(symbol); ok {
		if kLinesData, ok := store.KLinesOfInterval(interval); ok {
			for _, k := range *kLinesData {
				inc.EmitUpdate(k)
			}
		}
	}
}

func (s *Strategy) initializeEMA(session *bbgo.ExchangeSession) {
	kLines := KLines(session.MarketDataStream, s.Symbol, s.EMAIW.Interval)
	s.ewma = EWMA2(ClosePrices(kLines), s.EMAIW.Window)

	s.preloadKLines(kLines, session, s.Symbol, s.EMAIW.Interval)
}

func (s *Strategy) initializeRSI(session *bbgo.ExchangeSession) {
	kLines := KLines(session.MarketDataStream, s.Symbol, s.RSIIW.Interval)
	s.rsi = RSI2(ClosePrices(kLines), s.RSIIW.Window)

	s.preloadKLines(kLines, session, s.Symbol, s.RSIIW.Interval)
}

func (s *Strategy) initializeMACD(session *bbgo.ExchangeSession) {
	kLines := KLines(session.MarketDataStream, s.Symbol, s.MACDSetting.Interval)
	s.macd = MACD2(ClosePrices(kLines), s.MACDSetting.ShortWindow, s.MACDSetting.LongWindow, s.MACDSetting.SignalWindow)

	s.preloadKLines(kLines, session, s.Symbol, s.MACDSetting.Interval)
}
