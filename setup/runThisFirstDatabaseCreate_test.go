package main

import (
	"testing"
)

func TestReadData(t *testing.T) {
	records := readData()
	if len(records) == 0 {
		t.Errorf("Expected records to be read, got none.")
	}
}
