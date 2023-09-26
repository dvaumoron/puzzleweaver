// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package saltimpl

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
		Name:  "github.com/dvaumoron/puzzleweaver/serviceimpl/salt/SaltService",
		Iface: reflect.TypeOf((*SaltService)(nil)).Elem(),
		Impl:  reflect.TypeOf(saltImpl{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return saltService_local_stub{impl: impl.(SaltService), tracer: tracer, loadOrGenerateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/salt/SaltService", Method: "LoadOrGenerate", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return saltService_client_stub{stub: stub, loadOrGenerateMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/salt/SaltService", Method: "LoadOrGenerate", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return saltService_server_stub{impl: impl.(SaltService), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return saltService_reflect_stub{caller: caller}
		},
		RefData: "",
	})
}

// weaver.InstanceOf checks.
var _ weaver.InstanceOf[SaltService] = (*saltImpl)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*saltImpl)(nil)

// Local stub implementations.

type saltService_local_stub struct {
	impl                  SaltService
	tracer                trace.Tracer
	loadOrGenerateMetrics *codegen.MethodMetrics
}

// Check that saltService_local_stub implements the SaltService interface.
var _ SaltService = (*saltService_local_stub)(nil)

func (s saltService_local_stub) LoadOrGenerate(ctx context.Context, a0 ...string) (r0 [][]byte, err error) {
	// Update metrics.
	begin := s.loadOrGenerateMetrics.Begin()
	defer func() { s.loadOrGenerateMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "saltimpl.SaltService.LoadOrGenerate", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.LoadOrGenerate(ctx, a0...)
}

// Client stub implementations.

type saltService_client_stub struct {
	stub                  codegen.Stub
	loadOrGenerateMetrics *codegen.MethodMetrics
}

// Check that saltService_client_stub implements the SaltService interface.
var _ SaltService = (*saltService_client_stub)(nil)

func (s saltService_client_stub) LoadOrGenerate(ctx context.Context, a0 ...string) (r0 [][]byte, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.loadOrGenerateMetrics.Begin()
	defer func() { s.loadOrGenerateMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "saltimpl.SaltService.LoadOrGenerate", trace.WithSpanKind(trace.SpanKindClient))
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
	serviceweaver_enc_slice_string_4af10117(enc, a0)
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
	r0 = serviceweaver_dec_slice_slice_byte_8acc26ee(dec)
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

type saltService_server_stub struct {
	impl    SaltService
	addLoad func(key uint64, load float64)
}

// Check that saltService_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*saltService_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s saltService_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "LoadOrGenerate":
		return s.loadOrGenerate
	default:
		return nil
	}
}

func (s saltService_server_stub) loadOrGenerate(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 []string
	a0 = serviceweaver_dec_slice_string_4af10117(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.LoadOrGenerate(ctx, a0...)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_slice_slice_byte_8acc26ee(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

// Reflect stub implementations.

type saltService_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that saltService_reflect_stub implements the SaltService interface.
var _ SaltService = (*saltService_reflect_stub)(nil)

func (s saltService_reflect_stub) LoadOrGenerate(ctx context.Context, a0 ...string) (r0 [][]byte, err error) {
	err = s.caller("LoadOrGenerate", ctx, []any{a0}, []any{&r0})
	return
}

// Encoding/decoding implementations.

func serviceweaver_enc_slice_string_4af10117(enc *codegen.Encoder, arg []string) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		enc.String(arg[i])
	}
}

func serviceweaver_dec_slice_string_4af10117(dec *codegen.Decoder) []string {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = dec.String()
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

func serviceweaver_enc_slice_slice_byte_8acc26ee(enc *codegen.Encoder, arg [][]byte) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		serviceweaver_enc_slice_byte_87461245(enc, arg[i])
	}
}

func serviceweaver_dec_slice_slice_byte_8acc26ee(dec *codegen.Decoder) [][]byte {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([][]byte, n)
	for i := 0; i < n; i++ {
		res[i] = serviceweaver_dec_slice_byte_87461245(dec)
	}
	return res
}