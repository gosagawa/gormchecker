package a

import (
	"github.com/jinzhu/gorm"
)

// expected style
func f1() (*User, error) {

	db, err := getConnection()
	if err != nil {
		return nil, err
	}

	u := User{}
	db = db.Where("column_a = xxx")
	db = db.Where("column_b = xxx")
	db = db.Find(&u)

	if db.RecordNotFound() {
		return nil, nil
	}
	return &u, db.Error
}

// using pipe
func f2() (*User, error) {

	db, err := getConnection()
	if err != nil {
		return nil, err
	}

	u := User{}
	db = db.Where("column_a = xxx").Where("column_b = xxx") // want "do not use pipe"
	db = db.Find(&u)

	if db.RecordNotFound() {
		return nil, nil
	}
	return &u, db.Error
}

// not db pipe
// TODO:(sagawa) strictly check gorm.DB or not
func f3() (*User, error) {

	xdb, err := getConnection()
	if err != nil {
		return nil, err
	}

	u := User{}
	xdb = xdb.Where("column_a = xxx").Where("column_b = xxx") // ok
	xdb = xdb.Find(&u)

	if xdb.RecordNotFound() {
		return nil, nil
	}
	return &u, xdb.Error
}

func getConnection() (*gorm.DB, error) {
	return nil, nil
}

// User represents a row from 'test.users'.
type User struct {
	ID      int `json:"id"`       // id
	ColumnA int `json:"column_a"` // column_a
	ColumnB int `json:"column_b"` // column_b
}
