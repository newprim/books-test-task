package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/newprim/books-test-task/config"
	"github.com/newprim/books-test-task/internal/controller/http/v1/bookhandler"
	"github.com/newprim/books-test-task/internal/repository/runtimerep"
	"github.com/newprim/books-test-task/internal/usecase"
	"github.com/newprim/books-test-task/pkg/httpserver"
	"github.com/newprim/books-test-task/pkg/log"
	"github.com/newprim/books-test-task/pkg/middlewares"
)

func Run(cfg config.Config) {
	logger := log.New(cfg.Log.Level)
	if err := run(cfg, logger); err != nil {
		logger.Fatal("starting service: %v", err)
	}

	logger.Info("service stopped")
}

func run(cfg config.Config, logger log.Interface) error {
	booksRep, err := runtimerep.NewFakeRepository(10)
	if err != nil {
		return fmt.Errorf("creating repository: %w", err)
	}

	bookUseCase := usecase.NewBookUseCase(booksRep)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	throtMid := middlewares.NewRejectionThrottling(ctx, cfg.HTTP.MaxPRS, cfg.HTTP.Duration)
	handler := http.NewServeMux()
	bookhandler.AddHandlersToMux(bookUseCase, handler, logger, throtMid)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger.Info("starting service om port %v", cfg.HTTP.Port)

	select {
	case s := <-interrupt:
		logger.Info("stopping: %v", s)
	case err = <-httpServer.Notify():
		logger.Error("server notifies: %v", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	return nil
}
