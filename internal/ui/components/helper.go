package teacomponents

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math"
	"strings"
)

var (
	keyWidth = 14
)

func NewHelper() Helper {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})

	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#3C3C3C",
	})

	return Helper{
		Separator:      "    ",
		ShortSeparator: " • ",
		Ellipsis:       "…",
		Title:          "帮助",
		Styles: Styles{
			ShortKey:       keyStyle,
			ShortDesc:      descStyle,
			ShortSeparator: sepStyle,
			Ellipsis:       sepStyle.Copy(),
			Key: keyStyle.Copy().
				Bold(true).
				MarginRight(2).
				Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}),
			Desc: descStyle.Copy().
				Foreground(lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"}).
				Copy(),
			Separator: sepStyle.Copy(),
		},
		HelpKey: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "帮助"),
		),
		EscQuitKey: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "返回主页"),
		),
		QQuitKey: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "返回主页"),
		),
		ForceQuitKey: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "强制退出"),
		),
	}
}

type Styles struct {
	Ellipsis       lipgloss.Style
	TitleStyle     lipgloss.Style
	ShortKey       lipgloss.Style
	ShortDesc      lipgloss.Style
	ShortSeparator lipgloss.Style
	Key            lipgloss.Style
	Desc           lipgloss.Style
	Separator      lipgloss.Style
}

type Helper struct {
	Width          int
	Show           bool
	Separator      string
	ShortSeparator string
	Ellipsis       string
	Styles         Styles
	keys           []*key.Binding
	HelpKey        key.Binding
	QQuitKey       key.Binding
	EscQuitKey     key.Binding
	ForceQuitKey   key.Binding
	Title          string
}

func (h *Helper) Handler(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.SetSize(msg.Width, msg.Height)
	}
	return nil
}

func (h *Helper) Render() string {
	if h.Show {
		return lipgloss.NewStyle().
			Render(lipgloss.JoinVertical(
				lipgloss.Top,
				h.renderTitle(),
				h.renderFull(),
			))
	}
	return h.renderShort()
}

func (h *Helper) renderTitle() string {
	return h.Styles.TitleStyle.Render(h.Title)
}

func (h *Helper) renderFull() string {
	if len(h.keys) == 0 {
		return ""
	}
	var (
		out      []string
		sep      = h.Styles.Separator.Render(h.Separator)
		sepWidth = lipgloss.Width(sep)
		maxWidth = 0
	)

	maxKeyWidth := h.resolveMaxKey()

	for _, kb := range h.keys {
		if !kb.Enabled() {
			continue
		}
		helpStr, helpWidth := h.renderKey(kb.Help().Key, kb.Help().Desc, maxKeyWidth)
		if helpWidth > maxWidth {
			maxWidth = helpWidth
		}
		out = append(out, helpStr)
	}

	column := int(math.Floor(float64(h.Width+sepWidth) / float64(maxWidth+sepWidth)))

	if column < 1 {
		column = 1
	} else if column > 3 {
		column = 3
	}

	row := int(math.Floor(float64(len(out)) / float64(column)))
	rest := len(out) % column
	columns := make([]string, column)

	columnStyle := lipgloss.NewStyle().Width(maxWidth + sepWidth)

	j := 0
	current := 0

	for _, str := range out {
		columns[j] = lipgloss.JoinVertical(lipgloss.Left, columns[j], columnStyle.Render(str))
		current++

		if j < rest {
			if current == row+1 {
				j++
				current = 0
			}
		} else {
			if current == row {
				j++
				current = 0
			}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Left, columns...))
}

func (h *Helper) resolveMaxKey() int {
	maxLen := 0
	for _, kb := range h.keys {
		keyText := h.Styles.Key.Copy().Render(kb.Help().Key)
		w := lipgloss.Width(keyText)
		if w > keyWidth {
			return keyWidth
		}
		if w > maxLen {
			maxLen = w
		}
	}
	return maxLen
}

func (h *Helper) renderKey(key, desc string, maxKeyWidth int) (str string, width int) {
	keyText := h.Styles.Key.Copy().Width(maxKeyWidth).
		Render(key)

	descriptionText := h.Styles.Desc.
		Render(desc)

	str = lipgloss.JoinHorizontal(lipgloss.Top, keyText, descriptionText)
	width = lipgloss.Width(str)
	return
}

func (h *Helper) renderShort() string {
	bindings := []key.Binding{h.HelpKey, h.EscQuitKey, h.QQuitKey, h.ForceQuitKey}

	if len(bindings) == 0 {
		return ""
	}

	var b strings.Builder
	var totalWidth int
	var separator = h.Styles.ShortSeparator.Inline(true).Render(h.ShortSeparator)

	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		var sep string
		if totalWidth > 0 && i < len(bindings) {
			sep = separator
		}

		str := sep +
			h.Styles.ShortKey.Inline(true).Render(kb.Help().Key) + " " +
			h.Styles.ShortDesc.Inline(true).Render(kb.Help().Desc)

		w := lipgloss.Width(str)

		// If adding this help item would go over the available width, stop
		// drawing.
		if h.Width > 0 && totalWidth+w > h.Width {
			// Although if there's room for an ellipsis, print that.
			tail := " " + h.Styles.Ellipsis.Inline(true).Render(h.Ellipsis)
			tailWidth := lipgloss.Width(tail)

			if totalWidth+tailWidth < h.Width {
				b.WriteString(tail)
			}

			break
		}

		totalWidth += w
		b.WriteString(str)
	}

	return b.String()
}

func (h *Helper) SetSize(width, _ int) {
	h.Width = width
}

func (h *Helper) SetKeys(keys []*key.Binding) {
	h.keys = keys
}

func (h *Helper) ToggleShowAll() tea.Cmd {
	h.Show = !h.Show
	return nil
}

func (h *Helper) SetEscQuitEnable(enable bool) {
	h.EscQuitKey.SetEnabled(enable)
}

func (h *Helper) SetQQuitEnable(enable bool) {
	h.QQuitKey.SetEnabled(enable)
}

func (h *Helper) ResetQuit() {
	h.QQuitKey.SetEnabled(true)
	h.EscQuitKey.SetEnabled(true)
}

func (h *Helper) ClearKeys() {
	h.ResetQuit()
	h.keys = []*key.Binding{}
}
