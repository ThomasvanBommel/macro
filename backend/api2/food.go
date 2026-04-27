package api2

import "github.com/gin-gonic/gin"

type FoodInput struct {
	Name         string  `json:"name" binding:"required"`
	Brand        string  `json:"brand"`
	ServingSize  string  `json:"serving_size" binding:"required"`
	Calories     float64 `json:"calories" binding:"required, gte=0"`
	Carbs        float64 `json:"carbs" binding:"required, gte=0"`
	Protein      float64 `json:"protein" binding:"required, gte=0"`
	Fat          float64 `json:"fat" binding:"required, gte=0"`
	ServingCount float64 `json:"serving_count" binding:"required, gte=0.01"`
}

func forStorage(in float64) int {
	return int(in * 100)
}

func (db *DB) NewFood(c *gin.Context) {
	var input FoodInput
	if err := BindJSON(c, &input); err != nil {
		ErrorResponse(c, err)
		return
	}

	if input.Name == "" {
		Unprocessable(c, "Name cannot be empty")
		return
	}

	// calories := forStorage(input.Calories)
	// carbs := forStorage(input.Carbs)
	// protein := forStorage(input.Protein)
	// fat := forStorage(input.Fat)
	// servingCount := forStorage(input.ServingCount)

	// err, id, username, created
}
