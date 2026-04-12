package api

import (
	"errors"
	"macro/db"
	"strings"

	"github.com/gin-gonic/gin"
)

// toEntryParams converts a CreateEntryInput to EntryParams, applying necessary scaling.
func toEntryParams(in *CreateEntryInput) db.EntryParams {
	return db.EntryParams{
		FoodId:   *in.FoodId,
		MealName: in.MealName,
		Date:     in.Date,
		Servings: scaleForStorage(in.Servings),
	}
}

// toEntryWithFoodResponse converts an EntryWithFood DB record to an EntryWithFoodResponse
func toEntryWithFoodResponse(e *db.EntryWithFood) EntryWithFoodResponse {
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

// handleCreateEntry creates a food entry for the authenticated user.
func (api *API) handleCreateEntry(c *gin.Context) APIResponse {
	var in CreateEntryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	t, exists := getSessionToken(c)
	if !exists {
		return api.Unauthorized(nil)
	}

	p := toEntryParams(&in)
	if p.Servings < 1 {
		return api.BadRequest(errors.New("Servings must be greater than 0"))
	}

	e, err := api.db.CreateEntryByToken(p, t)
	if err != nil {
		if strings.Contains(err.Error(), "NOT NULL constraint failed: entry.user_name") {
			return api.Unauthorized(err)
		}

		return api.DBError(err, DBErrorMessages{NotFound: "Invalid food ID."})
	}

	return api.OK(toEntryWithFoodResponse(e))
}

// handleListUserEntries returns entries with food details for a user and date.
func (api *API) handleListUserEntries(c *gin.Context) APIResponse {
	var in ListUserEntriesInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	res, err := api.db.ListUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		return api.DBError(err)
	}

	r := make([]EntryWithFoodResponse, len(res))
	for i := range res {
		r[i] = toEntryWithFoodResponse(&res[i])
	}
	return api.OK(r)
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

	e, err := api.db.EditEntryAuthByToken(in.ID, p, t)
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

	if err := api.db.DeleteEntryAuthByToken(in.ID, t); err != nil {
		return api.DBError(err)
	}

	return api.NoContent()
}
