package main

import (
	"context"
	"encoding/json"
	"errors"
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
	Id      string `form:"id" json:"id" xml:"id" yaml:"id" redis:"id"`
	Content string `form:"content" json:"content" xml:"content" yaml:"content" redis:"content"`
}

type Env struct {
	redis *redis.Client
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

var acceptedFormats = map[string]bool{
	"JSON": true,
	"YML":  true,
	"YAML": true,
	"XML":  true,
}

var ctx = context.Background()

func main() {
	env := &Env{
		redis: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
			DB:   0,
		}),
	}

	router := gin.Default()
	router.POST("/create-message", env.createMessage)
	router.GET("/topic/:topic/limit/:limit/token/:token/format/:format", env.readMessage)
	router.Run("0.0.0.0:9999")
}

func (env *Env) createMessage(c *gin.Context) {
	token := c.Query("token")
	topic := c.Query("topic")
	message := Message{}
	c.Bind(&message)
	for _, user := range users {
		if token == user.Token {
			commonFormat, err := json.Marshal(message)
			if err != nil {
				log.Fatal(err)
			}
			redisErr := env.redis.RPush(ctx, topic, commonFormat).Err()
			if redisErr != nil {
				log.Fatal(redisErr)
			}
			return
		}
	}
	c.JSON(http.StatusForbidden, gin.H{"info": "Invalid token"})
}

func (env *Env) readMessage(c *gin.Context) {
	consumerToken := c.Param("token")
	topic := c.Param("topic")
	format := c.Param("format")
	msgLimit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		log.Fatal(err)
	}

	if !acceptedFormats[format] {
		keys := []string{}
		for format := range acceptedFormats {
			keys = append(keys, format)
		}
		c.JSON(http.StatusBadRequest, gin.H{"info": fmt.Sprintf("ESB only accepts: %s", keys)})
		return
	}
	if msgLimit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"info": "Limit is 0 or less"})
		return
	}
	for _, v := range users {
		if v.Token == consumerToken {
			// transform message based on consumer
			redisMessages := env.redis.LRange(ctx, topic, 0, int64(msgLimit)-1)
			messages := []Message{}
			for _, message := range redisMessages.Val() {
				outputMessage := Message{}
				err := json.Unmarshal([]byte(message), &outputMessage)
				if err != nil {
					log.Fatal(err)
				}

				messages = append(messages, outputMessage)
			}

			switch format {
			case "XML":
				c.XML(http.StatusOK, messages)
			case "JSON":
				c.JSON(http.StatusOK, messages)
			case "YML", "YAML":
				c.YAML(http.StatusOK, messages)
			}
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"info": "Invalid token"})
}

func transformMessage(c *gin.Context, message Message, format string) ([]byte, error) {
	switch format {
	case "JSON":
		return JSONTransformer(message)
	case "XML":
		return XMLTransformer(message)
	case "YML", "YAML":
		return YMLTransformer(message)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"info": "Invalid format"})
	}

	return nil, errors.New("Did not get an expected message format")
}
