package api

import (
	"encoding/json"
	"fmt"
	"io"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository"
	"net/http"

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
			a.log.Warn(fmt.Sprintf("error getting response from %s: %s", randomwordapi, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting random word"))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			a.log.Warn(fmt.Sprintf("error reading response body from %s: %s", randomwordapi, err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting random word"))
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
		if err != nil {
			err = fmt.Errorf("error reading request body: %w", err)
			a.log.Warn(err)
			a.response(w, err)
		}
		if err = json.Unmarshal(body, &req); err != nil {
			err = fmt.Errorf("error parsing request body: %w", err)
			a.log.Warn(err)
			a.response(w, err)
		}

		fmt.Println(req)
	}
}

// func (a *API) getFlashcard() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
