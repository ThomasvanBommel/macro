package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func initAPI(r *gin.Engine) {
	// Load session secret & middleware
	secret, err := os.ReadFile("/run/secrets/session_secret")
	if err != nil {
		log.Printf("Failed to read session secret: %v", err)
		secret = []byte(os.Getenv("SESSION_SECRET"))
	}

	if len(secret) == 0 {
		log.Fatal("Docker secret and SESSION_SECRET environment variable are both empty.")
	}

	r.Use(sessions.Sessions("macro_session", cookie.NewStore(secret)))

	// Define API endpoints
	api := r.Group("/api")
	api.POST("/register", register)
	api.POST("/login", login)
	api.POST("/logout", logout)
	api.POST("/entries", getEntries)
	api.POST("/entry", addEntry)
}

/* Bind request context to a model struct */
func bindModel(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

/* Extract the database wrapper from the request context, panic if it doesn't exist. */
func getDB(c *gin.Context) *DatabaseWrapper {
	return c.MustGet("db").(*DatabaseWrapper)
}

func register(c *gin.Context) {
	// Parse form data
	var req RequestUserModel
	if !bindModel(c, &req) {
		return
	}

	// Insert user into database and create session, fetch session info
	token, res, err := getDB(c).register(req)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set session cookie
	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()

	c.JSON(http.StatusOK, res)
}

func login(c *gin.Context) {
	// Check for existing session cookie, return session info if valid
	session := sessions.Default(c)
	if session.Get("token") != nil {
		res, err := getDB(c).getSessionInfo(session.Get("token").(string))
		if err == nil {
			c.JSON(http.StatusOK, res)
			return
		}
	}

	// Parse form data
	var req RequestUserModel
	if !bindModel(c, &req) {
		return
	}

	// Insert session into database and fetch session info
	token, res, err := getDB(c).login(req)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set session cookie
	session.Set("token", token)
	session.Save()

	c.JSON(http.StatusOK, res)
}

func logout(c *gin.Context) {
	// Delete session from database
	session := sessions.Default(c)
	if token := session.Get("token"); token != nil {
		err := getDB(c).deleteSession(token.(string))
		if err != nil {
			log.Printf("Error deleting session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
			return
		}
	}

	// Clear session cookie, expire immediately
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func getEntries(c *gin.Context) {
	var req RequestEntriesModel
	if !bindModel(c, &req) {
		return
	}

	entries, err := getDB(c).getEntries(req)
	if err != nil {
		log.Printf("Error retrieving entries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func addEntry(c *gin.Context) {
	// Get session token from cookie
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind request body to model
	var req RequestAddEntryModel
	if !bindModel(c, &req) {
		return
	}

	// Add food if FoodId is not provided
	if req.FoodId == nil {
		if req.Food == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Either food_id or food details must be provided"})
			return
		}

		food_id, err := getDB(c).addFood(req.Food, token.(string))
		if err != nil {
			log.Printf("Error adding food: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add food"})
			return
		}

		req.FoodId = &food_id
	}

	// Add entry and return it with food info populated
	entry, err := getDB(c).addEntry(req, token.(string))
	if err != nil {
		log.Printf("Error adding entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add entry"})
		return
	}

	c.JSON(http.StatusOK, entry)
}
