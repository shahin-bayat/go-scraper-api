package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
	c "github.com/shahin-bayat/scraper-api/internal/config"
)

var (
	pgDatabase = c.PostgresConf.PgDatabase
	pgPassword = c.PostgresConf.PgPassword
	pgUser     = c.PostgresConf.PgUser
	pgPort     = c.PostgresConf.PgPort
	pgHost     = c.PostgresConf.PgHost
)

type question struct {
	ID       uint   `db:"id"`
	Question string `db:"extracted_text"`
}

type answer struct {
	ID   uint   `db:"id"`
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
	var wg sync.WaitGroup

	connStrLocal := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDatabase,
	)

	//connectStrDev := c.AppConf.PostgresDevURL
	_ = c.AppConf.PostgresProdURL

	deeplAPIKey := c.AppConf.DeeplAPIKey
	client := http.Client{}

	start := time.Now()
	translateQuestions(&client, deeplAPIKey, connStrLocal, "TR", &wg)
	//translateQuestions(&client, deeplAPIKey, connectStrDev, "TR", &wg)
	// translateQuestions(&client, deeplAPIKey, connStrProd, "TR", &wg)

	translateAnswers(&client, deeplAPIKey, connStrLocal, "TR", &wg)
	//translateAnswers(&client, deeplAPIKey, connectStrDev, "TR", &wg)
	// translateAnswers(&client, deeplAPIKey, connStrProd, "TR", &wg)

	wg.Wait()
	fmt.Println("Time elapsed:", time.Since(start))
}
