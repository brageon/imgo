package main
import ( "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "net/http" )

func main() {
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  // Routes
  e.GET("/", hello)
  e.POST("/users", saveUser)
  e.GET("/users/:id", getUser)
  e.PUT("/users/:id", updateUser)
  e.DELETE("/users/:id", deleteUser)
  e.Logger.Fatal(e.Start(":1323")) }
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!") }
  
func getUser(c echo.Context) error {
  	id := c.Param("id")
	return c.String(http.StatusOK, id) }

func saveUser(c echo.Context) error {
    // Implement logic to save user data from request body (e.g., using a database)
    return c.String(http.StatusCreated, "User saved successfully")
}

func updateUser(c echo.Context) error {
    // Implement logic to update user data based on ID and request body (e.g., using a database)
    return c.String(http.StatusOK, "User updated successfully")
}

func deleteUser(c echo.Context) error {
    // Implement logic to delete user based on ID (e.g., using a database)
    return c.String(http.StatusOK, "User deleted successfully")
}
