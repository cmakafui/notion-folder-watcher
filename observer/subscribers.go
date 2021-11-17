package observer

import (
	"fmt"
	"path/filepath"

	"github.com/gen2brain/beeep"
)

type PathIndexer struct{}

func (pi *PathIndexer) receive(filePath, event string) {
	_, basename := filepath.Split(filePath)
	err := beeep.Notify(event, "Resource: "+basename, "assets/watch_icon.png")
	if err != nil {
		panic(err)
	}
}

type PathFileMD5 struct{}

func (pfm *PathFileMD5) receive(path, event string) {
	fmt.Printf("Syncing: %v, %v\n", path, event)
}
