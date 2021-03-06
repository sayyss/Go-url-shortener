package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type url struct {
	gorm.Model
	URL string
	Key string
}

type link struct {
	url string
}

func initMigration() {

	db, err = gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Connected to DB brrrrrr")
	}

	db.AutoMigrate(&url{})
}
func generateKey() string {

	key := make([]byte, 6)
	for i := range key {
		key[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(key)
}

func addURL(RecURL string) string {

	var randomKey = generateKey()

	db.Create(&url{URL: RecURL, Key: randomKey})

	var shortened string = "sigh.gq/" + randomKey
	return shortened
}

func createURL(c *gin.Context) {

	//requestURL := body[1]
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])

	requestURL := reqBody[4:]

	shortenedURL := addURL(requestURL)

	c.JSON(201, gin.H{
		"status":   "shortened link",
		"url":      shortenedURL,
		"original": requestURL})
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", nil)
}

func redirect(c *gin.Context) {

	para := c.Param("para")
	println(para)
	var URL url
	db.Where("Key = ?", para).First(&URL)

	c.Redirect(302, URL.URL)
	c.Abort()
}

func main() {

	initMigration()

	r := gin.Default()

	r.LoadHTMLGlob("templates/html/*")
	//r.Static("css/", "templates/css/")
	r.GET("/", home)
	r.GET("/:para", redirect)
	r.POST("/add", createURL)
	/*
		r.GET("/home", home)
		r.POST("/add", createURL)

		r.NoRoute(redirect)
	*/
	r.Run()
}
