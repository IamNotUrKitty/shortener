package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	env := getEnv("SERVER_PORT", "8080")

	assert.Equal(t, env, "8080")
}
