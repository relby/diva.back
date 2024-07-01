package interceptor

import (
	"context"
	"errors"

	"github.com/relby/diva.back/internal/domainerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorInterceptor struct{}

func NewErrorInterceptor() *ErrorInterceptor {
	return &ErrorInterceptor{}
}

func (interceptor *ErrorInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err == nil {
		return res, nil
	}
	if _, ok := status.FromError(err); ok {
		return nil, err
	}

	var notFoundError *domainerrors.NotFoundError
	var validationError *domainerrors.ValidationError

	if errors.As(err, &notFoundError) {
		return nil, status.Errorf(codes.NotFound, "not found: %v", err)
	}

	if errors.As(err, &validationError) {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	return nil, status.Errorf(codes.Unknown, "unknown error: %v", err)
}
