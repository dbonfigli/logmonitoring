package commonlog

import (
	"testing"
)

type testInput struct {
	plaintextLog      string
	expectedCommonLog CommonLog
	expectingError    bool
}

var testGoodInputs = []testInput{
	testInput{
		plaintextLog: "127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif/sda HTTP/1.0\" 200 2326",
		expectedCommonLog: CommonLog{
			RemoteHost: "127.0.0.1",
			Rfc931:     "user-identifier",
			AuthUser:   "frank",
			Timestamp:  "10/Oct/2000:13:55:36 -0700",
			Request:    "GET /apache_pb.gif/sda HTTP/1.0",
			Status:     200,
			Size:       2326,
		},
		expectingError: false,
	},
	testInput{
		plaintextLog: "127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326",
		expectedCommonLog: CommonLog{
			RemoteHost: "127.0.0.1",
			Rfc931:     "user-identifier",
			AuthUser:   "frank",
			Timestamp:  "10/Oct/2000:13:55:36 -0700",
			Request:    "GET /apache_pb.gif HTTP/1.0",
			Status:     200,
			Size:       2326,
		},
		expectingError: false,
	},
	testInput{
		plaintextLog: "127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] \"POST /section/apache_pb.php HTTP/1.1\" 404 14040",
		expectedCommonLog: CommonLog{
			RemoteHost: "127.0.0.1",
			Rfc931:     "-",
			AuthUser:   "-",
			Timestamp:  "10/Oct/2000:13:55:36 -0700",
			Request:    "POST /section/apache_pb.php HTTP/1.1",
			Status:     404,
			Size:       14040,
		},
		expectingError: false,
	},
}

var testBadInputs = []testInput{
	// empty string
	testInput{
		plaintextLog:   "",
		expectingError: true,
	},
	// \n at the end of string
	testInput{
		plaintextLog:   "127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] \"POST /section/apache_pb.php HTTP/1.1\" 404 14040\n",
		expectingError: true,
	},
	// space at end of string
	testInput{
		plaintextLog:   "127.0.0.2 - - [10/Oct/2000:13:55:36 -0700] \"POST /section/apache_pb.php HTTP/1.1\" 404 14040 ",
		expectingError: true,
	},
	// double space somewhere
	testInput{
		plaintextLog:   "127.0.0.3 - - [10/Oct/2000:13:55:36  -0700] \"POST /section/apache_pb.php HTTP/1.1\" 404 14040",
		expectingError: true,
	},
	// bad status code
	testInput{
		plaintextLog:   "127.0.0.4 - - [10/Oct/2000:13:55:36  -0700] \"POST /section/apache_pb.php HTTP/1.1\" x1 14040",
		expectingError: true,
	},
	// bad size
	testInput{
		plaintextLog:   "127.0.0.5 - - [10/Oct/2000:13:55:36  -0700] \"POST /section/apache_pb.php HTTP/1.1\" 404 a2",
		expectingError: true,
	},
	//missing timestamp
	testInput{
		plaintextLog:   "127.0.0.6 - -  \"POST /section/apache_pb.php HTTP/1.1\" 404 a2",
		expectingError: true,
	},
	// bad timestamp
	testInput{
		plaintextLog:   "127.0.0.7 - - [10/Oct/2000:13:55:36 -0700 \"POST /section/apache_pb.php HTTP/1.1\" 404 14040",
		expectingError: true,
	},
	// bad timestamp
	testInput{
		plaintextLog:   "127.0.0.8 - - 10/Oct/2000:13:55:36 -0700] \"POST /section/apache_pb.php HTTP/1.1\" 404 14040",
		expectingError: true,
	},
	//bad request
	testInput{
		plaintextLog:   "127.0.0.9 - - [10/Oct/2000:13:55:36 -0700] \"POST HTTP/1.1\" 404 14040",
		expectingError: true,
	},
	// bad request
	testInput{
		plaintextLog:   "127.0.0.10 - - 10/Oct/2000:13:55:36 -0700] \"POST /section/apache_pb.php HTTP/1.1 404 14040",
		expectingError: true,
	},
	// bad request
	testInput{
		plaintextLog:   "127.0.0.11 - - 10/Oct/2000:13:55:36 -0700] POST /section/apache_pb.php HTTP/1.1\" 404 14040",
		expectingError: true,
	},
	// bad request
	testInput{
		plaintextLog:   "127.0.0.12 - - 10/Oct/2000:13:55:36 -0700] \"POST /section/apache_pb.php\" HTTP/1.1 404 14040",
		expectingError: true,
	},
}

func validateInput(input testInput, t *testing.T) {
	clog, err := Parse(input.plaintextLog)
	if err == nil && input.expectingError {
		t.Errorf("error on parsing bad log string:\n%v\nExpected parse error but got no error", input.plaintextLog)
	} else if err != nil && !input.expectingError {
		t.Errorf("error on parsing valid log string:\n%v\nExpected no parse error but got error:\n%v", input.plaintextLog, err)
	} else if err == nil && clog != input.expectedCommonLog {
		t.Errorf("error on parsing valid log string:\n%v\nExpected:\n%v\nGot:\n%v\n", input.plaintextLog, input.expectedCommonLog, clog)
	}
}

func TestGoodInputs(t *testing.T) {
	for _, testInput := range testGoodInputs {
		validateInput(testInput, t)
	}
}

func TestBadInputs(t *testing.T) {
	for _, testInput := range testBadInputs {
		validateInput(testInput, t)
	}
}
