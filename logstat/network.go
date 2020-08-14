package logstat

import (
	"fmt"
	"strconv"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

// NetworkStat is the type that hold statistics about network traffic seen, implements interface LogStat.
type NetworkStat struct {
	bytes int
}

// NewInterval clears the statistic for the current interval.
func (s *NetworkStat) NewInterval() {
	s.bytes = 0
}

// Update updates the statistic with the new log information.
func (s *NetworkStat) Update(clog commonlog.CommonLog) error {
	if clog.Size < 0 {
		return fmt.Errorf("size of request reported as less than zero: %v, skipping this log for statistics", clog.Size)
	}
	s.bytes += clog.Size
	return nil
}

// String returns a string used to display this statistic.
func (s *NetworkStat) String() string {
	return "traffic seen in bytes: " + strconv.Itoa(s.bytes)
}
