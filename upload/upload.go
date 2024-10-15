package main
import ("fmt"
  "io"
  "net/http"
  "os"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware")

func uploadHandler(c echo.Context) error {
  file, err := c.FormFile("upload")
  if err != nil {
    return echo.NewHTTPError(http.StatusUnprocessableEntity, "No file uploaded") }
  uploadDir := "upload"
  if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
    err = os.Mkdir(uploadDir, os.ModePerm)
    if err != nil {
      return echo.NewHTTPError(http.StatusInternalServerError, "Error creating upload directory") }}
  src, err := file.Open()
  if err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, "Error opening uploaded file")}
  defer src.Close()
  fileName := fmt.Sprintf("%s.%s", file.Filename, c.Param("filename"))
  dst, err := os.Create(fmt.Sprintf("%s/%s", uploadDir, fileName))
  if err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, "Error creating destination file") }
  defer dst.Close()
  if _, err := io.Copy(dst, src); err != nil {
    return echo.NewHTTPError(http.StatusInternalServerError, "Error copying uploaded file") }
  return c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully!", fileName)) }

func main() {
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.POST("/upload/:filename", uploadHandler)
  e.Logger.Fatal(e.Start(":8080")) }
