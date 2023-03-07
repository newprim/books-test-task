package bookhandler

import (
	"net/http"

	"github.com/newprim/books-test-task/internal/usecase"
	"github.com/newprim/books-test-task/pkg/log"
)

type Handlers struct {
	book usecase.Book
	l    log.Interface
}

const (
	_getSome3 = "/v1/getSome3"
	_delete   = "/v1/delete"
)

func InitHandlers(book usecase.Book, mux *http.ServeMux, logger log.Interface) {
	h := &Handlers{
		book: book,
		l:    logger,
	}

	mux.HandleFunc(_getSome3, h.GetSome3)
	mux.HandleFunc(_delete, h.Delete)
}
