package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RunServer(h *Handlers) error {
	r := httprouter.New()
	r.Handle("POST", "/employee", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.PutEmployee(r.Context(), w, r)
	})

	return http.ListenAndServe(":8080", r)
}
