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

func InitHandlers(
	book usecase.Book,
	mux *http.ServeMux,
	logger log.Interface,
	throtMid func(root http.Handler) http.Handler,
) {
	h := &Handlers{
		book: book,
		l:    logger,
	}

	mux.Handle(_getSome3, throtMid(http.HandlerFunc(h.GetSome3)))
	mux.Handle(_delete, throtMid(http.HandlerFunc(h.Delete)))
}
