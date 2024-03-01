package tuimessage

type (
	TickMsg    struct{}
	FrameMsg   struct{}
	CreatedMsg struct{}

	AuthTokenInvalid struct{}
	AuthFailed       struct{}
	AuthSuccess      struct{}

	AuthPageEnter struct{}
	AuthPageExit  struct{}

	SearchPageEnter struct{}
	SearchPageExit  struct{}

	RankPageEnter struct{}
	RankPageExit  struct{}

	SettingPageEnter struct{}
	SettingPageExit  struct{}
)
