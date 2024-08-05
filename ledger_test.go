package ledger_test

import (
    "github.com/timhugh/ledger"
    "testing"
)

func TestLedger_Balances(t *testing.T) {
    t.Run("single transaction", func(t *testing.T) {
        ledger := ledger.Ledger{
            Transactions: []ledger.Transaction{
                {
                    LineItems: []ledger.LineItem{
                        {Account: "a", Amount: 10},
                        {Account: "b", Amount: -10},
                    },
                },
            },
        }
        balances := ledger.Balances()
        if balances["a"] != 10 {
            t.Error("expected balance of a to be 10")
        }
        if balances["b"] != -10 {
            t.Error("expected balance of b to be -10")
        }
    })
    t.Run("multiple transactions", func(t *testing.T) {
        ledger := ledger.Ledger{
            Transactions: []ledger.Transaction{
                {
                    LineItems: []ledger.LineItem{
                        {Account: "a", Amount: 10},
                        {Account: "b", Amount: -10},
                    },
                },
                {
                    LineItems: []ledger.LineItem{
                        {Account: "a", Amount: 5},
                        {Account: "b", Amount: -5},
                    },
                },
            },
        }
        balances := ledger.Balances()
        if balances["a"] != 15 {
            t.Error("expected balance of a to be 15")
        }
        if balances["b"] != -15 {
            t.Error("expected balance of b to be -15")
        }
    })
}

func TestLedger_Register(t *testing.T) {
    t.Run("single transaction", func(t *testing.T) {
        ledger := ledger.Ledger{
            Transactions: []ledger.Transaction{
                {
                    LineItems: []ledger.LineItem{
                        {Account: "a", Amount: 10},
                        {Account: "b", Amount: -10},
                    },
                },
            },
        }
        transactions := ledger.Register()
        if len(transactions) != 1 {
            t.Error("expected 1 transaction")
        }
    })
    t.Run("multiple transactions", func(t *testing.T) {
        ledger := ledger.Ledger{
            Transactions: []ledger.Transaction{
                {
                    LineItems: []ledger.LineItem{
                        {Account: "a", Amount: 10},
                        {Account: "b", Amount: -10},
                    },
                },
                {
                    LineItems: []ledger.LineItem{
                        {Account: "a", Amount: 5},
                        {Account: "b", Amount: -5},
                    },
                },
            },
        }
        transactions := ledger.Register()
        if len(transactions) != 2 {
            t.Error("expected 2 transactions")
        }
    })
}
