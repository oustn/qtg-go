package tuiview

import (
	"fmt"
	tuicommon "github.com/oustn/qtg/internal/tui/common"
	"strings"
)

var (
	authBlurredButton = fmt.Sprintf("[ %s ]", tuicommon.BlurredStyle("认证"))
	authFocusedButton = tuicommon.FocusedStyle.Copy().Render("[ 认证 ]")
)

func AuthView(q tuicommon.Q) string {
	var b strings.Builder
	info := q.Info()

	b.WriteString(tuicommon.Dot + tuicommon.SubtitleStyle("抓取 https://user.qtfm.cn/u2/api/v4/auth 接口获取"))
	b.WriteString("\n\n")
	for i := range info.AuthInputs {
		if i == 0 {
			b.WriteString(tuicommon.SubtitleStyle("Refresh Token:\n"))
		} else {
			b.WriteString(tuicommon.SubtitleStyle("Qing Ting Id:\n"))
		}
		b.WriteString(info.AuthInputs[i].View())
		if i < len(info.AuthInputs)-1 {
			b.WriteRune('\n')
			b.WriteRune('\n')
		}
	}

	button := &authBlurredButton
	if info.AuthFocusIndex == len(info.AuthInputs) {
		button = &authFocusedButton
	}
	_, _ = fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
