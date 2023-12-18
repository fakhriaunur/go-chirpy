package database

import (
	"errors"
	"log"
	"time"
)

type Session struct {
	ID          string    `json:"id"`
	IsRevoked   bool      `json:"is_revoked"`
	TimeRevoked time.Time `json:"time_revoked"`
}

func (db *DB) CreateSession(id string) (Session, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Println(err)
		return Session{}, err
	}

	session := Session{
		ID:          id,
		IsRevoked:   false,
		TimeRevoked: time.Time{},
	}
	dbStructure.Sessions[id] = session

	if err := db.writeDB(dbStructure); err != nil {
		log.Println(err)
		return Session{}, err
	}

	return session, nil

	// if _, ok := dbStructure.Sessions[id]; !ok {
	// 	session := Session{
	// 		ID:        id,
	// 		IsRevoked: false,
	// 	}
	// 	dbStructure.Sessions[id] = session

	// 	if err := db.writeDB(dbStructure); err != nil {
	// 		return Session{}, err
	// 	}

	// 	return session, nil
	// }

	// return Session{}, ErrAlreadyExist
}

func (db *DB) GetSession(id string) (Session, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Session{}, err
	}

	session, ok := dbStructure.Sessions[id]
	if !ok {
		return Session{}, ErrNotExist
	}

	return session, nil
}

func (db *DB) RevokeSession(id string) (Session, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Session{}, err
	}

	session, ok := dbStructure.Sessions[id]
	if !ok {
		return Session{}, ErrNotExist
	}

	if session.IsRevoked {
		return Session{}, errors.New("session already revoked")
	}

	session.IsRevoked = true
	session.TimeRevoked = time.Now()
	dbStructure.Sessions[id] = session

	if err := db.writeDB(dbStructure); err != nil {
		return Session{}, err
	}

	return session, nil

}
