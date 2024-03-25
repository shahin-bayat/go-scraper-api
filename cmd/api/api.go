package main

import (
	"fmt"

	"github.com/shahin-bayat/scraper-api/internal/ecosystem"
	"github.com/shahin-bayat/scraper-api/internal/server"
)

func main() {
	eco := ecosystem.Require()

	server := server.Create(eco)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
