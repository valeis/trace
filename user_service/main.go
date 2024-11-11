package main

import (
	"Booking_system/user_service/internal/config"
	"Booking_system/user_service/internal/controller"
	"Booking_system/user_service/internal/pb"
	"Booking_system/user_service/internal/repository"
	"Booking_system/user_service/internal/service"
	"Booking_system/user_service/pkg/postgres"
	"fmt"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
)

var DBPostgres *gorm.DB

func main() {
	DBPostgres = postgres.Connect()
	userRepo := repository.NewUserRepository(DBPostgres)
	userService := service.NewUserService(userRepo)

	if err := startGRPCServer(userService); err != nil {
		slog.Errorf("Failed to start gRPC server: %v", err)
	}
}

func startGRPCServer(userService *service.UserService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Cfg.GrpcPort))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	UserController := controller.NewUserController(userService)
	pb.RegisterUserServiceServer(grpcServer, UserController)

	slog.Info("gRPC server listen on port", config.Cfg.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("Failed to server: #{err}")
	}
	return nil
}
