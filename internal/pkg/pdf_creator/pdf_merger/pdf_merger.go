package pdfmerger

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"

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
	files, err := ioutil.ReadDir(imagesDirPath)
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
		return sortName(result[i]) < sortName(result[j])
	})

	return result, nil
}

// sortName returns a filename sort key with
// non-negative integer suffixes in numeric order.
// For example, amt, amt0, amt2, amt10, amt099, amt100, ...
func sortName(filename string) string {
    ext := filepath.Ext(filename)
    name := filename[:len(filename)-len(ext)]
    // split numeric suffix
    i := len(name) - 1
    for ; i >= 0; i-- {
        if '0' > name[i] || name[i] > '9' {
            break
        }
    }
    i++
    // string numeric suffix to uint64 bytes
    // empty string is zero, so integers are plus one
    b64 := make([]byte, 64/8)
    s64 := name[i:]
    if len(s64) > 0 {
        u64, err := strconv.ParseUint(s64, 10, 64)
        if err == nil {
            binary.BigEndian.PutUint64(b64, u64+1)
        }
    }
    // prefix + numeric-suffix + ext
    return name[:i] + string(b64) + ext
}