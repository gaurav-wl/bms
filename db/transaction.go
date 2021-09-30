package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(tx *sqlx.Tx) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *sqlx.DB, fn TxFn) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// QuestionToDollar takes a string and converts ? to $ parameterized variables.
func QuestionToDollar(str string) string {
	var (
		newStr     []byte
		paramCount = 1
	)
	for _, r := range str {
		if r == '?' {
			newStr = append(newStr, []byte(fmt.Sprintf("$%v", paramCount))...)
			paramCount++
		} else {
			newStr = append(newStr, byte(r))
		}
	}
	return string(newStr)
}
