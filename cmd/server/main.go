package main

import (
	"github.com/gin-gonic/gin"
	"github.com/timhugh/ledger/cmd/server/controllers"
	"github.com/timhugh/ledger/db/sqlite"
	"log"
)

func main() {
	repo, err := sqlite.Open("development.db")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", controllers.Ping())
	router.GET("/journals/:journal_id", controllers.GetJournal(repo))

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
