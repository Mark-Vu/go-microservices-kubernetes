package grpc

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
	pb "ride-sharing/shared/proto/driver/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	service domain.DriverService
}

func NewGRPCHandler(server *grpc.Server, service domain.DriverService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	if req.DriverID == "" {
		return nil, status.Error(codes.InvalidArgument, "driverID is required")
	}
	if req.PackageSlug == "" {
		return nil, status.Error(codes.InvalidArgument, "packageSlug is required")
	}

	// Call the service layer
	driver, err := h.service.RegisterDriver(ctx, req.DriverID, req.PackageSlug)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver: %v", err)
	}

	// Convert domain model to proto response
	return &pb.RegisterDriverResponse{
		Driver: ToProtoDriver(driver),
	}, nil
}

func (h *gRPCHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	// TODO: Implement unregister logic
	return nil, status.Error(codes.Unimplemented, "method UnregisterDriver not implemented")
}
