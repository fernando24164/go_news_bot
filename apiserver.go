package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Get the items from de database with the flag for only not sent ones or everything
func getitems(onlynew bool) string {
	rows, err := db.Query("SELECT * FROM items")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var response string
	for rows.Next() {
		var item, link, pubdate, channel string
		var sent int
		rows.Scan(&item, &link, &pubdate, &channel, &sent)
		response = response + "Item: " + item + " - Link : " + link + " - Publication date : " + pubdate + " - Channel : " + channel + "&&"
	}
	return response
}

func itemsresponse(w http.ResponseWriter, r *http.Request) {
	//user := r.URL.Path[len("/getfeed/"):]
	//fmt.Println(user)
	//getitems(false)
	fmt.Fprintf(w, getitems(false), "")
}

// API function to query the database
func apiserver() {
	http.HandleFunc("/getfeed/", itemsresponse)
	http.ListenAndServe(":9096", nil)
}
