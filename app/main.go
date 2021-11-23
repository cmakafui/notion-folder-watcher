package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/cmakafui/notion-folder-watcher/api"
	"github.com/cmakafui/notion-folder-watcher/observer"
	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
)

// Global variable to save state on exit. (Find a better design)
var box map[string]*observer.ListItem
var notion_token string
var page_id string

func main() {

	// Initialize config env variables if they don't exist
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	onExit := func() {

		file, _ := json.MarshalIndent(box, "", " ")

		_ = ioutil.WriteFile("config.json", file, 0644)

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

	// Notion Integration variables
	notion_token = os.Getenv("NOTION_TOKEN")
	page_id = os.Getenv("PAGE_ID")

	// Initialize box to track notion lists
	box = make(map[string]*observer.ListItem)

	// Initialize string slice to track database names
	db_names := []string{}

	// Initialize Notion client
	client := notionapi.NewClient(notionapi.Token(notion_token))

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
				response, response_bool, err := dlgs.Entry("Notion Folder Watcher", "Enter Database Name", "")
				if err != nil {
					panic(err)
				}
				if !response_bool {
					SendNotification("Error: Try again", "No name inputted")
					continue
				}

				// Create database in notion
				db_id, err1 := api.CreateDB(response, page_id, client)
				if err1 != nil {
					SendNotification("Error", "Could not create database on Notion")
					continue
				}

				// Update database slice
				db_names = append(db_names, response)

				// Add database to struct
				item := observer.ListItem{Name: response, DatabaseID: db_id, Client: client}
				box[response] = &item

				mLists.Enable()
				mAddFolder.Enable()

				mLists.AddSubMenuItem(response, "Notion Database")
				SendNotification("Success", response+" created successfully")
			case <-mAddFolder.ClickedCh:
				// Select database option to use from list
				db_name, _, err := dlgs.List("List", "Select item from list:", db_names)
				if err != nil {
					// panic(err)
					continue
				}

				// Add a folder to watch for changes
				folder, _, err := dlgs.File("Select folder to watch", "", true)
				if err != nil {
					continue
				}

				// Add folder path to global box
				box[db_name].AddFolder(folder)

				// Add folder to database
				mFolders.Enable()
				mFolders.AddSubMenuItem(folder, "Watching")

				SendNotification("Success", folder+" added to "+db_name)

			case <-mQuit.ClickedCh:
				fmt.Println("Requesting quit")
				systray.Quit()
				fmt.Println("Finished quitting")
				return
			}
		}
	}()

}

func SendNotification(title, message string) {
	err := beeep.Notify(title, message, "assets/watch_icon.png")
	if err != nil {
		panic(err)
	}
}
