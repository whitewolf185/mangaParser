package mangalib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

// GetPicsUrlInChapter выдает список url с картинками из главы манги. Принимает url на главу манги
func (mlb mangaLibController) GetPicsUrlInChapter(ctx context.Context, chapterUrl string) ([]string, error) {
	res, err := http.Get(chapterUrl)
	if err != nil {
		return nil, errors.Wrap(err, "chapter http get failure")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.Wrap(err, "chapter http get failure: status is not 200")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "parsing into goquery failure")
	}

	parsedPageList, err := mlb.getPageInfo(doc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	chapterID, err := mlb.getChapterID(doc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mangaName := mlb.getMangaName(chapterUrl)

	result := make([]string, 0, len(parsedPageList))
	for _, page := range parsedPageList {
		result = append(result, fmt.Sprintf("https://img33.imgslib.link//manga/%s/chapters/%s/%s", mangaName, chapterID, page.Url))
	}
	return result, nil
}

func (mlb mangaLibController) getMangaName(chapterUrl string) string {
	result := strings.ReplaceAll(chapterUrl, "https://mangalib.me/", "")
	indexToSlice := strings.Index(result, "/")
	return result[:indexToSlice]
}

func (mlb mangaLibController) getPageInfo(doc *goquery.Document) ([]pageList, error) {
	pageInfo := doc.Find("#pg").Text()
	pageInfo = strings.ReplaceAll(pageInfo, "window.__pg = ", "")
	pageInfo = strings.ReplaceAll(pageInfo, "\n", "")
	pageInfo = strings.ReplaceAll(pageInfo, "\t", "")
	pageInfo = strings.ReplaceAll(pageInfo, " ", "")
	pageInfo = pageInfo[:len(pageInfo)-1]

	var parsedPageList []pageList
	err := json.Unmarshal([]byte(pageInfo), &parsedPageList)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal error")
	}

	sort.Slice(parsedPageList, func(i, j int) bool {
		return parsedPageList[i].Page < parsedPageList[j].Page
	})

	return parsedPageList, nil
}

func (mlb mangaLibController) getChapterID(doc *goquery.Document) (string, error) {
	result, ok := doc.Find("#comments").Attr("data-post-id")
	if !ok {
		return "", errors.Wrap(customerrors.ErrEmptyAttr, "mangalib: data-post-id")
	}
	return result, nil
}

