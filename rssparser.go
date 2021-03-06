package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// RSS xml structure to parse it with the xml lib
type item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	PubDate string   `xml:"pubDate"`
}
type channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Items   []item   `xml:"item"`
}
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel *channel `xml:"channel"`
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Parse the feed data into the struct tags implemented before
func feedreader(data []byte) rss {
	feed := rss{}
	xmlreader := xml.NewDecoder(bytes.NewReader(data))
	err := xmlreader.Decode(&feed)
	checkError(err)
	return feed
}

// Get the rss data through http request
func getfeed(url string) []byte {
	res, err := http.Get(url)
	defer res.Body.Close()
	checkError(err)
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// Save all the items get from the subcriptions on the database
func dbsaver(feed rss) {
	tx, err := db.Begin()
	checkError(err)
	// TABLE structure : create table items(item text, link text, pubdate integer, channel text, sent integer)
	stmt, err := tx.Prepare("INSERT INTO items(item, link, pubdate, channel, sent) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, item := range feed.Channel.Items {
		_, err := stmt.Exec(item.Title, item.Link, item.PubDate, feed.Channel.Title, 0)
		checkError(err)
	}
	tx.Commit()
}

// Search for the list of subcriptions for a user
func getfeedsubcriptions(user string) []string {
	// TABLE structure : create table subcriptions(id int, name text, list text)
	// INSERT : INSERT INTO subcriptions(id, name, list) VALUES("1", "agpatag", "http://feeds.weblogssl.com/genbeta,http://www.eldiario.es/rss")
	var list string
	err := db.QueryRow("SELECT list FROM subcriptions WHERE name = ?", user).Scan(&list)
	checkError(err)
	result := make([]string, 0)
	for _, source := range strings.Split(list, ",") {
		result = append(result, source)
	}
	return result
}

// Sleep a search for the feed
func scheduler(sleeper int, user string) {
	fmt.Println(" # Starting the schedule . . .")
	for {
		fmt.Println(" > Repolling Subcriptions!")
		subcriptions := getfeedsubcriptions(user)
		for _, url := range subcriptions {
			data := getfeed(url)
			feed := feedreader(data)
			dbsaver(feed)
		}
		time.Sleep(time.Duration(sleeper) * time.Second)
	}
}

// connection with the database as a global variable
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./feedb.db")
	checkError(err)
	defer db.Close()
	sleeper := 60
	user := "agpatag"
	go scheduler(sleeper, user)
	apiserver()
	// getitems(false)
	// urls := getfeedsubcriptions(user)
}
