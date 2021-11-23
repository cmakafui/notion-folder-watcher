package api

import (
	"context"

	"github.com/jomei/notionapi"
)

// Entry in Notion Database

func CreateDB(name string, page_id notionapi.PageID, client *notionapi.Client) (string, error) {

	db := struct {
		request *notionapi.DatabaseCreateRequest
	}{
		request: &notionapi.DatabaseCreateRequest{
			Parent: notionapi.Parent{
				Type:   notionapi.ParentTypePageID,
				PageID: page_id,
			},
			Title: []notionapi.RichText{

				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{Content: name},
				},
			},
			Properties: notionapi.PropertyConfigs{
				"Name": notionapi.TitlePropertyConfig{
					Type: notionapi.PropertyConfigTypeTitle,
				},
				"Tags": notionapi.MultiSelectPropertyConfig{
					ID:          ";s|V",
					Type:        notionapi.PropertyConfigTypeMultiSelect,
					MultiSelect: notionapi.Select{Options: []notionapi.Option{{ID: "id", Name: "tag", Color: "blue"}}},
				},
				"Status": notionapi.SelectPropertyConfig{
					ID:     "status",
					Type:   notionapi.PropertyConfigTypeSelect,
					Select: notionapi.Select{Options: []notionapi.Option{{ID: "id", Name: "Status"}}},
				},
				"Rating": notionapi.SelectPropertyConfig{
					ID:     "rating",
					Type:   notionapi.PropertyConfigTypeSelect,
					Select: notionapi.Select{Options: []notionapi.Option{{ID: "id", Name: "Rating"}}},
				},
			},
		},
	}

	// Make Create Request
	got, err := client.Database.Create(context.Background(), db.request)
	if err != nil {
		return "", err
	}

	got.Properties = nil
	return got.Title[0].PlainText, nil

}

func UpdateDB(name string, db_id notionapi.DatabaseID, client *notionapi.Client) (string, error) {

	db := struct {
		request *notionapi.DatabaseUpdateRequest
	}{
		request: &notionapi.DatabaseUpdateRequest{
			Title: []notionapi.RichText{

				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{Content: name},
				},
			},
		},
	}

	// Change Database Name
	got, err := client.Database.Update(context.Background(), db_id, db.request)
	if err != nil {
		return "", err
	}

	got.Properties = nil
	return got.Title[0].PlainText, nil

}

func CreatePage(db_id string, name string, client *notionapi.Client) (string, error) {

	page := struct {
		request *notionapi.PageCreateRequest
	}{
		request: &notionapi.PageCreateRequest{
			Parent: notionapi.Parent{
				// Type:       notionapi.ParentTypeDatabaseID,
				DatabaseID: notionapi.DatabaseID(db_id),
			},
			Properties: notionapi.Properties{
				"Name": notionapi.TitleProperty{
					Title: []notionapi.RichText{
						{Text: notionapi.Text{Content: name}},
					},
				},
			},
		},
	}

	got, err := client.Page.Create(context.Background(), page.request)

	if err != nil {
		return "", err
	}

	got.Properties = nil

	return got.URL, nil

}
