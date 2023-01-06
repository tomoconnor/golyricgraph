package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type GraphRequest struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type GraphResponse struct {
	Artist      string `json:"artist"`
	Title       string `json:"title"`
	PNGFilename string `json:"png_filename"`
	SVGFilename string `json:"svg_filename"`
}
type AcceptLyricsRequest struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Lyrics string `json:"lyrics"`
}

func AcceptLyrics(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		a := AcceptLyricsRequest{}
		filename := uuid.New().String()

		if err := c.Bind(&a); err != nil {
			log.Println("Error binding request")
			return err
		}

		LyricGenerator := Lyrics{
			Artist: a.Artist,
			Title:  a.Title,
			Lyrics: a.Lyrics,
		}
		LyricGenerator.GetLyricsAsArray()
		LyricGenerator.GetWordMap()
		LyricGenerator.CreateLyricGraph(filename)
		graphResponse := GraphResponse{
			Artist:      a.Artist,
			Title:       a.Title,
			PNGFilename: filename + ".png",
			SVGFilename: filename + ".svg",
		}
		dbo := LyricGraph{
			Artist:   a.Artist,
			Title:    a.Title,
			Filename: filename + ".png",
		}

		db.Create(&dbo)
		return c.JSON(200, graphResponse)

	}
}

func GraphLyrics(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// POST convert to GraphRequest
		var graphRequest GraphRequest
		if err := c.Bind(&graphRequest); err != nil {
			log.Println("Error binding request")
			return err
		}
		log.Println("Artist: ", graphRequest.Artist)
		log.Println("Title: ", graphRequest.Title)
		filename := uuid.New().String()

		LyricGenerator := Lyrics{
			Artist: graphRequest.Artist,
			Title:  graphRequest.Title,
		}
		err := LyricGenerator.RetrieveLyrics()
		if err != nil {
			log.Println("Error retrieving lyrics")
			return c.JSON(500, "Error retrieving lyrics")
		}

		LyricGenerator.GetLyricsAsArray()
		LyricGenerator.GetWordMap()
		LyricGenerator.CreateLyricGraph(filename)

		graphResponse := GraphResponse{
			Artist:      graphRequest.Artist,
			Title:       graphRequest.Title,
			PNGFilename: filename + ".png",
			SVGFilename: filename + ".svg",
		}
		dbo := LyricGraph{
			Artist:   graphRequest.Artist,
			Title:    graphRequest.Title,
			Filename: filename + ".png",
		}

		db.Create(&dbo)
		return c.JSON(200, graphResponse)

	}
}

func CompareLyrics(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(501, "Not Implemented Yet")
	}
}

func GetGraph(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(200, "Graph Lyrics")
	}
}

func StartServer(db *gorm.DB) {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatal("HTTP_PORT not set")
	}
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	p := prometheus.NewPrometheus("lyrics", nil)
	p.Use(e)

	e.POST("/api/v1/graph", GraphLyrics(db))
	e.POST("/api/v1/compare", CompareLyrics(db))
	e.GET("/api/v1/retrieve", GetGraph(db))
	e.POST("/api/v1/accept", AcceptLyrics(db))

	e.Logger.Fatal(e.Start(":" + httpPort))

}
