package utils

import (
	"context"
	"net/http"
)

// create new context for application
func NewContext(r *http.Request) context.Context {
	return r.Context()
}

func GetContext() context.Context {
	return context.Background()
}
