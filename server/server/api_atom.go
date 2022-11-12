package server

import (
	"errors"
	"fmt"
	"net/http"
	"status-api/database"
	"status-api/structs"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func rssShowHandler(w http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title:       "Status API messages for leon.wtf",
		Link:        &feeds.Link{Href: "https://leon.wtf"},
		Description: "Messages about incidents, maintenance, and others for the services running on leon.wtf",
		Author:      &feeds.Author{Name: "Leon Schmidt", Email: "admin@leon.wtf"},
		Created:     time.Now(),
	}

	var items []structs.AtomFeedItemModel
	database.Con.Find(&items)

	for _, item := range items {
		// workaround: feed.ToAtom() panics, if there
		// is no feeds.Link element, so we will add
		// an empty feeds.Link
		if item.Data.Link == nil {
			item.Data.Link = &feeds.Link{}
		}
		feed.Add((*feeds.Item)(&item.Data))
	}

	atomFeedString, err := feed.ToAtom()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write([]byte(atomFeedString))
}

func rssListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var items []structs.AtomFeedItemModel
	database.Con.Find(&items)

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
	newFeedItem := &structs.AtomFeedItemModel{
		Data: structs.AtomFeedItem{
			Title:       "penis",
			Description: "penis2",
			Created:     time.Now(),
		},
	}
	database.Con.Create(newFeedItem)

	respondPlainJSON(&w, []byte(`{"response": "created"}`), 201) // TODO if existent
}

func rssChangeMessageHandler(w http.ResponseWriter, r *http.Request) {
	type respBody struct {
		Done                 bool
		structs.AtomFeedItem // allow all fields to be adjusted (except "Db_Id" of course)
	}

	// TODO document: can be used to mark message as done
	respondPlainJSON(&w, []byte(`{"response": "marked as done"}`), 200)
}

func rssDeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	respondInstance(
		&w,
		fmt.Sprintf("deleted message '%s'", messageToDelete.Data.Title),
		200,
	)
}
