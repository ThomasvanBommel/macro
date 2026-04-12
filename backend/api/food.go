package api

import (
	"errors"
	"macro/db"
	"strings"

	"github.com/gin-gonic/gin"
)

// toFoodResponse converts a Food DB record to a FoodResponse for API responses.
func toFoodResponse(f *db.Food) FoodResponse {
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

// toCreateFoodParams converts a CreateFoodInput to CreateFoodParams, applying necessary scaling.
func toCreateFoodParams(in *CreateFoodInput) db.FoodParams {
	return db.FoodParams{
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

// handleCreateFood creates a food record for the authenticated user.
func (api *API) handleCreateFood(c *gin.Context) APIResponse {
	var in CreateFoodInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	t, exists := getSessionToken(c)
	if !exists {
		return api.Unauthorized(nil)
	}

	p := toCreateFoodParams(&in)
	if p.ServingCount < 1 {
		return api.BadRequest(errors.New("Serving count must be greater than 0"))
	}

	if p.Calories < 0 || p.Carbs < 0 || p.Protein < 0 || p.Fat < 0 {
		return api.BadRequest(errors.New("Calories, carbs, protein, and fat must be non-negative"))
	}

	f, err := api.db.CreateFoodByToken(p, t)
	if err != nil {
		if strings.Contains(err.Error(), "NOT NULL constraint failed: food.user_name") {
			return api.Unauthorized(nil)
		}
		return api.DBError(err)
	}

	return api.OK(toFoodResponse(f))
}

// handleListFoods returns all foods.
func (api *API) handleListFoods(c *gin.Context) APIResponse {
	foods, err := api.db.ListFoods()
	if err != nil {
		return api.DBError(err)
	}

	f := make([]FoodResponse, len(foods))
	for i := range foods {
		f[i] = toFoodResponse(&foods[i])
	}
	return api.OK(f)
}

// handleFoodSearch returns foods matching a search query.
func (api *API) handleFoodSearch(c *gin.Context) APIResponse {
	var in struct {
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	t, exists := getSessionToken(c)

	var foods []db.Food
	var err error
	if exists {
		foods, err = api.db.SearchFoodsByNameSortedUserFromTokenFirst(in.Query, t)
	} else {
		foods, err = api.db.SearchFoodsByName(in.Query)
	}

	if err != nil {
		return api.DBError(err)
	}

	f := make([]FoodResponse, len(foods))
	for i := range foods {
		f[i] = toFoodResponse(&foods[i])
	}

	return api.OK(f)
}
