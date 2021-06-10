package services

import (
	"testing"
)

func TestHashSample(t *testing.T) {
	password := "random!pass"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error("Error in hashing password")
	}

	if CheckPasswordHash(password, hash) != true {
		t.Error("Failed to match password with hash")
	}

	hashWrong, err := HashPassword("random!|pass")
	if CheckPasswordHash(password, hashWrong) == true {
		t.Error("Matched wrong hash with password")
	}
}
