package storage

import "languago/internal/pkg/repository"

type RepositoryInteractor struct {
	repo repository.DatabaseInteractor
}

func NewRepositoryInteractor(r repository.DatabaseInteractor) *RepositoryInteractor {
	return &RepositoryInteractor{repo: r}
}

//func (r *RepositoryInteractor)
