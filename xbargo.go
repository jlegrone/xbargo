// Package xbargo provides helpers for rendering xbar (https://github.com/matryer/xbar)
// plugins with Go.
package xbargo

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	_ XbarElement = Separator{}
	_ XbarElement = &MenuItem{}
	_ Action      = HrefAction{}
	_ Action      = ShellAction{}
)

var (
	//go:embed assets/Status_Available.png
	statusAvailableBytes []byte
	//go:embed assets/Status_None.png
	statusNoneBytes []byte
	//go:embed assets/Status_Partially.png
	statusPartiallyBytes []byte
	//go:embed assets/Status_Unavailable.png
	statusUnavailableBytes []byte
	// IconStatusAvailable is a small green indicator, similar to iChat’s available image.
	IconStatusAvailable = bytes.NewReader(statusAvailableBytes)
	// IconStatusNone is a small clear indicator.
	IconStatusNone = bytes.NewReader(statusNoneBytes)
	// IconStatusPartially is a small yellow indicator, similar to iChat’s idle image.
	IconStatusPartially = bytes.NewReader(statusPartiallyBytes)
	// IconStatusUnavailable is a small red indicator, similar to iChat’s unavailable image.
	IconStatusUnavailable = bytes.NewReader(statusUnavailableBytes)
)

// An XbarElement may be either a MenuItem or Separator.
type XbarElement interface {
	renderSelf() string
	renderAlt() string
	children() []XbarElement
}

// A Separator can be used to create visually distinct groups of related
// menu items.
type Separator struct{}

func (Separator) renderSelf() string {
	return "---"
}

func (Separator) renderAlt() string {
	return ""
}

func (Separator) children() []XbarElement {
	return nil
}

type Style struct {
	// Setting MaxLength will truncate text to the specified number of characters.
	//
	// A … will be added to any truncated strings, as well as a tooltip displaying
	// the full string.
	MaxLength uint
	// Change the Title color, e.g. "red" or "#ff0000"
	Color string
	// An Icon used in plugin titles should have IconImageTemplate enabled.
	//
	// A template image discards color information and uses a mask to produce the
	// appearance you see onscreen. Template images automatically adapt to the user’s
	// appearance settings, so they look good on both dark and light menu bars, and
	// when your menu bar extra is selected.
	IconImageTemplate bool
}

// An Action may be either an HrefAction or ShellAction.
type Action interface {
	actionMarker()
}

// ShellAction runs a shell command on click.
type ShellAction struct {
	Command      string
	Args         []string
	OpenTerminal bool
}

// actionMarker implements Action.
func (ShellAction) actionMarker() {}

// WithTerminal opens the shell command output in a new terminal window.
func (sa ShellAction) WithTerminal() ShellAction {
	sa.OpenTerminal = true
	return sa
}

func NewShellAction(command string, args ...string) ShellAction {
	return ShellAction{
		Command:      command,
		Args:         args,
		OpenTerminal: false,
	}
}

// NewCopyAction copies the given text to the user's clipboard.
func NewCopyAction(text string) ShellAction {
	return ShellAction{
		Command:      "/bin/bash",
		Args:         []string{"-c", fmt.Sprintf("echo -n %s | pbcopy", text)},
		OpenTerminal: false,
	}
}

// HrefAction opens a URI on click.
type HrefAction struct {
	URI string
}

// actionMarker implements Action.
func (HrefAction) actionMarker() {}

func NewHrefAction(uri string) HrefAction {
	return HrefAction{
		URI: uri,
	}
}

// A MenuItem may be configured to initiate an action, toggle a state on or off,
// or display a submenu of additional menu items when selected or in response to
// an associated keyboard shortcut.
type MenuItem struct {
	// The text of the menu item. May be empty if Icon is set.
	Title string
	// An Icon can be used to help people recognize menu items and associate them with
	// content.
	//
	// Even when using an icon you think people will recognize, it’s best to reinforce
	// the icon’s meaning with a textual title.
	//
	// The image format can be any of the formats supported by macOS.
	// The recommended size for images in the statusbar and dropdown is 16x16 pixels.
	//
	// More information:
	// https://developer.apple.com/design/human-interface-guidelines/macos/menus/menu-anatomy/#using-icons-in-menus
	Icon io.Reader
	// Style configures the Title and/or Icon style.
	Style Style
	// Shortcut sets a keyboard shortcut for the item.
	//
	// Use + to create combinations, e.g. "shift+k". Example options:
	// - CmdOrCtrl
	// - OptionOrAlt
	// - shift
	// - ctrl
	// - super
	// - tab
	// - plus
	// - return
	// - escape
	// - f12
	// - up
	// - down
	// - space
	Shortcut string
	// Action determines what the MenuItem will do on click or when the Shortcut
	// keys are pressed.
	//
	// MenuItems that have no action configured will appear disabled.
	Action Action
	// Make the item refresh the plugin it belongs to.
	// If the item runs a script, refresh is performed after the script finishes.
	Refresh bool
	// Replace with an alternate item when the Option key is pressed in the dropdown.
	// A menu item is dynamic when its behavior changes with the addition of a modifier key (Control, Option, Shift, or Command). For example, the Minimize item in the Window menu changes to Minimize All when pressing the Option key.
	Alt *MenuItem
	// Items to nest in a submenu under the current item.
	SubMenu []*MenuItem
}

func NewMenuItem(title string) *MenuItem {
	return &MenuItem{Title: title}
}

func (m *MenuItem) WithStyle(style Style) *MenuItem {
	m.Style = style
	return m
}

func (m *MenuItem) WithRefresh() *MenuItem {
	m.Refresh = true
	return m
}

func (m *MenuItem) WithAction(action Action) *MenuItem {
	m.Action = action
	return m
}

func (m *MenuItem) WithHref(uri string) *MenuItem {
	m.Action = NewHrefAction(uri)
	return m
}

func (m *MenuItem) WithShell(command string, args ...string) *MenuItem {
	m.Action = NewShellAction(command, args...)
	return m
}

func (m *MenuItem) WithIcon(icon io.Reader) *MenuItem {
	m.Icon = icon
	return m
}

func (m *MenuItem) WithAlt(item *MenuItem) *MenuItem {
	m.Alt = item
	return m
}

// A ModifierKey may be used to assign a shortcut to a MenuItem's action.
type ModifierKey string

const (
	CommandKey = ModifierKey("CmdOrCtrl")
	OptionKey  = ModifierKey("OptionOrAlt")
	ControlKey = ModifierKey("ctrl")
	ShiftKey   = ModifierKey("shift")
)

func (m *MenuItem) WithShortcut(key string, modifiers ...ModifierKey) *MenuItem {
	var modStrings []string
	for _, m := range modifiers {
		modStrings = append(modStrings, string(m))
	}
	m.Shortcut = strings.Join(append(modStrings, key), "+")
	return m
}

func (m *MenuItem) WithSubMenu(items ...*MenuItem) *MenuItem {
	m.SubMenu = append(m.SubMenu, items...)
	return m
}

func (m *MenuItem) renderSelf() string {
	parts := []string{
		fmt.Sprintf("%s|", m.Title),
	}
	if m.Shortcut != "" {
		parts = append(parts, fmt.Sprintf("key=%s", m.Shortcut))
	}
	if m.Style.MaxLength > 0 {
		parts = append(parts, fmt.Sprintf("length=%d", m.Style.MaxLength))
	}
	if m.Style.Color != "" {
		parts = append(parts, fmt.Sprintf("color=%s", m.Style.Color))
	}
	if m.Action != nil {
		switch action := m.Action.(type) {
		case HrefAction:
			parts = append(parts, fmt.Sprintf("href=%s", action.URI))
		case ShellAction:
			part := fmt.Sprintf("terminal=%t shell=%q", action.OpenTerminal, action.Command)
			for i, arg := range action.Args {
				part = fmt.Sprintf("%s param%d='%s'", part, i+1, arg)
			}
			parts = append(parts, part)
		}
	}
	if m.Icon != nil {
		b, err := io.ReadAll(m.Icon)
		if err != nil {
			panic(err)
		}
		imageType := "image"
		if m.Style.IconImageTemplate {
			imageType = "templateImage"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", imageType, base64.StdEncoding.EncodeToString(b)))
	}
	parts = append(parts,
		fmt.Sprintf("refresh=%t", m.Refresh),
		"trim=false",
	)

	return strings.Join(parts, " ")
}

func (m *MenuItem) renderAlt() string {
	if m.Alt == nil {
		return ""
	}
	return m.Alt.renderSelf()
}

func (m *MenuItem) children() []XbarElement {
	var children []XbarElement
	for _, menuItem := range m.SubMenu {
		children = append(children, menuItem)
	}
	return children
}

// A Plugin is used to create an executable Go program that prints lines of
// text which will be converted into a macOS menu by xbar.
//
// More information on creating xbar plugins is available here:
// https://github.com/matryer/xbar-plugins/blob/main/CONTRIBUTING.md
//
// Plugins should follow the Apple Human Interface Guidelines:
// https://developer.apple.com/design/human-interface-guidelines/macos/menus/menu-anatomy/
type Plugin struct {
	Title    *MenuItem
	Elements []XbarElement
}

func NewPlugin() *Plugin {
	return &Plugin{
		Title: NewMenuItem("").WithStyle(Style{IconImageTemplate: true}),
	}
}

func (p *Plugin) WithIcon(icon io.Reader) *Plugin {
	p.Title = p.Title.WithIcon(icon)
	return p
}

func (p *Plugin) WithText(title string) *Plugin {
	p.Title.Title = title
	return p
}

func (p *Plugin) WithElements(elements ...XbarElement) *Plugin {
	p.Elements = append(p.Elements, elements...)
	return p
}

// Run implements the Plugin API of xbar by rendering its configuration to the standard output.
func (p *Plugin) Run() {
	if err := p.RunW(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// RunW renders the plugin configuration to the specified writer.
//
// This is provided for testing purposes; in other cases the Run function may
// be more convenient.
func (p *Plugin) RunW(w io.Writer) error {
	if _, err := fmt.Fprintln(w, p.Title.renderSelf()); err != nil {
		return err
	}
	if len(p.Elements) > 0 {
		if _, err := fmt.Fprintln(w, "---"); err != nil {
			return err
		}
		for _, item := range p.Elements {
			if err := printElement(w, item, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

func printElement(w io.Writer, el XbarElement, level int) error {
	prefix := strings.Repeat("--", level)
	if _, err := fmt.Fprintf(w, "%s%s\n", prefix, el.renderSelf()); err != nil {
		return err
	}
	// It's important that the child items come before the alt item, otherwise they'll
	// be attached to the alt.
	for _, child := range el.children() {
		if err := printElement(w, child, level+1); err != nil {
			return err
		}
	}
	if alt := el.renderAlt(); alt != "" {
		if _, err := fmt.Fprintf(w, "%s%s alternate=true\n", prefix, alt); err != nil {
			return err
		}
	}
	return nil
}
