package bookhandler

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/newprim/books-test-task/internal/entity"
	"github.com/newprim/books-test-task/pkg/middlewares"
)

const _seed = 1

func TestHandlers(t *testing.T) {
	type params struct {
		method       string
		path         string
		expectedCode int
		expectedBody string
		Headers      http.Header
		timeout      time.Duration
	}
	tests := []struct {
		name   string
		params params
	}{
		{
			name: _getSome3 + "Correct",
			params: params{
				method:       http.MethodGet,
				path:         _getSome3,
				expectedCode: http.StatusOK,
				expectedBody: bodyOfGetSome3(),
			},
		},
		{
			name: _getSome3 + "ErrorStatusMethodNotAllowed",
			params: params{
				method:       http.MethodPost,
				path:         _getSome3,
				expectedCode: http.StatusMethodNotAllowed,
				expectedBody: "",
			},
		},
		{
			name: _delete + "Correct",
			params: params{
				method:       http.MethodDelete,
				path:         _delete,
				expectedCode: http.StatusOK,
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"Accept":       {"application/json"},
				},
				expectedBody: "",
			},
		},
		{
			name: _delete + "ErrorStatusMethodNotAllowed",
			params: params{
				method:       http.MethodPost,
				path:         _delete,
				expectedCode: http.StatusMethodNotAllowed,
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"Accept":       {"application/json"},
				},
				expectedBody: "",
			},
		},
		{
			name: _delete + "CorrectButRated", // пятый тест в таблице, RPS настроен на 4
			params: params{
				method:       http.MethodDelete,
				path:         _delete,
				expectedCode: http.StatusTooManyRequests,
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"Accept":       {"application/json"},
				},
				expectedBody: "",
			},
		},
		{
			name: _delete + "correctAfterWait",
			params: params{
				method:       http.MethodDelete,
				path:         _delete,
				expectedCode: http.StatusOK,
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"Accept":       {"application/json"},
				},
				expectedBody: "",
				timeout:      time.Second,
			},
		},
	}

	const ratePerDuration = 4
	throttling := middlewares.NewRejectionThrottling(context.Background(), ratePerDuration, time.Second/2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			AddHandlersToMux(bookMock{rand.New(rand.NewSource(_seed))}, mux, logMock{}, throttling)

			req := httptest.NewRequest(tt.params.method, tt.params.path, strings.NewReader(bodyOfGetSome3()))
			req.Header = tt.params.Headers
			resp := httptest.NewRecorder()

			if tt.params.timeout != 0 {
				time.Sleep(tt.params.timeout)
			}
			mux.ServeHTTP(resp, req)

			body := resp.Body.String()
			if body != tt.params.expectedBody && tt.params.expectedBody != "" {
				t.Errorf("different body:\ngot:  %s\nwant:%s\n", body, tt.params.expectedBody)
			}
			if resp.Code != tt.params.expectedCode {
				t.Errorf("different code: got: %d, want: %d\n", resp.Code, tt.params.expectedCode)
			}

		})
	}
}

func bodyOfGetSome3() string {
	some, _ := bookMock{rand.New(rand.NewSource(_seed))}.GetSome(nil, _getSome3NeededCount)
	marshal, _ := json.Marshal(getSome3RepackToResponseDTO(some))
	return string(marshal)
}

// Mocks

type bookMock struct {
	rand *rand.Rand
}

func (b bookMock) randBook() entity.Book {
	return entity.Book{
		ID: b.rand.Int(),
	}
}

func (b bookMock) Get(_ context.Context, _ int) (entity.Book, error) {
	b.rand = rand.New(rand.NewSource(1))
	return b.randBook(), nil
}

func (b bookMock) GetSome(_ context.Context, count int) ([]entity.Book, error) {
	result := make([]entity.Book, count)
	for i := range result {
		result[i] = b.randBook()
	}
	return result, nil
}

func (b bookMock) Delete(_ context.Context, _ int) error {
	return nil
}

type logMock struct{}

func (logMock) Debug(_ interface{}, _ ...interface{}) {}
func (logMock) Info(_ string, _ ...interface{})       {}
func (logMock) Warn(_ string, _ ...interface{})       {}
func (logMock) Error(_ interface{}, _ ...interface{}) {}
func (logMock) Fatal(_ interface{}, _ ...interface{}) {}
