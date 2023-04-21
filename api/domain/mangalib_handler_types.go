package domain

import (
	"context"
	"net/http"
)

type HandlerType string

const (
	GetChapterList     = HandlerType("GetChapterList")
	GetChapterPages    = HandlerType("GetChapterPages")
	GetChapterPagesPDF = HandlerType("GetChapterPagesPDF")
)

// Интерфейс со всеми ручками сервиса
type Handlers interface {
	// GetChapterList Ручка получения ссылок на главы
	GetChapterList(ctx context.Context, req *http.Request) (*GetChapterListResponse, error)

	// GetChapterPages ручка получения самих страниц манги определенной главы
	GetChapterPages(ctx context.Context, req *http.Request) (*GetChapterPagesResponse, error)

	// GetChapterPagesPDF ручка получения страниц главы манги, объединенные в pdf
	GetChapterPagesPDF(ctx context.Context, req *http.Request) (*GetChapterPagesPDFResponse, error)
}

// GetChapterListResponse структура ответа
type GetChapterListResponse struct {
	ChapterURLs []string
	MangaName   string
	Total       int
}

// ImageBody структура, которая нужна для того, чтобы хранить картинку и название файла, которая эта картинка была записана
type ImageBody struct {
	Image    []byte
	FileName string
}
type GetChapterPagesResponse struct {
	Pages     []ImageBody
	Total     int
	MangaName string
}

type GetChapterPagesPDFResponse struct {
	PdfPath string
}
