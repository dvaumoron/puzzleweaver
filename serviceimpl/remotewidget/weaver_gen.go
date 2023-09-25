// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package remotewidgetimpl

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/service"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:  "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/RemoteWidgetService",
		Iface: reflect.TypeOf((*RemoteWidgetService)(nil)).Elem(),
		Impl:  reflect.TypeOf(remoteWidgetImpl{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return remoteWidgetService_local_stub{impl: impl.(RemoteWidgetService), tracer: tracer, getDescMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/RemoteWidgetService", Method: "GetDesc", Remote: false}), processMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/RemoteWidgetService", Method: "Process", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return remoteWidgetService_client_stub{stub: stub, getDescMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/RemoteWidgetService", Method: "GetDesc", Remote: true}), processMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/remotewidget/RemoteWidgetService", Method: "Process", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return remoteWidgetService_server_stub{impl: impl.(RemoteWidgetService), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return remoteWidgetService_reflect_stub{caller: caller}
		},
		RefData: "",
	})
}

// weaver.InstanceOf checks.
var _ weaver.InstanceOf[RemoteWidgetService] = (*remoteWidgetImpl)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*remoteWidgetImpl)(nil)

// Local stub implementations.

type remoteWidgetService_local_stub struct {
	impl           RemoteWidgetService
	tracer         trace.Tracer
	getDescMetrics *codegen.MethodMetrics
	processMetrics *codegen.MethodMetrics
}

// Check that remoteWidgetService_local_stub implements the RemoteWidgetService interface.
var _ RemoteWidgetService = (*remoteWidgetService_local_stub)(nil)

func (s remoteWidgetService_local_stub) GetDesc(ctx context.Context, a0 string) (r0 []remotewidgetservice.RawWidgetAction, err error) {
	// Update metrics.
	begin := s.getDescMetrics.Begin()
	defer func() { s.getDescMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "remotewidgetimpl.RemoteWidgetService.GetDesc", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetDesc(ctx, a0)
}

func (s remoteWidgetService_local_stub) Process(ctx context.Context, a0 string, a1 string, a2 map[string][]byte) (r0 string, r1 string, r2 []byte, err error) {
	// Update metrics.
	begin := s.processMetrics.Begin()
	defer func() { s.processMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "remotewidgetimpl.RemoteWidgetService.Process", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Process(ctx, a0, a1, a2)
}

// Client stub implementations.

type remoteWidgetService_client_stub struct {
	stub           codegen.Stub
	getDescMetrics *codegen.MethodMetrics
	processMetrics *codegen.MethodMetrics
}

// Check that remoteWidgetService_client_stub implements the RemoteWidgetService interface.
var _ RemoteWidgetService = (*remoteWidgetService_client_stub)(nil)

func (s remoteWidgetService_client_stub) GetDesc(ctx context.Context, a0 string) (r0 []remotewidgetservice.RawWidgetAction, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getDescMetrics.Begin()
	defer func() { s.getDescMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "remotewidgetimpl.RemoteWidgetService.GetDesc", trace.WithSpanKind(trace.SpanKindClient))
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
	r0 = serviceweaver_dec_slice_RawWidgetAction_845186e1(dec)
	err = dec.Error()
	return
}

func (s remoteWidgetService_client_stub) Process(ctx context.Context, a0 string, a1 string, a2 map[string][]byte) (r0 string, r1 string, r2 []byte, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.processMetrics.Begin()
	defer func() { s.processMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "remotewidgetimpl.RemoteWidgetService.Process", trace.WithSpanKind(trace.SpanKindClient))
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
	enc.String(a0)
	enc.String(a1)
	serviceweaver_enc_map_string_slice_byte_7ebbaefa(enc, a2)
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
	r0 = dec.String()
	r1 = dec.String()
	r2 = serviceweaver_dec_slice_byte_87461245(dec)
	err = dec.Error()
	return
}

// Note that "weaver generate" will always generate the error message below.
// Everything is okay. The error message is only relevant if you see it when
// you run "go build" or "go run".
var _ codegen.LatestVersion = codegen.Version[[0][20]struct{}](`

ERROR: You generated this file with 'weaver generate' v0.21.2 (codegen
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

type remoteWidgetService_server_stub struct {
	impl    RemoteWidgetService
	addLoad func(key uint64, load float64)
}

// Check that remoteWidgetService_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*remoteWidgetService_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s remoteWidgetService_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "GetDesc":
		return s.getDesc
	case "Process":
		return s.process
	default:
		return nil
	}
}

func (s remoteWidgetService_server_stub) getDesc(ctx context.Context, args []byte) (res []byte, err error) {
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
	r0, appErr := s.impl.GetDesc(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_slice_RawWidgetAction_845186e1(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s remoteWidgetService_server_stub) process(ctx context.Context, args []byte) (res []byte, err error) {
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
	var a1 string
	a1 = dec.String()
	var a2 map[string][]byte
	a2 = serviceweaver_dec_map_string_slice_byte_7ebbaefa(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, r1, r2, appErr := s.impl.Process(ctx, a0, a1, a2)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.String(r0)
	enc.String(r1)
	serviceweaver_enc_slice_byte_87461245(enc, r2)
	enc.Error(appErr)
	return enc.Data(), nil
}

// Reflect stub implementations.

type remoteWidgetService_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that remoteWidgetService_reflect_stub implements the RemoteWidgetService interface.
var _ RemoteWidgetService = (*remoteWidgetService_reflect_stub)(nil)

func (s remoteWidgetService_reflect_stub) GetDesc(ctx context.Context, a0 string) (r0 []remotewidgetservice.RawWidgetAction, err error) {
	err = s.caller("GetDesc", ctx, []any{a0}, []any{&r0})
	return
}

func (s remoteWidgetService_reflect_stub) Process(ctx context.Context, a0 string, a1 string, a2 map[string][]byte) (r0 string, r1 string, r2 []byte, err error) {
	err = s.caller("Process", ctx, []any{a0, a1, a2}, []any{&r0, &r1, &r2})
	return
}

// Encoding/decoding implementations.

func serviceweaver_enc_slice_RawWidgetAction_845186e1(enc *codegen.Encoder, arg []remotewidgetservice.RawWidgetAction) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_RawWidgetAction_845186e1(dec *codegen.Decoder) []remotewidgetservice.RawWidgetAction {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]remotewidgetservice.RawWidgetAction, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

func serviceweaver_enc_slice_byte_87461245(enc *codegen.Encoder, arg []byte) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		enc.Byte(arg[i])
	}
}

func serviceweaver_dec_slice_byte_87461245(dec *codegen.Decoder) []byte {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = dec.Byte()
	}
	return res
}

func serviceweaver_enc_map_string_slice_byte_7ebbaefa(enc *codegen.Encoder, arg map[string][]byte) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for k, v := range arg {
		enc.String(k)
		serviceweaver_enc_slice_byte_87461245(enc, v)
	}
}

func serviceweaver_dec_map_string_slice_byte_7ebbaefa(dec *codegen.Decoder) map[string][]byte {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make(map[string][]byte, n)
	var k string
	var v []byte
	for i := 0; i < n; i++ {
		k = dec.String()
		v = serviceweaver_dec_slice_byte_87461245(dec)
		res[k] = v
	}
	return res
}
