/*
Copyright 2017 gin-tokenauth authors.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
