package commonlog

import (
	"errors"
	"strconv"
	"strings"
)

// CommonLog is a structured representation of a log line of the "common log format" as defined in
// https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format.
type CommonLog struct {
	RemoteHost string
	Rfc931     string
	AuthUser   string
	Timestamp  string
	Request    string
	Status     int
	Size       int
}

// Parse return a CommonLog structure (i.e. structured information) from a plaintext string that represent
// a log of the common log format.
//
// Parsing a textual log is useful to validate the log format and obtain basic information to be later used as stats.
//
// Returns a non nil error in case parsing was not possible, in this case the returned CommonLog must be discarded.
func Parse(plaintextLog string) (CommonLog, error) {
	var clog CommonLog
	fields := strings.Split(plaintextLog, " ")

	if len(fields) != 10 {
		return clog, errors.New("log format not valid, not all fields are defined")
	}

	clog.RemoteHost = fields[0]

	clog.Rfc931 = fields[1]

	clog.AuthUser = fields[2]

	if len(fields[3]) < 1 || len(fields[4]) < 1 || fields[3][0] != '[' || fields[4][len(fields[4])-1] != ']' {
		return clog, errors.New("log format not valid, date field is invalid")
	}
	clog.Timestamp = fields[3][1:] + " " + fields[4][:len(fields[4])-1]

	if len(fields[5]) < 1 || len(fields[7]) < 1 || fields[5][0] != '"' || fields[7][len(fields[7])-1] != '"' {
		return clog, errors.New("log format not valid, request field is invalid")
	}

	clog.Request = fields[5][1:] + " " + fields[6] + " " + fields[7][:len(fields[7])-1]

	status, err := strconv.Atoi(fields[8])
	if err != nil {
		return clog, errors.New("log format not valid, cannot convert status code field to integer: " + fields[8])
	}
	clog.Status = status

	size, err := strconv.Atoi(fields[9])
	if err != nil {
		return clog, errors.New("log format not valid, cannot convert size field to integer: " + fields[9])
	}
	clog.Size = size

	return clog, nil

}
