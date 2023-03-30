package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "testuser"
	DBPassword = "testT12345"
	DBName     = "testdb"
)

var (
	db *sql.DB
)

type Page struct {
	Articles    []string `json:"articles"`
	NextPageKey string   `json:"nextPageKey,omitempty"`
}

type ListItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {

	router := gin.Default()
	router.GET("/list/head/:listKey", GetHeadHandler)
	router.GET("/list/page/:listKey/:pageKey", GetPageHandler)
	router.POST("/list/set/", SetHandler)

	log.Fatal(router.Run(":8080"))
}

func GetHeadHandler(c *gin.Context) {
	listKey := c.Param("listKey")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	var pageKey string
	err = db.QueryRow("SELECT page_key FROM pages WHERE list_key = $1 ORDER BY created_at DESC LIMIT 1", listKey).Scan(&pageKey)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"nextPageKey": pageKey})
}

func GetPageHandler(c *gin.Context) {

	listKey := c.Param("listKey")
	pageKey := c.Param("pageKey")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	var articles string
	var nextPageKey sql.NullString
	err = db.QueryRow("SELECT articles, next_page_key FROM pages WHERE list_key=$1 AND page_key=$2", listKey, pageKey).Scan(&articles, &nextPageKey)

	if err != nil {
		//報錯處理
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := gin.H{
		"articles": articles,
	}
	if nextPageKey.Valid {
		response["nextPageKey"] = nextPageKey.String
	}

	// Return Response
	c.JSON(http.StatusOK, response)
}

func SetHandler(c *gin.Context) {

	// 利用Gin Library 的ShouldBindJSON來檢查是否符合格式
	var data ListItem
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Insert or update existing values
	result, err := db.Exec(`
		INSERT INTO list_items (list_key, value) 
		VALUES ($1, $2) 
		ON CONFLICT (list_key)  DO UPDATE 
		SET value = EXCLUDED.value
		`, data.Key, data.Value)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 回傳格式
	rowCount, _ := result.RowsAffected()
	response := gin.H{
		"status":    "success",
		"rowCount":  rowCount,
		"listKey":   data.Key,
		"listValue": data.Value,
	}

	//回傳Response
	c.JSON(http.StatusOK, response)
}
