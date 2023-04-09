package mangalib

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	mangalibHost = "https://mangalib.me/"
	regexpCompilePattern = `[\s]`
)

// структура управления получения ссылок на картинки из mangalib
type mangaLibController struct {
	urlGetter UrlGetter
	regExpScriptCleaner *regexp.Regexp
}

// NewMangaLibController конструктор mangaLibController
func NewMangaLibController(urlGetter UrlGetter) (mangaLibController, error) {
	r, err := regexp.Compile(regexpCompilePattern)
	if err != nil {
		return mangaLibController{}, errors.Wrap(err, "regexp compile failure")
	}
	return mangaLibController{
		urlGetter: urlGetter,
		regExpScriptCleaner: r,
	}, nil
}

// структура, которая используется для парсинга данных об имеющихся картинках в главе для получения ссылок на них
type pageList struct {
	Page int `json:"p"`
	Url string `json:"u"`
}

func (mlc mangaLibController) cleanScript(scriptStr string) string {
	tmpStr := mlc.regExpScriptCleaner.ReplaceAllString(scriptStr, "")
	return tmpStr[:len(tmpStr)-1]
}

func (mlс mangaLibController) getMangaName(chapterUrl string) string {
	result := strings.ReplaceAll(chapterUrl, mangalibHost, "")
	indexToSlice := strings.Index(result, "/")
	if indexToSlice == -1 {
		indexToSlice = strings.Index(result, "?")
	}
	return result[:indexToSlice]
}