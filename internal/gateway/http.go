package gateway

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	pb "url-shortener/protos/gen/go"
)

func RunHTTPServer(httpAddr, grpcAddr string) error {
	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterUrlShortenerHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}

	log.Printf("HTTP server listening on %s", httpAddr)
	return http.ListenAndServe(httpAddr, mux)
}
