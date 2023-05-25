package main

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/whitewolf185/mangaparser/api/middleware"
	"github.com/whitewolf185/mangaparser/api/router"
	"github.com/whitewolf185/mangaparser/internal/app"
	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/config/flags"
	"github.com/whitewolf185/mangaparser/internal/pkg/mailer"
	"github.com/whitewolf185/mangaparser/internal/pkg/parse/mangalib"
	"github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator"
	"github.com/whitewolf185/mangaparser/internal/repository"
)

// @title Manga parser
// @version @0.9
// @description swagger для api к мангапарсеру

// @host 95.165.166.169
// @basePath /api

func main() {
	ctx := context.Background()
	flags.InitServiceFlags()
	db, err := config.ConnectPostgres(ctx)
	if err != nil {
		logrus.Fatalln("cannot connect to postgresql ", err)
	}

	// подготавливаем разные структуры
	urlGetter := repository.NewUrlRepo(db)
	mangalibController, err := mangalib.NewMangaLibController(urlGetter)
	if err != nil {
		logrus.Fatalln("cannot configure mangalib controller ", err)
	}

	// библиотеки для скачивания картинок
	imageGetter := pdf_creator.NewImageGetter()
	imageController := pdf_creator.NewImageController(imageGetter)

	ebookSender := mailer.NewEbookMailer()

	personRepo := repository.NewPersonController(db)
	
	// Имплементация API
	application, err := app.NewImplementation(mangalibController, imageController, ebookSender, personRepo)
	if err != nil {
		logrus.Fatalln("cannot configure implementation")
	}

	root := router.NewRouter(middleware.NewErrorHandler(application))

	logrus.Info("app successfully started")
	http.ListenAndServe(":"+config.GetValue(config.ListenPort), root)
}
