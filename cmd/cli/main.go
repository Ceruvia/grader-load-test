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
	useBuilder := flag.Int("builder", 0, "1 if builder is used, -1 if not")
	problemIdx := flag.Int("problem", 1, "Index of problem being tested")

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

	case "compare":
		missing := []string{}
		if *useBuilder == 0 {
			missing = append(missing, "=builder <bool>")
		}
		if *simCount == 0 {
			missing = append(missing, "=-count <n>")
		}
		if *problemIdx == 0 {
			missing = append(missing, "=problemIdx <idx>")
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
	case "compare":
		start := time.Now()
		if *useBuilder == -1 {
			ComparisonTest(server, *simCount, *problemIdx, false)
		} else {
			ComparisonTest(server, *simCount, *problemIdx, true)
		}
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

func ComparisonTest(server *machinery.Server, total, problemIdx int, useBuilder bool) {

	var wg sync.WaitGroup

	testcaseTimes := [][]int{}
	fmt.Println(total, problemIdx, useBuilder)

	for i := 1; i <= total; i++ {
		wg.Add(1)
		go func(taskNum int) {
			defer wg.Done()
			fmt.Printf("[%d] Sending request...\n", taskNum)
			accepted, result, duration := tests.TestBlackbox(server, useBuilder, problemIdx)

			if accepted {
				testcaseTime := []int{}
				for _, tcRes := range result.TestcaseGradingResult {
					testcaseTime = append(testcaseTime, tcRes.TimeToRunInMilliseconds)
				}
				testcaseTimes = append(testcaseTimes, testcaseTime)
			}

			fmt.Printf("[%d] Finished: OK=%t, time=%s\n", taskNum, accepted, duration)
		}(i)
	}

	wg.Wait()

	cols := len(testcaseTimes[0])
	sums := make([]int, cols)

	for _, row := range testcaseTimes {
		for i, val := range row {
			sums[i] += val
		}
	}

	means := make([]float64, cols)
	for i, sum := range sums {
		means[i] = float64(sum) / float64(len(testcaseTimes))
	}

	fmt.Printf("%+v\n", testcaseTimes)

	fmt.Println("Means:", means)

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
