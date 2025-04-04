package auth

import (
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

			if (err != nil && tc.want == nil) || (err == nil && tc.want != nil) {
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
