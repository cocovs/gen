package handler

import (
	"context"
	pb "library/proto"
	svc "library/svc"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type LibraryGrpcHandler struct {
	LibrarySvc svc.LibraryGrpcService
}

func (h *LibraryGrpcHandler) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksRequest, error) {
	h.LibrarySvc.ListBooks(ctx)
	return nil, nil
}

func (h *LibraryGrpcHandler) GetBook(ctx context.Context, in *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	h.LibrarySvc.GetBook(ctx)
	return nil, nil
}

func (h *LibraryGrpcHandler) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.CreateBookRequest, error) {
	h.LibrarySvc.CreateBook(ctx)
	return nil, nil
}

func (h *LibraryGrpcHandler) UpdateBook(ctx context.Context, in *pb.UpdateBookRequest) (*pb.UpdateBookRequest, error) {
	h.LibrarySvc.UpdateBook(ctx)
	return nil, nil
}

func (h *LibraryGrpcHandler) DeleteBook(ctx context.Context, in *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	h.LibrarySvc.DeleteBook(ctx)
	return nil, nil
}
