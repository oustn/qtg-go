package teacomponents

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"strings"
)

// from https://www.patorjk.com/software/taag/#p=display&f=BlurVision%20ASCII&t=QTG
var (
	banner = `┏┓┏┳┓┏┓
┃┃ ┃ ┃┓
┗┻ ┻ ┗┛
一个有趣的蜻蜓FM下载工具`
)

func NewBanner() *Banner {
	b := Banner{}
	return &b
}

type Banner struct {
	width int
}

func (b *Banner) Render() string {
	str := strings.Split(banner, "\n")
	colors := colorGrid(len([]rune(str[3])), len(str))
	sb := strings.Builder{}
	for i, x := range colors {
		for j, y := range x {
			charset := []rune(str[i])
			if j < len(charset) {
				s := lipgloss.NewStyle().SetString(string([]rune(str[i])[j])).Foreground(lipgloss.Color(y))
				sb.WriteString(s.String())
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (b *Banner) Handler(msg tea.Msg) tea.Cmd {
	return nil
}

func (b *Banner) SetSize(width, _ int) {
	b.width = width
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}
