package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	scraper_db_database = os.Getenv("SCRAPER_DB_DATABASE")
	scraper_db_password = os.Getenv("SCRAPER_DB_PASSWORD")
	scraper_db_username = os.Getenv("SCRAPER_DB_USERNAME")
	scraper_db_port     = os.Getenv("SCRAPER_DB_PORT")
	scraper_db_host     = os.Getenv("SCRAPER_DB_HOST")
)
var (
	main_db_database = os.Getenv("DB_DATABASE")
	main_db_password = os.Getenv("DB_PASSWORD")
	main_db_username = os.Getenv("DB_USERNAME")
	main_db_port     = os.Getenv("DB_PORT")
	main_db_host     = os.Getenv("DB_HOST")
)

func main() {

	sdb := connectScraperDB()
	mdb := connectMainDB()
	defer sdb.Close()
	defer mdb.Close()

	if err := syncQuestionsTable(sdb.DB, mdb.DB); err != nil {
		log.Fatal(err)
	}

}

func connectScraperDB() *sqlx.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", scraper_db_username, scraper_db_password, scraper_db_host, scraper_db_port, scraper_db_database)
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectMainDB() *sqlx.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", main_db_username, main_db_password, main_db_host, main_db_port, main_db_database)
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func syncQuestionsTable(sdb *sql.DB, mdb *sql.DB) error {
	query := `
		SELECT q.question_key, q.question_number, i.has_image, i.file_name, i.extracted_text
		FROM questions AS q
		JOIN images AS i ON q.id = i.question_id
	`
	rows, err := sdb.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var questionKey string
		var questionNumber string
		var hasImage bool
		var fileName string
		var extractedText string

		err := rows.Scan(&questionKey, &questionNumber, &hasImage, &fileName, &extractedText)
		if err != nil {
			return err
		}

		_, err = mdb.Exec("INSERT INTO questions (question_key, question_number, has_image, file_name, question) VALUES ($1, $2, $3, $4, $5)", questionKey, questionNumber, hasImage, fileName, extractedText)
		if err != nil {
			return err
		}
	}

	return nil
}
