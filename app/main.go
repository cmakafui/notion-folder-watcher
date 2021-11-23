package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/joho/godotenv"
)

func main() {

	onExit := func() {
		now := time.Now()
		fmt.Println(now)
	}

	systray.Run(onReady, onExit)
	// var pathWatcher observer.Publisher = &observer.PathWatcher{
	// 	Path: "../",
	// }

	// var pathIndexer observer.Subscriber = &observer.PathIndexer{}
	// pathWatcher.Register(&pathIndexer)

	// pathWatcher.Observe()
}

func onReady() {

	// Initialize config env variables if they don't exist
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	systray.SetIcon(icon.Data)
	systray.SetTitle("Notion Folder Watcher")
	systray.SetTooltip("Sync folders to Notion")

	mAddList := systray.AddMenuItem("Create Database", "Create a Notion Database")
	mAddFolder := systray.AddMenuItem("Watch Folder", "Watch a folder for content")
	systray.AddSeparator()
	mLists := systray.AddMenuItem("Databases", "Notion Databases")
	mFolders := systray.AddMenuItem("Folders", "Watched Folders")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)

	mAddFolder.Disable()
	mLists.Disable()
	mFolders.Disable()

	go func() {
		for {
			select {
			// Create a new database on notion
			case <-mAddList.ClickedCh:
				response, response_bool, err := dlgs.Entry("Notion Folder Watcher", "Enter New List Name", "")
				if err != nil {
					panic(err)
				}
				if !response_bool {
					_, err := dlgs.Error("Error", "No name inputted")
					if err != nil {
						log.Fatal(err)
					}
				}
				mLists.Enable()
				mAddFolder.Enable()

				mLists.AddSubMenuItem(response, "Notion Database")
				_, err1 := dlgs.Info("Info", response+" created Successfully")
				if err1 != nil {
					log.Fatal(err1)
				}
			case <-mAddFolder.ClickedCh:
				// Add a folder to watch for changes
				folder, _, err := dlgs.File("Select folder to watch", "", true)
				if err != nil {
					continue
				}
				// Add folder to database
				mFolders.Enable()
				mFolders.AddSubMenuItem(folder, "Watching")
				_, err1 := dlgs.Info("Info", "Added Successfully")
				if err1 != nil {
					log.Fatal(err1)
				}

			case <-mQuit.ClickedCh:
				fmt.Println("Requesting quit")
				systray.Quit()
				fmt.Println("Finished quitting")
				return
			}
		}
	}()

}
