package externalfunc

import (
	"testing"
)

func TestGetRandomCharacter(t *testing.T) {
	character, err := GetRandomCharacter()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if character == nil {
		t.Fatal("Expected character, got nil")
	}
	if character.ID == 0 {
		t.Error("Expected character ID to be non-zero")
	}
	if character.Name == "" {
		t.Error("Expected character Name to be non-empty")
	}
	if character.Image == "" {
		t.Error("Expected character Image to be non-empty")
	}
}
