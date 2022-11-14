package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"status-api/database"
	"status-api/structs"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var feedCreatedUpdated = time.Now()

func rssShowHandler(w http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title:       "Status API messages for leon.wtf",
		Link:        &feeds.Link{Href: "https://leon.wtf"},
		Description: "Messages about incidents, maintenance, and others for the services running on leon.wtf",
		Author:      &feeds.Author{Name: "Leon Schmidt", Email: "admin@leon.wtf"},
		Created:     feedCreatedUpdated,
	}

	var items []structs.AtomFeedItemModel
	database.Con.Find(&items)

	//feedsToAdd := make([]*feeds.Item, 0)
	for index, item := range items {
		// workaround: feed.ToAtom() panics, if there
		// is no feeds.Link element, so we will add
		// an empty feeds.Link
		if item.Data.Link == nil {
			items[index].Data.Link = &feeds.Link{}
		}
		// need to address Data via index, otherwise the
		// pointer would always point to the last item
		toAdd := &items[index].Data
		feed.Add((*feeds.Item)(toAdd))
	}

	atomFeedString, err := feed.ToAtom()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write([]byte(atomFeedString))
}

const (
	atomDescriptionStatusResolved   = "Status: RESOLVED"
	atomDescriptionStatusUnresolved = "Status: UNRESOLVED"
)

func rssListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var items []structs.AtomFeedItemModel

	// get "status" parameter
	query := r.URL.Query()
	statusParam := query.Get("status")
	if statusParam == "" { // if not set, retrieve all items
		database.Con.Find(&items)
	} else {
		if strings.ToLower(statusParam) == "resolved" {
			database.Con.Where("Data LIKE ?", `%"Description":"`+atomDescriptionStatusResolved+`"%`).Find(&items)
		} else if strings.ToLower(statusParam) == "unresolved" {
			database.Con.Where("Data LIKE ?", `%"Description":"`+atomDescriptionStatusUnresolved+`"%`).Find(&items)
		} else { // return error
			respondError(&w, errors.New("parameter 'status' can only be one of [resolved, unresolved]"), 400)
			return
		}
	}

	// enrich structs.AtomFeedItem with the database ID (for referencing)
	//
	// we cannot use a map[uint]structs.AtomFeedItem here, as uints are not
	// allowed as property in a JSON dict an are thus converted to string
	type ItemWithDBID struct {
		DbId                 uint `json:"Db_Id"`
		structs.AtomFeedItem `json:",inline"`
	}
	itemsAsSlice := make([]ItemWithDBID, 0)
	for _, item := range items {
		itemsAsSlice = append(itemsAsSlice, ItemWithDBID{
			item.ID,
			item.Data,
		})
	}

	respondInstance(&w, itemsAsSlice, 200)
}

func rssCreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// decode response body
	var respBody struct {
		Title    string
		Content  string
		Resolved bool
	}
	if err := json.NewDecoder(r.Body).Decode(&respBody); err != nil {
		respondError(&w, err, 400) // decoding errors means malformed request
		return
	}

	newFeedItem := &structs.AtomFeedItemModel{
		Data: structs.AtomFeedItem{
			Title:   respBody.Title,
			Content: respBody.Content,
			Created: time.Now(),
		},
	}
	if respBody.Resolved {
		newFeedItem.Data.Description = atomDescriptionStatusResolved
	} else {
		newFeedItem.Data.Description = atomDescriptionStatusUnresolved
	}

	// generate an ID
	titleTruncated := strings.ReplaceAll(newFeedItem.Data.Title, " ", "")
	titleTruncated = strings.ToLower(titleTruncated)
	newFeedItem.Data.Id = fmt.Sprintf("leon.wtf:%s:%s",
		newFeedItem.Data.Created.Format("2006-01-02"),
		titleTruncated,
	)

	res := database.Con.Create(&newFeedItem)
	if res.Error != nil {
		respondError(&w, res.Error, 500)
		return
	}

	// update feedCreatedUpdated time
	feedCreatedUpdated = time.Now()

	respondInstance(&w, fmt.Sprintf("message created with ID '%d'", newFeedItem.ID), 201)
}

func rssChangeMessageHandler(w http.ResponseWriter, r *http.Request) {
	// get path paramter
	vars := mux.Vars(r)
	id, ok := vars["db_id"]
	if !ok {
		respondError(&w, errors.New(`'db_id' is missing`), 400)
		return
	}

	// decode response body
	var respBody struct {
		Title    string
		Content  string
		Resolved *bool
	}
	if err := json.NewDecoder(r.Body).Decode(&respBody); err != nil {
		respondError(&w, err, 400) // decoding errors means malformed request
		return
	}

	// retrieve message to be changed from database
	var messageToChange structs.AtomFeedItemModel
	// SELECT * FROM AtomFeedItemModels WHERE id = 'id'
	res := database.Con.First(&messageToChange, id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			respondError(&w, errors.New("message not found"), 404)
		} else {
			respondError(&w, res.Error, 500)
		}
		return
	}

	// make changes
	changed := make([]string, 0)
	if respBody.Title != "" {
		messageToChange.Data.Title = respBody.Title
		changed = append(changed, "Title")
	}
	if respBody.Content != "" {
		messageToChange.Data.Content = respBody.Content
		changed = append(changed, "Content")
	}
	if respBody.Resolved != nil {
		if *respBody.Resolved {
			messageToChange.Data.Description = atomDescriptionStatusResolved
		} else {
			messageToChange.Data.Description = atomDescriptionStatusUnresolved
		}
		changed = append(changed, "Resolved")
	}

	// update "updated" time
	messageToChange.Data.Updated = time.Now()
	changed = append(changed, "Updated")

	// save the modified model back to the database
	database.Con.Save(&messageToChange)

	// update feedCreatedUpdated time
	feedCreatedUpdated = time.Now()

	respondInstance(
		&w,
		fmt.Sprintf("successfully changed the following attributes for message with ID '%d': %v",
			messageToChange.ID,
			"["+strings.Join(changed, ", ")+"]",
		),
		200,
	)
}

func rssDeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	// get path paramter
	vars := mux.Vars(r)
	id, ok := vars["db_id"]
	if !ok {
		respondError(&w, errors.New(`'db_id' is missing`), 400)
		return
	}

	var messageToDelete structs.AtomFeedItemModel
	// SELECT * FROM AtomFeedItemModels WHERE id = 'id'
	res := database.Con.First(&messageToDelete, id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			respondError(&w, errors.New("message not found"), 404)
		} else {
			respondError(&w, res.Error, 500)
		}
		return
	}

	database.Con.Delete(&messageToDelete)

	// update feedCreatedUpdated time
	feedCreatedUpdated = time.Now()

	respondInstance(
		&w,
		fmt.Sprintf("deleted message '%s'",
			messageToDelete.Data.Title,
		),
		200,
	)
}
