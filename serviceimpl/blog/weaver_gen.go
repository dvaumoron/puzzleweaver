// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package blogimpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:  "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService",
		Iface: reflect.TypeOf((*RemoteBlogService)(nil)).Elem(),
		Impl:  reflect.TypeOf(remoteBlogImpl{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return remoteBlogService_local_stub{impl: impl.(RemoteBlogService), tracer: tracer, createPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "CreatePost", Remote: false}), deleteMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "Delete", Remote: false}), getPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "GetPost", Remote: false}), getPostsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "GetPosts", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return remoteBlogService_client_stub{stub: stub, createPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "CreatePost", Remote: true}), deleteMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "Delete", Remote: true}), getPostMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "GetPost", Remote: true}), getPostsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/dvaumoron/puzzleweaver/serviceimpl/blog/RemoteBlogService", Method: "GetPosts", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return remoteBlogService_server_stub{impl: impl.(RemoteBlogService), addLoad: addLoad}
		},
		ReflectStubFn: func(caller func(string, context.Context, []any, []any) error) any {
			return remoteBlogService_reflect_stub{caller: caller}
		},
		RefData: "",
	})
}

// weaver.InstanceOf checks.
var _ weaver.InstanceOf[RemoteBlogService] = (*remoteBlogImpl)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*remoteBlogImpl)(nil)

// Local stub implementations.

type remoteBlogService_local_stub struct {
	impl              RemoteBlogService
	tracer            trace.Tracer
	createPostMetrics *codegen.MethodMetrics
	deleteMetrics     *codegen.MethodMetrics
	getPostMetrics    *codegen.MethodMetrics
	getPostsMetrics   *codegen.MethodMetrics
}

// Check that remoteBlogService_local_stub implements the RemoteBlogService interface.
var _ RemoteBlogService = (*remoteBlogService_local_stub)(nil)

func (s remoteBlogService_local_stub) CreatePost(ctx context.Context, a0 uint64, a1 uint64, a2 string, a3 string) (r0 uint64, err error) {
	// Update metrics.
	begin := s.createPostMetrics.Begin()
	defer func() { s.createPostMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "blogimpl.RemoteBlogService.CreatePost", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.CreatePost(ctx, a0, a1, a2, a3)
}

func (s remoteBlogService_local_stub) Delete(ctx context.Context, a0 uint64, a1 uint64) (err error) {
	// Update metrics.
	begin := s.deleteMetrics.Begin()
	defer func() { s.deleteMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "blogimpl.RemoteBlogService.Delete", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Delete(ctx, a0, a1)
}

func (s remoteBlogService_local_stub) GetPost(ctx context.Context, a0 uint64, a1 uint64) (r0 RawBlogPost, err error) {
	// Update metrics.
	begin := s.getPostMetrics.Begin()
	defer func() { s.getPostMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "blogimpl.RemoteBlogService.GetPost", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetPost(ctx, a0, a1)
}

func (s remoteBlogService_local_stub) GetPosts(ctx context.Context, a0 uint64, a1 uint64, a2 uint64, a3 string) (r0 uint64, r1 []RawBlogPost, err error) {
	// Update metrics.
	begin := s.getPostsMetrics.Begin()
	defer func() { s.getPostsMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "blogimpl.RemoteBlogService.GetPosts", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.GetPosts(ctx, a0, a1, a2, a3)
}

// Client stub implementations.

type remoteBlogService_client_stub struct {
	stub              codegen.Stub
	createPostMetrics *codegen.MethodMetrics
	deleteMetrics     *codegen.MethodMetrics
	getPostMetrics    *codegen.MethodMetrics
	getPostsMetrics   *codegen.MethodMetrics
}

// Check that remoteBlogService_client_stub implements the RemoteBlogService interface.
var _ RemoteBlogService = (*remoteBlogService_client_stub)(nil)

func (s remoteBlogService_client_stub) CreatePost(ctx context.Context, a0 uint64, a1 uint64, a2 string, a3 string) (r0 uint64, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.createPostMetrics.Begin()
	defer func() { s.createPostMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "blogimpl.RemoteBlogService.CreatePost", trace.WithSpanKind(trace.SpanKindClient))
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
	size += 8
	size += (4 + len(a2))
	size += (4 + len(a3))
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Uint64(a0)
	enc.Uint64(a1)
	enc.String(a2)
	enc.String(a3)
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
	r0 = dec.Uint64()
	err = dec.Error()
	return
}

func (s remoteBlogService_client_stub) Delete(ctx context.Context, a0 uint64, a1 uint64) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.deleteMetrics.Begin()
	defer func() { s.deleteMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "blogimpl.RemoteBlogService.Delete", trace.WithSpanKind(trace.SpanKindClient))
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
	size += 8
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Uint64(a0)
	enc.Uint64(a1)
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

func (s remoteBlogService_client_stub) GetPost(ctx context.Context, a0 uint64, a1 uint64) (r0 RawBlogPost, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getPostMetrics.Begin()
	defer func() { s.getPostMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "blogimpl.RemoteBlogService.GetPost", trace.WithSpanKind(trace.SpanKindClient))
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
	size += 8
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Uint64(a0)
	enc.Uint64(a1)
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
	(&r0).WeaverUnmarshal(dec)
	err = dec.Error()
	return
}

func (s remoteBlogService_client_stub) GetPosts(ctx context.Context, a0 uint64, a1 uint64, a2 uint64, a3 string) (r0 uint64, r1 []RawBlogPost, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.getPostsMetrics.Begin()
	defer func() { s.getPostsMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "blogimpl.RemoteBlogService.GetPosts", trace.WithSpanKind(trace.SpanKindClient))
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
	size += 8
	size += 8
	size += (4 + len(a3))
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.Uint64(a0)
	enc.Uint64(a1)
	enc.Uint64(a2)
	enc.String(a3)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 3, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = dec.Uint64()
	r1 = serviceweaver_dec_slice_RawBlogPost_43e1e8b0(dec)
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

type remoteBlogService_server_stub struct {
	impl    RemoteBlogService
	addLoad func(key uint64, load float64)
}

// Check that remoteBlogService_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*remoteBlogService_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s remoteBlogService_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "CreatePost":
		return s.createPost
	case "Delete":
		return s.delete
	case "GetPost":
		return s.getPost
	case "GetPosts":
		return s.getPosts
	default:
		return nil
	}
}

func (s remoteBlogService_server_stub) createPost(ctx context.Context, args []byte) (res []byte, err error) {
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
	var a1 uint64
	a1 = dec.Uint64()
	var a2 string
	a2 = dec.String()
	var a3 string
	a3 = dec.String()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.CreatePost(ctx, a0, a1, a2, a3)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Uint64(r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s remoteBlogService_server_stub) delete(ctx context.Context, args []byte) (res []byte, err error) {
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
	var a1 uint64
	a1 = dec.Uint64()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.Delete(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s remoteBlogService_server_stub) getPost(ctx context.Context, args []byte) (res []byte, err error) {
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
	var a1 uint64
	a1 = dec.Uint64()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.GetPost(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	(r0).WeaverMarshal(enc)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s remoteBlogService_server_stub) getPosts(ctx context.Context, args []byte) (res []byte, err error) {
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
	var a1 uint64
	a1 = dec.Uint64()
	var a2 uint64
	a2 = dec.Uint64()
	var a3 string
	a3 = dec.String()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, r1, appErr := s.impl.GetPosts(ctx, a0, a1, a2, a3)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Uint64(r0)
	serviceweaver_enc_slice_RawBlogPost_43e1e8b0(enc, r1)
	enc.Error(appErr)
	return enc.Data(), nil
}

// Reflect stub implementations.

type remoteBlogService_reflect_stub struct {
	caller func(string, context.Context, []any, []any) error
}

// Check that remoteBlogService_reflect_stub implements the RemoteBlogService interface.
var _ RemoteBlogService = (*remoteBlogService_reflect_stub)(nil)

func (s remoteBlogService_reflect_stub) CreatePost(ctx context.Context, a0 uint64, a1 uint64, a2 string, a3 string) (r0 uint64, err error) {
	err = s.caller("CreatePost", ctx, []any{a0, a1, a2, a3}, []any{&r0})
	return
}

func (s remoteBlogService_reflect_stub) Delete(ctx context.Context, a0 uint64, a1 uint64) (err error) {
	err = s.caller("Delete", ctx, []any{a0, a1}, []any{})
	return
}

func (s remoteBlogService_reflect_stub) GetPost(ctx context.Context, a0 uint64, a1 uint64) (r0 RawBlogPost, err error) {
	err = s.caller("GetPost", ctx, []any{a0, a1}, []any{&r0})
	return
}

func (s remoteBlogService_reflect_stub) GetPosts(ctx context.Context, a0 uint64, a1 uint64, a2 uint64, a3 string) (r0 uint64, r1 []RawBlogPost, err error) {
	err = s.caller("GetPosts", ctx, []any{a0, a1, a2, a3}, []any{&r0, &r1})
	return
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = (*RawBlogPost)(nil)

type __is_RawBlogPost[T ~struct {
	weaver.AutoMarshal
	Id        uint64
	CreatorId uint64
	CreatedAt int64
	Title     string
	Content   string
}] struct{}

var _ __is_RawBlogPost[RawBlogPost]

func (x *RawBlogPost) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("RawBlogPost.WeaverMarshal: nil receiver"))
	}
	enc.Uint64(x.Id)
	enc.Uint64(x.CreatorId)
	enc.Int64(x.CreatedAt)
	enc.String(x.Title)
	enc.String(x.Content)
}

func (x *RawBlogPost) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("RawBlogPost.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.Uint64()
	x.CreatorId = dec.Uint64()
	x.CreatedAt = dec.Int64()
	x.Title = dec.String()
	x.Content = dec.String()
}

var _ codegen.AutoMarshal = (*RawForumContent)(nil)

type __is_RawForumContent[T ~struct {
	weaver.AutoMarshal
	Id        uint64
	CreatorId uint64
	CreatedAt int64
	Text      string
}] struct{}

var _ __is_RawForumContent[RawForumContent]

func (x *RawForumContent) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("RawForumContent.WeaverMarshal: nil receiver"))
	}
	enc.Uint64(x.Id)
	enc.Uint64(x.CreatorId)
	enc.Int64(x.CreatedAt)
	enc.String(x.Text)
}

func (x *RawForumContent) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("RawForumContent.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.Uint64()
	x.CreatorId = dec.Uint64()
	x.CreatedAt = dec.Int64()
	x.Text = dec.String()
}

var _ codegen.AutoMarshal = (*RawWikiContent)(nil)

type __is_RawWikiContent[T ~struct {
	weaver.AutoMarshal
	Version   uint64
	CreatorId uint64
	CreatedAt int64
	Markdown  string
}] struct{}

var _ __is_RawWikiContent[RawWikiContent]

func (x *RawWikiContent) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("RawWikiContent.WeaverMarshal: nil receiver"))
	}
	enc.Uint64(x.Version)
	enc.Uint64(x.CreatorId)
	enc.Int64(x.CreatedAt)
	enc.String(x.Markdown)
}

func (x *RawWikiContent) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("RawWikiContent.WeaverUnmarshal: nil receiver"))
	}
	x.Version = dec.Uint64()
	x.CreatorId = dec.Uint64()
	x.CreatedAt = dec.Int64()
	x.Markdown = dec.String()
}

// Encoding/decoding implementations.

func serviceweaver_enc_slice_RawBlogPost_43e1e8b0(enc *codegen.Encoder, arg []RawBlogPost) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_RawBlogPost_43e1e8b0(dec *codegen.Decoder) []RawBlogPost {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]RawBlogPost, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}
