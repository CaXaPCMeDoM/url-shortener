package transport

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "url-shortener/protos/gen/go"
)

type MockShortenerService struct {
	pb.UnimplementedUrlShortenerServer
	mock.Mock
}

func TestRunGRPCServer(t *testing.T) {
	listener := bufconn.Listen(1024 * 1024)
	mockService := &MockShortenerService{}

	go func() {
		s := grpc.NewServer()
		pb.RegisterUrlShortenerServer(s, mockService)
		if err := s.Serve(listener); err != nil {
			t.Errorf("Failed to serve: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	assert.NoError(t, err, "Should connect without error")
	defer conn.Close()

	client := pb.NewUrlShortenerClient(conn)
	_, err = client.GetOriginalURL(ctx, &pb.GetOriginalURLRequest{Alias: "test"})

	assert.Error(t, err, "Should return error for unimplemented method")
	assert.Contains(t, err.Error(), "not implemented")
}
