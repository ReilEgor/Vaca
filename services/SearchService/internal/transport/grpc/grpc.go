package grpc

import (
	"fmt"
	"net"
	"strconv"

	elastic "github.com/ReilEgor/Vaca/services/SearchService/api/proto/gen/elastic"
	"github.com/ReilEgor/Vaca/services/SearchService/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	grpcServer *grpc.Server
	port       int
}

func NewServer(port config.GRPCPort, handler elastic.SearchServiceServer) *Server {
	s := grpc.NewServer()

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	elastic.RegisterSearchServiceServer(s, handler)
	healthpb.RegisterHealthServer(s, healthServer)

	grpcPortInt, _ := strconv.Atoi(string(port))
	return &Server{
		grpcServer: s,
		port:       grpcPortInt,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
