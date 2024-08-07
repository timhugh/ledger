package main

import (
    "github.com/timhugh/ledger"
    "github.com/urfave/cli/v2"
    "fmt"
    "os"
    "text/tabwriter"
)

var registerCmd = cli.Command{
    Name: "register",
    Category: "report",
    Aliases: []string{"reg", "r"},
    Usage: "print all transactions",
    Action: func(c *cli.Context) error {
        l := loadLedger()
        Register(l)
        return nil
    },
}

func Register(ledger ledger.Ledger) {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
    fmt.Fprintln(w, "Date\tAccount\tAmount")
    for _, transaction := range ledger.Register() {
        for _, lineItem := range transaction.LineItems {
            fmt.Fprintf(w, "%s\t%s\t%d\n", lineItem.Date, lineItem.Account, lineItem.Amount)
        }
    }
    w.Flush()
}
