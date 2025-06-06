package gormRepo

import "ec-wallet/internal/domain"

type repository struct {
}

func NewRepository() domain.Repository {
	return &repository{}
}
