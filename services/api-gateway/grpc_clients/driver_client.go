package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/driver/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DriverServiceClient struct {
	Client pb.DriverServiceClient
	conn   *grpc.ClientConn
}

func NewDriverServiceClient() (*DriverServiceClient, error) {
	driverServiceUrl := os.Getenv("DRIVER_SERVICE_URL")
	if driverServiceUrl == "" {
		driverServiceUrl = "driver-service:8083"
	}

	conn, err := grpc.NewClient(driverServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := pb.NewDriverServiceClient(conn)

	return &DriverServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *DriverServiceClient) Close() error {
	return c.conn.Close()
}
