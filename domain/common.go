package domain

import (
	"github.com/cjlapao/http-loadtester-go/common"
)

// Interval Entity
type Interval struct {
	value int
}

// Value Gets an interval value
func (s Interval) Value() int {
	return s.value
}

// NewInterval Creates a new interval value
func NewInterval(value int) Interval {
	return Interval{value: value}
}

// ResponseDetails Entity
type ResponseDetails struct {
	IP            string
	TLSCipher     string
	TLSVersion    string
	TLSServerName string
	Body          string
}

func GetRandomBlockInterval(maxInterval Interval, minInterval Interval) int {
	max := maxInterval.Value()
	min := minInterval.Value()

	return common.GetRandomNum(min, max)
}
