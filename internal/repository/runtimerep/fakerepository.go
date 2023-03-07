package runtimerep

import (
	"context"
	"fmt"
	"strconv"

	"github.com/newprim/books-test-task/internal/dto"
	"github.com/newprim/books-test-task/internal/usecase"
	runtime "github.com/newprim/books-test-task/pkg/mutexmap"
)

type FakeRepository struct {
	cache *runtime.MutexMap[int, dto.Book]
}

var _ usecase.BookRepository = (*FakeRepository)(nil)

func NewFakeRepository(countOfBooks int) (*FakeRepository, error) {
	if countOfBooks < 0 {
		return nil, fmt.Errorf("books count is less than 0: %d", countOfBooks)
	}

	cache := make(map[int]dto.Book, countOfBooks)
	for i := 0; i < countOfBooks; i++ {
		strI := strconv.Itoa(i)
		cache[i] = dto.Book{
			ID:            i,
			Author:        "Author " + strI,
			Title:         "Title " + strI,
			PublisherYear: 2022 - i,
		}
	}

	return &FakeRepository{
		cache: runtime.NewMutexMapFilled(cache),
	}, nil
}

func (c *FakeRepository) Get(_ context.Context, id int) (dto.Book, error) {
	value, ok := c.cache.GetOK(id)
	if !ok {
		return dto.Book{}, fmt.Errorf("no rows in result set")
	}
	return value, nil
}

func (c *FakeRepository) GetRandomN(_ context.Context, n int) ([]dto.Book, error) {
	values := c.cache.GetSomeValues(n)
	return values, nil
}

func (c *FakeRepository) Delete(_ context.Context, id int) error {
	if deleted := c.cache.DeleteIfExist(id); !deleted {
		return fmt.Errorf("no rows in result set")
	}

	return nil
}
