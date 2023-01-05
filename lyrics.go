package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Lyrics struct {
	Artist     string         `json:"artist"`
	Title      string         `json:"title"`
	Lyrics     string         `json:"lyrics"`
	LyricArray []string       `json:"lyric_array"`
	WordMap    map[string]int `json:"word_map"`
	// LyricGraph *dag.Graph     `json:"lyric_graph"`
	// NodeMap    map[string]Drawable
}

type LResponse struct {
	TrackId       int    `xml:"TrackId"`
	Checksum      string `xml:"LyricChecksum"`
	Id            int    `xml:"LyricId"`
	Song          string `xml:"LyricSong"`
	Artist        string `xml:"LyricArtist"`
	Url           string `xml:"LyricUrl"`
	CoverArtUrl   string `xml:"LyricCovertArtUrl"`
	Rank          int    `xml:"LyricRank"`
	CorrectionURL string `xml:"LyricCorrectUrl"`
	Lyrics        string `xml:"Lyric"`
}

func removeNewlines(s string) string {
	f := strings.Replace(s, "\n", " ", -1)
	return f
}

func (l *Lyrics) GetLyricsAsArray() {
	lyrics := removeNewlines(l.Lyrics)
	lyrics_array := strings.Split(lyrics, " ")
	l.LyricArray = lyrics_array
}

func (l *Lyrics) GetWordMap() {
	word_map := make(map[string]int)
	for _, word := range l.LyricArray {
		word_map[word]++
	}
	l.WordMap = word_map
	// l.NodeMap = make(map[string]Drawable)
}

func (l *Lyrics) RetrieveLyrics() {

	params := url.Values{}
	params.Add("artist", l.Artist)
	params.Add("song", l.Title)

	request_url := fmt.Sprintf("http://api.chartlyrics.com/apiv1.asmx/SearchLyricDirect?%s", params.Encode())
	// make request to chartlyrics api
	resp, err := http.Get(request_url)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// parse xml response
	var responseObject LResponse
	xml.Unmarshal(body, &responseObject)

	// set lyrics
	l.Lyrics = responseObject.Lyrics
	l.Artist = responseObject.Artist
	l.Title = responseObject.Song

}

func (l *Lyrics) CreateLyricGraph() {}

func (l *Lyrics) CreateLyricGraphDot(filename string) {}
