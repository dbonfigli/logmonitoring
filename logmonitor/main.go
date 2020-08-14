package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dbonfigli/logmonitoring/commonlog"
	"github.com/dbonfigli/logmonitoring/logstat"
)

func main() {

	// command flags
	logFilePath := flag.String("f", "/tmp/access.log", "location of the log file to analyze")
	statInterval := flag.Int("i", 10, "interval, in seconds, after which to display log statistics")
	requestRateThreshold := flag.Float64("r", 10, "threshold, in request/seconds, to be used for traffic alerting")
	requestRateInterval := flag.Int("t", 120, "seconds the traffic must be above threshold to trigger the alert")
	printHelp := flag.Bool("h", false, "print this help")
	flag.Parse()
	if *printHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// channel to signal we are ready to generate statistics
	infoTickerChan := time.NewTicker(time.Duration(*statInterval) * time.Second).C

	// initialize all the structs that will hold info statistics
	networkStat := new(logstat.NetworkStat)
	sectionStat := logstat.NewSectionStat(5)
	reqStat := new(logstat.ReqStat)
	logInfoStats := make([]logstat.Stat, 0)
	logInfoStats = append(logInfoStats, reqStat, networkStat, sectionStat)

	// channel to signal we want to check for alerts, i.e. every second
	alertTickerChan := time.NewTicker(1 * time.Second).C

	// initialize the struct that will hold alerting statistics
	alertingReqRateStat := logstat.NewAlertingReqRateStat(*requestRateInterval, *requestRateThreshold)

	// channel where logs lines will be put
	logChan := make(chan string, 1)

	// start the goroutine to countinuously read log lines from file and put logs in channel
	go readLogLines(*logFilePath, logChan)

	// for every given moment, heither read a new log line, generate stat or show alerts, never do that concurrently
	for {
		select {
		case logLine := <-logChan:
			// there is a new log line, update statistics and alerting
			clog, err := commonlog.Parse(logLine)
			if err != nil {
				log.Println(err)
				break
			}
			updateInfoStats(clog, logInfoStats)
			updateAlertingStats(clog, alertingReqRateStat)
		case <-infoTickerChan:
			// show the stats for the past interval
			printInfoStats(logInfoStats, *statInterval)
			newIntervalInfoStats(logInfoStats)
		case <-alertTickerChan:
			// check alerts and if needed show alerting information
			checkAlerts(alertingReqRateStat, *requestRateInterval)
			newIntervalAlertingStats(alertingReqRateStat)
		}
	}
}

// readLogLines continuously read complete lines (i.e a string of character terminating with \n) from the end of a file and
// put them in the channel logChan, if EOF is reached it polls the file until the line is complete.
// This function never returns.
func readLogLines(filePath string, logChan chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// start to read logs from the end of file
	if _, err := file.Seek(0, 2); err != nil {
		log.Fatal(err)
	}

	// continuously read a line on the file, wait and poll if the last character is not \n
	reader := bufio.NewReader(file)
	buf := ""
	for {
		// handle file renaming
		curFileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}
		if !os.SameFile(curFileInfo, fileInfo) {
			log.Println("info: log file was replaced")
			file.Close()
			fileInfo = curFileInfo
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			reader = bufio.NewReader(file)
		}

		//TODO: handle file truncation

		l, err := reader.ReadString('\n')
		line := string(l)
		if err != nil {
			if err == io.EOF {
				buf += line
				time.Sleep(200 * time.Millisecond)
			} else {
				break
			}
		} else {
			wholeLine := buf + line
			//remove trailing \n
			wholeLine = wholeLine[:len(wholeLine)-1]
			buf = ""
			logChan <- wholeLine
		}
	}
}

// updateInfoStats updates any logStat struct in stats with the new clog log information
func updateInfoStats(clog commonlog.CommonLog, stats []logstat.Stat) {
	for _, stat := range stats {
		err := stat.Update(clog)
		if err != nil {
			log.Println(err)
		}
	}
}

// updateAlertingStats updates any alertingStat struct (currently only alertingReqRateStat) with the new clog log information
func updateAlertingStats(clog commonlog.CommonLog, alertingReqRateStat *logstat.AlertingReqRateStat) {
	alertingReqRateStat.Update(clog)
}

// printInfoStats prints the log statistics in stats for the current interval
func printInfoStats(stats []logstat.Stat, interval int) {
	s := make([]string, 0)
	for _, stat := range stats {
		s = append(s, stat.String())
	}
	fmt.Printf("%ds stats >> %v\n", interval, strings.Join(s, " | "))
}

// newIntervalInfoStats signals the log statistics in stats that a new interval is beginning
func newIntervalInfoStats(stats []logstat.Stat) {
	for _, stat := range stats {
		stat.NewInterval()
	}
}

// checkAlerts prints alerting messages when an alert status has changed
func checkAlerts(alertingReqRateStat *logstat.AlertingReqRateStat, intervals int) {
	isFiring, changed := alertingReqRateStat.CheckAndUpdateAlert()
	requests := alertingReqRateStat.Requests()
	currentReqRate := alertingReqRateStat.CurrentReqRate()
	now := time.Now().Format("15:04:05")
	if changed {
		if isFiring {
			fmt.Printf("===> High traffic generated an alert - hits = %d in %ds (req/s: %.2f), triggered at %v <===\n", requests, intervals, currentReqRate, now)
		} else {
			fmt.Printf("===> High traffic alert recovered - hits = %d in %ds (req/s: %.2f), recovered at %v <===\n", requests, intervals, currentReqRate, now)
		}
	}
}

func newIntervalAlertingStats(alertingReqRateStat *logstat.AlertingReqRateStat) {
	alertingReqRateStat.NewInterval()
}
