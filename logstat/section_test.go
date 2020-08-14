package logstat

import (
	"testing"

	"github.com/dbonfigli/logmonitoring/commonlog"
)

func TestSectionStatLessSections(t *testing.T) {
	sectionStat := NewSectionStat(10)
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 2326})
	s := sectionStat.String()
	expected := "site hits by section: /: 2; /as: 1"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

// test when sections to show are less than actual sections
func TestSectionStatMoreSections(t *testing.T) {
	sectionStat := NewSectionStat(1)
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 2326})
	s := sectionStat.String()
	expected := "site hits by section: /: 2"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestSectionStatBadSections(t *testing.T) {
	sectionStat := NewSectionStat(10)
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET / HTTP/1.0", 200, 2326})
	s := sectionStat.String()
	expected := "site hits by section: /: 2; /as: 1"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestSectionStatMultipleSections(t *testing.T) {
	sectionStat := NewSectionStat(10)
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /first/second/apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /first/third/ HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET / HTTP/1.0", 200, 2326})
	s := sectionStat.String()
	expected := "site hits by section: /: 2; /first: 2"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}

func TestSectionStatNewInterval(t *testing.T) {
	sectionStat := NewSectionStat(10)
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /as/apache_pb.gif HTTP/1.0", 200, 2326})
	sectionStat.NewInterval()
	sectionStat.Update(commonlog.CommonLog{"127.0.0.1", "user-identifier", "frank", "10/Oct/2000:13:55:36 -0700", "GET /on/apache_pb.gif HTTP/1.0", 200, 2326})
	s := sectionStat.String()
	expected := "site hits by section: /on: 1"
	if s != expected {
		t.Errorf("String method result not as expected:\nExpected:\n%v\nGot:\n%v", expected, s)
	}
}
