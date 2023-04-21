package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/whitewolf185/mangaparser/api/domain"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

const (
	mangaUrlQuery = "mangaURL"
)

func (i *Implementation) GetChapterList(ctx context.Context, req *http.Request) (*domain.GetChapterListResponse, error) {
	mangaURL, err := getAndUnescapeStrFromUrlQuery(req, mangaUrlQuery)
	if err != nil {
		return nil, customerrors.CodesBadRequest(err)
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
		MangaName: i.mangaConfigurator.GetMangaName(mangaURL),
	}

	return &result, nil
}
