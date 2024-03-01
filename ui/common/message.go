package uicommon

type (
	TickMsg         struct{}
	ReplaceRouteMsg struct {
		Name string
	}
	AuthSuccessMsg struct {
		Token string
		Id    string
	}
)
