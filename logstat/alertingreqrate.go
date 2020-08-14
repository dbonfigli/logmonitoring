package logstat

// AlertingReqRateStat is the struct holding the information about the request rate statistic and the alert on it
// NewAlertingReqRateStat must be initialized with the function NewAlertingReqRateStat.
type AlertingReqRateStat struct {
	// ReqRateStat is the embedded type representing the request rate statistic
	*ReqRateStat
	// threshold is the request rate above which the alert is in firing state
	threshold float64
	// status is true if the alarm was firing last time it was checked
	status bool
}

// NewAlertingReqRateStat create and initialize a new AlertingReqRateStat struct
func NewAlertingReqRateStat(intervals int, threshold float64) *AlertingReqRateStat {
	return &AlertingReqRateStat{
		ReqRateStat: NewReqRateStat(intervals),
		threshold:   threshold,
		status:      false,
	}
}

// CheckAndUpdateAlert check if the request rate over several intervals is above a predefined threshold and update the current alert status.
// The first return value represents the alert current status after the check.
// The second return value is true if the alert changed its status, i.e. was off and now on or vice versa.
func (a *AlertingReqRateStat) CheckAndUpdateAlert() (bool, bool) {
	changed := false
	firing := a.CurrentReqRate() > a.threshold // if a.CurrentReqRate() is NaN then firing is always false
	if firing != a.status {
		changed = true
	}
	a.status = firing
	return firing, changed
}
