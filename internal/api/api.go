package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juliofilizzola/server/internal/store/pgstore"
)

type apiHandler struct {
	queries *pgstore.Queries
	r       *chi.Mux
}

func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		queries: q,
	}
	r := chi.NewRouter()

	a.r = r
	return a
}
