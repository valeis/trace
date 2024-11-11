package security

import (
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/internal/security/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthServiceClient(c *config.Config) (pb.SecurityServiceClient, error) {
	cc, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewSecurityServiceClient(cc), nil
}
