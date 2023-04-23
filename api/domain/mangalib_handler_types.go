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
	SendToEbook        = HandlerType("SendToEbook")
)

// Интерфейс со всеми ручками сервиса
type Handlers interface {
	// GetChapterList Ручка получения ссылок на главы
	GetChapterList(ctx context.Context, req *GetChapterListRequest) (*GetChapterListResponse, error)
	// GetChapterPages ручка получения самих страниц манги определенной главы
	GetChapterPages(ctx context.Context, req *GetChapterPagesRequest) (*GetChapterPagesResponse, error)
	// GetChapterPagesPDF ручка получения страниц главы манги, объединенные в pdf
	GetChapterPagesPDF(ctx context.Context, req *GetChapterPagesRequest) (*GetChapterPagesPDFResponse, error)
	// SendToEbook отправляет pdf манги на электронную книгу. POST запрос, требующий в своем теле SendToEbookRequest структуру
	SendToEbook(ctx context.Context, req *http.Request) (*SendToEbookResponse, error)
}

type GetChapterListRequest struct {
	MangaUrl string `json:"manga_url" in:"query=manga_url"`
}
// GetChapterListResponse структура ответа
type GetChapterListResponse struct {
	ChapterURLs []string `json:"chapter_urls"`
	MangaName   string   `json:"manga_name"`
	Total       int      `json:"total"`
}

// ImageBody структура, которая нужна для того, чтобы хранить картинку и название файла, которая эта картинка была записана
type ImageBody struct {
	Images   []byte `json:"images"`
	FileName string `json:"file_name"`
}
type GetChapterPagesRequest struct {
	ChapterUrl string `json:"chapter_url" in:"query=chapter_url"`
	PersonID string `json:"person_id" in:"query=person_id"`
}
type GetChapterPagesResponse struct {
	Pages     []ImageBody `json:"pages"`
	Total     int `json:"total"`
	MangaName string `json:"manga_name"`
}


type GetChapterPagesPDFResponse struct {
	PdfPath string
}

type PersonInfo struct {
	PersonID string `json:"person_id"`
	TelegramID int64 `json:"telegram_id"`
}
type SendToEbookRequest struct {
	ChapterUrl string    `json:"chapter_url"`
	Person   PersonInfo `json:"person"`
}
type SendToEbookResponse struct {
	MangaName  string `json:"manga_name"`
	TotalPages int `json:"total_pages"`
}
