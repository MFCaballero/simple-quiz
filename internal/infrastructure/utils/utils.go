package utils

import (
	"context"
	"errors"
	"io"
)

type contextKey string

const (
	readerKey contextKey = "reader"
)

// WithReader sets the io.Reader in the context.
func WithReader(ctx context.Context, reader io.Reader) context.Context {
	return context.WithValue(ctx, readerKey, reader)
}

// ReaderFromContext retrieves the io.Reader from the context.
func ReaderFromContext(ctx context.Context) (io.Reader, error) {
	reader, ok := ctx.Value(readerKey).(io.Reader)
	if !ok {
		return nil, errors.New("io.Reader not found in context")
	}
	return reader, nil
}
