package testutil

import "github.com/BrunoQuaresma/openticket/api"

func HasValidationError(errors []api.ValidationError, field string, validator string) bool {
	for _, e := range errors {
		if e.Field == field && e.Validator == validator {
			return true
		}
	}
	return false
}
