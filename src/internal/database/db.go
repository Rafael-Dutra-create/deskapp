package database

import (
	"database/sql"
	"deskapp/src/internal/utils"
	_ "github.com/lib/pq"
	"sync"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	logger     *utils.Logger
)

func InitDB(dsn string) (*sql.DB, error) {
	logger = utils.NewLogger()
	var err error
	once.Do(func() {
		dbInstance, err = sql.Open("postgres", dsn)
		if err != nil {
			return
		}

		err = dbInstance.Ping()
		if err != nil {
			dbInstance.Close()
			dbInstance = nil
			return
		}
		dbInstance.SetMaxOpenConns(20)
		logger.Infof("Database connection pool initialized successfully: %s", dsn)
	})
	return dbInstance, err
}

func GetDB() *sql.DB {
	if dbInstance == nil {
		panic("Database connection has not been initialized. Call InitDB() first.")
	}
	return dbInstance
}
