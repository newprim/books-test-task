package usecase

import (
	"context"
	"fmt"

	"github.com/newprim/books-test-task/internal/entity"
)

type BookUseCase struct {
	rep BookRepository
}

var _ Book = (*BookUseCase)(nil)

func (b *BookUseCase) Get(ctx context.Context, id int) (entity.Book, error) {
	bookDTO, err := b.rep.Get(ctx, id)
	if err != nil {
		return entity.Book{}, fmt.Errorf("getting book from repository: %w", err)
	}

	result := entity.Book{
		ID:            bookDTO.ID,
		Author:        bookDTO.Author,
		Title:         bookDTO.Title,
		PublisherYear: bookDTO.PublisherYear,
	}
	return result, nil
}

func (b *BookUseCase) GetSome(ctx context.Context, count int) ([]entity.Book, error) {
	someBooksDTO, err := b.rep.GetRandomN(ctx, count)
	if err != nil {
		return nil, fmt.Errorf("getting books from repository: %w", err)
	}

	result := make([]entity.Book, 0, len(someBooksDTO))
	for _, bookDTO := range someBooksDTO {
		result = append(result, entity.Book{
			ID:            bookDTO.ID,
			Author:        bookDTO.Author,
			Title:         bookDTO.Title,
			PublisherYear: bookDTO.PublisherYear,
		})
	}

	return result, nil
}

func (b *BookUseCase) Delete(ctx context.Context, id int) error {
	if err := b.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting book from repository: %w", err)
	}

	return nil
}
