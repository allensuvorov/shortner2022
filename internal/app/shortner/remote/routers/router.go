package routers

import (
	"github.com/allensuvorov/urlshortner/internal/app/shortner/remote/handlers/url"
	"github.com/go-chi/chi/v5"
)

func NewRouter(url url.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{hash}", url.Get)
	r.Post("/", url.Create)
	r.Post("/api/shorten", url.Create)
	return r
}
