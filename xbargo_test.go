package xbargo_test

import (
	"bytes"
	_ "embed"

	"github.com/jlegrone/xbargo"
)

func userHomeDir() string {
	// hardcoded for testing purposes, normally we'd use os.UserHomeDir()
	return "/tmp/xbargo_test"
}

// Example of how to run shell commands when menu items are clicked.
func ExamplePlugin_runShell() {
	xbargo.NewPlugin().WithText("🐌").WithElements(
		xbargo.NewMenuItem("🗣️ Say Hello").
			WithShell("say", "hello", "world").
			// Switch to "goodbye" when the alt/option key is pressed
			WithAlt(xbargo.NewMenuItem("🗣️ Say Goodbye").
				WithShell("say", "goodbye", "world")),
		xbargo.Separator{},
		xbargo.NewMenuItem("🔋 Battery Preferences").WithShell(
			"open", "-b", "com.apple.systempreferences",
			"/System/Library/PreferencePanes/Battery.prefPane",
		).WithShortcut("b", xbargo.ShiftKey),
		xbargo.Separator{},
		xbargo.NewMenuItem("🏠 Home Directory").WithSubMenu(
			xbargo.NewMenuItem("View Tree").WithAction(
				xbargo.NewShellAction(
					// Shell expansion isn't supported, so environment variables like $HOME
					// will need to be evaluated in Go before templating out the shell command.
					"tree", "-d", "-L", "1", userHomeDir(),
				).WithTerminal(),
			),
			xbargo.NewMenuItem("Copy Path").WithAction(
				xbargo.NewCopyAction(userHomeDir()),
			).WithShortcut("c", xbargo.CommandKey),
		),
		xbargo.Separator{},
		xbargo.NewMenuItem("ℹ️ Send Notification").WithShell(
			"osascript", "-e",
			`display notification "This is a notification" with title "Example" subtitle "Thanks for clicking!"`,
		).WithShortcut("n", xbargo.ControlKey, xbargo.OptionKey),
	).Run()
	// Output:
	// 🐌| refresh=false trim=false
	// ---
	// 🗣️ Say Hello| terminal=false shell="say" param1='hello' param2='world' refresh=false trim=false
	// 🗣️ Say Goodbye| terminal=false shell="say" param1='goodbye' param2='world' refresh=false trim=false alternate=true
	// ---
	// 🔋 Battery Preferences| key=shift+b terminal=false shell="open" param1='-b' param2='com.apple.systempreferences' param3='/System/Library/PreferencePanes/Battery.prefPane' refresh=false trim=false
	// ---
	// 🏠 Home Directory| refresh=false trim=false
	// --View Tree| terminal=true shell="tree" param1='-d' param2='-L' param3='1' param4='/tmp/xbargo_test' refresh=false trim=false
	// --Copy Path| key=CmdOrCtrl+c terminal=false shell="/bin/bash" param1='-c' param2='echo -n /tmp/xbargo_test | pbcopy' refresh=false trim=false
	// ---
	// ℹ️ Send Notification| key=ctrl+OptionOrAlt+n terminal=false shell="osascript" param1='-e' param2='display notification "This is a notification" with title "Example" subtitle "Thanks for clicking!"' refresh=false trim=false
}

// Example of how to include alternate items that replace their parent item when the Option key is pressed.
//
// Modeled after https://github.com/matryer/xbar-plugins/blob/main/Dev/Tutorial/alternate_options.sh
func ExamplePlugin_alternateOptions() {
	xbargo.NewPlugin().WithText("Alternate Options").WithElements(
		xbargo.NewMenuItem("Hello").WithAlt(xbargo.NewMenuItem("Option key is pressed")),
		xbargo.NewMenuItem("Another"),
	).Run()
	// Output:
	// Alternate Options| refresh=false trim=false
	// ---
	// Hello| refresh=false trim=false
	// Option key is pressed| refresh=false trim=false alternate=true
	// Another| refresh=false trim=false
}

// Example of how to include multiple levels of menu items.
//
// Modeled after https://github.com/matryer/xbar-plugins/blob/main/Dev/Tutorial/submenus.sh
func ExamplePlugin_submenus() {
	xbargo.NewPlugin().WithText("Submenu").WithElements(
		xbargo.NewMenuItem("Places").WithSubMenu(
			xbargo.NewMenuItem("London"),
			xbargo.NewMenuItem("Paris"),
			xbargo.NewMenuItem("Tokyo"),
		),
		xbargo.Separator{},
		xbargo.NewMenuItem("Fruit").WithSubMenu(
			xbargo.NewMenuItem("Apple"),
			xbargo.NewMenuItem("Orange"),
			xbargo.NewMenuItem("Melon").WithSubMenu(
				xbargo.NewMenuItem("Watermelon"),
				xbargo.NewMenuItem("Honeydew"),
			),
		),
	).Run()
	// Output:
	// Submenu| refresh=false trim=false
	// ---
	// Places| refresh=false trim=false
	// --London| refresh=false trim=false
	// --Paris| refresh=false trim=false
	// --Tokyo| refresh=false trim=false
	// ---
	// Fruit| refresh=false trim=false
	// --Apple| refresh=false trim=false
	// --Orange| refresh=false trim=false
	// --Melon| refresh=false trim=false
	// ----Watermelon| refresh=false trim=false
	// ----Honeydew| refresh=false trim=false
}

//go:embed internal/beaker.png
var beakerImage []byte

// Demonstrates embedding images in menu items. Image files
//
// Icon files may be in any of the formats supported by macOS.
// The recommended size for images in the statusbar and dropdown is 16x16 pixels.
//
// Modeled after https://github.com/matryer/xbar-plugins/blob/main/Dev/Tutorial/images.sh
func ExamplePlugin_imagesAndLinks() {
	xbargo.NewPlugin().
		WithIcon(bytes.NewReader(beakerImage)).
		WithElements(
			xbargo.NewMenuItem("Icon by Daniel Bruce"),
			xbargo.NewMenuItem("View Source").WithHref("https://iconscout.com/icon/lab-152"),
		).
		Run()
	// Output:
	// | templateImage=iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAhGVYSWZNTQAqAAAACAAFARIAAwAAAAEAAQAAARoABQAAAAEAAABKARsABQAAAAEAAABSASgAAwAAAAEAAgAAh2kABAAAAAEAAABaAAAAAAAAAEgAAAABAAAASAAAAAEAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAEKADAAQAAAABAAAAEAAAAADHbxzxAAAACXBIWXMAAAsTAAALEwEAmpwYAAABWWlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNi4wLjAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWRvYmUuY29tL3RpZmYvMS4wLyI+CiAgICAgICAgIDx0aWZmOk9yaWVudGF0aW9uPjE8L3RpZmY6T3JpZW50YXRpb24+CiAgICAgIDwvcmRmOkRlc2NyaXB0aW9uPgogICA8L3JkZjpSREY+CjwveDp4bXBtZXRhPgoZXuEHAAABU0lEQVQ4EXWTuy4FURSG93BcGhQkp3EeQKP3GKLTKHS8gkSiEY+hk5N4A61GVBpPoFCQ0Lgl+L49e80ZY86f/Hvd/rVmnzVzqvQXFeEPXIYbxcck8/fwtfhqejEo2UOsojbNidDUUef0SWIIr2EM0DcnQpODmXxODhvm4SO8naSzb86amgbdARa+S3WxUaUUftSaUt+AKK6Fg237rfR/NwaOKL3B2IG+ORGaOuqcsyW+wNr8Dr+KP8aKqW8hCtuI4slaBzhIfweK0NYRZ7yaOfxneAd34RVsD3sgXoAienIQE4+IbFjJ2fo4wNzAfehNjqGInmaSy3FZp1bBEmxEOZPSGVZNLDLfIkR7FHy6H8s0xIekVgxsjt+ylVMpnWDdwRP0ysIPaRVuGgC15zD3xoB1EpfwE3qTPlpTo1ZU0ay1QXjNIXSRsfEP/Bfo/8EhIvf8AlcrSrqa23krAAAAAElFTkSuQmCC refresh=false trim=false
	// ---
	// Icon by Daniel Bruce| refresh=false trim=false
	// View Source| href=https://iconscout.com/icon/lab-152 refresh=false trim=false
}
