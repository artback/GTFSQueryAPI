package query

import (
	"database/sql"
	"fmt"
	"github.com/artback/gtfsQueryGoApi/config"
	_ "github.com/lib/pq"
)

type Repository struct{ Db *sql.DB }

func (r *Repository) Connect(c config.DatabaseConfiguration) error {
	passwordArg := ""
	if len(c.Password) > 0 {
		passwordArg = "password=" + c.Password
	}

	var err error
	r.Db, err = sql.Open(c.Driver, fmt.Sprintf("host=%s port=%d user=%s %s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, passwordArg, c.Database))

	return err
}
