package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/alexey-medvedchikov/go-heapview/cmd/heapview/dumpcmd"
	"github.com/alexey-medvedchikov/go-heapview/cmd/heapview/ownedcmd"
	"github.com/alexey-medvedchikov/go-heapview/internal/profile"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalf("%+v", err)
	}
}

func newApp() *cli.App {
	var cpuFinish func()
	var heapFinish func()

	return &cli.App{
		Name: "heapview",
		Commands: cli.Commands{
			dumpcmd.Command(),
			ownedcmd.Command(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cpu-profile",
				Usage: "File to save CPU profile to",
				Action: func(_ *cli.Context, s string) error {
					var err error
					cpuFinish, err = profile.CPU(s)
					return err
				},
			},
			&cli.StringFlag{
				Name:  "heap-profile",
				Usage: "File to save heap profile to",
				Action: func(_ *cli.Context, s string) error {
					var err error
					heapFinish, err = profile.Heap(s)
					return err
				},
			},
		},
		After: func(_ *cli.Context) error {
			if cpuFinish != nil {
				cpuFinish()
			}

			if heapFinish != nil {
				heapFinish()
			}

			return nil
		},
	}
}
