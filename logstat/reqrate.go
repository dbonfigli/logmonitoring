package logstat

import (
	"log"
	"math"
	"strconv"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

// ReqRateStat is the type that holds the statistics about the number of requests on the previous n intervals of time.
//
// ReqRateStat must be initialized with the function NewReqRateStat.
// The intervals must be of the same lenght, i.e. the calculated request rate makes sense only if the time between calls to NewInterval() is constant.
type ReqRateStat struct {

	//requests is a slice of ints, each element the number of requests made in a given interval.
	requests []int

	//intevals is the number of intervals we want to keep tabs on.
	intervals int
}

// NewReqRateStat must be used to initialize the ReqRateStat struct.
func NewReqRateStat(intervals int) *ReqRateStat {
	if intervals < 1 {
		log.Fatal("n cannot be < 1!")
	}
	return &ReqRateStat{
		requests:  make([]int, 1),
		intervals: intervals,
	}
}

// NewInterval signals the start of a new interval.
func (s *ReqRateStat) NewInterval() {
	//add the placeholder for the new interval
	s.requests = append(s.requests, 0)

	//remove old interval data
	if len(s.requests) > s.intervals {
		s.requests = s.requests[1:]
	}
}

// Update updates the statistic with the new log line.
func (s *ReqRateStat) Update(clog commonlog.CommonLog) error {
	s.requests[len(s.requests)-1]++
	return nil
}

// Requests return the total number of seen request on the previous intervals.
func (s *ReqRateStat) Requests() int {
	sum := 0
	for _, r := range s.requests {
		sum += r
	}
	return sum
}

// CurrentReqRate return the request rate over the intervals.
func (s *ReqRateStat) CurrentReqRate() float64 {
	if len(s.requests) < s.intervals {
		return math.NaN()
	}

	reqRate := float64(s.Requests()) / float64(s.intervals)
	return reqRate
}

// String returns a string used to display the request rate over the intervals approximated to the nearest interger.
func (s *ReqRateStat) String() string {
	reqRate := s.CurrentReqRate()
	if math.IsNaN(reqRate) {
		return "request rate: not enough data to show a meaningful request rate"
	}
	return "request rate: " + strconv.FormatFloat(reqRate, 'f', 0, 32)
}
