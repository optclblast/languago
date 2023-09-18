package repository

type (
	DBCredentials interface {
		GetAddress() string
		GetDriver() string
		GetUser() string
		GetSecret() string
	}

	DBCred struct {
		DbAddress string
		Driver    string
		User      string
		Secret    string
	}

	// redisCred struct {
	// 	DbAddr   string
	// 	User     string
	// 	Password string
	// }
)

// TODO
func buildConnString(cred DBCredentials) string {
	var connStr string
	switch cred.GetDriver() {
	case "pgsql":
	case "mysql":
	case "sqlite":
	case "cassandra":
	case "redis":
	default:
	}
	return connStr
}

// Relational database impl
func (c *DBCred) GetAddress() string {
	return c.DbAddress
}

func (c *DBCred) GetDriver() string {
	return c.Driver
}

func (c *DBCred) GetUser() string {
	return c.User
}

func (c *DBCred) GetSecret() string {
	return c.Secret
}

// Redis cred impl
// func (c *redisCred) GetAddress() string {
// 	return c.DbAddr
// }

// func (c *redisCred) GetDriver() string {
// 	return "redis"
// }

// func (c *redisCred) GetUser() string {
// 	return c.User
// }

// func (c *redisCred) GetSecret() string {
// 	return c.Password
// }
