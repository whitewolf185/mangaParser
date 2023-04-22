package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/whitewolf185/mangaparser/api/domain"
	"github.com/whitewolf185/mangaparser/internal/config"
	pdfmerger "github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator/pdf_merger"
	"github.com/whitewolf185/mangaparser/internal/pkg/utils"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

func (i *Implementation) GetChapterPagesPDF(ctx context.Context, req *http.Request) (*domain.GetChapterPagesPDFResponse, error) {
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
	// получим название манги
	mangaName := i.mangaConfigurator.GetMangaName(chapterUrl)

	pageUrls, err := i.mangaConfigurator.GetPicsUrlInChapter(ctx, chapterUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get pages url for manga %s", mangaName)
	}

	preparedPathToDownload := fmt.Sprintf(config.ParentPathToDownloadPattern, personID, mangaName)

	err = i.imageController.GetImagesFromURLs(ctx, preparedPathToDownload, pageUrls)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot download manga for %s", mangaName)
	}

	// создаем папку, если она еще не была создана
	pdfFilePathFolder := fmt.Sprintf(config.ParentPathPdfFolderPattern, personID)
	err = utils.FolderController(pdfFilePathFolder)
	if err != nil{
		return nil, err
	}
	// формируем путь к файлу так, чтобы файл назывался названием манги
	pdfFilePathFolder += "/"+mangaName+".pdf"
	err = pdfmerger.CreatePDFFromImagesDir(preparedPathToDownload, pdfFilePathFolder)
	if err != nil {
		return nil, errors.Wrap(err, "pdf generate failed")
	}
	
	// Почистим папку, в которую скачали мангу по частям
	os.RemoveAll(preparedPathToDownload)

	return &domain.GetChapterPagesPDFResponse{
		PdfPath: pdfFilePathFolder,
	}, nil
}