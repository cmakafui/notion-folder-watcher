package observer

import (
	"path/filepath"

	"github.com/gen2brain/beeep"
)

type PathIndexer struct{}

func (pi *PathIndexer) receive(filePath string) {
	_, basename := filepath.Split(filePath)
	err := beeep.Notify("Watching", "Resource: "+basename, "assets/watch_icon.png")
	if err != nil {
		panic(err)
	}
}
