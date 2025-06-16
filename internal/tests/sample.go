package tests

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Ceruvia/grader-load-test/internal/models"
	"github.com/Ceruvia/grader-load-test/internal/tasks"
	"github.com/RichardKnop/machinery/v2"
)

func TestSampleBlackbox(server *machinery.Server, useBuilder bool, sampleIdx int) (bool, time.Duration) {
	signature := tasks.GetBlackboxSignature(useBuilder, sampleIdx)

	start := time.Now()
	asyncResult, err := server.SendTask(signature)
	elapsed := time.Since(start)

	if err != nil {
		panic(err)
	}

	reflectResult, _ := asyncResult.Get(10)

	if err != nil {
		panic(err)
	}

	tmp := reflectResult[0]

	var result models.GradingResult

	err = json.Unmarshal([]byte(tmp.String()), &result)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return assertGradingResult(result, tasks.GetBlackboxExpectedResult(useBuilder, sampleIdx)), elapsed
}
