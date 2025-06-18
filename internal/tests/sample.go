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
		log.Printf("ERROR IN GRADING: %s", err.Error())
		return false, elapsed
	}

	reflectResult, _ := asyncResult.Get(10)

	if err != nil {
		log.Printf("ERROR IN GRADING: %s", err.Error())
		return false, elapsed
	}

	tmp := reflectResult[0]

	var result models.GradingResult

	err = json.Unmarshal([]byte(tmp.String()), &result)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return assertGradingResult(result, tasks.GetBlackboxExpectedResult(useBuilder, sampleIdx)), elapsed
}

func TestBlackbox(server *machinery.Server, useBuilder bool, sampleIdx int) (bool, models.GradingResult, time.Duration) {
	signature := tasks.GetBlackboxSignature(useBuilder, sampleIdx)

	start := time.Now()
	asyncResult, err := server.SendTask(signature)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("ERROR IN GRADING: %s", err.Error())
		return false, models.GradingResult{}, elapsed
	}

	reflectResult, _ := asyncResult.Get(10)

	if err != nil {
		log.Printf("ERROR IN GRADING: %s", err.Error())
		return false, models.GradingResult{}, elapsed
	}

	tmp := reflectResult[0]

	var result models.GradingResult

	err = json.Unmarshal([]byte(tmp.String()), &result)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return assertGradingResult(result, tasks.GetBlackboxExpectedResult(useBuilder, sampleIdx)), result, elapsed
}
