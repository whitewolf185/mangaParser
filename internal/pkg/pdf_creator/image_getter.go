package pdf_creator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/whitewolf185/mangaparser/internal/config"
	"github.com/whitewolf185/mangaparser/internal/pkg/err_controller"
	"github.com/whitewolf185/mangaparser/internal/pkg/utils"
)

type imageGetter struct{}

// NewImageGetter конструктор imageGetter
func NewImageGetter() imageGetter {
	return imageGetter{}
}

const formatToSaveFile = "%s/img%d.%s"

// GetImageAndSave достает картинку по ссылке и пишет ее в канал в асинхронном режиме
func (ig imageGetter) GetImageAndSave(
	ctx context.Context, wg *sync.WaitGroup, ec *err_controller.ErrController,
	folderPathToSave string, page int, url string) {
	wg.Add(1)
	if ctxErr := ctx.Err(); ctxErr != nil { // проверяем контекст на истечение дедлайна
		ec.PutError(ctxErr)
		wg.Done()
		return
	}
	go func() {
		defer wg.Done()
		res, err := http.Get(url)
		if err != nil {
			ec.PutError(err)
			return
		}
		defer res.Body.Close()
		if ctxErr := ctx.Err(); ctxErr != nil { // проверяем контекст на истечение дедлайна
			ec.PutError(ctxErr)
			return
		}
		extension, err := ig.extensionParser(url)
		if err != nil {
			ec.PutError(err)
			return
		}
		pathToSave := fmt.Sprintf(formatToSaveFile, folderPathToSave, page, extension)
		file, err := os.Create(pathToSave) // открываем файл
		if err != nil {
			ec.PutError(fmt.Errorf("cannot open file %d\nError: %w", page, err))
			return
		}
		defer file.Close()
		if ctxErr := ctx.Err(); ctxErr != nil { // проверяем контекст на истечение дедлайна
			ec.PutError(ctxErr)
			return
		}
		_, err = io.Copy(file, res.Body) // копируем файл
		if err != nil {
			ec.PutError(fmt.Errorf("cannot save file %d\nError: %w", page, err))
			return
		}
	}()
}

func (ig imageGetter) extensionParser(url string) (string, error) {
	result := filepath.Ext(url)
	if result == "" || !utils.StringMatching(result, config.SliceOfFileExtensions...) {
		return "", fmt.Errorf("this is not an image")
	}
	return result, nil
}
