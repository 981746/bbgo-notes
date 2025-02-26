package indicator

import (
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

//go:generate callbackgen -type PriceCrossMa
type PriceCrossMa struct {
	Values          floats.Slice
	updateCallbacks []func(value float64)
	types.IntervalWindow
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

func (inc *PriceCrossMa) Update(close64 float64) {
	// indicator calculation here...
	// push value...

	calculatedValue := close64 / 2
	inc.Values.Push(calculatedValue)
}

func (inc *PriceCrossMa) PushK(k types.KLine) {
	inc.Update(k.Close.Float64())
}

// custom function
func (inc *PriceCrossMa) CalculateAndUpdate(allKLines []types.KLine) {
	if len(inc.Values) == 0 {
		// preload or initialization
		for _, k := range allKLines {
			inc.PushK(k)

		}

		inc.EmitUpdate(inc.Last(0))
	} else {
		// update new value only
		k := allKLines[len(allKLines)-1]
		inc.PushK(k)
		inc.EmitUpdate(inc.Last(0)) // produce data, broadcast to the subscribers
	}
}

// custom function
func (inc *PriceCrossMa) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	// filter on interval
	inc.CalculateAndUpdate(window)
}

// required
func (inc *PriceCrossMa) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
