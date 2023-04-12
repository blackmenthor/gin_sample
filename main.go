package main

import (
	"context"
	"fmt"
	pb "github.com/blackmenthor/gin_sample/publish"
	pbt "github.com/blackmenthor/gin_sample/tutorial"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
)

type ResponseType int

const (
	Json         ResponseType = 0
	XML                       = 1
	Protobuf                  = 2
	Protobuf_new              = 3
	YAML                      = 4
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
var albumsProto []*pb.ListOfAlbums_Album

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
		listOfAlbum := &pb.ListOfAlbums{
			Albums: albumsProto,
		}
		c.ProtoBuf(http.StatusOK, listOfAlbum)
	case Protobuf_new:
		kb := 1024
		mb := 1024 * kb
		gb := 1024 * mb

		conn, err := grpc.Dial(
			"localhost:9000",
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1*gb)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}
		defer conn.Close()

		// Code removed for brevity

		client := pbt.NewAlbumServiceClient(conn)

		// Note how we are calling the GetBookList method on the server
		// This is available to us through the auto-generated code
		albumz, errz := client.GetAlbum(context.Background(), &pbt.AlbumRequest{})

		if errz != nil {
			log.Printf("grpc error %v", errz)
			c.ProtoBuf(http.StatusInternalServerError, errz)
		}

		resp := albumz.GetAlbums()

		fmt.Printf("test albums %v", len(resp))
		c.ProtoBuf(http.StatusOK, albumz)
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
	var listOfData []album
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

func initializeDataProto() []*pb.ListOfAlbums_Album {
	var maximumDataSize = 1000000
	data := make([]*pb.ListOfAlbums_Album, maximumDataSize)
	for i := 1; i < maximumDataSize; i++ {
		var newAlbum = &pb.ListOfAlbums_Album{
			Id:     strconv.Itoa(i),
			Title:  fmt.Sprintf("Album %d", i),
			Artist: fmt.Sprintf("Artist %d", i),
			Price:  56.99,
		}
		data = append(data, newAlbum)
	}
	return data
}

func main() {
	albums = initializeData()
	albumsProto = initializeDataProto()

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
	router.GET("/proto-2/albums", func(c *gin.Context) {
		getAlbums(c, Protobuf_new)
	})
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run(":8081")
}
