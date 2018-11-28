// Copyright 2017 David Ackroyd. All Rights Reserved.
// See LICENSE for licensing terms.

package recovery

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (err error)

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(r, o.recoveryHandlerFunc)
			}
		}()

		return handler(ctx, req)
	}
}

func recoverFrom(p interface{}, r RecoveryHandlerFunc) error {
	if r == nil {
		return grpc.Errorf(codes.Internal, "%s", p)
	}
	return r(p)
}