package logstat

import (
	"math"
	"testing"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

func TestRequestRate(t *testing.T) {
	rstat := NewReqRateStat(3)
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	rstat.NewInterval()
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.NewInterval()
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})

	r := rstat.CurrentReqRate()
	var expectedr float64 = 2
	if r != expectedr {
		t.Errorf("expected request rate of %v but got %v", expectedr, r)
	}

	s := rstat.String()
	expected := "request rate: 2"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestRequestRate2(t *testing.T) {
	rstat := NewReqRateStat(3)
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	rstat.NewInterval()
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.NewInterval()
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.NewInterval()
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	s := rstat.String()
	expected := "request rate: 3"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestRequestRateNotEnoughData(t *testing.T) {
	rstat := NewReqRateStat(3)
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	rstat.NewInterval()
	rstat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	r := rstat.CurrentReqRate()
	if !math.IsNaN(r) {
		t.Errorf("CurrentReqRate method result not as expected:\nExpected:\n%v\nGot:\n%v", math.NaN(), r)
	}
	s := rstat.String()
	expected := "request rate: not enough data to show a meaningful request rate"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}
