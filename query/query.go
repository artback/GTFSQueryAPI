package query

import (
	"database/sql"
	"fmt"
	"github.com/artback/gtfsQueryGoApi/config"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

type Repository struct{ Db *sql.DB }

func (r *Repository) Connect(c config.DatabaseConfiguration) error {
	passwordArg := ""
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = c.Password
	}
	host := os.Getenv("POSTGRESS_HOST")
	if host == "" {
		host = c.Host
	}
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if port == 0 {
		port = c.Port
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = c.User
	}
	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = c.Database
	}
	if len(pass) > 0 {
		passwordArg = "password=" + pass
	}
	var err error
	db_string := fmt.Sprintf("host=%s port=%d user=%s %s dbname=%s sslmode=disable",
		host, port, user, passwordArg, db)
	r.Db, err = sql.Open(c.Driver, db_string)
	return err
}
