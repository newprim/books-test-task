package books

import (
	"context"

	"github.com/newprim/books-test-task/internal/entity"
)

// TODO(07.03.2023, Гурьянов Роман): replace.
type Repository interface {
	Get(ctx context.Context, id int) (entity.Book, error)
	GetSome(ctx context.Context, count int) ([]entity.Book, error)
	Delete(ctx context.Context, id int) error
}

type Handlers struct {
	rep Repository
}
