package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
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
		fmt.Printf("error: %v", err)
	}
	return feed
}

// Get the rss data through http request
func getfeed(url string) []byte {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return data
}

// TODO use log instead of prints
func main() {
	//url := "http://feeds.weblogssl.com/genbeta"
	url := "http://www.eldiario.es/rss/"
	data := getfeed(url)
	feed := feedreader(data)
	for _, item := range feed.Channel.Items {
		fmt.Println(item)
	}
}
