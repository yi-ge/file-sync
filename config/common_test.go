package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	// Ensure that calling Instance() returns the same instance each time.
	instance := Instance()

	assert.Equal(t, instance.name, "file-sync")
}
