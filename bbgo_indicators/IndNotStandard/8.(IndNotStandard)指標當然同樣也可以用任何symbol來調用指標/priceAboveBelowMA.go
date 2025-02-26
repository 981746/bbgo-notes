package indicator

import (
	"time"

	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

const MaxNumOfPriceAboveBelowMA = 5_000
const MaxNumOfPriceAboveBelowMATruncateSize = 100

//go:generate callbackgen -type PriceAboveBelowMA
type PriceAboveBelowMA struct {
	Values             floats.Slice
	MAValues           floats.Slice
	PricesAboveBelowMA floats.Slice
	updateCallbacks    []func(value float64)
	types.IntervalWindow
	EndTime time.Time
	types.SeriesBase

	// Ma type for using SMA or EWMA
	MAType string
	maEWMA *EWMA
	maSMA  *SMA
}

func (inc *PriceAboveBelowMA) Last(i int) float64 {
	return inc.Values.Last(i)
}

func (inc *PriceAboveBelowMA) Index(i int) float64 {
	return inc.Last(i)
}

func (inc *PriceAboveBelowMA) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &PriceAboveBelowMA{}

func (inc *PriceAboveBelowMA) Update(value float64) {
	// indicator calculation here...
	// push value...

	if len(inc.Values) == 0 {
		switch inc.MAType {
		case "SMA":
			inc.maSMA = &SMA{IntervalWindow: types.IntervalWindow{Window: inc.Window, Interval: inc.Interval}}
		case "EWMA", "EMA":
			inc.maEWMA = &EWMA{IntervalWindow: types.IntervalWindow{Window: inc.Window, Interval: inc.Interval}}
		default:
			// default we use EWMA
			// if outside didn't assign ma type, then use EWMA
			inc.maEWMA = &EWMA{IntervalWindow: types.IntervalWindow{Window: inc.Window, Interval: inc.Interval}}

		}
	}

	var ma float64
	switch inc.MAType {
	case "SMA":
		inc.maSMA.Update(value)
		ma = inc.maSMA.Last(0)
	case "EWMA", "EMA":
		inc.maEWMA.Update(value)
		ma = inc.maEWMA.Last(0)
	default:
		inc.maEWMA.Update(value)
		ma = inc.maEWMA.Last(0)
	}

	// now I can get ma value and close price from strategy
	inc.MAValues.Push(ma)
	inc.PricesAboveBelowMA.Push(value)

	if value > ma {
		inc.Values.Push(1)

	} else if value < ma {
		inc.Values.Push(-1)
	} else {
		inc.Values.Push(0)
	}

	if len(inc.Values) > MaxNumOfPriceAboveBelowMA {
		inc.Values = inc.Values[MaxNumOfPriceAboveBelowMATruncateSize-1:]
	}
}

func (inc *PriceAboveBelowMA) PushK(k types.KLine) {

	if inc.EndTime != zeroTime && k.EndTime.Before(inc.EndTime) {
		return
	}

	inc.Update(k.Close.Float64())
	inc.EndTime = k.EndTime.Time()
	inc.EmitUpdate(inc.Values.Last(0))

}

// 這樣外部可以以收盤價的K線來bind
func (inc *PriceAboveBelowMA) BindK(target KLineClosedEmitter, symbol string, interval types.Interval) {
	target.OnKLineClosed(types.KLineWith(symbol, interval, inc.PushK))
}

// 在外部載入K線
func (inc *PriceAboveBelowMA) LoadK(allKLines []types.KLine) {
	for _, k := range allKLines {
		inc.PushK(k)
	}
}
