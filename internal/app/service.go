package app

import (
	"context"
)

type (
	MangaConfigurator interface {
		GetChapterListUrlByMangaUrl(ctx context.Context, mainMangaUrl string) ([]string, error)
	}
)

// Implementation структура для реализации различных ручек
type Implementation struct {
	mangaConfigurator MangaConfigurator
}

// NewImplementation конструктор для Implementation
func NewImplementation(mangaConfigurator MangaConfigurator) *Implementation {
	return &Implementation{
		mangaConfigurator: mangaConfigurator,
	}
}
