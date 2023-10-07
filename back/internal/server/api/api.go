package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"languago/internal/pkg/config"
	"languago/internal/pkg/http/middleware"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/entities"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type (
	API struct {
		*mux.Router
		Repo repository.DatabaseInteractor
		log  logger.Logger
		stop chan struct{}
	}
)

func NewAPI(cfg config.AbstractLoggerConfig, interactor repository.DatabaseInteractor) *API {
	api := API{
		Repo: interactor,
		log:  logger.ProvideLogger(cfg),
		stop: make(chan struct{}),
	}

	router := mux.NewRouter()

	mw := middleware.NewMiddleware(api.log, api.stop)
	router.Use(mw.Close)
	router.Use(mw.LoggingMiddleware)
	router.Use(mw.AuthMiddleware)
	router.Use(mw.RequestValidationMiddleware)
	router.Use(mw.Recovery)

	api.Router = router

	return &api
}

func (a *API) Stop() {
	a.log.Warn("closing api", logger.LogFieldPair("time", time.Now()))
	a.stop <- struct{}{}
}

func (api *API) routes() {
	// Flashcards
	// Returns an array of random words. For now, len(arr) == 1
	// Example:
	// [
	// 	"hemps"
	// ]
	api.HandleFunc("/randomword", api.randomWordHandler()).Methods(http.MethodGet)

	// Example: /flashcard?word=truck&deck_id=60dd7b63-cc25-4928-9229-41bed0f33bd3. This will retern you
	// a flashcards included in the deck with id == "60dd7b63-cc25-4928-9229-41bed0f33bd3" and word
	// parameter == "truck". If there are no such cards response will be 404.
	api.HandleFunc("/flashcard", api.getFlashcardHandler()).Methods(http.MethodGet)

	// Endpoint for creating new flashcards.
	// Request body example:
	// {
	// 	  "native_lang": "ru",
	// 	  "target_lang": "en",
	// 	  "content": {
	// 		  "word_in_native": "Слово",
	// 		  "word_in_target": "Word",
	// 		  "usage": ["Woooordsss"]
	// 	  }
	// }
	// 200 if created successfully, 400 if body is invalid and 500 if there are any server-side issues
	api.HandleFunc("/flashcard", api.newFlashcardHandler()).Methods(http.MethodPost)

	// Endpoind for deleting flashcards. Request example:
	// /flashcard?id=d33c0c35-6549-4978-bd6e-c8a22bbd2f6b
	// If ^ is ok, then flashcard with id specified will be deleted.
	api.HandleFunc("/flashcard", api.deleteFlashcardHandler()).Methods(http.MethodDelete)

	// Endpoint for editing existing flashcards.
	// Body example:
	// {
	// 	  "id":"d33c0c35-6549-4978-bd6e-c8a22bbd2f6b"
	// 	  "word_in_native":"Cow"
	// 	  "word_in_target":"Корова"
	// 	  "usage":[
	// 		  "Вот станет жарко, корова хвост поднимет и побежит.",
	// 		  "Недалеко от места, где легла корова, протекал ручей, бравший начало из Кастальского ключа.",
	// 		  "Треневы держали корову, и мой отец всегда спрашивал: \“Леночка, а чья будет корова, когда ты переедешь к нам?\” \“Ваша\”."
	//     ]
	// }
	// Any of this fields may be empty, except "id" field.
	api.HandleFunc("/flashcard", api.editFlashcardHandler()).Methods(http.MethodPut)

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
	a.log.Info("api initialized", nil)
}

func (a *API) randomWordHandler() http.HandlerFunc {
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

		a.log.Debug(string(body), nil)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func (a *API) newFlashcardHandler() http.HandlerFunc {
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

		err = a.Repo.Database().CreateFlashcard(ctx, repository.CreateFlashcardParams{
			ID:      uuid.New(),
			Word:    req.Content.WordInTarget,
			Meaning: req.Content.WordInNative,
			Usage:   req.Content.UsageExamples,
		})
		if err != nil {
			a.responseError(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *API) getFlashcardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		deckId := r.URL.Query().Get("deck_id")
		word := r.URL.Query().Get("word")
		meaning := r.URL.Query().Get("meaning")

		ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		defer c()

		// TODO refactoring of this abomination
		var response *rest.GetFlashcardResponse
		if id != "" {
			id, err := uuid.Parse(id)
			if err != nil {
				a.responseError(w, err, http.StatusBadRequest)
				return
			}
			card, err := a.Repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
				ID: id,
			})
			if err != nil {
				a.responseError(w, err, http.StatusBadRequest)
				return
			}

			if card == nil {
				a.responseError(w, nil, http.StatusNotFound)
				return
			}

			response.Flashcards = []*entities.Flashcard{card}
			resp, err := json.Marshal(response)
			if err != nil {
				a.responseError(w, fmt.Errorf("internal error"), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else if deckId != "" {
			deckId, err := uuid.Parse(deckId)
			if err != nil {
				a.responseError(w, err, http.StatusBadRequest)
				return
			}
			if word != "" {
				card, err := a.Repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
					DeckID: deckId,
					Word:   word,
				})
				if err != nil {
					a.responseError(w, err, http.StatusBadRequest)
					return
				}

				if card == nil {
					a.responseError(w, nil, http.StatusNotFound)
					return
				}

				response.Flashcards = []*entities.Flashcard{card}
				if err != nil {
					a.responseError(w, fmt.Errorf("empty flashcard"), http.StatusBadRequest)
					return
				}

				resp, err := json.Marshal(response)
				if err != nil {
					a.responseError(w, fmt.Errorf("internal error"), http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(resp)
			} else if meaning != "" {
				card, err := a.Repo.Database().SelectFlashcard(ctx, repository.SelectFlashcardParams{
					DeckID:  deckId,
					Meaning: meaning,
				})
				if err != nil {
					a.responseError(w, err, http.StatusBadRequest)
					return
				}

				if card == nil {
					a.responseError(w, nil, http.StatusNotFound)
					return
				}

				response.Flashcards = []*entities.Flashcard{card}
				if err != nil {
					a.responseError(w, fmt.Errorf("empty flashcard"), http.StatusBadRequest)
					return
				}

				resp, err := json.Marshal(response)
				if err != nil {
					a.responseError(w, fmt.Errorf("internal error"), http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write(resp)
			} else {
				a.responseError(w, nil, http.StatusBadRequest)
				return
			}
		} else {
			a.responseError(w, nil, http.StatusBadRequest)
			return
		}
	}
}

func (a *API) deleteFlashcardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			a.responseError(w, err, http.StatusBadRequest)
			return
		}

		ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		defer c()

		err = a.Repo.Database().DeleteFlashcard(ctx, uuid)
		if err != nil {
			a.responseError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (a *API) editFlashcardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request *rest.EditFlashcardRequest
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			a.responseError(w, fmt.Errorf("error reading request body: %w", err), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(body, &request); err != nil {
			a.responseError(w, fmt.Errorf("error parsing request body: %w", err), http.StatusBadRequest)
			return
		}

		ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
		defer c()

		id, err := uuid.Parse(request.Id)
		if err != nil {
			a.responseError(w, fmt.Errorf("error invalid id: %w", err), http.StatusBadRequest)
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
			a.responseError(w, fmt.Errorf("error editing flashcard: %w", err), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
