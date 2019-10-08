package goerrors

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
)

func GRPCUnaryServerRecoverer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer OnRecover(func(err error) {
			Log().
				WithContext(CreateContext(ctx).AddStack(debug.Stack())).
				WithError(err).
				Error("unary GRPC recovered")
		})

		return handler(ctx, req)
	}
}

func GRPCStreamServerRecoverer() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer OnRecover(func(err error) {
			Log().
				WithContext(CreateContext(ss.Context()).AddStack(debug.Stack())).
				WithError(err).
				Error("stream GRPC recovered")
		})

		return handler(srv, ss)
	}
}
