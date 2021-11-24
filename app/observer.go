package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

type Publisher interface {
	Register(subscriber *Subscriber)
	Unregister(subscriber *Subscriber)
	notify(path string)
	AddPath(path string)
	Observe()
}

type Subscriber interface {
	receive(path string)
}

// PathWatcher observes changes in the file system and works as a Publisher for
// the application by notifying subscribers, which will perform other operations.
type PathWatcher struct {
	subscribers []*Subscriber
	watcher     *fsnotify.Watcher
}

// register subscribers to the publisher
func (pw *PathWatcher) Register(subscriber *Subscriber) {
	pw.subscribers = append(pw.subscribers, subscriber)
}

// unregister subscribers from the publisher
func (pw *PathWatcher) Unregister(subscriber *Subscriber) {
	length := len(pw.subscribers)

	for i, sub := range pw.subscribers {
		if sub == subscriber {
			pw.subscribers[i] = pw.subscribers[length-1]
			pw.subscribers = pw.subscribers[:length-1]
			break
		}
	}
}

// notify subscribers that a event has happened, passing the path and the type
// of event as message.
func (pw *PathWatcher) notify(path string) {
	for _, sub := range pw.subscribers {
		(*sub).receive(path)
	}
}

func (pw *PathWatcher) AddPath(path string) {

	pw.watcher.Add(path)
}

// observe changes to the file system using the fsnotify library
func (pw *PathWatcher) Observe() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// watcher.Add(pw.path)

	pw.watcher = watcher

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					pw.notify(event.Name)
				}
			case err := <-watcher.Errors:
				fmt.Println("Error", err)
			}
		}
	}()

	<-done
}
