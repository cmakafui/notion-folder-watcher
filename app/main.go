package main

import "github.com/cmakafui/notion-folder-watcher/observer"

func main() {

	var pathWatcher observer.Publisher = &observer.PathWatcher{
		Path: "../",
	}

	var pathIndexer observer.Subscriber = &observer.PathIndexer{}
	pathWatcher.Register(&pathIndexer)

	var pathFileMD5 observer.Subscriber = &observer.PathFileMD5{}
	pathWatcher.Register(&pathFileMD5)

	pathWatcher.Observe()
}
