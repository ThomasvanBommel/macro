package main

import (
	"embed"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/build/*
var frontendBuildFiles embed.FS

// API is a struct that holds a reference to the database and provides methods for handling API
// routes. It serves as a central point for defining the API endpoints and their associated logic,
// allowing for better organization and separation of concerns in the application. By encapsulating
// the API logic within this struct, it becomes easier to manage and extend the API functionality.
type API struct {
	db *Database
}

// InitAPI initializes the API by setting up the necessary middleware and routes for the Gin engine.
// It takes a Gin engine and a Database instance as arguments, and returns an API instance. The
// function applies session management and logging middleware, registers routes for serving static
// assets, defines API routes, and sets up a catch-all route for handling undefined paths. This
// function is responsible for configuring the API endpoints and ensuring that the necessary
// middleware is in place for handling requests and managing sessions.
func InitAPI(r *gin.Engine, db *Database) *API {
	defer Trace("InitAPI(r *gin.Engine, db *Database)")()

	api := &API{db: db}

	applySessionMiddleware(r)
	applyLoggingMiddleware(r)

	registerStaticAssetRoute(r)
	api.registerAPIRoutes(r)
	registerNoRoute(r)

	return api
}

// applySessionMiddleware sets up session management for the Gin engine using cookie-based sessions.
// It reads the session secret from an environment variable, and configures the session middleware
// to use this secret for encrypting session cookies. This ensures that user sessions are securely
// managed across requests. In the scenario where the session secret is not provided, it generates a
// random secret to ensure that sessions are still secure. This will invalidate all existing
// sessions on restart, but prevents the application from running with an insecure configuration.
func applySessionMiddleware(r *gin.Engine) {
	defer Trace("applySessionMiddleware(r *gin.Engine)")()

	s := []byte(os.Getenv("SESSION_SECRET"))
	if len(s) == 0 {
		s = GenerateRandomBytes(128)
	}

	r.Use(sessions.Sessions("macro_session", cookie.NewStore(s)))
}

// applyLoggingMiddleware sets up a logging middleware for the Gin engine that logs details about
// each request, including the HTTP method, path, query parameters, client IP, user agent, and
// the duration of the request. It also logs any errors that occur during the handling of the
// request.
func applyLoggingMiddleware(r *gin.Engine) {
	defer Trace("applyLoggingMiddleware(r *gin.Engine)")()

	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		dur := time.Since(start)

		attr := []slog.Attr{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.Duration("duration", dur),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) == 0 {
			slog.Info("Handled request", "details", attr)
			return
		}

		errs := c.Errors.Errors()
		slog.Error("Request resulted in errors", "count", len(errs), "errors", errs,
			"details", attr)
	})
}

// registerNoRoute sets up a catch-all route for the Gin engine that handles requests to undefined
// routes. It checks if the request path starts with "/api" and returns a 404 JSON response if it
// does, indicating that the API endpoint was not found. For all other paths, it serves the
// "index.html" file from the embedded frontend build files, allowing the frontend application to
// handle client-side routing for non-API requests. This ensures that API requests receive
// appropriate error responses, while still allowing the frontend to function correctly.
func registerNoRoute(r *gin.Engine) {
	defer Trace("registerNoRoute(r *gin.Engine)")()

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		f, err := frontendBuildFiles.ReadFile("frontend/build/index.html")
		if err != nil {
			slog.Error("Failed to load index.html", "error", err)
			c.String(http.StatusInternalServerError, "Failed to load index.html")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", f)
	})
}

// registerStaticAssetRoute sets up a route for serving static assets from the embedded filesystem.
// It creates a sub-filesystem from the embedded frontend build files, and configures the Gin engine
// to serve files from this sub-filesystem at the "/assets" URL path. This allows the application to
// serve static assets like JavaScript, CSS, and image files. Panics if there is an error creating
// the sub-filesystem, ensuring that the application does not run with an invalid configuration.
func registerStaticAssetRoute(r *gin.Engine) {
	defer Trace("registerStaticAssetRoute(r *gin.Engine)")()

	f, err := fs.Sub(frontendBuildFiles, "frontend/build/assets")
	FatalOnError(err, "Failed to create sub filesystem for frontend assets")

	r.StaticFS("/assets", http.FS(f))
}

// registerAPIRoutes defines the API routes for the application. It groups all API endpoints under
// the "/api" URL path.
func (a *API) registerAPIRoutes(r *gin.Engine) {
	defer Trace("registerAPIRoutes(r *gin.Engine)")()

	// naming scheme: handle<action><resource>

	api := r.Group("/api")
	api.POST("/register", a.handleRegisterUser)
	api.POST("/login", a.handleLoginUser)
	api.POST("/logout", a.handleLogoutUser)
	api.POST("/profile", a.handleGetUserProfile)
}

// bindInput is a helper function that binds the JSON body of a request to a specified struct. It
// takes a Gin context and a pointer to a struct as arguments, and attempts to bind the JSON body
// of the request to the struct using Gin's ShouldBindJSON method. If the binding fails, it returns
// a 400 Bad Request response with the error message, and returns false. If the binding is
// successful, it returns true, allowing the caller to proceed with the now-populated struct.
func bindInput(c *gin.Context, obj any) bool {
	defer Trace("bindInput(c *gin.Context, obj any)")()

	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

// setSessionToken is a helper function that sets a session token in the user's session cookie. It
// takes a Gin context and a token string as arguments, retrieves the current session using the
// sessions.Default function, sets the "token" key in the session to the provided token value, and
// saves the session. This is typically used after a user successfully logs in or registers to store
// their session token for subsequent authenticated requests.
func setSessionToken(c *gin.Context, token string) {
	defer Trace("setSessionToken(c *gin.Context, token string)", "token", token)()

	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()
}

// getSessionToken is a helper function that retrieves the session token from the user's session
// cookie. It returns the token string and a boolean indicating whether the token was found. It
// accesses the current session using the sessions.Default function, attempts to get the "token" key
// from the session, and checks if it is nil. If the token is not found (i.e., nil), it returns an
// empty string and false. If the token is found, it asserts that the token is a string and returns
// it along with true. This function is typically used to check if a user has an active session and
// to retrieve their session token for authentication purposes.
func getSessionToken(c *gin.Context) (string, bool) {
	defer Trace("getSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		return "", false
	}

	return token.(string), true
}

// clearSessionToken is a helper function that clears the session token from the user's session
// cookie. It retrieves the current session using the sessions.Default function, calls the Clear
// method to remove all session data, and saves the session. This is typically used when a user logs
// out to ensure that their session token is removed and the session is invalidated.
func clearSessionToken(c *gin.Context) {
	defer Trace("clearSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// handleCreateSession is a helper function that creates a new session for a given username and sets
// the session token in the user's session cookie. It takes a username and a Gin context as
// arguments, creates a new session in the database for the specified username, and sets the
// resulting session token in the user's session cookie using the setSessionToken function. If there
// is an error creating the session, it returns the error to the caller for handling.
func (a *API) handleCreateSession(name string, c *gin.Context) error {
	defer Trace("handleCreateSession(name string, c *gin.Context)", "name", name)()

	s, err := a.db.createSession(name)
	if err != nil {
		return err
	}

	setSessionToken(c, s.Token)

	return nil
}

// handleRegisterUser is a handler function for the user registration endpoint. It binds the
// incoming JSON request to a UserCredentialInput struct, creates a new user in the database using
// the provided credentials, and sets a session token for the newly registered user. If any step
// fails, it returns an appropriate error response to the client. On success, it returns a JSON
// response indicating that the user was registered successfully and logged in.
func (a *API) handleRegisterUser(c *gin.Context) {
	defer Trace("handleRegisterUser(c *gin.Context)")()

	var in UserCredentialInput
	if !bindInput(c, &in) {
		return
	}

	err := a.db.createUser(in.Name, in.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = a.handleCreateSession(in.Name, c)
	if err != nil {
		msg := "User created but failed to create session. Please try logging in."
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully and logged in"})
}

// handleLoginUser is a handler function for the user login endpoint. It binds the incoming JSON
// request to a UserCredentialInput struct, verifies the user's credentials against the database,
// and sets a session token if the credentials are valid. If the credentials are invalid, it returns
// a 401 Unauthorized response. If there is an error during the process, it returns an appropriate
// error response to the client. On success, it returns a JSON response indicating that the user was
// logged in successfully.
func (a *API) handleLoginUser(c *gin.Context) {
	defer Trace("handleLoginUser(c *gin.Context)")()

	var in UserCredentialInput
	if !bindInput(c, &in) {
		return
	}

	u, err := a.db.getUserByNameAndPassword(in.Name, in.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = a.handleCreateSession(u.Name, c)
	if err != nil {
		msg := "Failed to create session. Please try again."
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

// handleLogoutUser is a handler function for the user logout endpoint. It retrieves the session
// token from the user's session, deletes the corresponding session from the database, and clears
// the session token from the user's session cookie. If any step fails, it returns an appropriate
// error response to the client. On success, it returns a JSON response indicating that the user
// logged out successfully.
func (a *API) handleLogoutUser(c *gin.Context) {
	defer Trace("handleLogoutUser(c *gin.Context)")()

	t, exists := getSessionToken(c)
	if !exists {
		c.JSON(http.StatusOK, gin.H{"message": "No active session"})
		return
	}

	err := a.db.deleteSession(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	clearSessionToken(c)

	c.JSON(http.StatusOK, gin.H{"message": "User logged out and session cleared successfully"})
}

// handleGetUserProfile is a handler function for the user profile endpoint. It binds the incoming
// JSON request to a UserProfileRequestInput struct, retrieves the user's entries with food details
// from the database based on the provided name and date, and returns the result as a JSON response.
// If any step fails, it returns an appropriate error response to the client.
func (a *API) handleGetUserProfile(c *gin.Context) {
	defer Trace("handleGetUserProfile(c *gin.Context)")()

	var in UserProfileRequestInput
	if !bindInput(c, &in) {
		return
	}

	res, err := a.db.listUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// -------------------------------------------- OLD:

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
	api.GET("/foods", getFoods)
	api.POST("/food", addFood)
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

func getFoods(c *gin.Context) {
	foods, err := getDB(c).getFoods()
	if err != nil {
		log.Printf("Error retrieving foods: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

func addFood(c *gin.Context) {
	var req RequestCreateFoodModel
	if !bindModel(c, &req) {
		return
	}

	// Get session token from cookie
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res, err := getDB(c).createFood(&req, token.(string))
	if err != nil {
		log.Printf("Error creating food: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func addEntry(c *gin.Context) {
	var req RequestAddEntryModel
	if !bindModel(c, &req) {
		return
	}

	// Get session token from cookie
	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	res, err := getDB(c).createEntry(&req, token.(string))
	if err != nil {
		log.Printf("Error creating entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	c.JSON(http.StatusOK, res)
}
