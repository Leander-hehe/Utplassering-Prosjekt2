package data

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"
)

// Thing represents a row in the "things" table
type Thing struct {
	ID          string
	Name        string
	Description string
}

func logError(data any, err error) {
	fmt.Println(data, err)
}

func Setup(db *sql.DB) (error, int) {
	// Ensure the connection is valid
	if err := db.Ping(); err != nil {
		logError("Database connection error:", err)
		return err, http.StatusInternalServerError
	}

	// Create table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS things (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL
	);`
	if _, err := db.Exec(createTable); err != nil {
		logError("Failed to create table:", err)
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func Insert(db *sql.DB, thing Thing) (error, int, int) {

	result, err := db.Exec("INSERT INTO things (name, description) VALUES (?, ?)", thing.Name, thing.Description)
	if err != nil {
		logError("Failed to insert data:", err)
		return err, http.StatusBadRequest, 0
	}

	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		logError("Failed to get new id:", err)
		return err, http.StatusInternalServerError, 0
	}
	return nil, http.StatusOK, int(id)
}

func Select(db *sql.DB, id string) (*Thing, error, int) {

	// Retrieve data
	rows, err := db.Query("SELECT id, name, description FROM things")
	if err != nil {
		logError("Failed to query data:", err)
		return nil, err, http.StatusInternalServerError
	}
	defer rows.Close()

	// Iterate over rows
	var things []Thing
	for rows.Next() {
		var t Thing
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			logError("Failed to scan row:", err)
			return nil, err, http.StatusInternalServerError
		}
		things = append(things, t)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		logError("Row iteration error:", err)
		return nil, err, http.StatusInternalServerError
	}

	// Return results
	for _, t := range things {
		if t.ID == id {
			return &t, nil, http.StatusOK
		}

	}

	return nil, errors.New("no item found with ID: " + id), http.StatusNotFound

}

func SelectAll(db *sql.DB) (*[]Thing, error, int) {
	// Retrieve data
	rows, err := db.Query("SELECT id, name, description FROM things")
	if err != nil {
		logError("Failed to query data:", err)
		return nil, err, http.StatusInternalServerError
	}
	defer rows.Close()

	// Iterate over rows
	var things []Thing
	for rows.Next() {
		var t Thing
		if err := rows.Scan(&t.ID, &t.Name, &t.Description); err != nil {
			logError("Failed to scan row:", err)
			return nil, err, http.StatusInternalServerError
		}
		things = append(things, t)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		logError("Row iteration error:", err)
		return nil, err, http.StatusInternalServerError
	}

	return &things, nil, http.StatusOK

}

func Delete(db *sql.DB, id string) (error, int) {

	_, err := db.Exec("DELETE FROM things WHERE id = " + id)
	if err != nil {
		logError("Failed to delete data:", err)
		return err, http.StatusNotFound
	}

	return nil, http.StatusOK
}

func Update(db *sql.DB, id string, thing Thing) (error, int) {
	// "UPDATE things SET (name, description) WHERE id = "+id, thing.Name, thing.Description
	_, err := db.Exec("UPDATE things SET name = ?, description = ? WHERE id = ?", thing.Name, thing.Description, id)
	if err != nil {
		logError("Failed to update data:", err)
		return err, http.StatusNotFound
	}
	return nil, http.StatusOK
}
