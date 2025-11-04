package repository

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
)

type ModuloRespository struct {
	*repository.BaseRepository
}

func NewModuloRepository(db *sql.DB) *ModuloRespository {
	base := repository.NewBaseRepository(db, "modulo", "ep_dw")
	return &ModuloRespository{
		base,
	}
}