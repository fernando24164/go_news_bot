package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"

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

// Parse the feed data into the struct tags implemented before
func feedreader(data []byte) rss {
	feed := rss{}
	xmlreader := xml.NewDecoder(bytes.NewReader(data))
	err := xmlreader.Decode(&feed)
	if err != nil {
		log.Fatal(err)
	}
	return feed
}

// Get the rss data through http request
func getfeed(url string) []byte {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func dbsaver(feed rss) {
	db, err := sql.Open("sqlite3", "./feedb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// TABLE structure : create table items(item text, link text, pubdate integer, channel text)
	stmt, err := tx.Prepare("INSERT INTO items(item, link, pubdate, channel) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, item := range feed.Channel.Items {
		_, err = stmt.Exec(item.Title, item.Link, item.PubDate, feed.Channel.Title)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func main() {
	//url := "http://feeds.weblogssl.com/genbeta"
	url := "http://www.eldiario.es/rss/"
	data := getfeed(url)
	feed := feedreader(data)
	dbsaver(feed)
}
