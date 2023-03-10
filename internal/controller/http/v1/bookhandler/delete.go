package bookhandler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	const handler = _delete

	if isValid, code := h.validateDelete(r); !isValid {
		http.Error(w, "validation error", code)
		return
	}

	rawReq, err := io.ReadAll(r.Body)
	if err != nil {
		h.l.Error("reading request books on %s: %v", handler, err)
		http.Error(w, "reading request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req deleteRequest
	if err = json.Unmarshal(rawReq, &req); err != nil {
		h.l.Error("marshaling response on %s: %v", handler, err)
		http.Error(w, "unmarshaling: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.book.Delete(r.Context(), req.BookId); err != nil {
		h.l.Error("deleting book on %s: %v", handler, err)
		http.Error(w, "deleting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) validateDelete(r *http.Request) (bool, int) {
	if r.Method != http.MethodDelete {
		return false, http.StatusMethodNotAllowed
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		return false, http.StatusUnsupportedMediaType
	}

	if accept := r.Header.Get("Accept"); accept != "application/json" {
		return false, http.StatusNotAcceptable
	}

	return true, http.StatusOK
}

type deleteRequest struct {
	BookId int `json:"book_Id"`
}
