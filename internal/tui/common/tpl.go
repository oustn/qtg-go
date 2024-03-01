package tuicommon

import "fmt"

var (
	Dot = ColorFg(" • ", "236")
)

func Checkbox(label string, checked bool) string {
	if checked {
		return ColorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}
