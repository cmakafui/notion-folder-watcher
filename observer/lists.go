package observer

import (
	"log"
	"path/filepath"
	"time"

	"github.com/cmakafui/notion-folder-watcher/api"
	"github.com/gen2brain/beeep"
	"github.com/jomei/notionapi"
)

type ListWatcher interface {
	ChangeName(name string)
	AddFolder(path string)
	RemoveFolder(path string)
	receive(path, event string)
}

type List struct {
	name        string
	database_id string
	folders     []string
	client      *notionapi.Client
}

func (ll *List) ChangeName(name string) {
	ll.name = name
}

func (ll *List) AddFolder(path string) {
	ll.folders = append(ll.folders, path)
}

func (ll *List) RemoveFolder(path string) {
	length := len(ll.folders)

	for i, folder := range ll.folders {
		if folder == path {
			ll.folders[i] = ll.folders[length-1]
			ll.folders = ll.folders[:length-1]
			break
		}
	}
}

func (ll *List) receive(filePath string) {
	path, basename := filepath.Split(filePath)

	for _, folder := range ll.folders {
		if folder == path {
			time.Sleep(300 * time.Millisecond)

			url, err := api.CreatePage(ll.database_id, basename, ll.client)
			if err != nil {
				log.Println(err)
			}
			err1 := beeep.Notify("Entry Created", url, "assets/watch_icon.png")
			if err1 != nil {
				log.Println(err1)
			}
		}
	}
}
