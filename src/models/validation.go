package models

const (
	Success Result = "Success"
	Failure Result = "Failure"
)

// Result type is an enum representing all possible validation result types
type Result string

// ValidationCandidate holds data to be validated and metadata about it
type ValidationCandidate struct {
	Type string                 `json:"type"`
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
} // @Name ValidationCandidate

// ValidationResult represents the result of validating a single ValidationCandidate against a single validator
type ValidationResult struct {
	ValidatorSource string  `json:"validatorSource"`
	Result          Result  `json:"result"`
	CandidateID     string  `json:"candidateId"`
	Error           *string `json:"error,omitempty"`
} // @Name ValidationResult

// ValidationSummary represents the summary of multiple ValidationResult(s)
type ValidationSummary struct {
	Result    Result              `json:"result"`
	Successes []*ValidationResult `json:"successes"`
	Failures  []*ValidationResult `json:"failures"`
} // @Name ValidationSummary
