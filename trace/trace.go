package trace

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/tracer"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type TraceQueryFinishFunc = tracer.QueryFinishFunc
type TraceFieldFinishFunc = tracer.FieldFinishFunc

type Tracer = tracer.Tracer

// Deprecated: <reason> ?
type OpenTracingTracer struct{}

func (OpenTracingTracer) TraceQuery(ctx context.Context, queryString string, operationName string, variables map[string]interface{}, varTypes map[string]*introspection.Type) (context.Context, TraceQueryFinishFunc) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "GraphQL request")
	span.SetTag("graphql.query", queryString)

	if operationName != "" {
		span.SetTag("graphql.operationName", operationName)
	}

	if len(variables) != 0 {
		span.LogFields(log.Object("graphql.variables", variables))
	}

	return spanCtx, func(errs []*errors.QueryError) {
		if len(errs) > 0 {
			msg := errs[0].Error()
			if len(errs) > 1 {
				msg += fmt.Sprintf(" (and %d more errors)", len(errs)-1)
			}
			ext.Error.Set(span, true)
			span.SetTag("graphql.error", msg)
		}
		span.Finish()
	}
}

func (OpenTracingTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]interface{}) (context.Context, TraceFieldFinishFunc) {
	if trivial {
		return ctx, noop
	}

	span, spanCtx := opentracing.StartSpanFromContext(ctx, label)
	span.SetTag("graphql.type", typeName)
	span.SetTag("graphql.field", fieldName)
	for name, value := range args {
		span.SetTag("graphql.args."+name, value)
	}

	return spanCtx, func(err *errors.QueryError) {
		if err != nil {
			ext.Error.Set(span, true)
			span.SetTag("graphql.error", err.Error())
		}
		span.Finish()
	}
}

func (OpenTracingTracer) TraceValidation(ctx context.Context) TraceValidationFinishFunc {
	span, _ := opentracing.StartSpanFromContext(ctx, "Validate Query")

	return func(errs []*errors.QueryError) {
		if len(errs) > 0 {
			msg := errs[0].Error()
			if len(errs) > 1 {
				msg += fmt.Sprintf(" (and %d more errors)", len(errs)-1)
			}
			ext.Error.Set(span, true)
			span.SetTag("graphql.error", msg)
		}
		span.Finish()
	}
}

func noop(*errors.QueryError) {}

type NoopTracer = tracer.Noop
