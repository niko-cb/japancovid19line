package utils

import (
	"context"
	"net/http"
)

// create new context for application
func NewContext(r *http.Request) context.Context {
	ctx := r.Context()
	return ctx
}
