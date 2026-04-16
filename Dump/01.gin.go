package main

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/getUser", getUser)
	router.POST("/submitForm", submitForm)
	router.GET("/getQuery", getRouteQueries)
	router.GET("/getParam/:name/:age", getRouteParams)

	router.Run()
}

func getUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": "Hello World...",
	})
}

func submitForm(c *gin.Context) {
	IP := c.ClientIP()
	remoteIP := c.RemoteIP()

	header := c.Request.Header
	method := c.Request.Method

	body := c.Request.Body // Body becomes empty {} because it's a stream — once read, it's consumed.
	value, _ := io.ReadAll(body)
	convertedValue := string(value)
	var jsonData map[string]interface{}
	_ = json.Unmarshal(value, &jsonData)

	c.JSON(200, gin.H{
		"IP":              IP,
		"RemoteIP":        remoteIP,
		"Header":          header,
		"Method":          method,
		"Body":            body,
		"Value":           value,
		"Converted Value": convertedValue,
		"JSON Data":       jsonData,
	})
}

func getRouteQueries(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")

	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}

func getRouteParams(c *gin.Context) {
	name := c.Param("name")
	age := c.Param("age")

	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}
