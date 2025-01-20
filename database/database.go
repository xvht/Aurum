package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initTables() {
	// Have we already initialized the database?
	rows, err := db.Query("SELECT * FROM init")
	if err == nil {
		rows.Close()
		log.Println("Database: init table found")
		return
	}

	instructions := generateSetupInstructions()

	// sort tableMaps by key to ensure order of execution
	// 0 -> users
	// 1 -> invites
	// 2 -> init
	sortedInstructions := make([]map[string][]string, len(instructions))
	for key, tableMap := range instructions {
		sortedInstructions[key] = tableMap
	}

	for _, tableMap := range sortedInstructions {
		for tableName, queries := range tableMap {
			rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
			if err == nil {
				rows.Close()
				log.Printf("Instruct: find -> %s: found", tableName)
				continue
			} else {
				log.Printf("Instruct: find -> %s: not found", tableName)
			}

			log.Printf("Instruct: create -> %s", tableName)

			for indx, query := range queries {
				_, err := db.Exec(query)
				if err != nil {
					log.Fatalf("Instruct: cannot execute query: %v", err)
				}
				trimmedQuery := query
				if idx := strings.Index(query, "("); idx != -1 {
					trimmedQuery = query[:idx]
				}

				log.Printf("Instruct: exec %d -> %s: %s", indx, tableName, trimmedQuery)
			}
		}
	}
}

func Connect(connStr string) {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database: ping fail: %v", err)
	}

	log.Println("Database: connected to postgresql")
	initTables()
}
