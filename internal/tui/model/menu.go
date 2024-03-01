package tuimodel

import (
	tea "github.com/charmbracelet/bubbletea"
	tuicommon "github.com/oustn/qtg/internal/tui/common"
)

type menu struct {
	Home      bool   // Whether the home view is active
	Choice    string // The user's choice
	Ticks     int    // Auto close ticks
	AutoClose bool
}

func (m *menu) reset() {
	m.Home = true
	m.Choice = tuicommon.SearchPage
	m.Ticks = 10
}

func (m *menu) indexOfChoice() int {
	for i, choice := range tuicommon.MenuList {
		if choice == m.Choice {
			return i
		}
	}
	return -1
}

func (m *menu) SelectUp() {
	m.AutoClose = false
	choice := m.indexOfChoice()
	choice--
	if choice < 0 {
		choice = 0
	}
	m.Choice = tuicommon.MenuList[choice]
}

func (m *menu) SelectDown() {
	m.AutoClose = false
	choice := m.indexOfChoice()
	choice++
	if choice > len(tuicommon.MenuList)-1 {
		choice = len(tuicommon.MenuList) - 1
	}
	m.Choice = tuicommon.MenuList[choice]
}

func (m *menu) SelectEnter() tea.Cmd {
	m.AutoClose = false
	m.Home = false
	return tuicommon.PageEnter(m.indexOfChoice())
}

func (m *menu) Countdown() {
	m.Ticks--
}

func (m *menu) GetChoiceTitle() string {
	choice := m.indexOfChoice()

	return tuicommon.MenuTitle[choice]
}

func initMenu() menu {
	return menu{
		Home:      true,
		Choice:    tuicommon.SearchPage,
		Ticks:     10,
		AutoClose: true,
	}
}
