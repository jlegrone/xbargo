package main

import (
	"log"
	"os"
	"os/user"

	"github.com/jlegrone/xbargo"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	xbargo.NewPlugin().WithText("👋🌎").WithElements(
		xbargo.NewMenuItem("Greet").
			WithShell("say", os.Getenv("GREETING"), usr.Username).
			WithShortcut("G", xbargo.CommandKey),
	).Run()
}
