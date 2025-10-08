package model

import "time"

type User struct {
	ID        int64      `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	CreatedBy *string    `db:"created_by" json:"created_by"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by"`
	DeletedBy *string    `db:"deleted_by" json:"deleted_by"`
}
