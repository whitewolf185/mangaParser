package pdfmerger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pkg/errors"
)

// CreatePDFFromImagesDir создает pdf файл.
// На вход подается путь к папке, где лежат картинки вида 1.png 2.jpg 3.tiff и путь к выходному pdf файлу
func CreatePDFFromImagesDir(imagesDirPath string, outputPath string) error {
	imagesPath, err := getImagesPathStr(imagesDirPath)
	if err != nil {
		return errors.Wrap(err, "something wrong with creating images path")
	}
	err = api.ImportImagesFile(imagesPath, outputPath, nil, nil)
	if err != nil {
		return errors.Wrap(err, "pdf creation filure")
	}
	return nil
}

func getImagesPathStr(imagesDirPath string) ([]string, error) {
	if imagesDirPath == "" {
		return nil, errors.Wrap(customerrors.ErrEmptyStr, "countFilesInDir")
	}
	workingDir, err := os.Open(imagesDirPath)
	if err != nil{
		return nil, errors.WithStack(err)
	}
	files, err := workingDir.ReadDir(0)
	switch {
	case err != nil:
		return nil, errors.WithStack(err)
	case len(files) == 0:
		return nil, errors.Wrap(customerrors.ErrEmptyDir, "countFilesInDir")
	}
	
	var result []string
	for _, file := range files {
		result = append(result, fmt.Sprintf("%s/%s", imagesDirPath, file.Name()))
	}

	sort.Slice(result, func(i, j int) bool {
		return prepareFileNameForSort(result[i]) < prepareFileNameForSort(result[j])
	})

	return result, nil
}

// prepareFileNameForSort функция используется для изменения имени файла так, чтобы можно было их сравнить "правильно".
// Например, обычно file1.ext < file12.ext < file2.ext, но с помощью функции теперь file1.ext < file2.ext < file12.ext.
// Функция паникует, если на вход приходит файл вида file12dtf.ext (после цифр есть буквы)
func prepareFileNameForSort(filename string) int {
    ext := filepath.Ext(filename)
    name := filename[:len(filename)-len(ext)]
    
	index := strings.IndexAny(name, "0123456789")
	result, err := strconv.Atoi(name[index:])
	if err != nil {
		panic(err)
	}
    
    return result
}