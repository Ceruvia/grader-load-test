package tasks

import (
	"github.com/Ceruvia/grader-load-test/internal/models"
	"github.com/RichardKnop/machinery/v2/tasks"
)

func GetBlackboxSignature(useBuilder bool, sampleIdx int) *tasks.Signature {
	if useBuilder {
		sampleTask := BuilderSamples[sampleIdx]
		return &tasks.Signature{
			Name: "blackbox_with_builder",
			Args: []tasks.Arg{
				{Type: "string", Value: sampleTask.id},
				{Type: "string", Value: sampleTask.graderURL},
				{Type: "string", Value: sampleTask.submissionURL},
				{Type: "[]string", Value: sampleTask.inputTestcases},
				{Type: "[]string", Value: sampleTask.outputTestcases},
				{Type: "int", Value: sampleTask.timeLimit},
				{Type: "int", Value: sampleTask.memoryLimit},
				{Type: "string", Value: sampleTask.language},
				{Type: "string", Value: "Makefile"},
				{Type: "string", Value: sampleTask.compileScript},
				{Type: "string", Value: sampleTask.runScript},
			},
		}
	} else {
		sampleTask := LanguageSamples[sampleIdx]
		return &tasks.Signature{
			Name: "blackbox",
			Args: []tasks.Arg{
				{Type: "string", Value: sampleTask.id},
				{Type: "string", Value: sampleTask.graderURL},
				{Type: "string", Value: sampleTask.submissionURL},
				{Type: "[]string", Value: sampleTask.inputTestcases},
				{Type: "[]string", Value: sampleTask.outputTestcases},
				{Type: "int", Value: sampleTask.timeLimit},
				{Type: "int", Value: sampleTask.memoryLimit},
				{Type: "string", Value: sampleTask.language},
				{Type: "string", Value: sampleTask.mainSourceFilename},
			},
		}
	}
}

func GetBlackboxExpectedResult(useBuilder bool, sampleIdx int) models.GradingResult {
	if useBuilder {
		return BuilderSamples[sampleIdx].expectedResponse
	} else {
		return LanguageSamples[sampleIdx].expectedResponse
	}
}
