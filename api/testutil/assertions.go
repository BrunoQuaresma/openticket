package testutil

import (
	"fmt"
	"testing"

	"github.com/BrunoQuaresma/openticket/api"
	"github.com/stretchr/testify/require"
)

func RequireValidationError(t *testing.T, errors []api.ValidationError, field string, validator string) {
	hasError := false

	if len(errors) == 0 {
		hasError = false
	}

	for _, e := range errors {
		if e.Field == field && e.Validator == validator {
			hasError = true
			break
		}
	}

	require.True(t, hasError, fmt.Sprintf("expected error for field %s with validator %s", field, validator))
}
