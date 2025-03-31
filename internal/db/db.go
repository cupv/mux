package db

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

var DB *sql.DB

func ConnectDB() {
    dsn := "user:password@tcp(localhost:3306)/vocab_cards?parseTime=true"
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Cannot connect to database:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Database not responding:", err)
    }
    log.Println("Connected to MySQL successfully")
}
