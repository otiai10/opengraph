package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/otiai10/opengraph/v2"
)

func main() {
	flagset := flag.CommandLine
	flagset.Usage = func() {
		fmt.Println("Fetch URL and extract OpenGraph meta informations.")
	}
	abs := flagset.Bool("A", false, "populate relative URLs to absolute URLs")
	flagset.Parse(os.Args[1:])
	if err := run(flagset.Args(), *abs); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

func run(args []string, absolute bool) error {
	if len(args) == 0 {
		return fmt.Errorf("URL must be specified")
	}
	rawurl := args[0]
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	og := opengraph.New(u.String())
	if err := og.Fetch(); err != nil {
		return err
	}
	if absolute {
		if err := og.ToAbs(); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(og, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", string(b))
	return nil
}
