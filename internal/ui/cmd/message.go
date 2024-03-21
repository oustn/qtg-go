package teacmd

import (
	"github.com/charmbracelet/bubbles/key"
	teacommon "github.com/oustn/qtg/internal/ui/common"
)

type (
	UpdateStatusBarMsg struct {
		Title string
		Msg   string
	}

	AuthMsg struct {
		Id       string
		Token    string
		Callback teacommon.Callback
	}

	AuthSuccessMsg struct{}

	AuthFailMsg struct {
		Err error
	}

	UpdateQuitKeyMsg struct {
		Q   bool
		Esc bool
	}

	RouteMsg struct {
		Route string
	}

	ViewCreatedMsg struct {
	}

	HelpMsg struct {
		Keys []*key.Binding
	}

	SearchMsg struct {
		Keyword string
		Type    string
		Page    int
	}

	SearchResultMsg struct {
		Error  error
		Result teacommon.SearchResult
	}

	DownloadMsg struct {
		Channel teacommon.Channel
	}

	TickMsg struct{}
)
