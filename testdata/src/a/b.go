package a

// This file is for check it work on multiple files.

// using pipe
func fb1() (*User, error) {

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

// have two First
func fb2() (*User, error) { // want "have two more select function like Find, First"

	db, err := getConnection()
	if err != nil {
		return nil, err
	}

	u := User{}
	db = db.Where("column_a = xxx")
	db = db.Where("column_b = xxx")
	db = db.First(&u)
	db = db.First(&u)

	if db.RecordNotFound() {
		return nil, nil
	}
	return &u, db.Error
}
