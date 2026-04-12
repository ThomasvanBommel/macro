package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// handleGetDiary returns all entries for a user and date, grouped by meal with totals.
func (api *API) handleGetDiary(c *gin.Context) APIResponse {
	var in ListUserEntriesInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return api.BadRequest(err)
	}

	res, err := api.db.ListUserEntriesWithFoodByNameAndDate(in.Name, in.Date)
	if err != nil {
		return api.DBError(err)
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

	return api.OK(response)
}
