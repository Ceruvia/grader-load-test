package models

type EngineRunResult struct {
	Verdict                 Verdict `json:"verdict"`
	HasErrorMessage         bool           `json:"has_error_message"`
	ErrorMessage            string         `json:"error_message"`
	InputFilename           string         `json:"input_filename"`
	OutputFilename          string         `json:"output_filename"`
	TimeToRunInMilliseconds int            `json:"time_to_run_ms"`
	MemoryUsedInKilobytes   int            `json:"memory_used_kb"`
}

type GradingResult struct {
	IsSuccess             bool              `json:"is_success"`
	Status                string            `json:"status"`
	ErrorMessage          string            `json:"error_message"`
	TestcaseGradingResult []EngineRunResult `json:"testcase_result"`
}

type Verdict struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var (
	VerdictAC  = Verdict{Name: "Accepted", Code: "AC"}
	VerdictWA  = Verdict{Name: "Wrong Answer", Code: "WA"}
	VerdictCE  = Verdict{Name: "Compilation Error", Code: "CE"}
	VerdictRE  = Verdict{Name: "Runtime Error", Code: "RE"}
	VerdictTLE = Verdict{Name: "Time Limit Exceeded", Code: "TLE"}
	VerdictXX  = Verdict{Name: "Internal Error", Code: "XX"}
)