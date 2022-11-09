package server

import (
	"fmt"
	"net/http"
	"status-api/database"
	"status-api/structs"
	"time"

	"github.com/gorilla/feeds"
)

func rssShowHandler(w http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title:       "Status API messages for leon.wtf",
		Link:        &feeds.Link{Href: "https://leon.wtf"},
		Description: "Messages about incidents, maintenance, and others for the public services running on leon.wtf",
		Author:      &feeds.Author{Name: "Leon Schmidt", Email: "admin@leon.wtf"},
		Created:     time.Now(),
	}

	var items []structs.AtomFeedItemModel
	database.Con.Find(&items)

	fmt.Println(items)

	feed.Items = []*feeds.Item{
		{
			Title:       "Limiting Concurrency in Go",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Description: "A discussion on controlled parallelism in golang",
			Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     time.Now(),
		},
		{
			Title:       "Logic-less Template Redux",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
			Description: "More thoughts on logicless templates",
			Created:     time.Now(),
		},
		{
			Title:       "Idiomatic Code Reuse in Go",
			Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
			Description: "How to use interfaces <em>effectively</em>",
			Created:     time.Now(),
		},
	}

	atom, err := feed.ToAtom()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write([]byte(atom))
}

func rssListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var items []structs.AtomFeedItemModel
	database.Con.Find(&items)
	// TODO das mal richtig machen
	respondInstance(&w, items, 200)
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

	respondJSON(&w, []byte(`{"response": "created"}`), 200)
}

func rssDeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(&w, []byte(`{"response": "deleted"}`), 200)
}
