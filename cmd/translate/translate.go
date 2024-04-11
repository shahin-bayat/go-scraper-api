package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	pgDatabase = os.Getenv("PG_DATABASE")
	pgPassword = os.Getenv("PG_PASSWORD")
	pgUser     = os.Getenv("PG_USER")
	pgPort     = os.Getenv("PG_PORT")
	pgHost     = os.Getenv("PG_HOST")
)

type question struct {
	ID       int    `db:"id"`
	Question string `db:"extracted_text"`
}

type answer struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
}

type translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type translationResponse struct {
	Translations []translation
}

func main() {

	connStrLocal := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDatabase)
	connectStrDev := os.Getenv("PG_DEV_URL")
	connStrProd := os.Getenv("PG_PROD_URL")

	deeplAPIKey := os.Getenv("DEEPL_API_KEY")
	if deeplAPIKey == "" {
		fmt.Println("DEEPL_API_KEY is not set")
		return
	}

	client := http.Client{}

	translateQuestions(&client, deeplAPIKey, connStrLocal)
	translateQuestions(&client, deeplAPIKey, connectStrDev)
	translateQuestions(&client, deeplAPIKey, connStrProd)
	translateAnswers(&client, deeplAPIKey, connStrLocal)
	translateAnswers(&client, deeplAPIKey, connectStrDev)
	translateAnswers(&client, deeplAPIKey, connStrProd)
}

func translateQuestions(client *http.Client, deeplAPIKey, connStr string) {
	questions := []question{}
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return
	}

	err = db.Select(&questions, `
			SELECT q.id , i.extracted_text 
			FROM questions AS q
			JOIN images AS i ON i.question_id = q.id
	`)

	if err != nil {
		fmt.Println("failed to get local questions:", err)
		return
	}

	for _, q := range questions {
		text := q.Question
		payloadData := map[string]interface{}{
			"text":        []string{text},
			"target_lang": "EN",
		}

		payloadBytes, err := json.Marshal(payloadData)
		if err != nil {
			fmt.Println("Error encoding payload:", err)
			return
		}
		payload := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", payload)
		if err != nil {
			fmt.Println("failed to create request:", err)
			return
		}
		req.Header.Set("Authorization", "DeepL-Auth-Key "+deeplAPIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("failed to send request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body:", err)
			return
		}

		translationResponse := translationResponse{}
		err = json.Unmarshal(body, &translationResponse)
		if err != nil {
			fmt.Println("failed to unmarshal response body:", err)
			return
		}

		t := translationResponse.Translations[0].Text

		db.MustExec(`
			INSERT INTO translations (refer_id, type, lang, translation)
			VALUES ($1, $2, $3, $4)
		`, q.ID, "question", "en", t)
	}
}

func translateAnswers(client *http.Client, deeplAPIKey, connStr string) {
	answers := []answer{}
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return
	}

	err = db.Select(&answers, `
			SELECT id, text 
			FROM answers
	`)

	if err != nil {
		fmt.Println("failed to get local questions:", err)
		return
	}

	for _, a := range answers {
		text := a.Text
		payloadData := map[string]interface{}{
			"text":        []string{text},
			"target_lang": "EN",
		}

		payloadBytes, err := json.Marshal(payloadData)
		if err != nil {
			fmt.Println("Error encoding payload:", err)
			return
		}
		payload := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", payload)
		if err != nil {
			fmt.Println("failed to create request:", err)
			return
		}
		req.Header.Set("Authorization", "DeepL-Auth-Key "+deeplAPIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("failed to send request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body:", err)
			return
		}

		translationResponse := translationResponse{}
		err = json.Unmarshal(body, &translationResponse)
		if err != nil {
			fmt.Println("failed to unmarshal response body:", err)
			return
		}

		t := translationResponse.Translations[0].Text

		db.MustExec(`
			INSERT INTO translations (refer_id, type, lang, translation)
			VALUES ($1, $2, $3, $4)
		`, a.ID, "answer", "en", t)
	}
}
