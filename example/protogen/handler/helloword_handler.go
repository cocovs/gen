package handler

import (
	"context"
	pb "library/proto"
	svc "library/svc"
)

type HelloWordGrpcHandler struct {
	HelloWordSvc svc.HelloWordGrpcService
}

func (h *HelloWordGrpcHandler) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloRequest, error) {
	h.HelloWordSvc.Hello(ctx)
	return nil, nil
}
