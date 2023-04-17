package domain

import (
	"context"
	"net/http"
)

type HandlerType string

const (
	GetChapterList = HandlerType("GetChapterList")
)

// Интерфейс со всеми ручками сервиса
type Handlers interface {
	// GetChapterList Ручка получения ссылок на главы
	GetChapterList(ctx context.Context, req *http.Request) (*GetChapterListResponse, error)
}

// GetChapterListResponse структура ответа
type GetChapterListResponse struct {
	ChapterURLs []string
	Total       int
}
