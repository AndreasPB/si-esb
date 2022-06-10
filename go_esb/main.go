package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type Message struct {
	Id      string `form:"id" json:"id" xml:"id" yaml:"id" redis:"id"`
	Content string `form:"content" json:"content" xml:"content" yaml:"content" redis:"content"`
	Exp     int64  `form:"exp" json:"exp" xml:"exp" yaml:"exp" redis:"exp"`
}

type Env struct {
	redis *redis.Client
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
	router.GET("/topic/:topic/skip/:skip/limit/:limit/format/:format", env.readMessage)
	router.GET("/cleanup", env.handleMessageExpiration)
	router.Run("0.0.0.0:9999")
}

<<<<<<< HEAD
func (env *Env) handleMessageExpiration(c *gin.Context) {
	message := Message{}

	for _, topic := range env.redis.Keys(ctx, "*").Val() {
		fmt.Println("Cleaning:", topic)
		for {
			firstIndex := env.redis.LIndex(ctx, topic, 0)

			err := json.Unmarshal([]byte(firstIndex.Val()), &message)
			if err != nil {
				fmt.Println("No messages")
				break
			}

			if message.Exp < time.Now().Unix() {
				fmt.Println("Expired", message.Exp)
				poppedMessage := env.redis.LPop(ctx, topic)
				fmt.Println("Popped message:", poppedMessage)
			} else {
				fmt.Println("No expired messages")
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"info": "Done cleaning up"})
}
func verifyAuth(c *gin.Context) {
	auth := c.GetHeader("auth")
	secret := "6hest9"

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
}

func (env *Env) createMessage(c *gin.Context) {
	topic := c.Query("topic")
	auth := c.GetHeader("auth")
	message := Message{}
	c.Bind(&message)

<<<<<<< HEAD
	expirationTime := time.Hour * 24
	message.Exp = time.Now().Add(expirationTime).Unix()

	commonFormat, err := json.Marshal(message)
=======
	secret := "6hest9"

	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
>>>>>>> 5f184b1 (Add JWT verification for create_message)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if token.Valid {
		commonFormat, err := json.Marshal(message)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if len(message.Content) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"info": "Wrongly formatted message"})
			return
		}

		redisErr := env.redis.RPush(ctx, topic, commonFormat).Err()
		if redisErr != nil {
			log.Fatal(redisErr)
		}
		c.JSON(http.StatusOK, gin.H{"info": "Message was saved", "message": message})
		return
	} else {
		fmt.Println("Couldn't handle this token:", err)
		panic("Some JWT thing happened that we did not handle")
	}
}

func (env *Env) readMessage(c *gin.Context) {
	topic := c.Param("topic")
	format := c.Param("format")
	skip, err := strconv.Atoi(c.Param("skip"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"info": "skip must be an integer"})
		return
	}
	msgLimit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"info": "limit be must an integer"})
		return
	}

	if !acceptedFormats[format] {
		keys := []string{}
		for format := range acceptedFormats {
			keys = append(keys, format)
		}
		c.JSON(http.StatusBadRequest, gin.H{"info": fmt.Sprintf("ESB only accepts: %s", keys)})
		return
	}

	if skip < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"info": "skip offset must be 0 or above"})
	}
	if msgLimit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"info": "limit must be 1 or above"})
		return
	}

	redisMessages := env.redis.LRange(ctx, topic, int64(skip), int64(msgLimit)-1)
	messages := []Message{}
	for _, message := range redisMessages.Val() {
		outputMessage := Message{}
		err := json.Unmarshal([]byte(message), &outputMessage)
		if err != nil {
			log.Fatal(err)
		}

		messages = append(messages, outputMessage)
	}

	// transform message based on consumer
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
