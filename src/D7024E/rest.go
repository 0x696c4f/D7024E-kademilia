package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type msg struct {
    Data     string  `json:"data"`
}
// getAlbums responds with the list of all albums as JSON.
func getObject(c *gin.Context) {
	hash := c.Param("hash")
	//TODO: load data
	data := "data with hash "+hash

    var newMsg msg
	newMsg.Data = string(data)

    c.IndentedJSON(http.StatusOK, newMsg)
}

func postData(c *gin.Context) {
    var newMsg msg


	data,_:=c.GetRawData()
	newMsg.Data = string(data)

	//TODO: use newMsg to store data, get hash
	hash:= "abcd"

    c.Writer.Header().Set("Location","/objects/"+hash)

    // Add the new album to the slice.
    c.IndentedJSON(http.StatusCreated, newMsg)
}

func RestApi() {
    router := gin.Default()
	router.GET("/objects/:hash", getObject)
	router.POST("/objects",postData)

    router.Run("localhost:8080")
}