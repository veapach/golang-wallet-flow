package app

import (
	"fmt"
	"log"

	"github.com/veapach/golang-wallet-flow/internal/closer"
	"github.com/veapach/golang-wallet-flow/internal/config"
	"github.com/veapach/golang-wallet-flow/internal/database"
)

type DIContainer struct {
	cfg *config.Config

	db *database.DB
}

func NewDIContainer(cfg *config.Config) *DIContainer {
	return &DIContainer{cfg: cfg}
}

func (c *DIContainer) DB() *database.DB {
	if c.db == nil {
		db, err := database.New(config.DSN(c.cfg))
		log.Println("dsn", config.DSN(c.cfg))
		if err != nil {
			panic(fmt.Sprintf("database: %v", err))
		}

		closer.Add("база данных", db.Close)

		c.db = db
	}
	return c.db
}
