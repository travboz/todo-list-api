package env_test

import (
	"os"
	"testing"

	"github.com/travboz/backend-projects/todo-list-api/internal/assert"
	"github.com/travboz/backend-projects/todo-list-api/internal/env"
)

func TestGetString(t *testing.T) {
	t.Run("environment variable is set", func(t *testing.T) {
		key := "SOME_SET_KEY"

		want := "value"
		os.Setenv(key, want)
		defer os.Unsetenv(key) // clean up env

		got := env.GetString(key, "fallback")

		assert.Equal(t, got, want)

	})

	t.Run("environment variable is NOT set", func(t *testing.T) {
		key := "SOME_SET_KEY"

		want := "fallback"
		got := env.GetString(key, "fallback")

		assert.Equal(t, got, want)
	})
}
