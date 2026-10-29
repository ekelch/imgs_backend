package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

var router = gin.Default()

const sqlDriverName string = "sqlite"
const dataSourceName string = "/Users/evankelch/Programming/brdBase/brdBase.db"

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
	router.PUT("/updateImg", ginUpdate)
	router.POST("/upload", ginUpload)

	router.Run("localhost:8080")
}

func ginListimagesFromHost(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, listImagesCF())
}

func ginUpdateImages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, saveImagesToDatabase(listImagesCF()))
}

func ginListimagesDB(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, listImagesDB())
}

func ginUpdate(c *gin.Context) {
	var updatedImg Image
	err := c.BindJSON(&updatedImg)
	if err != nil {
		log.Fatalf("Failed to bind json: %v/n", err)
	}
	c.JSON(http.StatusOK, updateImage(updatedImg))
}

func ginUpload(c *gin.Context) {
	var uploadDTO UploadDTO
	if err := c.ShouldBind(&uploadDTO); err != nil {
		log.Printf("Failed to bind json: %v\n", err)
	}
	saveNewImage(uploadDTO)
}

func saveImagesToDatabase(images []Image) int {
	db, err := sql.Open(sqlDriverName, dataSourceName)
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
	db, err := sql.Open(sqlDriverName, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer db.Close()

	selectQry := "SELECT filename, uploaded, image_id, category, title, description FROM IMAGES order by id desc"

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

func listImagesCF() []Image {
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

	var imageResponse CFListResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal %v", err)
	}

	return imageResponse.Result.Images
}

func updateImage(updatedImg Image) error {
	db, err := sql.Open(sqlDriverName, dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	updateQry := "UPDATE IMAGES SET (title, description, category) = (?,?,?) WHERE image_id = ?"

	stmt, err := db.Prepare(updateQry)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
		return err
	}
	_, err = stmt.Exec(updatedImg.Title, updatedImg.Description, updatedImg.Category, updatedImg.ImageId)
	if err != nil {
		log.Fatalf("Failed to execute statement: %v", err)
		return err
	}
	return nil
}

func uploadImageToHost(img *multipart.FileHeader) (Image, error) {
	const hostUrl = "https://api.cloudflare.com/client/v4/accounts/400b3f859551ad3d06330ed21b0630dc/images/v1"

	imgFile, err := img.Open()
	if err != nil {
		log.Printf("Failed to open img file: %v", err)
		return Image{}, err
	}
	defer imgFile.Close()

	bodyBuffer := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuffer)
	part, err := writer.CreateFormFile("file", img.Filename)
	if err != nil {
		log.Printf("Failed to create form file part: %v", err)
		return Image{}, err
	}
	if _, err = io.Copy(part, imgFile); err != nil {
		log.Printf("Failed to copy file: %v", err)
		return Image{}, err
	}
	if err = writer.Close(); err != nil {
		log.Printf("Failed to close writer: %v", err)
		return Image{}, err
	}

	postReq, err := http.NewRequest("POST", hostUrl, bodyBuffer)
	if err != nil {
		log.Printf("Failed to create post request: %v", err)
		return Image{}, err
	}
	postReq.Header.Add("Authorization", os.Getenv("imgs_bearer"))
	postReq.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(postReq)
	if err != nil {
		log.Printf("Failed to execute post request to cloudflare: %v", err)
		return Image{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to execute post request to cloudflare: %v", err)
		return Image{}, err
	}

	cfResponse, err := io.ReadAll(resp.Body)

	fmt.Println(string(cfResponse))
	var body CFPostResponse
	err = json.Unmarshal(cfResponse, &body)

	return body.Result, nil
}

func saveNewImage(dto UploadDTO) {
	newImage, err := uploadImageToHost(dto.File)
	if err != nil {
		fmt.Println("Failed to upload image to cloudflare, returning")
	}

	newImage.Title = dto.Title
	newImage.Description = dto.Description
	newImage.Category = dto.Category
	saveImageDB(newImage)
}

func saveImageDB(i Image) {
	db, err := sql.Open(sqlDriverName, dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer db.Close()

	insertQry := "INSERT INTO images (image_id, title, description, filename, uploaded, requireSignedURLs, category) values (?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(insertQry)

	_, err = stmt.Exec(i.ImageId, i.Title, i.Description, i.Filename, i.Uploaded, i.RequireSignedURLs, i.Category)
	if err != nil {
		log.Printf("Failed to execute insert query %v", err)
	}
}
