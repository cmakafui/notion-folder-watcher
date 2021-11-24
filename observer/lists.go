package observer

import (
	"github.com/jomei/notionapi"
)

type ListWatcher interface {
	ChangeName(name string)
	AddFolder(path string)
	RemoveFolder(path string)
	receive(path, event string)
}

type ListItem struct {
	Name       string
	DatabaseID string
	Folders    []string
	Client     *notionapi.Client
}

func (ll *ListItem) ChangeName(name string) {
	ll.Name = name
}

func (ll *ListItem) AddFolder(path string) {
	ll.Folders = append(ll.Folders, path)
}

func (ll *ListItem) RemoveFolder(path string) {
	length := len(ll.Folders)

	for i, folder := range ll.Folders {
		if folder == path {
			ll.Folders[i] = ll.Folders[length-1]
			ll.Folders = ll.Folders[:length-1]
			break
		}
	}
}
