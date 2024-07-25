//go:build windows
// +build windows

package selector

import (
	"fmt"
	"log"
	"os"

	"github.com/buger/goterm"
	"golang.org/x/sys/windows"
)

type Menu struct {
	Prompt    string
	CursorPos int
	MenuItems []*MenuItem
}

type MenuItem struct {
	Text string
	ID   string
}

var up byte = 65
var down byte = 66
var escape byte = 27
var enter byte = 13
var ctrl_c byte = 3

var keys = map[byte]bool{
	up:   true,
	down: true,
}

func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID:   id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)

	return m
}

func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"

		if index == len(m.MenuItems)-1 {
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		if index == m.CursorPos {
			cursor = goterm.Color("> ", goterm.YELLOW)
			menuItemText = goterm.Color(menuItemText, goterm.YELLOW)
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

func (m *Menu) Display() string {
	defer func() {
		fmt.Printf("\033[?25h")
	}()

	fmt.Printf("%s\n", goterm.Color(goterm.Bold(m.Prompt)+":", goterm.CYAN))

	m.renderMenuItems(false)

	fmt.Printf("\033[?25l")

	for {
		keyCode := getInput()
		if keyCode == escape {
			return ""
		} else if keyCode == enter {
			menuItem := m.MenuItems[m.CursorPos]

			fmt.Println("\r")

			return menuItem.ID
		} else if keyCode == up {
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if keyCode == down {
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if keyCode == ctrl_c {
			return ""
		}
	}
}

func (m *Menu) Clear() {
	for i := 0; i < (len(m.MenuItems) + 1); i++ {
		fmt.Print("\033[F\033[K")
	}
}

func getInput() byte {
	hConsole := windows.Handle(os.Stdin.Fd())

	var originalMode uint32

	err := windows.GetConsoleMode(hConsole, &originalMode)
	if err != nil {
		log.Fatal(err)
	}

	rawMode := originalMode &^ (windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_INPUT)

	err = windows.SetConsoleMode(hConsole, rawMode)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		windows.SetConsoleMode(hConsole, originalMode)
	}()

	readBytes := make([]byte, 3)
	read, err := os.Stdin.Read(readBytes)
	if err != nil {
		log.Fatal(err)
	}

	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		Prompt:    prompt,
		MenuItems: make([]*MenuItem, 0),
	}
}
