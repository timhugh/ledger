package main

import (
    "github.com/timhugh/ledger"
    "github.com/urfave/cli/v2"
    "fmt"
    "os"
    "text/tabwriter"
)


var BalancesCmd =  &cli.Command{
    Name: "balances",
    Aliases: []string{"bal", "b"},
    Category: "report",
    Usage: "print all account balances",
    Action: func(c *cli.Context) error {
        l := loadLedger()
        return Balances(l)
    },
}

func Balances(ledger ledger.Ledger) error {
    balances := ledger.Balances()
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
    fmt.Fprintln(w, "Account\tBalance")
    for account, balance := range balances {
        fmt.Fprintf(w, "%s\t%d\n", account, balance)
    }
    w.Flush()
    return nil
}
