package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/timhugh/ledger"
)

func GetJournal(repo ledger.JournalGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		journalID := c.Param("journal_id")
		journal, err := repo.GetJournal(journalID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, journal)
	}
}
