package application

import (
	"log/slog"
	"os"

	grpcserver "github.com/neandrson/go-daev2-final/orchestrator/internal/application/grpc"
	httpserver "github.com/neandrson/go-daev2-final/orchestrator/internal/application/http"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/config"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/services/auth"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/services/expression"
	"github.com/neandrson/go-daev2-final/orchestrator/internal/storage"
)

func setUpLogger(logFile *os.File) error {
	// Инициализация логера в файл
	opts := slog.HandlerOptions{}
	var logger = slog.New(slog.NewTextHandler(logFile, &opts))
	slog.SetDefault(logger)
	return nil
}

type Application struct {
	config   *config.Config
	service  *expression.ExpressionService // здесь только для graceful shutdown
	forTests bool
}

func New(forTests bool) *Application {
	config, err := config.ConfigFromEnv()
	if err != nil {
		panic(err)
	}

	return &Application{
		config:   config,
		forTests: forTests,
	}
}

func (a *Application) RunServer() error {
	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error while opening log file", "error", err)
	}
	defer logFile.Close()
	err = setUpLogger(logFile)
	if err != nil {
		slog.Error("Error while setting up logger", "error", err)
	}

	// Создаем хранилище, которое будет передаваться вглубь приложения по ссылке,
	// то есть все сервисы будут работать с одним и тем же хранилищем
	storage := storage.NewStorage(a.forTests)

	// А вот и сервис по работе с выражениями. Он используется в хендлерах для обработки запросов
	expressionService := expression.NewExpressionService(storage, a.config.TimeConf)
	a.service = expressionService
	// Сервис авторизации
	authService := auth.NewAuthService(storage, []byte(a.config.SecretKey))

	// Запуск HTTP и gRPC серверов в разных горутинах
	go httpserver.RunHTTPServer(authService, expressionService, *a.config)
	grpcserver.RunGRPCServer(expressionService, *a.config)
	return nil
}

func (a *Application) Close() {
	slog.Info("Application shutdown")
	a.service.Close()
}
