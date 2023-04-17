package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/whitewolf185/mangaparser/api/domain"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

const contextDeadline = time.Minute * 5

// ErrHandler структура используется для того, чтобы возможно было хендлить ошибки из ручек
type ErrHandler struct {
	mangaHandler domain.Handlers
}

// Конструктор для ErrHandler
func NewErrorHandler(mangaHandler domain.Handlers) ErrHandler {
	return ErrHandler{
		mangaHandler: mangaHandler,
	}
}

// ErrMiddleware - функция-хендлер. Принимает в себя тип ручки, которая используется в хендлере
func (em ErrHandler) ErrMiddleware(handleType domain.HandlerType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Method", string(handleType))
		logrus.Infof("method: %v", handleType)
		ctx, cancel := context.WithTimeout(context.Background(), contextDeadline)
		defer cancel()

		res, err := em.handleTypeSwitcher(ctx, r, handleType)

		if err != nil {
			logrus.Errorln(err.Error())
			w.WriteHeader(err.(customerrors.ErrCodes).Code)
			return
		}
		if res == nil {
			w.Write(nil)
			w.WriteHeader(http.StatusOK)
		}
		toSend, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Errorln("unmarshal response error ", err.Error())
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(toSend)
		w.WriteHeader(http.StatusOK)
	}
}

func (em ErrHandler) handleTypeSwitcher(ctx context.Context, r *http.Request, handleType domain.HandlerType) (interface{}, error) {
	switch handleType {
	case domain.GetChapterList:
		return em.mangaHandler.GetChapterList(ctx, r)
	}
	return nil, customerrors.ErrUnknownType
}
