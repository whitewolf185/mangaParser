package router

import (
	"github.com/whitewolf185/mangaparser/api/middleware"

	"github.com/go-chi/chi"
)

func newLoginRouter(errHandler middleware.ErrHandler) chi.Router {
	loginRouter := chi.NewRouter()
	return loginRouter
}