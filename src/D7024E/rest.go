package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type msg struct {
	Data string `json:"data"`
}

// getAlbums responds with the list of all albums as JSON.
func getObject(c *gin.Context) {
	hash := c.Param("hash")
	fmt.Println("[REST] loading", hash)

	var net = NewNetwork("127.0.0.1:54321")
	data := net.SendLocalGet(hash)

	var newMsg msg
	newMsg.Data = string(data)

	c.IndentedJSON(http.StatusOK, newMsg)
}

func postData(c *gin.Context) {
	var net = NewNetwork("127.0.0.1:54321")
	var newMsg msg

	data, _ := c.GetRawData()
	newMsg.Data = string(data)
	fmt.Println("[REST] storing", data)

	//TODO: use newMsg to store data, get hash

	hash, err := net.SendLocalPut(data)
	if err == nil {
		c.Writer.Header().Set("Location", "/objects/"+hash)

		// Add the new album to the slice.
		c.IndentedJSON(http.StatusCreated, newMsg)
	}

}

func RestApi() {
	router := gin.Default()
	router.GET("/objects/:hash", getObject)
	router.POST("/objects", postData)

	router.Run("0.0.0.0:8080")
}
