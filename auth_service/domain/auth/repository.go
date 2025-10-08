package auth

import (
	"fmt"
	"time"

	"github.com/helyus1412/auth-service/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Insert(*model.User) error
	GetByEmail(string) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(*model.User) error
	GetByID(int64) (*model.User, error)
	SoftDelete(int64) error
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
	queryPrep, err := r.db.Preparex(fmt.Sprintf(`INSERT INTO %s.users (email, password) VALUES ($1, $2)`, r.schema))
	if err != nil {
		return err
	}

	_, err = queryPrep.Exec(user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByEmail(email string) (user *model.User, err error) {
	var res model.User

	query, err := r.db.Preparex(fmt.Sprintf(`SELECT * from %s."users" where email = $1`, r.schema))
	if err != nil {
		return nil, err
	}

	err = query.Get(&res, email)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *repository) GetAll() (users []model.User, err error) {
	query, err := r.db.Preparex(fmt.Sprintf(`select * from %s."users" where deleted_at is null`, r.schema))
	if err != nil {
		return nil, err
	}

	err = query.Select(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Update(user *model.User) error {
	queryPrep, err := r.db.Preparex(fmt.Sprintf(`update %s.users set email=$1, password=$2,updated_at=$3 WHERE id = $4`, r.schema))
	if err != nil {
		return err
	}

	_, err = queryPrep.Exec(user.Email, user.Password, time.Now(), user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(id int64) (user *model.User, err error) {
	var res model.User

	query, err := r.db.Preparex(fmt.Sprintf(`SELECT * from %s."users" where id = $1`, r.schema))
	if err != nil {
		return nil, err
	}

	err = query.Get(&res, id)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *repository) SoftDelete(id int64) error {
	queryPrep, err := r.db.Preparex(fmt.Sprintf(`update %s.users set deleted_at = $1 where id = $2`, r.schema))
	if err != nil {
		return err
	}

	_, err = queryPrep.Exec(time.Now(), id)

	if err != nil {
		return err
	}

	return nil
}
