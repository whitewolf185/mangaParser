package customerrors

import "github.com/pkg/errors"

var (
	// ErrEmptyDir ошибка сигналисирует о том, что в папке нет файлов
	ErrEmptyDir = errors.New("empty directory")
	// ErrEmptyStr ошибка сигнализирует о том, что пришеджая в функцию строка пустая
	ErrEmptyStr = errors.New("empty string input")
)