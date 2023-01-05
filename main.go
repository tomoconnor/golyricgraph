package main

import (
	// "github.com/getsentry/sentry-go"
	// sentryecho "github.com/getsentry/sentry-go/echo"

	"flag"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type LyricGraph struct {
	gorm.Model
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
}

func main() {
	log.Println("Lyric Grapher")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	db.AutoMigrate(&LyricGraph{})

	artistName := flag.String("artist", "", "Artist Name")
	songName := flag.String("song", "", "Song Name")
	serverFlag := flag.Bool("server", false, "Run as server")
	flag.Parse()
	if *serverFlag {
		log.Println("Running as server")
		StartServer(db)
	} else {
		log.Println("Artist: ", *artistName)
		log.Println("Song: ", *songName)

		if *artistName == "" || *songName == "" {
			log.Fatal("Please provide an artist and song name")
			return
		}
		filename := *artistName + "_" + *songName

		LyricGenerator := Lyrics{
			Artist: *artistName,
			Title:  *songName,
		}
		LyricGenerator.RetrieveLyrics()
		LyricGenerator.GetLyricsAsArray()
		LyricGenerator.GetWordMap()
		LyricGenerator.CreateLyricGraph(filename)
		dbo := LyricGraph{
			Artist:   *artistName,
			Title:    *songName,
			Filename: filename,
		}
		db.Create(&dbo)

	}
}

// type LyricGraphComparison struct {
// 	gorm.Model
// 	artist1  string `json:"artist1"`
// 	artist2  string `json:"artist2"`
// 	title1   string `json:"title1"`
// 	title2   string `json:"title2"`
// 	filename string `json:"filename"`
// 	filepath string `json:"filepath"`
// 	rank     int    `json:"rank"`
// }
