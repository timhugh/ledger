package main

import (
    "github.com/urfave/cli/v2"
    "fmt"
)

var registerCmd = cli.Command{
    Name: "register",
    Category: "report",
    Aliases: []string{"reg", "r"},
    Usage: "print all transactions",
    Action: func(c *cli.Context) error {
        l := loadLedger()
        for _, transaction := range l.Register() {
            fmt.Println(transaction)
        }
        return nil
    },
}

