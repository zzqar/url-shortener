package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	resp "url-shortener/internal/http-server/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

type Saver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alis,omitempty"`
}

func New(log *slog.Logger, urlSaver Saver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.url.save.New"
		w.Header().Set("Content-Type", "application/json")

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		var request Request

		if err := render.DecodeJSON(r.Body, &request); err != nil {
			log.Error("failed to decode request", sl.Err(err))
			responseError(w, r, http.StatusInternalServerError, "failed to decode request")
			return
		}
		log.Info("request decoded", slog.Any("request", request))

		if err := validator.New().Struct(request); err != nil {
			var validatorErr validator.ValidationErrors
			errors.As(err, &validatorErr)

			log.Error("failed to validate request", sl.Err(err))
			responseError(w, r, http.StatusInternalServerError, resp.ValidationErrors(validatorErr))
			return
		}

		alias := request.Alias
		if alias == "" {
			alias = random.RandStringBytes(cfg.RandomString.Length)
		}
		id, err := urlSaver.SaveURL(request.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("alias already exists", slog.String("alias", alias))

			responseError(w, r, http.StatusConflict, "alias already exists")
			return
		}
		if err != nil {
			log.Error("failed to save URL", sl.Err(err))

			responseError(w, r, http.StatusInternalServerError, "failed to save URL")
			return
		}
		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

func responseError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)
	render.JSON(w, r, resp.Error(status, msg))
}
