package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coderavels/airtablegolangcli/client"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "airtablecli",
		Usage: "make a cli interface for airtable",
		Action: func(c *cli.Context) error {
			apiToken := c.Args().Get(0)

			airtableClient := client.NewAirtableClient(apiToken)
			records, _, err := airtableClient.ListRecords("app123", "tbl123", client.OptionalParams{
				PageSize: 20,
			})
			if err != nil {
				return err
			}
			fmt.Println(records)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
