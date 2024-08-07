package main

import (
    "github.com/timhugh/ledger"
    "github.com/urfave/cli/v2"
    "os"
    "log"
)

func main() {
    if err := Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func Run(args []string) error {
    app := &cli.App{
        Name: "ledger",
        Commands: []*cli.Command{
            &registerCmd,
            BalancesCmd,
        },
    }
    return app.Run(args)
}

func loadLedger() ledger.Ledger {
    // TODO: placeholder data
    return ledger.Ledger{
        Transactions: []ledger.Transaction{
            {
                LineItems: []ledger.LineItem{
                    {Account: "Assets:Checking", Amount: 10},
                    {Account: "Equity:Opening balance", Amount: -10},
                },
            },
            {
                LineItems: []ledger.LineItem{
                    {Account: "Expenses:Hamburgers", Amount: 5},
                    {Account: "Assets:Checking", Amount: -5},
                },
            },
        },
    }
}

