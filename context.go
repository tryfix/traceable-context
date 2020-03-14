package traceable_context

import (
	"context"
	"github.com/google/uuid"
	"time"
)

var uuidPrefix = `uuid`

// TraceableContext wrap the go context with a traceable uuid
type TraceableContext interface {
	context.Context
	// UUID returns the uuid inside the context
	UUID() uuid.UUID
}

type traceableContext struct {
	context.Context
	uuid uuid.UUID
}

func WithCancel(parent context.Context) (ctx TraceableContext, cancel context.CancelFunc) {
	c, cancel := context.WithCancel(parent)
	return &traceableContext{
		Context: c,
	}, cancel
}

func WithDeadline(parent context.Context, deadline time.Time) (ctx TraceableContext, cancel context.CancelFunc) {
	c, cancel := context.WithDeadline(parent, deadline)
	return &traceableContext{
		Context: c,
	}, cancel
}

func WithTimeout(parent context.Context, timeout time.Duration) (ctx TraceableContext, cancel context.CancelFunc) {
	c, cancel := context.WithTimeout(parent, timeout)
	return &traceableContext{
		Context: c,
	}, cancel
}

func WithValue(parent context.Context, key, val interface{}) TraceableContext {
	return &traceableContext{
		Context: context.WithValue(parent, key, val),
	}
}

// WithUUID creates a new traceable context from a given UUID
func WithUUID(uuid uuid.UUID) TraceableContext {
	return &traceableContext{
		Context: context.WithValue(context.Background(), &uuidPrefix, uuid),
		uuid:    uuid,
	}
}

// FromContextWithUUID creates a new traceable context from a given parent context and a UUID
func FromContextWithUUID(parent context.Context, uuid uuid.UUID) TraceableContext {
	return &traceableContext{
		Context: context.WithValue(parent, &uuidPrefix, uuid),
		uuid:    uuid,
	}
}

func Background() context.Context {
	return &traceableContext{
		Context: context.Background(),
	}
}

// FromContext extracts the UUID from a given context
func FromContext(ctx context.Context) uuid.UUID {
	uid, ok := ctx.Value(&uuidPrefix).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return uid
}

func (c *traceableContext) Deadline() (deadline time.Time, ok bool) {
	return c.Context.Deadline()
}

func (c *traceableContext) Done() <-chan struct{} {
	return c.Context.Done()
}

func (c *traceableContext) Err() error {
	return c.Context.Err()
}

func (c *traceableContext) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}

func (c *traceableContext) UUID() uuid.UUID {
	u, ok := c.Value(&uuidPrefix).(uuid.UUID)
	if !ok {
		panic(`traceableContext.uuid dose not exist`)
	}

	return u
}
