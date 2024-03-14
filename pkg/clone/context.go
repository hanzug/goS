package clone

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type ContextWithoutDeadline struct {
	ctx context.Context
}

func (c *ContextWithoutDeadline) Value(key any) any {
	return key
}

func (*ContextWithoutDeadline) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*ContextWithoutDeadline) Done() <-chan struct{}       { return nil }
func (*ContextWithoutDeadline) Err() error                  { return nil }

func NewContextWithoutDeadline() *ContextWithoutDeadline {
	return &ContextWithoutDeadline{ctx: context.Background()}
}

func (c *ContextWithoutDeadline) Clone(ctx context.Context, keys ...interface{}) {
	_, span := otel.GetTracerProvider().Tracer("goS"+"/clone").
		Start(ctx, "clone", oteltrace.WithAttributes(attribute.Int("sync", 1)))
	defer span.End()

	c.ctx = oteltrace.ContextWithSpan(c.ctx, span)

	for _, key := range keys {
		if v := ctx.Value(key); v != nil {
			c.ctx = context.WithValue(c.ctx, key, v)
		}
	}
}
