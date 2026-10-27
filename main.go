package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

var router = gin.Default()

func main() {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://br.evankel.ch"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	router.GET("/listFromHost", ginListimagesFromHost)
	router.GET("/list", ginListimagesDB)
	router.GET("/update", ginUpdateImages)
	router.Run("localhost:8080")
}

func ginListimagesFromHost(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, listImages())
}

func ginUpdateImages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, saveImagesToDatabase(listImages()))
}

func ginListimagesDB(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, listImagesDB())
}

func getPublicUrl(image Image) string {
	return fmt.Sprintf("https://imagedelivery.net/lSznIf8GtkPnenDXu_oGvg/%s/public", image.ImageId)
}

func saveImagesToDatabase(images []Image) int {
	db, err := sql.Open("sqlite", "/Users/evankelch/Programming/brdBase/brdBase.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer db.Close()

	insertSQL := `INSERT OR IGNORE INTO IMAGES (filename, uploaded, requireSignedURLs, image_id) VALUES `
	var vals []any

	for _, image := range images {
		insertSQL += "(?,?,?,?),"
		vals = append(vals, image.Filename, image.Uploaded, image.RequireSignedURLs, image.ImageId)
	}
	insertSQL = strings.TrimSuffix(insertSQL, ",")
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v\n", err)
	}
	res, err := stmt.Exec(vals...)
	if err != nil {
		log.Fatalf("Failed to execute database query: %v\n", err)
	}
	rows, _ := res.RowsAffected()
	return int(rows)
}

func listImagesDB() []Image {
	db, err := sql.Open("sqlite", "/Users/evankelch/Programming/brdBase/brdBase.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer db.Close()

	selectQry := "SELECT filename, uploaded, image_id, category, title, description FROM IMAGES"

	rows, err := db.Query(selectQry)
	if err != nil {
		log.Fatalf("Failed to select rows from database: %v\n", err)
	}

	var images []Image
	for rows.Next() {
		var i Image
		err = rows.Scan(&i.Filename, &i.Uploaded, &i.ImageId, &i.Category, &i.Title, &i.Description)
		if err != nil {
			log.Fatalf("Failed to scan row: %v\n", err)
		}
		images = append(images, i)
	}
	return images
}

func listImages() []Image {
	listImgsUrl := "https://api.cloudflare.com/client/v4/accounts/400b3f859551ad3d06330ed21b0630dc/images/v2"
	req, err := http.NewRequest("GET", listImgsUrl, nil)
	req.Header.Add("Authorization", os.Getenv("imgs_bearer"))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed GET request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed with HTTP error code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var imageResponse ImageResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal %v", err)
	}

	return imageResponse.Result.Images
}

func uploadImage() {

}

func uploadImageToHost() {

}
