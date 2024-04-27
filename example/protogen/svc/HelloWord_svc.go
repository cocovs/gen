package svc

import "context"

type HelloWordGrpcService struct {
}

func (svc *HelloWordGrpcService) Hello(ctx context.Context) error {
	return nil
}
