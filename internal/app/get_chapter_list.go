package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/whitewolf185/mangaparser/api/domain"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

// @Tags manga
// @Description Получения ссылок на главы манги. На вход подается ссылка на мангу с секцией chaptets. Обязательно ссылку надо прогнать через url encoder, например https://www.urlencoder.org/. На выходе должна получиться ссылка вида, https%3A%2F%2Fmangalib.me%2Fo-ju-yesuyeo%3Fsection%3Dchapters
// @ID manga-get-chapter-list
// @Accept json
// @Produce json
// @Param input query domain.GetChapterListRequest true "manga url with chapters section. Url must be escaped"
// @Success 200 {object} domain.GetChapterListResponse
// @Router /manga/GetChapterList [get]
func (i *Implementation) GetChapterList(ctx context.Context, req *domain.GetChapterListRequest) (*domain.GetChapterListResponse, error) {
	if req == nil {
		return nil, customerrors.CodesBadRequest(fmt.Errorf("empty request"))
	}
	
	mangaURL, err := unescapeUrl(req.MangaUrl)
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
