package uicommon

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Render() string
	Handle(msg tea.Msg) tea.Cmd
	ShouldRenderUserInfo() bool
	KeyBinds() map[string]key.Binding
}

type InitView func() View

type UserInfo struct {
	Nickname  string `json:"nick_name"`
	AvatarUrl string `json:"avatar"`
	Vip       string `json:"private_info.vip_info.vip"`
	VipExpire string `json:"private_info.vip_info.expire_time"`
}

type ViewPayload struct {
	UserInfo     *UserInfo
	RefreshToken string
	QingTingId   string
	Created      bool
}
