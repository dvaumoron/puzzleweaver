// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package sessionimpl

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
		Name:  "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService",
		Iface: reflect.TypeOf((*SessionService)(nil)).Elem(),
		Impl:  reflect.TypeOf(sessionImpl{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return sessionService_local_stub{impl: impl.(SessionService), tracer: tracer, generateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService", Method: "Generate", Remote: false}), getMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService", Method: "Get", Remote: false}), updateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService", Method: "Update", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return sessionService_client_stub{stub: stub, generateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService", Method: "Generate", Remote: true}), getMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService", Method: "Get", Remote: true}), updateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/session/SessionService", Method: "Update", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return sessionService_server_stub{impl: impl.(SessionService), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return sessionService_reflect_stub{caller: caller}
		},
		RefData: "",
	})
}

// weaver.InstanceOf checks.
var _ weaver.InstanceOf[SessionService] = (*sessionImpl)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*sessionImpl)(nil)

// Local stub implementations.

type sessionService_local_stub struct {
	impl            SessionService
	tracer          trace.Tracer
	generateMetrics *codegen.MethodMetrics
	getMetrics      *codegen.MethodMetrics
	updateMetrics   *codegen.MethodMetrics
}

// Check that sessionService_local_stub implements the SessionService interface.
var _ SessionService = (*sessionService_local_stub)(nil)

func (s sessionService_local_stub) Generate(ctx context.Context) (r0 uint64, err error) {
	// Update metrics.
	begin := s.generateMetrics.Begin()
	defer func() { s.generateMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "sessionimpl.SessionService.Generate", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Generate(ctx)
}

func (s sessionService_local_stub) Get(ctx context.Context, a0 uint64) (r0 map[string]string, err error) {
	// Update metrics.
	begin := s.getMetrics.Begin()
	defer func() { s.getMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "sessionimpl.SessionService.Get", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Get(ctx, a0)
}

func (s sessionService_local_stub) Update(ctx context.Context, a0 uint64, a1 map[string]string) (err error) {
	// Update metrics.
	begin := s.updateMetrics.Begin()
	defer func() { s.updateMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "sessionimpl.SessionService.Update", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Update(ctx, a0, a1)
}

// Client stub implementations.

type sessionService_client_stub struct {
	stub            codegen.Stub
	generateMetrics *codegen.MethodMetrics
	getMetrics      *codegen.MethodMetrics
	updateMetrics   *codegen.MethodMetrics
}

// Check that sessionService_client_stub implements the SessionService interface.
var _ SessionService = (*sessionService_client_stub)(nil)

func (s sessionService_client_stub) Generate(ctx context.Context) (r0 uint64, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.generateMetrics.Begin()
	defer func() { s.generateMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "sessionimpl.SessionService.Generate", trace.WithSpanKind(trace.SpanKindClient))
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

	var shardKey uint64

	// Call the remote method.
	var results []byte
	results, err = s.stub.Run(ctx, 0, nil, shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = dec.Uint64()
	err = dec.Error()
	return
}

func (s sessionService_client_stub) Get(ctx context.Context, a0 uint64) (r0 map[string]string, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getMetrics.Begin()
	defer func() { s.getMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "sessionimpl.SessionService.Get", trace.WithSpanKind(trace.SpanKindClient))
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
	size += 8
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Uint64(a0)
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
	r0 = serviceweaver_dec_map_string_string_219dd46d(dec)
	err = dec.Error()
	return
}

func (s sessionService_client_stub) Update(ctx context.Context, a0 uint64, a1 map[string]string) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.updateMetrics.Begin()
	defer func() { s.updateMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "sessionimpl.SessionService.Update", trace.WithSpanKind(trace.SpanKindClient))
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

	// Encode arguments.
	enc := codegen.NewEncoder()
	enc.Uint64(a0)
	serviceweaver_enc_map_string_string_219dd46d(enc, a1)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 2, enc.Data(), shardKey)
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

ERROR: You generated this file with 'weaver generate' v0.23.0 (codegen
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

type sessionService_server_stub struct {
	impl    SessionService
	addLoad func(key uint64, load float64)
}

// Check that sessionService_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*sessionService_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s sessionService_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "Generate":
		return s.generate
	case "Get":
		return s.get
	case "Update":
		return s.update
	default:
		return nil
	}
}

func (s sessionService_server_stub) generate(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.Generate(ctx)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Uint64(r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s sessionService_server_stub) get(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 uint64
	a0 = dec.Uint64()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.Get(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_map_string_string_219dd46d(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s sessionService_server_stub) update(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 uint64
	a0 = dec.Uint64()
	var a1 map[string]string
	a1 = serviceweaver_dec_map_string_string_219dd46d(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.Update(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

// Reflect stub implementations.

type sessionService_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that sessionService_reflect_stub implements the SessionService interface.
var _ SessionService = (*sessionService_reflect_stub)(nil)

func (s sessionService_reflect_stub) Generate(ctx context.Context) (r0 uint64, err error) {
	err = s.caller("Generate", ctx, []any{}, []any{&r0})
	return
}

func (s sessionService_reflect_stub) Get(ctx context.Context, a0 uint64) (r0 map[string]string, err error) {
	err = s.caller("Get", ctx, []any{a0}, []any{&r0})
	return
}

func (s sessionService_reflect_stub) Update(ctx context.Context, a0 uint64, a1 map[string]string) (err error) {
	err = s.caller("Update", ctx, []any{a0, a1}, []any{})
	return
}

// Encoding/decoding implementations.

func serviceweaver_enc_map_string_string_219dd46d(enc *codegen.Encoder, arg map[string]string) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for k, v := range arg {
		enc.String(k)
		enc.String(v)
	}
}

func serviceweaver_dec_map_string_string_219dd46d(dec *codegen.Decoder) map[string]string {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make(map[string]string, n)
	var k string
	var v string
	for i := 0; i < n; i++ {
		k = dec.String()
		v = dec.String()
		res[k] = v
	}
	return res
}
