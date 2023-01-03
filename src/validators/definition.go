// package validators defines the Validator interface and other common entities used across validators
package validators

import (
	"github.com/disneystreaming/Hora/src/models"
)

// Validator is what all module validators need to implement to be allowed in the suite
type Validator interface {
	// ShouldValidate determines if the given candidate should be validated by this validator
	ShouldValidate(models.ValidationCandidate) bool
	// Validate validates the given candidate against this validator and returns the result to the given channel
	Validate(models.ValidationCandidate, chan *models.ValidationResult)
}
