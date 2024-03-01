package tuicommon

import (
	tea "github.com/charmbracelet/bubbletea"
	tuimessage "github.com/oustn/qtg/internal/tui/message"
)

const (
	HomePage    = "HomePage"
	AuthPage    = "AuthPage"
	SearchPage  = "SearchPage"
	RankPage    = "RankPage"
	SettingPage = "SettingPage"
)

var (
	MenuList = []string{
		SearchPage,
		RankPage,
		AuthPage,
		SettingPage,
	}

	MenuTitle = []string{
		"搜索 🔍",
		"排行榜 🔥",
		"Token 配置 🔑",
		"设置 ⚙️",
	}

	MenuEnterMsg = []tea.Msg{
		tuimessage.SearchPageEnter{},
		tuimessage.RankPageEnter{},
		tuimessage.AuthPageEnter{},
		tuimessage.SettingPageEnter{},
	}
)
