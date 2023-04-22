package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
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
			switch v := err.(type) {
			case customerrors.ErrCodes:
				w.WriteHeader(v.Code)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			err = errors.Wrapf(err, "Method %s ->", string(handleType))
			logrus.Errorln(err.Error())

			return
		}
		em.checkExclusiveFiles(res, w, handleType)
	}
}

// checkExclusiveFiles функция делает специфическую отправку файлов (не json), если того требует логика
func (em ErrHandler) checkExclusiveFiles(res interface{}, w http.ResponseWriter, handleType domain.HandlerType) {
	if res == nil {
		w.Write(nil)
		w.WriteHeader(http.StatusOK)
		return
	}

	switch value := res.(type) {
	case *domain.GetChapterPagesPDFResponse:
		defer os.RemoveAll(value.PdfPath)
		w.Header().Add("Content-Type", "application/pdf")
		f, err := os.Open(value.PdfPath)
		if err != nil {
			logrus.Errorln("cannot open pdf file %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err := io.Copy(w,f); err != nil {
			logrus.Errorln("cannot copy pdf file %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		toSend, err := json.Marshal(res)
		if err != nil {
			err = errors.Wrapf(err, "Method %s ->", string(handleType))
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Errorln("unmarshal response error ", err.Error())
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(toSend)
	}
	w.WriteHeader(http.StatusOK)
}

func (em ErrHandler) handleTypeSwitcher(ctx context.Context, r *http.Request, handleType domain.HandlerType) (interface{}, error) {
	switch handleType {
	case domain.GetChapterList:
		return em.mangaHandler.GetChapterList(ctx, r)
	case domain.GetChapterPages:
		return em.mangaHandler.GetChapterPages(ctx, r)
	case domain.GetChapterPagesPDF:
		return em.mangaHandler.GetChapterPagesPDF(ctx, r)
	case domain.SendToEbook:
		return em.mangaHandler.SendToEbook(ctx, r)
	}
	return nil, customerrors.ErrUnknownType
}
