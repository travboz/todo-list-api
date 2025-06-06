package token_test

import (
	"testing"

	"github.com/travboz/backend-projects/todo-list-api/internal/assert"
	"github.com/travboz/backend-projects/todo-list-api/internal/token"
)

func TestToken(t *testing.T) {
	t.Run("correct length", func(t *testing.T) {
		token_length := 16

		token, err := token.GenerateToken(token_length)
		if err != nil {
			t.Fatalf("unexpected errors: %v", err)
		}

		want := token_length * 2
		got := len(token)

		assert.Equal(t, got, want)
	})

	t.Run("creates different tokens", func(t *testing.T) {
		token_length := 16

		token1, err := token.GenerateToken(token_length)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		token2, err := token.GenerateToken(token_length)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := false
		got := token1 == token2

		assert.Equal(t, got, want)
	})

	t.Run("empty token for length 0", func(t *testing.T) {
		token_length := 0

		token, err := token.GenerateToken(token_length)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := 0
		got := len(token)

		assert.Equal(t, got, want)
	})
}
