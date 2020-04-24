package context

import (
	"context"
	"net/http"
)

// create new context for application
func New(r *http.Request) context.Context {
	return r.Context()
}

func Get() context.Context {
	return context.Background()
}
