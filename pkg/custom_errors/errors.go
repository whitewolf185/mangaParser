package customerrors

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	// ErrEmptyDir ошибка сигналисирует о том, что в папке нет файлов
	ErrEmptyDir = errors.New("empty directory")
	// ErrEmptyStr ошибка сигнализирует о том, что пришеджая в функцию строка пустая
	ErrEmptyStr = errors.New("empty string input")
	// ErrEmptyImage ошибка сигнализирует о том, что полученная картинка оказалось пустой.
	ErrEmptyImage = errors.New("no image detected")

	// ErrEmailsNotFound ошибка сигнализирует о том, что не было найдено email по id пользователя
	ErrEmailsNotFound = errors.New("emails not found")
	// ErrUrlIsEmpty ошибка сигнализирует о том, что не было найдено url для манги
	ErrUrlIsEmpty = errors.New("no url for manga found")
	// ErrWrongUrl ошибка сигнализирует о том, что ссылка не соответствует регулярным выражениям
	ErrWrongUrl = errors.New("url does not watch reg exp")

	// ErrEmptyAttr ошибка говорит о том, что искомы атрибут на сайте не существует
	ErrEmptyAttr = errors.New("attribute do not exists")

	// ErrUnknownType ошбика сигнализирует о том, что был пойман неизвестный тип
	ErrUnknownType = errors.New("unkown type")

	// ErrHttpRetry ошибка означает, что после некоторого количества ретраев достучаться до сайта не получилось
	ErrHttpRetry = errors.New("http request retry failure")
)

type ErrCodes struct {
	Err  error
	Code int
}

func (e ErrCodes) Error() string {
	return e.Err.Error()
}
func (e ErrCodes) StatusCode() int {
	return e.Code
}

func CodesNotFound(err error) error {
	return ErrCodes{
		Err:  err,
		Code: http.StatusNotFound,
	}
}

func CodesBadRequest(err error) error {
	return ErrCodes{
		Err:  err,
		Code: http.StatusBadRequest,
	}
}
