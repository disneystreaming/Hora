package example

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/disneystreaming/Hora/src/models"
)

const (
	// schema keys
	missionKey string = "mission"
	crewKey    string = "crew"
	rocketKey  string = "rocket"

	validatorSource string = "example"
)

// sharedValidator to be used across tests
var sharedValidator ExampleValidator = NewExampleValidator()

func TestShouldValidate(t *testing.T) {
	testCases := []struct {
		candidate models.ValidationCandidate
		expected  bool
	}{
		{
			models.ValidationCandidate{Type: "xyz"},
			true,
		},
		{
			models.ValidationCandidate{Type: "abc"},
			true,
		},
		{
			models.ValidationCandidate{Type: "def"},
			false,
		},
	}

	for _, tCase := range testCases {
		actual := sharedValidator.ShouldValidate(tCase.candidate)

		if tCase.expected != actual {
			t.Fatalf("expected != actual\ncandidate type: %s\nexpected: %t\nactual: %t",
				tCase.candidate.Type, tCase.expected, actual)
		}
	}
}

func TestValidate(t *testing.T) {
	// getRef is used as a helper to return a reference to the given string
	getRef := func(input string) *string {
		return &input
	}

	testCases := []struct {
		candidate       models.ValidationCandidate
		expected        models.ValidationResult
		alterSchema     bool
		alternateSchema []byte
	}{
		// success
		{
			models.ValidationCandidate{
				Type: "xyz",
				ID:   "1",
				Data: map[string]interface{}{
					missionKey: "Apollo 11",
					crewKey:    []string{"Neil", "Buzz", "Mike"},
					rocketKey:  "Saturn V",
				},
			},
			models.ValidationResult{
				ValidatorSource: validatorSource,
				Result:          models.Success,
				CandidateID:     "1",
			},
			false,
			nil,
		},
		// failure - bad mission
		{
			models.ValidationCandidate{
				Type: "xyz",
				ID:   "2",
				Data: map[string]interface{}{
					missionKey: 11,
					crewKey:    []string{"Neil", "Buzz", "Mike"},
					rocketKey:  "Saturn V",
				},
			},
			models.ValidationResult{
				ValidatorSource: validatorSource,
				Result:          models.Failure,
				CandidateID:     "2",
				Error:           getRef("[mission: Invalid type. Expected: string, given: integer]"),
			},
			false,
			nil,
		},
		// failure - bad crew
		{
			models.ValidationCandidate{
				Type: "xyz",
				ID:   "3",
				Data: map[string]interface{}{
					missionKey: "Apollo 11",
					crewKey:    []int{1, 2, 3},
					rocketKey:  "Saturn V",
				},
			},
			models.ValidationResult{
				ValidatorSource: validatorSource,
				Result:          models.Failure,
				CandidateID:     "3",
				Error: getRef(
					"[crew.0: Invalid type. Expected: string, given: integer crew.1: Invalid type. Expected: string, " +
						"given: integer crew.2: Invalid type. Expected: string, given: integer]",
				),
			},
			false,
			nil,
		},
		// failure - bad rocket
		{
			models.ValidationCandidate{
				Type: "xyz",
				ID:   "4",
				Data: map[string]interface{}{
					missionKey: "Apollo 11",
					crewKey:    []string{"Neil", "Buzz", "Mike"},
					rocketKey:  "Jupiter IV",
				},
			},
			models.ValidationResult{
				ValidatorSource: validatorSource,
				Result:          models.Failure,
				CandidateID:     "4",
				Error: getRef(
					"[rocket: rocket must be one of the following: \"Titan II\", \"Saturn V\", \"Falcon Heavy\", " +
						"\"Falcon 9\"]",
				),
			},
			false,
			nil,
		},
		// failure - bad schema
		{
			models.ValidationCandidate{
				Type: "xyz",
				ID:   "5",
				Data: map[string]interface{}{
					missionKey: "Apollo 11",
					crewKey:    []string{"Neil", "Buzz", "Mike"},
					rocketKey:  "Saturn V",
				},
			},
			models.ValidationResult{
				ValidatorSource: validatorSource,
				Result:          models.Failure,
				CandidateID:     "5",
				Error: getRef(
					"has a primitive type that is NOT VALID -- given: /junk/ Expected valid values are:[array " +
						"boolean integer number null object string]",
				),
			},
			true,
			[]byte("{\"type\": \"junk\"}"),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// schema is saved in the case we want to override and then revert
	initialSchema := schemaBytes
	for _, tCase := range testCases {
		channel := make(chan *models.ValidationResult)
		if tCase.alterSchema {
			schemaBytes = tCase.alternateSchema
		} else {
			schemaBytes = initialSchema
		}
		go sharedValidator.Validate(tCase.candidate, channel)

		// we are waiting on each case so that we have access to relevant tCase values
		select {
		case <-ctx.Done():
			t.Fatalf("context timeout")
		case actual := <-channel:
			if !reflect.DeepEqual(tCase.expected, *actual) {
				// gracefully handle error string pointers
				expectedErrStr := "nil"
				if tCase.expected.Error != nil {
					expectedErrStr = *tCase.expected.Error
				}
				actualErrStr := "nil"
				if actual.Error != nil {
					actualErrStr = *actual.Error
				}

				failureStr := `
				expected != actual
				___
				expected:
					ValidatorSource: %s
					Result: %s
					CandidateId: %s
					Error: %s
				___
				actual:
					ValidatorSource: %s
					Result: %s
					CandidateId: %s
					Error: %s
				`
				t.Fatalf(
					failureStr,
					tCase.expected.ValidatorSource,
					tCase.expected.Result,
					tCase.expected.CandidateID,
					expectedErrStr,
					actual.ValidatorSource,
					actual.Result,
					actual.CandidateID,
					actualErrStr,
				)
			}
		}
	}
}
