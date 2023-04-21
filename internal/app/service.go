package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
	
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

type (
	// MangaConfigurator интерфейс для управления информацией о манге
	MangaConfigurator interface {
		// GetChapterListUrlByMangaUrl получает главы манги и значение total
		GetChapterListUrlByMangaUrl(ctx context.Context, mainMangaUrl string) ([]string, error)
		// GetPicsUrlInChapter получает url картинок
		GetPicsUrlInChapter(ctx context.Context, chapterUrl string) ([]string, error)
		// Получает название манги по url главы
		GetMangaName(chapterUrl string) string
	}

	// интерфейс манипуляций с картинками
	ImageController interface {
		// GetImagesFromURLs скачивание картинок из url
		GetImagesFromURLs(ctx context.Context, folderPathToSave string, urlImages []string) error
	}
)

const regexpChapterCheckerPattern = `https://[a-z\.]+/[a-zA-Z\-\_]+/v[\d\.]/c[\d\.]+\?*`

// Implementation структура для реализации различных ручек
type Implementation struct {
	mangaConfigurator MangaConfigurator
	imageController ImageController

	chapterChecker *regexp.Regexp
}

// NewImplementation конструктор для Implementation
func NewImplementation(mangaConfigurator MangaConfigurator, imageController ImageController) (*Implementation, error) {
	r, err := regexp.Compile(regexpChapterCheckerPattern)
	if err != nil {
		return nil, errors.Wrap(err, "regexp compile failure")
	}
	return &Implementation{
		mangaConfigurator: mangaConfigurator,
		imageController: imageController,
		chapterChecker: r,
	}, nil
}

// getAndUnescapeStrFromUrlQuery функция собирает информацию из определенного query
func getAndUnescapeStrFromUrlQuery(req *http.Request, queryString string) (string, error) {
	resultUnescaped := req.URL.Query().Get(queryString)
	if resultUnescaped == "" {
		return "", customerrors.ErrUrlIsEmpty
	}

	result, err := url.QueryUnescape(resultUnescaped)
	if err != nil {
		return "", fmt.Errorf("unescape error: %w", err)
	}
	return result, nil
}
