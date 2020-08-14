package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	logFilePath := flag.String("f", "/tmp/access.log", "location of the log file where to generate logs")
	highTrafficReqRate := flag.Int("R", 15, "request/seconds for the \"high traffic\" phase")
	lowTrafficReqRate := flag.Int("r", 1, "request/seconds for the \"low traffic\" phase")
	highTrafficLenght := flag.Int("T", 125, "lenght in seconds of the \"high traffic\" phase")
	lowTrafficLenght := flag.Int("t", 10, "lenght in seconds of the \"low traffic\" phase")
	cycles := flag.Int("c", 3, "number of cycles of the high -> low traffic phases")
	printHelp := flag.Bool("h", false, "print this help")
	flag.Parse()
	if *printHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *highTrafficReqRate > 20 || *lowTrafficReqRate > 20 {
		fmt.Println("warning: request rate is fairly high, probably the log generation will not be accurate")
	}

	waitMsecHigh := 1000 / *highTrafficReqRate
	repetitionsHigh := *highTrafficLenght * *highTrafficReqRate

	waitMsecLow := 1000 / *lowTrafficReqRate
	repetitionsLow := *lowTrafficLenght * *lowTrafficReqRate

	file, err := os.Create(*logFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < *cycles; i++ {
		generateTraffic(file, waitMsecHigh, repetitionsHigh)
		generateTraffic(file, waitMsecLow, repetitionsLow)
	}
}

func generateTraffic(file *os.File, waitMsec int, repetitions int) {
	fmt.Printf("starting to generate logs every %d milliseconds for %d times\n", waitMsec, repetitions)
	for i := 0; i < repetitions; i++ {
		logLine := generateNewLog()
		_, err := file.WriteString(logLine)
		if err != nil {
			log.Fatal(err)
		}
		file.Sync()
		time.Sleep(time.Duration(waitMsec) * time.Millisecond)
	}
}

func generateNewLog() string {
	now := time.Now().Format("2/Jan/2006:15:04:05 -0700")
	strs := []string{
		"127.0.0.1 - frank [" + now + "] \"GET /section1/apache_pb.gif HTTP/1.0\" 200 10\n",
		"127.0.0.2 - nikki [" + now + "] \"GET /section2/apache_pb.gif HTTP/1.0\" 500 10\n",
	}
	return strs[rand.Intn(2)]
}
