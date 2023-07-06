package orm

import (
	"database/sql"

	"github.com/lz-nsc/cupcake/orm/log"
	"github.com/lz-nsc/cupcake/orm/session"
)

type ORMEngine struct {
	db *sql.DB
}

func NewORMEngine(driver string, source string) (engine *ORMEngine, err error) {
	// Open a database
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Verify connection to the database
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	engine = &ORMEngine{
		db: db,
	}

	log.Infof("Successfully connect to %s database %s", driver, source)
	return
}

func (oe *ORMEngine) Close() {
	if err := oe.db.Close(); err != nil {
		log.Error("Failed to close connection to database")
	}
	log.Info("Successfully close connection to database")
}

func (oe *ORMEngine) NewSession() *session.Session {
	return session.New(oe.db)
}
