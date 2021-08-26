package jobs

import (
	"math/rand"
	"time"

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

var callRandom *rand.Rand

func GetRandomBlockInterval(maxInterval Interval, minInterval Interval) int {
	max := maxInterval.Value()
	min := minInterval.Value()

	randomBlockNumber := callRandom.Intn(max-min) + min

	return randomBlockNumber
}

func NewRand() *rand.Rand {
	if callRandom == nil {
		rand.Seed(time.Now().UnixNano())
		someSalt := int64(rand.Intn(10000))
		saltSource := rand.NewSource(time.Now().UnixNano() * someSalt)
		saltRandom := rand.New(saltSource)
		randomSalt := saltRandom.Intn(1000000)
		BlockSource := rand.NewSource(someSalt * int64(randomSalt))
		callRandom = rand.New(BlockSource)
	}

	return callRandom
}
