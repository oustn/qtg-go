package uicommon

const (
	HomeRoute   = "home"
	AuthRoute   = "auth"
	SearchRoute = "search"
)

type Route struct {
	Path  string
	Title string
}

var (
	Routes = []Route{
		{Path: HomeRoute, Title: "主页 🏠"},
		{Path: AuthRoute, Title: "Token 配置 🔑"},
		{Path: SearchRoute, Title: "搜索 🔍"},
	}
)
