package logstat

import (
	"testing"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

func TestNetworkStatTotal(t *testing.T) {
	var networkStat NetworkStat
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	s := networkStat.String()
	expected := "traffic seen in bytes: 60"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestNetworkStatSkipNegative(t *testing.T) {
	var networkStat NetworkStat
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, -30})
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, -40})
	s := networkStat.String()
	expected := "traffic seen in bytes: 10"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestNetworkStatNewInterval(t *testing.T) {
	var networkStat NetworkStat
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 10})
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 20})
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 30})
	networkStat.NewInterval()
	networkStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /on/apache_pb.gif HTTP/1.0", 200, 10})
	s := networkStat.String()
	expected := "traffic seen in bytes: 10"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}
