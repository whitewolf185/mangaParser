package app

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/pkg/errors"

	"github.com/whitewolf185/mangaparser/api/domain"
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

	// интерфейс для отправки манги на электронную книгу
	EbookSender interface {
		SendManga(ctx context.Context, email string, mangaFilePath string) error
	}

	PersonRepo interface {
		GetEmailByID(ctx context.Context, person domain.PersonInfo) (string, error)
	}
)

const regexpChapterCheckerPattern = `https://[a-z\.]+/[a-zA-Z\-\_\d]+/v[\d\.]+/c[\d\.]+\?*`

// Implementation структура для реализации различных ручек
type Implementation struct {
	mangaConfigurator MangaConfigurator
	imageController ImageController
	ebookSender EbookSender
	personRepo PersonRepo

	chapterChecker *regexp.Regexp
}

// NewImplementation конструктор для Implementation
func NewImplementation(
	mangaConfigurator MangaConfigurator,
	imageController ImageController,
	ebookSender EbookSender,
	personRepo PersonRepo) (*Implementation, error) {
	r, err := regexp.Compile(regexpChapterCheckerPattern)
	if err != nil {
		return nil, errors.Wrap(err, "regexp compile failure")
	}
	return &Implementation{
		mangaConfigurator: mangaConfigurator,
		imageController: imageController,
		ebookSender: ebookSender,
		personRepo: personRepo,
		chapterChecker: r,
	}, nil
}

func unescapeUrl(escapedUrl string) (string, error) {
	result, err := url.QueryUnescape(escapedUrl)
	if err != nil {
		return "", fmt.Errorf("unescape error: %w", err)
	}
	return result, nil
}
