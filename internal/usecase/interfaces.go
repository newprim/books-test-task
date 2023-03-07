package usecase

import (
	"context"

	"github.com/newprim/books-test-task/internal/dto"
	"github.com/newprim/books-test-task/internal/entity"
)

// Book представляет интерфейс для взаимодействия с use cases.
type Book interface {
	Get(ctx context.Context, id int) (entity.Book, error)
	GetSome(ctx context.Context, count int) ([]entity.Book, error)
	Delete(ctx context.Context, id int) error
}

// BookRepository представляет интерфейс для получения данных из репозитория.
type BookRepository interface {
	Get(ctx context.Context, id int) (dto.Book, error)
	GetRandomN(ctx context.Context, n int) ([]dto.Book, error)
	Delete(ctx context.Context, id int) error
}
