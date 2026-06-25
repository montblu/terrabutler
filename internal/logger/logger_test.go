package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test if Logger config is Valid
func TestLoggerConfig(t *testing.T) {

	_, err := configZap.Build()
	assert.NoError(t, err)

}
