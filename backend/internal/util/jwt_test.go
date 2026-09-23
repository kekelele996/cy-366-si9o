package util

import "testing"

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test_secret_for_jwt_unit_test"
	token, err := GenerateToken(secret, 3600, 1, "admin", "admin")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if claims.UserID != 1 || claims.Username != "admin" || claims.Role != "admin" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
	if _, err := ParseToken(secret, token+"bad"); err == nil {
		t.Fatal("expected error for tampered token")
	}
}
