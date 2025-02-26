package indicatorv2

import (
	"github.com/c9s/bbgo/pkg/datatype/bools"
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

type CustomCrossType float64

const (
	CustomCrossOver  CustomCrossType = 1.0
	CustomCrossUnder CustomCrossType = -1.0
)

// CrossStream subscribes 2 upstreams, and calculate the cross signal
type CustomCrossStream struct {
	*types.Float64Series

	a, b    floats.Slice
	IsCross bools.BoolSlice
}

// Cross creates the CrossStream object:
//
// cross := Cross(fastEWMA, slowEWMA)
func CustomCross(a, b types.Float64Source) *CustomCrossStream {
	s := &CustomCrossStream{
		Float64Series: types.NewFloat64Series(),
	}
	a.OnUpdate(func(v float64) {
		s.a.Push(v)
		s.calculate()
	})
	b.OnUpdate(func(v float64) {
		s.b.Push(v)
		s.calculate()
	})
	return s
}

func (s *CustomCrossStream) calculate() {

	if s.a.Length() != s.b.Length() {
		return
	}

	current := s.a.Last(0) - s.b.Last(0)
	previous := s.a.Last(1) - s.b.Last(1)

	if previous == 0.0 {
		return
	}

	// cross over or cross under
	if current*previous < 0 {

		s.IsCross.Push(true)
		if current > 0 {
			s.PushAndEmit(float64(CustomCrossOver))
		} else {
			s.PushAndEmit(float64(CustomCrossUnder))
		}
	} else {
		s.IsCross.Push(false)
	}

}
