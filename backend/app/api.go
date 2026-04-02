package main

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"os"
	"time"

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
	applyLoggingMiddleware(r)

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

// applyLoggingMiddleware logs request metadata, latency, and handler errors.
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

// setSessionToken stores the session token in the cookie store.
func setSessionToken(c *gin.Context, token string) {
	defer Trace("setSessionToken(c *gin.Context, token string)", "token", token)()

	session := sessions.Default(c)
	session.Set("token", token)
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

// clearSessionToken clears and persists session state.
func clearSessionToken(c *gin.Context) {
	defer Trace("clearSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// handleCreateSession creates a DB session and persists its token in the cookie.
func (a *API) handleCreateSession(name string, c *gin.Context) error {
	defer Trace("handleCreateSession(name string, c *gin.Context)", "name", name)()

	s, err := a.db.createSession(name)
	if err != nil {
		return err
	}

	setSessionToken(c, s.Token)

	return nil
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	clearSessionToken(c)

	c.JSON(http.StatusOK, gin.H{"message": "User logged out and session cleared successfully"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// scaleFoodForResponse converts fixed-point nutrition values back to response units.
func scaleFoodForResponse(f *Food) {
	f.Calories /= 100
	f.Carbs /= 100
	f.Protein /= 100
	f.Fat /= 100
	f.ServingCount /= 100
}

// scaleFoodForStorage converts nutrition values to fixed-point for storage.
func scaleFoodForStorage(f *Food) {
	f.Calories *= 100
	f.Carbs *= 100
	f.Protein *= 100
	f.Fat *= 100
	f.ServingCount *= 100
}

// scaleEntryWithFoodForResponse scales embedded food fields and serving amount.
func scaleEntryWithFoodForResponse(e *EntryWithFood) {
	scaleFoodForResponse(&e.Food)
	e.Servings /= 100
}

// scaleEntryWithFoodForStorage scales embedded food fields and serving amount for storage.
func scaleEntryWithFoodForStorage(e *EntryWithFood) {
	scaleFoodForStorage(&e.Food)
	e.Servings *= 100
}

// jsonNumberToFloat32 safely converts a JSON number to float32, returning 0 on error.
func jsonNumberToFloat32(n json.Number) float32 {
	f, err := n.Float64()
	if err != nil {
		return 0
	}
	return float32(f)
}

func roundToTwoDecimalPlaces(f float64) float64 {
	return math.Round(f*100) / 100
}

func scaleForClient(n int) float64 {
	return roundToTwoDecimalPlaces(float64(n) / 100)
}

func scaleForDB(n json.Number) int {
	f, err := n.Float64()
	if err != nil {
		return 0
	}

	return int(roundToTwoDecimalPlaces(f) * 100)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	p := CreateEntryParams{
		FoodId:   *in.FoodId,
		MealName: in.MealName,
		Date:     in.Date,
		Servings: int(jsonNumberToFloat32(in.Servings) * 100),
	}

	if p.Servings < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Servings must be greater than 0"})
		return
	}

	e, err := a.db.createEntryByToken(p, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	f := FoodResponse{
		ID:           e.Food.ID,
		Name:         e.Food.Name,
		Brand:        e.Food.Brand,
		Created:      e.Food.Created,
		UserName:     e.Food.UserName,
		Calories:     scaleForClient(e.Food.Calories),
		Carbs:        scaleForClient(e.Food.Carbs),
		Protein:      scaleForClient(e.Food.Protein),
		Fat:          scaleForClient(e.Food.Fat),
		ServingSize:  e.Food.ServingSize,
		ServingCount: scaleForClient(e.Food.ServingCount),
	}

	r := EntryWithFoodResponse{
		ID:       e.ID,
		UserName: e.UserName,
		Food:     f,
		MealName: e.MealName,
		Date:     e.Date,
		Servings: scaleForClient(e.Servings),
		Created:  e.Created,
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	r := make([]EntryWithFoodResponse, len(res))

	for i := range res {
		f := FoodResponse{
			ID:           res[i].Food.ID,
			Name:         res[i].Food.Name,
			Brand:        res[i].Food.Brand,
			Created:      res[i].Food.Created,
			UserName:     res[i].Food.UserName,
			Calories:     scaleForClient(res[i].Food.Calories),
			Carbs:        scaleForClient(res[i].Food.Carbs),
			Protein:      scaleForClient(res[i].Food.Protein),
			Fat:          scaleForClient(res[i].Food.Fat),
			ServingSize:  res[i].Food.ServingSize,
			ServingCount: scaleForClient(res[i].Food.ServingCount),
		}

		r[i] = EntryWithFoodResponse{
			ID:       res[i].ID,
			UserName: res[i].UserName,
			Food:     f,
			MealName: res[i].MealName,
			Date:     res[i].Date,
			Servings: scaleForClient(res[i].Servings),
			Created:  res[i].Created,
		}
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

	p := CreateFoodParams{
		Name:         in.Name,
		Brand:        in.Brand,
		Calories:     scaleForDB(in.Calories),
		Carbs:        scaleForDB(in.Carbs),
		Protein:      scaleForDB(in.Protein),
		Fat:          scaleForDB(in.Fat),
		ServingSize:  in.ServingSize,
		ServingCount: scaleForDB(in.ServingCount),
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food"})
		return
	}

	res := FoodResponse{
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

	c.JSON(http.StatusOK, res)
}

// handleListFoods returns all foods.
func (a *API) handleListFoods(c *gin.Context) {
	defer Trace("handleListFoods(c *gin.Context)")()

	foods, err := a.db.listFoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve foods"})
		return
	}

	f := make([]FoodResponse, len(foods))

	for i := range foods {
		f[i] = FoodResponse{
			ID:           foods[i].ID,
			Name:         foods[i].Name,
			Brand:        foods[i].Brand,
			Created:      foods[i].Created,
			UserName:     foods[i].UserName,
			Calories:     scaleForClient(foods[i].Calories),
			Carbs:        scaleForClient(foods[i].Carbs),
			Protein:      scaleForClient(foods[i].Protein),
			Fat:          scaleForClient(foods[i].Fat),
			ServingSize:  foods[i].ServingSize,
			ServingCount: scaleForClient(foods[i].ServingCount),
		}
	}

	c.JSON(http.StatusOK, f)
}
