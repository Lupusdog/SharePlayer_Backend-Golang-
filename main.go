package main

import (
	"fmt"
	"net/http"

	"time"

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

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://share-player-frontend.vercel.app",
		},

		AllowMethods: []string{
			"POST",
			"GET",
			"OPTION",
		},

		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},

		AllowCredentials: true,
		
		MaxAge: 24 * time.Hour,
	}))

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
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("user",json.Name,7200,"/","",true,true)	

	})

	router.POST("/chat", func(c *gin.Context) {
		var json Chat

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		Name, json_err := c.Cookie("user")
		if json_err != nil {
			fmt.Print(json_err)
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
