package handler

import (
	"context"
	"encoding/json"
	"errors"

	"auditor.z9fr.xyz/server/internal/lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorHandler struct {
	logger lib.Logger
	env    *lib.Env
}

type GenerateErrorWithGrpcCodesRequest struct {
	Err     error
	Payload []byte
	Method  string
}

func NewErrorHandler(log lib.Logger, env *lib.Env) *ErrorHandler {
	return &ErrorHandler{
		logger: log,
		env:    env,
	}
}

func (h *ErrorHandler) GenerateErrorWithGrpcCodes(data GenerateErrorWithGrpcCodesRequest) error {
	if data.Err != nil {
		h.logger.Errorw("Error", "details", data.Err.Error(), "payload", data.Payload, "method", data.Method)

		if errors.Is(data.Err, context.Canceled) {
			return status.Error(codes.Canceled, "Request cancelled by client")
		} else if errors.Is(data.Err, context.DeadlineExceeded) {
			return status.Error(codes.DeadlineExceeded, "Request deadline exceeded")
		}
		return status.Error(codes.Unknown, "Unknown error occurred")
	}

	return nil
}

func (h *ErrorHandler) WithErrorHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	payload, err := json.Marshal(req)
	resp, err := handler(ctx, req)

	return resp, h.GenerateErrorWithGrpcCodes(GenerateErrorWithGrpcCodesRequest{
		Err:     err,
		Payload: payload,
		Method:  info.FullMethod,
	})
}
