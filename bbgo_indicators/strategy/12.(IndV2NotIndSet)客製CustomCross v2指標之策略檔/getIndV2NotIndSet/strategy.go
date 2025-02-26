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

	isPriceCrossMa int
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

	// 客製CustomCross v2指標
	price := session.Indicators(s.Symbol).CLOSE(types.Interval(s.Interval))
	ma := session.Indicators(s.Symbol).EWMA(types.IntervalWindow{Interval: s.Interval, Window: s.Window})
	CrossInd := CustomCross(price, ma)
	CrossInd.OnUpdate(func(v float64) {
		// 這邊是只有發生上穿或下穿才會推播
		// v can be 1.0, -1.0
		// fmt.Println("發生穿越：", CrossInd.IsCross.Last()) // this is always true in OnUpdate()
		// fmt.Println("CustomCross: ", v)

		if v == 1.0 {
			s.isPriceCrossMa = 1
		} else if v == -1.0 {
			s.isPriceCrossMa = -1
		}

	})
	count := 0
	session.MarketDataStream.OnKLineClosed(types.KLineWith(s.Symbol, s.Interval, func(k types.KLine) {

		count++

		// 1.這個是bbgo官方使用Cross的方式，s.isPriceCrossMa值有可能會是0喔，因為前面幾筆資料尚未推播指標值
		// 他這個方式是上穿了，如果還沒下穿，s.isPriceCrossMa值就一直是1，就會一直買
		// 這樣的方式在CrossInd.OnUpdate()裡面要先把1 或 -1的值存起來，才能在這邊判斷
		// fmt.Println(count, ". 是否發生穿越：", s.isPriceCrossMa)

		// 2.其實如果用原版的Cross()，可以這樣寫
		// 這樣也能判斷此根收盤K棒是否發生穿越
		// 這樣的方式在CrossInd.OnUpdate()裡面要先把1 或 -1的值存起來，才能在這邊判斷
		// if s.isPriceCrossMa == 1 || s.isPriceCrossMa == -1 {
		// 	if s.isPriceCrossMa == 1 {
		// 		fmt.Println(count, ". 發生穿越-上穿")

		// 	} else if s.isPriceCrossMa == -1 {
		// 		fmt.Println(count, ". 發生穿越-下穿")

		// 	}
		// 	s.isPriceCrossMa = 0
		// } else {
		// 	fmt.Println(count, ". 沒有發生穿越")
		// }

		// 3.印出相關資訊，Debug用
		// 如果CrossInd.Last(0)的值為0，代表指標值還沒有計算出來，無法判斷價格與均線狀況
		// fmt.Println(count, ". 是否發生穿越：", CrossInd.IsCross.Last())
		// fmt.Println("CrossInd.Last(0): ", CrossInd.Last(0))
		// fmt.Println()

		// 4.要判斷此根收盤K棒是否發生穿越，發生上穿還是下穿請用以下方式(使用指標內的IsCross boolean slice)
		fmt.Println(count, ".")
		if CrossInd.IsCross.Last() { // 此根收盤K棒是否發生穿越
			if CrossInd.Last(0) == 1.0 {
				fmt.Println("發生CrossOver")
			} else if CrossInd.Last(0) == -1.0 {
				fmt.Println("發生CrossUnder")
			}
		}

	}))
	return nil
}
