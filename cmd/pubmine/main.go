package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/tenkoh/go-pubmine"
	"github.com/urfave/cli/v2"
)

const version = "v0.0.1"
const helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`

func run(c *cli.Context) error {
	args := c.Args()
	if !args.Present() {
		return cli.ShowAppHelp(c)
	}

	g, err := pubmine.NewGenerator(args.First(), int64(runtime.NumCPU()))
	if err != nil {
		return fmt.Errorf("failed to initialize a miner: %w", err)
	}

	fmt.Println("process started. this may take a long time.")
	fmt.Print("mining...")
	type result struct {
		keyPair *pubmine.KeyPair
		err     error
	}
	ch := make(chan result)
	go func() {
		kp, err := g.Mine(c.Context)
		ch <- result{kp, err}
	}()

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Print(".")
		}
	}()

	r := <-ch
	if r.err != nil {
		return fmt.Errorf("failed to mine a keypair: %w", r.err)
	}
	ticker.Stop()
	fmt.Print("\n\n")

	fmt.Printf("Public key:\n%s\n", r.keyPair.Public)
	fmt.Printf("Private key:\n%s\n", r.keyPair.Private)
	return nil
}

func main() {
	app := &cli.App{
		Name:                  "pubmine",
		Usage:                 "mine a keypair with a pretty public key",
		Version:               version,
		Action:                run,
		ArgsUsage:             "prefix",
		Commands:              nil,
		CustomAppHelpTemplate: helpTemplate,
	}
	var exit int
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		exit = 1
	}
	os.Exit(exit)
}
