package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Steps:
// - Pass message to the ESB, (XML, JSON, TSV)
// - Register users and provide them with a token,
//   ESB will transform the message depending on who is consuming it
// - Consumer needs to identify itself, this is done via token
// - Consumer informs ESB to get messages from a specific provider

type User struct {
	Id    string `json:"id" xml:"id"`
	Token string `json:"token" xml:"token"`
}

type Message struct {
	Id      string `form:"id" json:"id" xml:"id" yaml:"id"`
	Content string `form:"content" json:"content" xml:"content" yaml:"content"`
}

var users = []User{
	{
		Id:    "1",
		Token: "1212",
	},
	{
		Id:    "2",
		Token: "3333",
	},
}

var messages = map[string][]Message{
	"1": {
		{Id: "1a76658a-7e4c-4f24-9a96-f68ef3526008", Content: "I am message 1"},
		{Id: "87a17661-13ae-45c3-bfb2-041afa15234a", Content: "I am message 2"},
		{Id: "aa229257-48f3-46e2-88dc-beb210e7f9e4", Content: "I am message 3"},
		{Id: "410cc421-71c0-414a-99d6-6e603e741692", Content: "I am message 4"},
	},
}

var ctx = context.Background()

func main() {
	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	router := gin.Default()
	router.POST("/create-message", messageHandler)
	router.GET("/provider/:id/limit/:limit/token/:token", esb)
	router.Run("0.0.0.0:9999")
}

func messageHandler(c *gin.Context) {
	token := c.Query("token")
	format := c.Query("format")
	message := Message{}
	c.Bind(&message)
	for _, user := range users {
		if token == user.Token {
			transformMessage(c, message, format)
			return
		}
	}
	c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid token"})
}

func esb(c *gin.Context) {
	consumerToken := c.Param("token")
	//msgId := c.Param("id")
	msgLimit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		log.Fatal("Could not convert limit to int")
	}

	if msgLimit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"info": "Limit is 0 or less"})
		return
	}
	for _, v := range users {
		if v.Token == consumerToken {
			// transform message based on consumer
			// TODO: handle properly
			//transformMessage(c, consumerToken, message)
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"info": "Invalid token"})
}

func transformMessage(c *gin.Context, message Message, format string) {
	switch format {
	case "JSON":
		JSONTransformer(c, message)
	case "XML":
		XMLTransformer(c, message)
	case "YML", "YAML":
		YMLTransformer(c, message)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid format"})
	}
}
