package mangalib

import (
	"context"

	"github.com/google/uuid"

	"github.com/whitewolf185/mangaparser/internal/config"
)

// UrlGetter интерфейс для получения url для парсинга
//
//go:generate mockgen -destination=./mock/url_getter_mock.go -package=mock github.com/whitewolf185/mangaparser/internal/pkg/parse/mangalib UrlGetter
type UrlGetter interface {
	GetUrlByID(ctx context.Context, mangaID uuid.UUID, sourceType config.MangaSourceType) (string, error)
}