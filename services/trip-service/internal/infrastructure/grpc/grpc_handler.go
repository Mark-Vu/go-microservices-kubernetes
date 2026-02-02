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
	userId := req.GetUserID()

	route, err := h.service.GetRoute(ctx, pickUpCoordinate, destinationCoordinate)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares, err := h.service.EstimatePackagesPriceWithRoute(ctx, route)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to estimate packages price: %v", err)
	}

	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userId, route)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to generate trip fares: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     osrmToProtoRoute(route),
		RideFares: ToProtoRideFares(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareId := req.GetRideFareID()
	userId := req.GetUserID()
	// 1. Fetch and validate ride fare
	fare, err := h.service.GetRideFareByID(ctx, fareId, userId)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get ride fare: %v", err)
	}
	log.Printf("starting GenerateTripFares %v", fare)
	// 2. Create trip
	trip, err := h.service.CreateTrip(ctx, fare)
	log.Printf("created trip %v", trip)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}
	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
