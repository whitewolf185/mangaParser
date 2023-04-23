package router

import (
	"fmt"
	"time"

	custommiddleware "github.com/whitewolf185/mangaparser/api/middleware"
	_ "github.com/whitewolf185/mangaparser/docs"
	"github.com/whitewolf185/mangaparser/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/swaggo/http-swagger/v2"
)

const (
	defaultTimeout = 2 * time.Minute
	swaggerUrlPattern = "http://%s:%s/docs/doc.json"
)

func NewRouter(
	errHandler custommiddleware.ErrHandler,
) chi.Router {
	root := chi.NewRouter()

	// useful log info
	root.Use(middleware.RequestID)
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)

	// автоматический таймаут для всех запросов
	root.Use(middleware.Timeout(defaultTimeout))

	// обязательный mount на все. Нужен для создания запросов на /api
	// Например:
	// root.Mount("/api", newCusmomRouter)
	root.Mount("/manga", newMangaRouter(errHandler))
	root.Mount("/login", newLoginRouter(errHandler))

	rootApi := chi.NewRouter()
	rootApi.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf(swaggerUrlPattern, config.GetValue(config.ServerIP), config.GetValue(config.ListenPort))),
	))
	rootApi.Get("/docs-local/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf(swaggerUrlPattern, "localhost", config.GetValue(config.ListenPort))),
	))

	rootApi.Mount("/api", root)
	return rootApi
}
