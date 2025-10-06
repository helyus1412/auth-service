package auth

import (
	"fmt"

	"github.com/helyus1412/auth-service/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Insert(*model.User) error
}

type repository struct {
	db     *sqlx.DB
	schema string
}

func NewRepository(db *sqlx.DB, schema string) Repository {
	if schema == "" {
		schema = "public"
	}

	return &repository{db, schema}
}

func (r *repository) Insert(user *model.User) error {
	query := fmt.Sprintf(`INSERT INTO %s.users (email, password) VALUES ($1, $2)`, r.schema)

	_, err := r.db.Exec(query, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}
