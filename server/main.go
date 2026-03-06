package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Row map[string]interface{}

func main() {
	// 1. Initialize DB and Seed if necessary
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := seedDB(db); err != nil {
		log.Fatal("Failed to seed DB:", err)
	}

	// 2. Setup HTTP handlers
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/api/rows", func(w http.ResponseWriter, r *http.Request) {
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		start, _ := strconv.Atoi(startStr)
		end, _ := strconv.Atoi(endStr)
		if end == 0 {
			end = 100 // Default limit if not specified
		}
		limit := end - start

		rows, err := db.Query("SELECT * FROM users LIMIT ? OFFSET ?", limit, start)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		result := []Row{}

		for rows.Next() {
			// Create a slice of interface{} to hold pointer to values
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			// Scan the result into the column pointers
			if err := rows.Scan(columnPointers...); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Create our map, and retrieve the value for each column
			m := make(Row)
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				m[colName] = *val
			}
			result = append(result, m)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rows":    result,
			"lastRow": getTotalCount(db),
		})
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTotalCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count
}

func seedDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, age INTEGER)`)
	if err != nil {
		return err
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		log.Println("Seeding database...")
		tx, _ := db.Begin()
		stmt, _ := tx.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
		for i := 0; i < 1000; i++ {
			stmt.Exec(fmt.Sprintf("User %d", i), fmt.Sprintf("user%d@example.com", i), 20+(i%50))
		}
		tx.Commit()
		log.Println("Seeded 1000 rows.")
	}
	return nil
}
