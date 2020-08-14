unit_test:
	go test ./...

e2e_test:
	e2e/test.sh

extended_e2e_test:
	GENERATED_REQ_RATE_HIGH="15" \
	GENERATED_REQ_RATE_HIGH_SECONDS="130" \
	GENERATED_REQ_RATE_LOW="1" \
	GENERATED_REQ_RATE_LOW_SECONDS="60" \
	THRESHOLD_REQ_RATE_HIGH="10" \
	THRESHOLD_REQ_RATE_HIGH_SECONDS="120" \
	ANALYSIS_SECONDS="190" \
	CYCLES="1" \
	EXPECTED_HIGH_TRAFFIC_ALERTS="1" \
	EXPECTED_HIGH_TRAFFIC_RECOVERIES="1" \
	e2e/test.sh

short_test: unit_test e2e_test

test: unit_test e2e_test extended_e2e_test

build: 
	go build -v -o bin/logmonitor github.com/dbonfigli/logmonitoring/logmonitor
	go build -v -o bin/loggenerator github.com/dbonfigli/logmonitoring/loggenerator

install:
	go install github.com/dbonfigli/logmonitoring/logmonitor
	go install github.com/dbonfigli/logmonitoring/loggenerator


clean:
	rm -rf bin
	go clean -i github.com/dbonfigli/logmonitoring/logmonitor
	go clean -i github.com/dbonfigli/logmonitoring/loggenerator

all: short_test build