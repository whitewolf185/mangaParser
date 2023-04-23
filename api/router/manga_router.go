package router

import (

	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/api/middleware"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi"
)

func newMangaRouter(errHandler middleware.ErrHandler) chi.Router {
	mangaRouter := chi.NewRouter()
	mangaRouter.With(httpin.NewInput(domain.GetChapterListRequest{})).Get("/GetChapterList", 
		errHandler.ErrMiddleware(domain.GetChapterList),
	)
	mangaRouter.With(httpin.NewInput(domain.GetChapterPagesRequest{})).Get("/GetChapterPages", 
		errHandler.ErrMiddleware(domain.GetChapterPages),
	)
	mangaRouter.With(httpin.NewInput(domain.GetChapterPagesRequest{})).Get("/GetChapterPagesPDF",
		errHandler.ErrMiddleware(domain.GetChapterPagesPDF),
	)
	mangaRouter.Post("/SendToEbook",
		errHandler.ErrMiddleware(domain.SendToEbook),
	)
	return mangaRouter
}