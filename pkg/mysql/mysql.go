package mysql

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Database struct {
	Conn *sql.DB
}

func Serve(user, password, host, dbname string) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to MySQL database successfully")
	return &Database{Conn: db}, nil
}

func (d *Database) Close() error {
	return d.Conn.Close()
}
