# gormcheck

`gormcheck` is a static analysis tool which can find wrong gorm using logic.
Current state is alphe version. Only check pipeed function( make pipe is originally not wrong but my working team stop to use to make review easier).


## Install

You can get `gormcheck` by `go get` command.

```bash
$ go get -u github.com/gosagawa/gormcheck
```

## How to use

`gormcheck` run with `go vet` as below when Go is 1.12 and higher.

```bash
$ go vet -vettool=$(which gormcheck) your_check_path
```

When Go is lower than 1.12, just run `gormcheck` command with the package name (import path).

```bash
$ gorm check your_check_path
```

## Example

This is expected function.

```go
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
```

This code is piped. It is grammatically correct, but want stop to use to make review easier.

```go
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
```

Allow pipe function except db. Now only check paramater name is "db" or not...

```
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
```


