package teacomponents

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Height represents the height of the statusbar.
const Height = 1

type ColorConfig struct {
	Foreground lipgloss.AdaptiveColor
	Background lipgloss.AdaptiveColor
}

type Column struct {
	Content      string
	Foreground   lipgloss.AdaptiveColor
	Background   lipgloss.AdaptiveColor
	Width        int
	ContentWidth int
	Flex         bool
}

type renderColumn struct {
	content string
	style   lipgloss.Style
	width   int
	flex    bool
	index   int
}

type StatusBar struct {
	Width   int
	Height  int
	Columns []Column
}

func NewStatusBar(column int) StatusBar {
	return StatusBar{
		Height:  Height,
		Columns: make([]Column, column),
	}
}

func (m *StatusBar) SetSize(width, height int) {
	m.Width = width
}

func (m *StatusBar) Handler(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	return nil
}

func (m *StatusBar) SetColumns(columns []Column) {
	m.Columns = columns[0:len(m.Columns)]
}

func (m *StatusBar) UpdateColumn(index int, content string) {
	if index > len(m.Columns) {
		return
	}
	m.Columns[index].Content = content
}

func (m *StatusBar) UpdateColumnContent(index int, content string) {
	if index > len(m.Columns) {
		return
	}
	m.Columns[index].Content = content
}

func (m *StatusBar) UpdateColumnForeground(index int, color lipgloss.AdaptiveColor) {
	if index > len(m.Columns) {
		return
	}
	m.Columns[index].Foreground = color
}

func (m *StatusBar) UpdateColumnBackground(index int, color lipgloss.AdaptiveColor) {
	if index > len(m.Columns) {
		return
	}
	m.Columns[index].Background = color
}

func (m *StatusBar) UpdateColumnWidth(index int, width int) {
	if index > len(m.Columns) {
		return
	}
	m.Columns[index].ContentWidth = width
}

func (m *StatusBar) Render() string {
	width := lipgloss.Width

	var columns []renderColumn
	var contents []string

	for i, column := range m.Columns {
		style := lipgloss.NewStyle().Padding(0, 1).Height(m.Height)
		if column.Foreground != (lipgloss.AdaptiveColor{}) {
			style = style.Foreground(column.Foreground)
		}
		if column.Background != (lipgloss.AdaptiveColor{}) {
			style = style.Background(column.Background)
		}
		if column.Width > 0 {
			style.Width(column.Width)
		}
		if column.Flex {
			columns = append(columns, renderColumn{
				content: column.Content,
				style:   style,
				index:   i,
			})
			contents = append(contents, column.Content)
			continue
		}
		if column.ContentWidth > 0 {
			content := style.Render(truncate.StringWithTail(column.Content, uint(column.ContentWidth), "..."))
			columns = append(columns, renderColumn{
				content: content,
				width:   width(content),
				index:   i,
			})
			contents = append(contents, content)
			continue
		}
		content := style.Render(column.Content)
		columns = append(columns, renderColumn{
			content: content,
			width:   width(content),
			index:   i,
		})
		contents = append(contents, content)
	}

	availableWidth := m.Width
	flexColumns := make([]int, 0)
	for i, column := range columns {
		if column.width == 0 {
			flexColumns = append(flexColumns, i)
		} else {
			availableWidth -= column.width
		}
	}

	if len(flexColumns) > 0 {
		flexWidth := availableWidth / len(flexColumns)
		for _, i := range flexColumns {
			columns[i].style.Width(flexWidth)
			contents[i] = columns[i].style.Render(truncate.StringWithTail(contents[i], uint(flexWidth), "..."))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, contents...)
}
