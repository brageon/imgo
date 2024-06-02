package main
import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/s3" )

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const savePath = "./public/"
const authKey = "4b22x[:Va4"
const fileSize = 50

const ( letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits )

const ( bucketName = "ddarwin"
        region = "eu-north-1" )

var src = rand.NewSource(time.Now().UnixNano())
var supportedTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/gif":  true,
	"image/png":  true,
	"text/plain": true,
	"video/mp4": true,
	"video/webm": true,
	"application/zip": true, }

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax }
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i-- }
		cache >>= letterIdxBits
		remain-- }
	return string(b) }

func Upload(c echo.Context) error {
	req := c.Request()
	auth := req.Header.Get("Authorization")
	if auth != authKey {
		return c.String(http.StatusForbidden, "Authorization failed") }
	err := req.ParseMultipartForm(fileSize << 20)
	if err != nil {
		return c.String(http.StatusForbidden, "File size is too big") }
	file, header, err := req.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request") }
	defer file.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	file.Seek(0, 0)
	fileType := http.DetectContentType(buff)
	if _, ok := supportedTypes[fileType]; !ok {
		return c.String(http.StatusForbidden, "Unsupported file type") }
	extension := filepath.Ext(header.Filename)
	fileName := RandStringBytesMaskImprSrc(5) + extension
	dst, err := os.Create(savePath + fileName)
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request") }
	return c.String(http.StatusCreated, fmt.Sprintf("https://www.rjtve.com/%s", fileName)) }

func showImage(c echo.Context) error {
	path := c.Param("id")
	return c.File(savePath + path) }

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    file, err := r.FormFile("upload")
    if err != nil {
        http.Error(w, "Error reading uploaded file", http.StatusInternalServerError)
        return }
    defer file.Close()
    fileName := file.Filename
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region), })
    if err != nil {
        http.Error(w, "Error creating AWS session", http.StatusInternalServerError)
        return }
    svc := s3.New(sess)
    _, err = svc.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key: aws.String(fileName),
        Body: file, })
    if err != nil {
        http.Error(w. "Error uploading file to S3", http.StatusInternalServerError)
        return }
    fmt.Fprintf(w, "File %s uploaded successfully to S3!", fileName) }

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.File("/", "index.html")
	e.File("/favicon.ico", "favicon.ico")
	e.GET("/:id", showImage)
	e.POST("/upload", Upload)
	http.HandleFunc("/upload", uploadHandler)
        fmt.Println("Server listening on port 8081")
