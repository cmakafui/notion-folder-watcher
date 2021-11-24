package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

var client *notionapi.Client

var pathWatcher Publisher = &PathWatcher{}
var watchIndexer Subscriber = &WatchIndexer{}

type WatchIndexer struct{}

func (wi *WatchIndexer) receive(filePath string) {
	folderpath, basename := filepath.Split(filePath)
	// Remove trailing \
	folderpath = folderpath[:len(folderpath)-1]
	for _, element := range box {
		for _, folder := range element.Folders {

			if folder == folderpath {
				// Create page in a database in notion
				time.Sleep(300 * time.Millisecond)
				_, err := api.CreatePage(element.DatabaseID, basename, client)
				if err != nil {
					SendNotification("Error", "Could not send update to notion database")
					continue
				}
				SendNotification("Success", "Pushed "+basename+" to notion")
			}

		}
	}

}

func main() {

	// Initialize config env variables if they don't exist
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		// Notion Integration variables
		notion_token = os.Getenv("NOTION_TOKEN")
		page_id = os.Getenv("PAGE_ID")
	} else {
		response_token, response_token_bool, err := dlgs.Entry("Notion Folder Watcher", "Paste Notion Integration Token here", "")
		if err != nil {
			panic(err)
		}
		response_page, response_page_bool, err := dlgs.Entry("Notion Folder Watcher", "Paste Main Page ID here", "")
		if err != nil {
			panic(err)
		}
		if response_token_bool && response_page_bool {
			notion_token = response_token
			page_id = response_page

			env, err := godotenv.Unmarshal("NOTION_TOKEN=" + notion_token + "\nPAGE_ID=" + page_id)
			if err != nil {
				panic(err)
			}
			err1 := godotenv.Write(env, "./.env")
			if err1 != nil {
				panic(err)
			}
		} else {
			// Exit
			systray.Quit()
		}

	}

	onExit := func() {

		file, _ := json.MarshalIndent(box, "", " ")

		_ = ioutil.WriteFile("config.json", file, 0644)

		now := time.Now()
		fmt.Println(now)
	}

	pathWatcher.Register(&watchIndexer)

	go pathWatcher.Observe()

	systray.Run(onReady, onExit)
}

func onReady() {

	// Initialize box to track notion lists if config.json doesn't exist
	if _, err := os.Stat("config.json"); err == nil {
		jsonFile, err := os.Open("config.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &box)
	} else {
		box = make(map[string]*observer.ListItem)
	}

	// Initialize string slice to track database names
	db_names := []string{}
	for key := range box {
		db_names = append(db_names, key)
	}

	// Initialize Notion client
	client = notionapi.NewClient(notionapi.Token(notion_token))

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

	mLists.Disable()
	mAddFolder.Disable()
	mFolders.Disable()

	// Add submenu items to systray if they already exist
	for key, element := range box {
		mLists.Enable()
		mAddFolder.Enable()
		mLists.AddSubMenuItem(key, "Notion Database")
		for _, folder := range element.Folders {
			pathWatcher.AddPath(folder)
			mFolders.Enable()
			mFolders.AddSubMenuItem(folder, "Watching")
		}
	}

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
				db_name, db_name_bool, err := dlgs.List("List", "Select item from list:", db_names)
				if err != nil {
					// panic(err)
					continue
				}

				// Add a folder to watch for changes
				folder, folderbool, err := dlgs.File("Select folder to watch", "", true)
				if err != nil {
					continue
				}

				if folderbool && db_name_bool {
					// Add folder path to global box
					box[db_name].AddFolder(folder)

					// Add folder to watcher
					pathWatcher.AddPath(folder)

					// Add folder to database
					mFolders.Enable()
					mFolders.AddSubMenuItem(folder, "Watching")

					SendNotification("Success", folder+" added to "+db_name)
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

func SendNotification(title, message string) {
	err := beeep.Notify(title, message, "assets/watch_icon.png")
	if err != nil {
		panic(err)
	}
}
