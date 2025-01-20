package database

import "database/sql"

// Query executes the given SQL query with optional arguments and returns the resulting rows.
// It uses the underlying database connection to perform the query.
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

// Exec executes the given SQL query with the provided arguments and returns the result.
// It delegates the execution to the underlying database connection.
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
