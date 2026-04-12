package api

import (
	"macro/db"
	"macro/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// initSession initializes session state with the provided Session struct.
func initSession(s *db.Session, c *gin.Context) {
	defer util.Trace("initSession(s *Session, c *gin.Context)", "session", s)()

	session := sessions.Default(c)
	session.Set("token", s.Token)
	session.Set("username", s.UserName)
	session.Save()
}

// clearSession clears and persists session state.
func clearSession(c *gin.Context) {
	defer util.Trace("clearSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	session.Set("token", nil)
	session.Set("username", nil)
	session.Save()
}

// getSessionToken returns the session token if present.
func getSessionToken(c *gin.Context) (string, bool) {
	defer util.Trace("getSessionToken(c *gin.Context)")()

	session := sessions.Default(c)
	token := session.Get("token")
	if token == nil {
		return "", false
	}

	return token.(string), true
}

// createSession creates a DB session and persists its token in the cookie.
func (api *API) createSession(name string, c *gin.Context) error {
	defer util.Trace("createSession(name string, c *gin.Context)", "name", name)()

	s, err := api.db.CreateSession(name)
	if err != nil {
		c.Set("error", err.Error())
		return err
	}

	initSession(s, c)

	return nil
}

// handleSessionValidation verifies the current session token.
func (api *API) handleSessionValidation(c *gin.Context) APIResponse {
	t, exists := getSessionToken(c)
	if !exists {
		return api.Unauthorized(nil)
	}

	s, err := api.db.GetSessionByToken(t)
	if err != nil {
		return api.Unauthorized(nil)
	}

	initSession(s, c)
	return api.OK(s)
}

// handleRegisterUser registers a new user.
func (api *API) handleRegisterUser(c *gin.Context) APIResponse {
	var in UserCredentialInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	if err := api.db.CreateUser(in.Name, in.Password); err != nil {
		return api.DBError(err, DBErrorMessages{Unique: "Username already exists."})
	}

	return api.Created(nil)
}

// handleLoginUser authenticates a user and creates a session.
func (api *API) handleLoginUser(c *gin.Context) APIResponse {
	var in UserCredentialInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	u, err := api.db.GetUserByNameAndPassword(in.Name, in.Password)
	if err != nil {
		return api.Unauthorized(nil)
	}

	if err := api.createSession(u.Name, c); err != nil {
		return api.InternalServerError(nil)
	}

	return api.OK(nil)
}

// handleLogoutUser invalidates the active session and clears local session state.
func (api *API) handleLogoutUser(c *gin.Context) APIResponse {
	t, exists := getSessionToken(c)
	if !exists {
		return api.OK(nil)
	}

	if err := api.db.DeleteSession(t); err != nil {
		return api.DBError(err)
	}

	clearSession(c)
	return api.OK(nil)
}
