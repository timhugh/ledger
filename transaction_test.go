package ledger_test

import (
    "github.com/timhugh/ledger"
    "testing"
)

func TestTransaction_Cleared(t *testing.T) {
    t.Run("all cleared", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Status: ledger.Cleared},
                {Status: ledger.Cleared},
            },
        }
        if !transaction.Cleared() {
            t.Error("expected transaction to be cleared")
        }
    })
    t.Run("not all cleared", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Status: ledger.Cleared},
                {Status: ledger.Pending},
            },
        }
        if transaction.Cleared() {
            t.Error("expected transaction to not be cleared")
        }
    })
}

func TestTransaction_Date(t *testing.T) {
    t.Run("single line item", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Date: "2020-01-01"},
            },
        }
        if transaction.Date() != "2020-01-01" {
            t.Error("expected transaction date to be 2020-01-01")
        }
    })
    t.Run("multiple line items", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Date: "2020-01-01"},
                {Date: "2020-01-02"},
            },
        }
        if transaction.Date() != "2020-01-01" {
            t.Error("expected transaction date to be 2020-01-01")
        }
    })
    t.Run("multiple line items any order", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Date: "2020-01-02"},
                {Date: "2020-01-01"},
            },
        }
        if transaction.Date() != "2020-01-01" {
            t.Error("expected transaction date to be 2020-01-01")
        }
    })
}

func TestTransaction_Valid(t *testing.T) {
    t.Run("valid", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Amount: 10},
                {Amount: -10},
            },
        }
        if !transaction.Valid() {
            t.Error("expected transaction to be valid")
        }
    })
    t.Run("not valid", func(t *testing.T) {
        transaction := ledger.Transaction{
            LineItems: []ledger.LineItem{
                {Amount: 10},
                {Amount: -5},
            },
        }
        if transaction.Valid() {
            t.Error("expected transaction to not be valid")
        }
    })
}
