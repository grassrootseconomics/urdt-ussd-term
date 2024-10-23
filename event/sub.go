package event

import (
	"context"
	"io"
)

type Subscription interface {
	io.Closer
	Connect(ctx context.Context, connStr string) error
	Next() error
}
