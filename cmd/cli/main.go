package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	appConfig "github.com/Ceruvia/grader-load-test/internal/config"
	"github.com/Ceruvia/grader-load-test/internal/tests"
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"

	amqpbackend "github.com/RichardKnop/machinery/v2/backends/amqp"
	amqpbroker "github.com/RichardKnop/machinery/v2/brokers/amqp"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
)

func main() {
	// CLI flags
	testType := flag.String("test", "", "Type of test to run: simultaneous or real (required)")
	simCount := flag.Int("count", 1500, "Number of requests for simultaneous test")
	realUsers := flag.Int("users", 0, "Number of virtual users for real-condition test (required for real)")
	lowerBound := flag.Int("lower", 0, "Lower bound delay in seconds for real-condition test (required for real)")
	upperBound := flag.Int("upper", 0, "Upper bound delay in seconds for real-condition test (required for real)")
	flag.Parse()

	// Validate required flags
	if *testType == "" {
		fmt.Println("error: -test flag is required")
		flag.Usage()
		os.Exit(1)
	}

	switch *testType {
	case "simultaneous":
		// nothing extra required

	case "real":
		missing := []string{}
		if *realUsers <= 0 {
			missing = append(missing, "-users <n>")
		}
		if *lowerBound <= 0 {
			missing = append(missing, "-lower <seconds>")
		}
		if *upperBound <= 0 {
			missing = append(missing, "-upper <seconds>")
		}
		if len(missing) > 0 {
			fmt.Printf("error: missing required flags for real test: %v\n", missing)
			flag.Usage()
			os.Exit(1)
		}

	default:
		fmt.Printf("error: unknown test type: %s\n", *testType)
		flag.Usage()
		os.Exit(1)
	}

	// Setup Machinery server
	cfg := appConfig.GetAppConfig()
	cnf := &config.Config{
		Broker:          cfg.MachineryCfg.BrokerURL,
		DefaultQueue:    cfg.MachineryCfg.QueueName,
		ResultBackend:   cfg.MachineryCfg.ResultBackendURL,
		ResultsExpireIn: cfg.MachineryCfg.ResultsExpireIn,
		AMQP: &config.AMQPConfig{
			Exchange:      "machinery_exchange",
			ExchangeType:  "direct",
			BindingKey:    "machinery_task",
			PrefetchCount: 3,
		},
	}

	broker := amqpbroker.New(cnf)
	backend := amqpbackend.New(cnf)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	switch *testType {
	case "simultaneous":
		start := time.Now()
		SimultaneousLoadTest(server, *simCount)
		fmt.Printf("Total elapsed time for simultaneous test: %s\n", time.Since(start))

	case "real":
		start := time.Now()
		SimulateRealCondition(server, *realUsers, *lowerBound, *upperBound)
		fmt.Printf("Total elapsed time for real-condition test: %s\n", time.Since(start))
	}
}

// SimultaneousLoadTest sends `total` concurrent requests as fast as possible.
func SimultaneousLoadTest(server *machinery.Server, total int) {
	var wg sync.WaitGroup

	for i := 1; i <= total; i++ {
		wg.Add(1)
		go func(taskNum int) {
			defer wg.Done()
			fmt.Printf("[%d] Sending request...\n", taskNum)
			accepted, duration := tests.TestSampleBlackbox(server, true, 0)
			fmt.Printf("[%d] Finished: OK=%t, time=%s\n", taskNum, accepted, duration)
		}(i)
	}

	wg.Wait()
}

// SimulateRealCondition runs `numUsers` virtual users sending requests
// with random delays between lower and upper seconds.
func SimulateRealCondition(server *machinery.Server, numUsers, lower, upper int) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	for u := 1; u <= numUsers; u++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// Random delay between lower and upper bounds
			delta := upper - lower
			d := lower
			if delta > 0 {
				d = lower + rand.Intn(delta+1)
			}
			time.Sleep(time.Duration(d) * time.Second)

			fmt.Printf("[User %d] Sending request after %ds...\n", userID, d)
			accepted, duration := tests.TestSampleBlackbox(server, true, 0)
			fmt.Printf("[User %d] Finished: OK=%t, time=%s\n", userID, accepted, duration)
		}(u)
	}

	wg.Wait()
}
