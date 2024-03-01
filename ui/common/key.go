package uicommon

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Help        key.Binding
	Quit        key.Binding
	DynamicKeys []key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.DynamicKeys,    // first column
		{k.Help, k.Quit}, // second column
	}
}

func (k KeyMap) BindDynamicKeys(keys map[string]key.Binding) KeyMap {
	k.DynamicKeys = make([]key.Binding, 0)

	for _, v := range keys {
		k.DynamicKeys = append(k.DynamicKeys, v)
	}

	return k
}
