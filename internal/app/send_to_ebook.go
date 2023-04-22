package app

import (
	"context"
	"encoding/json"
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

func (i *Implementation) SendToEbook(ctx context.Context, req *http.Request) (*domain.SendToEbookResponse, error) {
	// Валидируем пришедшие данные
	defer req.Body.Close()

	var requestBody domain.SendToEbookRequest
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		return nil, customerrors.CodesBadRequest(err)
	}

	email, err := i.personRepo.GetEmailByID(ctx, requestBody.Person)
	switch {
	case errors.Is(err, customerrors.ErrEmailsNotFound):
		return nil, customerrors.CodesNotFound(fmt.Errorf("person not exists or email is empty"))
	case err != nil:
		return nil, err
	}
	
	// Начинаем заполнять результат 
	// получим название манги
	mangaName := i.mangaConfigurator.GetMangaName(requestBody.ChapterUrl)

	pageUrls, err := i.mangaConfigurator.GetPicsUrlInChapter(ctx, requestBody.ChapterUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get pages url for manga %s", mangaName)
	}

	preparedPathToDownload := fmt.Sprintf(config.ParentPathToDownloadPattern, email, mangaName)

	err = i.imageController.GetImagesFromURLs(ctx, preparedPathToDownload, pageUrls)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot download manga for %s", mangaName)
	}
	// Почистим папку, в которую скачали мангу по частям
	defer os.RemoveAll(preparedPathToDownload)

	// создаем папку, если она еще не была создана
	pdfFilePathFolder := fmt.Sprintf(config.ParentPathPdfFolderPattern, email)
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
	defer os.RemoveAll(pdfFilePathFolder)
	
	// отправляем мангу на электронную книку
	i.ebookSender.SendManga(ctx, email, pdfFilePathFolder)

	return &domain.SendToEbookResponse{
		MangaName: mangaName,
		TotalPages: len(pageUrls),
	}, nil
}