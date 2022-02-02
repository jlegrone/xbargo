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
		xbargo.Separator{},
		xbargo.NewMenuItem("Statuses").WithSubMenu(
			xbargo.NewMenuItem("Available").WithIcon(xbargo.IconStatusAvailable),
			xbargo.NewMenuItem("None").WithIcon(xbargo.IconStatusNone),
			xbargo.NewMenuItem("Partially").WithIcon(xbargo.IconStatusPartially),
			xbargo.NewMenuItem("Unavailable").WithIcon(xbargo.IconStatusUnavailable),
		),
	).Run()
}
