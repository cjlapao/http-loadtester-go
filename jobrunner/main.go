package jobrunner

import "github.com/cjlapao/common-go/log"

var logger = log.Get()

type Interval struct {
	value int
}

func (s Interval) Value() int {
	return s.value
}

func NewInterval(value int) Interval {
	return Interval{value: value}
}
