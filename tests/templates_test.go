package tests

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vigneshrajj/gofind/internal/templates"
)

func TestBase64Template(t *testing.T) {
	defer setupServerTest()()
	data := "test"
	w := httptest.NewRecorder()
	templates.Base64Template(w, data)
	bodyBytes, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	body := string(bodyBytes)
	if !strings.Contains(body, "test") {
		t.Fatalf("Expected test, but got %v", body)
	}
}

func TestBase64DecodeTemplate(t *testing.T) {
	defer setupServerTest()()
	data := "test"
	w := httptest.NewRecorder()
	templates.Base64DecodeTemplate(w, data)
	bodyBytes, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	body := string(bodyBytes)
	if !strings.Contains(body, "test") {
		t.Fatalf("Expected test, but got %v", body)
	}
}

func TestSha256Template(t *testing.T) {
	defer setupServerTest()()
	data := "test"
	w := httptest.NewRecorder()
	templates.Sha256Template(w, data)
	bodyBytes, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	body := string(bodyBytes)
	if !strings.Contains(body, "test") {
		t.Fatalf("Expected test, but got %v", body)
	}
}

func TestMessageTemplate(t *testing.T) {
	defer setupServerTest()()
	message := "test"
	w := httptest.NewRecorder()
	templates.MessageTemplate(w, message)
	bodyBytes, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	body := string(bodyBytes)
	if !strings.Contains(body, "test") {
		t.Fatalf("Expected test, but got %v", body)
	}
}
