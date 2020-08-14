package logstat

import (
	"testing"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

func TestAlertingReqRateStat(t *testing.T) {
	a := NewAlertingReqRateStat(3, 3)
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	a.NewInterval()
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	a.NewInterval()
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})

	isFiring, changed := a.CheckAndUpdateAlert()
	if !isFiring {
		t.Errorf("expected alarm in firing state but it is not")
	}
	if !changed {
		t.Errorf("expected alarm changed compared to the previous state but it is not")
	}

	isFiring, changed = a.CheckAndUpdateAlert()
	if !isFiring {
		t.Errorf("expected alarm in firing state but it is not")
	}
	if changed {
		t.Errorf("expected alarm not changed compared to the previous state but it is")
	}
}

func TestAlertingReqRateStatWithNotEnoughData(t *testing.T) {
	a := NewAlertingReqRateStat(4, 3)
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	a.NewInterval()
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	a.NewInterval()
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})

	isFiring, changed := a.CheckAndUpdateAlert()
	if isFiring {
		t.Errorf("expected alarm not in firing state but it is")
	}

	a.NewInterval()
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	a.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	isFiring, changed = a.CheckAndUpdateAlert()
	if !isFiring {
		t.Errorf("expected alarm in firing state but it is not")
	}
	if !changed {
		t.Errorf("expected alarm changed compared to the previous state but it is not")
	}
}
