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

type retrier struct {
	currentTry int
}
func newRetrier () retrier {
	return retrier{
		currentTry: 0,
	}
}

// retryHTTPGet повторяет totalTry запрос, если ошибка != 200 и == 500
func (r retrier) retryHTTPGet(url string, timeToSleep time.Duration, totalTry int) (*http.Response, error) {
	if r.currentTry >= totalTry {
		return nil, errors.Wrap(customerrors.ErrHttpRetry, fmt.Sprintf("retring failure tries %d", totalTry))
	}
	if r.currentTry != 0{
		logrus.Infof("try %d", r.currentTry)
		time.Sleep(timeToSleep)
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "http get failure")
	}
	switch {
	case res.StatusCode >= 500 && res.StatusCode < 600:
		r.currentTry++
		return r.retryHTTPGet(url, timeToSleep, totalTry)
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
	retrier := newRetrier()
	res, err := retrier.retryHTTPGet(url, timeToSleep, 2)
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