package main

import (
	"reflect"
	"testing"
)

func TestLyricGet(t *testing.T) {
	lyric := Lyrics{
		Artist: "The Beatles",
		Title:  "Yellow Submarine",
	}
	lyric.RetrieveLyrics()
	expected := "In the town where I was born\nLived a man who sailed to sea\nAnd he told us of his life\nIn the land of submarines\n\nSo we sailed on to the sun\nTill we found the sea of green\nAnd we lived beneath the waves\nIn our yellow submarine\n\nWe all live in a yellow submarine\nYellow submarine, yellow submarine\nWe all live in a yellow submarine\nYellow submarine, yellow submarine\n\nAnd our friends are all aboard\nMany more of them live next door\nAnd the band begins to play\n\nWe all live in a yellow submarine\nYellow submarine, yellow submarine\nWe all live in a yellow submarine\nYellow submarine, yellow submarine\n\nFull steam ahead, Mister Boatswain, full steam ahead\nFull steam ahead it is, Sergeant\nCut the cable! Drop the cable!\nAye-aye, sir, aye-aye\nCaptain! Captain!\n\n[Ringo Starr & *Paul McCartney*]\nAs we live a life of ease\nEvery one of us (*Every one of us*) has all we need (*Has all we need*)\nSky of blue (*Sky of blue*) and sea of green (*Sea of green*)\nIn our yellow (*In our yellow) submarine (*Submarine, ha-ha!*)\n\nWe all live in a yellow submarine\nYellow submarine, yellow submarine\nWe all live in a yellow submarine\nYellow submarine, yellow submarine\nWe all live in a yellow submarine\nYellow submarine, yellow submarine"
	if lyric.Lyrics != expected {
		t.Errorf("Lyrics were incorrect, got: %s, want: %s.", lyric.Lyrics, expected)
	}
}

func TestRemoveNewLines(t *testing.T) {
	s := "hello\nworld"
	expected := "helloworld"
	if removeNewlines(s) != expected {
		t.Errorf("Newlines were not removed, got: %s, want: %s.", s, expected)
	}
}

func TestLyricArray(t *testing.T) {
	lyric := Lyrics{
		Lyrics: "A\nB\nC\nD\nE\nF\nG\nH\nI\nJ\nK\nL\nM\nN\nO\nP\nQ\nR\nS\nT\nU\nV\nW\nX\nY\nZ",
	}
	// lyric.RetrieveLyrics()
	expected := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	if reflect.DeepEqual(lyric.GetLyricsAsArray(), expected) {
		t.Errorf("Lyrics were incorrect, got: %s, want: %s.", lyric.GetLyricsAsArray(), expected)
	}
}
