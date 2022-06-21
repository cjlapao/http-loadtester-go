package jobs

import (
	"crypto/rand"
	"math/big"

	"github.com/cjlapao/common-go/log"
)

var logger = log.Get()

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

	return GetRandomNum(min, max)
}

func GetRandomNum(min, max int) int {
	bg := big.NewInt(int64(max) - int64(min))

	n, err := rand.Int(rand.Reader, bg)

	if err != nil {
		return 0
	}

	return int(n.Int64() + int64(min))
}
