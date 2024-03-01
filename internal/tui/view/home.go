package tuiview

import (
	"fmt"
	tuicommon "github.com/oustn/qtg/internal/tui/common"
	"strconv"
	"strings"
)

func HomeView(q tuicommon.Q) string {
	c, tick, autoClose := q.GetMenuInfo()
	tpl := "%s\n\n"
	if autoClose {
		tpl += "应用将会在 %s 秒后退出\n\n"
	}
	tpl += tuicommon.SubtitleStyle("j/k, up/down: 选择") + tuicommon.Dot + tuicommon.SubtitleStyle("enter: 进入") + tuicommon.Dot + tuicommon.SubtitleStyle("q, esc: 退出")

	boxes := make([]interface{}, len(tuicommon.MenuTitle))
	formats := make([]string, len(tuicommon.MenuTitle))
	for i, title := range tuicommon.MenuTitle {
		formats[i] = "%s"
		boxes[i] = tuicommon.Checkbox(title, c == tuicommon.MenuList[i])
	}

	choices := fmt.Sprintf(
		strings.Join(formats, "\n"),
		boxes...,
	)

	if autoClose {
		return fmt.Sprintf(tpl, choices, tuicommon.ColorFg(strconv.Itoa(tick), "79"))
	} else {
		return fmt.Sprintf(tpl, choices)
	}
}
