package tasks

import (
	"fmt"

	"github.com/Ceruvia/grader-load-test/internal/models"
)

type SampleLanguageSubmission struct {
	id                 string
	graderURL          string
	submissionURL      string
	inputTestcases     []string
	outputTestcases    []string
	timeLimit          int
	memoryLimit        int
	language           string
	mainSourceFilename string
	expectedResponse   models.GradingResult
}

type SampleBuilderSubmission struct {
	id               string
	graderURL        string
	submissionURL    string
	inputTestcases   []string
	outputTestcases  []string
	timeLimit        int
	memoryLimit      int
	language         string
	compileScript    string
	runScript        string
	expectedResponse models.GradingResult
}

var (
	BuilderSamples []SampleBuilderSubmission = []SampleBuilderSubmission{
		{
			id:               "c_1_reverseprime",
			graderURL:        "https://pub-aa14e9fb26a94974a23c01cf74108727.r2.dev/c_1_reverseprime_grader.zip",
			submissionURL:    "https://pub-aa14e9fb26a94974a23c01cf74108727.r2.dev/c_1_reverseprime_submission.zip",
			inputTestcases:   createInputTestcases(10),
			outputTestcases:  createOutputTestcases(10),
			timeLimit:        1000,
			memoryLimit:      10240,
			language:         "C",
			compileScript:    "compile",
			runScript:        "prog",
			expectedResponse: createExpectedResult(true, "Success", "", []string{"AC", "AC", "AC", "AC", "AC", "AC", "WA", "AC", "WA", "WA"}),
		},
		{
			id:               "java_1_inventory_management",
			graderURL:        "https://pub-aa14e9fb26a94974a23c01cf74108727.r2.dev/java_1_inventory_management_grader.zip",
			submissionURL:    "https://pub-aa14e9fb26a94974a23c01cf74108727.r2.dev/java_1_inventory_management_submission.zip",
			inputTestcases:   createInputTestcases(5),
			outputTestcases:  createOutputTestcases(5),
			timeLimit:        1000,
			memoryLimit:      10240,
			language:         "Java",
			compileScript:    "Main.class",
			runScript:        "Main",
			expectedResponse: createExpectedResult(true, "Success", "", []string{"AC", "AC", "AC", "AC", "AC"}),
		},
	}

	LanguageSamples []SampleLanguageSubmission = []SampleLanguageSubmission{
		{
			id:                 "python_1_representasi",
			graderURL:          "https://pub-aa14e9fb26a94974a23c01cf74108727.r2.dev/python_1_representasi_grader.zip",
			submissionURL:      "https://pub-aa14e9fb26a94974a23c01cf74108727.r2.dev/python_1_representasi_submission.zip",
			inputTestcases:     createInputTestcases(10),
			outputTestcases:    createOutputTestcases(10),
			timeLimit:          1000,
			memoryLimit:        10240,
			language:           "Python 3",
			mainSourceFilename: "unique_digit.py",
			expectedResponse:   createExpectedResult(true, "Success", "", []string{"AC", "AC", "AC", "AC", "AC", "AC", "AC", "AC", "AC", "AC"}),
		},
	}
)

func createInputTestcases(count int) []string {
	var inputTestcases []string
	for i := range count {
		inputTestcases = append(inputTestcases, fmt.Sprintf("%d.in", i+1))
	}
	return inputTestcases
}

func createOutputTestcases(count int) []string {
	var inputTestcases []string
	for i := range count {
		inputTestcases = append(inputTestcases, fmt.Sprintf("%d.out", i+1))
	}
	return inputTestcases
}

func createExpectedResult(isSuccess bool, status, errorMessage string, verdicts []string) models.GradingResult {
	var verdictToEngineRunResult []models.EngineRunResult
	for i, verdict := range verdicts {
		var v models.Verdict
		switch verdict {
		case "AC":
			v = models.VerdictAC
		case "RE":
			v = models.VerdictRE
		case "WA":
			v = models.VerdictWA
		case "CE":
			v = models.VerdictCE
		case "TLE":
			v = models.VerdictTLE
		default:
			v = models.VerdictXX
		}

		verdictToEngineRunResult = append(verdictToEngineRunResult, models.EngineRunResult{
			Verdict:        v,
			InputFilename:  fmt.Sprintf("%d.in", i+1),
			OutputFilename: fmt.Sprintf("%d.out", i+1),
		})
	}

	return models.GradingResult{
		IsSuccess:             isSuccess,
		Status:                status,
		ErrorMessage:          errorMessage,
		TestcaseGradingResult: verdictToEngineRunResult,
	}
}
