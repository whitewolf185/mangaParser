package mangalib

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
	httprequester "github.com/whitewolf185/mangaparser/internal/pkg/parse/http_requester"
)

const imageDownloadTemplate = "https://img33.imgslib.link//manga/%s/chapters/%s/%s"

// GetPicsUrlInChapter выдает список url с картинками из главы манги. Принимает url на главу манги
func (mlb mangaLibController) GetPicsUrlInChapter(ctx context.Context, chapterUrl string) ([]string, error) {
	doc, err := httprequester.GetDOM(chapterUrl)
	if err != nil {
		return nil, errors.Wrap(err, "chapter pages")
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
		result = append(result, fmt.Sprintf(imageDownloadTemplate, mangaName, chapterID, page.Url))
	}
	return result, nil
}

func (mlb mangaLibController) getPageInfo(doc *goquery.Document) ([]pageList, error) {
	pageInfo := doc.Find("#pg").Text()
	pageInfo = strings.ReplaceAll(pageInfo, "window.__pg = ", "")
	pageInfo = mlb.cleanScript(pageInfo)

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

