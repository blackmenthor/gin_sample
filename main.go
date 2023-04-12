package main

import (
	"fmt"
	pb "github.com/blackmenthor/protobuf-presentation/gin_sample/publish"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResponseType int

const (
	Json     ResponseType = 0
	XML                   = 1
	Protobuf              = 2
	YAML                  = 3
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums []album

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context, responseType ResponseType) {
	switch responseType {
	case Json:
		c.IndentedJSON(http.StatusOK, albums)
	case XML:
		c.XML(http.StatusOK, albums)
	case YAML:
		c.YAML(http.StatusOK, albums)
	case Protobuf:
		//reps := []int64{int64(1), int64(2)}
		//label := "Album"
		//// The specific definition of protobuf is written in the testdata/protoexample file.
		data := &pb.Album{}
		//// Note that data becomes binary data in the response
		//// Will output protoexample.Test protobuf serialized data
		c.ProtoBuf(http.StatusOK, data)
		//c.IndentedJSON(http.StatusOK, albums)
	}
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func initializeData() []album {
	listOfData := []album{}
	var maximumDataSize = 1000000
	for i := 1; i < maximumDataSize; i++ {
		var newAlbum = album{
			ID:     strconv.Itoa(i),
			Title:  fmt.Sprintf("Album %d", i),
			Artist: fmt.Sprintf("Artist %d", i),
			Price:  56.99,
		}
		listOfData = append(listOfData, newAlbum)
	}
	return listOfData
}

func main() {
	albums = initializeData()

	router := gin.Default()
	router.GET("/json/albums", func(c *gin.Context) {
		getAlbums(c, Json)
	})
	router.GET("/xml/albums", func(c *gin.Context) {
		getAlbums(c, XML)
	})
	router.GET("/yaml/albums", func(c *gin.Context) {
		getAlbums(c, YAML)
	})
	router.GET("/proto/albums", func(c *gin.Context) {
		getAlbums(c, Protobuf)
	})
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
