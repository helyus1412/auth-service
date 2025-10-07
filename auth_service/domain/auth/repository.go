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
	queryPrep, err := r.db.Prepare(fmt.Sprintf(`INSERTt INTO %s.users (email, password) VALUES ($1, $2)`, r.schema))
	if err != nil {
		return err
	}

	_, err = queryPrep.Exec(user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}
