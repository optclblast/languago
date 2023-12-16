package test

import (
	"context"
	"languago/infrastructure/config"
	"languago/infrastructure/repository"
	"languago/test/generators"
	"sync"
	"testing"

	"github.com/google/uuid"
)

// Script to fill the database with data
func TestFillDb(t *testing.T) {
	var errs []error

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()
			cfg := &config.DatabaseConfig{
				DatabaseAddress: "localhost:5432",
				DatabaseDriver:  "postgres",
				DatabaseUser:    "languago",
				DatabaseSecret:  "languago",
			}

			db, err := repository.NewDatabaseInteractor(cfg)
			if err != nil {
				panic(err)
			}

			defer db.Database().Close()

			// var wg sync.WaitGroup

			// var (
			// 	usersBatch      []models.User      = make([]models.User, 100)
			// 	decksBatch      []models.Deck      = make([]models.Deck, 5000)
			// 	flashcardsBatch []models.Flashcard = make([]models.Flashcard, 125000)
			// )

			ctx := context.Background()

			for i := 0; i < 20; i++ {
				uid := uuid.New()
				// usersBatch[i] = models.User{
				// 	Id:    uid,
				// 	Login: generators.RandStringRunes(30),
				// }

				err := db.Database().CreateUser(ctx, repository.CreateUserParams{
					ID:       uid,
					Login:    generators.RandStringRunes(30),
					Password: generators.RandStringRunes(30),
				})

				if err != nil {
					t.Logf("CreateUser err: %s", err.Error())
					errs = append(errs, err)
				}

				for j := 0; j < 5; j++ {
					did := uuid.New()
					// decksBatch[j] = models.Deck{
					// 	Id:    did.String(),
					// 	Name:  generators.RandStringRunes(30),
					// 	Owner: uid,
					// }

					err := db.Database().CreateDeck(ctx, repository.CreateDeckParams{
						ID:    did,
						Name:  generators.RandStringRunes(30),
						Owner: uid,
					})
					if err != nil {
						t.Logf("CreateDeck err: %s", err.Error())
						errs = append(errs, err)
					}

					for p := 0; p < 20; p++ {
						cid := uuid.New()

						// 	ID      uuid.UUID `db:"id" json:"id"`
						// 	Word    string    `db:"word" json:"word"`
						// 	Meaning string    `db:"meaning" json:"meaning"`
						// 	Usage   []string  `db:"usage" json:"usage"`
						//       }

						err := db.Database().CreateFlashcard(ctx, repository.CreateFlashcardParams{
							ID:      cid,
							Word:    generators.RandStringRunes(8),
							Meaning: generators.RandStringRunes(8),
							Usage:   generators.RandStringSlice(0, 10),
						})
						if err != nil {
							t.Logf("CreateFlashcard err: %s", err.Error())
							errs = append(errs, err)
						}

						err = db.Database().AddToDeck(ctx, repository.AddToDeckParams{
							DeckID:      did,
							FlashcardID: cid,
						})
						if err != nil {
							t.Logf("AddToDeck err: %s", err.Error())
							errs = append(errs, err)
						}
					}
				}
			}
		}()
	}

	wg.Wait()

	t.Logf("errors: %v", errs)
}

func TestFillDbCards(t *testing.T) {
	cfg := &config.DatabaseConfig{
		DatabaseAddress: "localhost:5432",
		DatabaseDriver:  "postgres",
		DatabaseUser:    "languago",
		DatabaseSecret:  "languago",
	}

	db, err := repository.NewDatabaseInteractor(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Database().Close()

	// var wg sync.WaitGroup

	// var (
	// 	usersBatch      []models.User      = make([]models.User, 100)
	// 	decksBatch      []models.Deck      = make([]models.Deck, 5000)
	// 	flashcardsBatch []models.Flashcard = make([]models.Flashcard, 125000)
	// )

	var errs []error

	ctx := context.Background()

	for i := 0; i < 25; i++ {
		uid := uuid.New()
		// usersBatch[i] = models.User{
		// 	Id:    uid,
		// 	Login: generators.RandStringRunes(30),
		// }

		err := db.Database().CreateUser(ctx, repository.CreateUserParams{
			ID:       uid,
			Login:    generators.RandStringRunes(30),
			Password: generators.RandStringRunes(30),
		})

		if err != nil {
			t.Logf("CreateUser err: %s", err.Error())
			errs = append(errs, err)
		}

		for j := 0; j < 10; j++ {
			did := uuid.New()
			// decksBatch[j] = models.Deck{
			// 	Id:    did.String(),
			// 	Name:  generators.RandStringRunes(30),
			// 	Owner: uid,
			// }

			err := db.Database().CreateDeck(ctx, repository.CreateDeckParams{
				ID:    did,
				Name:  generators.RandStringRunes(30),
				Owner: uid,
			})
			if err != nil {
				t.Logf("CreateDeck err: %s", err.Error())
				errs = append(errs, err)
			}

			for p := 0; p < 25; p++ {
				cid := uuid.New()

				// 	ID      uuid.UUID `db:"id" json:"id"`
				// 	Word    string    `db:"word" json:"word"`
				// 	Meaning string    `db:"meaning" json:"meaning"`
				// 	Usage   []string  `db:"usage" json:"usage"`
				//       }

				err := db.Database().CreateFlashcard(ctx, repository.CreateFlashcardParams{
					ID:      cid,
					Word:    generators.RandStringRunes(8),
					Meaning: generators.RandStringRunes(8),
					Usage:   generators.RandStringSlice(0, 10),
				})
				if err != nil {
					t.Logf("CreateFlashcard err: %s", err.Error())
					errs = append(errs, err)
				}

				err = db.Database().AddToDeck(ctx, repository.AddToDeckParams{
					DeckID:      did,
					FlashcardID: cid,
				})
				if err != nil {
					t.Logf("AddToDeck err: %s", err.Error())
					errs = append(errs, err)
				}
			}
		}
	}

	t.Logf("errors: %v", errs)
}
