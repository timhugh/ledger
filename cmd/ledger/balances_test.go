package main_test

import (
    "github.com/timhugh/ledger"
    "github.com/timhugh/ledger/cmd/ledger"
)

func ExampleBalances() {
    ledger := ledger.Ledger{
        Transactions: []ledger.Transaction{
            {
                LineItems: []ledger.LineItem{
                    { Account: "Assets:Checking", Amount: 1000 },
                    { Account: "Assets:Savings", Amount: 10000 },
                    { Account: "Expenses:Groceries", Amount: -200 },
                    { Account: "Expenses:Rent", Amount: -1000 },
                },
            },
        },
    }
    main.Balances(ledger)

    // Output:
    // Account               Balance
    // Assets:Checking       1000
    // Assets:Savings        10000
    // Expenses:Groceries    -200
    // Expenses:Rent         -1000
}
