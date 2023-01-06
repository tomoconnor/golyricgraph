package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Lyrics struct {
	Artist     string         `json:"artist"`
	Title      string         `json:"title"`
	Lyrics     string         `json:"lyrics"`
	LyricArray []string       `json:"lyric_array"`
	WordMap    map[string]int `json:"word_map"`
	LyricGraph *graphviz.Graphviz
	NodeMap    map[string]*cgraph.Node
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

func TidyUpLyrics(s string) string {
	f := strings.Replace(s, "\n", " ", -1)
	f = strings.Replace(f, ",", "", -1)
	f = strings.Replace(f, ".", "", -1)
	f = strings.Replace(f, "!", "", -1)
	f = strings.Replace(f, "?", "", -1)
	f = strings.Replace(f, "(", "", -1)
	f = strings.Replace(f, ")", "", -1)
	f = strings.Replace(f, "[", "", -1)
	f = strings.Replace(f, "]", "", -1)
	f = strings.Replace(f, "{", "", -1)
	f = strings.Replace(f, "}", "", -1)
	f = strings.Replace(f, ":", "", -1)
	f = strings.Replace(f, ";", "", -1)
	f = strings.Replace(f, "-", "", -1)
	f = strings.Replace(f, "_", "", -1)
	f = strings.Replace(f, "/", "", -1)
	f = strings.Replace(f, "\"", "", -1)
	f = strings.Replace(f, "'", "", -1)
	f = strings.Replace(f, "\n", " ", -1)
	f = strings.Replace(f, "\n", " ", -1)
	f = strings.Replace(f, "\n", " ", -1)
	f = strings.Replace(f, "\t", " ", -1)

	f = strings.ToLower(f)
	return f
}

func (l *Lyrics) GetLyricsAsArray() {
	lyrics := TidyUpLyrics(l.Lyrics)
	lyrics_array := strings.Split(lyrics, " ")
	l.LyricArray = lyrics_array
}

func (l *Lyrics) GetWordMap() {
	word_map := make(map[string]int)
	for _, word := range l.LyricArray {
		word_map[word]++
	}
	l.WordMap = word_map
	l.NodeMap = make(map[string]*cgraph.Node)
}

func (l *Lyrics) RetrieveLyrics() error {

	params := url.Values{}
	params.Add("artist", l.Artist)
	params.Add("song", l.Title)

	request_url := fmt.Sprintf("http://api.chartlyrics.com/apiv1.asmx/SearchLyricDirect?%s", params.Encode())
	// make request to chartlyrics api
	resp, err := http.Get(request_url)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	// parse xml response
	var responseObject LResponse
	xml.Unmarshal(body, &responseObject)

	// set lyrics
	l.Lyrics = responseObject.Lyrics
	l.Artist = responseObject.Artist
	l.Title = responseObject.Song
	return nil
}

func (l *Lyrics) CreateLyricGraph(filename string) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()
	// Create Start Node
	startNode, err := graph.CreateNode("START")
	if err != nil {
		log.Fatal(err)
	}
	startNode.SetLabel("START")

	// Create End Node
	endNode, err := graph.CreateNode("END")
	if err != nil {
		log.Fatal(err)
	}
	endNode.SetLabel("END")

	// create nodes
	for word, count := range l.WordMap {
		node, err := graph.CreateNode(word)
		if err != nil {
			log.Fatal(err)
		}
		node.SetLabel(fmt.Sprintf("%s (%d)", word, count))
		l.NodeMap[word] = node
	}

	// create edges
	for i := 0; i < len(l.LyricArray)-1; i++ {
		from := l.NodeMap[l.LyricArray[i]]
		to := l.NodeMap[l.LyricArray[i+1]]

		_, err := graph.CreateEdge("", from, to)
		if err != nil {
			log.Fatal(err)
		}
		// edge.SetLabel("1")

	}
	// join start node to first word
	_, err = graph.CreateEdge("", startNode, l.NodeMap[l.LyricArray[0]])
	if err != nil {
		log.Fatal(err)
	}
	// join last word to end node
	_, err = graph.CreateEdge("", l.NodeMap[l.LyricArray[len(l.LyricArray)-1]], endNode)
	if err != nil {
		log.Fatal(err)
	}

	l.LyricGraph = g
	imageFileName := fmt.Sprintf("/srv/images/imagefiles/%s.png", filename)
	if err := g.RenderFilename(graph, graphviz.PNG, imageFileName); err != nil {
		log.Fatal(err)
	}
	svgFileName := fmt.Sprintf("/srv/images/imagefiles/%s.svg", filename)
	if err := g.RenderFilename(graph, graphviz.SVG, svgFileName); err != nil {
		log.Fatal(err)
	}
	dotFileName := fmt.Sprintf("/srv/images/dotfiles/%s.dot", filename)
	if err := g.RenderFilename(graph, graphviz.XDOT, dotFileName); err != nil {
		log.Fatal(err)
	}

}
