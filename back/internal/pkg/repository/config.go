package repository

type (
	abstractDatabaseConfig interface {
		GetCredentials() DBCredentials
	}
)
