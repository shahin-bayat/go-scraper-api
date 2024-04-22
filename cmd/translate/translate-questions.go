package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"sync"
)

func translateQuestions(client *http.Client, deeplAPIKey, connStr, lang string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	var questions []question
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return
	}

	err = db.Select(
		&questions, `
			SELECT q.id , i.extracted_text
			FROM questions AS q
			JOIN images AS i ON i.question_id = q.id
	`,
	)

	if err != nil {
		fmt.Println("failed to get local questions:", err)
		return
	}

	for _, q := range questions {
		fmt.Printf("Translating question: %s\n", q.Question)
		text := q.Question
		payloadData := map[string]interface{}{
			"text":        []string{text},
			"target_lang": lang,
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

		db.MustExec(
			`
			INSERT INTO translations (refer_id, type, lang, translation)
			VALUES ($1, $2, $3, $4)
		`, q.ID, "question", "en", t,
		)
	}
}
