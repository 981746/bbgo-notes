package indicatorv2

import (
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

const MaxNumOfPriceAboveBelowMA = 5_000

type AboveBelowType float64

const (
	priceAboveMa AboveBelowType = 1.0
	priceBelowMa AboveBelowType = -1.0
	priceEuqalMa AboveBelowType = 0.0
)

type PriceAboveBelowMAStream struct {
	*types.Float64Series
	price, ma floats.Slice
}

func PriceAboveBelowMA2(price, ma types.Float64Source) *PriceAboveBelowMAStream {
	// 本來想要使用stream對stream去做運算，但無法，目前不知道怎麼用
	// maEWMA := EWMA2(source, window)
	// closePrices := ClosePrices(source)
	// 這兩個source 類型不一樣
	// EWMA2的source是types.Float64Source
	// ClosePrices是source KLineSubscription

	// 所以最後直接從外部拿ema stream跟close price stream進來
	// 使用OnUpdate()拿到這兩個stream的最新值，再push進slice裡
	// 接著在calculate()裡面分別取出最新的ema跟close price，做判斷

	s := &PriceAboveBelowMAStream{
		Float64Series: types.NewFloat64Series(),
	}
	price.OnUpdate(func(v float64) {
		s.price.Push(v)
		s.calculate()
	})
	ma.OnUpdate(func(v float64) {
		s.ma.Push(v)
		s.calculate()
	})
	// s.Bind(source, s)
	// 這邊就不bind了，因為直接從外部拿ema stream跟close price stream進來
	return s
}

func (s *PriceAboveBelowMAStream) calculate() { // Calculate在cross function為小寫，是因為沒有回傳值還是不想給外部用?
	if s.price.Length() != s.ma.Length() {
		return
	}

	current := s.price.Last(0) - s.ma.Last(0)
	if current == 0.0 {
		s.PushAndEmit(float64(priceEuqalMa))
	} else if current > 0 {
		s.PushAndEmit(float64(priceAboveMa))
	} else {
		s.PushAndEmit(float64(priceBelowMa))
	}

}

func (s *PriceAboveBelowMAStream) Truncate() {
	s.Slice = s.Slice.Truncate(MaxNumOfPriceAboveBelowMA)
}
