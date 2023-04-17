package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/whitewolf185/mangaparser/api/middleware"
	"github.com/whitewolf185/mangaparser/api/router"
	"github.com/whitewolf185/mangaparser/internal/app"
	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/config/flags"
	"github.com/whitewolf185/mangaparser/internal/pkg/parse/mangalib"
	"github.com/whitewolf185/mangaparser/internal/repository"
)

func main() {
	ctx := context.Background()
	flags.InitServiceFlags()
	db, err := config.ConnectPostgres(ctx)
	if err != nil {
		logrus.Fatalln("cannot connect to postgresql ", err)
	}

	urlGetter := repository.NewUrlRepo(db)
	mangalibController, err := mangalib.NewMangaLibController(urlGetter)
	if err != nil {
		logrus.Fatalln("cannot configure mangalib controller ", err)
	}

	application := app.NewImplementation(mangalibController)

	root := router.NewRouter(middleware.NewErrorHandler(application))

	fmt.Println("app successfully started")
	http.ListenAndServe(":"+config.GetValue(config.ListenPort), root)
}
