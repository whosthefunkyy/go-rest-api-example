package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

var DB *sql.DB

func InitDB() {
    var err error
    connStr := "host=localhost user=artem password=password dbname=goapp_db sslmode=disable"
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("DB open error:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("DB ping failed:", err)
    }

    log.Println("Database connected")
}
