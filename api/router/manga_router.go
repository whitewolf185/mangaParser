package router

import (

	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/api/middleware"

	"github.com/go-chi/chi"
)

func newMangaRouter(errHandler middleware.ErrHandler) chi.Router {
	mangaRouter := chi.NewRouter()
	mangaRouter.Get("/GetChapterList", 
		errHandler.ErrMiddleware(domain.GetChapterList),
	)
	return mangaRouter
}