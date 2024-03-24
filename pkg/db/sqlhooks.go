package db

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

type hook struct{}

// Before implements sqlhooks.Hooks.
func (h *hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		return ctx, nil
	}
	span := opentracing.StartSpan("sql",
		opentracing.ChildOf(parent.Context()),
		ext.SpanKindRPCClient)
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, "sql")
	span.LogFields(
		otlog.Object("args", args),
	)

	return opentracing.ContextWithSpan(ctx, span), nil
}

// After implements sqlhooks.Hooks.
func (h *hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.Finish()
	}
	return ctx, nil
}

// After implements sqlhooks.OnErroer.
func (h *hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		ext.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		span.Finish()
	}
	return err
}
