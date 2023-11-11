package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"languago/internal/pkg/auth"
	"languago/internal/pkg/config"
	"languago/internal/pkg/controllers/flashcards"
	"languago/internal/pkg/controllers/users"
	errors2 "languago/internal/pkg/errors"
	"languago/internal/pkg/http/middleware"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository"
	"os"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type (
	API struct {
		*chi.Mux
		Repo            repository.DatabaseInteractor
		log             logger.Logger
		errorsPresenter errors2.ErrorsPersenter

		usersController      users.UsersController
		flashcardsController flashcards.FlashcardsController
	}
)

func NewAPI(cfg config.AbstractLoggerConfig, interactor repository.DatabaseInteractor) *API {
	logger := logger.ProvideLogger(cfg)
	errorsPresenter := errors2.NewErrorPresenter(logger)

	api := API{
		Repo:            interactor,
		log:             logger,
		errorsPresenter: errorsPresenter,
		flashcardsController: flashcards.NewFlashcardsController(
			logger,
			interactor,
		),
		usersController: users.NewUsersController(
			logger,
			interactor,
		),
	}

	router := chi.NewRouter()

	mw := middleware.NewMiddleware(api.log, auth.NewAuthorizer(
		logger,
		interactor.Database(),
		[]byte(os.Getenv("LANGUAGO_SECRET")),
	))

	router.Use(chimw.RequestID)
	router.Use(mw.LoggingMiddleware)
	router.Use(mw.AuthMiddleware)
	router.Use(mw.Recovery)

	router.Post("/signup", api.signUpHandler)

	router.Get("/randomword", api.randomWordHandler)

	router.Get("/flashcard", api.getFlashcardHandler)
	router.Post("/flashcard", api.newFlashcardHandler)
	router.Delete("/flashcard", api.deleteFlashcardHandler)
	router.Put("/flashcard", api.editFlashcardHandler)

	api.Mux = router

	return &api
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

func (a *API) signUpHandler(w http.ResponseWriter, r *http.Request) {
	req := new(rest.SignUpRequest)

	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		a.Response(
			w,
			fmt.Errorf("error read request body: %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	err = json.Unmarshal(rawBody, &req)
	if err != nil {
		a.Response(
			w,
			fmt.Errorf("error read request body: %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	ctx, close := context.WithTimeout(r.Context(), 500*time.Second)
	defer close()
	err = a.usersController.CreateUser(ctx, req)
	if err != nil {
		a.Response(
			w,
			fmt.Errorf("error create user: %w", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) randomWordHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(randomwordapi)
	if err != nil {
		a.Response(
			w,
			fmt.Errorf("error getting response from %s: %w", randomwordapi, err),
			http.StatusInternalServerError,
		)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.Response(
			w,
			fmt.Errorf("error reading response body from %s: %w", randomwordapi, err),
			http.StatusInternalServerError,
		)
		return
	}

	a.log.Debug(string(body), nil)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (a *API) newFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	req := new(rest.NewFlashcardRequest)
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		a.Response(w, fmt.Errorf("error reading request body: %w", err), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &req); err != nil {
		a.Response(w, fmt.Errorf("error parsing request body: %w", err), http.StatusBadRequest)
		return
	}

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	err = a.flashcardsController.CreateFlashcard(ctx, req)
	if err != nil {
		a.Response(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
		return
	}

	a.Response(w, nil, http.StatusOK)
}

func (a *API) getFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	deckId := r.URL.Query().Get("deck_id")
	word := r.URL.Query().Get("word")
	meaning := r.URL.Query().Get("meaning")

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	w.Header().Add("Content-Type", "application/json")

	// TODO refactoring of this abomination
	response := new(rest.GetFlashcardResponse)
	if id != "" {
		id, err := uuid.Parse(id)
		if err != nil {
			a.Response(w, err, http.StatusBadRequest)
			return
		}
		cards, err := a.Repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
			ID: id,
		})
		if err != nil {
			a.Response(w, err, http.StatusBadRequest)
			return
		}

		if cards == nil {
			a.Response(w, nil, http.StatusNotFound)
			return
		}

		response.Flashcards = cards
		resp, err := json.Marshal(response)
		if err != nil {
			a.Response(w, fmt.Errorf("internal error"), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else if deckId != "" {
		deckId, err := uuid.Parse(deckId)
		if err != nil {
			a.Response(w, err, http.StatusBadRequest)
			return
		}
		if word != "" {
			cards, err := a.Repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
				DeckID: deckId,
				Word:   word,
			})
			if err != nil {
				a.Response(w, err, http.StatusBadRequest)
				return
			}

			if cards == nil {
				a.Response(w, nil, http.StatusNotFound)
				return
			}

			response.Flashcards = cards
			if err != nil {
				a.Response(w, fmt.Errorf("empty flashcard"), http.StatusBadRequest)
				return
			}

			resp, err := json.Marshal(response)
			if err != nil {
				a.Response(w, fmt.Errorf("internal error"), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else if meaning != "" {
			cards, err := a.Repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
				DeckID:  deckId,
				Meaning: meaning,
			})
			if err != nil {
				a.Response(w, err, http.StatusBadRequest)
				return
			}

			if cards == nil {
				a.Response(w, nil, http.StatusNotFound)
				return
			}

			response.Flashcards = cards
			if err != nil {
				a.Response(w, fmt.Errorf("empty flashcard"), http.StatusBadRequest)
				return
			}

			a.Response(w, response, http.StatusOK)
		} else {
			a.Response(w, nil, http.StatusBadRequest)
			return
		}
	} else {
		a.Response(w, nil, http.StatusBadRequest)
		return
	}
}

func (a *API) deleteFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		a.Response(w, err, http.StatusBadRequest)
		return
	}

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	err = a.Repo.Database().DeleteFlashcard(ctx, uuid)
	if err != nil {
		a.Response(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) editFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	var request *rest.EditFlashcardRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		a.Response(w, fmt.Errorf("error reading request body: %w", err), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		a.Response(w, fmt.Errorf("error parsing request body: %w", err), http.StatusBadRequest)
		return
	}

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.Response(w, fmt.Errorf("error invalid id: %w", err), http.StatusBadRequest)
		return
	}
	params := repository.UpdateFlashcardParams{
		ID: id,
	}

	switch {
	case request.WordInNative != "":
		params.Meaning = request.WordInNative
	case request.WordInTarget != "":
		params.Word = request.WordInTarget
	case request.UsageExamples != nil:
		params.Usage = request.UsageExamples
	default:
		a.responseError(w, nil, http.StatusExpectationFailed)
		return
	}

	err = a.Repo.Database().UpdateFlashcard(ctx, params)
	if err != nil {
		a.Response(w, fmt.Errorf("error editing flashcard: %w", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) Response(w http.ResponseWriter, resp any, status int) {
	if status != 200 {
		a.responseError(w, resp.(error), status)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	respRaw, err := json.Marshal(resp)
	if err != nil {
		a.responseError(w, fmt.Errorf("internal error"), http.StatusBadRequest)
		return
	}

	w.WriteHeader(status)
	w.Write(respRaw)
}
