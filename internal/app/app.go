package app

import (
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

	handler := http.NewServeMux()
	bookhandler.InitHandlers(bookUseCase, handler, logger)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("stopping: %v", s)
	case err = <-httpServer.Notify():
		logger.Error("server notifies: %v", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	return nil
}
