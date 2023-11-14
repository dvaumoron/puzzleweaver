// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package passwordstrengthimpl

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:  "github.com/dvaumoron/puzzleweaver/serviceimpl/passwordstrength/PasswordStrengthService",
		Iface: reflect.TypeOf((*PasswordStrengthService)(nil)).Elem(),
		Impl:  reflect.TypeOf(strengthImpl{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return passwordStrengthService_local_stub{impl: impl.(PasswordStrengthService), tracer: tracer, getRulesMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/passwordstrength/PasswordStrengthService", Method: "GetRules", Remote: false}), validateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/passwordstrength/PasswordStrengthService", Method: "Validate", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return passwordStrengthService_client_stub{stub: stub, getRulesMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/passwordstrength/PasswordStrengthService", Method: "GetRules", Remote: true}), validateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/passwordstrength/PasswordStrengthService", Method: "Validate", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return passwordStrengthService_server_stub{impl: impl.(PasswordStrengthService), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return passwordStrengthService_reflect_stub{caller: caller}
		},
		RefData: "",
	})
}

// weaver.InstanceOf checks.
var _ weaver.InstanceOf[PasswordStrengthService] = (*strengthImpl)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*strengthImpl)(nil)

// Local stub implementations.

type passwordStrengthService_local_stub struct {
	impl            PasswordStrengthService
	tracer          trace.Tracer
	getRulesMetrics *codegen.MethodMetrics
	validateMetrics *codegen.MethodMetrics
}

// Check that passwordStrengthService_local_stub implements the PasswordStrengthService interface.
var _ PasswordStrengthService = (*passwordStrengthService_local_stub)(nil)

func (s passwordStrengthService_local_stub) GetRules(ctx context.Context, a0 string) (r0 string, err error) {
	// Update metrics.
	begin := s.getRulesMetrics.Begin()
	defer func() { s.getRulesMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "passwordstrengthimpl.PasswordStrengthService.GetRules", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetRules(ctx, a0)
}

func (s passwordStrengthService_local_stub) Validate(ctx context.Context, a0 string) (err error) {
	// Update metrics.
	begin := s.validateMetrics.Begin()
	defer func() { s.validateMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "passwordstrengthimpl.PasswordStrengthService.Validate", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Validate(ctx, a0)
}

// Client stub implementations.

type passwordStrengthService_client_stub struct {
	stub            codegen.Stub
	getRulesMetrics *codegen.MethodMetrics
	validateMetrics *codegen.MethodMetrics
}

// Check that passwordStrengthService_client_stub implements the PasswordStrengthService interface.
var _ PasswordStrengthService = (*passwordStrengthService_client_stub)(nil)

func (s passwordStrengthService_client_stub) GetRules(ctx context.Context, a0 string) (r0 string, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getRulesMetrics.Begin()
	defer func() { s.getRulesMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "passwordstrengthimpl.PasswordStrengthService.GetRules", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()

	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += (4 + len(a0))
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.String(a0)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = dec.String()
	err = dec.Error()
	return
}

func (s passwordStrengthService_client_stub) Validate(ctx context.Context, a0 string) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.validateMetrics.Begin()
	defer func() { s.validateMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "passwordstrengthimpl.PasswordStrengthService.Validate", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()

	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += (4 + len(a0))
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.String(a0)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 1, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

// Note that "weaver generate" will always generate the error message below.
// Everything is okay. The error message is only relevant if you see it when
// you run "go build" or "go run".
var _ codegen.LatestVersion = codegen.Version[[0][20]struct{}](`

ERROR: You generated this file with 'weaver generate' v0.22.0 (codegen
version v0.20.0). The generated code is incompatible with the version of the
github.com/ServiceWeaver/weaver module that you're using. The weaver module
version can be found in your go.mod file or by running the following command.

    go list -m github.com/ServiceWeaver/weaver

We recommend updating the weaver module and the 'weaver generate' command by
running the following.

    go get github.com/ServiceWeaver/weaver@latest
    go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

Then, re-run 'weaver generate' and re-build your code. If the problem persists,
please file an issue at https://github.com/ServiceWeaver/weaver/issues.

`)

// Server stub implementations.

type passwordStrengthService_server_stub struct {
	impl    PasswordStrengthService
	addLoad func(key uint64, load float64)
}

// Check that passwordStrengthService_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*passwordStrengthService_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s passwordStrengthService_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "GetRules":
		return s.getRules
	case "Validate":
		return s.validate
	default:
		return nil
	}
}

func (s passwordStrengthService_server_stub) getRules(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 string
	a0 = dec.String()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.GetRules(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.String(r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s passwordStrengthService_server_stub) validate(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 string
	a0 = dec.String()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.Validate(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

// Reflect stub implementations.

type passwordStrengthService_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that passwordStrengthService_reflect_stub implements the PasswordStrengthService interface.
var _ PasswordStrengthService = (*passwordStrengthService_reflect_stub)(nil)

func (s passwordStrengthService_reflect_stub) GetRules(ctx context.Context, a0 string) (r0 string, err error) {
	err = s.caller("GetRules", ctx, []any{a0}, []any{&r0})
	return
}

func (s passwordStrengthService_reflect_stub) Validate(ctx context.Context, a0 string) (err error) {
	err = s.caller("Validate", ctx, []any{a0}, []any{})
	return
}
