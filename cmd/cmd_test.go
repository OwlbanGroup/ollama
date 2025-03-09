package cmd

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus" // Correctly importing logrus for logging
)

func setup() error {
	// Initialize any necessary state for the tests
	return nil
}



func TestMain(m *testing.M) {
	// Ensure setup is called only once
	if err := setup(); err != nil {
		logrus.Fatalf("Failed to setup: %v", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}
