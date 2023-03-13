package pdf_creator

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/pkg/err_controller"
)

//go:generate mockgen -destination=./mock/image_getter_mock.go -package=mock github.com/whitewolf185/mangaparser/internal/pkg/pdf_creator ImageGetter
type (
	// ImageGetter - интерфейс получения картинки по url
	ImageGetter interface {
		// GetImageAndSave достает картинку по ссылке и пишет ее в канал в асинхронном режиме
		GetImageAndSave(ctx context.Context, wg *sync.WaitGroup, ec *err_controller.ErrController, folderPathToSave string, page int, url string)
	}
)

type imageController struct {
	imageGetter   ImageGetter
	errController *err_controller.ErrController
}

// NewImageController - конструктор для imageController
func NewImageController(ig ImageGetter) imageController {
	return imageController{imageGetter: ig, errController: err_controller.NewErrController()}
}

// GetImagesFromURLs запускает процесс сохранения картинок. Картинки должны быть отсортированы в нужном порядке.
func (ic imageController) GetImagesFromURLs(ctx context.Context, folderPathToSave string, urlImages []string) error {
	ctx, cancel := context.WithTimeout(ctx, config.TimeOutImageDownloading)
	defer cancel()
	err := ic.folderController(folderPathToSave)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for page, image := range urlImages {
		ic.imageGetter.GetImageAndSave(ctx, &wg, ic.errController, folderPathToSave, page, image)
	}
	wg.Wait()
	if erResult := ic.errController.IsNul(); erResult != nil {
		// TODO сделать здесь возврат своей ошибки
		return fmt.Errorf("something wrong with download images \nErrors list:\n%v", erResult)
	}
	return nil
}

func (ic imageController) folderController(folderPathToSave string) error {
	err := os.MkdirAll(folderPathToSave, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "unexpected error in folderController")
	}
	return nil
}
