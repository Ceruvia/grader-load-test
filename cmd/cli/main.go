package main

import (
	"fmt"
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

	start := time.Now()
	SimultaneousLoadTest(server)
	elapsed := time.Since(start)

	fmt.Printf("Total elapsed time to run all: %s\n", elapsed)
}

func SimultaneousLoadTest(server *machinery.Server) {
	var wg sync.WaitGroup
	total := 100

	for i := 1; i <= total; i++ {
		wg.Add(1)
		go func(taskNum int) {
			defer wg.Done()
			fmt.Printf("[%d] Sending request\n", taskNum)
			accepted, timeElapsed := tests.TestSampleBlackbox(server, true, 0)
			fmt.Printf("[%d] Finish request, OK:%t, time:%s\n", taskNum, accepted, timeElapsed)
		}(i)
	}

	wg.Wait()
}
