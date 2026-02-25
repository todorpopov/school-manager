package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todorpopov/school-manager/internal"
)

func TestBCryptService_PasswordsMatch_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		testPassword string
		shouldMatch  bool
	}{
		{
			name:         "exact match",
			password:     "password123",
			testPassword: "password123",
			shouldMatch:  true,
		},
		{
			name:         "case sensitive mismatch",
			password:     "Password123",
			testPassword: "password123",
			shouldMatch:  false,
		},
		{
			name:         "completely different password",
			password:     "correctPassword",
			testPassword: "wrongPassword",
			shouldMatch:  false,
		},
		{
			name:         "special characters match",
			password:     "p@ssw0rd!#$%",
			testPassword: "p@ssw0rd!#$%",
			shouldMatch:  true,
		},
		{
			name:         "long password match",
			password:     "thisIsAVeryLongPasswordThatShouldStillWork123!@#",
			testPassword: "thisIsAVeryLongPasswordThatShouldStillWork123!@#",
			shouldMatch:  true,
		},
		{
			name:         "empty string no match",
			password:     "password",
			testPassword: "",
			shouldMatch:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := internal.NewBCryptService()
			hash, success := bs.HashPassword(tt.password)
			assert.True(t, success)
			matches := bs.PasswordsMatch(hash, tt.testPassword)
			assert.Equal(t, tt.shouldMatch, matches)
		})
	}
}
