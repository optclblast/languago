package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository"
	"languago/internal/pkg/repository/postgresql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type (
	API struct {
		*mux.Router
		AuthProvider *Authorizer
		Repo         repository.DatabaseInteractor
		log          logger.Logger
	}
)

func NewAPI(cfg config.AbstractLoggerConfig, interactor repository.DatabaseInteractor) *API {
	return &API{
		Router: mux.NewRouter(),
		Repo:   interactor,
		log:    logger.ProvideLogger(cfg),
	}
}
func (api *API) routes() {
	api.HandleFunc("/randomword", api.randomWord()).Methods(http.MethodGet) // test
	// api.HandleFunc("/flashcard", api.getFlashcard()).Methods(http.MethodGet)
	api.HandleFunc("/flashcard", api.newFlashcard()).Methods(http.MethodPost)
	// api.HandleFunc("/flashcard", api.deleteFlashcard()).Methods(http.MethodPost)
	// api.HandleFunc("/flashcard", api.editFlashcard()).Methods(http.MethodPost)
}

const (
	// https://dictionaryapi.dev/ API docs
	dictionaryapi = "https://api.dictionaryapi.dev/api/v2/"
	// random word api
	randomwordapi = "http://random-words-api.vercel.app/word"
	// word types
	Noun      = 0
	Verb      = 1
	Adjective = 2
	Any       = 3
)

func (a *API) Init() {
	a.routes()
	a.log.Info("api initialized")
}

func (a *API) randomWord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(randomwordapi)
		if err != nil {
			a.responseError(w, fmt.Errorf("error getting response from %s: %w", randomwordapi, err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			a.responseError(w, fmt.Errorf("error reading response body from %s: %w", randomwordapi, err), http.StatusInternalServerError)
			return
		}

		a.log.Debug(string(body))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func (a *API) newFlashcard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req rest.NewFlashcardRequest
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			a.responseError(w, fmt.Errorf("error reading request body: %w", err), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(body, &req); err != nil {
			a.responseError(w, fmt.Errorf("error parsing request body: %w", err), http.StatusBadRequest)
			return
		}

		ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		defer c()

		_, err = a.Repo.Database().CreateFlashcard(ctx, postgresql.CreateFlashcardParams{
			ID:      uuid.New(),
			Word:    sql.NullString{String: req.Content.WordInTarget, Valid: true},
			Meaning: sql.NullString{String: req.Content.WordInNative, Valid: true},
			Usage:   req.Content.UsageExamples,
		})
		if err != nil {
			a.responseError(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// func (a *API) getFlashcard() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
