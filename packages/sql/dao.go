package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

func dbConfig() string {
	host, presentHost := os.LookupEnv("DB_HOST")
	if !presentHost {
		log.Fatalf("DB_HOST env var is mandatory")
	}
	port, presentPort := os.LookupEnv("DB_PORT")
	if !presentPort {
		log.Fatalf("DB_PORT env var is mandatory")
	}
	user, presentUser := os.LookupEnv("DB_USER")
	if !presentUser {
		log.Fatalf("DB_USER env var is mandatory")
	}
	password, presentPass := os.LookupEnv("DB_PASS")
	if !presentPass {
		log.Fatalf("DB_PASS env var is mandatory")
	}
	dbname, presentDbname := os.LookupEnv("DB_NAME")
	if !presentDbname {
		log.Fatalf("DB_NAME env var is mandatory")
	}
	portNumber, errP := strconv.Atoi(port)
	if errP != nil {
		panic(errP)
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, portNumber, user, password, dbname)
}

func Connect() *sql.DB {
	db, err := sql.Open("postgres", dbConfig())
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
func Persist(db *sql.DB, fecha string, from string, subject string) string {
	var lastInsertId string
	err := db.QueryRow("INSERT INTO challenge(fecha, from_, subject) VALUES($1,$2,$3) returning m_id;", fecha, from, subject).Scan(&lastInsertId)
	if err != nil {
		panic(err)
	}
	return lastInsertId
}
