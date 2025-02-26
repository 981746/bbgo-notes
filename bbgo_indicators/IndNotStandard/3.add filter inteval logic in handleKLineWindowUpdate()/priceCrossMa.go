package indicator

import (
	"time"

	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

// These two constants are already declared in pkg/indicator/sma.go
// so I don't need to declare again here
// const MaxNumOfSMA = 5_000
// const MaxNumOfSMATruncateSize = 100

//go:generate callbackgen -type PriceCrossMa
type PriceCrossMa struct {
	Values          floats.Slice
	updateCallbacks []func(value float64)
	types.IntervalWindow

	types.SeriesBase
	rawValues *types.Queue
	EndTime   time.Time
}

func (inc *PriceCrossMa) Last(i int) float64 {
	return inc.Values.Last(i)
}

func (inc *PriceCrossMa) Index(i int) float64 {
	return inc.Last(i)
}

func (inc *PriceCrossMa) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &PriceCrossMa{}

func (inc *PriceCrossMa) Update(value float64) {
	// indicator calculation here...
	// push value...

	if inc.rawValues == nil {
		inc.rawValues = types.NewQueue(inc.Window)
		inc.SeriesBase.Series = inc
	}

	inc.rawValues.Update(value)
	if inc.rawValues.Length() < inc.Window {
		return
	}

	inc.Values.Push(types.Mean(inc.rawValues))
	if len(inc.Values) > MaxNumOfSMA {
		inc.Values = inc.Values[MaxNumOfSMATruncateSize-1:]
	}

}

func (inc *PriceCrossMa) PushK(k types.KLine) {

	if inc.EndTime != zeroTime && k.EndTime.Before(inc.EndTime) {
		return
	}

	inc.Update(k.Close.Float64())
	inc.EndTime = k.EndTime.Time()
	inc.EmitUpdate(inc.Values.Last(0))
}

// custom function
func (inc *PriceCrossMa) CalculateAndUpdate(allKLines []types.KLine) {

	if inc.rawValues == nil {
		for _, k := range allKLines {
			inc.PushK(k)
		}

	} else {
		var last = allKLines[len(allKLines)-1]
		inc.PushK(last)
	}

}

// custom function
func (inc *PriceCrossMa) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {

	// Three ways of filter on interval
	// (1)Use this PriceCrossMa indicator's interval
	if interval != inc.Interval {
		return
	}

	// (2)Use the interval of the types.IntervalWindow
	// filteriw := types.IntervalWindow{Window: 20, Interval: "5m"}
	// if interval != filteriw.Interval {
	// 	return
	// }

	// (3)Use the interval of the types.Interval
	// var i types.Interval = "5m"
	// if interval != i {
	// 	return
	// }
	inc.CalculateAndUpdate(window)
}

// required
func (inc *PriceCrossMa) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}

// 這樣外部可以以收盤價的K線來bind
func (inc *PriceCrossMa) BindK(target KLineClosedEmitter, symbol string, interval types.Interval) {
	target.OnKLineClosed(types.KLineWith(symbol, interval, inc.PushK))
}

// 在外部載入K線
func (inc *PriceCrossMa) LoadK(allKLines []types.KLine) {
	for _, k := range allKLines {
		inc.PushK(k)
	}
}
