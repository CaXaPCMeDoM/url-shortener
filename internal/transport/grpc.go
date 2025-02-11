package transport

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"url-shortener/internal/config"
	"url-shortener/internal/service/shortener"
	pb "url-shortener/protos/gen/go"
)

func RunGRPCServer(grpcAddr string, shortenerService *shortener.Service, cfg config.Config) error {
	configuration(cfg)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUrlShortenerServer(grpcServer, shortenerService)

	log.Printf("gRPC server listening on %s", grpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
		return err
	}

	return nil
}

func configuration(cfg config.Config) {
	runtime.DefaultContextTimeout = cfg.Grpc.Timeout
}
