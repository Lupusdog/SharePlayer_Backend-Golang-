package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	type ChatList struct {
		Name    string
		Comment string
	}

	type Chat struct {
		Comment string `json:"comment" binding:"required"`
	}

	type MovieData struct {
		Url  string  `json:"url"`
		Time float64 `json:"time"`
	}

	type UserName struct {
		Name string `json:"name" binding:"required"`
	}

	var Url string
	var Time float64
	var ChatLine []ChatList = []ChatList{{"Lupusdog", "Let's Chatting!!"}}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/share", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"url":  Url,
			"time": Time,
		})
	})

	router.POST("/share", func(c *gin.Context) {
		var json MovieData
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		Url = json.Url
		Time = json.Time
	})

	router.POST("/user", func(c *gin.Context) {
		var json UserName
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.SetCookie("user", json.Name, 7200, "/", "https://shareplayer-backend.herokuapp.com", true, false)

	})

	router.POST("/chat", func(c *gin.Context) {
		var json Chat

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		Name, json_err := c.Cookie("user")
		if json_err != nil {
			Name = "Guest"
		}

		SendedChat := ChatList{Name, json.Comment}
		ChatLine = append(ChatLine, SendedChat)
		if len(ChatLine) > 15 {
			ChatLine = ChatLine[1:]
		}

	})

	router.GET("/chat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"chat": ChatLine,
		})
	})

	router.Run()
}
