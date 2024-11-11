package main

import (
	"Booking_system/security_service/internal/config"
	security_controller "Booking_system/security_service/internal/controller"
	"Booking_system/security_service/internal/pb"
	security_repository "Booking_system/security_service/internal/repository"
	security_service "Booking_system/security_service/internal/service"
	"Booking_system/security_service/pkg/database/redis_repository"
	"Booking_system/user_service/pkg/postgres"
	"fmt"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"os"
)

func main() {
	RunSecurityService()
}

func RunSecurityService() {
	cfg := config.LoadConfig()
	redisDB, postgresDB := dbConnection(cfg)
	userRepo := security_repository.NewUserRepository(postgresDB)
	redisRepo := security_repository.NewRedisRepository(redisDB)
	sService := security_service.NewSecurityService(userRepo, redisRepo)
	sRpcServer := security_controller.NewSecurityController(sService, cfg)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.GrpcPort))
	if err != nil {
		slog.Fatalf("Failed to listen to security service on GRPC port %v: %v, cfg.GrpcPort, err")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSecurityServiceServer(grpcServer, sRpcServer)

	slog.Infof("Listening security on %v", cfg.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		slog.Fatalf("failed to serve security service on %v: %v", cfg.GrpcPort, err)
	}
}

func dbConnection(cfg *config.Config) (*redis.Client, *gorm.DB) {
	postgresDB := postgres.Connect()
	redisDB := redis_repository.NewRedisDBConnection()
	return redisDB, postgresDB
}

func fetchEnvVariables() (string, string) {
	if err := godotenv.Load(".env"); err != nil {
		slog.Fatalf("Failed to read .env file: %v", err)
	}
	return os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS")
}
