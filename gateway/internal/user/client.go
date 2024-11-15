package user

import (
	"Booking_system/gateway/internal/config"
	"Booking_system/gateway/internal/user/pb"
	"fmt"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserServiceClient(c *config.Config) pb.UserServiceClient {
	fmt.Println(c.UserPort)
	cc, err := grpc.Dial(c.UserPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Fatalf("Could not connect: %v", err)
		return nil
	}

	return pb.NewUserServiceClient(cc)
}
