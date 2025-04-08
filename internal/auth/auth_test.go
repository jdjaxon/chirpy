package auth

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	mockPassword := "supersecretpassword123!@#"
	mockLongPassword := mockPassword + mockPassword + mockPassword

	tests := map[string]struct {
		input string
		want  error
	}{
		"normal password": {
			input: mockPassword,
			want:  nil,
		},
		"empty password": {
			input: "",
			want:  nil,
		},
		"long password": {
			input: mockLongPassword,
			want:  bcrypt.ErrPasswordTooLong,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hash, err := HashPassword(tc.input)

			// if (err != nil && tc.want == nil) || (err == nil && tc.want != nil) {
			if !errors.Is(err, tc.want) {
				t.Fatalf("unexpected error\n\twant: %v, got: %v\n", tc.want, err)
			}

			if err == nil && hash == "" {
				t.Fatalf("expected non-empty hash, but got empty string\n")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	mockPassword := "mysupersecretpassword"
	mockHashedPassword, _ := HashPassword(mockPassword)

	tests := map[string]struct {
		inputPassword string
		inputHash     string
		want          error
	}{
		"correct password": {
			inputPassword: mockPassword,
			inputHash:     mockHashedPassword,
			want:          nil,
		},
		"wrong password": {
			inputPassword: "wrong",
			inputHash:     mockHashedPassword,
			want:          bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := CheckPasswordHash(tc.inputHash, tc.inputPassword)
			if (err != nil && tc.want == nil) || (err == nil && tc.want != nil) {
				t.Fatalf("unexpected error\n\twant: %v, got: %v\n", tc.want, err)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	headerWithAuth := http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer test-token"},
	}
	authHeader := headerWithAuth.Get("Authorization")

	malformedHeader := http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"SomeToken heresATokenForYou"},
	}

	headerWithoutAuth := http.Header{
		"Content-Type": {"application/json"},
	}

	emptyHeader := http.Header{}

	tests := map[string]struct {
		input         http.Header
		want          string
		expectedError error
	}{
		"auth header": {
			input:         headerWithAuth,
			want:          strings.Split(authHeader, " ")[1],
			expectedError: nil,
		},
		"no auth header": {
			input:         headerWithoutAuth,
			want:          "",
			expectedError: ErrNoAuthHeader,
		},
		"malformed header": {
			input:         malformedHeader,
			want:          "",
			expectedError: ErrMalformedAuthHeader,
		},
		"empty header": {
			input:         emptyHeader,
			want:          "",
			expectedError: ErrNoAuthHeader,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := GetBearerToken(tc.input)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("Name: %s\n\tExpected error: %#v, got: %#v\n", name, tc.want, output)
			}
			if !reflect.DeepEqual(tc.want, output) {
				t.Fatalf("Name: %s\n\tExpected output: %#v, got: %#v\n", name, tc.want, output)
			}
		})
	}
}
