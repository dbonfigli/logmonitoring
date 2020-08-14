package logstat

import (
	"strconv"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

// ReqStat is the type that hold statistics about requests seen in an interval, implements interface LogStat.
type ReqStat struct {
	n int
}

// Update updates the statistic with the new log information.
func (s *ReqStat) Update(clog commonlog.CommonLog) error {
	s.n++
	return nil
}

// NewInterval clears the statistic for the current interval.
func (s *ReqStat) NewInterval() {
	s.n = 0
}

// String returns a string used to display this statistic.
func (s *ReqStat) String() string {
	return "request seen: " + strconv.Itoa(s.n)
}
