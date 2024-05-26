package testutil

import (
	"github.com/BrunoQuaresma/openticket/sdk"
)

func HasValidationError(res sdk.RequestResult[any], field string, validator string) bool {
	if len(res.Error.Errors) == 0 {
		return false
	}

	for _, e := range res.Error.Errors {
		if e.Field == field && e.Validator == validator {
			return true
		}
	}
	return false
}
