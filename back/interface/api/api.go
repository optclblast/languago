package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
	"languago/infrastructure/repository"
	"languago/pkg/auth"
	"languago/pkg/controllers/flashcards"
	"languago/pkg/controllers/users"
	errors2 "languago/pkg/errors"
	"languago/pkg/http/middleware"
	"languago/pkg/models/entities"
	"languago/pkg/models/requests/rest"
	"os"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type (
	API struct {
		*chi.Mux

		ID   uuid.UUID
		repo repository.DatabaseInteractor
		log  zerolog.Logger
		//auth                 auth.Authorizer
		errorsPresenter      errors2.ErrorsPersenter
		usersController      users.UsersController
		flashcardsController flashcards.FlashcardsController
	}
)

func NewAPI(cfg config.AbstractLoggerConfig, interactor repository.DatabaseInteractor) *API {
	logger := logger.ProvideLogger(cfg)
	errorsPresenter := errors2.NewErrorPresenter(logger)

	api := API{
		ID:              uuid.New(),
		repo:            interactor,
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

	auth.NewAuthorizer(
		logger,
		interactor.Database(),
		[]byte(os.Getenv("LANGUAGO_SECRET")),
	)

	mw := middleware.NewMiddleware(api.log)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"X-PINGOTHER",
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
			"Cache-Control",
			"Connection",
		},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(chimw.RequestID)
	//router.Use(mw.Options)
	router.Use(mw.LoggingMiddleware)
	router.Use(mw.Recovery)
	router.Use(mw.AuthMiddleware)

	router.Post("/signup", api.signUpHandler)
	router.Post("/signin", api.signInHandler)

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
	// dictionaryapi = "https://api.dictionaryapi.dev/api/v2/"
	// random word api
	randomwordapi = "http://random-words-api.vercel.app/word"
	// word types
	// Noun      = 0
	// Verb      = 1
	// Adjective = 2
	// Any       = 3
)

func (a *API) signUpHandler(w http.ResponseWriter, r *http.Request) {
	req := new(rest.SignUpRequest)

	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error read request body", err, http.StatusInternalServerError))
		return
	}

	err = json.Unmarshal(rawBody, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error bind request body to request model", err, http.StatusBadRequest))
		return
	}

	ctx, close := context.WithTimeout(r.Context(), 5*time.Second)
	defer close()
	err = a.usersController.CreateUser(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error create user", err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) signInHandler(w http.ResponseWriter, r *http.Request) {
	req := new(rest.SignInRequest)

	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error read request body", err, http.StatusInternalServerError))
		return
	}

	err = json.Unmarshal(rawBody, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(
			a.responseError(
				"error bind request body to request model", err,
				http.StatusBadRequest,
			),
		)
		return
	}

	ctx, close := context.WithTimeout(r.Context(), 5*time.Second)
	defer close()

	var user *entities.User = new(entities.User)

	err = a.repo.Database().WithTransaction(ctx, nil, func(tx *sql.Tx) error {
		user, err = a.repo.Database().SelectUser(ctx, repository.SelectUserParams{
			Login: req.Login,
		})

		return err
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.log.Info().Any("user: ", user)

	// todo password validation

	// token, err := auth.CreateToken(auth.ClaimJWTParams{
	// 	UserId: user.Id.String(),
	// })
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	//todo return token
}

func (a *API) randomWordHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(randomwordapi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error get random word", err, http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error read request body", err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (a *API) newFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	req := new(rest.NewFlashcardRequest)
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error read request body", err, http.StatusInternalServerError))
		return
	}
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error bind request body to a request model", err, http.StatusBadRequest))
		return
	}

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	err = a.flashcardsController.CreateFlashcard(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(a.responseError("error create flashcard", err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write(a.responseError("error parse card id", err, http.StatusBadRequest))
			return
		}
		cards, err := a.repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
			ID: id,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(a.responseError("not found", errors2.ErrNotFound, http.StatusNotFound))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(a.responseError("error select flashcard", err, http.StatusInternalServerError))
			return
		}

		response.Flashcards = cards
		resp, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(a.responseError("error marshal response body", err, http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else if deckId != "" {
		deckId, err := uuid.Parse(deckId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(a.responseError("error parse deck id", err, http.StatusBadRequest))
			return
		}
		if word != "" {
			cards, err := a.repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
				DeckID: deckId,
				Word:   word,
			})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(a.responseError("not found", errors2.ErrNotFound, http.StatusNotFound))
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(a.responseError("error select flashcard", err, http.StatusInternalServerError))
				return
			}

			response.Flashcards = cards

			resp, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(a.responseError("error marshal response body", err, http.StatusInternalServerError))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else if meaning != "" {
			cards, err := a.repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
				DeckID:  deckId,
				Meaning: meaning,
			})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(a.responseError("not found", errors2.ErrNotFound, http.StatusNotFound))
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(a.responseError("error select flashcard", err, http.StatusInternalServerError))
				return
			}

			response.Flashcards = cards

			resp, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(a.responseError("error marshal response body", err, http.StatusInternalServerError))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(a.responseError("error missing required fields", nil, http.StatusBadRequest))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error missing required fields", nil, http.StatusBadRequest))
		return
	}
}

func (a *API) deleteFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error parse card id", err, http.StatusBadRequest))
		return
	}

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	err = a.repo.Database().DeleteFlashcard(ctx, uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(a.responseError("error delete flashcard", err, http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) editFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	var request *rest.EditFlashcardRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error read card body", err, http.StatusBadRequest))
		return
	}

	if err = json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error bind request body to a request model", err, http.StatusBadRequest))
		return
	}

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()

	id, err := uuid.Parse(request.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error parse flashcard uuid", err, http.StatusBadRequest))
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(a.responseError("error missing required fields", err, http.StatusBadRequest))
		return
	}

	err = a.repo.Database().UpdateFlashcard(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(a.responseError("error update flashcard", err, http.StatusInternalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
}
