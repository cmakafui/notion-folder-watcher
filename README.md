# notion-folder-watcher

A cross platform desktop service that watches custom folders for file changes and updates the corresponding database in Notion. Perfect for tracking reading lists

## Main Features

1. Watches folders for content changes and syncs to notion database
2. Track Status of each resource in the folder
3. Organize resources in Notion

## Components

1. Developed as a service
2. A file/folder watcher monitoring for changes
3. Notion API for creating and updating tables with a worker pool
4. Desktop Notifications for event triggers
5. Dialog Boxes to input custom information
6. System Tray to quick management of the service

## Installation

The app was developed to be cross platform. currently supports windows and linux.
On windows you use the release binary as is.

On linux, some prerequesites need to be installed.

`sudo apt-get install gcc libgtk-3-dev libappindicator3-dev`

On Linux Mint, `libxapp-dev` is also required

## Usage

Before using the app, integrations needs to be generated in Notion. Instructions on
how to set up an internal integration token and add to a page are in the link below

[Notion Integrations](https://www.notion.so/help/add-and-manage-integrations-with-the-api)

When starting the app for the first time, it would ask for both the integration token and the
page id.
After that, you can use the app by accessing it in the system tray to create tables and watch custom folders.
