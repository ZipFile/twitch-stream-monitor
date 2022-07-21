package cli

import (
	"testing"
)

func TestEmergencyLoggerFactory(t *testing.T) {
	log := emergencyLoggerFactory()

	if log == nil {
		t.Error("log: nil; expected: not nil")
	}
}
