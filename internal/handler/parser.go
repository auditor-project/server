package handler

import (
	"context"

	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/proto"
	"auditor.z9fr.xyz/server/internal/worker"
	"github.com/google/uuid"
)

type ParserHandler struct {
	log         lib.Logger
	env         *lib.Env
	distributor worker.TaskDistributor
	proto.UnsafeParserHandlerServiceServer
}

func NewParserHandlerImpl(
	log lib.Logger,
	env *lib.Env,
	distributor worker.TaskDistributor,
) *ParserHandler {
	return &ParserHandler{
		log:         log,
		env:         env,
		distributor: distributor,
	}
}

func (c *ParserHandler) AuditStartProcessor(ctx context.Context, req *proto.AuditStartRequest) (*proto.AuditStartResponse, error) {
	requestId := uuid.New()
	ctx = context.WithValue(ctx, "requestId", requestId.String())

	c.distributor.DistributeTaskAuditStartRequest(ctx, req)

	return &proto.AuditStartResponse{
		RequestId: requestId.String(),
	}, nil

}
