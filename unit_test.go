package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var router = gin.Default()

func TestGetHeadHandler(t *testing.T) {

	// 連接上Database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName))

	if err != nil {
		t.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	// 初始化資料庫
	err = initializeSchema(db)

	if err != nil {
		log.Fatalf("Error %q", err)
	}

	// 開啓HTTP Server用於測試
	router.GET("/list/head/:listKey", GetHeadHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// 寫進測試資料
	_, err = db.Exec("INSERT INTO pages (list_key, page_key, articles, created_at) VALUES ($1, $2, $3, NOW())", "testlist", "testpage", "testarticle")
	assert.NoError(t, err)

	// 發送HTTP Request 用於測試
	resp, err := http.Get(ts.URL + "/list/head/testlist")

	if err != nil {
		t.Errorf("Http Get Method error : %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d; return %d", http.StatusOK, resp.StatusCode)
	}

	// 檢查Response 是否一致
	expected := gin.H{"nextPageKey": "testpage"}
	var actual gin.H
	err = json.NewDecoder(resp.Body).Decode(&actual)

	if err != nil {
		t.Errorf("Json newDecoder error %s", err)
	}

	assert.Equal(t, expected, actual)

	// 刪除測試資料

	_, err = db.Exec("DELETE FROM pages WHERE list_key = $1 AND page_key = $2", "testlist", "testpage")
	if err != nil {
		t.Errorf("Query error : %s", err)
	}
}

func TestGetPageHandler(t *testing.T) {

	// 連接上Database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		t.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	// 初始化資料庫
	err = initializeSchema(db)

	if err != nil {
		log.Fatalf("Error %q", err)
	}

	// 開啓HTTP Server用於測試
	router.GET("/list/page/:listKey/:pageKey", GetPageHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// 寫進測試資料
	_, err = db.Exec("INSERT INTO pages (list_key, page_key, articles, next_page_key, created_at) VALUES ($1, $2, $3, $4, NOW())", "testlist", "testpage", "testarticle", "nextpage")
	if err != nil {
		log.Fatalf("Error query : %q", err)
	}

	// 發送HTTP Request 用於測試
	resp, err := http.Get(ts.URL + "/list/page/testlist/testpage")

	if err != nil {
		t.Errorf("Http Get Method error : %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d; return %d", http.StatusOK, resp.StatusCode)
	}

	// 檢查Response 是否一致
	expected := gin.H{"articles": "testarticle", "nextPageKey": "nextpage"}
	var actual gin.H
	err = json.NewDecoder(resp.Body).Decode(&actual)

	if err != nil {
		t.Errorf("Json newDecoder error %s", err)
	}
	assert.Equal(t, expected, actual)

	// 刪除測試資料

	_, err = db.Exec("DELETE FROM pages WHERE list_key = $1 AND page_key = $2", "testlist", "testpage")
	if err != nil {
		t.Errorf("Query error : %s", err)
	}

}

func TestSetHandler(t *testing.T) {

	//連接到Database
	testDB, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer testDB.Close()

	db = testDB

	// 初始化資料庫
	err = initializeSchema(testDB)

	if err != nil {
		log.Fatalf("Error %q", err)
	}

	//開啓HTTP Server用於測試
	router.POST("/list/set/", SetHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	body := ListItem{
		Key:   "key1",
		Value: "value1",
	}
	bodyJson, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/list/set/", bytes.NewBuffer(bodyJson))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; return %d", http.StatusOK, recorder.Code)
	}

	// 測試有沒有儲存成功
	var value string
	err = db.QueryRow("SELECT value FROM list_items WHERE list_key=$1", "key1").Scan(&value)
	if err != nil {
		t.Fatalf("Error querying test database: %q", err)
	}
	if value != "value1" {
		t.Errorf("Expected value %s; got %s", "value1", value)
	}
}

// 創建TABLE若不存在
func initializeSchema(db *sql.DB) error {

	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS pages (
            list_key text NOT NULL,
            page_key text NOT NULL,
            articles text NOT NULL,
            next_page_key text,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 HOURS',
            PRIMARY KEY (list_key, page_key)
        );
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS list_items (
            list_key text NOT NULL,
            value text NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 HOURS',
            PRIMARY KEY (list_key)
        );
    `)
	if err != nil {
		return err
	}

	return nil
}
