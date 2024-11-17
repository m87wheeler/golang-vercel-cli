package menu

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/m87wheeler/golang-vercel-cli/pkg/utils"

	"github.com/buger/goterm"
	"github.com/pkg/term"
)

// Raw input keycodes
var up byte = 65
var down byte = 66
var escape byte = 27
var enter byte = 13
var space byte = 32
var keys = map[byte]bool{
	up:   true,
	down: true,
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		Prompt:    prompt,
		MenuItems: make([]*MenuItem, 0),
	}
}

// AddItem will add a new menu option to the menu list
func (m *Menu) AddItem(id string, option string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID:   id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

// Display will display the current menu options and awaits user selection
// It returns the users selected choice
func (m *Menu) Display() string {
	defer func() {
		// Show cursor again.
		fmt.Printf("\033[?25h")
	}()

	fmt.Printf("%s\n", goterm.Color(goterm.Bold(m.Prompt)+":", goterm.CYAN))

	m.renderMenuItems(false, false, []string{})

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	for {
		keyCode := getInput()
		switch keyCode {
		case escape:
			return ""
		case enter:
			menuItem := m.MenuItems[m.CursorPos]
			fmt.Println("\r")
			return menuItem.ID
		case up:
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true, false, []string{})
		case down:
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true, false, []string{})
		}
	}
}

func (m *Menu) DisplayMultiChoice(f func(c string) []string) string {
	defer func() {
		// Show cursor again.
		fmt.Printf("\033[?25h")
	}()

	fmt.Printf("%s\n", goterm.Color(goterm.Bold(m.Prompt)+":", goterm.CYAN))

	// Store multi-choice selection
	selection := f("")
	m.renderMenuItems(false, true, selection)

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	for {
		keyCode := getInput()
		switch keyCode {
		case escape:
			return ""
		case space:
			menuItem := m.MenuItems[m.CursorPos]
			selection = f(menuItem.ID)
			m.renderMenuItems(true, true, selection)
		case enter:
			menuItem := m.MenuItems[m.CursorPos]
			fmt.Println("\r")
			return menuItem.ID
		case up:
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true, true, selection)
		case down:
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true, true, selection)
		}
	}
}

func (m *Menu) DisplayInfoTable(data []InfoTableData) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Header row
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "Field      \tValue         \t\n")
	fmt.Fprintf(w, "-----------\t-------------\t\n")

	// Body rows
	for _, d := range data {
		fmt.Fprintf(w, "%s       \t%s           \t\n", d.Label, d.Data)
	}
	fmt.Fprintf(w, "\n")

	w.Flush()
}

// getInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func getInput() byte {
	t, _ := term.Open("/dev/tty")

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)

	t.Restore()
	t.Close()

	// Arrow keys are prefixed with the ANSI escape code which take up the first two bytes.
	// The third byte is the key specific value we are looking for.
	// For example the left arrow key is '<esc>[A' while the right is '<esc>[C'
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}

// renderMenuItems prints the menu item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (m *Menu) renderMenuItems(redraw bool, multi bool, selection []string) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		checkbox := "\u2610"
		if utils.Contains(selection, menuItemText) {
			checkbox = "\u2612"
		}

		if index == m.CursorPos {
			checkbox = goterm.Color(checkbox, goterm.YELLOW)
			cursor = goterm.Color("> ", goterm.YELLOW)
			menuItemText = goterm.Color(menuItemText, goterm.YELLOW)
		}

		if multi {
			fmt.Printf("\r%s %s %s%s", cursor, checkbox, menuItemText, newline)
		} else {
			fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
		}
	}
}
