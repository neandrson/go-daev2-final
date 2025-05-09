package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/clients/sso/grpc"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/domain/models"
	middleware "github.com/neandrson/go-daev2-final/orchestrator/internal/http-server/middlewares"
	httpservice "github.com/neandrson/go-daev2-final/orchestrator/internal/service/http"
)

type Server struct {
	log         *slog.Logger
	storage     httpservice.ExpressionStorage
	secret      string
	server      *http.Server
	httpService HttpService
}

type HttpService interface {
	GetExpressionById(ctx context.Context, id string, uid int) (*models.Expression, error)
	EvaluateExpression(ctx context.Context, expression *models.Expression, uid int) (string, error) //float32, error
	GetAgentStates(ctx context.Context) ([]models.Agent, error)
	GetExpressionsForUser(ctx context.Context, uid int) ([]models.Expression, error)
	Login(ctx context.Context, login string, password string) (string, error)
	Register(ctx context.Context, login string, password string) (int, error)
}

func New(log *slog.Logger, storage httpservice.ExpressionStorage, port int, secret string, client *grpc.Client) *Server {
	httpService := httpservice.New(log, storage, client)

	server := &Server{
		log:         log,
		storage:     storage,
		secret:      secret,
		httpService: httpService,
	}
	serveMux := mux.NewRouter()

	// Хэндлеры для запросов с сайта
	serveMux.HandleFunc("/", server.MainPage)
	serveMux.Handle("/api/v1/calculate", middleware.ValidateToken(middleware.ValidateExpressionMiddleware(http.HandlerFunc(server.EvaluateExpression)), server.secret)).Methods("POST")
	serveMux.Handle("/api/v1/expression", middleware.ValidateToken(http.HandlerFunc(server.GetExpressionById), server.secret)).Methods("GET")
	serveMux.Handle("/api/v1/expressions", middleware.ValidateToken(http.HandlerFunc(server.GetExpressionsForUser), server.secret)).Methods("GET")
	serveMux.HandleFunc("/api/v1/login", server.Login).Methods("POST")
	serveMux.HandleFunc("/api/v1/register", server.Register).Methods("POST")
	serveMux.HandleFunc("/internal/task", server.GetAgentStates).Methods("GET")

	// Хэндлеры для API
	serveMux.HandleFunc("/internal/task", server.GetImpodenceKeyHandler).Methods("POST")

	http.Handle("/", serveMux)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serveMux,
	}

	server.server = s
	return server
}

func (s *Server) Run() error {
	const op = "server.Run"

	l, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = s.server.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
