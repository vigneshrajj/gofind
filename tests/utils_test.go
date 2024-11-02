package tests

import (
	"testing"

	"github.com/vigneshrajj/gofind/internal/helpers"
)

func TestB64EncodeHelper(t *testing.T) {
	data := "test"
	encoded := helpers.GetB64(data)
	if encoded != "dGVzdA==" {
		t.Fatalf("Expected dGVzdA==, but got %v", encoded)
	}
}

func TestB64DecodeHelper(t *testing.T) {
	data := "dGVzdA=="
	decoded := helpers.GetB64Decode(data)
	if decoded != "test" {
		t.Fatalf("Expected test, but got %v", decoded)
	}
}

func TestSha256Helper(t *testing.T) {
	data := "test"
	hashed := helpers.Sha256(data)
	if hashed != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Fatalf("Expected 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08, but got %v", hashed)
	}
}
