package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/neandrson/go-daev2-final/orchestrator/internal/config"
	grpctasks "github.com/neandrson/go-daev2-final/orchestrator/internal/grpc"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/services/expression"
)

func RunGRPCServer(expressionService *expression.ExpressionService, config config.Config) {
	host := "localhost"
	port := config.GRPCPort

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	slog.Info("tcp listener started", "port", port)
	grpcServer := grpc.NewServer()
	grpctasks.Register(grpcServer, expressionService)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
