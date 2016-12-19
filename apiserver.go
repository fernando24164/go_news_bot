package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type jList struct {
	items []jItem
}

type jItem struct {
	Title   string
	Link    string
	PubDate string
}

// Get the items from de database with the flag for only not sent ones or everything
func getitems(onlynew bool) []byte {
	rows, err := db.Query("SELECT * FROM items")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var responseList []jItem
	for rows.Next() {
		var item, link, pubdate, channel string
		var sent int
		rows.Scan(&item, &link, &pubdate, &channel, &sent)
		responseList = append(responseList, jItem{item, link, pubdate})
	}
	response := jList{responseList}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	return jsonResp
}

func itemsresponse(w http.ResponseWriter, r *http.Request) {
	//user := r.URL.Path[len("/getfeed/"):]
	//fmt.Println(user)
	//getitems(false)
	fmt.Println(getitems(false))
}

// API function to query the database
func apiserver() {
	http.HandleFunc("/getfeed/", itemsresponse)
	http.ListenAndServe(":9096", nil)
}
