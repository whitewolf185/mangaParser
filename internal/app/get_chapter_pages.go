package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/internal/config"
	pdfmerger "github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator/pdf_merger"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

const (
	chapterUrlQuery = "chapterURL"
	personIDQuery   = "personID"
)

func (i *Implementation) GetChapterPages(ctx context.Context, req *http.Request) (*domain.GetChapterPagesResponse, error) {
	// Валидируем пришедшие данные
	chapterUrl, err := getAndUnescapeStrFromUrlQuery(req, chapterUrlQuery)
	if err != nil {
		return nil, customerrors.CodesBadRequest(err)
	}

	if !i.chapterChecker.Match([]byte(chapterUrl)) {
		return nil, customerrors.CodesBadRequest(fmt.Errorf("wrong chapter url"))
	}

	personID, err := getAndUnescapeStrFromUrlQuery(req, personIDQuery)
	if err != nil {
		return nil, customerrors.CodesBadRequest(err)
	}

	// Начинаем заполнять результат
	var result domain.GetChapterPagesResponse
	// получим название манги
	result.MangaName = i.mangaConfigurator.GetMangaName(chapterUrl)

	pageUrls, err := i.mangaConfigurator.GetPicsUrlInChapter(ctx, chapterUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get pages url for manga %s", result.MangaName)
	}

	preparedPathToDownload := fmt.Sprintf(config.ParentPathToDownloadPattern, personID, result.MangaName)

	err = i.imageController.GetImagesFromURLs(ctx, preparedPathToDownload, pageUrls)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot download manga for %s", result.MangaName)
	}

	pathesToImgs, err := pdfmerger.GetImagesPathStr(preparedPathToDownload)
	if err != nil {
		return nil, errors.Wrap(err, "get images pathes failed")
	}
	// Подсчитаем колпчество страниц, которые получилось скачать
	result.Total = len(pathesToImgs)

	// запишем байтовое представление файла в итоговый массив
	result.Pages = make([]domain.ImageBody, 0, result.Total)
	for iteration, path := range pathesToImgs {
		img, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "get bytes failed on %d iteration", iteration)
		}
		result.Pages = append(result.Pages, domain.ImageBody{
			Images:   img,
			FileName: filepath.Base(path),
		})
	}

	// Почистим папку, в которую скачали мангу по частям
	os.RemoveAll(preparedPathToDownload)

	return &result, nil
}
