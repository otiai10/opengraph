package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/otiai10/opengraph"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Action = func(ctx *cli.Context) error {
		rawurl := ctx.Args().First()
		if rawurl == "" {
			return fmt.Errorf("URL must be specified")
		}
		og, err := opengraph.Fetch(rawurl)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(og, "", "\t")
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", string(b))
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
