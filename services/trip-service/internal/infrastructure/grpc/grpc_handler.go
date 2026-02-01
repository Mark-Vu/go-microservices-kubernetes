package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickUpCoordinate := protoToCoordinate(req.GetStartLocation())
	destinationCoordinate := protoToCoordinate(req.GetEndLocation())

	trip, err := h.service.GetRoute(ctx, pickUpCoordinate, destinationCoordinate)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares, err := h.service.EstimatePackagesPriceWithRoute(ctx, trip)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to estimate packages price: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     osrmToProtoRoute(trip),
		RideFares: ToProtoRideFares(estimatedFares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTrip not implemented")
}
