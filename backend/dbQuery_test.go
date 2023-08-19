package backend

import (
	"testing"
)

func TestTermExistsEmbed(t *testing.T) {
	// Test with a known term from the database
	term := "func"
	answer, exists, err := TermExistsEmbed(term)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !exists {
		t.Errorf("Expected term %s to exist", term)
	}
	if answer == "" {
		t.Errorf("Expected an answer for term %s", term)
	}

	// Test with a term that doesn't exist in the database
	term = "nonExistentTerm"
	answer, exists, err = TermExistsEmbed(term)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if exists {
		t.Errorf("Did not expect term %s to exist", term)
	}
	if answer != "" {
		t.Errorf("Did not expect an answer for term %s", term)
	}
}
