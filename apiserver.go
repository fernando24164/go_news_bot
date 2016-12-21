package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Struc to define item's json object
type jItem struct {
	Title, Link, PubDate string
}

// Get the items from de database with the flag for only not sent ones or everything
func getitems(onlynew bool) []byte {
	rows, err := db.Query("SELECT * FROM items LIMIT 0,2")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var responseList []jItem
	for rows.Next() {
		var item, link, pubdate, channel string
		var sent int
		rows.Scan(&item, &link, &pubdate, &channel, &sent)
		responseList = append(responseList, jItem{Title: item, Link: link, PubDate: pubdate})
	}
	jsonResp, err := json.Marshal(responseList)
	if err != nil {
		log.Fatal(err)
	}
	return jsonResp
}

// Print items in json
func itemsresponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", getitems(false))
}

// API function to query the database
func apiserver() {
	http.HandleFunc("/getfeed/", itemsresponse)
	http.ListenAndServe(":9096", nil)
}
