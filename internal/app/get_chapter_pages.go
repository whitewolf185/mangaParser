package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/internal/config"
	pdfmerger "github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator/pdf_merger"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

// @Tags manga
// @Description Получения массива байт картинок с метаданными
// @ID manga-get-chapter-pages
// @Accept json
// @Produce json
// @Param input query domain.GetChapterPagesRequest true "chapter url and person id info. Chapter must be url encoded"
// @Success 200 {object} domain.GetChapterPagesResponse
// @Router /manga/GetChapterPages [get]
func (i *Implementation) GetChapterPages(ctx context.Context, req *domain.GetChapterPagesRequest) (*domain.GetChapterPagesResponse, error) {
	if req == nil {
		return nil, customerrors.CodesBadRequest(fmt.Errorf("empty request"))
	}
	// unescaping url
	chapterUrl, err := unescapeUrl(req.ChapterUrl) 
	if err != nil {
		return nil, customerrors.CodesBadRequest(err)
	}

	// Валидируем пришедшие данные
	if !i.chapterChecker.Match([]byte(chapterUrl)) {
		return nil, customerrors.CodesBadRequest(fmt.Errorf("wrong chapter url"))
	}

	if req.PersonID == "" {
		return nil, customerrors.CodesBadRequest(fmt.Errorf("empty person ID"))
	}

	// Начинаем заполнять результат
	var result domain.GetChapterPagesResponse
	// получим название манги
	result.MangaName = i.mangaConfigurator.GetMangaName(chapterUrl)

	pageUrls, err := i.mangaConfigurator.GetPicsUrlInChapter(ctx, chapterUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get pages url for manga %s", result.MangaName)
	}

	preparedPathToDownload := fmt.Sprintf(config.ParentPathToDownloadPattern, req.PersonID, result.MangaName)

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
		if err != nil || len(img) == 0{
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
