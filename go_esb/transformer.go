package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONTransformer(c *gin.Context, message Message) {
	c.JSON(http.StatusOK, message)
}

func XMLTransformer(c *gin.Context, message Message) {
	c.XML(http.StatusOK, message)
}

func YMLTransformer(c *gin.Context, message Message) {
	c.YAML(http.StatusOK, message)
}

func TSVTransformer() {

}
