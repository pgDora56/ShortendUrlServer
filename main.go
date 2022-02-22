package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Data struct {
	Key   string
	Value string
}

const (
	DB_FILE = "kvdata.db"
)

func main() {
	// if len(os.Args) > 1 {
	// 	if os.Args[1] == "--get" {
	// 		return
	// 	}
	// }

	log.Println("Shorten URL service started")
	initDb()

	// routerの初期設定
	router := gin.Default()

	// js,css,faviconなどを読み込むためのasstes設定
	router.LoadHTMLGlob("view/*.tmpl")
	router.Static("/resource", "./resource")
	router.StaticFile("/favicon.ico", "./resource/favicon.ico")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "main.tmpl", gin.H{})
	})

	router.GET("/:key", func(ctx *gin.Context) {
		key := ctx.Param("key")
		db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
		if err != nil {
			log.Println("gorm.Open error:", err)
			ctx.HTML(http.StatusOK, "error.tmpl", gin.H{})
		}

		db.Where("key = ?", key).First(&Data{})
	})

	router.Run(":50417")

	log.Println("Shorten URL service stopped")
}

func initDb() {
	db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
	if err != nil {
		panic("gorm.Open error:" + err.Error())
	}

	db.AutoMigrate(&Data{})
}

type Config struct {
}

func getConfig() (cfg Config) {
	js, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Can't read config.json: %v\n", err)
	}

	err = json.Unmarshal(js, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal error config.json: %v\n", err)
	}

	return cfg
}
