package testutil

import "github.com/brianvoe/gofakeit"

func FakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 15)
}
