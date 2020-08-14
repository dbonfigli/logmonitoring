package logstat

import (
	"github.com/dbonfigli/logmonitoring/commonlog"
)

// Stat is the interface that is used to keep tabs on log statistics on a given inteval.
type Stat interface {
	// NewInterval prepares the statistics to hold information for a new interval.
	NewInterval()

	// Update updates the statistics for the current time interval using the information in the new log line "clog".
	// If there are problems on input (e.g. format not valid) return an error.
	Update(clog commonlog.CommonLog) error

	// String returns a string rapresenting the information about the statistic for the current interval, ready to be displayed.
	String() string
}

// Alert represents information about an alarm.
type Alert interface {
	// CheckAndUpdateAlert checks for the alert condition.
	// It returns the first value as true if the alert is firing, the second value as true if the alert changed status compared to the previous check.
	CheckAndUpdateAlert() (bool, bool)
}
