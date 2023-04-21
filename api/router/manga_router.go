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
	mangaRouter.Get("/GetChapterPages", 
		errHandler.ErrMiddleware(domain.GetChapterPages),
	)
	mangaRouter.Get("/GetChapterPagesPDF",
		errHandler.ErrMiddleware(domain.GetChapterPagesPDF),
	)
	return mangaRouter
}