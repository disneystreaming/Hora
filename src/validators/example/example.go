package example

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	"github.com/disneystreaming/Hora/src/models"

	_ "embed"
)

//go:embed schema.json
var schemaBytes []byte

// ExampleValidator is a simple json schema validator
type ExampleValidator struct {
	source      string
	targetTypes map[string]interface{}
}

// NewExampleValidator returns a configured ExampleValidator
func NewExampleValidator() ExampleValidator {
	return ExampleValidator{source: "example", targetTypes: map[string]interface{}{"abc": nil, "xyz": nil}}
}

// ShouldValidate determines if the given candidate should be validated by this validator
func (v ExampleValidator) ShouldValidate(candidate models.ValidationCandidate) bool {
	_, exists := v.targetTypes[candidate.Type]
	return exists
}

// Validate validates the given candidate against this validator and returns the result to the given channel
func (v ExampleValidator) Validate(candidate models.ValidationCandidate, channel chan *models.ValidationResult) {
	ret := models.ValidationResult{ValidatorSource: v.source, Result: models.Failure, CandidateID: candidate.ID}
	if schemaBytes != nil {
		schema := gojsonschema.NewBytesLoader(schemaBytes)
		data := gojsonschema.NewGoLoader(candidate.Data)
		result, err := gojsonschema.Validate(schema, data)
		if err != nil {
			errStr := err.Error()
			ret.Error = &errStr
		} else if !result.Valid() {
			errors := []string{}
			for _, err := range result.Errors() {
				errors = append(errors, err.String())
			}
			errStr := fmt.Sprint(errors)
			ret.Error = &errStr
		} else {
			ret.Result = models.Success
		}
	}
	channel <- &ret
}
