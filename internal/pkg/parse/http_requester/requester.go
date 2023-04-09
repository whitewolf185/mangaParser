package httprequester

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/whitewolf185/mangaparser/internal/config"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

// сколько ретраев должен делать http request при 500 ошибке
const totalRetries = 2

// getResponse повторяет totalTry запрос, если ошибка != 200 и == 500
func getResponse(url string, timeToSleep time.Duration) (*http.Response, error) {
	for currentTry := 0; currentTry < totalRetries; currentTry++ {
		res, err := http.Get(url)
		if err != nil {
			return nil, errors.Wrap(err, "http get failure")
		}
		
		switch {
		case res.StatusCode >= 500 && res.StatusCode < 600:
			logrus.Infof("try %d", currentTry+1)
			time.Sleep(timeToSleep)
		case res.StatusCode == http.StatusOK:
			return res, nil
		}
	}

	return nil, errors.Wrapf(customerrors.ErrHttpRetry, "tries %d", totalRetries)
}

// GetDOM получает подготовленное DOM дерево при помощи получения сайта из приходящей ссылки
func GetDOM(url string) (*goquery.Document, error) {
	timeToSleep, err := time.ParseDuration(config.GetValue(config.RetryDuration))
	if err != nil {
		return nil, errors.Wrap(err, "retry duration parsing failure")
	}
	res, err := getResponse(url, timeToSleep)
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