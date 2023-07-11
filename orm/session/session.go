package session

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lz-nsc/cupcake/orm/log"
	"github.com/lz-nsc/cupcake/orm/schema"
	"github.com/lz-nsc/cupcake/orm/translator"
)

type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ DB = (*sql.DB)(nil)
var _ DB = (*sql.Tx)(nil)

type Session struct {
	db        DB
	sql       strings.Builder
	sqlVars   []interface{}
	trans     translator.Translator
	schema    *schema.Schema
	statement *statement
}

func New(db *sql.DB, trans translator.Translator) *Session {
	return &Session{
		db:        db,
		trans:     trans,
		statement: &statement{},
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.statement = &statement{}
}

func (s Session) DB() DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	// Reset session after successfully execute previous query
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	// Reset session after successfully execute previous query
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	// Reset session after successfully execute previous query
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Begin() (err error) {
	log.Debug("begin transaction")
	db, ok := s.DB().(*sql.DB)
	if !ok {
		err = errors.New("transaction has alread begun")
		log.Error(err)
		return
	}
	s.db, err = db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("begin transaction succeed")
	return
}

func (s *Session) Commit() (err error) {
	log.Debug("commit transaction")
	tx, ok := s.DB().(*sql.Tx)
	if !ok {
		err = errors.New("need to begin transaction first")
		log.Error(err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("commit transaction succeed")
	return
}

func (s *Session) Rollback() (err error) {
	log.Debug("rollback transaction")
	tx, ok := s.DB().(*sql.Tx)
	if !ok {
		err = errors.New("need to begin transaction first")
		log.Error(err)
		return
	}
	err = tx.Rollback()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("rollback transaction succeed")
	return
}
