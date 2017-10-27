package tokenauth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenAuth is token auth.
type TokenAuth struct {
	store TokenStore
}

// TokenStore is a token store that stores and validates tokens.
type TokenStore interface {
	IsTokenValid(t string) bool
}

// Token represents a single auth token.
type Token struct {
	Name        string `yaml:"name"`
	Token       string `yaml:"token"`
	Description string `yaml:"description"`
	IsDisabled  bool   `yaml:"is_disabled"`
}

// Authenticate validates given token t and returns a bool to indicate whether
// the token is valid or not.
func (a *TokenAuth) Authenticate(t string) bool {
	return a.store.IsTokenValid(t)
}

// New initializes a new middleware that handles gin requests and validates
// Authorization header Bearer token against the store.
func New(store TokenStore) gin.HandlerFunc {
	auth := &TokenAuth{
		store: store,
	}

	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")

		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !auth.Authenticate(strings.TrimPrefix(h, "Bearer ")) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
