package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/whitewolf185/mangaparser/api/domain"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

func (i *Implementation) GetChapterList(ctx context.Context, req *http.Request) (*domain.GetChapterListResponse, error) {
	mangaURLEscaped := req.URL.Query().Get("mangaURL")
	if mangaURLEscaped == "" {
		return nil, customerrors.CodesBadRequest(customerrors.ErrUrlIsEmpty)
	}

	mangaURL, err := url.QueryUnescape(mangaURLEscaped)
	if err != nil {
		return nil, fmt.Errorf("unescape error: %w", err)
	}

	if !strings.Contains(mangaURL, "chapters") {
		return nil, customerrors.CodesBadRequest(errors.Wrap(customerrors.ErrUrlIsEmpty, "this is not an chapters section"))
	}

	urls, err := i.mangaConfigurator.GetChapterListUrlByMangaUrl(ctx, mangaURL)
	if err != nil {
		return nil, err
	}

	result := domain.GetChapterListResponse{
		ChapterURLs: urls,
		Total:       len(urls),
	}

	return &result, nil
}
