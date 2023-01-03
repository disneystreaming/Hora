package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/disneystreaming/Hora/src/models"
	"github.com/disneystreaming/Hora/src/validators"
	"github.com/disneystreaming/Hora/src/validators/example"
)

const (
	// sample context timeout
	timeoutSeconds int = 60
)

// Validate asynchronously validates the given candidates against the configured validators and returns a
// ValidationSummary
func Validate(ctx context.Context, candidates []models.ValidationCandidate) (*models.ValidationSummary, error) {
	// in production this configuration should be brought into a dynamic setting
	configuredValidators := []validators.Validator{example.NewExampleValidator()}
	channel := make(chan *models.ValidationResult)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// runCount keeps track of how many validation runs were actually kicked off
	runCount := 0
	for _, candidate := range candidates {
		for _, validator := range configuredValidators {
			if validator.ShouldValidate(candidate) {
				runCount++
				go validator.Validate(candidate, channel)
			}
		}
	}

	summary := models.ValidationSummary{
		Result:    models.Failure,
		Successes: []*models.ValidationResult{},
		Failures:  []*models.ValidationResult{},
	}
	// runTracker keeps track of how many validation runs have returned
	runTracker := 0
	for runTracker < runCount {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("validation exceeded timeout")
		case result := <-channel:
			if result.Result == models.Success {
				summary.Successes = append(summary.Successes, result)
			} else {
				summary.Failures = append(summary.Failures, result)
			}
		}
		runTracker++
	}
	if len(summary.Failures) == 0 {
		summary.Result = models.Success
	}

	return &summary, nil
}
