package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	password := "supersecretpassword123!@#"
	expiresIn := 15 * time.Minute

	tests := map[string]struct {
		userID   uuid.UUID
		password string
		expIn    time.Duration
		want     error
	}{
		"create token": {
			userID:   userID,
			password: password,
			expIn:    expiresIn,
			want:     nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenString, err := MakeJWT(tc.userID, tc.password, tc.expIn)
			if err != tc.want {
				t.Fatalf("unexpected error\n\twant: %v, got: %v\n", tc.want, err)
			}
			// Using to check IssuedAt later.
			now := time.Now()

			if err == nil && tokenString == "" {
				t.Fatalf("expected non-empty token, but got empty string\n")
			}

			claims := &jwt.RegisteredClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return []byte(password), nil
			})
			if err != nil {
				t.Fatalf("Failed to parse JWT: %s\n", err)
			}
			if !token.Valid {
				t.Fatalf("Invalid token\n")
			}

			if claims.Issuer != "chirpy" {
				t.Fatalf("Expected issuer chirpy, got %s\n", claims.Issuer)
			}
			if claims.Subject != userID.String() {
				t.Fatalf("Expected subject %s, got %s\n", userID.String(), claims.Subject)
			}

			allowedTimeDrift := time.Second * 5
			if claims.IssuedAt == nil {
				t.Fatal("IssuedAt is nil")
			}
			if now.Sub(claims.IssuedAt.Time).Abs() > allowedTimeDrift {
				t.Fatalf(
					"IssuedAt time is too far from now: %v\n\tAllowed drift: %v, actual: %v",
					claims.IssuedAt,
					allowedTimeDrift,
					now.Sub(claims.IssuedAt.Time),
				)
			}
			if claims.ExpiresAt == nil {
				t.Fatal("ExpiresAt is nil")
			}
			if claims.ExpiresAt.Time.Sub(now.Add(expiresIn)).Abs() > allowedTimeDrift {
				t.Fatalf(
					"IssuedAt time is too far from now: %v\n\tAllowed drift: %v, actual: %v",
					claims.IssuedAt,
					allowedTimeDrift,
					now.Sub(claims.IssuedAt.Time),
				)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	goodPassword := "goodpassword"
	badPassword := "badpassword"
	normalExpiry := 15 * time.Minute
	shortExpiry := 1 * time.Second

	tests := map[string]struct {
		userID    uuid.UUID
		password  string
		expIn     time.Duration
		expUserID string
		expError  error
	}{
		"normal token": {
			userID:    userID,
			password:  goodPassword,
			expIn:     normalExpiry,
			expUserID: userID.String(),
			expError:  nil,
		},
		"expired token": {
			userID:    userID,
			password:  goodPassword,
			expIn:     shortExpiry,
			expUserID: uuid.Nil.String(),
			expError:  jwt.ErrTokenExpired,
		},
		"incorrect secret": {
			userID:    userID,
			password:  badPassword,
			expIn:     normalExpiry,
			expUserID: uuid.Nil.String(),
			expError:  jwt.ErrTokenSignatureInvalid,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenString, err := MakeJWT(tc.userID, tc.password, tc.expIn)
			if err != nil {
				t.Fatalf("unexpected error\n")
			}

			if name == "expired token" {
				time.Sleep(shortExpiry)
			}

			userID, err := ValidateJWT(tokenString, goodPassword)
			if !errors.Is(err, tc.expError) {
				t.Fatalf("Expected error %s, got %s\n", tc.expError, err)
			}
			if userID.String() != tc.expUserID {
				t.Fatalf("Expected user ID %s, got %s\n", tc.expUserID, userID.String())
			}
		})

	}
}
