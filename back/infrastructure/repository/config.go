package repository

type (
	abstractDatabaseConfig interface {
		GetCredentials() DBCredentials
		IsMock() bool
	}
)
