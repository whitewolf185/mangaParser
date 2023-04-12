package customerrors

import "github.com/pkg/errors"

var (
	// ErrEmptyDir ошибка сигналисирует о том, что в папке нет файлов
	ErrEmptyDir = errors.New("empty directory")
	// ErrEmptyStr ошибка сигнализирует о том, что пришеджая в функцию строка пустая
	ErrEmptyStr = errors.New("empty string input")

	// ErrEmailsNotFound ошибка сигнализирует о том, что не было найдено email по id пользователя
	ErrEmailsNotFound = errors.New("emails not found")
	// ErrUrlNotFound ошибка сигнализирует о том, что не было найдено url для манги
	ErrUrlNotFound = errors.New("no url for manga found")

	// ErrEmptyAttr ошибка говорит о том, что искомы атрибут на сайте не существует
	ErrEmptyAttr = errors.New("attribute do not exists")

	// ErrHttpRetry ошибка означает, что после некоторого количества ретраев достучаться до сайта не получилось
	ErrHttpRetry = errors.New("http request retry failure")
)
