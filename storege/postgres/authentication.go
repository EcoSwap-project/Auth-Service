package postgres

import (
	user "authentication_service/genproto/authentication_service"

	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Register(userReq *user.RegisterRequest) (*user.RegisterResponse, error){
	return nil, nil
}
