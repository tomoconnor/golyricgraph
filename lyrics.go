package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	dag "github.com/hashicorp/terraform/dag"
)

type DrawableOrigin struct {
	VertexName string
}

type Drawable struct {
	VertexName string
}

func (node *Drawable) Name() string {
	return node.VertexName
}

func (node *Drawable) DotNode(n string, opts *dag.DotOpts) *dag.DotNode {
	return &dag.DotNode{
		Name:  n,
		Attrs: map[string]string{},
	}
}

// type GraphVertex struct {
// 	DotNodeTitle  string
// 	DotNodeOpts   *dag.DotOpts
// 	DotNodeReturn *dag.DotNode
// }

// func (v *GraphVertex) MakeDotNode(title string, opts *dag.DotOpts) *dag.DotNode {
// 	v.DotNodeTitle = title
// 	v.DotNodeOpts = opts
// 	v.DotNodeReturn = &dag.DotNode{
// 		Name: title,
// 		// Attrs: map[string]string{
// 		// 	"label": title,
// 		// },
// 	}

// 	return v.DotNodeReturn
// }

type Lyrics struct {
	Artist     string         `json:"artist"`
	Title      string         `json:"title"`
	Lyrics     string         `json:"lyrics"`
	LyricArray []string       `json:"lyric_array"`
	WordMap    map[string]int `json:"word_map"`
	LyricGraph *dag.Graph     `json:"lyric_graph"`
	NodeMap    map[string]Drawable
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
	l.NodeMap = make(map[string]Drawable)
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

func (l *Lyrics) CreateLyricGraph() {
	g := dag.Graph{}
	o := DrawableOrigin{VertexName: "Song Start"}
	g.Add(&o)

	for _, word := range l.LyricArray {
		w := Drawable{VertexName: word}
		// w := dag.DotNode{
		// 	Name: word,
		// }

		g.Add(w)
		l.NodeMap[word] = w
	}
	for i, word := range l.LyricArray {
		if i < len(l.LyricArray)-1 {
			a := l.NodeMap[word]
			b := l.NodeMap[l.LyricArray[i+1]]
			g.Connect(dag.BasicEdge(a, b))
		}
	}
	g.Connect(dag.BasicEdge(&o, l.NodeMap[l.LyricArray[0]]))

	// 		g.Connect(dag.BasicEdge(word, l.LyricArray[i+1]))
	// 	}
	// }
	l.LyricGraph = &g
}

func (l *Lyrics) CreateLyricGraphDot(filename string) {
	g := l.LyricGraph
	dotOptions := dag.DotOpts{
		Verbose:    false,
		DrawCycles: false,
		MaxDepth:   10,
	}

	dot := g.Dot(&dotOptions)
	// fmt.Println(dot) // Meaningless fucking dot file
	// output to file.
	err := ioutil.WriteFile(filename, dot, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}
