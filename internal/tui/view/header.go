package tuiview

import tuicommon "github.com/oustn/qtg/internal/tui/common"

func HeaderView(q tuicommon.Q) string {
	config, authed, nickName := q.GetAuthInfo()
	info := q.Info()

	tpl := "蜻蜓FM下载工具"
	s := tpl
	if info.Home {
		s = "欢迎使用" + tpl
	}
	if info.Choice != tuicommon.AuthPage || info.Home {
		if !config {
			s += tuicommon.BlurredStyle(" 请先配置 Token")
		} else {
			if authed {
				s += "，" + nickName
			} else {
				s += tuicommon.BlurredStyle(" (认证失败)")
			}
		}

	} else {
		s += "-"
		s += info.ChoiceTitle
	}
	s += "\n\n"
	return s
}
