package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "super-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash to not be empty")
	}

	if hash == password {
		t.Fatal("hash should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "super-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error checking password hash, got: %v", err)
	}

	if !match {
		t.Fatal("expected password hashes to match")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	password := "super-secret-password"
	wrong_password := "incorrect-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}

	match, err := CheckPasswordHash(wrong_password, hash)
	if err != nil {
		t.Fatalf("expected no error checking password hash, got: %v", err)
	}

	if match {
		t.Fatal("expected password hashes to not match")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer test-token"}

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected no error getting bearer token, got: %v", err)
	}

	if token != "test-token" {
		t.Fatalf("expected token to be 'test-token', got: %s", token)
	}
}

func TestGetBearerTokenMissingHeader(t *testing.T) {
	headers := make(map[string][]string)

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error when Authorization header is missing")
	}
}

func TestGetBearerTokenInvalidFormat(t *testing.T) {
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"InvalidFormat"}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected error when Authorization header format is invalid")
	}
}

func TestMakeJWT(t *testing.T) {
	user_id := uuid.New()
	token_secret := "test-secret"

	token, err := MakeJWT(
		user_id,
		token_secret,
		time.Hour,
	)

	if err != nil {
		t.Fatalf("expected no error creating JWT, got: %v", err)
	}

	if token == "" {
		t.Fatal("expected token to not be empty")
	}
}

func TestValidateJWT(t *testing.T) {
	user_id := uuid.New()
	token_secret := "test-secret"

	token, err := MakeJWT(
		user_id,
		token_secret,
		time.Hour,
	)

	if err != nil {
		t.Fatalf("expected no error creating JWT, got: %v", err)
	}

	validated_user_id, err := ValidateJWT(
		token,
		token_secret,
	)

	if err != nil {
		t.Fatalf("expected no error validating JWT, got: %v", err)
	}

	if validated_user_id != user_id {
		t.Fatalf(
			"expected user ID %v, got %v",
			user_id,
			validated_user_id,
		)
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	user_id := uuid.New()

	token, err := MakeJWT(
		user_id,
		"correct-secret",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("expected no error creating JWT, got: %v", err)
	}

	_, err = ValidateJWT(
		token,
		"wrong-secret",
	)

	if err == nil {
		t.Fatal("expected validation to fail with wrong secret")
	}
}

func TestValidateJWTExpired(t *testing.T) {
	user_id := uuid.New()
	token_secret := "test-secret"

	token, err := MakeJWT(
		user_id,
		token_secret,
		-time.Hour,
	)

	if err != nil {
		t.Fatalf("expected no error creating JWT, got: %v", err)
	}

	_, err = ValidateJWT(
		token,
		token_secret,
	)

	if err == nil {
		t.Fatal("expected validation to fail for expired token")
	}
}

func TestValidateJWTMalformedToken(t *testing.T) {
	token_secret := "test-secret"

	_, err := ValidateJWT(
		"this-is-not-a-valid-token",
		token_secret,
	)

	if err == nil {
		t.Fatal("expected validation to fail for malformed token")
	}
}
