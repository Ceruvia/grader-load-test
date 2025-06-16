package tests

import (
	"strings"

	"github.com/Ceruvia/grader-load-test/internal/models"
)

func assertGradingResult(got, want models.GradingResult) bool {
	if got.IsSuccess != want.IsSuccess {
		return false
	}

	if !strings.Contains(got.ErrorMessage, want.ErrorMessage) {
		return false
	}

	if got.Status != "Compile Error" {
		if len(got.TestcaseGradingResult) != len(want.TestcaseGradingResult) {
			return false
		}

		for i, _ := range got.TestcaseGradingResult {
			if got.TestcaseGradingResult[i].Verdict != want.TestcaseGradingResult[i].Verdict {
				return false
			}
		}
	}

	return true
}
