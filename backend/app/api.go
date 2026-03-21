package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
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
	api.GET("/session", getSession)
	// api.GET("/profile", profile)
	// api.POST("/session", getSession)
	// api.DELETE("/session", delSession)
	// api.PUT("/food", addFood)
	// api.GET("/food", getFoods)
	// api.PUT("/entry", addEntry)
	// api.GET("/entry", getEntries)
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

	// Insert user into database
	err := getDB(c).insertUser(req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func login(c *gin.Context) {
	// Parse form data
	var req RequestUserModel
	if !bindModel(c, &req) {
		return
	}

	// Insert session into database
	token, err := getDB(c).insertSession(req)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set session cookie
	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func getSession(c *gin.Context) {
	// Get session token from cookie
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Select session from database
	s, err := getDB(c).selectSession(token.(string))
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// func (h *Handler) registerUser(c *gin.Context) {
// 	var req RegisterUserRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}

// 	err := createUser(h.db, req.Username, req.Password)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
// 			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
// }

// func (h *Handler) loginUser(c *gin.Context) {
// 	var req LoginUserRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}

// 	token, err := loginUser(h.db, req.Username, req.Password)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
// 		return
// 	}
// 	if token == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}

// 	session := sessions.Default(c)
// 	session.Set("token", token)
// 	session.Save()

// 	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
// }

// func (h *Handler) getSession(c *gin.Context) {
// 	session := sessions.Default(c)
// 	token := session.Get("token")
// 	if token == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	log.Printf("Session token: %s", token.(string))
// 	s, err := selectSession(h.db, token.(string))
// 	if err != nil {
// 		log.Printf("Error retrieving session token: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
// 		return
// 	}
// 	if s == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or invalid"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, s)
// }

// func (h *Handler) delSession(c *gin.Context) {
// 	session := sessions.Default(c)
// 	token := session.Get("token")
// 	if token == nil {
// 		c.JSON(http.StatusOK, gin.H{"message": "No active session"})
// 		return
// 	}

// 	err := deleteSession(h.db, token.(string))
// 	if err != nil {
// 		log.Printf("Error deleting session: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove session"})
// 	}

// 	c.JSON(http.StatusOK, session)
// }

// func (h *Handler) addFood(c *gin.Context) {
// 	var req PutFoodRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Printf("Error parsing JSON: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}

// 	session := sessions.Default(c)
// 	token := session.Get("token")

// 	if token == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	req.SessionToken = token.(string)
// 	err := insertFood(h.db, req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add food"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Food added successfully"})
// }

// func (h *Handler) getFoods(c *gin.Context) {
// 	foods, err := selectFoods(h.db)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, foods)
// }

// func (h *Handler) addEntry(c *gin.Context) {
// 	var req PutEntryRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Printf("Error parsing JSON: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}

// 	session := sessions.Default(c)
// 	token := session.Get("token")

// 	if token == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	req.SessionToken = token.(string)
// 	err := insertEntry(h.db, req)
// 	if err != nil {
// 		log.Printf("Error inserting entry: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add entry"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Entry added successfully"})
// }

// func (h *Handler) getEntries(c *gin.Context) {
// 	var req EntryRequest
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		log.Printf("Error parsing query parameters: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
// 		return
// 	}

// 	entries, err := selectEntries(h.db, req.UserID, req.Date)
// 	if err != nil {
// 		log.Printf("Error retrieving entries: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve entries"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, entries)
// }
