package err_controller

import "sync"

// ErrController структура по управлению ошибок при асинхронной работе
type ErrController struct {
	mx     sync.Mutex
	errors []error
}

// NewErrController конструктор ErrController
func NewErrController() *ErrController {
	return &ErrController{mx: sync.Mutex{}}
}

// IsNul проверяет, что нет ошибок.
// Если ошибки есть, то возвращает эти ошибки и, ВНИМАНИЕ, обнуляет массив ошибок.
// После повторного вызова IsNul будет всегда возвращаться nil.
func (e *ErrController) IsNul() []error {
	e.mx.Lock()
	defer e.mx.Unlock()

	if e.errors == nil {
		return nil
	}
	result := make([]error, len(e.errors))
	copy(result, e.errors)
	e.errors = nil
	return result
}

// PutError кладет ошибки в стэк ошибок для дальнейшей их обработки
func (e *ErrController) PutError(err error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.errors = append(e.errors, err)
}
