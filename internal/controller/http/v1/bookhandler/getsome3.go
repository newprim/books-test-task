package bookhandler

import (
	"encoding/json"
	"net/http"
)

func (h *Handlers) GetSome3(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if isValid, code := h.validateGetSome3(r); !isValid {
		http.Error(w, "validation error", code)
		return
	}

	const (
		neededCount = 3
		handler     = "GetSome3"
	)

	firstBooks, err := h.book.GetSome(r.Context(), neededCount)
	if err != nil {
		h.l.Error("getting books on %s: %v", handler, err)
		http.Error(w, "getting some books: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respData := getResponse{
		Books: make([]book, 0, len(firstBooks)),
	}
	for _, b := range firstBooks {
		respData.Books = append(respData.Books, book{
			ID:            b.ID,
			Title:         b.Title,
			Author:        b.Author,
			PublisherYear: b.PublisherYear,
		})
	}

	marshaled, err := json.Marshal(respData)
	if err != nil {
		h.l.Error("marshaling response on %s: %v", handler, err)
		http.Error(w, "marshaling: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(marshaled); err != nil {
		h.l.Error("writing response on %s: %v", handler, err)
	}
}

func (h *Handlers) validateGetSome3(r *http.Request) (bool, int) {
	if r.Method != http.MethodGet {
		return false, http.StatusMethodNotAllowed
	}

	// В задании сказано принимать эти заголовки, но если мы не парсим тело - незачем это делать.
	// if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
	// 	return false, http.StatusUnsupportedMediaType
	// }
	//
	// if accept := r.Header.Get("Accept"); accept != "application/json" {
	// 	return false, http.StatusNotAcceptable
	// }

	return true, http.StatusOK
}

type getResponse struct {
	Books []book `json:"books"`
}

type book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublisherYear int    `json:"publisher_year"`
}
