package orm

import (
	"database/sql"
	"fmt"

	"github.com/lz-nsc/cupcake/orm/log"
	"github.com/lz-nsc/cupcake/orm/session"
	"github.com/lz-nsc/cupcake/orm/translator"
)

type ORMEngine struct {
	db    *sql.DB
	trans translator.Translator
}
type TransactionFunc func(*session.Session) (interface{}, error)

func NewORMEngine(driver string, source string) (engine *ORMEngine, err error) {
	trans, ok := translator.GetTranslator(driver)
	if !ok {
		err = fmt.Errorf("driver %s is not supported", driver)
		log.Error(err)
		return
	}
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
		db:    db,
		trans: trans,
	}

	log.Infof("successfully connect to %s database %s", driver, source)
	return
}

func (oe *ORMEngine) Close() {
	if err := oe.db.Close(); err != nil {
		log.Error("failed to close connection to database")
	}
	log.Info("successfully close connection to database")
}

func (oe *ORMEngine) NewSession() *session.Session {
	return session.New(oe.db, oe.trans)
}

func (oe *ORMEngine) Transaction(fn TransactionFunc) (result interface{}, err error) {
	// Create new session for transaction
	s := oe.NewSession()
	err = s.Begin()
	if err != nil {
		log.Errorf("failed to begin transaction, err: %v", err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			// Rollback of commit fails
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()
	return fn(s)
}
