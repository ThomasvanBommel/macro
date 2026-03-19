package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) registerUser(c *gin.Context) {
	var req RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := createUser(h.db, req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) loginUser(c *gin.Context) {
	var req LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	token, err := loginUser(h.db, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *Handler) getSession(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("Session token: %s", token.(string))
	s, err := selectSession(h.db, token.(string))
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		log.Printf("Error retrieving session token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}
	if s == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or invalid"})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *Handler) delSession(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No active session"})
		return
	}

	err := deleteSession(h.db, token.(string))
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove session"})
	}

	c.JSON(http.StatusOK, session)
}

func (h *Handler) addFood(c *gin.Context) {
	var food Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := insertFood(h.db, food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add food"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Food added successfully"})
}