package httprequester

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/whitewolf185/mangaparser/internal/config"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

// retryHTTPGet повторяет totalTry запрос, если ошибка != 200 и == 500
func retryHTTPGet(url string, timeToSleep time.Duration, currentTry, totalTry int) (*http.Response, error) {
	if currentTry >= totalTry {
		return nil, errors.Wrap(customerrors.ErrHttpRetry, fmt.Sprintf("retring failure tries %d", totalTry))
	}
	if currentTry != 0{
		logrus.Infof("try %d", currentTry)
		time.Sleep(timeToSleep)
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "http get failure")
	}
	switch {
	case res.StatusCode >= 500 && res.StatusCode < 600:
		return retryHTTPGet(url, timeToSleep, currentTry+1, totalTry)
	case res.StatusCode != 200:
		return nil, errors.Wrapf(err, "chapter list http get failure: status is %d", res.StatusCode)
	}

	return res, nil
}

// GetDOM получает подготовленное DOM дерево при помощи получения сайта из приходящей ссылки
func GetDOM(url string) (*goquery.Document, error) {
	timeToSleep, err := time.ParseDuration(config.GetValue(config.RetryDuration))
	if err != nil {
		return nil, errors.Wrap(err, "retry duration parsing failure")
	}
	res, err := retryHTTPGet(url, timeToSleep, 0, 2)
	if err != nil{
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "parsing into goquery failure")
	}

	return doc, nil
}