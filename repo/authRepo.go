package repo

import "example.com/app/domain"

type AuthRepo interface {
	Login(username string, password string) (*domain.User, string, error)
}

