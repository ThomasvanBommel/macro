package main

import (
	"embed"
	"io/fs"
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

	api := r.Group("/api")
	api.POST("/register", a.handleRegisterUser)
	api.POST("/login", a.handleLoginUser)
	api.POST("/logout", a.handleLogoutUser)
	api.POST("/validate-session", a.handleSessionValidation)
	api.POST("/food", a.handleCreateFood)
	api.POST("/foods", a.handleListFoods)
	api.POST("/entry", a.handleCreateEntry)
	api.POST("/entries", a.handleListUserEntries)
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
// the provided credentials, and returns an appropriate response to the client. If any step fails,
// it returns an error response. On success, it returns a JSON response indicating that the user
// was registered successfully.
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

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
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

// handleSessionValidation is a handler function for validating the user's session. It retrieves the
// session token from the user's session cookie, checks if the token is valid by querying the
// database, and returns the session information if the token is valid. If the token is not found or
// invalid, it returns a 401 Unauthorized response. This endpoint can be used by the frontend to
// check if the user has an active session and to retrieve session details if needed.
func (a *API) handleSessionValidation(c *gin.Context) {
	defer Trace("handleSessionValidation(c *gin.Context)")()

	t, exists := getSessionToken(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	s, err := a.db.getSessionByToken(t)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// scaleFoodForResponse is a helper function that scales the nutritional values of a Food struct for
// the API response. It divides the calories, carbs, protein, fat, and serving count by 100 to
// convert them back to their original values before they were multiplied for storage in the
// database. This is necessary because the database stores these values as integers multiplied by
// 100 to avoid floating-point precision issues, but the API response should return the actual
// nutritional values to the client.
func scaleFoodForResponse(f *Food) {
	f.Calories /= 100
	f.Carbs /= 100
	f.Protein /= 100
	f.Fat /= 100
	f.ServingCount /= 100
}

// scaleEntryWithFoodForResponse is a helper function that scales the nutritional values of an
// EntryWithFood struct for the API response. It calls the scaleFoodForResponse function to scale
// the food values and also scales the servings value. It divides the servings by 100 to convert it
// back to its original value before it was multiplied for storage in the database. This ensures
// that the API response returns the actual servings value.
func scaleEntryWithFoodForResponse(e *EntryWithFood) {
	scaleFoodForResponse(&e.Food)
	e.Servings /= 100
}

// handleListUserEntries is a handler function for the endpoint that retrieves a list of user
// entries for a specific user on a specific date. It binds the incoming JSON request to a
// ListUserEntriesInput struct, retrieves the user's entries with food details from the database
// based on the provided name and date, and returns the result as a JSON response. If any step
// fails, it returns an appropriate error response to the client.
func (a *API) handleListUserEntries(c *gin.Context) {
	defer Trace("handleListUserEntries(c *gin.Context)")()

	var in ListUserEntriesInput
	if !bindInput(c, &in) {
		return
	}

	res, err := a.db.listUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	for i := range res {
		scaleEntryWithFoodForResponse(&res[i])
	}

	c.JSON(http.StatusOK, res)
}

// handleCreateFood is a handler function for the food creation endpoint. It binds the incoming JSON
// request to a CreateFoodInput struct, retrieves the session token from the user's session, creates
// a new food entry in the database using the provided details and the session token for
// authentication, and returns the created food entry as a JSON response. If any step fails, it
// returns an appropriate error response to the client. On success, it returns the newly created
// food entry in the response.
func (a *API) handleCreateFood(c *gin.Context) {
	defer Trace("handleCreateFood(c *gin.Context)")()

	var in CreateFoodInput
	if !bindInput(c, &in) {
		return
	}

	t, exists := getSessionToken(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	p := CreateFoodParams{
		Name:         in.Name,
		Brand:        in.Brand,
		Calories:     in.Calories,
		Carbs:        in.Carbs,
		Protein:      in.Protein,
		Fat:          in.Fat,
		ServingSize:  in.ServingSize,
		ServingCount: in.ServingCount,
	}

	f, err := a.db.createFoodByToken(p, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	c.JSON(http.StatusOK, f)
}

// handleListFoods is a handler function for the endpoint that retrieves a list of all food entries.
// It queries the database for all food entries and returns them as a JSON response. If any error
// occurs during the query, it returns an appropriate error response to the client.
func (a *API) handleListFoods(c *gin.Context) {
	defer Trace("handleListFoods(c *gin.Context)")()

	foods, err := a.db.listFoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// handleCreateEntry is a handler function for the entry creation endpoint. It binds the incoming
// JSON request to a CreateEntryInput struct, retrieves the session token from the user's session,
// creates a new entry in the database using the provided details and the session token for
// authentication, and returns the created entry as a JSON response. If any step fails, it returns
// an appropriate error response to the client. On success, it returns the newly created entry in
// the response.
func (a *API) handleCreateEntry(c *gin.Context) {
	defer Trace("handleCreateEntry(c *gin.Context)")()

	var in CreateEntryInput
	if !bindInput(c, &in) {
		return
	}

	t, exists := getSessionToken(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	p := CreateEntryParams{
		FoodId:   *in.FoodId,
		MealName: in.MealName,
		Date:     in.Date,
		Servings: in.Servings,
	}

	e, err := a.db.createEntryByToken(p, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	c.JSON(http.StatusOK, e)
}
