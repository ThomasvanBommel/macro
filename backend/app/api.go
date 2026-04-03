package main

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// API handles HTTP routes and database-backed request operations.
type API struct {
	db *Database
}

// InitAPI configures middleware and registers API routes.
func InitAPI(r *gin.Engine, db *Database) *API {
	defer Trace("InitAPI(r *gin.Engine, db *Database)")()

	api := &API{db: db}

	applySessionMiddleware(r)

	api.registerAPIRoutes(r)

	return api
}

// applySessionMiddleware configures cookie-backed sessions.
// If SESSION_SECRET is missing, a random in-memory secret is used.
func applySessionMiddleware(r *gin.Engine) {
	defer Trace("applySessionMiddleware(r *gin.Engine)")()

	s := []byte(os.Getenv("SESSION_SECRET"))
	if len(s) == 0 {
		s = GenerateRandomBytes(128)
	}

	r.Use(sessions.Sessions("macro_session", cookie.NewStore(s)))
}

// registerAPIRoutes mounts all endpoints under /api.
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

// bindInput binds and validates a JSON body; on failure it writes 400 and returns false.
func bindInput(c *gin.Context, obj any) bool {
	defer Trace("bindInput(c *gin.Context, obj any)")()

	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(err)
		c.Set("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

// initSession initializes session state with the provided Session struct.
func initSession(s *Session, c *gin.Context) {
	defer Trace("initSession(s *Session, c *gin.Context)", "session", s)()

	session := sessions.Default(c)
	session.Set("token", s.Token)
	session.Set("username", s.UserName)
	session.Save()
}

// clearSession clears and persists session state.
func clearSession(c *gin.Context) {
	defer Trace("clearSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	session.Set("token", nil)
	session.Set("username", nil)
	session.Save()
}

// getSessionToken returns the session token if present.
func getSessionToken(c *gin.Context) (string, bool) {
	defer Trace("getSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		return "", false
	}

	return token.(string), true
}

// createSession creates a DB session and persists its token in the cookie.
func (a *API) createSession(name string, c *gin.Context) error {
	defer Trace("createSession(name string, c *gin.Context)", "name", name)()

	s, err := a.db.createSession(name)
	if err != nil {
		c.Set("error", err.Error())
		return err
	}

	initSession(s, c)

	return nil
}

// handleSessionValidation verifies the current session token.
func (a *API) handleSessionValidation(c *gin.Context) {
	defer Trace("handleSessionValidation(c *gin.Context)")()

	t, exists := getSessionToken(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	s, err := a.db.getSessionByToken(t)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	initSession(s, c)

	c.JSON(http.StatusOK, s)
}

// handleRegisterUser registers a new user.
func (a *API) handleRegisterUser(c *gin.Context) {
	defer Trace("handleRegisterUser(c *gin.Context)")()

	var in UserCredentialInput
	if !bindInput(c, &in) {
		return
	}

	err := a.db.createUser(in.Name, in.Password)
	if err != nil {
		c.Set("error", err.Error())

		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// handleLoginUser authenticates a user and creates a session.
func (a *API) handleLoginUser(c *gin.Context) {
	defer Trace("handleLoginUser(c *gin.Context)")()

	var in UserCredentialInput
	if !bindInput(c, &in) {
		return
	}

	u, err := a.db.getUserByNameAndPassword(in.Name, in.Password)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = a.createSession(u.Name, c)
	if err != nil {
		c.Set("error", err.Error())
		msg := "Failed to create session. Please try again."
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

// handleLogoutUser invalidates the active session and clears local session state.
func (a *API) handleLogoutUser(c *gin.Context) {
	defer Trace("handleLogoutUser(c *gin.Context)")()

	t, exists := getSessionToken(c)
	if !exists {
		c.JSON(http.StatusOK, gin.H{"message": "No active session"})
		return
	}

	err := a.db.deleteSession(t)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	clearSession(c)

	c.JSON(http.StatusOK, gin.H{"message": "User logged out and session cleared successfully"})
}

// roundToTwoDecimalPlaces rounds a float to two decimal places.
func roundToTwoDecimalPlaces(f float64) float64 {
	return math.Round(f*100) / 100
}

// scaleForClient converts an integer DB value to a float for client responses.
func scaleForClient(n int) float64 {
	return roundToTwoDecimalPlaces(float64(n) / 100)
}

// scaleForStorage converts a JSON number to an integer for DB storage.
func scaleForStorage(n json.Number) int {
	f, err := n.Float64()
	if err != nil {
		return 0
	}

	return int(roundToTwoDecimalPlaces(f) * 100)
}

// toCreateEntryParams converts a CreateEntryInput to CreateEntryParams, applying necessary scaling.
func toCreateEntryParams(in *CreateEntryInput) CreateEntryParams {
	return CreateEntryParams{
		FoodId:   *in.FoodId,
		MealName: in.MealName,
		Date:     in.Date,
		Servings: scaleForStorage(in.Servings),
	}
}

// toFoodResponse converts a Food DB record to a FoodResponse for API responses.
func toFoodResponse(f *Food) FoodResponse {
	return FoodResponse{
		ID:           f.ID,
		Name:         f.Name,
		Brand:        f.Brand,
		Created:      f.Created,
		UserName:     f.UserName,
		Calories:     scaleForClient(f.Calories),
		Carbs:        scaleForClient(f.Carbs),
		Protein:      scaleForClient(f.Protein),
		Fat:          scaleForClient(f.Fat),
		ServingSize:  f.ServingSize,
		ServingCount: scaleForClient(f.ServingCount),
	}
}

// toEntryWithFoodResponse converts an EntryWithFood DB record to an EntryWithFoodResponse
func toEntryWithFoodResponse(e *EntryWithFood) EntryWithFoodResponse {
	return EntryWithFoodResponse{
		ID:       e.ID,
		UserName: e.UserName,
		Food:     toFoodResponse(&e.Food),
		MealName: e.MealName,
		Date:     e.Date,
		Servings: scaleForClient(e.Servings),
		Created:  e.Created,
	}
}

// toCreateFoodParams converts a CreateFoodInput to CreateFoodParams, applying necessary scaling.
func toCreateFoodParams(in *CreateFoodInput) CreateFoodParams {
	return CreateFoodParams{
		Name:         in.Name,
		Brand:        in.Brand,
		Calories:     scaleForStorage(in.Calories),
		Carbs:        scaleForStorage(in.Carbs),
		Protein:      scaleForStorage(in.Protein),
		Fat:          scaleForStorage(in.Fat),
		ServingSize:  in.ServingSize,
		ServingCount: scaleForStorage(in.ServingCount),
	}
}

// handleCreateEntry creates a food entry for the authenticated user.
func (a *API) handleCreateEntry(c *gin.Context) {
	defer Trace("handleCreateEntry(c *gin.Context)")()

	var in CreateEntryInput
	if !bindInput(c, &in) {
		return
	}

	t, exists := getSessionToken(c)
	if !exists {
		c.Set("error", "No active session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	p := toCreateEntryParams(&in)
	if p.Servings < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Servings must be greater than 0"})
		return
	}

	e, err := a.db.createEntryByToken(p, t)
	if err != nil {
		c.Set("error", err.Error())

		if strings.Contains(err.Error(), "NOT NULL constraint failed: entry.user_name") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if strings.Contains(err.Error(), "no rows in result set") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food ID"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	r := toEntryWithFoodResponse(e)
	c.JSON(http.StatusOK, r)
}

// handleListUserEntries returns entries with food details for a user and date.
func (a *API) handleListUserEntries(c *gin.Context) {
	defer Trace("handleListUserEntries(c *gin.Context)")()

	var in ListUserEntriesInput
	if !bindInput(c, &in) {
		return
	}

	res, err := a.db.listUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	r := make([]EntryWithFoodResponse, len(res))
	for i := range res {
		r[i] = toEntryWithFoodResponse(&res[i])
	}

	c.JSON(http.StatusOK, r)
}

// handleCreateFood creates a food record for the authenticated user.
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

	p := toCreateFoodParams(&in)
	if p.ServingCount < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Serving count must be greater than 0"})
		return
	}

	if p.Calories < 0 || p.Carbs < 0 || p.Protein < 0 || p.Fat < 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Calories, carbs, protein, and fat must be non-negative"})
		return
	}

	f, err := a.db.createFoodByToken(p, t)
	if err != nil {
		c.Set("error", err.Error())

		if strings.Contains(err.Error(), "NOT NULL constraint failed: food.user_name") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	c.JSON(http.StatusOK, toFoodResponse(f))
}

// handleListFoods returns all foods.
func (a *API) handleListFoods(c *gin.Context) {
	defer Trace("handleListFoods(c *gin.Context)")()

	foods, err := a.db.listFoods()
	if err != nil {
		c.Set("error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	f := make([]FoodResponse, len(foods))

	for i := range foods {
		f[i] = toFoodResponse(&foods[i])
	}

	c.JSON(http.StatusOK, f)
}
