package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type APIHandler func(*gin.Context) APIResponse

// Wrap converts an APIHandler to a gin.HandlerFunc, handling APIResponse and errors uniformly.
func Wrap(fn APIHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := fn(c)
		if r == nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		quiet := false
		if res, ok := r.(*ErrorResponse); ok && res != nil {
			quiet = res.Quiet
			c.Error(res)
		}

		payload := r.Result()
		if quiet || payload == nil {
			c.Status(r.Status())
			return
		}

		c.JSON(r.Status(), payload)
	}
}

// InitAPI configures middleware and registers API routes.
func InitAPI(r *gin.Engine, db *Database) *API {
	defer Trace("InitAPI(r *gin.Engine, db *Database)")()

	api := &API{
		db: db,

		OK:             func(d any) APIResponse { return &DataResponse{200, d} },
		Message:        func(m string) APIResponse { return &DataResponse{200, gin.H{"message": m}} },
		Created:        func(d any) APIResponse { return &DataResponse{201, d} },
		Accepted:       func() APIResponse { return &DataResponse{202, nil} },
		NoContent:      func() APIResponse { return &DataResponse{204, nil} },
		PartialContent: func(d any) APIResponse { return &DataResponse{206, d} },

		BadRequest:      func(e error) APIResponse { return &ErrorResponse{400, false, e} },
		Unauthorized:    func(e error) APIResponse { return &ErrorResponse{401, true, e} },
		Forbidden:       func(e error) APIResponse { return &ErrorResponse{403, true, e} },
		NotFound:        func(e error) APIResponse { return &ErrorResponse{404, true, e} },
		Conflict:        func(e error) APIResponse { return &ErrorResponse{409, true, e} },
		TooManyRequests: func(e error) APIResponse { return &ErrorResponse{429, true, e} },

		InternalServerError: func(e error) APIResponse { return &ErrorResponse{500, true, e} },
		NotImplemented:      func(e error) APIResponse { return &ErrorResponse{501, true, e} },
	}

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
	api.POST("/register", Wrap(a.handleRegisterUser))
	api.POST("/login", Wrap(a.handleLoginUser))
	api.POST("/logout", Wrap(a.handleLogoutUser))
	api.POST("/validate-session", Wrap(a.handleSessionValidation))
	api.POST("/food", Wrap(a.handleCreateFood))
	api.POST("/foods", Wrap(a.handleListFoods))
	api.POST("/entry", Wrap(a.handleCreateEntry))
	api.POST("/entries", Wrap(a.handleListUserEntries))
	api.POST("/diary", Wrap(a.handleGetDiary))
	api.POST("/food/search", Wrap(a.handleFoodSearch))

	api.POST("/entry/edit", Wrap(a.handleEditEntry))
	api.POST("/entry/delete", Wrap(a.handleDeleteEntry))

	api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, nil) })
}

// DBErrorMessages holds overridable response messages for known DB error conditions.
// Zero-value fields fall back to their default.
type DBErrorMessages struct{ Unique, NotFound, Busy, Default string }

// DBError maps a database error to an APIResponse, with optional message overrides.
func (a *API) DBError(err error, msgs ...DBErrorMessages) APIResponse {
	m := DBErrorMessages{
		"Already exists.",
		"Not found.",
		"Database is busy, please try again.",
		"An unexpected database error occurred.",
	}

	if len(msgs) > 0 {
		o := msgs[0]
		if o.Unique != "" {
			m.Unique = o.Unique
		}
		if o.NotFound != "" {
			m.NotFound = o.NotFound
		}
		if o.Busy != "" {
			m.Busy = o.Busy
		}
		if o.Default != "" {
			m.Default = o.Default
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return a.NotFound(errors.New(m.NotFound))
	}

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE, sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return a.Conflict(errors.New(m.Unique))
		case sqlite3.SQLITE_BUSY:
			return a.TooManyRequests(errors.New(m.Busy))
		}
	}

	return a.InternalServerError(errors.New(m.Default))
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
func (a *API) handleSessionValidation(c *gin.Context) APIResponse {
	t, exists := getSessionToken(c)
	if !exists {
		return a.Unauthorized(nil)
	}

	s, err := a.db.getSessionByToken(t)
	if err != nil {
		return a.Unauthorized(nil)
	}

	initSession(s, c)
	return a.OK(s)
}

// handleRegisterUser registers a new user.
func (a *API) handleRegisterUser(c *gin.Context) APIResponse {
	var in UserCredentialInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	if err := a.db.createUser(in.Name, in.Password); err != nil {
		return a.DBError(err, DBErrorMessages{Unique: "Username already exists."})
	}

	return a.Created(nil)
}

// handleLoginUser authenticates a user and creates a session.
func (a *API) handleLoginUser(c *gin.Context) APIResponse {
	var in UserCredentialInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	u, err := a.db.getUserByNameAndPassword(in.Name, in.Password)
	if err != nil {
		return a.Unauthorized(nil)
	}

	if err := a.createSession(u.Name, c); err != nil {
		return a.InternalServerError(nil)
	}

	return a.OK(nil)
}

// handleLogoutUser invalidates the active session and clears local session state.
func (a *API) handleLogoutUser(c *gin.Context) APIResponse {
	t, exists := getSessionToken(c)
	if !exists {
		return a.OK(nil)
	}

	if err := a.db.deleteSession(t); err != nil {
		return a.DBError(err)
	}

	clearSession(c)
	return a.OK(nil)
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

// toEntryParams converts a CreateEntryInput to EntryParams, applying necessary scaling.
func toEntryParams(in *CreateEntryInput) EntryParams {
	return EntryParams{
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
func (a *API) handleCreateEntry(c *gin.Context) APIResponse {
	var in CreateEntryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	t, exists := getSessionToken(c)
	if !exists {
		return a.Unauthorized(nil)
	}

	p := toEntryParams(&in)
	if p.Servings < 1 {
		return a.BadRequest(errors.New("Servings must be greater than 0"))
	}

	e, err := a.db.createEntryByToken(p, t)
	if err != nil {
		if strings.Contains(err.Error(), "NOT NULL constraint failed: entry.user_name") {
			return a.Unauthorized(err)
		}

		return a.DBError(err, DBErrorMessages{NotFound: "Invalid food ID."})
	}

	return a.OK(toEntryWithFoodResponse(e))
}

// handleListUserEntries returns entries with food details for a user and date.
func (a *API) handleListUserEntries(c *gin.Context) APIResponse {
	var in ListUserEntriesInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	res, err := a.db.listUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		return a.DBError(err)
	}

	r := make([]EntryWithFoodResponse, len(res))
	for i := range res {
		r[i] = toEntryWithFoodResponse(&res[i])
	}
	return a.OK(r)
}

// handleCreateFood creates a food record for the authenticated user.
func (a *API) handleCreateFood(c *gin.Context) APIResponse {
	var in CreateFoodInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	t, exists := getSessionToken(c)
	if !exists {
		return a.Unauthorized(nil)
	}

	p := toCreateFoodParams(&in)
	if p.ServingCount < 1 {
		return a.BadRequest(errors.New("Serving count must be greater than 0"))
	}

	if p.Calories < 0 || p.Carbs < 0 || p.Protein < 0 || p.Fat < 0 {
		return a.BadRequest(errors.New("Calories, carbs, protein, and fat must be non-negative"))
	}

	f, err := a.db.createFoodByToken(p, t)
	if err != nil {
		if strings.Contains(err.Error(), "NOT NULL constraint failed: food.user_name") {
			return a.Unauthorized(nil)
		}
		return a.DBError(err)
	}

	return a.OK(toFoodResponse(f))
}

// handleListFoods returns all foods.
func (a *API) handleListFoods(c *gin.Context) APIResponse {
	foods, err := a.db.listFoods()
	if err != nil {
		return a.DBError(err)
	}

	f := make([]FoodResponse, len(foods))
	for i := range foods {
		f[i] = toFoodResponse(&foods[i])
	}
	return a.OK(f)
}

// handleGetDiary returns all entries for a user and date, grouped by meal with totals.
func (a *API) handleGetDiary(c *gin.Context) APIResponse {
	var in ListUserEntriesInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	res, err := a.db.listUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		return a.DBError(err)
	}

	type Totals struct {
		Calories float64 `json:"calories"`
		Carbs    float64 `json:"carbs"`
		Protein  float64 `json:"protein"`
		Fat      float64 `json:"fat"`
	}

	type Diary struct {
		Totals  Totals                  `json:"totals"`
		Entries []EntryWithFoodResponse `json:"entries"`
	}

	type Response struct {
		Totals    Totals `json:"totals"`
		Breakfast Diary  `json:"breakfast"`
		Lunch     Diary  `json:"lunch"`
		Dinner    Diary  `json:"dinner"`
		Snacks    Diary  `json:"snacks"`
	}

	addTotals := func(totals *Totals, entry EntryWithFoodResponse) {
		totals.Calories += entry.Food.Calories * entry.Servings / entry.Food.ServingCount
		totals.Carbs += entry.Food.Carbs * entry.Servings / entry.Food.ServingCount
		totals.Protein += entry.Food.Protein * entry.Servings / entry.Food.ServingCount
		totals.Fat += entry.Food.Fat * entry.Servings / entry.Food.ServingCount
	}

	response := Response{
		Totals:    Totals{},
		Breakfast: Diary{Totals: Totals{}, Entries: []EntryWithFoodResponse{}},
		Lunch:     Diary{Totals: Totals{}, Entries: []EntryWithFoodResponse{}},
		Dinner:    Diary{Totals: Totals{}, Entries: []EntryWithFoodResponse{}},
		Snacks:    Diary{Totals: Totals{}, Entries: []EntryWithFoodResponse{}},
	}

	for i := range res {
		entry := toEntryWithFoodResponse(&res[i])
		addTotals(&response.Totals, entry)

		var meal *Diary
		switch strings.ToLower(entry.MealName) {
		case "breakfast":
			meal = &response.Breakfast
		case "lunch":
			meal = &response.Lunch
		case "dinner":
			meal = &response.Dinner
		case "snacks":
			meal = &response.Snacks
		}

		if meal != nil {
			meal.Entries = append(meal.Entries, entry)
			addTotals(&meal.Totals, entry)
		}
	}

	return a.OK(response)
}

// handleFoodSearch returns foods matching a search query.
func (a *API) handleFoodSearch(c *gin.Context) APIResponse {
	var in struct {
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		return a.BadRequest(err)
	}

	t, exists := getSessionToken(c)

	var foods []Food
	var err error
	if exists {
		foods, err = a.db.searchFoodsByNameSortedUserFromTokenFirst(in.Query, t)
	} else {
		foods, err = a.db.searchFoodsByName(in.Query)
	}

	if err != nil {
		return a.DBError(err)
	}

	f := make([]FoodResponse, len(foods))
	for i := range foods {
		f[i] = toFoodResponse(&foods[i])
	}

	return a.OK(f)
}

// handleEditEntry edits an existing entry for the authenticated user.
func (api *API) handleEditEntry(c *gin.Context) APIResponse {
	t, exsits := getSessionToken(c)
	if !exsits {
		return api.Unauthorized(nil)
	}

	var in struct {
		ID int `json:"id" binding:"required"`
		CreateEntryInput
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	p := toEntryParams(&in.CreateEntryInput)
	if p.Servings < 1 {
		return api.BadRequest(errors.New("Servings must be greater than 0"))
	}

	e, err := api.db.editEntryAuthByToken(in.ID, p, t)
	if err != nil {
		return api.DBError(err)
	}

	r := toEntryWithFoodResponse(e)
	return api.OK(r)
}

// handleDeleteEntry deletes an existing entry for the authenticated user.
func (api *API) handleDeleteEntry(c *gin.Context) APIResponse {
	t, exsits := getSessionToken(c)
	if !exsits {
		return api.Unauthorized(nil)
	}

	var in struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	if err := api.db.deleteEntryAuthByToken(in.ID, t); err != nil {
		return api.DBError(err)
	}

	return api.NoContent()
}
