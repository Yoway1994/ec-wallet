package gormRepoImpl

import (
	gormrepo "ec-wallet/internal/domain/gorm_repo"

	"gorm.io/gorm"
)

type repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) gormrepo.Repo {
	return &repository{Db: db}
}

func (repo *repository) Begin() *gorm.DB {
	return repo.Db.Begin()
}

func (repo *repository) Commit(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (repo *repository) RollBack(tx *gorm.DB) *gorm.DB {
	return repo.Db.Rollback()
}
